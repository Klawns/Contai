package domain

import (
	"errors"
	"testing"
	"time"

	financedomain "contai/internal/finance/domain"
)

func TestCommitmentLifecycle(t *testing.T) {
	dueAt := time.Now().Add(24 * time.Hour)
	commitment, err := NewPayable("commitment-id", "user-id", validFields(dueAt))
	if err != nil {
		t.Fatalf("expected valid payable, got %v", err)
	}

	if commitment.Status != CommitmentStatusPending {
		t.Fatalf("expected pending status, got %s", commitment.Status)
	}
	if commitment.EffectiveStatus(time.Now()) != EffectiveStatusPending {
		t.Fatalf("expected effective pending, got %s", commitment.EffectiveStatus(time.Now()))
	}

	settledAt := time.Now()
	if err := commitment.MarkPaid("transaction-id", settledAt); err != nil {
		t.Fatalf("expected payable settlement to succeed, got %v", err)
	}
	if commitment.Status != CommitmentStatusPaid || commitment.SettlementTransactionID == nil {
		t.Fatalf("expected paid commitment with transaction, got %#v", commitment)
	}

	if err := commitment.MarkPaid("other-transaction-id", settledAt); !errors.Is(err, ErrCommitmentNotPending) {
		t.Fatalf("expected duplicate settlement to fail as not pending, got %v", err)
	}
}

func TestCommitmentCancelRequiresPending(t *testing.T) {
	commitment, err := NewReceivable("commitment-id", "user-id", validFields(time.Now()))
	if err != nil {
		t.Fatalf("expected valid receivable, got %v", err)
	}
	if err := commitment.Cancel(); err != nil {
		t.Fatalf("expected cancel to succeed, got %v", err)
	}
	if commitment.Status != CommitmentStatusCanceled || commitment.CanceledAt == nil {
		t.Fatalf("expected canceled commitment, got %#v", commitment)
	}
	if err := commitment.Cancel(); !errors.Is(err, ErrCommitmentNotPending) {
		t.Fatalf("expected second cancel to fail as not pending, got %v", err)
	}
}

func TestCommitmentValidationRejectsInvalidRecurrence(t *testing.T) {
	fields := validFields(time.Now())
	fields.Recurrence = &Recurrence{Frequency: RecurrenceFrequencyMonthly}

	_, err := NewPayable("commitment-id", "user-id", fields)

	if !errors.Is(err, ErrCommitmentInvalidRecurrence) {
		t.Fatalf("expected invalid recurrence error, got %v", err)
	}
}

func validFields(dueAt time.Time) EditableFields {
	return EditableFields{
		Description: "Invoice",
		Amount:      financedomain.NewMoney(1000),
		DueAt:       dueAt,
		AccountID:   "account-id",
		CategoryID:  "category-id",
	}
}
