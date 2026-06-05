package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	accountdomain "contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"
)

var _ reportports.ReportService = ReportService{}

var (
	ErrReportPeriodInvalid          = errors.New("report period is invalid")
	ErrReportTransactionTypeInvalid = errors.New("report transaction type is invalid")
	ErrReportMonthlyPeriodInvalid   = errors.New("report monthly period must be within one month")
	ErrReportAccountIDRequired      = errors.New("report account id is required")
	ErrReportAccountNotFound        = errors.New("report account not found")
)

type ReportService struct {
	repository reportports.ReportRepository
	renderer   reportports.PDFRenderer
}

func NewReportService(repository reportports.ReportRepository, renderer reportports.PDFRenderer) ReportService {
	return ReportService{
		repository: repository,
		renderer:   renderer,
	}
}

func (service ReportService) GenerateAccountsPDF(ctx context.Context, input reportports.GenerateAccountsReportInput) (reportports.PDFFile, error) {
	if input.UserID == "" {
		return reportports.PDFFile{}, accountdomain.ErrAccountUserIDRequired
	}

	generatedAt := input.Now
	if generatedAt.IsZero() {
		generatedAt = time.Now()
	}

	accounts, err := service.repository.ListAccounts(ctx, input.UserID)
	if err != nil {
		return reportports.PDFFile{}, err
	}

	report := buildAccountsReport(accounts, generatedAt)
	content, err := service.renderer.RenderAccountsReport(ctx, report)
	if err != nil {
		return reportports.PDFFile{}, err
	}

	return reportports.PDFFile{
		Filename: "contai-relatorio-contas.pdf",
		Content:  content,
	}, nil
}

func (service ReportService) GenerateTransactionsPDF(ctx context.Context, input reportports.GenerateTransactionsReportInput) (reportports.PDFFile, error) {
	if input.UserID == "" {
		return reportports.PDFFile{}, accountdomain.ErrAccountUserIDRequired
	}
	if err := validatePeriod(input.StartAt, input.EndAt); err != nil {
		return reportports.PDFFile{}, err
	}
	if input.Type != transactiondomain.TransactionTypeIncome && input.Type != transactiondomain.TransactionTypeExpense {
		return reportports.PDFFile{}, ErrReportTransactionTypeInvalid
	}

	transactions, err := service.repository.ListTransactions(ctx, reportports.ListReportTransactionsInput{
		UserID:  input.UserID,
		StartAt: input.StartAt,
		EndAt:   input.EndAt,
		Type:    &input.Type,
	})
	if err != nil {
		return reportports.PDFFile{}, err
	}

	title := "Relatorio de receitas por periodo"
	filename := "contai-relatorio-receitas.pdf"
	if input.Type == transactiondomain.TransactionTypeExpense {
		title = "Relatorio de despesas por periodo"
		filename = "contai-relatorio-despesas.pdf"
	}
	report := buildFinancialReport(title, "", transactions, input.StartAt, input.EndAt, reportGeneratedAt(input.Now))
	content, err := service.renderer.RenderFinancialReport(ctx, report)
	if err != nil {
		return reportports.PDFFile{}, err
	}

	return reportports.PDFFile{Filename: filename, Content: content}, nil
}

func (service ReportService) GeneratePeriodPDF(ctx context.Context, input reportports.PeriodReportInput) (reportports.PDFFile, error) {
	if input.UserID == "" {
		return reportports.PDFFile{}, accountdomain.ErrAccountUserIDRequired
	}
	if err := validatePeriod(input.StartAt, input.EndAt); err != nil {
		return reportports.PDFFile{}, err
	}

	transactions, err := service.repository.ListTransactions(ctx, reportports.ListReportTransactionsInput{
		UserID:  input.UserID,
		StartAt: input.StartAt,
		EndAt:   input.EndAt,
	})
	if err != nil {
		return reportports.PDFFile{}, err
	}

	report := buildFinancialReport(
		"Relatorio geral por periodo",
		"Consolidado de receitas, despesas e transferencias",
		transactions,
		input.StartAt,
		input.EndAt,
		reportGeneratedAt(input.Now),
	)
	content, err := service.renderer.RenderFinancialReport(ctx, report)
	if err != nil {
		return reportports.PDFFile{}, err
	}

	return reportports.PDFFile{Filename: "contai-relatorio-periodo.pdf", Content: content}, nil
}

