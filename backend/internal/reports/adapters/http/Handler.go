package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	accountdomain "contai/internal/account/domain"
	authhttp "contai/internal/auth/adapters/http"
	reportports "contai/internal/reports/app/ports"
	reportservices "contai/internal/reports/app/services"
	transactiondomain "contai/internal/transactions/domain"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	reportService reportports.ReportService
}

func NewHandler(reportService reportports.ReportService) Handler {
	return Handler{reportService: reportService}
}

func (handler Handler) DownloadAccountsPDF(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := handler.reportService.GenerateAccountsPDF(ctx.Request.Context(), reportports.GenerateAccountsReportInput{
		UserID: authenticatedUser.UserID,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	writePDF(ctx, file)
}

func (handler Handler) DownloadTransactionsPDF(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startAt, endAt, ok := parsePeriod(ctx)
	if !ok {
		return
	}

	file, err := handler.reportService.GenerateTransactionsPDF(ctx.Request.Context(), reportports.GenerateTransactionsReportInput{
		UserID:  authenticatedUser.UserID,
		StartAt: startAt,
		EndAt:   endAt,
		Type:    transactiondomain.TransactionType(ctx.Query("type")),
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	writePDF(ctx, file)
}

func (handler Handler) DownloadPeriodPDF(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startAt, endAt, ok := parsePeriod(ctx)
	if !ok {
		return
	}

	file, err := handler.reportService.GeneratePeriodPDF(ctx.Request.Context(), reportports.PeriodReportInput{
		UserID:  authenticatedUser.UserID,
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	writePDF(ctx, file)
}

func (handler Handler) DownloadMonthlyPDF(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startAt, endAt, ok := parsePeriod(ctx)
	if !ok {
		return
	}

	file, err := handler.reportService.GenerateMonthlyPDF(ctx.Request.Context(), reportports.PeriodReportInput{
		UserID:  authenticatedUser.UserID,
		StartAt: startAt,
		EndAt:   endAt,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	writePDF(ctx, file)
}

func (handler Handler) DownloadAccountPDF(ctx *gin.Context) {
	authenticatedUser, ok := authhttp.AuthenticatedUserFromContext(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startAt, endAt, ok := parsePeriod(ctx)
	if !ok {
		return
	}

	file, err := handler.reportService.GenerateAccountPDF(ctx.Request.Context(), reportports.GenerateAccountReportInput{
		UserID:    authenticatedUser.UserID,
		AccountID: accountdomain.AccountID(ctx.Param("accountID")),
		StartAt:   startAt,
		EndAt:     endAt,
	})
	if err != nil {
		writeError(ctx, err)
		return
	}

	writePDF(ctx, file)
}

func writePDF(ctx *gin.Context, file reportports.PDFFile) {
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Filename))
	ctx.Data(http.StatusOK, "application/pdf", file.Content)
}

func parsePeriod(ctx *gin.Context) (time.Time, time.Time, bool) {
	startAt, err := time.Parse(time.RFC3339, ctx.Query("startAt"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid startAt"})
		return time.Time{}, time.Time{}, false
	}
	endAt, err := time.Parse(time.RFC3339, ctx.Query("endAt"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid endAt"})
		return time.Time{}, time.Time{}, false
	}
	return startAt, endAt, true
}

func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, accountdomain.ErrAccountUserIDRequired),
		errors.Is(err, reportservices.ErrReportPeriodInvalid),
		errors.Is(err, reportservices.ErrReportTransactionTypeInvalid),
		errors.Is(err, reportservices.ErrReportMonthlyPeriodInvalid),
		errors.Is(err, reportservices.ErrReportAccountIDRequired):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid report"})
	case errors.Is(err, reportservices.ErrReportAccountNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "account not found"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
