package persistence

import (
	"testing"
	"time"

	"contai/internal/users/domain"
)

func TestUserMapper(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 2, 4, 5, 6, 0, time.UTC)
	user, err := domain.RehydrateUser(
		"2de2f56f-c0c5-48a2-a1ac-59a581f8da79",
		"John Doe",
		"john@example.com",
		"hashed-password",
		domain.UserStatusActive,
		createdAt,
		updatedAt,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	entity := toUserEntity(user)
	if entity.ID != string(user.ID) {
		t.Fatalf("expected entity id %s, got %s", user.ID, entity.ID)
	}

	mappedUser, err := toDomainUser(entity)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if mappedUser != user {
		t.Fatalf("expected mapped user %#v, got %#v", user, mappedUser)
	}
}
