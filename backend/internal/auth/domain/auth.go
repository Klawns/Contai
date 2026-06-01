package domain

import (
	"errors"
	"time"

	userdomain "contai/internal/users/domain"
)

type AuthTokenType string

const (
	AuthTokenTypeAccess AuthTokenType = "access"
)

type AuthClaims struct {
	Type      AuthTokenType
	UserID    userdomain.UserID
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type AuthenticatedUser struct {
	UserID userdomain.UserID
	Email  string
	Status userdomain.UserStatus
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("expired token")
)
