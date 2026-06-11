package persistence

import "time"

type TransactionEntity struct {
	ID                   string    `gorm:"type:uuid;primaryKey"`
	UserID               string    `gorm:"type:uuid;not null;index:idx_transactions_user_occurred"`
	Type                 string    `gorm:"not null;index:idx_transactions_user_type"`
	Description          string    `gorm:"not null"`
	Amount               int64     `gorm:"not null"`
	OccurredAt           time.Time `gorm:"not null;index:idx_transactions_user_occurred"`
	AccountID            *string   `gorm:"type:uuid;index"`
	SourceAccountID      *string   `gorm:"type:uuid;index"`
	DestinationAccountID *string   `gorm:"type:uuid;index"`
	CategoryID           *string   `gorm:"type:uuid;index"`
	Status               string    `gorm:"not null;index"`
	OriginType           string    `gorm:"not null;default:manual;index"`
	OriginID             *string   `gorm:"index"`
	SettlementStatus     string    `gorm:"not null;default:settled;index"`
	SettledAt            *time.Time
	RecurrenceType       string `gorm:"not null;default:none;index"`
	RecurrenceFrequency  *string
	RecurrenceQuantity   *int
	RecurrenceStartsAt   *time.Time
	RecurrenceEndsAt     *time.Time
	RecurrenceDayOfMonth *int
	Note                 string `gorm:"not null"`
	RemovedAt            *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (TransactionEntity) TableName() string {
	return "transactions"
}
