package domain

import "errors"

var (
	ErrTransactionIDRequired                      = errors.New("transaction id is required")
	ErrTransactionUserIDRequired                  = errors.New("transaction user id is required")
	ErrTransactionDescriptionRequired             = errors.New("transaction description is required")
	ErrTransactionAmountInvalid                   = errors.New("transaction amount must be positive")
	ErrTransactionOccurredAtRequired              = errors.New("transaction occurred at is required")
	ErrTransactionAccountIDRequired               = errors.New("transaction account id is required")
	ErrTransactionSourceAccountIDRequired         = errors.New("transaction source account id is required")
	ErrTransactionDestinationAccountIDRequired    = errors.New("transaction destination account id is required")
	ErrTransactionTransferAccountsMustBeDifferent = errors.New("transaction transfer accounts must be different")
	ErrTransactionCategoryIDRequired              = errors.New("transaction category id is required")
	ErrTransactionInvalidType                     = errors.New("transaction type is invalid")
	ErrTransactionInvalidStatus                   = errors.New("transaction status is invalid")
	ErrTransactionNotFound                        = errors.New("transaction not found")
	ErrTransactionRemoved                         = errors.New("transaction is removed")
	ErrTransactionAccountNotFound                 = errors.New("transaction account not found")
	ErrTransactionCategoryNotFound                = errors.New("transaction category not found")
	ErrTransactionCategoryTypeMismatch            = errors.New("transaction category type mismatch")
)
