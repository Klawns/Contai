package persistence

import (
	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/creditcards/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

func toCreditCardEntity(card domain.CreditCard) CreditCardEntity {
	return CreditCardEntity{
		ID:              string(card.ID),
		UserID:          string(card.UserID),
		Name:            card.Name,
		LinkedAccountID: string(card.LinkedAccountID),
		LimitTotal:      card.LimitTotal.Cents(),
		ClosingDay:      card.ClosingDay,
		DueDay:          card.DueDay,
		Status:          string(card.Status),
		CreatedAt:       card.CreatedAt,
		UpdatedAt:       card.UpdatedAt,
	}
}

func toDomainCreditCard(entity CreditCardEntity) (domain.CreditCard, error) {
	return domain.RehydrateCreditCard(
		domain.CreditCardID(entity.ID),
		userdomain.UserID(entity.UserID),
		entity.Name,
		accountdomain.AccountID(entity.LinkedAccountID),
		financedomain.NewMoney(entity.LimitTotal),
		entity.ClosingDay,
		entity.DueDay,
		domain.CreditCardStatus(entity.Status),
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}

func toPurchaseEntity(purchase domain.Purchase) CardPurchaseEntity {
	return CardPurchaseEntity{
		ID:               string(purchase.ID),
		UserID:           string(purchase.UserID),
		CardID:           string(purchase.CardID),
		CategoryID:       string(purchase.CategoryID),
		Description:      purchase.Description,
		TotalAmount:      purchase.TotalAmount.Cents(),
		PurchaseDate:     purchase.PurchaseDate,
		InstallmentCount: purchase.InstallmentCount,
		Note:             purchase.Note,
		Status:           string(purchase.Status),
		CanceledAt:       purchase.CanceledAt,
		CreatedAt:        purchase.CreatedAt,
		UpdatedAt:        purchase.UpdatedAt,
	}
}

func toDomainPurchase(entity CardPurchaseEntity) (domain.Purchase, error) {
	return domain.RehydratePurchase(
		domain.PurchaseID(entity.ID),
		userdomain.UserID(entity.UserID),
		domain.CreditCardID(entity.CardID),
		categorydomain.CategoryID(entity.CategoryID),
		entity.Description,
		financedomain.NewMoney(entity.TotalAmount),
		entity.PurchaseDate,
		entity.InstallmentCount,
		entity.Note,
		domain.PurchaseStatus(entity.Status),
		entity.CanceledAt,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}

func toInstallmentEntity(installment domain.Installment) CardInstallmentEntity {
	return CardInstallmentEntity{
		ID:             string(installment.ID),
		UserID:         string(installment.UserID),
		CardID:         string(installment.CardID),
		PurchaseID:     string(installment.PurchaseID),
		InvoiceID:      string(installment.InvoiceID),
		Number:         installment.Number,
		Amount:         installment.Amount.Cents(),
		Status:         string(installment.Status),
		ReferenceMonth: installment.ReferenceMonth,
		CreatedAt:      installment.CreatedAt,
		UpdatedAt:      installment.UpdatedAt,
	}
}

func toDomainInstallment(entity CardInstallmentEntity) (domain.Installment, error) {
	return domain.RehydrateInstallment(
		domain.InstallmentID(entity.ID),
		userdomain.UserID(entity.UserID),
		domain.CreditCardID(entity.CardID),
		domain.PurchaseID(entity.PurchaseID),
		domain.InvoiceID(entity.InvoiceID),
		entity.Number,
		financedomain.NewMoney(entity.Amount),
		domain.PurchaseStatus(entity.Status),
		entity.ReferenceMonth,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}

func toInvoiceEntity(invoice domain.Invoice) CardInvoiceEntity {
	return CardInvoiceEntity{
		ID:                   string(invoice.ID),
		UserID:               string(invoice.UserID),
		CardID:               string(invoice.CardID),
		ReferenceMonth:       invoice.ReferenceMonth,
		ClosingAt:            invoice.ClosingAt,
		DueAt:                invoice.DueAt,
		Status:               string(invoice.Status),
		PaidAt:               invoice.PaidAt,
		PaymentTransactionID: transactionIDToString(invoice.PaymentTransactionID),
		CreatedAt:            invoice.CreatedAt,
		UpdatedAt:            invoice.UpdatedAt,
	}
}

func toDomainInvoice(entity CardInvoiceEntity) (domain.Invoice, error) {
	return domain.RehydrateInvoice(
		domain.InvoiceID(entity.ID),
		userdomain.UserID(entity.UserID),
		domain.CreditCardID(entity.CardID),
		entity.ReferenceMonth,
		entity.ClosingAt,
		entity.DueAt,
		domain.InvoiceStatus(entity.Status),
		entity.PaidAt,
		stringToTransactionID(entity.PaymentTransactionID),
		entity.CreatedAt,
		entity.UpdatedAt,
	)
}

func transactionIDToString(value *transactiondomain.TransactionID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}

func stringToTransactionID(value *string) *transactiondomain.TransactionID {
	if value == nil {
		return nil
	}
	converted := transactiondomain.TransactionID(*value)
	return &converted
}
