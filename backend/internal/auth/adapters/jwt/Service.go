package jwt

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	authdomain "contai/internal/auth/domain"
	userdomain "contai/internal/users/domain"
)

type Service struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time
}

func NewService(secret string, ttl time.Duration) Service {
	if ttl <= 0 {
		ttl = 30 * time.Minute
	}

	return Service{
		secret: []byte(secret),
		ttl:    ttl,
		now:    time.Now,
	}
}

func (service Service) IssueAccessToken(ctx context.Context, user userdomain.User) (string, authdomain.AuthClaims, error) {
	issuedAt := service.now().UTC()
	claims := authdomain.AuthClaims{
		Type:      authdomain.AuthTokenTypeAccess,
		UserID:    user.ID,
		IssuedAt:  issuedAt,
		ExpiresAt: issuedAt.Add(service.ttl),
	}

	token, err := service.sign(claims)
	if err != nil {
		return "", authdomain.AuthClaims{}, err
	}

	return token, claims, nil
}

func (service Service) ValidateAccessToken(ctx context.Context, token string) (authdomain.AuthClaims, error) {
	claims, err := service.parse(token)
	if err != nil {
		return authdomain.AuthClaims{}, err
	}
	if claims.Type != authdomain.AuthTokenTypeAccess {
		return authdomain.AuthClaims{}, authdomain.ErrInvalidToken
	}
	if !service.now().UTC().Before(claims.ExpiresAt) {
		return authdomain.AuthClaims{}, authdomain.ErrExpiredToken
	}

	return claims, nil
}

func (service Service) sign(claims authdomain.AuthClaims) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	payload := tokenPayload{
		Type: string(claims.Type),
		Sub:  string(claims.UserID),
		Iat:  claims.IssuedAt.Unix(),
		Exp:  claims.ExpiresAt.Unix(),
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	unsigned := encodeSegment(headerJSON) + "." + encodeSegment(payloadJSON)
	return unsigned + "." + service.signature(unsigned), nil
}

func (service Service) parse(token string) (authdomain.AuthClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return authdomain.AuthClaims{}, authdomain.ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(parts[2]), []byte(service.signature(unsigned))) {
		return authdomain.AuthClaims{}, authdomain.ErrInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return authdomain.AuthClaims{}, authdomain.ErrInvalidToken
	}

	var payload tokenPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return authdomain.AuthClaims{}, authdomain.ErrInvalidToken
	}
	if payload.Sub == "" || payload.Type == "" || payload.Exp == 0 {
		return authdomain.AuthClaims{}, authdomain.ErrInvalidToken
	}

	return authdomain.AuthClaims{
		Type:      authdomain.AuthTokenType(payload.Type),
		UserID:    userdomain.UserID(payload.Sub),
		IssuedAt:  time.Unix(payload.Iat, 0).UTC(),
		ExpiresAt: time.Unix(payload.Exp, 0).UTC(),
	}, nil
}

func (service Service) signature(unsigned string) string {
	mac := hmac.New(sha256.New, service.secret)
	mac.Write([]byte(unsigned))
	return encodeSegment(mac.Sum(nil))
}

func encodeSegment(value []byte) string {
	return base64.RawURLEncoding.EncodeToString(value)
}

type tokenPayload struct {
	Type string `json:"typ"`
	Sub  string `json:"sub"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
}
