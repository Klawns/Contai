package http

import (
	"time"

	accountdomain "contai/internal/account/domain"
	"contai/internal/creditcards/app/ports"
	"contai/internal/creditcards/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type cardRequest struct {
	Name            string `json:"name" binding:"required"`
	LinkedAccountID string `json:"linkedAccountId" binding:"required"`
	LimitTotal      int64  `json:"limitTotal" binding:"required"`
	ClosingDay      int    `json:"closingDay" binding:"required"`
	DueDay          int    `json:"dueDay" binding:"required"`
	Status          string `json:"status"`
}

type purchaseRequest struct {
	CategoryID       string `json:"categoryId" binding:"required"`
	Description      string `json:"description" binding:"required"`
	TotalAmount      int64  `json:"totalAmount" binding:"required"`
	PurchaseDate     string `json:"purchaseDate" binding:"required"`
	InstallmentCount int    `json:"installmentCount" binding:"required"`
	Note             string `json:"note"`
}

type payInvoiceRequest struct {
	OccurredAt string `json:"occurredAt" binding:"required"`
	CategoryID string `json:"categoryId" binding:"required"`
	Note       string `json:"note"`
}

type cardResponse struct {
	ID              string `json:"id"`
	UserID          string `json:"userId"`
	Name            string `json:"name"`
	LinkedAccountID string `json:"linkedAccountId"`
	LimitTotal      int64  `json:"limitTotal"`
	LimitUsed       int64  `json:"limitUsed"`
	LimitAvailable  int64  `json:"limitAvailable"`
	ClosingDay      int    `json:"closingDay"`
	DueDay          int    `json:"dueDay"`
	Status          string `json:"status"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

type purchaseResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"userId"`
	CardID           string  `json:"cardId"`
	CategoryID       string  `json:"categoryId"`
	Description      string  `json:"description"`
	TotalAmount      int64   `json:"totalAmount"`
	PurchaseDate     string  `json:"purchaseDate"`
	InstallmentCount int     `json:"installmentCount"`
	Note             string  `json:"note"`
	Status           string  `json:"status"`
	CanceledAt       *string `json:"canceledAt"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type installmentResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	CardID         string `json:"cardId"`
	PurchaseID     string `json:"purchaseId"`
	InvoiceID      string `json:"invoiceId"`
	Number         int    `json:"number"`
	Amount         int64  `json:"amount"`
	Status         string `json:"status"`
	ReferenceMonth string `json:"referenceMonth"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type invoiceResponse struct {
	ID                   string                `json:"id"`
	UserID               string                `json:"userId"`
	CardID               string                `json:"cardId"`
	ReferenceMonth       string                `json:"referenceMonth"`
	ClosingAt            string                `json:"closingAt"`
	DueAt                string                `json:"dueAt"`
	Amount               int64                 `json:"amount"`
	Status               string                `json:"status"`
	EffectiveStatus      string                `json:"effectiveStatus"`
	PaidAt               *string               `json:"paidAt"`
	PaymentTransactionID *string               `json:"paymentTransactionId"`
	Installments         []installmentResponse `json:"installments"`
	CreatedAt            string                `json:"createdAt"`
	UpdatedAt            string                `json:"updatedAt"`
}

func toCreateCardInput(request cardRequest, userID string) ports.CreateCreditCardInput {
	return ports.CreateCreditCardInput{
		UserID:          userdomain.UserID(userID),
		Name:            request.Name,
		LinkedAccountID: accountdomain.AccountID(request.LinkedAccountID),
		LimitTotal:      financedomain.NewMoney(request.LimitTotal),
		ClosingDay:      request.ClosingDay,
		DueDay:          request.DueDay,
	}
}

func toUpdateCardInput(request cardRequest, userID string, cardID string) ports.UpdateCreditCardInput {
	return ports.UpdateCreditCardInput{
		UserID:          userdomain.UserID(userID),
		CardID:          domain.CreditCardID(cardID),
		Name:            request.Name,
		LinkedAccountID: accountdomain.AccountID(request.LinkedAccountID),
		LimitTotal:      financedomain.NewMoney(request.LimitTotal),
		ClosingDay:      request.ClosingDay,
		DueDay:          request.DueDay,
		Status:          domain.CreditCardStatus(request.Status),
	}
}

func toCardResponse(card ports.CreditCardDTO) cardResponse {
	return cardResponse{
		ID:              string(card.ID),
		UserID:          string(card.UserID),
		Name:            card.Name,
		LinkedAccountID: string(card.LinkedAccountID),
		LimitTotal:      card.LimitTotal.Cents(),
		LimitUsed:       card.LimitUsed.Cents(),
		LimitAvailable:  card.LimitAvailable.Cents(),
		ClosingDay:      card.ClosingDay,
		DueDay:          card.DueDay,
		Status:          string(card.Status),
		CreatedAt:       card.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:       card.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toCardResponses(cards []ports.CreditCardDTO) []cardResponse {
	responses := make([]cardResponse, 0, len(cards))
	for _, card := range cards {
		responses = append(responses, toCardResponse(card))
	}
	return responses
}

func toPurchaseResponse(purchase ports.PurchaseDTO) purchaseResponse {
	return purchaseResponse{
		ID:               string(purchase.ID),
		UserID:           string(purchase.UserID),
		CardID:           string(purchase.CardID),
		CategoryID:       string(purchase.CategoryID),
		Description:      purchase.Description,
		TotalAmount:      purchase.TotalAmount.Cents(),
		PurchaseDate:     purchase.PurchaseDate.Format(timeFormatRFC3339),
		InstallmentCount: purchase.InstallmentCount,
		Note:             purchase.Note,
		Status:           string(purchase.Status),
		CanceledAt:       timeToString(purchase.CanceledAt),
		CreatedAt:        purchase.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:        purchase.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toPurchaseResponses(purchases []ports.PurchaseDTO) []purchaseResponse {
	responses := make([]purchaseResponse, 0, len(purchases))
	for _, purchase := range purchases {
		responses = append(responses, toPurchaseResponse(purchase))
	}
	return responses
}

func toInvoiceResponse(invoice ports.InvoiceDTO) invoiceResponse {
	return invoiceResponse{
		ID:                   string(invoice.ID),
		UserID:               string(invoice.UserID),
		CardID:               string(invoice.CardID),
		ReferenceMonth:       invoice.ReferenceMonth.Format(timeFormatRFC3339),
		ClosingAt:            invoice.ClosingAt.Format(timeFormatRFC3339),
		DueAt:                invoice.DueAt.Format(timeFormatRFC3339),
		Amount:               invoice.Amount.Cents(),
		Status:               string(invoice.Status),
		EffectiveStatus:      string(invoice.EffectiveStatus),
		PaidAt:               timeToString(invoice.PaidAt),
		PaymentTransactionID: transactionIDToString(invoice.PaymentTransactionID),
		Installments:         toInstallmentResponses(invoice.Installments),
		CreatedAt:            invoice.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:            invoice.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toInvoiceResponses(invoices []ports.InvoiceDTO) []invoiceResponse {
	responses := make([]invoiceResponse, 0, len(invoices))
	for _, invoice := range invoices {
		responses = append(responses, toInvoiceResponse(invoice))
	}
	return responses
}

func toInstallmentResponses(installments []ports.InstallmentDTO) []installmentResponse {
	responses := make([]installmentResponse, 0, len(installments))
	for _, installment := range installments {
		responses = append(responses, installmentResponse{
			ID:             string(installment.ID),
			UserID:         string(installment.UserID),
			CardID:         string(installment.CardID),
			PurchaseID:     string(installment.PurchaseID),
			InvoiceID:      string(installment.InvoiceID),
			Number:         installment.Number,
			Amount:         installment.Amount.Cents(),
			Status:         string(installment.Status),
			ReferenceMonth: installment.ReferenceMonth.Format(timeFormatRFC3339),
			CreatedAt:      installment.CreatedAt.Format(timeFormatRFC3339),
			UpdatedAt:      installment.UpdatedAt.Format(timeFormatRFC3339),
		})
	}
	return responses
}

func parseRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func timeToString(value *time.Time) *string {
	if value == nil {
		return nil
	}
	converted := value.Format(timeFormatRFC3339)
	return &converted
}

func transactionIDToString(value *transactiondomain.TransactionID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}
