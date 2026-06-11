package services

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	accountports "contai/internal/account/app/ports"
	accountdomain "contai/internal/account/domain"
	categoryports "contai/internal/category/app/ports"
	categorydomain "contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	"contai/internal/transactions/app/ports"
	"contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

func TestTransactionService_CreateExpenseUsesUnitOfWorkAndUpdatesBalance(t *testing.T) {
	account := validAccount(t, "account-id", financedomain.NewMoney(10000))
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	transactionRepository := &fakeTransactionRepository{}
	accountRepository := &fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{"account-id": &account}}
	categoryRepository := &fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}}
	unitOfWork := &fakeUnitOfWork{}
	service := NewTransactionService(transactionRepository, accountRepository, categoryRepository, fakeTransactionIDGenerator{}, unitOfWork)

	transaction, err := service.CreateExpense(context.Background(), ports.CreateExpenseInput{
		UserID:           "user-id",
		Description:      "Groceries",
		Amount:           financedomain.NewMoney(3500),
		OccurredAt:       time.Now(),
		AccountID:        accountIDPtr("account-id"),
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusSettled,
		RecurrenceType:   domain.RecurrenceTypeNone,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if transaction.Type != domain.TransactionTypeExpense {
		t.Fatalf("expected expense transaction, got %#v", transaction)
	}
	if !unitOfWork.called {
		t.Fatal("expected UnitOfWork to be used")
	}
	if accountRepository.accounts["account-id"].CurrentBalance.Cents() != 6500 {
		t.Fatalf("expected balance 6500, got %d", accountRepository.accounts["account-id"].CurrentBalance.Cents())
	}
}

func TestTransactionService_CreateIncomeRejectsExpenseCategory(t *testing.T) {
	account := validAccount(t, "account-id", financedomain.NewMoney(10000))
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	service := NewTransactionService(
		&fakeTransactionRepository{},
		&fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{"account-id": &account}},
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		fakeTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	_, err := service.CreateIncome(context.Background(), ports.CreateIncomeInput{
		UserID:           "user-id",
		Description:      "Salary",
		Amount:           financedomain.NewMoney(10000),
		OccurredAt:       time.Now(),
		AccountID:        accountIDPtr("account-id"),
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusSettled,
		RecurrenceType:   domain.RecurrenceTypeNone,
	})

	if !errors.Is(err, domain.ErrTransactionCategoryTypeMismatch) {
		t.Fatalf("expected category type mismatch, got %v", err)
	}
}

func TestTransactionService_UpdateRejectsManagedOrigin(t *testing.T) {
	transaction, err := domain.NewExpense(
		"transaction-id",
		"user-id",
		"Invoice",
		financedomain.NewMoney(1000),
		time.Now(),
		accountIDPtr("account-id"),
		"category-id",
		domain.SettlementStatusSettled,
		nil,
		domain.RecurrenceTypeNone,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}
	if err := transaction.SetOrigin(domain.TransactionOriginTypePayable, "commitment-id"); err != nil {
		t.Fatalf("expected origin to be set, got %v", err)
	}
	service := NewTransactionService(
		&fakeTransactionRepository{found: &transaction},
		&fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{}},
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{}},
		fakeTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	_, err = service.UpdateTransaction(context.Background(), ports.UpdateTransactionInput{
		UserID:           "user-id",
		TransactionID:    "transaction-id",
		Description:      "Invoice",
		Amount:           financedomain.NewMoney(1000),
		OccurredAt:       time.Now(),
		AccountID:        accountIDPtr("account-id"),
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusSettled,
		RecurrenceType:   domain.RecurrenceTypeNone,
	})

	if !errors.Is(err, domain.ErrTransactionManagedOrigin) {
		t.Fatalf("expected managed origin error, got %v", err)
	}
}