func (service ReportService) GenerateMonthlyPDF(ctx context.Context, input reportports.PeriodReportInput) (reportports.PDFFile, error) {
	if input.UserID == "" {
		return reportports.PDFFile{}, accountdomain.ErrAccountUserIDRequired
	}
	if err := validatePeriod(input.StartAt, input.EndAt); err != nil {
		return reportports.PDFFile{}, err
	}
	if input.StartAt.Year() != input.EndAt.Year() || input.StartAt.Month() != input.EndAt.Month() {
		return reportports.PDFFile{}, ErrReportMonthlyPeriodInvalid
	}

	transactions, err := service.repository.ListTransactions(ctx, reportports.ListReportTransactionsInput{
		UserID:  input.UserID,
		StartAt: input.StartAt,
		EndAt:   input.EndAt,
	})
	if err != nil {
		return reportports.PDFFile{}, err
	}

	report := buildFinancialReport(
		"Relatorio mensal consolidado",
		input.StartAt.Format("01/2006"),
		transactions,
		input.StartAt,
		input.EndAt,
		reportGeneratedAt(input.Now),
	)
	content, err := service.renderer.RenderFinancialReport(ctx, report)
	if err != nil {
		return reportports.PDFFile{}, err
	}

	return reportports.PDFFile{Filename: fmt.Sprintf("contai-relatorio-mensal-%s.pdf", input.StartAt.Format("2006-01")), Content: content}, nil
}

func (service ReportService) GenerateAccountPDF(ctx context.Context, input reportports.GenerateAccountReportInput) (reportports.PDFFile, error) {
	if input.UserID == "" {
		return reportports.PDFFile{}, accountdomain.ErrAccountUserIDRequired
	}
	if input.AccountID == "" {
		return reportports.PDFFile{}, ErrReportAccountIDRequired
	}
	if err := validatePeriod(input.StartAt, input.EndAt); err != nil {
		return reportports.PDFFile{}, err
	}

	account, err := service.repository.FindAccountByID(ctx, input.UserID, input.AccountID)
	if err != nil {
		return reportports.PDFFile{}, err
	}
	if account == nil {
		return reportports.PDFFile{}, ErrReportAccountNotFound
	}

	transactions, err := service.repository.ListTransactions(ctx, reportports.ListReportTransactionsInput{
		UserID:    input.UserID,
		StartAt:   input.StartAt,
		EndAt:     input.EndAt,
		AccountID: &input.AccountID,
	})
	if err != nil {
		return reportports.PDFFile{}, err
	}

	report := buildFinancialReport(
		"Relatorio por conta bancaria",
		account.Name,
		transactions,
		input.StartAt,
		input.EndAt,
		reportGeneratedAt(input.Now),
	)
	report.AccountName = account.Name
	content, err := service.renderer.RenderFinancialReport(ctx, report)
	if err != nil {
		return reportports.PDFFile{}, err
	}

	return reportports.PDFFile{Filename: fmt.Sprintf("contai-relatorio-conta-%s.pdf", input.AccountID), Content: content}, nil
}

func buildAccountsReport(accounts []reportports.AccountReportRow, generatedAt time.Time) reportports.AccountsReportDTO {
	rows := make([]reportports.AccountReportRow, 0, len(accounts))
	var total financedomain.Money
	var dashboardTotal financedomain.Money

	for _, account := range accounts {
		rows = append(rows, reportports.AccountReportRow{
			ID:                      account.ID,
			Name:                    account.Name,
			Type:                    account.Type,
			Status:                  account.Status,
			InitialBalance:          account.InitialBalance,
			CurrentBalance:          account.CurrentBalance,
			IncludeInDashboardTotal: account.IncludeInDashboardTotal,
		})
		total = total.Add(account.CurrentBalance)
		if account.Status == accountdomain.AccountStatusActive && account.IncludeInDashboardTotal {
			dashboardTotal = dashboardTotal.Add(account.CurrentBalance)
		}
	}

	return reportports.AccountsReportDTO{
		GeneratedAt:    generatedAt,
		Accounts:       rows,
		TotalBalance:   total,
		DashboardTotal: dashboardTotal,
	}
}

func validatePeriod(startAt, endAt time.Time) error {
	if startAt.IsZero() || endAt.IsZero() || endAt.Before(startAt) {
		return ErrReportPeriodInvalid
	}
	return nil
}

func reportGeneratedAt(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now()
	}
	return value
}

func buildFinancialReport(title string, subtitle string, transactions []reportports.ReportTransactionRow, startAt, endAt, generatedAt time.Time) reportports.FinancialReportDTO {
	var incomeTotal financedomain.Money
	var expenseTotal financedomain.Money
	var transferInTotal financedomain.Money
	var transferOutTotal financedomain.Money

	for _, transaction := range transactions {
		switch transaction.Type {
		case transactiondomain.TransactionTypeIncome:
			incomeTotal = incomeTotal.Add(transaction.Amount)
		case transactiondomain.TransactionTypeExpense:
			expenseTotal = expenseTotal.Add(transaction.Amount)
		case transactiondomain.TransactionTypeTransfer:
			transferInTotal = transferInTotal.Add(transaction.Amount)
			transferOutTotal = transferOutTotal.Add(transaction.Amount)
		}
	}

	return reportports.FinancialReportDTO{
		Title:            title,
		Subtitle:         subtitle,
		GeneratedAt:      generatedAt,
		StartAt:          startAt,
		EndAt:            endAt,
		Transactions:     transactions,
		IncomeTotal:      incomeTotal,
		ExpenseTotal:     expenseTotal,
		TransferInTotal:  transferInTotal,
		TransferOutTotal: transferOutTotal,
		NetTotal:         incomeTotal.Sub(expenseTotal),
	}
}
