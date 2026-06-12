package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	accountdomain "contai/internal/account/domain"
	authhttp "contai/internal/auth/adapters/http"
	categorydomain "contai/internal/category/domain"
	reportports "contai/internal/reports/app/ports"
	reportservices "contai/internal/reports/app/services"

	"github.com/gin-gonic/gin"
)

const timeFormatRFC3339 = "2006-01-02T15:04:05Z07:00"

type Handler struct {
	reportService reportports.ReportService
}

func NewHandler(reportService reportports.ReportService) Handler {
	return Handler{reportService: reportService}
}

func (handler Handler) GetFinancialReport(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	input, ok := parseFinancialInput(ctx)
	if !ok {
		return
	}
	input.UserID = authenticatedUser.UserID

	report, err := handler.reportService.GetFinancialReport(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, toFinancialReportResponse(report))
}

func (handler Handler) DownloadFinancialPDF(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	input, ok := parseFinancialInput(ctx)
	if !ok {
		return
	}
	input.UserID = authenticatedUser.UserID

	file, err := handler.reportService.GenerateFinancialPDF(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	writePDF(ctx, file)
}

func parseFinancialInput(ctx *gin.Context) (reportports.FinancialReportInput, bool) {
	startAt, err := time.Parse(time.RFC3339, ctx.Query("startAt"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid startAt"})
		return reportports.FinancialReportInput{}, false
	}
	endAt, err := time.Parse(time.RFC3339, ctx.Query("endAt"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid endAt"})
		return reportports.FinancialReportInput{}, false
	}
	input := reportports.FinancialReportInput{
		StartAt:          startAt,
		EndAt:            endAt,
		MovementType:     reportports.MovementType(defaultString(ctx.Query("movementType"), string(reportports.MovementTypeAll))),
		SettlementStatus: reportports.SettlementStatusFilter(defaultString(ctx.Query("settlementStatus"), string(reportports.SettlementStatusAll))),
		GroupBy:          reportports.ReportGroupBy(defaultString(ctx.Query("groupBy"), string(reportports.ReportGroupByNone))),
	}
	if value := ctx.Query("categoryId"); value != "" {
		categoryID := categorydomain.CategoryID(value)
		input.CategoryID = &categoryID
	}
	if value := ctx.Query("accountId"); value != "" {
		accountID := accountdomain.AccountID(value)
		input.AccountID = &accountID
	}
	return input, true
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func writePDF(ctx *gin.Context, file reportports.PDFFile) {
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Filename))
	ctx.Data(http.StatusOK, "application/pdf", file.Content)
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, accountdomain.ErrAccountUserIDRequired),
		errors.Is(err, reportservices.ErrReportPeriodInvalid),
		errors.Is(err, reportservices.ErrReportMovementTypeInvalid),
		errors.Is(err, reportservices.ErrReportSettlementStatusInvalid),
		errors.Is(err, reportservices.ErrReportGroupByInvalid):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid report"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
