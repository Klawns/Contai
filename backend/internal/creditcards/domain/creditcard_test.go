package domain

import (
	"errors"
	"testing"
	"time"

	financedomain "contai/internal/finance/domain"
)

func TestSplitInstallmentsPutsRemainderInLastInstallment(t *testing.T) {
	installments := SplitInstallments(financedomain.NewMoney(10000), 3)

	if len(installments) != 3 {
		t.Fatalf("expected three installments, got %#v", installments)
	}
	if installments[0].Cents() != 3333 || installments[1].Cents() != 3333 || installments[2].Cents() != 3334 {
		t.Fatalf("expected remainder in last installment, got %#v", installments)
	}
}

func TestCycleForPurchaseUsesCurrentInvoiceOnClosingDay(t *testing.T) {
	purchaseDate := time.Date(2026, 1, 10, 14, 0, 0, 0, time.UTC)

	referenceMonth, closingAt, dueAt := CycleForPurchase(purchaseDate, 10, 5)

	if referenceMonth.Format("2006-01-02") != "2026-01-01" {
		t.Fatalf("expected january reference month, got %s", referenceMonth)
	}
	if closingAt.Format("2006-01-02") != "2026-01-10" {
		t.Fatalf("expected closing on purchase day, got %s", closingAt)
	}
	if dueAt.Format("2006-01-02") != "2026-02-05" {
		t.Fatalf("expected due date in next month, got %s", dueAt)
	}
}

func TestCycleForPurchaseMovesAfterClosingToNextMonthAndClampsDays(t *testing.T) {
	purchaseDate := time.Date(2026, 1, 31, 12, 0, 0, 0, time.UTC)

	referenceMonth, closingAt, dueAt := CycleForPurchase(purchaseDate, 30, 31)

	if referenceMonth.Format("2006-01-02") != "2026-02-01" {
		t.Fatalf("expected february reference month, got %s", referenceMonth)
	}
	if closingAt.Format("2006-01-02") != "2026-02-28" {
		t.Fatalf("expected clamped february closing, got %s", closingAt)
	}
	if dueAt.Format("2006-01-02") != "2026-02-28" {
		t.Fatalf("expected clamped february due date, got %s", dueAt)
	}
}

func TestNewPurchaseValidatesPurchaseTypeRules(t *testing.T) {
	firstInvoiceMonth := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)

	_, err := NewPurchase("purchase-id", "user-id", "card-id", "category-id", "Notebook", financedomain.NewMoney(120000), time.Now(), PurchaseTypeSingle, 2, firstInvoiceMonth, "")
	if !errors.Is(err, ErrPurchaseInstallmentCountInvalid) {
		t.Fatalf("expected single purchase with multiple installments to be invalid, got %v", err)
	}

	_, err = NewPurchase("purchase-id", "user-id", "card-id", "category-id", "Notebook", financedomain.NewMoney(120000), time.Now(), PurchaseTypeInstallment, 13, firstInvoiceMonth, "")
	if !errors.Is(err, ErrPurchaseInstallmentCountInvalid) {
		t.Fatalf("expected installment purchase above 12x to be invalid, got %v", err)
	}

	_, err = NewPurchase("purchase-id", "user-id", "card-id", "category-id", "Subscription", financedomain.NewMoney(12000), time.Now(), PurchaseTypeFixed, 2, firstInvoiceMonth, "")
	if !errors.Is(err, ErrPurchaseInstallmentCountInvalid) {
		t.Fatalf("expected fixed purchase with installments to be invalid, got %v", err)
	}
}

func TestNewPurchaseStoresFirstInvoiceMonth(t *testing.T) {
	firstInvoiceMonth := time.Date(2026, 2, 17, 10, 0, 0, 0, time.UTC)

	purchase, err := NewPurchase("purchase-id", "user-id", "card-id", "category-id", "Notebook", financedomain.NewMoney(120000), time.Now(), PurchaseTypeInstallment, 12, firstInvoiceMonth, "")

	if err != nil {
		t.Fatalf("expected purchase to be valid, got %v", err)
	}
	if purchase.FirstInvoiceMonth.Format("2006-01-02") != "2026-02-01" {
		t.Fatalf("expected first invoice month normalized to first day, got %s", purchase.FirstInvoiceMonth)
	}
}
