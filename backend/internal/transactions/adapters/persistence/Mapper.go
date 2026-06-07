package persistence

import (
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
		entity.Note,
		entity.RemovedAt,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
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
