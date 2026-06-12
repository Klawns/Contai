package persistence

import "time"

type AccountEntity struct {
	ID                      string `gorm:"type:uuid;primaryKey"`
	UserID                  string `gorm:"type:uuid;not null;index:idx_accounts_user_status;index:idx_accounts_user_status_name,priority:1;index:idx_accounts_user_status_dashboard,priority:1"`
	Name                    string `gorm:"not null;index:idx_accounts_user_status_name,priority:3"`
	Type                    string `gorm:"not null"`
	InitialBalance          int64  `gorm:"not null"`
	CurrentBalance          int64  `gorm:"not null"`
	BankIconID              string `gorm:"not null"`
	IncludeInDashboardTotal bool   `gorm:"not null;default:true;index:idx_accounts_user_status_dashboard,priority:3"`
	Status                  string `gorm:"not null;index:idx_accounts_user_status;index:idx_accounts_user_status_name,priority:2;index:idx_accounts_user_status_dashboard,priority:2"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (AccountEntity) TableName() string {
	return "accounts"
}
