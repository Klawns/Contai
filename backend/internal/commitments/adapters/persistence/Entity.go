package persistence

import "time"

type CommitmentEntity struct {
	ID                      string    `gorm:"type:uuid;primaryKey"`
	UserID                  string    `gorm:"type:uuid;not null;index:idx_commitments_user_due"`
	Type                    string    `gorm:"not null;index:idx_commitments_user_type"`
	Description             string    `gorm:"not null"`
	Amount                  int64     `gorm:"not null"`
	DueAt                   time.Time `gorm:"not null;index:idx_commitments_user_due"`
	AccountID               string    `gorm:"type:uuid;not null;index"`
	CategoryID              string    `gorm:"type:uuid;not null;index"`
	Note                    string    `gorm:"not null"`
	Status                  string    `gorm:"not null;index"`
	RecurrenceFrequency     *string
	RecurrenceInterval      *int
	RecurrenceEndsAt        *time.Time
	SettledAt               *time.Time
	SettlementTransactionID *string `gorm:"type:uuid;index"`
	CanceledAt              *time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

func (CommitmentEntity) TableName() string {
	return "commitments"
}
