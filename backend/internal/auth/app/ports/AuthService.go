package ports

import (
	"contai/internal/auth/app/contracts"
	"contai/internal/auth/domain"
	"context"
)

type AuthService interface {
	Login(ctx context.Context, input contracts.LoginInput) (domain.AuthenticatedUser, contracts.AuthTokens, error)
	Logout(ctx context.Context) error
	ValidateAccessToken(ctx context.Context, accessToken string) (domain.AuthenticatedUser, error)
}
