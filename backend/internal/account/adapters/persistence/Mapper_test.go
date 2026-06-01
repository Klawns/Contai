package persistence

import (
	"testing"

	"contai/internal/account/domain"
	financedomain "contai/internal/finance/domain"
)

func TestAccountMapperRoundTrip(t *testing.T) {
	account, err := domain.NewAccount("account-id", "user-id", "Checking", domain.AccountTypeChecking, financedomain.NewMoney(-500), "bank_1")
	if err != nil {
		t.Fatalf("expected valid account, got %v", err)
	}

	entity := toAccountEntity(account)
	mapped, err := toDomainAccount(entity)

	if err != nil {
		t.Fatalf("expected mapper to succeed, got %v", err)
	}
	if mapped.ID != account.ID || mapped.InitialBalance.Cents() != -500 || mapped.CurrentBalance.Cents() != -500 {
		t.Fatalf("expected round trip account, got %#v", mapped)
	}
}
