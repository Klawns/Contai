package persistence

import (
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/commitments/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

func toCommitmentEntity(commitment domain.Commitment) CommitmentEntity {
	recurrenceFrequency, recurrenceInterval, recurrenceEndsAt := recurrenceToFields(commitment.Recurrence)
	return CommitmentEntity{
		ID:                      string(commitment.ID),
		UserID:                  string(commitment.UserID),
		Type:                    string(commitment.Type),
		Description:             commitment.Description,
		Amount:                  commitment.Amount.Cents(),
		DueAt:                   commitment.DueAt,
		AccountID:               string(commitment.AccountID),
		CategoryID:              string(commitment.CategoryID),
		Note:                    commitment.Note,
		Status:                  string(commitment.Status),
		RecurrenceFrequency:     recurrenceFrequency,
		RecurrenceInterval:      recurrenceInterval,
		RecurrenceEndsAt:        recurrenceEndsAt,
		SettledAt:               commitment.SettledAt,
		SettlementTransactionID: transactionIDToString(commitment.SettlementTransactionID),
		CanceledAt:              commitment.CanceledAt,
		CreatedAt:               commitment.CreatedAt,
		UpdatedAt:               commitment.UpdatedAt,
	}
}

func toDomainCommitment(entity CommitmentEntity) (domain.Commitment, error) {
	return domain.RehydrateCommitment(
		domain.CommitmentID(entity.ID),
		userdomain.UserID(entity.UserID),
		domain.CommitmentType(entity.Type),
		domain.EditableFields{
			Description: entity.Description,
			Amount:      financedomain.NewMoney(entity.Amount),
			DueAt:       entity.DueAt,
			AccountID:   accountdomain.AccountID(entity.AccountID),
			CategoryID:  categorydomain.CategoryID(entity.CategoryID),
			Note:        entity.Note,
			Recurrence:  recurrenceFromFields(entity),
		},
		domain.CommitmentStatus(entity.Status),
		entity.SettledAt,
		stringToTransactionID(entity.SettlementTransactionID),
		entity.CanceledAt,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}

func recurrenceToFields(recurrence *domain.Recurrence) (*string, *int, *time.Time) {
	if recurrence == nil {
		return nil, nil, nil
	}
	frequency := string(recurrence.Frequency)
	interval := recurrence.Interval
	return &frequency, &interval, recurrence.EndsAt
}

func recurrenceFromFields(entity CommitmentEntity) *domain.Recurrence {
	if entity.RecurrenceFrequency == nil || entity.RecurrenceInterval == nil {
		return nil
	}
	return &domain.Recurrence{
		Frequency: domain.RecurrenceFrequency(*entity.RecurrenceFrequency),
		Interval:  *entity.RecurrenceInterval,
		EndsAt:    entity.RecurrenceEndsAt,
	}
}

func transactionIDToString(value *transactiondomain.TransactionID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}

func stringToTransactionID(value *string) *transactiondomain.TransactionID {
	if value == nil {
		return nil
	}
	converted := transactiondomain.TransactionID(*value)
	return &converted
}
