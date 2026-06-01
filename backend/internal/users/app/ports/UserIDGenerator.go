package ports

import "contai/internal/users/domain"

type UserIDGenerator interface {
	NewUserID() domain.UserID
}
