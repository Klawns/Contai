package ports

import (
	"context"

	userdomain "contai/internal/users/domain"
)

type UserValidator interface {
	ValidateActiveUser(ctx context.Context, userID userdomain.UserID) error
}
