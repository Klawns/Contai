package services

import (
	"context"
	"errors"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

func TestReportServiceGeneratesAccountsPDFWithTotals(t *testing.T) {
	generatedAt := time.Date(2026, 6, 5, 10, 30, 0, 0, time.UTC)
	repository := &fakeReportRepository{
		accounts: []reportports.AccountReportRow{
			{
				ID:                      "checking",
				Name:                    "Checking",
				Type:                    accountdomain.AccountTypeChecking,
				Status:                  accountdomain.AccountStatusActive,
				InitialBalance:          financedomain.NewMoney(1000),
				CurrentBalance:          financedomain.NewMoney(1500),
				IncludeInDashboardTotal: true,
			},
			{
				ID:                      "cash",
				Name:                    "Cash",
				Type:                    accountdomain.AccountTypeCash,
				Status:                  accountdomain.AccountStatusInactive,
				InitialBalance:          financedomain.NewMoney(500),
				CurrentBalance:          financedomain.NewMoney(700),
				IncludeInDashboardTotal: true,
			},
			{
				ID:                      "savings",
				Name:                    "Savings",
				Type:                    accountdomain.AccountTypeSavings,
				Status:                  accountdomain.AccountStatusActive,
				InitialBalance:          financedomain.NewMoney(200),
				CurrentBalance:          financedomain.NewMoney(300),
				IncludeInDashboardTotal: false,
			},
		},
	}
	renderer := &fakePDFRenderer{content: []byte("%PDF")}
	service := NewReportService(repository, renderer)

	file, err := service.GenerateAccountsPDF(context.Background(), reportports.GenerateAccountsReportInput{
		UserID: "user-id",
		Now:    generatedAt,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if file.Filename != "contai-relatorio-contas.pdf" {
		t.Fatalf("expected accounts filename, got %s", file.Filename)
	}
	if string(file.Content) != "%PDF" {
		t.Fatalf("expected renderer content, got %q", string(file.Content))
	}
	if repository.listAccountsUserID != "user-id" {
		t.Fatalf("expected repository to list accounts for user, got %s", repository.listAccountsUserID)
	}
	if renderer.report.TotalBalance.Cents() != 2500 {
		t.Fatalf("expected total balance 2500, got %d", renderer.report.TotalBalance.Cents())
	}
	if renderer.report.DashboardTotal.Cents() != 1500 {
		t.Fatalf("expected dashboard total 1500, got %d", renderer.report.DashboardTotal.Cents())
	}
	if !renderer.report.GeneratedAt.Equal(generatedAt) {
		t.Fatalf("expected generated at to be propagated, got %s", renderer.report.GeneratedAt)
	}
}

func TestReportServiceRequiresUserID(t *testing.T) {
	service := NewReportService(&fakeReportRepository{}, &fakePDFRenderer{})

	_, err := service.GenerateAccountsPDF(context.Background(), reportports.GenerateAccountsReportInput{})

	if !errors.Is(err, accountdomain.ErrAccountUserIDRequired) {
		t.Fatalf("expected user id error, got %v", err)
	}
}

func TestReportServiceGeneratesPeriodPDFWithTotals(t *testing.T) {
	startAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC)
	repository := &fakeReportRepository{
		transactions: []reportports.ReportTransactionRow{
			{
				ID:          "income-id",
				Type:        transactiondomain.TransactionTypeIncome,
				Description: "Salary",
				Amount:      financedomain.NewMoney(10000),
				OccurredAt:  startAt,
			},
			{
				ID:          "expense-id",
				Type:        transactiondomain.TransactionTypeExpense,
				Description: "Market",
				Amount:      financedomain.NewMoney(3500),
				OccurredAt:  endAt,
			},
		},
	}
	renderer := &fakePDFRenderer{content: []byte("%PDF")}
	service := NewReportService(repository, renderer)

	file, err := service.GeneratePeriodPDF(context.Background(), reportports.PeriodReportInput{
		UserID:  "user-id",
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if file.Filename != "contai-relatorio-periodo.pdf" {
		t.Fatalf("expected period filename, got %s", file.Filename)
	}
	if repository.listInput.UserID != "user-id" {
		t.Fatalf("expected repository to filter by user id, got %s", repository.listInput.UserID)
	}
	if renderer.financialReport.IncomeTotal.Cents() != 10000 {
		t.Fatalf("expected income total 10000, got %d", renderer.financialReport.IncomeTotal.Cents())
	}
	if renderer.financialReport.ExpenseTotal.Cents() != 3500 {
		t.Fatalf("expected expense total 3500, got %d", renderer.financialReport.ExpenseTotal.Cents())
	}
	if renderer.financialReport.NetTotal.Cents() != 6500 {
		t.Fatalf("expected net total 6500, got %d", renderer.financialReport.NetTotal.Cents())
	}
}

type fakePDFRenderer struct {
	report          reportports.AccountsReportDTO
	financialReport reportports.FinancialReportDTO
	content         []byte
	err             error
}

func (renderer *fakePDFRenderer) RenderAccountsReport(ctx context.Context, report reportports.AccountsReportDTO) ([]byte, error) {
	renderer.report = report
	if renderer.err != nil {
		return nil, renderer.err
	}
	return renderer.content, nil
}

func (renderer *fakePDFRenderer) RenderFinancialReport(ctx context.Context, report reportports.FinancialReportDTO) ([]byte, error) {
	renderer.financialReport = report
	if renderer.err != nil {
		return nil, renderer.err
	}
	return renderer.content, nil
}

type fakeReportRepository struct {
	account            *reportports.AccountReportRow
	accounts           []reportports.AccountReportRow
	transactions       []reportports.ReportTransactionRow
	listInput          reportports.ListReportTransactionsInput
	listAccountsUserID userdomain.UserID
	err                error
}

func (repository *fakeReportRepository) FindAccountByID(ctx context.Context, userID userdomain.UserID, accountID accountdomain.AccountID) (*reportports.AccountReportRow, error) {
	if repository.err != nil {
		return nil, repository.err
	}
	return repository.account, nil
}

func (repository *fakeReportRepository) ListAccounts(ctx context.Context, userID userdomain.UserID) ([]reportports.AccountReportRow, error) {
	repository.listAccountsUserID = userID
	if repository.err != nil {
		return nil, repository.err
	}
	return repository.accounts, nil
}

func (repository *fakeReportRepository) ListTransactions(ctx context.Context, input reportports.ListReportTransactionsInput) ([]reportports.ReportTransactionRow, error) {
	repository.listInput = input
	if repository.err != nil {
		return nil, repository.err
	}
	return repository.transactions, nil
}
