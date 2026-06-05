package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	authdomain "contai/internal/auth/domain"
	reportports "contai/internal/reports/app/ports"
	reportservices "contai/internal/reports/app/services"

	"github.com/gin-gonic/gin"
)

func TestHandlerRequiresAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(&fakeReportService{}).RegisterForTest(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/accounts/pdf", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestHandlerDownloadsAccountsPDF(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeReportService{file: reportports.PDFFile{
		Filename: "contai-relatorio-contas.pdf",
		Content:  []byte("%PDF"),
	}}
	router := authenticatedReportsRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/accounts/pdf", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.input.UserID != "authenticated-user" {
		t.Fatalf("expected authenticated user id, got %s", service.input.UserID)
	}
	if recorder.Header().Get("Content-Type") != "application/pdf" {
		t.Fatalf("expected pdf content type, got %s", recorder.Header().Get("Content-Type"))
	}
	expectedDisposition := `attachment; filename="contai-relatorio-contas.pdf"`
	if recorder.Header().Get("Content-Disposition") != expectedDisposition {
		t.Fatalf("expected content disposition %q, got %q", expectedDisposition, recorder.Header().Get("Content-Disposition"))
	}
	if recorder.Body.String() != "%PDF" {
		t.Fatalf("expected pdf body, got %q", recorder.Body.String())
	}
}

func TestHandlerReturnsInternalServerErrorWhenReportFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := authenticatedReportsRouter(&fakeReportService{err: errors.New("render failed")})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/accounts/pdf", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", recorder.Code)
	}
}

func TestHandlerReturnsBadRequestForInvalidPeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := authenticatedReportsRouter(&fakeReportService{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/period/pdf?startAt=bad&endAt=bad", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandlerReturnsBadRequestForInvalidTransactionType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := authenticatedReportsRouter(&fakeReportService{
		err: reportservices.ErrReportTransactionTypeInvalid,
	})
	query := "?startAt=2026-06-01T00:00:00Z&endAt=2026-06-30T23:59:59Z&type=transfer"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/transactions/pdf"+query, nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandlerReturnsNotFoundForMissingAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := authenticatedReportsRouter(&fakeReportService{
		err: reportservices.ErrReportAccountNotFound,
	})
	query := "?startAt=2026-06-01T00:00:00Z&endAt=2026-06-30T23:59:59Z"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/account/account-id/pdf"+query, nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recorder.Code)
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
	router.GET("/reports/accounts/pdf", handler.DownloadAccountsPDF)
	router.GET("/reports/transactions/pdf", handler.DownloadTransactionsPDF)
	router.GET("/reports/period/pdf", handler.DownloadPeriodPDF)
	router.GET("/reports/monthly/pdf", handler.DownloadMonthlyPDF)
	router.GET("/reports/account/:accountID/pdf", handler.DownloadAccountPDF)
}

type fakeReportService struct {
	input             reportports.GenerateAccountsReportInput
	transactionsInput reportports.GenerateTransactionsReportInput
	periodInput       reportports.PeriodReportInput
	accountInput      reportports.GenerateAccountReportInput
	file              reportports.PDFFile
	err               error
}

func (service *fakeReportService) GenerateAccountsPDF(ctx context.Context, input reportports.GenerateAccountsReportInput) (reportports.PDFFile, error) {
	service.input = input
	if service.err != nil {
		return reportports.PDFFile{}, service.err
	}
	return service.file, nil
}

func (service *fakeReportService) GenerateTransactionsPDF(ctx context.Context, input reportports.GenerateTransactionsReportInput) (reportports.PDFFile, error) {
	service.transactionsInput = input
	if service.err != nil {
		return reportports.PDFFile{}, service.err
	}
	return service.file, nil
}

func (service *fakeReportService) GeneratePeriodPDF(ctx context.Context, input reportports.PeriodReportInput) (reportports.PDFFile, error) {
	service.periodInput = input
	if service.err != nil {
		return reportports.PDFFile{}, service.err
	}
	return service.file, nil
}

func (service *fakeReportService) GenerateMonthlyPDF(ctx context.Context, input reportports.PeriodReportInput) (reportports.PDFFile, error) {
	service.periodInput = input
	if service.err != nil {
		return reportports.PDFFile{}, service.err
	}
	return service.file, nil
}

func (service *fakeReportService) GenerateAccountPDF(ctx context.Context, input reportports.GenerateAccountReportInput) (reportports.PDFFile, error) {
	service.accountInput = input
	if service.err != nil {
		return reportports.PDFFile{}, service.err
	}
	return service.file, nil
}

func TestHandlerParsesAccountReportParams(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeReportService{file: reportports.PDFFile{
		Filename: "account.pdf",
		Content:  []byte("%PDF"),
	}}
	router := authenticatedReportsRouter(service)
	query := "?startAt=2026-06-01T00:00:00Z&endAt=2026-06-30T23:59:59Z"

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/reports/account/account-id/pdf"+query, nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.accountInput.UserID != "authenticated-user" || service.accountInput.AccountID != "account-id" {
		t.Fatalf("expected user and account params, got %#v", service.accountInput)
	}
	expectedStart := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	if !service.accountInput.StartAt.Equal(expectedStart) {
		t.Fatalf("expected parsed start date, got %s", service.accountInput.StartAt)
	}
}
