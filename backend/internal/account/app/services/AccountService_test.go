package services

import (
	"context"
	"errors"
	"testing"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

func TestAccountService_CreateAccountValidatesActiveUser(t *testing.T) {
	service := NewAccountService(&fakeAccountRepository{}, fakeAccountIDGenerator{}, fakeUserValidator{err: domain.ErrAccountUserInactive})

	_, err := service.CreateAccount(context.Background(), createAccountInput())

	if !errors.Is(err, domain.ErrAccountUserInactive) {
		t.Fatalf("expected inactive user error, got %v", err)
	}
}

func TestAccountService_CreateAccountPersistsAccount(t *testing.T) {
	repository := &fakeAccountRepository{}
	service := NewAccountService(repository, fakeAccountIDGenerator{}, fakeUserValidator{})

	account, err := service.CreateAccount(context.Background(), createAccountInput())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account.ID != "account-id" || account.CurrentBalance.Cents() != 1000 {
		t.Fatalf("expected created account dto, got %#v", account)
	}
	if repository.created == nil || repository.created.UserID != "user-id" {
		t.Fatalf("expected persisted account, got %#v", repository.created)
	}
	if !account.IncludeInDashboardTotal {
		t.Fatal("expected created account to include in dashboard total by default")
	}
}

func TestAccountService_CreateAccountPersistsDashboardTotalFalse(t *testing.T) {
	includeInDashboardTotal := false
	repository := &fakeAccountRepository{}
	service := NewAccountService(repository, fakeAccountIDGenerator{}, fakeUserValidator{})
	input := createAccountInput()
	input.IncludeInDashboardTotal = &includeInDashboardTotal

	account, err := service.CreateAccount(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account.IncludeInDashboardTotal {
		t.Fatal("expected created account to preserve false dashboard total flag")
	}
}

func TestAccountService_UpdateAccountPreservesBalances(t *testing.T) {
	existing := validServiceAccount(t)
	repository := &fakeAccountRepository{found: &existing}
	service := NewAccountService(repository, fakeAccountIDGenerator{}, fakeUserValidator{})

	updated, err := service.UpdateAccount(context.Background(), ports.UpdateAccountInput{
		UserID:     "user-id",
		AccountID:  "account-id",
		Name:       "Savings",
		Type:       domain.AccountTypeSavings,
		BankIconID: "bank-2",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "Savings" || updated.InitialBalance.Cents() != 1000 || updated.CurrentBalance.Cents() != 1000 {
		t.Fatalf("expected updated account with preserved balances, got %#v", updated)
	}
	if !updated.IncludeInDashboardTotal {
		t.Fatal("expected omitted dashboard total flag to preserve current value")
	}
}

func TestAccountService_UpdateAccountChangesDashboardTotalFlag(t *testing.T) {
	existing := validServiceAccount(t)
	includeInDashboardTotal := false
	repository := &fakeAccountRepository{found: &existing}
	service := NewAccountService(repository, fakeAccountIDGenerator{}, fakeUserValidator{})

	updated, err := service.UpdateAccount(context.Background(), ports.UpdateAccountInput{
		UserID:                  "user-id",
		AccountID:               "account-id",
		Name:                    "Savings",
		Type:                    domain.AccountTypeSavings,
		BankIconID:              "bank-2",
		IncludeInDashboardTotal: &includeInDashboardTotal,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.IncludeInDashboardTotal {
		t.Fatal("expected dashboard total flag to be updated to false")
	}
}

func TestAccountService_InactivateAccountIsIdempotent(t *testing.T) {
	existing := validServiceAccount(t)
	if err := existing.Inactivate(); err != nil {
		t.Fatalf("expected fixture inactivation, got %v", err)
	}
	repository := &fakeAccountRepository{found: &existing}
	service := NewAccountService(repository, fakeAccountIDGenerator{}, fakeUserValidator{})

	err := service.InactivateAccount(context.Background(), ports.InactivateAccountInput{UserID: "user-id", AccountID: "account-id"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repository.updated != nil {
		t.Fatalf("expected no update for inactive account, got %#v", repository.updated)
	}
}

func TestAccountService_GetTotalBalance(t *testing.T) {
	service := NewAccountService(&fakeAccountRepository{total: 2500}, fakeAccountIDGenerator{}, fakeUserValidator{})

	total, err := service.GetTotalBalance(context.Background(), ports.GetTotalBalanceInput{UserID: "user-id"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if total.Cents() != 2500 {
		t.Fatalf("expected total balance, got %d", total.Cents())
	}
}

func createAccountInput() ports.CreateAccountInput {
	return ports.CreateAccountInput{
		UserID:         "user-id",
		Name:           "Checking",
		Type:           domain.AccountTypeChecking,
		InitialBalance: financedomain.NewMoney(1000),
		BankIconID:     "bank",
	}
}

func validServiceAccount(t *testing.T) domain.Account {
	t.Helper()
	account, err := domain.NewAccount("account-id", "user-id", "Checking", domain.AccountTypeChecking, financedomain.NewMoney(1000), "bank")
	if err != nil {
		t.Fatalf("expected valid account, got %v", err)
	}
	return account
}

type fakeAccountRepository struct {
	created *domain.Account
	updated *domain.Account
	found   *domain.Account
	total   int64
}

func (repository *fakeAccountRepository) WithTx(tx databaseports.TxHandle) ports.AccountRepository {
	return repository
}

func (repository *fakeAccountRepository) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	repository.created = account
	return account, nil
}

func (repository *fakeAccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	repository.updated = account
	return account, nil
}

func (repository *fakeAccountRepository) FindAccountByID(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID) (*domain.Account, error) {
	return repository.found, nil
}

func (repository *fakeAccountRepository) FindAccountByIDForUpdate(ctx context.Context, accountID domain.AccountID, userID userdomain.UserID) (*domain.Account, error) {
	return repository.found, nil
}

func (repository *fakeAccountRepository) FindAccountsByUserID(ctx context.Context, input ports.ListAccountsInput) ([]domain.Account, error) {
	if repository.found == nil {
		return nil, nil
	}
	return []domain.Account{*repository.found}, nil
}

func (repository *fakeAccountRepository) SumActiveAccountBalances(ctx context.Context, userID userdomain.UserID) (int64, error) {
	return repository.total, nil
}

type fakeAccountIDGenerator struct{}

func (generator fakeAccountIDGenerator) NewAccountID() domain.AccountID {
	return "account-id"
}

type fakeUserValidator struct {
	err error
}

func (validator fakeUserValidator) ValidateActiveUser(ctx context.Context, userID userdomain.UserID) error {
	return validator.err
}
