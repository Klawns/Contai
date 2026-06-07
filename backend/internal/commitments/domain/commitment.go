package domain

import (
	"strings"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type CommitmentID string

type CommitmentType string

const (
	CommitmentTypePayable    CommitmentType = "payable"
	CommitmentTypeReceivable CommitmentType = "receivable"
)

type CommitmentStatus string

const (
	CommitmentStatusPending  CommitmentStatus = "pending"
	CommitmentStatusPaid     CommitmentStatus = "paid"
	CommitmentStatusReceived CommitmentStatus = "received"
	CommitmentStatusCanceled CommitmentStatus = "canceled"
)

type EffectiveStatus string

const (
	EffectiveStatusPending  EffectiveStatus = "pending"
	EffectiveStatusOverdue  EffectiveStatus = "overdue"
	EffectiveStatusPaid     EffectiveStatus = "paid"
	EffectiveStatusReceived EffectiveStatus = "received"
	EffectiveStatusCanceled EffectiveStatus = "canceled"
)

type RecurrenceFrequency string

const (
	RecurrenceFrequencyDaily   RecurrenceFrequency = "daily"
	RecurrenceFrequencyWeekly  RecurrenceFrequency = "weekly"
	RecurrenceFrequencyMonthly RecurrenceFrequency = "monthly"
)

type Recurrence struct {
	Frequency RecurrenceFrequency
	Interval  int
	EndsAt    *time.Time
}

type Commitment struct {
	ID                      CommitmentID
	UserID                  userdomain.UserID
	Type                    CommitmentType
	Description             string
	Amount                  financedomain.Money
	DueAt                   time.Time
	AccountID               accountdomain.AccountID
	CategoryID              categorydomain.CategoryID
	Note                    string
	Status                  CommitmentStatus
	Recurrence              *Recurrence
	SettledAt               *time.Time
	SettlementTransactionID *transactiondomain.TransactionID
	CanceledAt              *time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type EditableFields struct {
	Description string
	Amount      financedomain.Money
	DueAt       time.Time
	AccountID   accountdomain.AccountID
	CategoryID  categorydomain.CategoryID
	Note        string
	Recurrence  *Recurrence
}

func NewPayable(
	id CommitmentID,
	userID userdomain.UserID,
	fields EditableFields,
) (Commitment, error) {
	return newCommitment(id, userID, CommitmentTypePayable, fields)
}

func NewReceivable(
	id CommitmentID,
	userID userdomain.UserID,
	fields EditableFields,
) (Commitment, error) {
	return newCommitment(id, userID, CommitmentTypeReceivable, fields)
}

func RehydrateCommitment(
	id CommitmentID,
	userID userdomain.UserID,
	commitmentType CommitmentType,
	fields EditableFields,
	status CommitmentStatus,
	settledAt *time.Time,
	settlementTransactionID *transactiondomain.TransactionID,
	canceledAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (Commitment, error) {
	commitment := Commitment{
		ID:                      CommitmentID(strings.TrimSpace(string(id))),
		UserID:                  userdomain.UserID(strings.TrimSpace(string(userID))),
		Type:                    commitmentType,
		Status:                  status,
		SettledAt:               settledAt,
		SettlementTransactionID: trimTransactionIDPointer(settlementTransactionID),
		CanceledAt:              canceledAt,
		CreatedAt:               createdAt,
		UpdatedAt:               updatedAt,
	}
	commitment.applyFields(fields)
	if err := commitment.validate(); err != nil {
		return Commitment{}, err
	}
	return commitment, nil
}

func newCommitment(id CommitmentID, userID userdomain.UserID, commitmentType CommitmentType, fields EditableFields) (Commitment, error) {
	now := time.Now()
	commitment := Commitment{
		ID:        CommitmentID(strings.TrimSpace(string(id))),
		UserID:    userdomain.UserID(strings.TrimSpace(string(userID))),
		Type:      commitmentType,
		Status:    CommitmentStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	commitment.applyFields(fields)
	if err := commitment.validate(); err != nil {
		return Commitment{}, err
	}
	return commitment, nil
}

func (commitment *Commitment) Edit(fields EditableFields) error {
	if commitment.Status != CommitmentStatusPending {
		return ErrCommitmentNotPending
	}
	commitment.applyFields(fields)
	commitment.UpdatedAt = time.Now()
	return commitment.validate()
}

func (commitment *Commitment) Cancel() error {
	if commitment.Status != CommitmentStatusPending {
		return ErrCommitmentNotPending
	}
	now := time.Now()
	commitment.Status = CommitmentStatusCanceled
	commitment.CanceledAt = &now
	commitment.UpdatedAt = now
	return commitment.validate()
}

func (commitment *Commitment) MarkPaid(transactionID transactiondomain.TransactionID, settledAt time.Time) error {
	return commitment.settle(CommitmentTypePayable, CommitmentStatusPaid, transactionID, settledAt)
}

func (commitment *Commitment) MarkReceived(transactionID transactiondomain.TransactionID, settledAt time.Time) error {
	return commitment.settle(CommitmentTypeReceivable, CommitmentStatusReceived, transactionID, settledAt)
}

func (commitment Commitment) EffectiveStatus(now time.Time) EffectiveStatus {
	switch commitment.Status {
	case CommitmentStatusPaid:
		return EffectiveStatusPaid
	case CommitmentStatusReceived:
		return EffectiveStatusReceived
	case CommitmentStatusCanceled:
		return EffectiveStatusCanceled
	}
	if commitment.DueAt.Before(now) {
		return EffectiveStatusOverdue
	}
	return EffectiveStatusPending
}

func (commitment *Commitment) settle(
	expectedType CommitmentType,
	status CommitmentStatus,
	transactionID transactiondomain.TransactionID,
	settledAt time.Time,
) error {
	if commitment.Type != expectedType {
		return ErrCommitmentSettlementTypeMismatch
	}
	if commitment.Status != CommitmentStatusPending {
		return ErrCommitmentNotPending
	}
	if strings.TrimSpace(string(transactionID)) == "" || settledAt.IsZero() {
		return ErrCommitmentInvalidStatus
	}
	trimmedTransactionID := transactiondomain.TransactionID(strings.TrimSpace(string(transactionID)))
	now := time.Now()
	commitment.Status = status
	commitment.SettledAt = &settledAt
	commitment.SettlementTransactionID = &trimmedTransactionID
	commitment.UpdatedAt = now
	return commitment.validate()
}

func (commitment *Commitment) applyFields(fields EditableFields) {
	commitment.Description = strings.TrimSpace(fields.Description)
	commitment.Amount = fields.Amount
	commitment.DueAt = fields.DueAt
	commitment.AccountID = accountdomain.AccountID(strings.TrimSpace(string(fields.AccountID)))
	commitment.CategoryID = categorydomain.CategoryID(strings.TrimSpace(string(fields.CategoryID)))
	commitment.Note = strings.TrimSpace(fields.Note)
	commitment.Recurrence = trimRecurrence(fields.Recurrence)
}

func (commitment Commitment) validate() error {
	if strings.TrimSpace(string(commitment.ID)) == "" {
		return ErrCommitmentIDRequired
	}
	if strings.TrimSpace(string(commitment.UserID)) == "" {
		return ErrCommitmentUserIDRequired
	}
	if commitment.Type != CommitmentTypePayable && commitment.Type != CommitmentTypeReceivable {
		return ErrCommitmentInvalidType
	}
	if strings.TrimSpace(commitment.Description) == "" {
		return ErrCommitmentDescriptionRequired
	}
	if !commitment.Amount.IsPositive() {
		return ErrCommitmentAmountInvalid
	}
	if commitment.DueAt.IsZero() {
		return ErrCommitmentDueAtRequired
	}
	if strings.TrimSpace(string(commitment.AccountID)) == "" {
		return ErrCommitmentAccountIDRequired
	}
	if strings.TrimSpace(string(commitment.CategoryID)) == "" {
		return ErrCommitmentCategoryIDRequired
	}
	if err := validateStatus(commitment); err != nil {
		return err
	}
	if err := validateRecurrence(commitment.Recurrence); err != nil {
		return err
	}
	return nil
}

func validateStatus(commitment Commitment) error {
	switch commitment.Status {
	case CommitmentStatusPending:
		if commitment.SettledAt != nil || commitment.SettlementTransactionID != nil || commitment.CanceledAt != nil {
			return ErrCommitmentInvalidStatus
		}
	case CommitmentStatusPaid:
		if commitment.Type != CommitmentTypePayable || commitment.SettledAt == nil || commitment.SettlementTransactionID == nil {
			return ErrCommitmentInvalidStatus
		}
	case CommitmentStatusReceived:
		if commitment.Type != CommitmentTypeReceivable || commitment.SettledAt == nil || commitment.SettlementTransactionID == nil {
			return ErrCommitmentInvalidStatus
		}
	case CommitmentStatusCanceled:
		if commitment.CanceledAt == nil || commitment.SettledAt != nil || commitment.SettlementTransactionID != nil {
			return ErrCommitmentInvalidStatus
		}
	default:
		return ErrCommitmentInvalidStatus
	}
	return nil
}

func validateRecurrence(recurrence *Recurrence) error {
	if recurrence == nil {
		return nil
	}
	if recurrence.Interval <= 0 {
		return ErrCommitmentInvalidRecurrence
	}
	switch recurrence.Frequency {
	case RecurrenceFrequencyDaily, RecurrenceFrequencyWeekly, RecurrenceFrequencyMonthly:
		return nil
	default:
		return ErrCommitmentInvalidRecurrence
	}
}

func trimRecurrence(recurrence *Recurrence) *Recurrence {
	if recurrence == nil {
		return nil
	}
	trimmed := Recurrence{
		Frequency: recurrence.Frequency,
		Interval:  recurrence.Interval,
		EndsAt:    recurrence.EndsAt,
	}
	return &trimmed
}

func trimTransactionIDPointer(value *transactiondomain.TransactionID) *transactiondomain.TransactionID {
	if value == nil {
		return nil
	}
	trimmed := transactiondomain.TransactionID(strings.TrimSpace(string(*value)))
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
