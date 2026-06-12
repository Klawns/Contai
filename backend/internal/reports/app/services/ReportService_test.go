package services

import (
	"context"
	"errors"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"
)

func TestReportServiceBuildsFinancialReportTotalsAndGroups(t *testing.T) {
	startAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC)
	categoryID := categorydomain.CategoryID("food")
	accountID := accountdomain.AccountID("checking")
	repository := &fakeReportRepository{
		movements: []reportports.FinancialMovementDTO{
			{
				ID:               "income-id",
				Type:             reportports.MovementTypeIncome,
				Description:      "Salary",
				Amount:           financedomain.NewMoney(10000),
				OccurredAt:       startAt,
				AccountID:        &accountID,
				AccountName:      "Checking",
				SettlementStatus: transactiondomain.SettlementStatusSettled,
			},
			{
				ID:               "expense-id",
				Type:             reportports.MovementTypeCreditCardExpense,
				Description:      "Market",
				Amount:           financedomain.NewMoney(3500),
				OccurredAt:       endAt,
				CategoryID:       &categoryID,
				CategoryName:     "Food",
				AccountID:        &accountID,
				AccountName:      "Checking",
				SettlementStatus: transactiondomain.SettlementStatusPending,
			},
			{
				ID:               "transfer-id",
				Type:             reportports.MovementTypeTransfer,
				Description:      "Move money",
				Amount:           financedomain.NewMoney(2500),
				OccurredAt:       endAt,
				AccountID:        &accountID,
				AccountName:      "Checking",
				SettlementStatus: transactiondomain.SettlementStatusSettled,
			},
		},
	}
	service := NewReportService(repository, &fakePDFRenderer{})

	report, err := service.GetFinancialReport(context.Background(), reportports.FinancialReportInput{
		UserID:  "user-id",
		StartAt: startAt,
		EndAt:   endAt,
		GroupBy: reportports.ReportGroupByAccount,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.input.UserID != "user-id" {
		t.Fatalf("expected repository to filter by user, got %s", repository.input.UserID)
	}
	if report.Summary.IncomeTotal.Cents() != 10000 {
		t.Fatalf("expected income total 10000, got %d", report.Summary.IncomeTotal.Cents())
	}
	if report.Summary.ExpenseTotal.Cents() != 3500 {
		t.Fatalf("expected expense total 3500, got %d", report.Summary.ExpenseTotal.Cents())
	}
	if report.Summary.PeriodResult.Cents() != 6500 {
		t.Fatalf("expected period result 6500, got %d", report.Summary.PeriodResult.Cents())
	}
	if report.Summary.PendingTotal.Cents() != 3500 || report.Summary.SettledTotal.Cents() != 10000 {
		t.Fatalf("expected status totals, got pending=%d settled=%d", report.Summary.PendingTotal.Cents(), report.Summary.SettledTotal.Cents())
	}
	if len(report.Groups) != 1 || report.Groups[0].Label != "Checking" {
		t.Fatalf("expected account group, got %#v", report.Groups)
	}
}

func TestReportServiceGeneratesFinancialPDFFromSameReport(t *testing.T) {
	startAt := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC)
	renderer := &fakePDFRenderer{content: []byte("%PDF")}
	service := NewReportService(&fakeReportRepository{}, renderer)

	file, err := service.GenerateFinancialPDF(context.Background(), reportports.FinancialReportInput{
		UserID:  "user-id",
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if file.Filename != "contai-relatorio-financeiro.pdf" || string(file.Content) != "%PDF" {
		t.Fatalf("expected pdf file, got %#v", file)
	}
	if renderer.financialReport.Title != "Relatorio financeiro" {
		t.Fatalf("expected financial report title, got %s", renderer.financialReport.Title)
	}
}

func TestReportServiceRequiresUserID(t *testing.T) {
	service := NewReportService(&fakeReportRepository{}, &fakePDFRenderer{})

	_, err := service.GetFinancialReport(context.Background(), reportports.FinancialReportInput{})

	if !errors.Is(err, accountdomain.ErrAccountUserIDRequired) {
		t.Fatalf("expected user id error, got %v", err)
	}
}

type fakePDFRenderer struct {
	financialReport reportports.FinancialReportDTO
	content         []byte
	err             error
}

func (renderer *fakePDFRenderer) RenderAccountsReport(ctx context.Context, report reportports.AccountsReportDTO) ([]byte, error) {
	return renderer.content, renderer.err
}

func (renderer *fakePDFRenderer) RenderFinancialReport(ctx context.Context, report reportports.FinancialReportDTO) ([]byte, error) {
	renderer.financialReport = report
	if renderer.err != nil {
		return nil, renderer.err
	}
	return renderer.content, nil
}

type fakeReportRepository struct {
	movements []reportports.FinancialMovementDTO
	input     reportports.ListFinancialMovementsInput
	err       error
}

func (repository *fakeReportRepository) ListFinancialMovements(ctx context.Context, input reportports.ListFinancialMovementsInput) ([]reportports.FinancialMovementDTO, error) {
	repository.input = input
	if repository.err != nil {
		return nil, repository.err
	}
	return repository.movements, nil
}
