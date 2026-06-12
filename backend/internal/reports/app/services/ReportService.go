package services

import (
	"context"
	"errors"
	"sort"
	"time"

	accountdomain "contai/internal/account/domain"
	reportports "contai/internal/reports/app/ports"
	transactiondomain "contai/internal/transactions/domain"
)

var _ reportports.ReportService = ReportService{}

var (
	ErrReportPeriodInvalid           = errors.New("report period is invalid")
	ErrReportMovementTypeInvalid     = errors.New("report movement type is invalid")
	ErrReportSettlementStatusInvalid = errors.New("report settlement status is invalid")
	ErrReportGroupByInvalid          = errors.New("report group by is invalid")
)

type ReportService struct {
	repository reportports.ReportRepository
	renderer   reportports.PDFRenderer
}

func NewReportService(repository reportports.ReportRepository, renderer reportports.PDFRenderer) ReportService {
	return ReportService{repository: repository, renderer: renderer}
}

func (service ReportService) GetFinancialReport(ctx context.Context, input reportports.FinancialReportInput) (reportports.FinancialReportDTO, error) {
	if err := validateFinancialInput(input); err != nil {
		return reportports.FinancialReportDTO{}, err
	}
	movements, err := service.repository.ListFinancialMovements(ctx, reportports.ListFinancialMovementsInput{
		UserID:           input.UserID,
		StartAt:          input.StartAt,
		EndAt:            input.EndAt,
		MovementType:     input.MovementType,
		CategoryID:       input.CategoryID,
		AccountID:        input.AccountID,
		SettlementStatus: input.SettlementStatus,
	})
	if err != nil {
		return reportports.FinancialReportDTO{}, err
	}
	return buildFinancialReport(input, movements), nil
}

func (service ReportService) GenerateFinancialPDF(ctx context.Context, input reportports.FinancialReportInput) (reportports.PDFFile, error) {
	report, err := service.GetFinancialReport(ctx, input)
	if err != nil {
		return reportports.PDFFile{}, err
	}
	content, err := service.renderer.RenderFinancialReport(ctx, report)
	if err != nil {
		return reportports.PDFFile{}, err
	}
	return reportports.PDFFile{Filename: "contai-relatorio-financeiro.pdf", Content: content}, nil
}

func validateFinancialInput(input reportports.FinancialReportInput) error {
	if input.UserID == "" {
		return accountdomain.ErrAccountUserIDRequired
	}
	if input.StartAt.IsZero() || input.EndAt.IsZero() || input.EndAt.Before(input.StartAt) {
		return ErrReportPeriodInvalid
	}
	switch input.MovementType {
	case "", reportports.MovementTypeAll, reportports.MovementTypeIncome, reportports.MovementTypeExpense, reportports.MovementTypeCreditCardExpense, reportports.MovementTypeTransfer:
	default:
		return ErrReportMovementTypeInvalid
	}
	switch input.SettlementStatus {
	case "", reportports.SettlementStatusAll, reportports.SettlementStatusSettled, reportports.SettlementStatusPending:
	default:
		return ErrReportSettlementStatusInvalid
	}
	switch input.GroupBy {
	case "", reportports.ReportGroupByNone, reportports.ReportGroupByCategory, reportports.ReportGroupByAccount, reportports.ReportGroupByDay, reportports.ReportGroupByMonth:
	default:
		return ErrReportGroupByInvalid
	}
	return nil
}

func buildFinancialReport(input reportports.FinancialReportInput, movements []reportports.FinancialMovementDTO) reportports.FinancialReportDTO {
	generatedAt := input.Now
	if generatedAt.IsZero() {
		generatedAt = time.Now()
	}
	groupBy := input.GroupBy
	if groupBy == "" {
		groupBy = reportports.ReportGroupByNone
	}

	sort.SliceStable(movements, func(i, j int) bool {
		if movements[i].OccurredAt.Equal(movements[j].OccurredAt) {
			return movements[i].Description < movements[j].Description
		}
		return movements[i].OccurredAt.Before(movements[j].OccurredAt)
	})

	summary := buildSummary(movements)
	return reportports.FinancialReportDTO{
		Title:        "Relatorio financeiro",
		GeneratedAt:  generatedAt,
		StartAt:      input.StartAt,
		EndAt:        input.EndAt,
		Summary:      summary,
		Movements:    movements,
		Groups:       buildGroups(movements, groupBy),
		Charts:       buildCharts(movements),
		IncomeTotal:  summary.IncomeTotal,
		ExpenseTotal: summary.ExpenseTotal,
		NetTotal:     summary.PeriodResult,
	}
}