func TestTransactionService_CreatePendingExpenseWithoutAccountDoesNotUpdateBalance(t *testing.T) {
	account := validAccount(t, "account-id", financedomain.NewMoney(10000))
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	accountRepository := &fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{"account-id": &account}}
	service := NewTransactionService(
		&fakeTransactionRepository{},
		accountRepository,
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		fakeTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	transaction, err := service.CreateExpense(context.Background(), ports.CreateExpenseInput{
		UserID:           "user-id",
		Description:      "Groceries",
		Amount:           financedomain.NewMoney(3500),
		OccurredAt:       time.Now(),
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusPending,
		RecurrenceType:   domain.RecurrenceTypeNone,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if transaction.AccountID != nil {
		t.Fatalf("expected nil account, got %#v", transaction.AccountID)
	}
	if accountRepository.accounts["account-id"].CurrentBalance.Cents() != 10000 {
		t.Fatalf("expected unchanged balance 10000, got %d", accountRepository.accounts["account-id"].CurrentBalance.Cents())
	}
}

func TestTransactionService_UpdateSettlesPendingExpenseWithAccount(t *testing.T) {
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	account := validAccount(t, "account-id", financedomain.NewMoney(10000))
	transaction, err := domain.NewExpense(
		"transaction-id",
		"user-id",
		"Groceries",
		financedomain.NewMoney(3500),
		time.Now(),
		nil,
		"category-id",
		domain.SettlementStatusPending,
		nil,
		domain.RecurrenceTypeNone,
		nil,
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}
	accountRepository := &fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{"account-id": &account}}
	service := NewTransactionService(
		&fakeTransactionRepository{found: &transaction},
		accountRepository,
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		fakeTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	updated, err := service.UpdateTransaction(context.Background(), ports.UpdateTransactionInput{
		UserID:           "user-id",
		TransactionID:    "transaction-id",
		Description:      "Groceries",
		Amount:           financedomain.NewMoney(3500),
		OccurredAt:       transaction.OccurredAt,
		AccountID:        accountIDPtr("account-id"),
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusSettled,
		RecurrenceType:   domain.RecurrenceTypeNone,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.SettlementStatus != domain.SettlementStatusSettled || updated.SettledAt == nil {
		t.Fatalf("expected settled transaction, got %#v", updated)
	}
	if accountRepository.accounts["account-id"].CurrentBalance.Cents() != 6500 {
		t.Fatalf("expected balance 6500, got %d", accountRepository.accounts["account-id"].CurrentBalance.Cents())
	}
}

func TestTransactionService_CreateRepeatExpenseCreatesFiniteOccurrences(t *testing.T) {
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	transactionRepository := &fakeTransactionRepository{}
	startsAt := time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC)
	quantity := 3
	dayOfMonth := 31
	service := NewTransactionService(
		transactionRepository,
		&fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{}},
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		&sequenceTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	created, err := service.CreateExpense(context.Background(), ports.CreateExpenseInput{
		UserID:           "user-id",
		Description:      "Subscription",
		Amount:           financedomain.NewMoney(2000),
		OccurredAt:       startsAt,
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusPending,
		RecurrenceType:   domain.RecurrenceTypeRepeat,
		Recurrence: &domain.Recurrence{
			Frequency:  domain.RecurrenceFrequencyMonthly,
			Quantity:   &quantity,
			StartsAt:   startsAt,
			DayOfMonth: &dayOfMonth,
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID != "transaction-1" {
		t.Fatalf("expected first created transaction to be returned, got %#v", created)
	}
	if len(transactionRepository.createdTransactions) != 3 {
		t.Fatalf("expected 3 created transactions, got %d", len(transactionRepository.createdTransactions))
	}
	expectedDates := []time.Time{
		time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 28, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 3, 31, 12, 0, 0, 0, time.UTC),
	}
	for index, expected := range expectedDates {
		if !transactionRepository.createdTransactions[index].OccurredAt.Equal(expected) {
			t.Fatalf("expected occurrence %d at %s, got %s", index, expected, transactionRepository.createdTransactions[index].OccurredAt)
		}
	}
}

func TestTransactionService_CreateFixedExpensePersistsOnlyRule(t *testing.T) {
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	transactionRepository := &fakeTransactionRepository{}
	startsAt := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	service := NewTransactionService(
		transactionRepository,
		&fakeAccountRepository{accounts: map[accountdomain.AccountID]*accountdomain.Account{}},
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		&sequenceTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	_, err := service.CreateExpense(context.Background(), ports.CreateExpenseInput{
		UserID:           "user-id",
		Description:      "Rent",
		Amount:           financedomain.NewMoney(200000),
		OccurredAt:       startsAt,
		CategoryID:       "category-id",
		SettlementStatus: domain.SettlementStatusPending,
		RecurrenceType:   domain.RecurrenceTypeFixed,
		Recurrence: &domain.Recurrence{
			Frequency: domain.RecurrenceFrequencyMonthly,
			StartsAt:  startsAt,
		},
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(transactionRepository.createdTransactions) != 1 {
		t.Fatalf("expected fixed recurrence to create one transaction, got %d", len(transactionRepository.createdTransactions))
	}
	if transactionRepository.createdTransactions[0].RecurrenceType != domain.RecurrenceTypeFixed {
		t.Fatalf("expected fixed recurrence rule to be persisted, got %#v", transactionRepository.createdTransactions[0])
	}
}

func validAccount(t *testing.T, id accountdomain.AccountID, balance financedomain.Money) accountdomain.Account {
	t.Helper()
	account, err := accountdomain.NewAccount(id, "user-id", "Checking", accountdomain.AccountTypeChecking, balance, "bank")
	if err != nil {
		t.Fatalf("expected valid account, got %v", err)
	}
	return account
}

func validCategory(t *testing.T, id categorydomain.CategoryID, categoryType categorydomain.CategoryType) categorydomain.Category {
	t.Helper()
	category, err := categorydomain.NewCategory(id, "user-id", "Category", categoryType, "#2563EB", "tag", false)
	if err != nil {
		t.Fatalf("expected valid category, got %v", err)
	}
	return category
}

func accountIDPtr(value accountdomain.AccountID) *accountdomain.AccountID {
	return &value
}

type fakeUnitOfWork struct {
	called bool
}

func (unit *fakeUnitOfWork) WithinTx(ctx context.Context, fn func(context.Context, databaseports.TxHandle) error) error {
	unit.called = true
	return fn(ctx, databaseports.NewTxHandle(nil))
}

type fakeTransactionIDGenerator struct{}

func (generator fakeTransactionIDGenerator) NewTransactionID() domain.TransactionID {
	return "transaction-id"
}

type sequenceTransactionIDGenerator struct {
	next int
}

func (generator *sequenceTransactionIDGenerator) NewTransactionID() domain.TransactionID {
	generator.next++
	return domain.TransactionID("transaction-" + strconv.Itoa(generator.next))
}

type fakeTransactionRepository struct {
	created             *domain.Transaction
	createdTransactions []domain.Transaction
	updated             *domain.Transaction
	found               *domain.Transaction
}

func (repository *fakeTransactionRepository) WithTx(tx databaseports.TxHandle) ports.TransactionRepository {
	return repository
}

func (repository *fakeTransactionRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	repository.created = transaction
	repository.createdTransactions = append(repository.createdTransactions, *transaction)
	return transaction, nil
}

func (repository *fakeTransactionRepository) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	repository.updated = transaction
	return transaction, nil
}

func (repository *fakeTransactionRepository) FindTransactionByID(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID) (*domain.Transaction, error) {
	return repository.found, nil
}

func (repository *fakeTransactionRepository) FindTransactionByIDForUpdate(ctx context.Context, transactionID domain.TransactionID, userID userdomain.UserID) (*domain.Transaction, error) {
	return repository.found, nil
}

func (repository *fakeTransactionRepository) FindTransactionsByUserID(ctx context.Context, input ports.ListTransactionsInput) ([]domain.Transaction, error) {
	return nil, nil
}

type fakeAccountRepository struct {
	accounts map[accountdomain.AccountID]*accountdomain.Account
}

func (repository *fakeAccountRepository) WithTx(tx databaseports.TxHandle) accountports.AccountRepository {
	return repository
}

func (repository *fakeAccountRepository) CreateAccount(ctx context.Context, account *accountdomain.Account) (*accountdomain.Account, error) {
	repository.accounts[account.ID] = account
	return account, nil
}

func (repository *fakeAccountRepository) UpdateAccount(ctx context.Context, account *accountdomain.Account) (*accountdomain.Account, error) {
	repository.accounts[account.ID] = account
	return account, nil
}

func (repository *fakeAccountRepository) FindAccountByID(ctx context.Context, accountID accountdomain.AccountID, userID userdomain.UserID) (*accountdomain.Account, error) {
	return repository.accounts[accountID], nil
}

func (repository *fakeAccountRepository) FindAccountByIDForUpdate(ctx context.Context, accountID accountdomain.AccountID, userID userdomain.UserID) (*accountdomain.Account, error) {
	return repository.accounts[accountID], nil
}

func (repository *fakeAccountRepository) FindAccountsByUserID(ctx context.Context, input accountports.ListAccountsInput) ([]accountdomain.Account, error) {
	return nil, nil
}

func (repository *fakeAccountRepository) SumActiveAccountBalances(ctx context.Context, userID userdomain.UserID) (int64, error) {
	return 0, nil
}

type fakeCategoryRepository struct {
	categories map[categorydomain.CategoryID]*categorydomain.Category
}

func (repository *fakeCategoryRepository) WithTx(tx databaseports.TxHandle) categoryports.CategoryRepository {
	return repository
}

func (repository *fakeCategoryRepository) CreateCategory(ctx context.Context, category *categorydomain.Category) (*categorydomain.Category, error) {
	repository.categories[category.ID] = category
	return category, nil
}

func (repository *fakeCategoryRepository) CreateCategories(ctx context.Context, categories []categorydomain.Category) ([]categorydomain.Category, error) {
	return categories, nil
}

func (repository *fakeCategoryRepository) UpdateCategory(ctx context.Context, category *categorydomain.Category) (*categorydomain.Category, error) {
	repository.categories[category.ID] = category
	return category, nil
}

func (repository *fakeCategoryRepository) FindCategoryByID(ctx context.Context, categoryID categorydomain.CategoryID, userID userdomain.UserID) (*categorydomain.Category, error) {
	return repository.categories[categoryID], nil
}

func (repository *fakeCategoryRepository) FindCategoriesByUserID(ctx context.Context, input categoryports.ListCategoriesInput) ([]categorydomain.Category, error) {
	return nil, nil
}

func (repository *fakeCategoryRepository) CategoryNameExistsByUserAndType(ctx context.Context, userID userdomain.UserID, categoryType categorydomain.CategoryType, normalizedName string, excludingCategoryID *categorydomain.CategoryID) (bool, error) {
	return false, nil
}
