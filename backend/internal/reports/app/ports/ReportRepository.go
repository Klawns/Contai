package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
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

type ListFinancialMovementsInput struct {
	UserID           userdomain.UserID
	StartAt          time.Time
	EndAt            time.Time
	MovementType     MovementType
	CategoryID       *categorydomain.CategoryID
	AccountID        *accountdomain.AccountID
	SettlementStatus SettlementStatusFilter
}

type ReportRepository interface {
	ListFinancialMovements(ctx context.Context, input ListFinancialMovementsInput) ([]FinancialMovementDTO, error)
}
