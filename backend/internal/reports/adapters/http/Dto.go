package http

import (
	accountdomain "contai/internal/account/domain"
	categorydomain "contai/internal/category/domain"
	reportports "contai/internal/reports/app/ports"
)

type financialReportResponse struct {
	Summary   financialReportSummaryResponse `json:"summary"`
	Movements []financialMovementResponse    `json:"movements"`
	Groups    []financialReportGroupResponse `json:"groups"`
	Charts    financialReportChartsResponse  `json:"charts"`
}

type financialReportSummaryResponse struct {
	IncomeTotal  int64 `json:"incomeTotal"`
	ExpenseTotal int64 `json:"expenseTotal"`
	PeriodResult int64 `json:"periodResult"`
	PendingTotal int64 `json:"pendingTotal"`
	SettledTotal int64 `json:"settledTotal"`
}

type financialMovementResponse struct {
	ID               string  `json:"id"`
	Source           string  `json:"source"`
	Type             string  `json:"type"`
	Description      string  `json:"description"`
	Amount           int64   `json:"amount"`
	OccurredAt       string  `json:"occurredAt"`
	CategoryID       *string `json:"categoryId"`
	CategoryName     string  `json:"categoryName"`
	AccountID        *string `json:"accountId"`
	AccountName      string  `json:"accountName"`
	SettlementStatus string  `json:"settlementStatus"`
}

type financialReportGroupResponse struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	IncomeTotal  int64  `json:"incomeTotal"`
	ExpenseTotal int64  `json:"expenseTotal"`
	NetTotal     int64  `json:"netTotal"`
	Total        int64  `json:"total"`
	Count        int    `json:"count"`
}

type financialReportSeriesPointResponse struct {
	Key          string `json:"key"`
	Label        string `json:"label"`
	IncomeTotal  int64  `json:"incomeTotal"`
	ExpenseTotal int64  `json:"expenseTotal"`
	NetTotal     int64  `json:"netTotal"`
}

type financialReportCategoryChartResponse struct {
	CategoryID string `json:"categoryId"`
	Name       string `json:"name"`
	Total      int64  `json:"total"`
}

type financialReportChartsResponse struct {
	IncomeVsExpense    []financialReportSeriesPointResponse   `json:"incomeVsExpense"`
	ExpensesByCategory []financialReportCategoryChartResponse `json:"expensesByCategory"`
	Evolution          []financialReportSeriesPointResponse   `json:"evolution"`
}

func toFinancialReportResponse(report reportports.FinancialReportDTO) financialReportResponse {
	return financialReportResponse{
		Summary: financialReportSummaryResponse{
			IncomeTotal:  report.Summary.IncomeTotal.Cents(),
			ExpenseTotal: report.Summary.ExpenseTotal.Cents(),
			PeriodResult: report.Summary.PeriodResult.Cents(),
			PendingTotal: report.Summary.PendingTotal.Cents(),
			SettledTotal: report.Summary.SettledTotal.Cents(),
		},
		Movements: toMovementResponses(report.Movements),
		Groups:    toGroupResponses(report.Groups),
		Charts: financialReportChartsResponse{
			IncomeVsExpense:    toSeriesPointResponses(report.Charts.IncomeVsExpense),
			ExpensesByCategory: toCategoryChartResponses(report.Charts.ExpensesByCategory),
			Evolution:          toSeriesPointResponses(report.Charts.Evolution),
		},
	}
}

func toMovementResponses(movements []reportports.FinancialMovementDTO) []financialMovementResponse {
	responses := make([]financialMovementResponse, 0, len(movements))
	for _, movement := range movements {
		responses = append(responses, financialMovementResponse{
			ID:               movement.ID,
			Source:           string(movement.Source),
			Type:             string(movement.Type),
			Description:      movement.Description,
			Amount:           movement.Amount.Cents(),
			OccurredAt:       movement.OccurredAt.Format(timeFormatRFC3339),
			CategoryID:       categoryIDToString(movement.CategoryID),
			CategoryName:     movement.CategoryName,
			AccountID:        accountIDToString(movement.AccountID),
			AccountName:      movement.AccountName,
			SettlementStatus: string(movement.SettlementStatus),
		})
	}
	return responses
}

func toGroupResponses(groups []reportports.FinancialReportGroupDTO) []financialReportGroupResponse {
	responses := make([]financialReportGroupResponse, 0, len(groups))
	for _, group := range groups {
		responses = append(responses, financialReportGroupResponse{
			Key:          group.Key,
			Label:        group.Label,
			IncomeTotal:  group.IncomeTotal.Cents(),
			ExpenseTotal: group.ExpenseTotal.Cents(),
			NetTotal:     group.NetTotal.Cents(),
			Total:        group.Total.Cents(),
			Count:        group.Count,
		})
	}
	return responses
}

func toSeriesPointResponses(points []reportports.FinancialReportSeriesPointDTO) []financialReportSeriesPointResponse {
	responses := make([]financialReportSeriesPointResponse, 0, len(points))
	for _, point := range points {
		responses = append(responses, financialReportSeriesPointResponse{
			Key:          point.Key,
			Label:        point.Label,
			IncomeTotal:  point.IncomeTotal.Cents(),
			ExpenseTotal: point.ExpenseTotal.Cents(),
			NetTotal:     point.NetTotal.Cents(),
		})
	}
	return responses
}

func toCategoryChartResponses(points []reportports.FinancialReportCategoryChartDTO) []financialReportCategoryChartResponse {
	responses := make([]financialReportCategoryChartResponse, 0, len(points))
	for _, point := range points {
		responses = append(responses, financialReportCategoryChartResponse{
			CategoryID: string(point.CategoryID),
			Name:       point.Name,
			Total:      point.Total.Cents(),
		})
	}
	return responses
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
