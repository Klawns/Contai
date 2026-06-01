package services

import (
	"context"
	"strings"

	"contai/internal/auth/app/contracts"
	"contai/internal/auth/app/ports"
	authdomain "contai/internal/auth/domain"
	userports "contai/internal/users/app/ports"
	userdomain "contai/internal/users/domain"
)

var _ ports.AuthService = AuthService{}

type AuthService struct {
	userRepository userports.UserRepository
	passwordHasher userports.PasswordHasher
	jwtService     ports.JWTService
}

func NewAuthService(
	userRepository userports.UserRepository,
	passwordHasher userports.PasswordHasher,
	jwtService ports.JWTService,
) AuthService {
	return AuthService{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		jwtService:     jwtService,
	}
}

func (service AuthService) Login(ctx context.Context, input contracts.LoginInput) (authdomain.AuthenticatedUser, contracts.AuthTokens, error) {
	email := normalizeEmail(input.Email)
	user, err := service.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return authdomain.AuthenticatedUser{}, contracts.AuthTokens{}, err
	}
	if user == nil {
		return authdomain.AuthenticatedUser{}, contracts.AuthTokens{}, authdomain.ErrInvalidCredentials
	}

	if _, err := user.CanAuthenticate(); err != nil {
		return authdomain.AuthenticatedUser{}, contracts.AuthTokens{}, err
	}

	if err := service.passwordHasher.ComparePassword(ctx, user.PasswordHash, input.PlainPassword); err != nil {
		return authdomain.AuthenticatedUser{}, contracts.AuthTokens{}, authdomain.ErrInvalidCredentials
	}

	accessToken, accessClaims, err := service.jwtService.IssueAccessToken(ctx, *user)
	if err != nil {
		return authdomain.AuthenticatedUser{}, contracts.AuthTokens{}, err
	}

	return toAuthenticatedUser(*user), contracts.AuthTokens{
		AccessToken:  accessToken,
		AccessClaims: accessClaims,
	}, nil
}

func (service AuthService) Logout(ctx context.Context) error {
	return nil
}

func (service AuthService) ValidateAccessToken(ctx context.Context, accessToken string) (authdomain.AuthenticatedUser, error) {
	claims, err := service.jwtService.ValidateAccessToken(ctx, accessToken)
	if err != nil {
		return authdomain.AuthenticatedUser{}, err
	}
	if claims.Type != authdomain.AuthTokenTypeAccess {
		return authdomain.AuthenticatedUser{}, authdomain.ErrInvalidToken
	}

	user, err := service.userRepository.FindUserById(ctx, claims.UserID)
	if err != nil {
		return authdomain.AuthenticatedUser{}, err
	}
	if user == nil {
		return authdomain.AuthenticatedUser{}, authdomain.ErrInvalidToken
	}
	if _, err := user.CanAuthenticate(); err != nil {
		return authdomain.AuthenticatedUser{}, err
	}

	return toAuthenticatedUser(*user), nil
}

func toAuthenticatedUser(user userdomain.User) authdomain.AuthenticatedUser {
	return authdomain.AuthenticatedUser{
		UserID: user.ID,
		Email:  user.Email,
		Status: user.Status,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
