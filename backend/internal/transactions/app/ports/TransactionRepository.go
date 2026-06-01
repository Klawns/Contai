package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	"contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type ListTransactionsInput struct {
	UserID     userdomain.UserID
	StartAt    *time.Time
	EndAt      *time.Time
	AccountID  *accountdomain.AccountID
	CategoryID *categorydomain.CategoryID
	Type       *domain.TransactionType
	Limit      int
	Offset     int
}

type TransactionRepository interface {
	WithTx(tx databaseports.TxHandle) TransactionRepository
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
	FindTransactionByID(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID) (*domain.Transaction, error)
	FindTransactionByIDForUpdate(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID) (*domain.Transaction, error)
	FindTransactionsByUserID(ctx context.Context, input ListTransactionsInput) ([]domain.Transaction, error)
}
