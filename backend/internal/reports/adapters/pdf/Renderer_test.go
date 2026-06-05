package pdf

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"
)

func TestRendererProducesAccountsPDF(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("expected renderer, got %v", err)
	}

	pdfBytes, err := renderer.RenderAccountsReport(context.Background(), reportports.AccountsReportDTO{
		GeneratedAt: time.Date(2026, 6, 5, 10, 30, 0, 0, time.UTC),
		Accounts: []reportports.AccountReportRow{
			{
				ID:                      "account-id",
				Name:                    "Checking",
				Type:                    accountdomain.AccountTypeChecking,
				Status:                  accountdomain.AccountStatusActive,
				InitialBalance:          financedomain.NewMoney(1000),
				CurrentBalance:          financedomain.NewMoney(1500),
				IncludeInDashboardTotal: true,
			},
		},
		TotalBalance:   financedomain.NewMoney(1500),
		DashboardTotal: financedomain.NewMoney(1500),
	})
	if err != nil {
		t.Fatalf("expected pdf, got %v", err)
	}
	assertPDF(t, pdfBytes)
}

func TestRendererProducesFinancialPDF(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("expected renderer, got %v", err)
	}

	accountID := accountdomain.AccountID("account-id")
	pdfBytes, err := renderer.RenderFinancialReport(context.Background(), reportports.FinancialReportDTO{
		Title:       "Relatorio geral por periodo",
		Subtitle:    "Consolidado",
		GeneratedAt: time.Date(2026, 6, 5, 10, 30, 0, 0, time.UTC),
		StartAt:     time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		EndAt:       time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
		Transactions: []reportports.ReportTransactionRow{
			{
				ID:          "transaction-id",
				Type:        transactiondomain.TransactionTypeIncome,
				Description: "Recebimento mensal",
				Amount:      financedomain.NewMoney(50000),
				OccurredAt:  time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC),
				AccountID:   &accountID,
				AccountName: "Conta principal",
			},
		},
		IncomeTotal: financedomain.NewMoney(50000),
		NetTotal:    financedomain.NewMoney(50000),
	})
	if err != nil {
		t.Fatalf("expected pdf, got %v", err)
	}
	assertPDF(t, pdfBytes)
}

func TestRendererProducesEmptyFinancialPDF(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("expected renderer, got %v", err)
	}

	pdfBytes, err := renderer.RenderFinancialReport(context.Background(), reportports.FinancialReportDTO{
		Title:       "Relatorio de despesas por periodo",
		GeneratedAt: time.Date(2026, 6, 5, 10, 30, 0, 0, time.UTC),
		StartAt:     time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		EndAt:       time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("expected pdf, got %v", err)
	}
	assertPDF(t, pdfBytes)
}

func TestRendererPaginatesFinancialPDF(t *testing.T) {
	renderer, err := NewRenderer()
	if err != nil {
		t.Fatalf("expected renderer, got %v", err)
	}

	transactions := make([]reportports.ReportTransactionRow, 0, 90)
	accountID := accountdomain.AccountID("account-id")
	for index := 0; index < 90; index++ {
		transactions = append(transactions, reportports.ReportTransactionRow{
			ID:          transactiondomain.TransactionID(fmt.Sprintf("transaction-%d", index)),
			Type:        transactiondomain.TransactionTypeExpense,
			Description: fmt.Sprintf("Despesa recorrente numero %d com descricao longa para truncamento", index),
			Amount:      financedomain.NewMoney(1000),
			OccurredAt:  time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, index%30),
			AccountID:   &accountID,
			AccountName: "Conta principal",
		})
	}

	pdfBytes, err := renderer.RenderFinancialReport(context.Background(), reportports.FinancialReportDTO{
		Title:        "Relatorio de despesas por periodo",
		GeneratedAt:  time.Date(2026, 6, 5, 10, 30, 0, 0, time.UTC),
		StartAt:      time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC),
		EndAt:        time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
		Transactions: transactions,
		ExpenseTotal: financedomain.NewMoney(90000),
		NetTotal:     financedomain.NewMoney(-90000),
	})
	if err != nil {
		t.Fatalf("expected pdf, got %v", err)
	}
	assertPDF(t, pdfBytes)
	pageCount := extractPDFPageCount(t, pdfBytes)
	if pageCount < 2 {
		t.Fatalf("expected paginated pdf, got %d page", pageCount)
	}
}

func assertPDF(t *testing.T, pdfBytes []byte) {
	t.Helper()
	if len(pdfBytes) < 4 {
		t.Fatalf("expected non-empty pdf, got %d bytes", len(pdfBytes))
	}
	if !strings.HasPrefix(string(pdfBytes[:4]), "%PDF") {
		t.Fatalf("expected pdf header, got %q", string(pdfBytes[:4]))
	}
}

func extractPDFPageCount(t *testing.T, pdfBytes []byte) int {
	t.Helper()
	matches := regexp.MustCompile(`/Count\s+(\d+)`).FindStringSubmatch(string(pdfBytes))
	if len(matches) != 2 {
		t.Fatal("expected pdf page count")
	}
	pageCount, err := strconv.Atoi(matches[1])
	if err != nil {
		t.Fatalf("expected numeric pdf page count, got %q", matches[1])
	}
	return pageCount
}
