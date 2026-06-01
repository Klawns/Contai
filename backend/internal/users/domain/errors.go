package domain

import "errors"

var (
	ErrUserIDRequired           = errors.New("user id is required")
	ErrUserNameRequired         = errors.New("user name is required")
	ErrUserEmailRequired        = errors.New("user email is required")
	ErrUserInvalidEmail         = errors.New("user email is invalid")
	ErrUserPasswordHashRequired = errors.New("user password hash is required")
	ErrUserInvalidStatus        = errors.New("user status is invalid")
	ErrUserEmailAlreadyExists   = errors.New("user email already exists")
	ErrUserPasswordTooWeak      = errors.New("user password is too weak")

	ErrUserAlreadyInactive = errors.New("user is already inactive")
	ErrUserAlreadyActive   = errors.New("user is already active")
	ErrUserInactive        = errors.New("user is inactive")
)
