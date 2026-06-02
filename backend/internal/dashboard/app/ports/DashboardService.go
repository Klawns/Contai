package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	dashboarddomain "contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type GetMonthlyDashboardInput struct {
	UserID userdomain.UserID
	Period dashboarddomain.Period
}

type MonthlyDashboardDTO struct {
	UserID             userdomain.UserID
	Period             dashboarddomain.Period
	TotalBalance       financedomain.Money
	MonthlyIncome      financedomain.Money
	MonthlyExpense     financedomain.Money
	MonthlyTransferIn  financedomain.Money
	MonthlyTransferOut financedomain.Money
	MonthlyNetBalance  financedomain.Money
	AccountBalances    []AccountBalanceDTO
	ExpensesByCategory []CategoryExpenseDTO
	RecentTransactions []TransactionDTO
}

type AccountBalanceDTO struct {
	AccountID               accountdomain.AccountID
	Name                    string
	Type                    accountdomain.AccountType
	Balance                 financedomain.Money
	BankIconID              string
	IncludeInDashboardTotal bool
}

type CategoryExpenseDTO struct {
	CategoryID categorydomain.CategoryID
	Name       string
	Color      string
	Icon       string
	Total      financedomain.Money
}

type TransactionDTO struct {
	ID                   transactiondomain.TransactionID
	UserID               userdomain.UserID
	Type                 transactiondomain.TransactionType
	Description          string
	Amount               financedomain.Money
	OccurredAt           time.Time
	AccountID            *accountdomain.AccountID
	SourceAccountID      *accountdomain.AccountID
	DestinationAccountID *accountdomain.AccountID
	CategoryID           *categorydomain.CategoryID
	Status               transactiondomain.TransactionStatus
	Note                 string
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type DashboardService interface {
	GetMonthlyDashboard(ctx context.Context, input GetMonthlyDashboardInput) (MonthlyDashboardDTO, error)
}