func buildSummary(movements []reportports.FinancialMovementDTO) reportports.FinancialReportSummaryDTO {
	var summary reportports.FinancialReportSummaryDTO
	for _, movement := range movements {
		if movement.Type == reportports.MovementTypeTransfer {
			continue
		}
		if movement.Type == reportports.MovementTypeIncome {
			summary.IncomeTotal = summary.IncomeTotal.Add(movement.Amount)
		} else {
			summary.ExpenseTotal = summary.ExpenseTotal.Add(movement.Amount)
		}
		if movement.SettlementStatus == transactiondomain.SettlementStatusSettled {
			summary.SettledTotal = summary.SettledTotal.Add(movement.Amount)
		} else {
			summary.PendingTotal = summary.PendingTotal.Add(movement.Amount)
		}
	}
	summary.PeriodResult = summary.IncomeTotal.Sub(summary.ExpenseTotal)
	return summary
}

func buildGroups(movements []reportports.FinancialMovementDTO, groupBy reportports.ReportGroupBy) []reportports.FinancialReportGroupDTO {
	if groupBy == reportports.ReportGroupByNone {
		return []reportports.FinancialReportGroupDTO{}
	}
	groups := make(map[string]*reportports.FinancialReportGroupDTO)
	for _, movement := range movements {
		key, label := groupKey(movement, groupBy)
		group := groups[key]
		if group == nil {
			group = &reportports.FinancialReportGroupDTO{Key: key, Label: label}
			groups[key] = group
		}
		group.Count++
		group.Total = group.Total.Add(movement.Amount)
		if movement.Type == reportports.MovementTypeIncome {
			group.IncomeTotal = group.IncomeTotal.Add(movement.Amount)
		} else if movement.Type != reportports.MovementTypeTransfer {
			group.ExpenseTotal = group.ExpenseTotal.Add(movement.Amount)
		}
		group.NetTotal = group.IncomeTotal.Sub(group.ExpenseTotal)
	}
	result := make([]reportports.FinancialReportGroupDTO, 0, len(groups))
	for _, group := range groups {
		result = append(result, *group)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result
}

func groupKey(movement reportports.FinancialMovementDTO, groupBy reportports.ReportGroupBy) (string, string) {
	switch groupBy {
	case reportports.ReportGroupByCategory:
		if movement.CategoryID != nil {
			return string(*movement.CategoryID), fallbackLabel(movement.CategoryName, "Sem categoria")
		}
		return "none", "Sem categoria"
	case reportports.ReportGroupByAccount:
		if movement.AccountID != nil {
			return string(*movement.AccountID), fallbackLabel(movement.AccountName, "Sem conta")
		}
		return "none", "Sem conta"
	case reportports.ReportGroupByMonth:
		return movement.OccurredAt.Format("2006-01"), movement.OccurredAt.Format("01/2006")
	default:
		return movement.OccurredAt.Format("2006-01-02"), movement.OccurredAt.Format("02/01/2006")
	}
}

func buildCharts(movements []reportports.FinancialMovementDTO) reportports.FinancialReportChartsDTO {
	return reportports.FinancialReportChartsDTO{
		IncomeVsExpense:    buildTimeSeries(movements, "2006-01", "01/2006"),
		ExpensesByCategory: buildCategoryExpenses(movements),
		Evolution:          buildTimeSeries(movements, "2006-01-02", "02/01"),
	}
}

func buildTimeSeries(movements []reportports.FinancialMovementDTO, keyFormat, labelFormat string) []reportports.FinancialReportSeriesPointDTO {
	points := make(map[string]*reportports.FinancialReportSeriesPointDTO)
	for _, movement := range movements {
		if movement.Type == reportports.MovementTypeTransfer {
			continue
		}
		key := movement.OccurredAt.Format(keyFormat)
		point := points[key]
		if point == nil {
			point = &reportports.FinancialReportSeriesPointDTO{Key: key, Label: movement.OccurredAt.Format(labelFormat)}
			points[key] = point
		}
		if movement.Type == reportports.MovementTypeIncome {
			point.IncomeTotal = point.IncomeTotal.Add(movement.Amount)
		} else {
			point.ExpenseTotal = point.ExpenseTotal.Add(movement.Amount)
		}
		point.NetTotal = point.IncomeTotal.Sub(point.ExpenseTotal)
	}
	result := make([]reportports.FinancialReportSeriesPointDTO, 0, len(points))
	for _, point := range points {
		result = append(result, *point)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
	return result
}

func buildCategoryExpenses(movements []reportports.FinancialMovementDTO) []reportports.FinancialReportCategoryChartDTO {
	totals := make(map[string]*reportports.FinancialReportCategoryChartDTO)
	for _, movement := range movements {
		if movement.Type == reportports.MovementTypeIncome || movement.Type == reportports.MovementTypeTransfer || movement.CategoryID == nil {
			continue
		}
		key := string(*movement.CategoryID)
		total := totals[key]
		if total == nil {
			total = &reportports.FinancialReportCategoryChartDTO{CategoryID: *movement.CategoryID, Name: fallbackLabel(movement.CategoryName, "Sem categoria")}
			totals[key] = total
		}
		total.Total = total.Total.Add(movement.Amount)
	}
	result := make([]reportports.FinancialReportCategoryChartDTO, 0, len(totals))
	for _, total := range totals {
		result = append(result, *total)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Total.Cents() > result[j].Total.Cents() })
	return result
}

func fallbackLabel(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
