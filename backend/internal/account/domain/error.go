package domain

import "errors"

var (
	ErrAccountIDRequired            = errors.New("account id is required")
	ErrAccountUserIDRequired        = errors.New("account user id is required")
	ErrAccountNameRequired          = errors.New("account name is required")
	ErrAccountInvalidType           = errors.New("account type is invalid")
	ErrAccountInvalidStatus         = errors.New("account status is invalid")
	ErrAccountBankIconIDRequired    = errors.New("account bank icon id is required")
	ErrAccountInvalidBankIconID     = errors.New("account bank icon id is invalid")
	ErrAccountMutationAmountInvalid = errors.New("account mutation amount must be positive")
	ErrAccountNotFound              = errors.New("account not found")
	ErrAccountUserInactive          = errors.New("account user is inactive")
)
