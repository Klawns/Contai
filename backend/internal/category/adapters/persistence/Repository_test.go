package persistence

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"contai/internal/category/app/ports"
	"contai/internal/category/domain"
	"contai/internal/database"
	userdomain "contai/internal/users/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestCategoryRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := database.OpenPostgres(dsn)
	if err != nil {
		t.Fatalf("expected database connection, got %v", err)
	}
	if err := db.AutoMigrate(&CategoryEntity{}); err != nil {
		t.Fatalf("expected migration to succeed, got %v", err)
	}

	repository := NewCategoryRepository(db)
	ctx := context.Background()
	userID := userdomain.UserID(uuid.NewString())
	name := fmt.Sprintf("Moradia %d", time.Now().UnixNano())

	category, err := domain.NewCategory(domain.CategoryID(uuid.NewString()), userID, name, domain.CategoryTypeExpense, "#2563EB", "house", false)
	if err != nil {
		t.Fatalf("expected no domain error, got %v", err)
	}

	created, err := repository.CreateCategory(ctx, &category)
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	exists, err := repository.CategoryNameExistsByUserAndType(ctx, userID, domain.CategoryTypeExpense, created.NormalizedName, nil)
	if err != nil {
		t.Fatalf("expected exists query to succeed, got %v", err)
	}
	if !exists {
		t.Fatal("expected category name to exist")
	}

	found, err := repository.FindCategoryByID(ctx, created.ID, userID)
	if err != nil {
		t.Fatalf("expected find to succeed, got %v", err)
	}
	if found == nil || found.ID != created.ID {
		t.Fatalf("expected category by id %s, got %#v", created.ID, found)
	}

	active := domain.CategoryStatusActive
	categories, err := repository.FindCategoriesByUserID(ctx, ports.ListCategoriesInput{UserID: userID, Status: &active})
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}
	if len(categories) == 0 {
		t.Fatal("expected listed category")
	}

	if err := created.Inactivate(); err != nil {
		t.Fatalf("expected inactivate to succeed, got %v", err)
	}
	updated, err := repository.UpdateCategory(ctx, created)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if updated.Status != domain.CategoryStatusInactive {
		t.Fatalf("expected inactive category, got %s", updated.Status)
	}

	duplicate, err := domain.NewCategory(domain.CategoryID(uuid.NewString()), userID, name, domain.CategoryTypeExpense, "#2563EB", "house", false)
	if err != nil {
		t.Fatalf("expected duplicate fixture to be valid, got %v", err)
	}
	_, err = repository.CreateCategory(ctx, &duplicate)
	if !errors.Is(err, domain.ErrCategoryNameAlreadyExists) {
		t.Fatalf("expected duplicate category name error, got %v", err)
	}
}

func TestCategoryRepositoryMapsUniqueViolation(t *testing.T) {
	err := &pgconn.PgError{Code: "23505"}

	if !isUniqueViolation(err) {
		t.Fatal("expected postgres unique violation to be detected")
	}
}
