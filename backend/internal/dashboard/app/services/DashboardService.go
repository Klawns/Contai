package services

import (
	"context"
	"time"

	"contai/internal/dashboard/app/ports"
	"contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
)

const recentTransactionsLimit = 5

var _ ports.DashboardService = DashboardService{}

type DashboardService struct {
	repository ports.DashboardRepository
}

func NewDashboardService(repository ports.DashboardRepository) DashboardService {
	return DashboardService{repository: repository}
}

func (service DashboardService) GetMonthlyDashboard(ctx context.Context, input ports.GetMonthlyDashboardInput) (ports.MonthlyDashboardDTO, error) {
	if input.UserID == "" {
		return ports.MonthlyDashboardDTO{}, domain.ErrDashboardUserIDRequired
	}
	if err := input.Period.Validate(); err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}

	accountBalances, err := service.repository.FindActiveAccountBalances(ctx, input.UserID)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	totalBalance := financedomain.NewMoney(0)
	for _, accountBalance := range accountBalances {
		if accountBalance.IncludeInDashboardTotal {
			totalBalance = totalBalance.Add(accountBalance.Balance)
		}
	}

	transactions, err := service.repository.FindTransactionsByPeriod(ctx, input.UserID, input.Period)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	transactionTotals := transactiondomain.CalculateTransactionTotals(transactions)
	expensesByCategory, err := service.repository.FindExpensesByCategory(ctx, input.UserID, input.Period)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	recentTransactions, err := service.repository.FindRecentTransactions(ctx, input.UserID, recentTransactionsLimit)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	creditCards, err := service.repository.FindCreditCards(ctx, input.UserID, time.Now())
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}

	return ports.MonthlyDashboardDTO{
		UserID:             input.UserID,
		Period:             input.Period,
		TotalBalance:       totalBalance,
		MonthlyIncome:      transactionTotals.IncomeTotal,
		MonthlyExpense:     transactionTotals.ExpenseTotal,
		MonthlyTransferIn:  transactionTotals.TransferInTotal,
		MonthlyTransferOut: transactionTotals.TransferOutTotal,
		MonthlyNetBalance:  transactionTotals.IncomeTotal.Sub(transactionTotals.ExpenseTotal),
		AccountBalances:    nonNilAccountBalances(accountBalances),
		CreditCards:        nonNilCreditCards(creditCards),
		ExpensesByCategory: nonNilExpensesByCategory(expensesByCategory),
		RecentTransactions: nonNilRecentTransactions(recentTransactions),
	}, nil
}

func (service DashboardService) GetMonthlySeries(ctx context.Context, input ports.GetMonthlySeriesInput) (ports.MonthlySeriesDTO, error) {
	if input.UserID == "" {
		return ports.MonthlySeriesDTO{}, domain.ErrDashboardUserIDRequired
	}
	if err := input.Period.Validate(); err != nil {
		return ports.MonthlySeriesDTO{}, err
	}

	months := monthlyPeriods(input.Period)
	monthEnds := make([]time.Time, 0, len(months))
	for _, month := range months {
		monthEnds = append(monthEnds, month.MonthEndAt)
	}

	incomeExpenses, err := service.repository.FindMonthlyIncomeExpense(ctx, input.UserID, input.Period)
	if err != nil {
		return ports.MonthlySeriesDTO{}, err
	}
	balances, err := service.repository.FindMonthlyBalances(ctx, input.UserID, monthEnds)
	if err != nil {
		return ports.MonthlySeriesDTO{}, err
	}

	incomeExpenseByMonth := map[string]ports.MonthlyIncomeExpenseDTO{}
	for _, value := range incomeExpenses {
		incomeExpenseByMonth[monthKey(value.MonthStartAt)] = value
	}
	balanceByMonthEnd := map[string]financedomain.Money{}
	for _, value := range balances {
		balanceByMonthEnd[monthKey(value.MonthEndAt)] = value.Balance
	}

	points := make([]ports.MonthlySeriesPointDTO, 0, len(months))
	for _, month := range months {
		incomeExpense := incomeExpenseByMonth[monthKey(month.MonthStartAt)]
		points = append(points, ports.MonthlySeriesPointDTO{
			MonthStartAt: month.MonthStartAt,
			MonthEndAt:   month.MonthEndAt,
			Income:       incomeExpense.Income,
			Expense:      incomeExpense.Expense,
			Balance:      balanceByMonthEnd[monthKey(month.MonthEndAt)],
		})
	}

	return ports.MonthlySeriesDTO{
		UserID: input.UserID,
		Period: input.Period,
		Points: points,
	}, nil
}

func monthlyPeriods(period domain.Period) []ports.MonthlySeriesPointDTO {
	location := period.StartAt.Location()
	cursor := time.Date(period.StartAt.In(location).Year(), period.StartAt.In(location).Month(), 1, 0, 0, 0, 0, location)
	last := period.EndAt.In(location)
	months := []ports.MonthlySeriesPointDTO{}
	for !cursor.After(last) {
		nextMonth := cursor.AddDate(0, 1, 0)
		monthEnd := nextMonth.Add(-time.Second)
		months = append(months, ports.MonthlySeriesPointDTO{
			MonthStartAt: cursor,
			MonthEndAt:   monthEnd,
		})
		cursor = nextMonth
	}
	return months
}

func monthKey(value time.Time) string {
	return value.Format("2006-01")
}

func nonNilAccountBalances(values []ports.AccountBalanceDTO) []ports.AccountBalanceDTO {
	if values == nil {
		return []ports.AccountBalanceDTO{}
	}
	return values
}

func nonNilExpensesByCategory(values []ports.CategoryExpenseDTO) []ports.CategoryExpenseDTO {
	if values == nil {
		return []ports.CategoryExpenseDTO{}
	}
	return values
}

func nonNilCreditCards(values []ports.CreditCardDashboardDTO) []ports.CreditCardDashboardDTO {
	if values == nil {
		return []ports.CreditCardDashboardDTO{}
	}
	return values
}

func nonNilRecentTransactions(values []ports.TransactionDTO) []ports.TransactionDTO {
	if values == nil {
		return []ports.TransactionDTO{}
	}
	return values
}
