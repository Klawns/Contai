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
	Note                 string
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type BalanceEffect struct {
	AccountID accountdomain.AccountID
	Amount    financedomain.Money
}

func NewIncome(id TransactionID, userID userdomain.UserID, description string, amount financedomain.Money, occurredAt time.Time, accountID accountdomain.AccountID, categoryID categorydomain.CategoryID, note string) (Transaction, error) {
	account := accountdomain.AccountID(strings.TrimSpace(string(accountID)))
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return newTransaction(id, userID, TransactionTypeIncome, description, amount, occurredAt, &account, nil, nil, &category, note)
}

func NewExpense(id TransactionID, userID userdomain.UserID, description string, amount financedomain.Money, occurredAt time.Time, accountID accountdomain.AccountID, categoryID categorydomain.CategoryID, note string) (Transaction, error) {
	account := accountdomain.AccountID(strings.TrimSpace(string(accountID)))
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return newTransaction(id, userID, TransactionTypeExpense, description, amount, occurredAt, &account, nil, nil, &category, note)
}

func NewTransfer(id TransactionID, userID userdomain.UserID, description string, amount financedomain.Money, occurredAt time.Time, sourceAccountID, destinationAccountID accountdomain.AccountID, note string) (Transaction, error) {
	source := accountdomain.AccountID(strings.TrimSpace(string(sourceAccountID)))
	destination := accountdomain.AccountID(strings.TrimSpace(string(destinationAccountID)))
	return newTransaction(id, userID, TransactionTypeTransfer, description, amount, occurredAt, nil, &source, &destination, nil, note)
}

func RehydrateTransaction(id TransactionID, userID userdomain.UserID, transactionType TransactionType, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, sourceAccountID *accountdomain.AccountID, destinationAccountID *accountdomain.AccountID, categoryID *categorydomain.CategoryID, status TransactionStatus, note string, removedAt *time.Time, createdAt, updatedAt time.Time) (Transaction, error) {
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

func newTransaction(id TransactionID, userID userdomain.UserID, transactionType TransactionType, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, sourceAccountID *accountdomain.AccountID, destinationAccountID *accountdomain.AccountID, categoryID *categorydomain.CategoryID, note string) (Transaction, error) {
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
		Note:                 strings.TrimSpace(note),
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if err := transaction.validate(); err != nil {
		return Transaction{}, err
	}
	return transaction, nil
}

func (transaction *Transaction) EditIncome(description string, amount financedomain.Money, occurredAt time.Time, accountID accountdomain.AccountID, categoryID categorydomain.CategoryID, note string) error {
	account := accountdomain.AccountID(strings.TrimSpace(string(accountID)))
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return transaction.edit(TransactionTypeIncome, description, amount, occurredAt, &account, nil, nil, &category, note)
}

func (transaction *Transaction) EditExpense(description string, amount financedomain.Money, occurredAt time.Time, accountID accountdomain.AccountID, categoryID categorydomain.CategoryID, note string) error {
	account := accountdomain.AccountID(strings.TrimSpace(string(accountID)))
	category := categorydomain.CategoryID(strings.TrimSpace(string(categoryID)))
	return transaction.edit(TransactionTypeExpense, description, amount, occurredAt, &account, nil, nil, &category, note)
}

func (transaction *Transaction) EditTransfer(description string, amount financedomain.Money, occurredAt time.Time, sourceAccountID, destinationAccountID accountdomain.AccountID, note string) error {
	source := accountdomain.AccountID(strings.TrimSpace(string(sourceAccountID)))
	destination := accountdomain.AccountID(strings.TrimSpace(string(destinationAccountID)))
	return transaction.edit(TransactionTypeTransfer, description, amount, occurredAt, nil, &source, &destination, nil, note)
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
		return []BalanceEffect{{AccountID: *transaction.AccountID, Amount: transaction.Amount}}
	case TransactionTypeExpense:
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

func (transaction *Transaction) edit(transactionType TransactionType, description string, amount financedomain.Money, occurredAt time.Time, accountID *accountdomain.AccountID, sourceAccountID *accountdomain.AccountID, destinationAccountID *accountdomain.AccountID, categoryID *categorydomain.CategoryID, note string) error {
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
	transaction.Note = strings.TrimSpace(note)
	transaction.UpdatedAt = time.Now()
	return transaction.validate()
}

func (transaction Transaction) validate() error {
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
	switch transaction.Type {
	case TransactionTypeIncome, TransactionTypeExpense:
		if transaction.AccountID == nil || strings.TrimSpace(string(*transaction.AccountID)) == "" {
			return ErrTransactionAccountIDRequired
		}
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
