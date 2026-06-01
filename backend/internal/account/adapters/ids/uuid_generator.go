package ids

import (
	"contai/internal/account/domain"

	"github.com/google/uuid"
)

type UUIDAccountIDGenerator struct{}

func NewUUIDAccountIDGenerator() UUIDAccountIDGenerator {
	return UUIDAccountIDGenerator{}
}

func (generator UUIDAccountIDGenerator) NewAccountID() domain.AccountID {
	return domain.AccountID(uuid.NewString())
}
