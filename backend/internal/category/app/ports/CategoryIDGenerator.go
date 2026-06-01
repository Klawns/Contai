package ports

import "contai/internal/category/domain"

type CategoryIDGenerator interface {
	NewCategoryID() domain.CategoryID
}
