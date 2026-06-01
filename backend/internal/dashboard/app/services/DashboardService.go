package services

import (
	"context"

	"contai/internal/dashboard/app/ports"
	"contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
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
		totalBalance = totalBalance.Add(accountBalance.Balance)
	}

	monthlyIncome, err := service.repository.SumIncome(ctx, input.UserID, input.Period)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	monthlyExpense, err := service.repository.SumExpense(ctx, input.UserID, input.Period)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	expensesByCategory, err := service.repository.FindExpensesByCategory(ctx, input.UserID, input.Period)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}
	recentTransactions, err := service.repository.FindRecentTransactions(ctx, input.UserID, recentTransactionsLimit)
	if err != nil {
		return ports.MonthlyDashboardDTO{}, err
	}

	return ports.MonthlyDashboardDTO{
		UserID:             input.UserID,
		Period:             input.Period,
		TotalBalance:       totalBalance,
		MonthlyIncome:      monthlyIncome,
		MonthlyExpense:     monthlyExpense,
		MonthlyNetBalance:  monthlyIncome.Sub(monthlyExpense),
		AccountBalances:    nonNilAccountBalances(accountBalances),
		ExpensesByCategory: nonNilExpensesByCategory(expensesByCategory),
		RecentTransactions: nonNilRecentTransactions(recentTransactions),
	}, nil
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

func nonNilRecentTransactions(values []ports.TransactionDTO) []ports.TransactionDTO {
	if values == nil {
		return []ports.TransactionDTO{}
	}
	return values
}
