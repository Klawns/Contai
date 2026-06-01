package ports

import (
	"context"

	"contai/internal/account/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"
)

type AccountRepository interface {
	WithTx(tx databaseports.TxHandle) AccountRepository
	CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error)
	FindAccountByID(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID) (*domain.Account, error)
	FindAccountByIDForUpdate(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID) (*domain.Account, error)
	FindAccountsByUserID(ctx context.Context, input ListAccountsInput) ([]domain.Account, error)
	SumActiveAccountBalances(ctx context.Context, userID userdomain.UserID) (int64, error)
}
