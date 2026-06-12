package persistence

import "time"

type CommitmentEntity struct {
	ID                      string    `gorm:"type:uuid;primaryKey"`
	UserID                  string    `gorm:"type:uuid;not null;index:idx_commitments_user_due;index:idx_commitments_user_type_due_created,priority:1;index:idx_commitments_user_type_status_due,priority:1;index:idx_commitments_user_type_account_due,priority:1;index:idx_commitments_user_type_category_due,priority:1"`
	Type                    string    `gorm:"not null;index:idx_commitments_user_type;index:idx_commitments_user_type_due_created,priority:2;index:idx_commitments_user_type_status_due,priority:2;index:idx_commitments_user_type_account_due,priority:2;index:idx_commitments_user_type_category_due,priority:2"`
	Description             string    `gorm:"not null"`
	Amount                  int64     `gorm:"not null"`
	DueAt                   time.Time `gorm:"not null;index:idx_commitments_user_due;index:idx_commitments_user_type_due_created,priority:3;index:idx_commitments_user_type_status_due,priority:4;index:idx_commitments_user_type_account_due,priority:4;index:idx_commitments_user_type_category_due,priority:4"`
	AccountID               string    `gorm:"type:uuid;not null;index;index:idx_commitments_user_type_account_due,priority:3"`
	CategoryID              string    `gorm:"type:uuid;not null;index;index:idx_commitments_user_type_category_due,priority:3"`
	Note                    string    `gorm:"not null"`
	Status                  string    `gorm:"not null;index;index:idx_commitments_user_type_status_due,priority:3"`
	RecurrenceFrequency     *string
	RecurrenceInterval      *int
	RecurrenceEndsAt        *time.Time
	SettledAt               *time.Time
	SettlementTransactionID *string `gorm:"type:uuid;index"`
	CanceledAt              *time.Time
	CreatedAt               time.Time `gorm:"index:idx_commitments_user_type_due_created,priority:4"`
	UpdatedAt               time.Time
}

func (CommitmentEntity) TableName() string {
	return "commitments"
}
