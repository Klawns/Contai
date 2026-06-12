package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	authdomain "contai/internal/auth/domain"
	financedomain "contai/internal/finance/domain"
	reportports "contai/internal/reports/app/ports"

	"github.com/gin-gonic/gin"
)

func TestHandlerRequiresAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(&fakeReportService{}).RegisterForTest(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/financial", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestHandlerReturnsFinancialReportJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeReportService{report: reportports.FinancialReportDTO{
		Summary: reportports.FinancialReportSummaryDTO{IncomeTotal: financedomain.NewMoney(1000)},
	}}
	router := authenticatedReportsRouter(service)
	query := "?startAt=2026-06-01T00:00:00Z&endAt=2026-06-30T23:59:59Z&movementType=income&settlementStatus=settled&groupBy=category"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/financial"+query, nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.input.UserID != "authenticated-user" || service.input.MovementType != reportports.MovementTypeIncome {
		t.Fatalf("expected parsed input, got %#v", service.input)
	}
	if recorder.Body.String() == "" {
		t.Fatalf("expected json body")
	}
}

func TestHandlerDownloadsFinancialPDF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeReportService{file: reportports.PDFFile{
		Filename: "contai-relatorio-financeiro.pdf",
		Content:  []byte("%PDF"),
	}}
	router := authenticatedReportsRouter(service)
	query := "?startAt=2026-06-01T00:00:00Z&endAt=2026-06-30T23:59:59Z"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/financial/pdf"+query, nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if recorder.Header().Get("Content-Type") != "application/pdf" {
		t.Fatalf("expected pdf content type, got %s", recorder.Header().Get("Content-Type"))
	}
	if recorder.Body.String() != "%PDF" {
		t.Fatalf("expected pdf body, got %q", recorder.Body.String())
	}
}

func TestHandlerReturnsBadRequestForInvalidPeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := authenticatedReportsRouter(&fakeReportService{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/financial?startAt=bad&endAt=bad", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandlerReturnsInternalServerErrorWhenReportFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	query := "?startAt=2026-06-01T00:00:00Z&endAt=2026-06-30T23:59:59Z"
	router := authenticatedReportsRouter(&fakeReportService{err: errors.New("failed")})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/financial"+query, nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
}

func authenticatedReportsRouter(service *fakeReportService) *gin.Engine {
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("authenticated_user", authdomain.AuthenticatedUser{UserID: "authenticated-user"})
		ctx.Next()
	})
	NewHandler(service).RegisterForTest(router)
	return router
}

func (handler Handler) RegisterForTest(router *gin.Engine) {
	router.GET("/reports/financial", handler.GetFinancialReport)
	router.GET("/reports/financial/pdf", handler.DownloadFinancialPDF)
}

type fakeReportService struct {
	input  reportports.FinancialReportInput
	report reportports.FinancialReportDTO
	file   reportports.PDFFile
	err    error
}

func (service *fakeReportService) GetFinancialReport(ctx context.Context, input reportports.FinancialReportInput) (reportports.FinancialReportDTO, error) {
	service.input = input
	if service.err != nil {
		return reportports.FinancialReportDTO{}, service.err
	}
	return service.report, nil
}

func (service *fakeReportService) GenerateFinancialPDF(ctx context.Context, input reportports.FinancialReportInput) (reportports.PDFFile, error) {
	service.input = input
	if service.err != nil {
		return reportports.PDFFile{}, service.err
	}
	return service.file, nil
}
