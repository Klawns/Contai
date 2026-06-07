package ports

import (
	"context"
	"time"

	dashboarddomain "contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type DashboardRepository interface {
	FindActiveAccountBalances(ctx context.Context, userID userdomain.UserID) ([]AccountBalanceDTO, error)
	FindCreditCards(ctx context.Context, userID userdomain.UserID, now time.Time) ([]CreditCardDashboardDTO, error)
	SumIncome(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) (financedomain.Money, error)
	SumExpense(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) (financedomain.Money, error)
	FindMonthlyIncomeExpense(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]MonthlyIncomeExpenseDTO, error)
	FindMonthlyBalances(ctx context.Context, userID userdomain.UserID, monthEnds []time.Time) ([]MonthlyBalanceDTO, error)
	FindTransactionsByPeriod(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]transactiondomain.Transaction, error)
	FindExpensesByCategory(ctx context.Context, userID userdomain.UserID, period dashboarddomain.Period) ([]CategoryExpenseDTO, error)
	FindRecentTransactions(ctx context.Context, userID userdomain.UserID, limit int) ([]TransactionDTO, error)
}

type MonthlyIncomeExpenseDTO struct {
	MonthStartAt time.Time
	Income       financedomain.Money
	Expense      financedomain.Money
}

type MonthlyBalanceDTO struct {
	MonthEndAt time.Time
	Balance    financedomain.Money
}
