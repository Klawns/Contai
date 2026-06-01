package services

import (
	"context"
	"errors"
	"testing"

	"contai/internal/category/app/ports"
	"contai/internal/category/domain"
	databaseports "contai/internal/database/ports"
	userdomain "contai/internal/users/domain"
)

func TestCategoryService_CreateCategoryRejectsDuplicateName(t *testing.T) {
	service := NewCategoryService(&fakeCategoryRepository{nameExists: true}, fakeCategoryIDGenerator{})

	_, err := service.CreateCategory(context.Background(), createCategoryInput("Alimentação"))

	if !errors.Is(err, domain.ErrCategoryNameAlreadyExists) {
		t.Fatalf("expected duplicate name error, got %v", err)
	}
}

func TestCategoryService_CreateCategoryPersistsCategory(t *testing.T) {
	repository := &fakeCategoryRepository{}
	service := NewCategoryService(repository, fakeCategoryIDGenerator{})

	category, err := service.CreateCategory(context.Background(), createCategoryInput("Alimentação"))

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if category.NormalizedName != "alimentacao" {
		t.Fatalf("expected normalized name, got %q", category.NormalizedName)
	}
	if repository.created == nil || repository.created.UserID != "user-id" {
		t.Fatalf("expected category to be persisted, got %#v", repository.created)
	}
}

func TestCategoryService_CreateDefaultCategories(t *testing.T) {
	repository := &fakeCategoryRepository{}
	service := NewCategoryService(repository, fakeCategoryIDGenerator{})

	err := service.CreateDefaultCategories(context.Background(), "user-id")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repository.createdMany) != 12 {
		t.Fatalf("expected 12 default categories, got %d", len(repository.createdMany))
	}
	for _, category := range repository.createdMany {
		if !category.IsDefault {
			t.Fatalf("expected default category, got %#v", category)
		}
	}
}

func TestCategoryService_ListCategoriesValidatesFilters(t *testing.T) {
	service := NewCategoryService(&fakeCategoryRepository{}, fakeCategoryIDGenerator{})
	invalidType := domain.CategoryType("invalid")

	_, err := service.ListCategories(context.Background(), ports.ListCategoriesInput{UserID: "user-id", Type: &invalidType})

	if !errors.Is(err, domain.ErrCategoryInvalidType) {
		t.Fatalf("expected invalid type error, got %v", err)
	}
}

func TestCategoryService_UpdateCategory(t *testing.T) {
	existing := validCategory(t, "category-id", "Moradia")
	repository := &fakeCategoryRepository{found: &existing}
	service := NewCategoryService(repository, fakeCategoryIDGenerator{})

	updated, err := service.UpdateCategory(context.Background(), ports.UpdateCategoryInput{
		UserID:     "user-id",
		CategoryID: "category-id",
		Name:       "Casa",
		Color:      "#2563EB",
		Icon:       "house",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Name != "Casa" || updated.NormalizedName != "casa" {
		t.Fatalf("expected updated category, got %#v", updated)
	}
}

func TestCategoryService_InactivateCategoryIsIdempotent(t *testing.T) {
	existing := validCategory(t, "category-id", "Moradia")
	if err := existing.Inactivate(); err != nil {
		t.Fatalf("expected fixture inactivation to succeed, got %v", err)
	}
	repository := &fakeCategoryRepository{found: &existing}
	service := NewCategoryService(repository, fakeCategoryIDGenerator{})

	err := service.InactivateCategory(context.Background(), ports.InactivateCategoryInput{UserID: "user-id", CategoryID: "category-id"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repository.updated != nil {
		t.Fatalf("expected no update for inactive category, got %#v", repository.updated)
	}
}

func createCategoryInput(name string) ports.CreateCategoryInput {
	return ports.CreateCategoryInput{
		UserID: "user-id",
		Name:   name,
		Type:   domain.CategoryTypeExpense,
		Color:  "#EA580C",
		Icon:   "utensils",
	}
}

func validCategory(t *testing.T, id domain.CategoryID, name string) domain.Category {
	t.Helper()
	category, err := domain.NewCategory(id, "user-id", name, domain.CategoryTypeExpense, "#2563EB", "house", false)
	if err != nil {
		t.Fatalf("expected valid category, got %v", err)
	}
	return category
}

type fakeCategoryRepository struct {
	nameExists  bool
	created     *domain.Category
	createdMany []domain.Category
	updated     *domain.Category
	found       *domain.Category
}

func (repository *fakeCategoryRepository) WithTx(tx databaseports.TxHandle) ports.CategoryRepository {
	return repository
}

func (repository *fakeCategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	repository.created = category
	return category, nil
}

func (repository *fakeCategoryRepository) CreateCategories(ctx context.Context, categories []domain.Category) ([]domain.Category, error) {
	repository.createdMany = categories
	return categories, nil
}

func (repository *fakeCategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	repository.updated = category
	return category, nil
}

func (repository *fakeCategoryRepository) FindCategoryByID(ctx context.Context, categoryID domain.CategoryID, userID userdomain.UserID) (*domain.Category, error) {
	return repository.found, nil
}

func (repository *fakeCategoryRepository) FindCategoriesByUserID(ctx context.Context, input ports.ListCategoriesInput) ([]domain.Category, error) {
	return nil, nil
}

func (repository *fakeCategoryRepository) CategoryNameExistsByUserAndType(ctx context.Context, userID userdomain.UserID, categoryType domain.CategoryType, normalizedName string, excludingCategoryID *domain.CategoryID) (bool, error) {
	return repository.nameExists, nil
}

type fakeCategoryIDGenerator struct{}

func (generator fakeCategoryIDGenerator) NewCategoryID() domain.CategoryID {
	return "category-id"
}
