package users

import (
	"context"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	userports "contai/internal/users/app/ports"
	userdomain "contai/internal/users/domain"
)

var _ ports.UserValidator = ActiveUserValidator{}

type ActiveUserValidator struct {
	userRepository userports.UserRepository
}

func NewActiveUserValidator(userRepository userports.UserRepository) ActiveUserValidator {
	return ActiveUserValidator{userRepository: userRepository}
}

func (validator ActiveUserValidator) ValidateActiveUser(ctx context.Context, userID userdomain.UserID) error {
	if userID == "" {
		return domain.ErrAccountUserIDRequired
	}

	user, err := validator.userRepository.FindUserById(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil || user.Status != userdomain.UserStatusActive {
		return domain.ErrAccountUserInactive
	}

	return nil
}
