package services

import (
	"context"
	"testing"
	"time"

	accountports "contai/internal/account/app/ports"
	accountdomain "contai/internal/account/domain"
	categoryports "contai/internal/category/app/ports"
	categorydomain "contai/internal/category/domain"
	"contai/internal/commitments/app/ports"
	"contai/internal/commitments/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	transactionports "contai/internal/transactions/app/ports"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

func TestCommitmentService_CreateDoesNotChangeBalance(t *testing.T) {
	account := validAccount(t, "account-id", financedomain.NewMoney(10000))
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	accountRepository := &fakeAccountRepository{
		accounts: map[accountdomain.AccountID]*accountdomain.Account{"account-id": &account},
	}
	service := NewCommitmentService(
		&fakeCommitmentRepository{},
		&fakeTransactionRepository{},
		accountRepository,
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		fakeCommitmentIDGenerator{},
		fakeTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	_, err := service.CreateCommitment(context.Background(), ports.CreateCommitmentInput{
		UserID:      "user-id",
		Type:        domain.CommitmentTypePayable,
		Description: "Invoice",
		Amount:      financedomain.NewMoney(2500),
		DueAt:       time.Now().Add(24 * time.Hour),
		AccountID:   "account-id",
		CategoryID:  "category-id",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if accountRepository.accounts["account-id"].CurrentBalance.Cents() != 10000 {
		t.Fatalf("expected unchanged balance, got %d", accountRepository.accounts["account-id"].CurrentBalance.Cents())
	}
}

func TestCommitmentService_SettlePayableCreatesManagedExpenseAndUpdatesBalance(t *testing.T) {
	account := validAccount(t, "account-id", financedomain.NewMoney(10000))
	category := validCategory(t, "category-id", categorydomain.CategoryTypeExpense)
	commitment := validPayable(t)
	commitmentRepository := &fakeCommitmentRepository{found: &commitment}
	transactionRepository := &fakeTransactionRepository{}
	accountRepository := &fakeAccountRepository{
		accounts: map[accountdomain.AccountID]*accountdomain.Account{"account-id": &account},
	}
	service := NewCommitmentService(
		commitmentRepository,
		transactionRepository,
		accountRepository,
		&fakeCategoryRepository{categories: map[categorydomain.CategoryID]*categorydomain.Category{"category-id": &category}},
		fakeCommitmentIDGenerator{},
		fakeTransactionIDGenerator{},
		&fakeUnitOfWork{},
	)

	settled, err := service.SettleCommitment(context.Background(), ports.SettleCommitmentInput{
		UserID:       "user-id",
		CommitmentID: "commitment-id",
		Type:         domain.CommitmentTypePayable,
		Amount:       financedomain.NewMoney(2500),
		SettledAt:    time.Now(),
		AccountID:    "account-id",
		CategoryID:   "category-id",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if settled.Status != domain.CommitmentStatusPaid {
		t.Fatalf("expected paid commitment, got %s", settled.Status)
	}
	if accountRepository.accounts["account-id"].CurrentBalance.Cents() != 7500 {
		t.Fatalf("expected balance 7500, got %d", accountRepository.accounts["account-id"].CurrentBalance.Cents())
	}
	if transactionRepository.created == nil || transactionRepository.created.OriginType != transactiondomain.TransactionOriginTypePayable {
		t.Fatalf("expected managed payable transaction, got %#v", transactionRepository.created)
	}
}

func validPayable(t *testing.T) domain.Commitment {
	t.Helper()
	commitment, err := domain.NewPayable("commitment-id", "user-id", domain.EditableFields{
		Description: "Invoice",
		Amount:      financedomain.NewMoney(2500),
		DueAt:       time.Now().Add(24 * time.Hour),
		AccountID:   "account-id",
		CategoryID:  "category-id",
	})
	if err != nil {
		t.Fatalf("expected valid payable, got %v", err)
	}
	return commitment
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

type fakeCommitmentIDGenerator struct{}

func (generator fakeCommitmentIDGenerator) NewCommitmentID() domain.CommitmentID {
	return "commitment-id"
}

type fakeTransactionIDGenerator struct{}

func (generator fakeTransactionIDGenerator) NewTransactionID() transactiondomain.TransactionID {
	return "transaction-id"
}

type fakeCommitmentRepository struct {
	created *domain.Commitment
	updated *domain.Commitment
	found   *domain.Commitment
}

func (repository *fakeCommitmentRepository) WithTx(tx databaseports.TxHandle) ports.CommitmentRepository {
	return repository
}

func (repository *fakeCommitmentRepository) CreateCommitment(ctx context.Context, commitment *domain.Commitment) (*domain.Commitment, error) {
	repository.created = commitment
	return commitment, nil
}

func (repository *fakeCommitmentRepository) UpdateCommitment(ctx context.Context, commitment *domain.Commitment) (*domain.Commitment, error) {
	repository.updated = commitment
	return commitment, nil
}

func (repository *fakeCommitmentRepository) FindCommitmentByIDForUpdate(ctx context.Context, commitmentID domain.CommitmentID, userID userdomain.UserID) (*domain.Commitment, error) {
	return repository.found, nil
}

func (repository *fakeCommitmentRepository) FindCommitmentsByUserID(ctx context.Context, input ports.ListCommitmentsInput) ([]domain.Commitment, error) {
	return nil, nil
}

type fakeTransactionRepository struct {
	created *transactiondomain.Transaction
}

func (repository *fakeTransactionRepository) WithTx(tx databaseports.TxHandle) transactionports.TransactionRepository {
	return repository
}

func (repository *fakeTransactionRepository) CreateTransaction(ctx context.Context, transaction *transactiondomain.Transaction) (*transactiondomain.Transaction, error) {
	repository.created = transaction
	return transaction, nil
}

func (repository *fakeTransactionRepository) UpdateTransaction(ctx context.Context, transaction *transactiondomain.Transaction) (*transactiondomain.Transaction, error) {
	return transaction, nil
}

func (repository *fakeTransactionRepository) FindTransactionByID(ctx context.Context, transactionID transactiondomain.TransactionID, userID userdomain.UserID) (*transactiondomain.Transaction, error) {
	return nil, nil
}

func (repository *fakeTransactionRepository) FindTransactionByIDForUpdate(ctx context.Context, transactionID transactiondomain.TransactionID, userID userdomain.UserID) (*transactiondomain.Transaction, error) {
	return nil, nil
}

func (repository *fakeTransactionRepository) FindTransactionsByUserID(ctx context.Context, input transactionports.ListTransactionsInput) ([]transactiondomain.Transaction, error) {
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
