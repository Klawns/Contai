package persistence

import "time"

type CategoryEntity struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	UserID         string `gorm:"type:uuid;not null;uniqueIndex:idx_categories_user_type_name;index:idx_categories_user_type_status_default_name,priority:1"`
	Name           string `gorm:"not null;index:idx_categories_user_type_status_default_name,priority:5"`
	NormalizedName string `gorm:"not null;uniqueIndex:idx_categories_user_type_name"`
	Type           string `gorm:"not null;uniqueIndex:idx_categories_user_type_name;index:idx_categories_user_type_status_default_name,priority:2"`
	Color          string `gorm:"not null"`
	Icon           string `gorm:"not null"`
	IsDefault      bool   `gorm:"not null;index:idx_categories_user_type_status_default_name,priority:4"`
	Status         string `gorm:"not null;index:idx_categories_user_type_status_default_name,priority:3"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (CategoryEntity) TableName() string {
	return "categories"
}
