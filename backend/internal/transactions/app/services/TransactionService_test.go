package services

import (
	"context"
	"errors"
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
		UserID:      "user-id",
		Description: "Groceries",
		Amount:      financedomain.NewMoney(3500),
		OccurredAt:  time.Now(),
		AccountID:   "account-id",
		CategoryID:  "category-id",
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
		UserID:      "user-id",
		Description: "Salary",
		Amount:      financedomain.NewMoney(10000),
		OccurredAt:  time.Now(),
		AccountID:   "account-id",
		CategoryID:  "category-id",
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
		"account-id",
		"category-id",
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
		UserID:        "user-id",
		TransactionID: "transaction-id",
		Description:   "Invoice",
		Amount:        financedomain.NewMoney(1000),
		OccurredAt:    time.Now(),
		AccountID:     "account-id",
		CategoryID:    "category-id",
	})

	if !errors.Is(err, domain.ErrTransactionManagedOrigin) {
		t.Fatalf("expected managed origin error, got %v", err)
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

type fakeTransactionRepository struct {
	created *domain.Transaction
	updated *domain.Transaction
	found   *domain.Transaction
}

func (repository *fakeTransactionRepository) WithTx(tx databaseports.TxHandle) ports.TransactionRepository {
	return repository
}

func (repository *fakeTransactionRepository) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	repository.created = transaction
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
