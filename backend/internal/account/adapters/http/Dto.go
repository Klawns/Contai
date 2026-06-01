package http

import (
	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type createAccountRequest struct {
	Name           string `json:"name" binding:"required"`
	Type           string `json:"type" binding:"required"`
	InitialBalance int64  `json:"initialBalance"`
	BankIconID     string `json:"bankIconId" binding:"required"`
}

type updateAccountRequest struct {
	Name       string `json:"name" binding:"required"`
	Type       string `json:"type" binding:"required"`
	BankIconID string `json:"bankIconId" binding:"required"`
}

type accountResponse struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	InitialBalance int64  `json:"initialBalance"`
	CurrentBalance int64  `json:"currentBalance"`
	BankIconID     string `json:"bankIconId"`
	Status         string `json:"status"`
	CreatedAt      string `json:"createdAt"`
	UpdatedAt      string `json:"updatedAt"`
}

type totalBalanceResponse struct {
	TotalBalance int64 `json:"totalBalance"`
}

func toAccountResponse(account ports.AccountDTO) accountResponse {
	return accountResponse{
		ID:             string(account.ID),
		UserID:         string(account.UserID),
		Name:           account.Name,
		Type:           string(account.Type),
		InitialBalance: account.InitialBalance.Cents(),
		CurrentBalance: account.CurrentBalance.Cents(),
		BankIconID:     account.BankIconID,
		Status:         string(account.Status),
		CreatedAt:      account.CreatedAt.Format(timeFormatRFC3339),
		UpdatedAt:      account.UpdatedAt.Format(timeFormatRFC3339),
	}
}

func toAccountResponses(accounts []ports.AccountDTO) []accountResponse {
	responses := make([]accountResponse, 0, len(accounts))
	for _, account := range accounts {
		responses = append(responses, toAccountResponse(account))
	}
	return responses
}

func parseAccountStatus(value string) (*domain.AccountStatus, error) {
	if value == "" {
		return nil, nil
	}
	status := domain.AccountStatus(value)
	if status != domain.AccountStatusActive && status != domain.AccountStatusInactive {
		return nil, domain.ErrAccountInvalidStatus
	}
	return &status, nil
}
