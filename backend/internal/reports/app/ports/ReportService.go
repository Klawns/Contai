package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type PDFFile struct {
	Filename string
	Content  []byte
}

type MovementType string

const (
	MovementTypeAll               MovementType = "all"
	MovementTypeIncome            MovementType = "income"
	MovementTypeExpense           MovementType = "expense"
	MovementTypeCreditCardExpense MovementType = "credit_card_expense"
	MovementTypeTransfer          MovementType = "transfer"
)

type SettlementStatusFilter string

const (
	SettlementStatusAll     SettlementStatusFilter = "all"
	SettlementStatusSettled SettlementStatusFilter = "settled"
	SettlementStatusPending SettlementStatusFilter = "pending"
)

type ReportGroupBy string

const (
	ReportGroupByNone     ReportGroupBy = "none"
	ReportGroupByCategory ReportGroupBy = "category"
	ReportGroupByAccount  ReportGroupBy = "account"
	ReportGroupByDay      ReportGroupBy = "day"
	ReportGroupByMonth    ReportGroupBy = "month"
)

type FinancialReportInput struct {
	UserID           userdomain.UserID
	StartAt          time.Time
	EndAt            time.Time
	MovementType     MovementType
	CategoryID       *categorydomain.CategoryID
	AccountID        *accountdomain.AccountID
	SettlementStatus SettlementStatusFilter
	GroupBy          ReportGroupBy
	Now              time.Time
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

type FinancialMovementDTO struct {
	ID               string
	Source           MovementType
	Type             MovementType
	Description      string
	Amount           financedomain.Money
	OccurredAt       time.Time
	CategoryID       *categorydomain.CategoryID
	CategoryName     string
	AccountID        *accountdomain.AccountID
	AccountName      string
	SettlementStatus transactiondomain.SettlementStatus
}

type FinancialReportSummaryDTO struct {
	IncomeTotal  financedomain.Money
	ExpenseTotal financedomain.Money
	PeriodResult financedomain.Money
	PendingTotal financedomain.Money
	SettledTotal financedomain.Money
}

type FinancialReportGroupDTO struct {
	Key          string
	Label        string
	IncomeTotal  financedomain.Money
	ExpenseTotal financedomain.Money
	NetTotal     financedomain.Money
	Total        financedomain.Money
	Count        int
}

type FinancialReportSeriesPointDTO struct {
	Key          string
	Label        string
	IncomeTotal  financedomain.Money
	ExpenseTotal financedomain.Money
	NetTotal     financedomain.Money
}

type FinancialReportCategoryChartDTO struct {
	CategoryID categorydomain.CategoryID
	Name       string
	Total      financedomain.Money
}

type FinancialReportChartsDTO struct {
	IncomeVsExpense    []FinancialReportSeriesPointDTO
	ExpensesByCategory []FinancialReportCategoryChartDTO
	Evolution          []FinancialReportSeriesPointDTO
}

type FinancialReportDTO struct {
	Title            string
	Subtitle         string
	GeneratedAt      time.Time
	StartAt          time.Time
	EndAt            time.Time
	AccountName      string
	Summary          FinancialReportSummaryDTO
	Movements        []FinancialMovementDTO
	Groups           []FinancialReportGroupDTO
	Charts           FinancialReportChartsDTO
	Transactions     []ReportTransactionRow
	IncomeTotal      financedomain.Money
	ExpenseTotal     financedomain.Money
	TransferInTotal  financedomain.Money
	TransferOutTotal financedomain.Money
	NetTotal         financedomain.Money
}

type ReportService interface {
	GetFinancialReport(ctx context.Context, input FinancialReportInput) (FinancialReportDTO, error)
	GenerateFinancialPDF(ctx context.Context, input FinancialReportInput) (PDFFile, error)
}
