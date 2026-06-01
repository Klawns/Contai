package domain

import (
	"errors"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
)

func TestNewIncomeBalanceEffect(t *testing.T) {
	transaction, err := NewIncome("transaction-id", "user-id", "Salary", financedomain.NewMoney(10000), time.Now(), "account-id", "category-id", "")
	if err != nil {
		t.Fatalf("expected income to be valid, got %v", err)
	}

	effects := transaction.BalanceEffects()
	if len(effects) != 1 || effects[0].AccountID != "account-id" || effects[0].Amount.Cents() != 10000 {
		t.Fatalf("expected positive account effect, got %#v", effects)
	}
}

func TestNewExpenseBalanceEffect(t *testing.T) {
	transaction, err := NewExpense("transaction-id", "user-id", "Groceries", financedomain.NewMoney(3500), time.Now(), "account-id", "category-id", "")
	if err != nil {
		t.Fatalf("expected expense to be valid, got %v", err)
	}

	effects := transaction.BalanceEffects()
	if len(effects) != 1 || effects[0].AccountID != "account-id" || effects[0].Amount.Cents() != -3500 {
		t.Fatalf("expected negative account effect, got %#v", effects)
	}
}

func TestNewTransferRequiresDifferentAccounts(t *testing.T) {
	_, err := NewTransfer("transaction-id", "user-id", "Move money", financedomain.NewMoney(1000), time.Now(), "same-account", "same-account", "")

	if !errors.Is(err, ErrTransactionTransferAccountsMustBeDifferent) {
		t.Fatalf("expected different accounts error, got %v", err)
	}
}

func TestMarkRemovedClearsBalanceEffects(t *testing.T) {
	transaction, err := NewTransfer("transaction-id", "user-id", "Move money", financedomain.NewMoney(1000), time.Now(), "source-account", "destination-account", "")
	if err != nil {
		t.Fatalf("expected transfer to be valid, got %v", err)
	}

	if err := transaction.MarkRemoved(); err != nil {
		t.Fatalf("expected soft delete to succeed, got %v", err)
	}

	if transaction.Status != TransactionStatusRemoved || transaction.RemovedAt == nil {
		t.Fatalf("expected removed transaction, got %#v", transaction)
	}
	if effects := transaction.BalanceEffects(); len(effects) != 0 {
		t.Fatalf("expected removed transaction to have no active effects, got %#v", effects)
	}
}

func TestRehydrateTransactionValidatesShape(t *testing.T) {
	accountID := accountdomain.AccountID("account-id")
	categoryID := categorydomain.CategoryID("category-id")

	_, err := RehydrateTransaction("transaction-id", "user-id", TransactionTypeTransfer, "Bad transfer", financedomain.NewMoney(1000), time.Now(), &accountID, nil, nil, &categoryID, TransactionStatusActive, "", nil, time.Now(), time.Now())

	if !errors.Is(err, ErrTransactionSourceAccountIDRequired) {
		t.Fatalf("expected transfer source account error, got %v", err)
	}
}
