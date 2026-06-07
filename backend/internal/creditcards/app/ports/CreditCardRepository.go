package ports

import (
	"context"
	"time"

	"contai/internal/creditcards/domain"
	databaseports "contai/internal/database/ports"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"
)

type CreditCardRepository interface {
	WithTx(tx databaseports.TxHandle) CreditCardRepository
	CreateCreditCard(ctx context.Context, card *domain.CreditCard) (*domain.CreditCard, error)
	UpdateCreditCard(ctx context.Context, card *domain.CreditCard) (*domain.CreditCard, error)
	FindCreditCardByID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) (*domain.CreditCard, error)
	FindCreditCardByIDForUpdate(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) (*domain.CreditCard, error)
	FindCreditCardsByUserID(ctx context.Context, userID userdomain.UserID) ([]domain.CreditCard, error)
	CreatePurchase(ctx context.Context, purchase *domain.Purchase) (*domain.Purchase, error)
	UpdatePurchase(ctx context.Context, purchase *domain.Purchase) (*domain.Purchase, error)
	FindPurchaseByIDForUpdate(ctx context.Context, purchaseID domain.PurchaseID, userID userdomain.UserID) (*domain.Purchase, error)
	FindPurchasesByCardID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) ([]domain.Purchase, error)
	CreateInstallments(ctx context.Context, installments []domain.Installment) ([]domain.Installment, error)
	UpdateInstallments(ctx context.Context, installments []domain.Installment) error
	FindInstallmentsByPurchaseID(ctx context.Context, purchaseID domain.PurchaseID, userID userdomain.UserID) ([]domain.Installment, error)
	FindInstallmentsByInvoiceID(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) ([]domain.Installment, error)
	FindInvoicesByPurchaseID(ctx context.Context, purchaseID domain.PurchaseID, userID userdomain.UserID) ([]domain.Invoice, error)
	CreateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error)
	FindInvoiceByID(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) (*domain.Invoice, error)
	FindInvoiceByIDForUpdate(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) (*domain.Invoice, error)
	FindInvoiceByCardAndReferenceMonth(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID, referenceMonth time.Time) (*domain.Invoice, error)
	FindInvoicesByCardID(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) ([]domain.Invoice, error)
	SumInvoiceAmount(ctx context.Context, invoiceID domain.InvoiceID, userID userdomain.UserID) (financedomain.Money, error)
	SumLimitUsed(ctx context.Context, cardID domain.CreditCardID, userID userdomain.UserID) (financedomain.Money, error)
}
