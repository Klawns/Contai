package persistence

import "time"

type CategoryEntity struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	UserID         string `gorm:"type:uuid;not null;uniqueIndex:idx_categories_user_type_name"`
	Name           string `gorm:"not null"`
	NormalizedName string `gorm:"not null;uniqueIndex:idx_categories_user_type_name"`
	Type           string `gorm:"not null;uniqueIndex:idx_categories_user_type_name"`
	Color          string `gorm:"not null"`
	Icon           string `gorm:"not null"`
	IsDefault      bool   `gorm:"not null"`
	Status         string `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (CategoryEntity) TableName() string {
	return "categories"
}
