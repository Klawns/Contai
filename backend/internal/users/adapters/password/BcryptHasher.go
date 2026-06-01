package password

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{cost: bcrypt.DefaultCost}
}

func (hasher BcryptHasher) HashPassword(ctx context.Context, plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), hasher.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (hasher BcryptHasher) ComparePassword(ctx context.Context, passwordHash string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
}
