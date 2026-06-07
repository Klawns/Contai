package services

import (
	"context"
	"errors"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	"contai/internal/dashboard/app/ports"
	"contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
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

func TestDashboardServiceMonthlySeriesRequiresUserID(t *testing.T) {
	service := NewDashboardService(&fakeDashboardRepository{})
	period := validPeriod(t)

	_, err := service.GetMonthlySeries(context.Background(), ports.GetMonthlySeriesInput{Period: period})

	if !errors.Is(err, domain.ErrDashboardUserIDRequired) {
		t.Fatalf("expected user id required, got %v", err)
	}
}

func TestDashboardServiceMonthlySeriesRejectsInvalidPeriod(t *testing.T) {
	service := NewDashboardService(&fakeDashboardRepository{})

	_, err := service.GetMonthlySeries(context.Background(), ports.GetMonthlySeriesInput{
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

func TestDashboardServiceMonthlySeriesFillsEmptyMonthsAndOrdersPoints(t *testing.T) {
	period, err := domain.NewPeriod(
		time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 20, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("expected period, got %v", err)
	}
	repository := &fakeDashboardRepository{
		monthlyIncomeExpense: []ports.MonthlyIncomeExpenseDTO{
			{MonthStartAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), Income: financedomain.NewMoney(5000), Expense: financedomain.NewMoney(1200)},
			{MonthStartAt: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC), Income: financedomain.NewMoney(7000), Expense: financedomain.NewMoney(2500)},
		},
		monthlyBalances: []ports.MonthlyBalanceDTO{
			{MonthEndAt: time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC), Balance: financedomain.NewMoney(3800)},
			{MonthEndAt: time.Date(2026, 2, 28, 23, 59, 59, 0, time.UTC), Balance: financedomain.NewMoney(3800)},
			{MonthEndAt: time.Date(2026, 3, 31, 23, 59, 59, 0, time.UTC), Balance: financedomain.NewMoney(8300)},
		},
	}
	service := NewDashboardService(repository)

	series, err := service.GetMonthlySeries(context.Background(), ports.GetMonthlySeriesInput{
		UserID: "user-id",
		Period: period,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(series.Points) != 3 {
		t.Fatalf("expected three points, got %#v", series.Points)
	}
	if series.Points[0].MonthStartAt.Format("2006-01-02") != "2026-01-01" ||
		series.Points[1].MonthStartAt.Format("2006-01-02") != "2026-02-01" ||
		series.Points[2].MonthStartAt.Format("2006-01-02") != "2026-03-01" {
		t.Fatalf("expected ordered months, got %#v", series.Points)
	}
	if series.Points[1].Income.Cents() != 0 || series.Points[1].Expense.Cents() != 0 || series.Points[1].Balance.Cents() != 3800 {
		t.Fatalf("expected empty february with balance, got %#v", series.Points[1])
	}
	if series.Points[2].Income.Cents() != 7000 || series.Points[2].Expense.Cents() != 2500 || series.Points[2].Balance.Cents() != 8300 {
		t.Fatalf("expected march totals, got %#v", series.Points[2])
	}
	if len(repository.monthEnds) != 3 {
		t.Fatalf("expected three month ends, got %#v", repository.monthEnds)
	}
}

func TestDashboardServiceCalculatesBalances(t *testing.T) {
	period := validPeriod(t)
	repository := &fakeDashboardRepository{
		accountBalances: []ports.AccountBalanceDTO{
			{AccountID: "checking", Balance: financedomain.NewMoney(1200), IncludeInDashboardTotal: true},
			{AccountID: "cash", Balance: financedomain.NewMoney(-200), IncludeInDashboardTotal: false},
		},
		transactions: []transactiondomain.Transaction{
			validIncomeTransaction(t, "income-id", 5000),
			validExpenseTransaction(t, "expense-id", 1250),
			validTransferTransaction(t, "transfer-id", 900),
		},
	}
	service := NewDashboardService(repository)

	dashboard, err := service.GetMonthlyDashboard(context.Background(), ports.GetMonthlyDashboardInput{
		UserID: "user-id",
		Period: period,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if dashboard.TotalBalance.Cents() != 1200 {
		t.Fatalf("expected total balance 1200, got %d", dashboard.TotalBalance.Cents())
	}
	if dashboard.MonthlyNetBalance.Cents() != 3750 {
		t.Fatalf("expected monthly net balance 3750, got %d", dashboard.MonthlyNetBalance.Cents())
	}
	if dashboard.MonthlyTransferIn.Cents() != 900 || dashboard.MonthlyTransferOut.Cents() != 900 {
		t.Fatalf("expected transfer totals 900/900, got %d/%d", dashboard.MonthlyTransferIn.Cents(), dashboard.MonthlyTransferOut.Cents())
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
	if dashboard.CreditCards == nil {
		t.Fatal("expected non-nil credit cards")
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
	accountBalances      []ports.AccountBalanceDTO
	creditCards          []ports.CreditCardDashboardDTO
	income               financedomain.Money
	expense              financedomain.Money
	expensesByCategory   []ports.CategoryExpenseDTO
	recentTransactions   []ports.TransactionDTO
	transactions         []transactiondomain.Transaction
	monthlyIncomeExpense []ports.MonthlyIncomeExpenseDTO
	monthlyBalances      []ports.MonthlyBalanceDTO
	monthEnds            []time.Time
	recentLimit          int
}

func (repository *fakeDashboardRepository) FindActiveAccountBalances(ctx context.Context, userID userdomain.UserID) ([]ports.AccountBalanceDTO, error) {
	return repository.accountBalances, nil
}

func (repository *fakeDashboardRepository) FindCreditCards(ctx context.Context, userID userdomain.UserID, now time.Time) ([]ports.CreditCardDashboardDTO, error) {
	return repository.creditCards, nil
}

func (repository *fakeDashboardRepository) SumIncome(ctx context.Context, userID userdomain.UserID, period domain.Period) (financedomain.Money, error) {
	return repository.income, nil
}

func (repository *fakeDashboardRepository) SumExpense(ctx context.Context, userID userdomain.UserID, period domain.Period) (financedomain.Money, error) {
	return repository.expense, nil
}

func (repository *fakeDashboardRepository) FindMonthlyIncomeExpense(ctx context.Context, userID userdomain.UserID, period domain.Period) ([]ports.MonthlyIncomeExpenseDTO, error) {
	return repository.monthlyIncomeExpense, nil
}

func (repository *fakeDashboardRepository) FindMonthlyBalances(ctx context.Context, userID userdomain.UserID, monthEnds []time.Time) ([]ports.MonthlyBalanceDTO, error) {
	repository.monthEnds = monthEnds
	return repository.monthlyBalances, nil
}

func (repository *fakeDashboardRepository) FindTransactionsByPeriod(ctx context.Context, userID userdomain.UserID, period domain.Period) ([]transactiondomain.Transaction, error) {
	return repository.transactions, nil
}

func (repository *fakeDashboardRepository) FindExpensesByCategory(ctx context.Context, userID userdomain.UserID, period domain.Period) ([]ports.CategoryExpenseDTO, error) {
	return repository.expensesByCategory, nil
}

func (repository *fakeDashboardRepository) FindRecentTransactions(ctx context.Context, userID userdomain.UserID, limit int) ([]ports.TransactionDTO, error) {
	repository.recentLimit = limit
	return repository.recentTransactions, nil
}

func validIncomeTransaction(t *testing.T, id transactiondomain.TransactionID, cents int64) transactiondomain.Transaction {
	t.Helper()
	transaction, err := transactiondomain.NewIncome(id, "user-id", "Income", financedomain.NewMoney(cents), time.Now(), "account-id", "category-id", "")
	if err != nil {
		t.Fatalf("expected valid income transaction, got %v", err)
	}
	return transaction
}

func validExpenseTransaction(t *testing.T, id transactiondomain.TransactionID, cents int64) transactiondomain.Transaction {
	t.Helper()
	transaction, err := transactiondomain.NewExpense(id, "user-id", "Expense", financedomain.NewMoney(cents), time.Now(), "account-id", "category-id", "")
	if err != nil {
		t.Fatalf("expected valid expense transaction, got %v", err)
	}
	return transaction
}

func validTransferTransaction(t *testing.T, id transactiondomain.TransactionID, cents int64) transactiondomain.Transaction {
	t.Helper()
	source := accountdomain.AccountID("source-account-id")
	destination := accountdomain.AccountID("destination-account-id")
	transaction, err := transactiondomain.NewTransfer(id, "user-id", "Transfer", financedomain.NewMoney(cents), time.Now(), source, destination, "")
	if err != nil {
		t.Fatalf("expected valid transfer transaction, got %v", err)
	}
	return transaction
}
