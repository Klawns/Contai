package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"contai/internal/dashboard/app/ports"
	"contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

func TestDashboardServiceRequiresUserID(t *testing.T) {
	service := NewDashboardService(&fakeDashboardRepository{})
	period := validPeriod(t)

	_, err := service.GetMonthlyDashboard(context.Background(), ports.GetMonthlyDashboardInput{Period: period})

	if !errors.Is(err, domain.ErrDashboardUserIDRequired) {
		t.Fatalf("expected user id required, got %v", err)
	}
}

func TestDashboardServiceRejectsInvalidPeriod(t *testing.T) {
	service := NewDashboardService(&fakeDashboardRepository{})

	_, err := service.GetMonthlyDashboard(context.Background(), ports.GetMonthlyDashboardInput{
		UserID: "user-id",
		Period: domain.Period{
			StartAt: time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			EndAt:   time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	})

	if !errors.Is(err, domain.ErrDashboardInvalidPeriod) {
		t.Fatalf("expected invalid period, got %v", err)
	}
}

func TestDashboardServiceCalculatesBalances(t *testing.T) {
	period := validPeriod(t)
	repository := &fakeDashboardRepository{
		accountBalances: []ports.AccountBalanceDTO{
			{AccountID: "checking", Balance: financedomain.NewMoney(1200)},
			{AccountID: "cash", Balance: financedomain.NewMoney(-200)},
		},
		income:  financedomain.NewMoney(5000),
		expense: financedomain.NewMoney(1250),
	}
	service := NewDashboardService(repository)

	dashboard, err := service.GetMonthlyDashboard(context.Background(), ports.GetMonthlyDashboardInput{
		UserID: "user-id",
		Period: period,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dashboard.TotalBalance.Cents() != 1000 {
		t.Fatalf("expected total balance 1000, got %d", dashboard.TotalBalance.Cents())
	}
	if dashboard.MonthlyNetBalance.Cents() != 3750 {
		t.Fatalf("expected monthly net balance 3750, got %d", dashboard.MonthlyNetBalance.Cents())
	}
}

func TestDashboardServiceUsesRecentTransactionsLimit(t *testing.T) {
	repository := &fakeDashboardRepository{}
	service := NewDashboardService(repository)

	_, err := service.GetMonthlyDashboard(context.Background(), ports.GetMonthlyDashboardInput{
		UserID: "user-id",
		Period: validPeriod(t),
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repository.recentLimit != 5 {
		t.Fatalf("expected recent transactions limit 5, got %d", repository.recentLimit)
	}
}

func TestDashboardServiceReturnsEmptySlices(t *testing.T) {
	service := NewDashboardService(&fakeDashboardRepository{})

	dashboard, err := service.GetMonthlyDashboard(context.Background(), ports.GetMonthlyDashboardInput{
		UserID: "user-id",
		Period: validPeriod(t),
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dashboard.AccountBalances == nil {
		t.Fatal("expected non-nil account balances")
	}
	if dashboard.ExpensesByCategory == nil {
		t.Fatal("expected non-nil expenses by category")
	}
	if dashboard.RecentTransactions == nil {
		t.Fatal("expected non-nil recent transactions")
	}
}

func validPeriod(t *testing.T) domain.Period {
	t.Helper()
	period, err := domain.NewPeriod(
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("expected valid period, got %v", err)
	}
	return period
}

type fakeDashboardRepository struct {
	accountBalances    []ports.AccountBalanceDTO
	income             financedomain.Money
	expense            financedomain.Money
	expensesByCategory []ports.CategoryExpenseDTO
	recentTransactions []ports.TransactionDTO
	recentLimit        int
}

func (repository *fakeDashboardRepository) FindActiveAccountBalances(ctx context.Context, userID userdomain.UserID) ([]ports.AccountBalanceDTO, error) {
	return repository.accountBalances, nil
}

func (repository *fakeDashboardRepository) SumIncome(ctx context.Context, userID userdomain.UserID, period domain.Period) (financedomain.Money, error) {
	return repository.income, nil
}

func (repository *fakeDashboardRepository) SumExpense(ctx context.Context, userID userdomain.UserID, period domain.Period) (financedomain.Money, error) {
	return repository.expense, nil
}

func (repository *fakeDashboardRepository) FindExpensesByCategory(ctx context.Context, userID userdomain.UserID, period domain.Period) ([]ports.CategoryExpenseDTO, error) {
	return repository.expensesByCategory, nil
}

func (repository *fakeDashboardRepository) FindRecentTransactions(ctx context.Context, userID userdomain.UserID, limit int) ([]ports.TransactionDTO, error) {
	repository.recentLimit = limit
	return repository.recentTransactions, nil
}
