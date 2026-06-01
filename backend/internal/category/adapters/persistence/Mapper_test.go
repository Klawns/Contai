package persistence

import (
	"testing"

	"contai/internal/category/domain"
)

func TestCategoryMapperRoundTrip(t *testing.T) {
	category, err := domain.NewCategory("category-id", "user-id", "Alimentação", domain.CategoryTypeExpense, "#EA580C", "utensils", true)
	if err != nil {
		t.Fatalf("expected fixture to be valid, got %v", err)
	}

	entity := toCategoryEntity(category)
	mapped, err := toDomainCategory(entity)

	if err != nil {
		t.Fatalf("expected mapper to succeed, got %v", err)
	}
	if mapped.ID != category.ID || mapped.NormalizedName != "alimentacao" || !mapped.IsDefault {
		t.Fatalf("expected mapped category to match, got %#v", mapped)
	}
}
