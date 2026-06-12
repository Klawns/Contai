package persistence

import "time"

type TransactionEntity struct {
	ID                   string    `gorm:"type:uuid;primaryKey"`
	UserID               string    `gorm:"type:uuid;not null;index:idx_transactions_user_occurred;index:idx_transactions_active_user_occurred_created,priority:1,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_settlement_type_occurred,priority:1,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_category_occurred,priority:1,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_account_occurred,priority:1,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_source_account_occurred,priority:1,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_destination_account_occurred,priority:1,where:status = 'active' AND removed_at IS NULL"`
	Type                 string    `gorm:"not null;index:idx_transactions_user_type;index:idx_transactions_active_user_settlement_type_occurred,priority:3,where:status = 'active' AND removed_at IS NULL"`
	Description          string    `gorm:"not null"`
	Amount               int64     `gorm:"not null"`
	OccurredAt           time.Time `gorm:"not null;index:idx_transactions_user_occurred;index:idx_transactions_active_user_occurred_created,priority:2,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_settlement_type_occurred,priority:4,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_category_occurred,priority:3,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_account_occurred,priority:3,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_source_account_occurred,priority:3,where:status = 'active' AND removed_at IS NULL;index:idx_transactions_active_user_destination_account_occurred,priority:3,where:status = 'active' AND removed_at IS NULL"`
	AccountID            *string   `gorm:"type:uuid;index;index:idx_transactions_active_user_account_occurred,priority:2,where:status = 'active' AND removed_at IS NULL"`
	SourceAccountID      *string   `gorm:"type:uuid;index;index:idx_transactions_active_user_source_account_occurred,priority:2,where:status = 'active' AND removed_at IS NULL"`
	DestinationAccountID *string   `gorm:"type:uuid;index;index:idx_transactions_active_user_destination_account_occurred,priority:2,where:status = 'active' AND removed_at IS NULL"`
	CategoryID           *string   `gorm:"type:uuid;index;index:idx_transactions_active_user_category_occurred,priority:2,where:status = 'active' AND removed_at IS NULL"`
	Status               string    `gorm:"not null;index"`
	OriginType           string    `gorm:"not null;default:manual;index"`
	OriginID             *string   `gorm:"index"`
	SettlementStatus     string    `gorm:"not null;default:settled;index;index:idx_transactions_active_user_settlement_type_occurred,priority:2,where:status = 'active' AND removed_at IS NULL"`
	SettledAt            *time.Time
	RecurrenceType       string `gorm:"not null;default:none;index"`
	RecurrenceFrequency  *string
	RecurrenceQuantity   *int
	RecurrenceStartsAt   *time.Time
	RecurrenceEndsAt     *time.Time
	RecurrenceDayOfMonth *int
	Note                 string `gorm:"not null"`
	RemovedAt            *time.Time
	CreatedAt            time.Time `gorm:"index:idx_transactions_active_user_occurred_created,priority:3,where:status = 'active' AND removed_at IS NULL"`
	UpdatedAt            time.Time
}

func (TransactionEntity) TableName() string {
	return "transactions"
}
