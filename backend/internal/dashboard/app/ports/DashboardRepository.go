package ports

import (
	"context"

	dashboarddomain "contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type DashboardRepository interface {
	FindActiveAccountBalances(ctx context.Context, userID userdomain.UserID) ([]AccountBalanceDTO, error)
	SumIncome(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) (financedomain.Money, error)
	SumExpense(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) (financedomain.Money, error)
	FindTransactionsByPeriod(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]transactiondomain.Transaction, error)
	FindExpensesByCategory(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]CategoryExpenseDTO, error)
	FindRecentTransactions(ctx context.Context, userID userdomain.UserID, limit int) ([]TransactionDTO, error)
}
