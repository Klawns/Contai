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
	ID                string    `gorm:"primaryKey;type:uuid"`
	UserID            string    `gorm:"type:uuid;not null;index;index:idx_card_purchases_user_card_date_created,priority:1;index:idx_card_purchases_user_status_date_category,priority:1"`
	CardID            string    `gorm:"type:uuid;not null;index;index:idx_card_purchases_user_card_date_created,priority:2"`
	CategoryID        string    `gorm:"type:uuid;not null;index;index:idx_card_purchases_user_status_date_category,priority:4"`
	Description       string    `gorm:"not null"`
	TotalAmount       int64     `gorm:"not null"`
	PurchaseDate      time.Time `gorm:"index:idx_card_purchases_user_card_date_created,priority:3;index:idx_card_purchases_user_status_date_category,priority:3"`
	PurchaseType      string    `gorm:"not null;default:single"`
	InstallmentCount  int
	FirstInvoiceMonth time.Time `gorm:"index"`
	Note              string
	Status            string `gorm:"not null;index;index:idx_card_purchases_user_status_date_category,priority:2"`
	CanceledAt        *time.Time
	CreatedAt         time.Time `gorm:"index:idx_card_purchases_user_card_date_created,priority:4"`
	UpdatedAt         time.Time
}

func (CardPurchaseEntity) TableName() string {
	return "card_purchases"
}

type CardInstallmentEntity struct {
	ID             string    `gorm:"primaryKey;type:uuid"`
	UserID         string    `gorm:"type:uuid;not null;index;index:idx_card_installments_user_purchase_number,priority:1;index:idx_card_installments_user_invoice_month_number,priority:1;index:idx_card_installments_user_invoice_status,priority:1;index:idx_card_installments_user_card_status_invoice,priority:1"`
	CardID         string    `gorm:"type:uuid;not null;index;index:idx_card_installments_user_card_status_invoice,priority:2"`
	PurchaseID     string    `gorm:"type:uuid;not null;index;index:idx_card_installments_user_purchase_number,priority:2"`
	InvoiceID      string    `gorm:"type:uuid;not null;index;index:idx_card_installments_user_invoice_month_number,priority:2;index:idx_card_installments_user_invoice_status,priority:2;index:idx_card_installments_user_card_status_invoice,priority:4"`
	Number         int       `gorm:"not null;index:idx_card_installments_user_purchase_number,priority:3;index:idx_card_installments_user_invoice_month_number,priority:4"`
	Amount         int64     `gorm:"not null"`
	Status         string    `gorm:"not null;index;index:idx_card_installments_user_invoice_status,priority:3;index:idx_card_installments_user_card_status_invoice,priority:3"`
	ReferenceMonth time.Time `gorm:"index:idx_card_installments_user_invoice_month_number,priority:3"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (CardInstallmentEntity) TableName() string {
	return "card_installments"
}

type CardInvoiceEntity struct {
	ID                   string    `gorm:"primaryKey;type:uuid"`
	UserID               string    `gorm:"type:uuid;not null;index;index:idx_card_invoices_user_card_reference_month,priority:1"`
	CardID               string    `gorm:"type:uuid;not null;index;index:idx_card_invoices_user_card_reference_month,priority:2"`
	ReferenceMonth       time.Time `gorm:"not null;index;index:idx_card_invoices_user_card_reference_month,priority:3"`
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
