package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type ListReportTransactionsInput struct {
	UserID    userdomain.UserID
	StartAt   time.Time
	EndAt     time.Time
	Type      *transactiondomain.TransactionType
	AccountID *accountdomain.AccountID
}

type ReportRepository interface {
	FindAccountByID(ctx context.Context, userID userdomain.UserID, accountID accountdomain.AccountID) (*AccountReportRow, error)
	ListAccounts(ctx context.Context, userID userdomain.UserID) ([]AccountReportRow, error)
	ListTransactions(ctx context.Context, input ListReportTransactionsInput) ([]ReportTransactionRow, error)
}
