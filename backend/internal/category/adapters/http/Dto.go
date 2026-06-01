package http

import (
	"contai/internal/category/app/ports"
	"contai/internal/category/domain"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type createCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Color string `json:"color" binding:"required"`
	Icon  string `json:"icon" binding:"required"`
}

type updateCategoryRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color" binding:"required"`
	Icon  string `json:"icon" binding:"required"`
}

type categoryResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	Name           string `json:"name"`
	NormalizedName string `json:"normalizedName"`
	Type           string `json:"type"`
	Color          string `json:"color"`
	Icon           string `json:"icon"`
	IsDefault      bool   `json:"isDefault"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

func toCategoryResponse(category ports.CategoryDTO) categoryResponse {
	return categoryResponse{
		ID:             string(category.ID),
		UserID:         string(category.UserID),
		Name:           category.Name,
		NormalizedName: category.NormalizedName,
		Type:           string(category.Type),
		Color:          category.Color,
		Icon:           category.Icon,
		IsDefault:      category.IsDefault,
		Status:         string(category.Status),
		CreatedAt:      category.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:      category.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toCategoryResponses(categories []ports.CategoryDTO) []categoryResponse {
	responses := make([]categoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, toCategoryResponse(category))
	}
	return responses
}

func parseCategoryType(value string) (*domain.CategoryType, error) {
	if value == "" {
		return nil, nil
	}
	categoryType := domain.CategoryType(value)
	if categoryType != domain.CategoryTypeIncome && categoryType != domain.CategoryTypeExpense {
		return nil, domain.ErrCategoryInvalidType
	}
	return &categoryType, nil
}

func parseCategoryStatus(value string) (*domain.CategoryStatus, error) {
	if value == "" {
		return nil, nil
	}
	status := domain.CategoryStatus(value)
	if status != domain.CategoryStatusActive && status != domain.CategoryStatusInactive {
		return nil, domain.ErrCategoryInvalidStatus
	}
	return &status, nil
}
