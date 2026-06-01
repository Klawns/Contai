package domain

import "errors"

var (
	ErrCategoryIDRequired        = errors.New("category id is required")
	ErrCategoryUserIDRequired    = errors.New("category user id is required")
	ErrCategoryNameRequired      = errors.New("category name is required")
	ErrCategoryInvalidType       = errors.New("category type is invalid")
	ErrCategoryInvalidStatus     = errors.New("category status is invalid")
	ErrCategoryInvalidColor      = errors.New("category color is invalid")
	ErrCategoryInvalidIcon       = errors.New("category icon is invalid")
	ErrCategoryNameAlreadyExists = errors.New("category name already exists")
	ErrCategoryNotFound          = errors.New("category not found")
)
