package ids

import (
	"contai/internal/users/domain"

	"github.com/google/uuid"
)

type UUIDUserIDGenerator struct{}

func NewUUIDUserIDGenerator() UUIDUserIDGenerator {
	return UUIDUserIDGenerator{}
}

func (UUIDUserIDGenerator) NewUserID() domain.UserID {
	return domain.UserID(uuid.NewString())
}
