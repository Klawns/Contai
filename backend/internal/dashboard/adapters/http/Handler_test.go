package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	accountdomain "contai/internal/account/domain"
	authdomain "contai/internal/auth/domain"
	categorydomain "contai/internal/category/domain"
	"contai/internal/dashboard/app/ports"
	"contai/internal/dashboard/domain"
	financedomain "contai/internal/finance/domain"
	transactiondomain "contai/internal/transactions/domain"

	"github.com/gin-gonic/gin"
)

func TestHandlerRequiresAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(&fakeDashboardService{}).RegisterForTest(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/dashboard/monthly?startAt=2026-01-01T00:00:00Z&endAt=2026-01-31T23:59:59Z", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestHandlerUsesAuthenticatedUserAndParsesPeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeDashboardService{}
	router := authenticatedDashboardRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/dashboard/monthly?startAt=2026-01-01T00:00:00Z&endAt=2026-01-31T23:59:59Z", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.input.UserID != "authenticated-user" {
		t.Fatalf("expected authenticated user id, got %s", service.input.UserID)
	}
	if service.input.Period.StartAt.Format(time.RFC3339) != "2026-01-01T00:00:00Z" {
		t.Fatalf("expected parsed startAt, got %s", service.input.Period.StartAt.Format(time.RFC3339))
	}
}

func TestHandlerSerializesCamelCaseMoneyInCents(t *testing.T) {
	gin.SetMode(gin.TestMode)
	accountID := accountdomain.AccountID("account-id")
	categoryID := categorydomain.CategoryID("category-id")
	sourceAccountID := accountdomain.AccountID("source-account-id")
	destinationAccountID := accountdomain.AccountID("destination-account-id")
	occurredAt := time.Date(2026, 1, 5, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 1, 5, 12, 1, 0, 0, time.UTC)
	service := &fakeDashboardService{dashboard: ports.MonthlyDashboardDTO{
		UserID:             "authenticated-user",
		Period:             mustPeriod(t),
		TotalBalance:       financedomain.NewMoney(10000),
		MonthlyIncome:      financedomain.NewMoney(7000),
		MonthlyExpense:     financedomain.NewMoney(2500),
		MonthlyTransferIn:  financedomain.NewMoney(1200),
		MonthlyTransferOut: financedomain.NewMoney(1200),
		MonthlyNetBalance:  financedomain.NewMoney(4500),
		AccountBalances: []ports.AccountBalanceDTO{
			{AccountID: accountID, Name: "Checking", Type: accountdomain.AccountTypeChecking, Balance: financedomain.NewMoney(10000), BankIconID: "bank_1"},
		},
		ExpensesByCategory: []ports.CategoryExpenseDTO{
			{CategoryID: categoryID, Name: "Groceries", Color: "#2563EB", Icon: "shopping-cart", Total: financedomain.NewMoney(2500)},
		},
		RecentTransactions: []ports.TransactionDTO{
			{
				ID:                   "transaction-id",
				UserID:               "authenticated-user",
				Type:                 transactiondomain.TransactionTypeTransfer,
				Description:          "Transfer",
				Amount:               financedomain.NewMoney(2500),
				OccurredAt:           occurredAt,
				SourceAccountID:      &sourceAccountID,
				DestinationAccountID: &destinationAccountID,
				Status:               transactiondomain.TransactionStatusActive,
				Note:                 "note",
				CreatedAt:            createdAt,
				UpdatedAt:            createdAt,
			},
		},
	}}
	router := authenticatedDashboardRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/dashboard/monthly?startAt=2026-01-01T00:00:00Z&endAt=2026-01-31T23:59:59Z", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, expected := range []string{
		`"userId":"authenticated-user"`,
		`"totalBalance":10000`,
		`"monthlyIncome":7000`,
		`"monthlyExpense":2500`,
		`"monthlyTransferIn":1200`,
		`"monthlyTransferOut":1200`,
		`"monthlyNetBalance":4500`,
		`"accountBalances":[{"accountId":"account-id","name":"Checking","type":"checking","balance":10000,"bankIconId":"bank_1"}]`,
		`"expensesByCategory":[{"categoryId":"category-id","name":"Groceries","color":"#2563EB","icon":"shopping-cart","total":2500}]`,
		`"recentTransactions":[{"id":"transaction-id","userId":"authenticated-user","type":"transfer","description":"Transfer","amount":2500,"occurredAt":"2026-01-05T12:00:00Z","accountId":null,"sourceAccountId":"source-account-id","destinationAccountId":"destination-account-id","categoryId":null,"status":"active","note":"note","removedAt":null,"createdAt":"2026-01-05T12:01:00Z","updatedAt":"2026-01-05T12:01:00Z"}]`,
	} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected response to contain %s, got %s", expected, body)
		}
	}
}

func TestHandlerReturnsBadRequestForMissingOrInvalidDates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := authenticatedDashboardRouter(&fakeDashboardService{})

	for _, path := range []string{
		"/dashboard/monthly?endAt=2026-01-31T23:59:59Z",
		"/dashboard/monthly?startAt=invalid&endAt=2026-01-31T23:59:59Z",
		"/dashboard/monthly?startAt=2026-02-01T00:00:00Z&endAt=2026-01-31T23:59:59Z",
	} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, path, nil)

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected 400 for %s, got %d", path, recorder.Code)
		}
	}
}

func authenticatedDashboardRouter(service *fakeDashboardService) *gin.Engine {
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("authenticated_user", authdomain.AuthenticatedUser{UserID: "authenticated-user"})
		ctx.Next()
	})
	NewHandler(service).RegisterForTest(router)
	return router
}

func (handler Handler) RegisterForTest(router *gin.Engine) {
	router.GET("/dashboard/monthly", handler.GetMonthlyDashboard)
}

type fakeDashboardService struct {
	input     ports.GetMonthlyDashboardInput
	dashboard ports.MonthlyDashboardDTO
	err       error
}

func (service *fakeDashboardService) GetMonthlyDashboard(ctx context.Context, input ports.GetMonthlyDashboardInput) (ports.MonthlyDashboardDTO, error) {
	service.input = input
	if service.err != nil {
		return ports.MonthlyDashboardDTO{}, service.err
	}
	if service.dashboard.UserID != "" {
		return service.dashboard, nil
	}
	return ports.MonthlyDashboardDTO{
		UserID:             input.UserID,
		Period:             input.Period,
		AccountBalances:    []ports.AccountBalanceDTO{},
		ExpensesByCategory: []ports.CategoryExpenseDTO{},
		RecentTransactions: []ports.TransactionDTO{},
	}, nil
}

func mustPeriod(t *testing.T) domain.Period {
	t.Helper()
	period, err := domain.NewPeriod(
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 1, 31, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("expected period, got %v", err)
	}
	return period
}
