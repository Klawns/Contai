package ports

import (
	"context"

	"contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"
)

type CategoryRepository interface {
	WithTx(tx databaseports.TxHandle) CategoryRepository
	CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	CreateCategories(ctx context.Context, categories []domain.Category) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error)
	FindCategoryByID(ctx context.Context, categoryID domain.CategoryID, userID userdomain.UserID) (*domain.Category, error)
	FindCategoriesByUserID(ctx context.Context, input ListCategoriesInput) ([]domain.Category, error)
	CategoryNameExistsByUserAndType(ctx context.Context, userID userdomain.UserID, categoryType domain.CategoryType, normalizedName string, excludingCategoryID *domain.CategoryID) (bool, error)
}
