package ports

import (
	authdomain "contai/internal/auth/domain"
	userdomain "contai/internal/users/domain"
	"context"
)

type JWTService interface {
	IssueAccessToken(ctx context.Context, user userdomain.User) (string, authdomain.AuthClaims, error)
	ValidateAccessToken(ctx context.Context, token string) (authdomain.AuthClaims, error)
}
