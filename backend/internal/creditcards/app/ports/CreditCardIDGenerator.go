package ports

import "contai/internal/creditcards/domain"

type CreditCardIDGenerator interface {
	NewCreditCardID() domain.CreditCardID
	NewPurchaseID() domain.PurchaseID
	NewInstallmentID() domain.InstallmentID
	NewInvoiceID() domain.InvoiceID
}
