package domain

import "errors"

var (
	ErrCommitmentIDRequired             = errors.New("commitment id is required")
	ErrCommitmentUserIDRequired         = errors.New("commitment user id is required")
	ErrCommitmentDescriptionRequired    = errors.New("commitment description is required")
	ErrCommitmentAmountInvalid          = errors.New("commitment amount must be positive")
	ErrCommitmentDueAtRequired          = errors.New("commitment due at is required")
	ErrCommitmentAccountIDRequired      = errors.New("commitment account id is required")
	ErrCommitmentCategoryIDRequired     = errors.New("commitment category id is required")
	ErrCommitmentInvalidType            = errors.New("commitment type is invalid")
	ErrCommitmentInvalidStatus          = errors.New("commitment status is invalid")
	ErrCommitmentInvalidRecurrence      = errors.New("commitment recurrence is invalid")
	ErrCommitmentInvalidPagination      = errors.New("commitment pagination is invalid")
	ErrCommitmentNotFound               = errors.New("commitment not found")
	ErrCommitmentNotPending             = errors.New("commitment is not pending")
	ErrCommitmentAccountNotFound        = errors.New("commitment account not found")
	ErrCommitmentCategoryNotFound       = errors.New("commitment category not found")
	ErrCommitmentCategoryTypeMismatch   = errors.New("commitment category type mismatch")
	ErrCommitmentSettlementTypeMismatch = errors.New("commitment settlement type mismatch")
)
