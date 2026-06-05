package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type PDFFile struct {
	Filename string
	Content  []byte
}

type GenerateAccountsReportInput struct {
	UserID userdomain.UserID
	Now    time.Time
}

type PeriodReportInput struct {
	UserID  userdomain.UserID
	StartAt time.Time
	EndAt   time.Time
	Now     time.Time
}

type GenerateTransactionsReportInput struct {
	UserID  userdomain.UserID
	StartAt time.Time
	EndAt   time.Time
	Type    transactiondomain.TransactionType
	Now     time.Time
}

type GenerateAccountReportInput struct {
	UserID    userdomain.UserID
	AccountID accountdomain.AccountID
	StartAt   time.Time
	EndAt     time.Time
	Now       time.Time
}

type AccountReportRow struct {
	ID                      accountdomain.AccountID
	Name                    string
	Type                    accountdomain.AccountType
	Status                  accountdomain.AccountStatus
	InitialBalance          financedomain.Money
	CurrentBalance          financedomain.Money
	IncludeInDashboardTotal bool
}

type AccountsReportDTO struct {
	GeneratedAt    time.Time
	Accounts       []AccountReportRow
	TotalBalance   financedomain.Money
	DashboardTotal financedomain.Money
}

type ReportTransactionRow struct {
	ID                   transactiondomain.TransactionID
	Type                 transactiondomain.TransactionType
	Description          string
	Amount               financedomain.Money
	OccurredAt           time.Time
	AccountID            *accountdomain.AccountID
	SourceAccountID      *accountdomain.AccountID
	DestinationAccountID *accountdomain.AccountID
	AccountName          string
}

type FinancialReportDTO struct {
	Title            string
	Subtitle         string
	GeneratedAt      time.Time
	StartAt          time.Time
	EndAt            time.Time
	AccountName      string
	Transactions     []ReportTransactionRow
	IncomeTotal      financedomain.Money
	ExpenseTotal     financedomain.Money
	TransferInTotal  financedomain.Money
	TransferOutTotal financedomain.Money
	NetTotal         financedomain.Money
}

type ReportService interface {
	GenerateAccountsPDF(ctx context.Context, input GenerateAccountsReportInput) (PDFFile, error)
	GenerateTransactionsPDF(ctx context.Context, input GenerateTransactionsReportInput) (PDFFile, error)
	GeneratePeriodPDF(ctx context.Context, input PeriodReportInput) (PDFFile, error)
	GenerateMonthlyPDF(ctx context.Context, input PeriodReportInput) (PDFFile, error)
	GenerateAccountPDF(ctx context.Context, input GenerateAccountReportInput) (PDFFile, error)
}
