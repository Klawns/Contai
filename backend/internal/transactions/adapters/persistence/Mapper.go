package persistence

import (
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	"contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

func toTransactionEntity(transaction domain.Transaction) TransactionEntity {
	return TransactionEntity{
		ID:                   string(transaction.ID),
		UserID:               string(transaction.UserID),
		Type:                 string(transaction.Type),
		Description:          transaction.Description,
		Amount:               transaction.Amount.Cents(),
		OccurredAt:           transaction.OccurredAt,
		AccountID:            accountIDToString(transaction.AccountID),
		SourceAccountID:      accountIDToString(transaction.SourceAccountID),
		DestinationAccountID: accountIDToString(transaction.DestinationAccountID),
		CategoryID:           categoryIDToString(transaction.CategoryID),
		Status:               string(transaction.Status),
		OriginType:           string(transaction.OriginType),
		OriginID:             transaction.OriginID,
		SettlementStatus:     string(transaction.SettlementStatus),
		SettledAt:            transaction.SettledAt,
		RecurrenceType:       string(transaction.RecurrenceType),
		RecurrenceFrequency:  recurrenceFrequencyToString(transaction.Recurrence),
		RecurrenceQuantity:   recurrenceQuantity(transaction.Recurrence),
		RecurrenceStartsAt:   recurrenceStartsAt(transaction.Recurrence),
		RecurrenceEndsAt:     recurrenceEndsAt(transaction.Recurrence),
		RecurrenceDayOfMonth: recurrenceDayOfMonth(transaction.Recurrence),
		Note:                 transaction.Note,
		RemovedAt:            transaction.RemovedAt,
		CreatedAt:            transaction.CreatedAt,
		UpdatedAt:            transaction.UpdatedAt,
	}
}

func toDomainTransaction(entity TransactionEntity) (domain.Transaction, error) {
	return domain.RehydrateTransaction(
		domain.TransactionID(entity.ID),
		userdomain.UserID(entity.UserID),
		domain.TransactionType(entity.Type),
		entity.Description,
		financedomain.NewMoney(entity.Amount),
		entity.OccurredAt,
		stringToAccountID(entity.AccountID),
		stringToAccountID(entity.SourceAccountID),
		stringToAccountID(entity.DestinationAccountID),
		stringToCategoryID(entity.CategoryID),
		domain.TransactionStatus(entity.Status),
		domain.TransactionOriginType(entity.OriginType),
		entity.OriginID,
		domain.SettlementStatus(entity.SettlementStatus),
		entity.SettledAt,
		domain.RecurrenceType(entity.RecurrenceType),
		entityToRecurrence(entity),
		entity.Note,
		entity.RemovedAt,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}

func recurrenceFrequencyToString(value *domain.Recurrence) *string {
	if value == nil {
		return nil
	}
	converted := string(value.Frequency)
	return &converted
}

func recurrenceQuantity(value *domain.Recurrence) *int {
	if value == nil {
		return nil
	}
	return value.Quantity
}

func recurrenceStartsAt(value *domain.Recurrence) *time.Time {
	if value == nil {
		return nil
	}
	return &value.StartsAt
}

func recurrenceEndsAt(value *domain.Recurrence) *time.Time {
	if value == nil {
		return nil
	}
	return value.EndsAt
}

func recurrenceDayOfMonth(value *domain.Recurrence) *int {
	if value == nil {
		return nil
	}
	return value.DayOfMonth
}

func entityToRecurrence(entity TransactionEntity) *domain.Recurrence {
	if entity.RecurrenceType == "" || entity.RecurrenceType == string(domain.RecurrenceTypeNone) {
		return nil
	}
	if entity.RecurrenceFrequency == nil || entity.RecurrenceStartsAt == nil {
		return nil
	}
	return &domain.Recurrence{
		Frequency:  domain.RecurrenceFrequency(*entity.RecurrenceFrequency),
		Quantity:   entity.RecurrenceQuantity,
		StartsAt:   *entity.RecurrenceStartsAt,
		EndsAt:     entity.RecurrenceEndsAt,
		DayOfMonth: entity.RecurrenceDayOfMonth,
	}
}

func accountIDToString(value *accountdomain.AccountID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}

func categoryIDToString(value *categorydomain.CategoryID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}

func stringToAccountID(value *string) *accountdomain.AccountID {
	if value == nil {
		return nil
	}
	converted := accountdomain.AccountID(*value)
	return &converted
}

func stringToCategoryID(value *string) *categorydomain.CategoryID {
	if value == nil {
		return nil
	}
	converted := categorydomain.CategoryID(*value)
	return &converted
}
