package ids

import (
	"contai/internal/transactions/domain"

	"github.com/google/uuid"
)

type UUIDTransactionIDGenerator struct{}

func NewUUIDTransactionIDGenerator() UUIDTransactionIDGenerator {
	return UUIDTransactionIDGenerator{}
}

func (generator UUIDTransactionIDGenerator) NewTransactionID() domain.TransactionID {
	return domain.TransactionID(uuid.NewString())
}
