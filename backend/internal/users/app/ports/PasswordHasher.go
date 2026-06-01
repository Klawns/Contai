package ports

import "context"

type PasswordHasher interface {
	HashPassword(ctx context.Context, plainPassword string) (string, error)
	ComparePassword(ctx context.Context, passwordHash string, plainPassword string) error
}
