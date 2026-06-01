package persistence

import (
	"context"
	"os"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	dashboarddomain "contai/internal/dashboard/domain"
	"contai/internal/database"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"

	"github.com/google/uuid"
)

func TestDashboardRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := database.OpenPostgres(dsn)
	if err != nil {
		t.Fatalf("expected database connection, got %v", err)
	}
	if err := db.AutoMigrate(&dashboardAccountEntity{}, &dashboardCategoryEntity{}, &dashboardTransactionEntity{}); err != nil {
		t.Fatalf("expected migration to succeed, got %v", err)
	}

	repository := NewRepository(db)
	ctx := context.Background()
	userID := userdomain.UserID(uuid.NewString())
	otherUserID := userdomain.UserID(uuid.NewString())
	startAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC)
	period, err := dashboarddomain.NewPeriod(startAt, endAt)
	if err != nil {
		t.Fatalf("expected period, got %v", err)
	}

	activeAccountID := uuid.NewString()
	inactiveAccountID := uuid.NewString()
	expenseCategoryID := uuid.NewString()
	incomeCategoryID := uuid.NewString()
	if err := db.Create([]dashboardAccountEntity{
		{ID: activeAccountID, UserID: string(userID), Name: "Checking", Type: string(accountdomain.AccountTypeChecking), InitialBalance: 0, CurrentBalance: 8500, BankIconID: "bank", Status: string(accountdomain.AccountStatusActive), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: inactiveAccountID, UserID: string(userID), Name: "Closed", Type: string(accountdomain.AccountTypeCash), InitialBalance: 0, CurrentBalance: 9900, BankIconID: "cash", Status: string(accountdomain.AccountStatusInactive), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}).Error; err != nil {
		t.Fatalf("expected accounts to be created, got %v", err)
	}
	if err := db.Create([]dashboardCategoryEntity{
		{ID: expenseCategoryID, UserID: string(userID), Name: "Groceries", NormalizedName: "groceries", Type: string(categorydomain.CategoryTypeExpense), Color: "#2563EB", Icon: "tag", IsDefault: false, Status: string(categorydomain.CategoryStatusActive), CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: incomeCategoryID, UserID: string(userID), Name: "Salary", NormalizedName: "salary", Type: string(categorydomain.CategoryTypeIncome), Color: "#16A34A", Icon: "briefcase", IsDefault: false, Status: string(categorydomain.CategoryStatusActive), CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}).Error; err != nil {
		t.Fatalf("expected categories to be created, got %v", err)
	}

	removedAt := time.Now()
	transactions := []dashboardTransactionEntity{
		{ID: uuid.NewString(), UserID: string(userID), Type: string(transactiondomain.TransactionTypeIncome), Description: "Salary", Amount: 10000, OccurredAt: startAt, AccountID: &activeAccountID, CategoryID: &incomeCategoryID, Status: string(transactiondomain.TransactionStatusActive), Note: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: string(userID), Type: string(transactiondomain.TransactionTypeExpense), Description: "Market", Amount: 2500, OccurredAt: startAt.Add(24 * time.Hour), AccountID: &activeAccountID, CategoryID: &expenseCategoryID, Status: string(transactiondomain.TransactionStatusActive), Note: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: string(userID), Type: string(transactiondomain.TransactionTypeExpense), Description: "Removed", Amount: 9999, OccurredAt: startAt.Add(48 * time.Hour), AccountID: &activeAccountID, CategoryID: &expenseCategoryID, Status: string(transactiondomain.TransactionStatusRemoved), Note: "", RemovedAt: &removedAt, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: string(userID), Type: string(transactiondomain.TransactionTypeExpense), Description: "Outside", Amount: 7777, OccurredAt: endAt.Add(time.Second), AccountID: &activeAccountID, CategoryID: &expenseCategoryID, Status: string(transactiondomain.TransactionStatusActive), Note: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.NewString(), UserID: string(otherUserID), Type: string(transactiondomain.TransactionTypeExpense), Description: "Other user", Amount: 8888, OccurredAt: startAt, AccountID: &activeAccountID, CategoryID: &expenseCategoryID, Status: string(transactiondomain.TransactionStatusActive), Note: "", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	if err := db.Create(&transactions).Error; err != nil {
		t.Fatalf("expected transactions to be created, got %v", err)
	}

	accounts, err := repository.FindActiveAccountBalances(ctx, userID)
	if err != nil {
		t.Fatalf("expected active balances, got %v", err)
	}
	if len(accounts) != 1 || accounts[0].Balance.Cents() != 8500 || accounts[0].BankIconID != "bank" {
		t.Fatalf("expected one active balance 8500, got %#v", accounts)
	}

	income, err := repository.SumIncome(ctx, userID, period)
	if err != nil {
		t.Fatalf("expected income sum, got %v", err)
	}
	if income.Cents() != 10000 {
		t.Fatalf("expected income 10000, got %d", income.Cents())
	}
	expense, err := repository.SumExpense(ctx, userID, period)
	if err != nil {
		t.Fatalf("expected expense sum, got %v", err)
	}
	if expense.Cents() != 2500 {
		t.Fatalf("expected expense 2500, got %d", expense.Cents())
	}

	expensesByCategory, err := repository.FindExpensesByCategory(ctx, userID, period)
	if err != nil {
		t.Fatalf("expected expenses by category, got %v", err)
	}
	if len(expensesByCategory) != 1 ||
		expensesByCategory[0].CategoryID != categorydomain.CategoryID(expenseCategoryID) ||
		expensesByCategory[0].Name != "Groceries" ||
		expensesByCategory[0].Color != "#2563EB" ||
		expensesByCategory[0].Icon != "tag" ||
		expensesByCategory[0].Total.Cents() != 2500 {
		t.Fatalf("expected groceries category total 2500, got %#v", expensesByCategory)
	}

	recent, err := repository.FindRecentTransactions(ctx, userID, 5)
	if err != nil {
		t.Fatalf("expected recent transactions, got %v", err)
	}
	if len(recent) != 3 {
		t.Fatalf("expected three active recent transactions, got %#v", recent)
	}
	if recent[0].Description != "Outside" {
		t.Fatalf("expected most recent transaction first, got %#v", recent)
	}
}

type dashboardAccountEntity struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	UserID         string `gorm:"type:uuid;not null;index"`
	Name           string `gorm:"not null"`
	Type           string `gorm:"not null"`
	InitialBalance int64  `gorm:"not null"`
	CurrentBalance int64  `gorm:"not null"`
	BankIconID     string `gorm:"not null"`
	Status         string `gorm:"not null;index"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (dashboardAccountEntity) TableName() string {
	return "accounts"
}

type dashboardCategoryEntity struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	UserID         string `gorm:"type:uuid;not null;index"`
	Name           string `gorm:"not null"`
	NormalizedName string `gorm:"not null"`
	Type           string `gorm:"not null"`
	Color          string `gorm:"not null"`
	Icon           string `gorm:"not null"`
	IsDefault      bool   `gorm:"not null"`
	Status         string `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (dashboardCategoryEntity) TableName() string {
	return "categories"
}

type dashboardTransactionEntity struct {
	ID                   string    `gorm:"type:uuid;primaryKey"`
	UserID               string    `gorm:"type:uuid;not null;index"`
	Type                 string    `gorm:"not null;index"`
	Description          string    `gorm:"not null"`
	Amount               int64     `gorm:"not null"`
	OccurredAt           time.Time `gorm:"not null;index"`
	AccountID            *string   `gorm:"type:uuid;index"`
	SourceAccountID      *string   `gorm:"type:uuid;index"`
	DestinationAccountID *string   `gorm:"type:uuid;index"`
	CategoryID           *string   `gorm:"type:uuid;index"`
	Status               string    `gorm:"not null;index"`
	Note                 string    `gorm:"not null"`
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (dashboardTransactionEntity) TableName() string {
	return "transactions"
}
