package domain

import "errors"

var (
	ErrUserNameRequired         = errors.New("user name is required")
	ErrUserEmailRequired        = errors.New("user email is required")
	ErrUserInvalidEmail         = errors.New("user email is invalid")
	ErrUserPasswordHashRequired = errors.New("user password hash is required")

	ErrUserAlreadyInactive = errors.New("user is already inactive")
	ErrUserAlreadyActive   = errors.New("user is already active")
	ErrUserInactive        = errors.New("user is inactive")
)
