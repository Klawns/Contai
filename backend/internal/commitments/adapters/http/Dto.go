package http

import (
	"time"

	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/commitments/app/ports"
	"contai/internal/commitments/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"
	userdomain "contai/internal/users/domain"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type recurrenceRequest struct {
	Frequency string `json:"frequency" binding:"required"`
	Interval  int    `json:"interval" binding:"required"`
	EndsAt    string `json:"endsAt"`
}

type commitmentRequest struct {
	Description string             `json:"description" binding:"required"`
	Amount      int64              `json:"amount" binding:"required"`
	DueAt       string             `json:"dueAt" binding:"required"`
	AccountID   string             `json:"accountId" binding:"required"`
	CategoryID  string             `json:"categoryId" binding:"required"`
	Note        string             `json:"note"`
	Recurrence  *recurrenceRequest `json:"recurrence"`
}

type settlementRequest struct {
	Amount     int64  `json:"amount" binding:"required"`
	OccurredAt string `json:"occurredAt" binding:"required"`
	AccountID  string `json:"accountId" binding:"required"`
	CategoryID string `json:"categoryId" binding:"required"`
	Note       string `json:"note"`
}

type recurrenceResponse struct {
	Frequency string  `json:"frequency"`
	Interval  int     `json:"interval"`
	EndsAt    *string `json:"endsAt"`
}

type commitmentResponse struct {
	ID                      string              `json:"id"`
	UserID                  string              `json:"userId"`
	Type                    string              `json:"type"`
	Description             string              `json:"description"`
	Amount                  int64               `json:"amount"`
	DueAt                   string              `json:"dueAt"`
	AccountID               string              `json:"accountId"`
	CategoryID              string              `json:"categoryId"`
	Note                    string              `json:"note"`
	Status                  string              `json:"status"`
	EffectiveStatus         string              `json:"effectiveStatus"`
	Recurrence              *recurrenceResponse `json:"recurrence"`
	SettledAt               *string             `json:"settledAt"`
	SettlementTransactionID *string             `json:"settlementTransactionId"`
	CanceledAt              *string             `json:"canceledAt"`
	CreatedAt               string              `json:"createdAt"`
	UpdatedAt               string              `json:"updatedAt"`
}

func requestToCreateInput(
	request commitmentRequest,
	userID string,
	commitmentType domain.CommitmentType,
) (ports.CreateCommitmentInput, error) {
	dueAt, err := parseTime(request.DueAt)
	if err != nil {
		return ports.CreateCommitmentInput{}, err
	}
	recurrence, err := parseRecurrence(request.Recurrence)
	if err != nil {
		return ports.CreateCommitmentInput{}, err
	}
	return ports.CreateCommitmentInput{
		UserID:      userdomain.UserID(userID),
		Type:        commitmentType,
		Description: request.Description,
		Amount:      financedomain.NewMoney(request.Amount),
		DueAt:       dueAt,
		AccountID:   accountdomain.AccountID(request.AccountID),
		CategoryID:  categorydomain.CategoryID(request.CategoryID),
		Note:        request.Note,
		Recurrence:  recurrence,
	}, nil
}

func requestToUpdateInput(
	request commitmentRequest,
	userID string,
	commitmentID string,
	commitmentType domain.CommitmentType,
) (ports.UpdateCommitmentInput, error) {
	createInput, err := requestToCreateInput(request, userID, commitmentType)
	if err != nil {
		return ports.UpdateCommitmentInput{}, err
	}
	return ports.UpdateCommitmentInput{
		UserID:       createInput.UserID,
		CommitmentID: domain.CommitmentID(commitmentID),
		Type:         createInput.Type,
		Description:  createInput.Description,
		Amount:       createInput.Amount,
		DueAt:        createInput.DueAt,
		AccountID:    createInput.AccountID,
		CategoryID:   createInput.CategoryID,
		Note:         createInput.Note,
		Recurrence:   createInput.Recurrence,
	}, nil
}

func requestToSettleInput(
	request settlementRequest,
	userID string,
	commitmentID string,
	commitmentType domain.CommitmentType,
) (ports.SettleCommitmentInput, error) {
	occurredAt, err := parseTime(request.OccurredAt)
	if err != nil {
		return ports.SettleCommitmentInput{}, err
	}
	return ports.SettleCommitmentInput{
		UserID:       userdomain.UserID(userID),
		CommitmentID: domain.CommitmentID(commitmentID),
		Type:         commitmentType,
		Amount:       financedomain.NewMoney(request.Amount),
		SettledAt:    occurredAt,
		AccountID:    accountdomain.AccountID(request.AccountID),
		CategoryID:   categorydomain.CategoryID(request.CategoryID),
		Note:         request.Note,
	}, nil
}

func toCommitmentResponse(commitment ports.CommitmentDTO) commitmentResponse {
	return commitmentResponse{
		ID:                      string(commitment.ID),
		UserID:                  string(commitment.UserID),
		Type:                    string(commitment.Type),
		Description:             commitment.Description,
		Amount:                  commitment.Amount.Cents(),
		DueAt:                   commitment.DueAt.Format(timeFormatRFC3339),
		AccountID:               string(commitment.AccountID),
		CategoryID:              string(commitment.CategoryID),
		Note:                    commitment.Note,
		Status:                  string(commitment.Status),
		EffectiveStatus:         string(commitment.EffectiveStatus),
		Recurrence:              toRecurrenceResponse(commitment.Recurrence),
		SettledAt:               timeToString(commitment.SettledAt),
		SettlementTransactionID: transactionIDToString(commitment.SettlementTransactionID),
		CanceledAt:              timeToString(commitment.CanceledAt),
		CreatedAt:               commitment.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:               commitment.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toCommitmentResponses(commitments []ports.CommitmentDTO) []commitmentResponse {
	responses := make([]commitmentResponse, 0, len(commitments))
	for _, commitment := range commitments {
		responses = append(responses, toCommitmentResponse(commitment))
	}
	return responses
}

func parseRecurrence(request *recurrenceRequest) (*domain.Recurrence, error) {
	if request == nil {
		return nil, nil
	}
	var endsAt *time.Time
	if request.EndsAt != "" {
		parsed, err := parseTime(request.EndsAt)
		if err != nil {
			return nil, err
		}
		endsAt = &parsed
	}
	return &domain.Recurrence{
		Frequency: domain.RecurrenceFrequency(request.Frequency),
		Interval:  request.Interval,
		EndsAt:    endsAt,
	}, nil
}

func toRecurrenceResponse(recurrence *domain.Recurrence) *recurrenceResponse {
	if recurrence == nil {
		return nil
	}
	return &recurrenceResponse{
		Frequency: string(recurrence.Frequency),
		Interval:  recurrence.Interval,
		EndsAt:    timeToString(recurrence.EndsAt),
	}
}

func parseTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, domain.ErrCommitmentDueAtRequired
	}
	return parsed, nil
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
