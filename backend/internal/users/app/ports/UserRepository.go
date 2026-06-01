package ports

import (
	"context"

	databaseports "contai/internal/database/ports"
	"contai/internal/users/domain"
)

type UserRepository interface {
	WithTx(tx databaseports.TxHandle) UserRepository
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindUserById(ctx context.Context, userID domain.UserID) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}
