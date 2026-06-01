package persistence

import "time"

type UserEntity struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Status       string `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (UserEntity) TableName() string {
	return "users"
}
