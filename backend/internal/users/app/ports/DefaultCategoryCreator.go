package ports

import (
	databaseports "contai/internal/database/ports"
	"contai/internal/users/domain"
	"context"
)

type DefaultCategoryCreator interface {
	WithTx(tx databaseports.TxHandle) DefaultCategoryCreator
	EnsureDefaultCategories(ctx context.Context, userID domain.UserID) error
}
