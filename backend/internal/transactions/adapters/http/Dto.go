package http

import (
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	financedomain "contai/internal/finance/domain"
	"contai/internal/transactions/app/ports"
	"contai/internal/transactions/domain"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type createTransactionRequest struct {
	Description      string             `json:"description" binding:"required"`
	Amount           int64              `json:"amount" binding:"required"`
	OccurredAt       string             `json:"occurredAt" binding:"required"`
	AccountID        *string            `json:"accountId"`
	CategoryID       string             `json:"categoryId"`
	SettlementStatus *string            `json:"settlementStatus"`
	SettledAt        *string            `json:"settledAt"`
	RecurrenceType   string             `json:"recurrenceType"`
	Recurrence       *recurrenceRequest `json:"recurrence"`
	Note             string             `json:"note"`
}

type createTransferRequest struct {
	Description          string `json:"description" binding:"required"`
	Amount               int64  `json:"amount" binding:"required"`
	OccurredAt           string `json:"occurredAt" binding:"required"`
	SourceAccountID      string `json:"sourceAccountId" binding:"required"`
	DestinationAccountID string `json:"destinationAccountId" binding:"required"`
	Note                 string `json:"note"`
}

type updateTransactionRequest struct {
	Description          string             `json:"description" binding:"required"`
	Amount               int64              `json:"amount" binding:"required"`
	OccurredAt           string             `json:"occurredAt" binding:"required"`
	AccountID            *string            `json:"accountId"`
	SourceAccountID      string             `json:"sourceAccountId"`
	DestinationAccountID string             `json:"destinationAccountId"`
	CategoryID           string             `json:"categoryId"`
	SettlementStatus     *string            `json:"settlementStatus"`
	SettledAt            *string            `json:"settledAt"`
	RecurrenceType       string             `json:"recurrenceType"`
	Recurrence           *recurrenceRequest `json:"recurrence"`
	Note                 string             `json:"note"`
}

type recurrenceRequest struct {
	Frequency  string  `json:"frequency"`
	Quantity   *int    `json:"quantity"`
	StartsAt   string  `json:"startsAt"`
	EndsAt     *string `json:"endsAt"`
	DayOfMonth *int    `json:"dayOfMonth"`
}

type transactionResponse struct {
	ID                   string              `json:"id"`
	UserID               string              `json:"userId"`
	Type                 string              `json:"type"`
	Description          string              `json:"description"`
	Amount               int64               `json:"amount"`
	OccurredAt           string              `json:"occurredAt"`
	AccountID            *string             `json:"accountId"`
	SourceAccountID      *string             `json:"sourceAccountId"`
	DestinationAccountID *string             `json:"destinationAccountId"`
	CategoryID           *string             `json:"categoryId"`
	Status               string              `json:"status"`
	OriginType           string              `json:"originType"`
	OriginID             *string             `json:"originId"`
	SettlementStatus     string              `json:"settlementStatus"`
	SettledAt            *string             `json:"settledAt"`
	RecurrenceType       string              `json:"recurrenceType"`
	Recurrence           *recurrenceResponse `json:"recurrence"`
	Note                 string              `json:"note"`
	RemovedAt            *string             `json:"removedAt"`
	CreatedAt            string              `json:"createdAt"`
	UpdatedAt            string              `json:"updatedAt"`
}

type recurrenceResponse struct {
	Frequency  string  `json:"frequency"`
	Quantity   *int    `json:"quantity"`
	StartsAt   string  `json:"startsAt"`
	EndsAt     *string `json:"endsAt"`
	DayOfMonth *int    `json:"dayOfMonth"`
}

func toTransactionResponse(transaction ports.TransactionDTO) transactionResponse {
	return transactionResponse{
		ID:                   string(transaction.ID),
		UserID:               string(transaction.UserID),
		Type:                 string(transaction.Type),
		Description:          transaction.Description,
		Amount:               transaction.Amount.Cents(),
		OccurredAt:           transaction.OccurredAt.Format(timeFormatRFC3339),
		AccountID:            accountIDToString(transaction.AccountID),
		SourceAccountID:      accountIDToString(transaction.SourceAccountID),
		DestinationAccountID: accountIDToString(transaction.DestinationAccountID),
		CategoryID:           categoryIDToString(transaction.CategoryID),
		Status:               string(transaction.Status),
		OriginType:           string(transaction.OriginType),
		OriginID:             transaction.OriginID,
		SettlementStatus:     string(transaction.SettlementStatus),
		SettledAt:            timeToString(transaction.SettledAt),
		RecurrenceType:       string(transaction.RecurrenceType),
		Recurrence:           toRecurrenceResponse(transaction.Recurrence),
		Note:                 transaction.Note,
		RemovedAt:            timeToString(transaction.RemovedAt),
		CreatedAt:            transaction.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:            transaction.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toTransactionResponses(transactions []ports.TransactionDTO) []transactionResponse {
	responses := make([]transactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		responses = append(responses, toTransactionResponse(transaction))
	}
	return responses
}

func parseOccurredAt(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, domain.ErrTransactionOccurredAtRequired
	}
	return parsed, nil
}

func parseTransactionType(value string) (*domain.TransactionType, error) {
	if value == "" {
		return nil, nil
	}
	transactionType := domain.TransactionType(value)
	if transactionType != domain.TransactionTypeIncome && transactionType != domain.TransactionTypeExpense && transactionType != domain.TransactionTypeTransfer {
		return nil, domain.ErrTransactionInvalidType
	}
	return &transactionType, nil
}

func parseSettlementStatus(value string) (*domain.SettlementStatus, error) {
	if value == "" {
		return nil, nil
	}
	status := domain.SettlementStatus(value)
	if status != domain.SettlementStatusSettled && status != domain.SettlementStatusPending {
		return nil, domain.ErrTransactionInvalidSettlementStatus
	}
	return &status, nil
}

func parseRequiredSettlementStatus(value *string) (domain.SettlementStatus, error) {
	if value == nil || *value == "" {
		return "", domain.ErrTransactionSettlementStatusRequired
	}
	status := domain.SettlementStatus(*value)
	if status != domain.SettlementStatusSettled && status != domain.SettlementStatusPending {
		return "", domain.ErrTransactionInvalidSettlementStatus
	}
	return status, nil
}

func parseOptionalTime(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return nil, domain.ErrTransactionInvalidSettledAt
	}
	return &parsed, nil
}

func parseRecurrenceType(value string) domain.RecurrenceType {
	if value == "" {
		return domain.RecurrenceTypeNone
	}
	return domain.RecurrenceType(value)
}

func parseRecurrencePayload(value string, request *recurrenceRequest) (domain.RecurrenceType, *domain.Recurrence, error) {
	recurrenceType := parseRecurrenceType(value)
	if recurrenceType == domain.RecurrenceTypeNone {
		return recurrenceType, nil, nil
	}
	recurrence, err := parseRecurrence(request)
	return recurrenceType, recurrence, err
}

func parseRecurrence(request *recurrenceRequest) (*domain.Recurrence, error) {
	if request == nil {
		return nil, nil
	}
	startsAt, err := time.Parse(time.RFC3339, request.StartsAt)
	if err != nil {
		return nil, domain.ErrTransactionInvalidRecurrence
	}
	var endsAt *time.Time
	if request.EndsAt != nil && *request.EndsAt != "" {
		parsed, err := time.Parse(time.RFC3339, *request.EndsAt)
		if err != nil {
			return nil, domain.ErrTransactionInvalidRecurrence
		}
		endsAt = &parsed
	}
	return &domain.Recurrence{
		Frequency:  domain.RecurrenceFrequency(request.Frequency),
		Quantity:   request.Quantity,
		StartsAt:   startsAt,
		EndsAt:     endsAt,
		DayOfMonth: request.DayOfMonth,
	}, nil
}

func parseOptionalAccountID(value *string) *accountdomain.AccountID {
	if value == nil || *value == "" {
		return nil
	}
	accountID := accountdomain.AccountID(*value)
	return &accountID
}

func moneyFromCents(value int64) financedomain.Money {
	return financedomain.NewMoney(value)
}

func accountIDToString(value *accountdomain.AccountID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}

func categoryIDToString(value *categorydomain.CategoryID) *string {
	if value == nil {
		return nil
	}
	converted := string(*value)
	return &converted
}

func timeToString(value *time.Time) *string {
	if value == nil {
		return nil
	}
	converted := value.Format(timeFormatRFC3339)
	return &converted
}

func toRecurrenceResponse(value *domain.Recurrence) *recurrenceResponse {
	if value == nil {
		return nil
	}
	return &recurrenceResponse{
		Frequency:  string(value.Frequency),
		Quantity:   value.Quantity,
		StartsAt:   value.StartsAt.Format(timeFormatRFC3339),
		EndsAt:     timeToString(value.EndsAt),
		DayOfMonth: value.DayOfMonth,
	}
}
