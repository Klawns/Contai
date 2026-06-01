package http

import (
	"time"

	"contai/internal/dashboard/app/ports"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type monthlyDashboardResponse struct {
	UserID             string                      `json:"userId"`
	Period             periodResponse              `json:"period"`
	TotalBalance       int64                       `json:"totalBalance"`
	MonthlyIncome      int64                       `json:"monthlyIncome"`
	MonthlyExpense     int64                       `json:"monthlyExpense"`
	MonthlyNetBalance  int64                       `json:"monthlyNetBalance"`
	AccountBalances    []accountBalanceResponse    `json:"accountBalances"`
	ExpensesByCategory []expenseByCategoryResponse `json:"expensesByCategory"`
	RecentTransactions []recentTransactionResponse `json:"recentTransactions"`
}

type periodResponse struct {
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}

type accountBalanceResponse struct {
	AccountID  string `json:"accountId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Balance    int64  `json:"balance"`
	BankIconID string `json:"bankIconId"`
}

type expenseByCategoryResponse struct {
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Icon       string `json:"icon"`
	Total      int64  `json:"total"`
}

type recentTransactionResponse struct {
	ID                   string  `json:"id"`
	UserID               string  `json:"userId"`
	Type                 string  `json:"type"`
	Description          string  `json:"description"`
	Amount               int64   `json:"amount"`
	OccurredAt           string  `json:"occurredAt"`
	AccountID            *string `json:"accountId"`
	SourceAccountID      *string `json:"sourceAccountId"`
	DestinationAccountID *string `json:"destinationAccountId"`
	CategoryID           *string `json:"categoryId"`
	Status               string  `json:"status"`
	Note                 string  `json:"note"`
	RemovedAt            *string `json:"removedAt"`
	CreatedAt            string  `json:"createdAt"`
	UpdatedAt            string  `json:"updatedAt"`
}

func toMonthlyDashboardResponse(dashboard ports.MonthlyDashboardDTO) monthlyDashboardResponse {
	return monthlyDashboardResponse{
		UserID: string(dashboard.UserID),
		Period: periodResponse{
			StartAt: dashboard.Period.StartAt.Format(timeFormatRFC3339),
			EndAt:   dashboard.Period.EndAt.Format(timeFormatRFC3339),
		},
		TotalBalance:       dashboard.TotalBalance.Cents(),
		MonthlyIncome:      dashboard.MonthlyIncome.Cents(),
		MonthlyExpense:     dashboard.MonthlyExpense.Cents(),
		MonthlyNetBalance:  dashboard.MonthlyNetBalance.Cents(),
		AccountBalances:    toAccountBalanceResponses(dashboard.AccountBalances),
		ExpensesByCategory: toExpenseByCategoryResponses(dashboard.ExpensesByCategory),
		RecentTransactions: toRecentTransactionResponses(dashboard.RecentTransactions),
	}
}

func toAccountBalanceResponses(values []ports.AccountBalanceDTO) []accountBalanceResponse {
	responses := make([]accountBalanceResponse, 0, len(values))
	for _, value := range values {
		responses = append(responses, accountBalanceResponse{
			AccountID:  string(value.AccountID),
			Name:       value.Name,
			Type:       string(value.Type),
			Balance:    value.Balance.Cents(),
			BankIconID: value.BankIconID,
		})
	}
	return responses
}

func toExpenseByCategoryResponses(values []ports.CategoryExpenseDTO) []expenseByCategoryResponse {
	responses := make([]expenseByCategoryResponse, 0, len(values))
	for _, value := range values {
		responses = append(responses, expenseByCategoryResponse{
			CategoryID: string(value.CategoryID),
			Name:       value.Name,
			Color:      value.Color,
			Icon:       value.Icon,
			Total:      value.Total.Cents(),
		})
	}
	return responses
}

func toRecentTransactionResponses(values []ports.TransactionDTO) []recentTransactionResponse {
	responses := make([]recentTransactionResponse, 0, len(values))
	for _, value := range values {
		responses = append(responses, recentTransactionResponse{
			ID:                   string(value.ID),
			UserID:               string(value.UserID),
			Type:                 string(value.Type),
			Description:          value.Description,
			Amount:               value.Amount.Cents(),
			OccurredAt:           value.OccurredAt.Format(timeFormatRFC3339),
			AccountID:            stringPtr(value.AccountID),
			SourceAccountID:      stringPtr(value.SourceAccountID),
			DestinationAccountID: stringPtr(value.DestinationAccountID),
			CategoryID:           stringPtr(value.CategoryID),
			Status:               string(value.Status),
			Note:                 value.Note,
			RemovedAt:            timeToString(value.RemovedAt),
			CreatedAt:            value.CreatedAt.Format(timeFormatRFC3339),
			UpdatedAt:            value.UpdatedAt.Format(timeFormatRFC3339),
		})
	}
	return responses
}

func parseRFC3339(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func stringPtr[T ~string](value *T) *string {
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
