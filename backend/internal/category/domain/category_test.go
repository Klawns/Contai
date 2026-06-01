package domain

import (
	"errors"
	"testing"
	"time"

	userdomain "contai/internal/users/domain"
)

func TestNewCategoryNormalizesNameAndValidatesFields(t *testing.T) {
	category, err := NewCategory("category-id", "user-id", "  Saúde  ", CategoryTypeExpense, "#DC2626", "heart-pulse", false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if category.Name != "Saúde" {
		t.Fatalf("expected trimmed name, got %q", category.Name)
	}
	if category.NormalizedName != "saude" {
		t.Fatalf("expected normalized name, got %q", category.NormalizedName)
	}
	if category.Status != CategoryStatusActive {
		t.Fatalf("expected active status, got %s", category.Status)
	}
}

func TestNewCategoryRejectsInvalidTypeColorAndIcon(t *testing.T) {
	tests := []struct {
		name string
		err  error
		fn   func() error
	}{
		{
			name: "type",
			err:  ErrCategoryInvalidType,
			fn: func() error {
				_, err := NewCategory("category-id", "user-id", "Food", "invalid", "#DC2626", "utensils", false)
				return err
			},
		},
		{
			name: "color",
			err:  ErrCategoryInvalidColor,
			fn: func() error {
				_, err := NewCategory("category-id", "user-id", "Food", CategoryTypeExpense, "red", "utensils", false)
				return err
			},
		},
		{
			name: "icon",
			err:  ErrCategoryInvalidIcon,
			fn: func() error {
				_, err := NewCategory("category-id", "user-id", "Food", CategoryTypeExpense, "#DC2626", "Utensils", false)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}

func TestRehydrateCategoryAndEdit(t *testing.T) {
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now().Add(-time.Minute)
	category, err := RehydrateCategory("category-id", "user-id", "Food", "food", CategoryTypeExpense, "#DC2626", "utensils", true, CategoryStatusInactive, createdAt, updatedAt)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := category.Edit(" Alimentação ", "#EA580C", "shopping-basket"); err != nil {
		t.Fatalf("expected edit to succeed, got %v", err)
	}
	if category.Name != "Alimentação" || category.NormalizedName != "alimentacao" {
		t.Fatalf("expected edited normalized name, got %#v", category)
	}
}

func TestCategoryActivationMethods(t *testing.T) {
	category, err := NewCategory("category-id", userdomain.UserID("user-id"), "Food", CategoryTypeExpense, "#DC2626", "utensils", false)
	if err != nil {
		t.Fatalf("expected fixture to be valid, got %v", err)
	}

	if err := category.Inactivate(); err != nil {
		t.Fatalf("expected inactivate to succeed, got %v", err)
	}
	if category.Status != CategoryStatusInactive {
		t.Fatalf("expected inactive status, got %s", category.Status)
	}
	if err := category.Activate(); err != nil {
		t.Fatalf("expected activate to succeed, got %v", err)
	}
	if category.Status != CategoryStatusActive {
		t.Fatalf("expected active status, got %s", category.Status)
	}
}
