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
	occurredAt := time.Now()
	transaction, err := NewIncome("transaction-id", "user-id", "Salary", financedomain.NewMoney(10000), occurredAt, accountIDPtr("account-id"), "category-id", SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected income to be valid, got %v", err)
	}

	effects := transaction.BalanceEffects()
	if len(effects) != 1 || effects[0].AccountID != "account-id" || effects[0].Amount.Cents() != 10000 {
		t.Fatalf("expected positive account effect, got %#v", effects)
	}
}

func TestNewExpenseBalanceEffect(t *testing.T) {
	transaction, err := NewExpense("transaction-id", "user-id", "Groceries", financedomain.NewMoney(3500), time.Now(), accountIDPtr("account-id"), "category-id", SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected expense to be valid, got %v", err)
	}

	effects := transaction.BalanceEffects()
	if len(effects) != 1 || effects[0].AccountID != "account-id" || effects[0].Amount.Cents() != -3500 {
		t.Fatalf("expected negative account effect, got %#v", effects)
	}
}

func TestNewExpenseAllowsMissingAccount(t *testing.T) {
	transaction, err := NewExpense("transaction-id", "user-id", "Groceries", financedomain.NewMoney(3500), time.Now(), nil, "category-id", SettlementStatusPending, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected expense without account to be valid, got %v", err)
	}
	if transaction.AccountID != nil {
		t.Fatalf("expected nil account, got %#v", transaction.AccountID)
	}
}

func TestPendingTransactionHasNoBalanceEffect(t *testing.T) {
	transaction, err := NewIncome("transaction-id", "user-id", "Salary", financedomain.NewMoney(10000), time.Now(), accountIDPtr("account-id"), "category-id", SettlementStatusPending, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected pending income to be valid, got %v", err)
	}
	if effects := transaction.BalanceEffects(); len(effects) != 0 {
		t.Fatalf("expected pending transaction to have no balance effects, got %#v", effects)
	}
}

func TestSettledTransactionWithoutAccountHasNoBalanceEffect(t *testing.T) {
	transaction, err := NewExpense("transaction-id", "user-id", "Groceries", financedomain.NewMoney(3500), time.Now(), nil, "category-id", SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected settled expense without account to be valid, got %v", err)
	}
	if effects := transaction.BalanceEffects(); len(effects) != 0 {
		t.Fatalf("expected transaction without account to have no balance effects, got %#v", effects)
	}
}

func TestInvalidRecurrence(t *testing.T) {
	quantity := 0
	_, err := NewExpense("transaction-id", "user-id", "Groceries", financedomain.NewMoney(3500), time.Now(), nil, "category-id", SettlementStatusPending, nil, RecurrenceTypeRepeat, &Recurrence{
		Frequency: RecurrenceFrequencyMonthly,
		Quantity:  &quantity,
		StartsAt:  time.Now(),
	}, "")
	if !errors.Is(err, ErrTransactionInvalidRecurrence) {
		t.Fatalf("expected invalid recurrence error, got %v", err)
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

	_, err := RehydrateTransaction("transaction-id", "user-id", TransactionTypeTransfer, "Bad transfer", financedomain.NewMoney(1000), time.Now(), &accountID, nil, nil, &categoryID, TransactionStatusActive, TransactionOriginTypeManual, nil, SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "", nil, time.Now(), time.Now())

	if !errors.Is(err, ErrTransactionSourceAccountIDRequired) {
		t.Fatalf("expected transfer source account error, got %v", err)
	}
}

func TestCalculateTransactionTotals(t *testing.T) {
	income, err := NewIncome("income-id", "user-id", "Salary", financedomain.NewMoney(10000), time.Now(), accountIDPtr("account-id"), "category-id", SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected valid income, got %v", err)
	}
	expense, err := NewExpense("expense-id", "user-id", "Market", financedomain.NewMoney(2500), time.Now(), accountIDPtr("account-id"), "category-id", SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected valid expense, got %v", err)
	}
	transfer, err := NewTransfer("transfer-id", "user-id", "Move money", financedomain.NewMoney(1500), time.Now(), "source-account", "destination-account", "")
	if err != nil {
		t.Fatalf("expected valid transfer, got %v", err)
	}
	removed, err := NewExpense("removed-id", "user-id", "Removed", financedomain.NewMoney(9999), time.Now(), accountIDPtr("account-id"), "category-id", SettlementStatusSettled, nil, RecurrenceTypeNone, nil, "")
	if err != nil {
		t.Fatalf("expected valid removed expense, got %v", err)
	}
	if err := removed.MarkRemoved(); err != nil {
		t.Fatalf("expected remove to succeed, got %v", err)
	}

	totals := CalculateTransactionTotals([]Transaction{income, expense, transfer, removed})

	if totals.IncomeTotal.Cents() != 10000 {
		t.Fatalf("expected income total 10000, got %d", totals.IncomeTotal.Cents())
	}
	if totals.ExpenseTotal.Cents() != 2500 {
		t.Fatalf("expected expense total 2500, got %d", totals.ExpenseTotal.Cents())
	}
	if totals.TransferInTotal.Cents() != 1500 || totals.TransferOutTotal.Cents() != 1500 {
		t.Fatalf("expected transfer totals 1500/1500, got %d/%d", totals.TransferInTotal.Cents(), totals.TransferOutTotal.Cents())
	}
}

func accountIDPtr(value accountdomain.AccountID) *accountdomain.AccountID {
	return &value
}

func TestCalculateTransactionTotalsEmptyList(t *testing.T) {
	totals := CalculateTransactionTotals(nil)

	if totals.IncomeTotal.Cents() != 0 || totals.ExpenseTotal.Cents() != 0 || totals.TransferInTotal.Cents() != 0 || totals.TransferOutTotal.Cents() != 0 {
		t.Fatalf("expected empty totals, got %#v", totals)
	}
}
