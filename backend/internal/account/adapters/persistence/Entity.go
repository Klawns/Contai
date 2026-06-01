package persistence

import "time"

type AccountEntity struct {
	ID             string `gorm:"type:uuid;primaryKey"`
	UserID         string `gorm:"type:uuid;not null;index:idx_accounts_user_status"`
	Name           string `gorm:"not null"`
	Type           string `gorm:"not null"`
	InitialBalance int64  `gorm:"not null"`
	CurrentBalance int64  `gorm:"not null"`
	BankIconID     string `gorm:"not null"`
	Status         string `gorm:"not null;index:idx_accounts_user_status"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (AccountEntity) TableName() string {
	return "accounts"
}
