package persistence

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	"contai/internal/database"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"

	"github.com/google/uuid"
)

func TestAccountRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TEST_DATABASE_DSN is not set")
	}

	db, err := database.OpenPostgres(dsn)
	if err != nil {
		t.Fatalf("expected database connection, got %v", err)
	}
	if err := db.AutoMigrate(&AccountEntity{}); err != nil {
		t.Fatalf("expected migration to succeed, got %v", err)
	}

	repository := NewAccountRepository(db)
	ctx := context.Background()
	userID := userdomain.UserID(uuid.NewString())
	name := fmt.Sprintf("Checking %d", time.Now().UnixNano())

	account, err := domain.NewAccount(domain.AccountID(uuid.NewString()), userID, name, domain.AccountTypeChecking, financedomain.NewMoney(1200), "bank")
	if err != nil {
		t.Fatalf("expected no domain error, got %v", err)
	}

	created, err := repository.CreateAccount(ctx, &account)
	if err != nil {
		t.Fatalf("expected create to succeed, got %v", err)
	}

	found, err := repository.FindAccountByID(ctx, created.ID, userID)
	if err != nil {
		t.Fatalf("expected find to succeed, got %v", err)
	}
	if found == nil || found.ID != created.ID {
		t.Fatalf("expected account by id %s, got %#v", created.ID, found)
	}

	locked, err := repository.FindAccountByIDForUpdate(ctx, created.ID, userID)
	if err != nil {
		t.Fatalf("expected locked find to succeed, got %v", err)
	}
	if locked == nil || locked.ID != created.ID {
		t.Fatalf("expected locked account by id %s, got %#v", created.ID, locked)
	}

	active := domain.AccountStatusActive
	accounts, err := repository.FindAccountsByUserID(ctx, ports.ListAccountsInput{UserID: userID, Status: &active})
	if err != nil {
		t.Fatalf("expected list to succeed, got %v", err)
	}
	if len(accounts) == 0 {
		t.Fatal("expected listed account")
	}

	total, err := repository.SumActiveAccountBalances(ctx, userID)
	if err != nil {
		t.Fatalf("expected sum to succeed, got %v", err)
	}
	if total != 1200 {
		t.Fatalf("expected total 1200, got %d", total)
	}

	if err := created.Inactivate(); err != nil {
		t.Fatalf("expected inactivate to succeed, got %v", err)
	}
	updated, err := repository.UpdateAccount(ctx, created)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}
	if updated.Status != domain.AccountStatusInactive {
		t.Fatalf("expected inactive account, got %s", updated.Status)
	}
}
