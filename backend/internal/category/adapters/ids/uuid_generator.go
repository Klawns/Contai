package ids

import (
	"contai/internal/category/domain"

	"github.com/google/uuid"
)

type UUIDCategoryIDGenerator struct{}

func NewUUIDCategoryIDGenerator() UUIDCategoryIDGenerator {
	return UUIDCategoryIDGenerator{}
}

func (generator UUIDCategoryIDGenerator) NewCategoryID() domain.CategoryID {
	return domain.CategoryID(uuid.NewString())
}
