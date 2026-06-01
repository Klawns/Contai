package contracts

import "contai/internal/auth/domain"

type LoginInput struct {
	Email         string
	PlainPassword string
}

type AuthTokens struct {
	AccessToken  string
	AccessClaims domain.AuthClaims
}
