package jwt

import (
	"context"
	"errors"
	"testing"
	"time"

	authdomain "contai/internal/auth/domain"
	userdomain "contai/internal/users/domain"
)

func TestService_IssueAndValidateAccessToken(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	service := Service{secret: []byte("secret"), ttl: 30 * time.Minute, now: func() time.Time { return now }}
	user := mustJWTUser(t)

	token, issuedClaims, err := service.IssueAccessToken(context.Background(), user)
	if err != nil {
		t.Fatalf("expected token, got %v", err)
	}

	claims, err := service.ValidateAccessToken(context.Background(), token)
	if err != nil {
		t.Fatalf("expected valid token, got %v", err)
	}
	if claims.UserID != user.ID || claims.Type != authdomain.AuthTokenTypeAccess {
		t.Fatalf("expected access claims for user %s, got %#v", user.ID, claims)
	}
	if !claims.ExpiresAt.Equal(issuedClaims.ExpiresAt) {
		t.Fatalf("expected expiration %v, got %v", issuedClaims.ExpiresAt, claims.ExpiresAt)
	}
}

func TestService_ValidateAccessTokenRejectsInvalidCases(t *testing.T) {
	now := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	user := mustJWTUser(t)

	t.Run("wrong secret", func(t *testing.T) {
		issuer := Service{secret: []byte("secret"), ttl: 30 * time.Minute, now: func() time.Time { return now }}
		validator := Service{secret: []byte("other-secret"), ttl: 30 * time.Minute, now: func() time.Time { return now }}
		token, _, err := issuer.IssueAccessToken(context.Background(), user)
		if err != nil {
			t.Fatalf("expected token, got %v", err)
		}

		_, err = validator.ValidateAccessToken(context.Background(), token)

		if !errors.Is(err, authdomain.ErrInvalidToken) {
			t.Fatalf("expected invalid token, got %v", err)
		}
	})

	t.Run("expired", func(t *testing.T) {
		issuer := Service{secret: []byte("secret"), ttl: time.Minute, now: func() time.Time { return now }}
		validator := Service{secret: []byte("secret"), ttl: time.Minute, now: func() time.Time { return now.Add(2 * time.Minute) }}
		token, _, err := issuer.IssueAccessToken(context.Background(), user)
		if err != nil {
			t.Fatalf("expected token, got %v", err)
		}

		_, err = validator.ValidateAccessToken(context.Background(), token)

		if !errors.Is(err, authdomain.ErrExpiredToken) {
			t.Fatalf("expected expired token, got %v", err)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		service := Service{secret: []byte("secret"), ttl: 30 * time.Minute, now: func() time.Time { return now }}
		token, err := service.sign(authdomain.AuthClaims{
			Type:      authdomain.AuthTokenType("refresh"),
			UserID:    user.ID,
			IssuedAt:  now,
			ExpiresAt: now.Add(time.Hour),
		})
		if err != nil {
			t.Fatalf("expected token, got %v", err)
		}

		_, err = service.ValidateAccessToken(context.Background(), token)

		if !errors.Is(err, authdomain.ErrInvalidToken) {
			t.Fatalf("expected invalid token, got %v", err)
		}
	})
}

func mustJWTUser(t *testing.T) userdomain.User {
	t.Helper()

	user, err := userdomain.NewUser("user-id", "John Doe", "john@example.com", "hash")
	if err != nil {
		t.Fatalf("expected user, got %v", err)
	}

	return user
}
