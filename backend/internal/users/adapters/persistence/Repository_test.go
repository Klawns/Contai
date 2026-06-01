package persistence

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"contai/internal/database"
	"contai/internal/users/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestUserRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := database.OpenPostgres(dsn)
	if err != nil {
		t.Fatalf("expected database connection, got %v", err)
	}
	if err := db.AutoMigrate(&UserEntity{}); err != nil {
		t.Fatalf("expected migration to succeed, got %v", err)
	}

	repository := NewUserRepository(db)
	ctx := context.Background()
	id := domain.UserID(uuid.NewString())
	email := fmt.Sprintf("john-%d@example.com", time.Now().UnixNano())

	user, err := domain.NewUser(id, "John Doe", email, "hashed-password")
	if err != nil {
		t.Fatalf("expected no domain error, got %v", err)
	}

	created, err := repository.CreateUser(ctx, &user)
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}
	if created.ID != id {
		t.Fatalf("expected created id %s, got %s", id, created.ID)
	}

	exists, err := repository.EmailExists(ctx, email)
	if err != nil {
		t.Fatalf("expected email exists query to succeed, got %v", err)
	}
	if !exists {
		t.Fatal("expected email to exist")
	}

	foundByID, err := repository.FindUserById(ctx, id)
	if err != nil {
		t.Fatalf("expected find by id to succeed, got %v", err)
	}
	if foundByID == nil || foundByID.ID != id {
		t.Fatalf("expected user by id %s, got %#v", id, foundByID)
	}

	foundByEmail, err := repository.FindUserByEmail(ctx, email)
	if err != nil {
		t.Fatalf("expected find by email to succeed, got %v", err)
	}
	if foundByEmail == nil || foundByEmail.Email != email {
		t.Fatalf("expected user by email %s, got %#v", email, foundByEmail)
	}

	duplicatedUser, err := domain.NewUser(domain.UserID(uuid.NewString()), "John Doe", email, "hashed-password")
	if err != nil {
		t.Fatalf("expected duplicate fixture to be valid, got %v", err)
	}

	_, err = repository.CreateUser(ctx, &duplicatedUser)
	if !errors.Is(err, domain.ErrUserEmailAlreadyExists) {
		t.Fatalf("expected duplicated email domain error, got %v", err)
	}
}

func TestUserRepositoryMapsUniqueEmailViolation(t *testing.T) {
	err := &pgconn.PgError{Code: "23505"}

	if !isUniqueViolation(err) {
		t.Fatal("expected postgres unique violation to be detected")
	}
}
