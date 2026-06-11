package domain

import (
	"strings"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

type TransactionID string

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type TransactionStatus string

const (
	TransactionStatusActive  TransactionStatus = "active"
	TransactionStatusRemoved TransactionStatus = "removed"
)

type TransactionOriginType string

const (
	TransactionOriginTypeManual            TransactionOriginType = "manual"
	TransactionOriginTypePayable           TransactionOriginType = "payable"
	TransactionOriginTypeReceivable        TransactionOriginType = "receivable"
	TransactionOriginTypeCreditCardInvoice TransactionOriginType = "credit_card_invoice"
)

type SettlementStatus string

const (
	SettlementStatusSettled SettlementStatus = "settled"
	SettlementStatusPending SettlementStatus = "pending"
)

type RecurrenceType string

const (
	RecurrenceTypeNone   RecurrenceType = "none"
	RecurrenceTypeFixed  RecurrenceType = "fixed"
	RecurrenceTypeRepeat RecurrenceType = "repeat"
)

type RecurrenceFrequency string

const (
	RecurrenceFrequencyDaily   RecurrenceFrequency = "daily"
	RecurrenceFrequencyWeekly  RecurrenceFrequency = "weekly"
	RecurrenceFrequencyMonthly RecurrenceFrequency = "monthly"
)

type Recurrence struct {
	Frequency  RecurrenceFrequency
	Quantity   *int
	StartsAt   time.Time
	EndsAt     *time.Time
	DayOfMonth *int
}

type Transaction struct {
	ID                   TransactionID
	UserID               userdomain.UserID
	Type                 TransactionType
	Description          string
	Amount               financedomain.Money
	OccurredAt           time.Time
	AccountID            *accountdomain.AccountID
	SourceAccountID      *accountdomain.AccountID
	DestinationAccountID *accountdomain.AccountID
	CategoryID           *categorydomain.CategoryID
	Status               TransactionStatus
	OriginType           TransactionOriginType
	OriginID             *string
	SettlementStatus     SettlementStatus
	SettledAt            *time.Time
	RecurrenceType       RecurrenceType
	Recurrence           *Recurrence
	Note                 string
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type BalanceEffect struct {
	AccountID accountdomain.AccountID
	Amount    financedomain.Money
}

func NewIncome(id TransactionID, userID userdomain.UserID, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, categoryID categorydomain.CategoryID, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string) (Transaction, error) {
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return newTransaction(id, userID, TransactionTypeIncome, description, amount, occurredAt, accountID, nil, nil, &category, settlementStatus, settledAt, recurrenceType, recurrence, note)
}

func NewExpense(id TransactionID, userID userdomain.UserID, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, categoryID categorydomain.CategoryID, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string) (Transaction, error) {
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return newTransaction(id, userID, TransactionTypeExpense, description, amount, occurredAt, accountID, nil, nil, &category, settlementStatus, settledAt, recurrenceType, recurrence, note)
}

func NewTransfer(id TransactionID, userID userdomain.UserID, description string, amount financedomain.Money, occurredAt time.Time, sourceAccountID, destinationAccountID accountdomain.AccountID, note string) (Transaction, error) {
	source := accountdomain.AccountID(strings.TrimSpace(string(sourceAccountID)))
	destination := accountdomain.AccountID(strings.TrimSpace(string(destinationAccountID)))
	return newTransaction(id, userID, TransactionTypeTransfer, description, amount, occurredAt, nil, &source, &destination, nil, SettlementStatusSettled, nil, RecurrenceTypeNone, nil, note)
}

func RehydrateTransaction(id TransactionID, userID userdomain.UserID, transactionType TransactionType, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, sourceAccountID *accountdomain.AccountID, destinationAccountID *accountdomain.AccountID, categoryID *categorydomain.CategoryID, status TransactionStatus, originType TransactionOriginType, originID *string, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string, removedAt *time.Time, createdAt, updatedAt time.Time) (Transaction, error) {
	if originType == "" {
		originType = TransactionOriginTypeManual
	}
	if settlementStatus == "" {
		settlementStatus = SettlementStatusSettled
	}
	if recurrenceType == "" {
		recurrenceType = RecurrenceTypeNone
	}
	transaction := Transaction{
		ID:                   TransactionID(strings.TrimSpace(string(id))),
		UserID:               userdomain.UserID(strings.TrimSpace(string(userID))),
		Type:                 transactionType,
		Description:          strings.TrimSpace(description),
		Amount:               amount,
		OccurredAt:           occurredAt,
		AccountID:            trimAccountIDPointer(accountID),
		SourceAccountID:      trimAccountIDPointer(sourceAccountID),
		DestinationAccountID: trimAccountIDPointer(destinationAccountID),
		CategoryID:           trimCategoryIDPointer(categoryID),
		Status:               status,
		OriginType:           originType,
		OriginID:             trimStringPointer(originID),
		SettlementStatus:     settlementStatus,
		SettledAt:            trimTimePointer(settledAt),
		RecurrenceType:       recurrenceType,
		Recurrence:           trimRecurrencePointer(recurrence),
		Note:                 strings.TrimSpace(note),
		RemovedAt:            removedAt,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
	if err := transaction.validate(); err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func newTransaction(id TransactionID, userID userdomain.UserID, transactionType TransactionType, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, sourceAccountID *accountdomain.AccountID, destinationAccountID *accountdomain.AccountID, categoryID *categorydomain.CategoryID, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string) (Transaction, error) {
	now := time.Now()
	transaction := Transaction{
		ID:                   TransactionID(strings.TrimSpace(string(id))),
		UserID:               userdomain.UserID(strings.TrimSpace(string(userID))),
		Type:                 transactionType,
		Description:          strings.TrimSpace(description),
		Amount:               amount,
		OccurredAt:           occurredAt,
		AccountID:            trimAccountIDPointer(accountID),
		SourceAccountID:      trimAccountIDPointer(sourceAccountID),
		DestinationAccountID: trimAccountIDPointer(destinationAccountID),
		CategoryID:           trimCategoryIDPointer(categoryID),
		Status:               TransactionStatusActive,
		OriginType:           TransactionOriginTypeManual,
		SettlementStatus:     settlementStatus,
		SettledAt:            trimTimePointer(settledAt),
		RecurrenceType:       recurrenceType,
		Recurrence:           trimRecurrencePointer(recurrence),
		Note:                 strings.TrimSpace(note),
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := transaction.validate(); err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func (transaction *Transaction) SetOrigin(originType TransactionOriginType, originID string) error {
	transaction.OriginType = originType
	if strings.TrimSpace(originID) == "" {
		transaction.OriginID = nil
	} else {
		trimmed := strings.TrimSpace(originID)
		transaction.OriginID = &trimmed
	}
	transaction.UpdatedAt = time.Now()
	return transaction.validate()
}

func (transaction Transaction) HasManagedOrigin() bool {
	return transaction.OriginType != TransactionOriginTypeManual
}

func (transaction *Transaction) EditIncome(description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, categoryID categorydomain.CategoryID, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string) error {
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return transaction.edit(TransactionTypeIncome, description, amount, occurredAt, accountID, nil, nil, &category, settlementStatus, settledAt, recurrenceType, recurrence, note)
}

func (transaction *Transaction) EditExpense(description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, categoryID categorydomain.CategoryID, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string) error {
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return transaction.edit(TransactionTypeExpense, description, amount, occurredAt, accountID, nil, nil, &category, settlementStatus, settledAt, recurrenceType, recurrence, note)
}

func (transaction *Transaction) EditTransfer(description string, amount financedomain.Money, occurredAt time.Time, sourceAccountID, destinationAccountID accountdomain.AccountID, note string) error {
	source := accountdomain.AccountID(strings.TrimSpace(string(sourceAccountID)))
	destination := accountdomain.AccountID(strings.TrimSpace(string(destinationAccountID)))
	return transaction.edit(TransactionTypeTransfer, description, amount, occurredAt, nil, &source, &destination, nil, SettlementStatusSettled, nil, RecurrenceTypeNone, nil, note)
}

func (transaction *Transaction) MarkRemoved() error {
	if transaction.Status == TransactionStatusRemoved {
		return nil
	}
	now := time.Now()
	transaction.Status = TransactionStatusRemoved
	transaction.RemovedAt = &now
	transaction.UpdatedAt = now
	return transaction.validate()
}

func (transaction Transaction) BalanceEffects() []BalanceEffect {
	if transaction.Status != TransactionStatusActive {
		return []BalanceEffect{}
	}
	switch transaction.Type {
	case TransactionTypeIncome:
		if transaction.SettlementStatus != SettlementStatusSettled || transaction.AccountID == nil {
			return []BalanceEffect{}
		}
		return []BalanceEffect{{AccountID: *transaction.AccountID, Amount: transaction.Amount}}
	case TransactionTypeExpense:
		if transaction.SettlementStatus != SettlementStatusSettled || transaction.AccountID == nil {
			return []BalanceEffect{}
		}
		return []BalanceEffect{{AccountID: *transaction.AccountID, Amount: transaction.Amount.Neg()}}
	case TransactionTypeTransfer:
		return []BalanceEffect{
			{AccountID: *transaction.SourceAccountID, Amount: transaction.Amount.Neg()},
			{AccountID: *transaction.DestinationAccountID, Amount: transaction.Amount},
		}
	default:
		return []BalanceEffect{}
	}
}

func (transaction Transaction) ReversalEffects() []BalanceEffect {
	effects := transaction.BalanceEffects()
	reversed := make([]BalanceEffect, 0, len(effects))
	for _, effect := range effects {
		reversed = append(reversed, BalanceEffect{AccountID: effect.AccountID, Amount: effect.Amount.Neg()})
	}
	return reversed
}

func (transaction *Transaction) edit(transactionType TransactionType, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, sourceAccountID *accountdomain.AccountID, destinationAccountID *accountdomain.AccountID, categoryID *categorydomain.CategoryID, settlementStatus SettlementStatus, settledAt *time.Time, recurrenceType RecurrenceType, recurrence *Recurrence, note string) error {
	if transaction.Status == TransactionStatusRemoved {
		return ErrTransactionRemoved
	}
	if transaction.Type != transactionType {
		return ErrTransactionInvalidType
	}
	transaction.Description = strings.TrimSpace(description)
	transaction.Amount = amount
	transaction.OccurredAt = occurredAt
	transaction.AccountID = trimAccountIDPointer(accountID)
	transaction.SourceAccountID = trimAccountIDPointer(sourceAccountID)
	transaction.DestinationAccountID = trimAccountIDPointer(destinationAccountID)
	transaction.CategoryID = trimCategoryIDPointer(categoryID)
	transaction.SettlementStatus = settlementStatus
	transaction.SettledAt = trimTimePointer(settledAt)
	transaction.RecurrenceType = recurrenceType
	transaction.Recurrence = trimRecurrencePointer(recurrence)
	transaction.Note = strings.TrimSpace(note)
	transaction.UpdatedAt = time.Now()
	return transaction.validate()
}

func (transaction *Transaction) validate() error {
	if strings.TrimSpace(string(transaction.ID)) == "" {
		return ErrTransactionIDRequired
	}
	if strings.TrimSpace(string(transaction.UserID)) == "" {
		return ErrTransactionUserIDRequired
	}
	if strings.TrimSpace(transaction.Description) == "" {
		return ErrTransactionDescriptionRequired
	}
	if !transaction.Amount.IsPositive() {
		return ErrTransactionAmountInvalid
	}
	if transaction.OccurredAt.IsZero() {
		return ErrTransactionOccurredAtRequired
	}
	if transaction.Status != TransactionStatusActive && transaction.Status != TransactionStatusRemoved {
		return ErrTransactionInvalidStatus
	}
	if transaction.OriginType != TransactionOriginTypeManual &&
		transaction.OriginType != TransactionOriginTypePayable &&
		transaction.OriginType != TransactionOriginTypeReceivable &&
		transaction.OriginType != TransactionOriginTypeCreditCardInvoice {
		return ErrTransactionInvalidType
	}
	if err := transaction.normalizeAndValidateSettlement(); err != nil {
		return err
	}
	if err := transaction.validateRecurrence(); err != nil {
		return err
	}
	switch transaction.Type {
	case TransactionTypeIncome, TransactionTypeExpense:
		if transaction.CategoryID == nil || strings.TrimSpace(string(*transaction.CategoryID)) == "" {
			return ErrTransactionCategoryIDRequired
		}
		if transaction.SourceAccountID != nil || transaction.DestinationAccountID != nil {
			return ErrTransactionInvalidType
		}
	case TransactionTypeTransfer:
		if transaction.SourceAccountID == nil || strings.TrimSpace(string(*transaction.SourceAccountID)) == "" {
			return ErrTransactionSourceAccountIDRequired
		}
		if transaction.DestinationAccountID == nil || strings.TrimSpace(string(*transaction.DestinationAccountID)) == "" {
			return ErrTransactionDestinationAccountIDRequired
		}
		if *transaction.SourceAccountID == *transaction.DestinationAccountID {
			return ErrTransactionTransferAccountsMustBeDifferent
		}
		if transaction.AccountID != nil || transaction.CategoryID != nil {
			return ErrTransactionInvalidType
		}
		if transaction.RecurrenceType != RecurrenceTypeNone || transaction.Recurrence != nil {
			return ErrTransactionInvalidRecurrence
		}
	default:
		return ErrTransactionInvalidType
	}
	if transaction.Status == TransactionStatusRemoved && transaction.RemovedAt == nil {
		return ErrTransactionInvalidStatus
	}
	if transaction.Status == TransactionStatusActive && transaction.RemovedAt != nil {
		return ErrTransactionInvalidStatus
	}
	return nil
}

func (transaction *Transaction) normalizeAndValidateSettlement() error {
	switch transaction.Type {
	case TransactionTypeIncome, TransactionTypeExpense:
		if transaction.SettlementStatus == "" {
			return ErrTransactionSettlementStatusRequired
		}
		switch transaction.SettlementStatus {
		case SettlementStatusSettled:
			if transaction.SettledAt == nil {
				settledAt := transaction.OccurredAt
				transaction.SettledAt = &settledAt
			}
		case SettlementStatusPending:
			if transaction.SettledAt != nil {
				return ErrTransactionInvalidSettledAt
			}
		default:
			return ErrTransactionInvalidSettlementStatus
		}
	case TransactionTypeTransfer:
		transaction.SettlementStatus = SettlementStatusSettled
		transaction.SettledAt = nil
	default:
		return ErrTransactionInvalidType
	}
	return nil
}

func (transaction Transaction) validateRecurrence() error {
	switch transaction.RecurrenceType {
	case "":
		return ErrTransactionInvalidRecurrence
	case RecurrenceTypeNone:
		if transaction.Recurrence != nil {
			return ErrTransactionInvalidRecurrence
		}
	case RecurrenceTypeFixed:
		if transaction.Recurrence == nil || transaction.Recurrence.Quantity != nil {
			return ErrTransactionInvalidRecurrence
		}
		return validateRecurrenceDetails(*transaction.Recurrence)
	case RecurrenceTypeRepeat:
		if transaction.Recurrence == nil || transaction.Recurrence.Quantity == nil || *transaction.Recurrence.Quantity <= 0 {
			return ErrTransactionInvalidRecurrence
		}
		return validateRecurrenceDetails(*transaction.Recurrence)
	default:
		return ErrTransactionInvalidRecurrence
	}
	return nil
}

func validateRecurrenceDetails(recurrence Recurrence) error {
	if recurrence.StartsAt.IsZero() {
		return ErrTransactionInvalidRecurrence
	}
	if recurrence.Frequency != RecurrenceFrequencyDaily &&
		recurrence.Frequency != RecurrenceFrequencyWeekly &&
		recurrence.Frequency != RecurrenceFrequencyMonthly {
		return ErrTransactionInvalidRecurrence
	}
	if recurrence.DayOfMonth != nil && (*recurrence.DayOfMonth < 1 || *recurrence.DayOfMonth > 31) {
		return ErrTransactionInvalidRecurrence
	}
	if recurrence.EndsAt != nil && recurrence.EndsAt.Before(recurrence.StartsAt) {
		return ErrTransactionInvalidRecurrence
	}
	return nil
}

func trimStringPointer(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func trimAccountIDPointer(value *accountdomain.AccountID) *accountdomain.AccountID {
	if value == nil {
		return nil
	}
	trimmed := accountdomain.AccountID(strings.TrimSpace(string(*value)))
	return &trimmed
}

func trimCategoryIDPointer(value *categorydomain.CategoryID) *categorydomain.CategoryID {
	if value == nil {
		return nil
	}
	trimmed := categorydomain.CategoryID(strings.TrimSpace(string(*value)))
	return &trimmed
}

func trimTimePointer(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	trimmed := *value
	if trimmed.IsZero() {
		return nil
	}
	return &trimmed
}

func trimRecurrencePointer(value *Recurrence) *Recurrence {
	if value == nil {
		return nil
	}
	trimmed := *value
	return &trimmed
}
