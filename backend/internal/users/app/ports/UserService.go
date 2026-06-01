package ports

import (
	"contai/internal/users/domain"
	"context"
	"time"
)

type UserDTO struct {
	ID        domain.UserID
	Name      string
	Email     string
	Status    domain.UserStatus
	CreatedAt time.Time
}

type CreateUserInput struct {
	Name          string
	Email         string
	PlainPassword string
}

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (UserDTO, error)
	GetUserByID(ctx context.Context, userID domain.UserID) (UserDTO, error)
}
