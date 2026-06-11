package ports

import (
	"context"
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/creditcards/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

type CreditCardDTO struct {
	ID              domain.CreditCardID
	UserID          userdomain.UserID
	Name            string
	LinkedAccountID accountdomain.AccountID
	LimitTotal      financedomain.Money
	LimitUsed       financedomain.Money
	LimitAvailable  financedomain.Money
	ClosingDay      int
	DueDay          int
	Status          domain.CreditCardStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PurchaseDTO struct {
	ID                domain.PurchaseID
	UserID            userdomain.UserID
	CardID            domain.CreditCardID
	CategoryID        categorydomain.CategoryID
	Description       string
	TotalAmount       financedomain.Money
	PurchaseDate      time.Time
	PurchaseType      domain.PurchaseType
	InstallmentCount  int
	FirstInvoiceMonth time.Time
	Note              string
	Status            domain.PurchaseStatus
	CanceledAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type InstallmentDTO struct {
	ID             domain.InstallmentID
	UserID         userdomain.UserID
	CardID         domain.CreditCardID
	PurchaseID     domain.PurchaseID
	InvoiceID      domain.InvoiceID
	Number         int
	Amount         financedomain.Money
	Status         domain.PurchaseStatus
	ReferenceMonth time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type InvoiceDTO struct {
	ID                   domain.InvoiceID
	UserID               userdomain.UserID
	CardID               domain.CreditCardID
	ReferenceMonth       time.Time
	ClosingAt            time.Time
	DueAt                time.Time
	Amount               financedomain.Money
	Status               domain.InvoiceStatus
	EffectiveStatus      domain.InvoiceEffectiveStatus
	PaidAt               *time.Time
	PaymentTransactionID *transactiondomain.TransactionID
	Installments         []InstallmentDTO
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CreateCreditCardInput struct {
	UserID          userdomain.UserID
	Name            string
	LinkedAccountID accountdomain.AccountID
	LimitTotal      financedomain.Money
	ClosingDay      int
	DueDay          int
}

type UpdateCreditCardInput struct {
	UserID          userdomain.UserID
	CardID          domain.CreditCardID
	Name            string
	LinkedAccountID accountdomain.AccountID
	LimitTotal      financedomain.Money
	ClosingDay      int
	DueDay          int
	Status          domain.CreditCardStatus
}

type CardIDInput struct {
	UserID userdomain.UserID
	CardID domain.CreditCardID
}

type CreatePurchaseInput struct {
	UserID            userdomain.UserID
	CardID            domain.CreditCardID
	CategoryID        categorydomain.CategoryID
	Description       string
	TotalAmount       financedomain.Money
	PurchaseDate      time.Time
	PurchaseType      domain.PurchaseType
	InstallmentCount  int
	FirstInvoiceMonth time.Time
	Note              string
}

type PurchaseIDInput struct {
	UserID     userdomain.UserID
	PurchaseID domain.PurchaseID
}

type InvoiceIDInput struct {
	UserID    userdomain.UserID
	InvoiceID domain.InvoiceID
}

type PayInvoiceInput struct {
	UserID     userdomain.UserID
	InvoiceID  domain.InvoiceID
	OccurredAt time.Time
	CategoryID categorydomain.CategoryID
	Note       string
}

type CreditCardService interface {
	ListCreditCards(ctx context.Context, userID userdomain.UserID) ([]CreditCardDTO, error)
	CreateCreditCard(ctx context.Context, input CreateCreditCardInput) (CreditCardDTO, error)
	UpdateCreditCard(ctx context.Context, input UpdateCreditCardInput) (CreditCardDTO, error)
	InactivateCreditCard(ctx context.Context, input CardIDInput) (CreditCardDTO, error)
	ListPurchases(ctx context.Context, input CardIDInput) ([]PurchaseDTO, error)
	CreatePurchase(ctx context.Context, input CreatePurchaseInput) (PurchaseDTO, error)
	CancelPurchase(ctx context.Context, input PurchaseIDInput) (PurchaseDTO, error)
	ListInvoices(ctx context.Context, input CardIDInput) ([]InvoiceDTO, error)
	GetInvoice(ctx context.Context, input InvoiceIDInput) (InvoiceDTO, error)
	CloseInvoice(ctx context.Context, input InvoiceIDInput) (InvoiceDTO, error)
	PayInvoice(ctx context.Context, input PayInvoiceInput) (InvoiceDTO, error)
}
