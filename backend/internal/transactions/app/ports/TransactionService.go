package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	"contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type TransactionDTO struct {
	ID                   domain.TransactionID
	UserID               userdomain.UserID
	Type                 domain.TransactionType
	Description          string
	Amount               financedomain.Money
	OccurredAt           time.Time
	AccountID            *accountdomain.AccountID
	SourceAccountID      *accountdomain.AccountID
	DestinationAccountID *accountdomain.AccountID
	CategoryID           *categorydomain.CategoryID
	Status               domain.TransactionStatus
	OriginType           domain.TransactionOriginType
	OriginID             *string
	Note                 string
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CreateIncomeInput struct {
	UserID      userdomain.UserID
	Description string
	Amount      financedomain.Money
	OccurredAt  time.Time
	AccountID   accountdomain.AccountID
	CategoryID  categorydomain.CategoryID
	OriginType  domain.TransactionOriginType
	OriginID    string
	Note        string
}

type CreateExpenseInput = CreateIncomeInput

type CreateTransferInput struct {
	UserID               userdomain.UserID
	Description          string
	Amount               financedomain.Money
	OccurredAt           time.Time
	SourceAccountID      accountdomain.AccountID
	DestinationAccountID accountdomain.AccountID
	Note                 string
}

type UpdateTransactionInput struct {
	UserID               userdomain.UserID
	TransactionID        domain.TransactionID
	Description          string
	Amount               financedomain.Money
	OccurredAt           time.Time
	AccountID            accountdomain.AccountID
	SourceAccountID      accountdomain.AccountID
	DestinationAccountID accountdomain.AccountID
	CategoryID           categorydomain.CategoryID
	Note                 string
}

type DeleteTransactionInput struct {
	UserID        userdomain.UserID
	TransactionID domain.TransactionID
}

type TransactionService interface {
	CreateIncome(ctx context.Context, input CreateIncomeInput) (TransactionDTO, error)
	CreateExpense(ctx context.Context, input CreateExpenseInput) (TransactionDTO, error)
	CreateTransfer(ctx context.Context, input CreateTransferInput) (TransactionDTO, error)
	ListTransactions(ctx context.Context, input ListTransactionsInput) ([]TransactionDTO, error)
	UpdateTransaction(ctx context.Context, input UpdateTransactionInput) (TransactionDTO, error)
	DeleteTransaction(ctx context.Context, input DeleteTransactionInput) error
}
