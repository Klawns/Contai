package ids

import (
	"contai/internal/creditcards/domain"

	"github.com/google/uuid"
)

type UUIDCreditCardIDGenerator struct{}

func NewUUIDCreditCardIDGenerator() UUIDCreditCardIDGenerator {
	return UUIDCreditCardIDGenerator{}
}

func (UUIDCreditCardIDGenerator) NewCreditCardID() domain.CreditCardID {
	return domain.CreditCardID(uuid.NewString())
}

func (UUIDCreditCardIDGenerator) NewPurchaseID() domain.PurchaseID {
	return domain.PurchaseID(uuid.NewString())
}

func (UUIDCreditCardIDGenerator) NewInstallmentID() domain.InstallmentID {
	return domain.InstallmentID(uuid.NewString())
}

func (UUIDCreditCardIDGenerator) NewInvoiceID() domain.InvoiceID {
	return domain.InvoiceID(uuid.NewString())
}
