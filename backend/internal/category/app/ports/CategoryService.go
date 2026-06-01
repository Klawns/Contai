package ports

import (
	"context"
	"time"

	"contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"
)

type CategoryDTO struct {
	ID             domain.CategoryID
	UserID         userdomain.UserID
	Name           string
	NormalizedName string
	Type           domain.CategoryType
	Color          string
	Icon           string
	IsDefault      bool
	Status         domain.CategoryStatus
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateCategoryInput struct {
	UserID userdomain.UserID
	Name   string
	Type   domain.CategoryType
	Color  string
	Icon   string
}

type UpdateCategoryInput struct {
	UserID     userdomain.UserID
	CategoryID domain.CategoryID
	Name       string
	Color      string
	Icon       string
}

type ListCategoriesInput struct {
	UserID userdomain.UserID
	Type   *domain.CategoryType
	Status *domain.CategoryStatus
}

type InactivateCategoryInput struct {
	UserID     userdomain.UserID
	CategoryID domain.CategoryID
}

type CategoryService interface {
	WithTx(tx databaseports.TxHandle) CategoryService
	CreateCategory(ctx context.Context, input CreateCategoryInput) (CategoryDTO, error)
	CreateDefaultCategories(ctx context.Context, userID userdomain.UserID) error
	ListCategories(ctx context.Context, input ListCategoriesInput) ([]CategoryDTO, error)
	UpdateCategory(ctx context.Context, input UpdateCategoryInput) (CategoryDTO, error)
	InactivateCategory(ctx context.Context, input InactivateCategoryInput) error
}
