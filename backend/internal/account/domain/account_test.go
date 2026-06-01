package domain

import (
	"errors"
	"testing"

	financedomain "contai/internal/finance/domain"
)

func TestNewAccountSetsCurrentBalanceFromInitialBalance(t *testing.T) {
	account, err := NewAccount("account-id", "user-id", " Checking ", AccountTypeChecking, financedomain.NewMoney(-1500), "bank_1")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account.Name != "Checking" {
		t.Fatalf("expected trimmed name, got %q", account.Name)
	}
	if account.CurrentBalance.Cents() != -1500 {
		t.Fatalf("expected current balance from initial balance, got %d", account.CurrentBalance.Cents())
	}
	if account.Status != AccountStatusActive {
		t.Fatalf("expected active account, got %s", account.Status)
	}
}

func TestNewAccountValidatesRequiredFields(t *testing.T) {
	_, err := NewAccount("", "user-id", "Checking", AccountTypeChecking, 0, "bank")

	if !errors.Is(err, ErrAccountIDRequired) {
		t.Fatalf("expected id required, got %v", err)
	}
}

func TestAccountEditDoesNotChangeBalances(t *testing.T) {
	account := validAccount(t)
	initial := account.InitialBalance
	current := account.CurrentBalance

	err := account.Edit("Savings", AccountTypeSavings, "bank-2")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if account.InitialBalance != initial || account.CurrentBalance != current {
		t.Fatalf("expected balances to be preserved, got initial=%d current=%d", account.InitialBalance.Cents(), account.CurrentBalance.Cents())
	}
}

func TestAccountBalanceMutationsRequirePositiveAmount(t *testing.T) {
	account := validAccount(t)

	if err := account.IncreaseBalance(financedomain.NewMoney(500)); err != nil {
		t.Fatalf("expected increase to succeed, got %v", err)
	}
	if account.CurrentBalance.Cents() != 1500 {
		t.Fatalf("expected increased balance, got %d", account.CurrentBalance.Cents())
	}

	err := account.DecreaseBalance(0)

	if !errors.Is(err, ErrAccountMutationAmountInvalid) {
		t.Fatalf("expected invalid mutation amount, got %v", err)
	}
}

func TestAccountStatusChanges(t *testing.T) {
	account := validAccount(t)

	if err := account.Inactivate(); err != nil {
		t.Fatalf("expected inactivate to succeed, got %v", err)
	}
	if account.Status != AccountStatusInactive {
		t.Fatalf("expected inactive status, got %s", account.Status)
	}

	if err := account.Activate(); err != nil {
		t.Fatalf("expected activate to succeed, got %v", err)
	}
	if account.Status != AccountStatusActive {
		t.Fatalf("expected active status, got %s", account.Status)
	}
}

func validAccount(t *testing.T) Account {
	t.Helper()
	account, err := NewAccount("account-id", "user-id", "Checking", AccountTypeChecking, financedomain.NewMoney(1000), "bank")
	if err != nil {
		t.Fatalf("expected valid account, got %v", err)
	}
	return account
}
