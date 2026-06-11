package persistence

import "time"

type CreditCardEntity struct {
	ID              string `gorm:"primaryKey;type:uuid"`
	UserID          string `gorm:"type:uuid;not null;index"`
	Name            string `gorm:"not null"`
	LinkedAccountID string `gorm:"type:uuid;not null;index"`
	LimitTotal      int64  `gorm:"not null"`
	ClosingDay      int    `gorm:"not null"`
	DueDay          int    `gorm:"not null"`
	Status          string `gorm:"not null;index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (CreditCardEntity) TableName() string {
	return "credit_cards"
}

type CardPurchaseEntity struct {
	ID                string `gorm:"primaryKey;type:uuid"`
	UserID            string `gorm:"type:uuid;not null;index"`
	CardID            string `gorm:"type:uuid;not null;index"`
	CategoryID        string `gorm:"type:uuid;not null;index"`
	Description       string `gorm:"not null"`
	TotalAmount       int64  `gorm:"not null"`
	PurchaseDate      time.Time
	PurchaseType      string `gorm:"not null;default:single"`
	InstallmentCount  int
	FirstInvoiceMonth time.Time `gorm:"index"`
	Note              string
	Status            string `gorm:"not null;index"`
	CanceledAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (CardPurchaseEntity) TableName() string {
	return "card_purchases"
}

type CardInstallmentEntity struct {
	ID             string `gorm:"primaryKey;type:uuid"`
	UserID         string `gorm:"type:uuid;not null;index"`
	CardID         string `gorm:"type:uuid;not null;index"`
	PurchaseID     string `gorm:"type:uuid;not null;index"`
	InvoiceID      string `gorm:"type:uuid;not null;index"`
	Number         int    `gorm:"not null"`
	Amount         int64  `gorm:"not null"`
	Status         string `gorm:"not null;index"`
	ReferenceMonth time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (CardInstallmentEntity) TableName() string {
	return "card_installments"
}

type CardInvoiceEntity struct {
	ID                   string    `gorm:"primaryKey;type:uuid"`
	UserID               string    `gorm:"type:uuid;not null;index"`
	CardID               string    `gorm:"type:uuid;not null;index"`
	ReferenceMonth       time.Time `gorm:"not null;index"`
	ClosingAt            time.Time `gorm:"not null"`
	DueAt                time.Time `gorm:"not null;index"`
	Status               string    `gorm:"not null;index"`
	PaidAt               *time.Time
	PaymentTransactionID *string `gorm:"type:uuid"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (CardInvoiceEntity) TableName() string {
	return "card_invoices"
}
