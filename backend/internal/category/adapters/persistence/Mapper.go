package persistence

import (
	"contai/internal/category/domain"
	userdomain "contai/internal/users/domain"
)

func toCategoryEntity(category domain.Category) CategoryEntity {
	return CategoryEntity{
		ID:             string(category.ID),
		UserID:         string(category.UserID),
		Name:           category.Name,
		NormalizedName: category.NormalizedName,
		Type:           string(category.Type),
		Color:          category.Color,
		Icon:           category.Icon,
		IsDefault:      category.IsDefault,
		Status:         string(category.Status),
		CreatedAt:      category.CreatedAt,
		UpdatedAt:      category.UpdatedAt,
	}
}

func toDomainCategory(entity CategoryEntity) (domain.Category, error) {
	return domain.RehydrateCategory(
		domain.CategoryID(entity.ID),
		userdomain.UserID(entity.UserID),
		entity.Name,
		entity.NormalizedName,
		domain.CategoryType(entity.Type),
		entity.Color,
		entity.Icon,
		entity.IsDefault,
		domain.CategoryStatus(entity.Status),
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}
