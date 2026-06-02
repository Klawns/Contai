package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"contai/internal/account/app/ports"
	"contai/internal/account/domain"
	authdomain "contai/internal/auth/domain"
	financedomain "contai/internal/finance/domain"
	userdomain "contai/internal/users/domain"

	"github.com/gin-gonic/gin"
)

func TestHandlerRequiresAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewHandler(&fakeAccountService{}).RegisterForTest(router)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/accounts", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", recorder.Code)
	}
}

func TestHandlerCreateAccountUsesAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAccountService{}
	router := authenticatedAccountRouter(service)
	body := bytes.NewBufferString(`{"name":"Checking","type":"checking","initialBalance":1500,"bankIconId":"bank_1"}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/accounts", body)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.createInput.UserID != "authenticated-user" {
		t.Fatalf("expected authenticated user id, got %s", service.createInput.UserID)
	}
	if service.createInput.InitialBalance.Cents() != 1500 {
		t.Fatalf("expected cents from request, got %d", service.createInput.InitialBalance.Cents())
	}
	if service.createInput.IncludeInDashboardTotal != nil {
		t.Fatal("expected omitted dashboard total flag to stay nil in service input")
	}
}

func TestHandlerCreateAccountAcceptsDashboardTotalFalse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAccountService{}
	router := authenticatedAccountRouter(service)
	body := bytes.NewBufferString(`{"name":"Checking","type":"checking","initialBalance":1500,"bankIconId":"bank_1","includeInDashboardTotal":false}`)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/accounts", body)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if service.createInput.IncludeInDashboardTotal == nil || *service.createInput.IncludeInDashboardTotal {
		t.Fatalf("expected false dashboard total flag, got %#v", service.createInput.IncludeInDashboardTotal)
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"includeInDashboardTotal":false`)) {
		t.Fatalf("expected response to include dashboard total flag, got %s", recorder.Body.String())
	}
}

func TestHandlerTotalBalanceResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAccountService{total: financedomain.NewMoney(2500)}
	router := authenticatedAccountRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/accounts/total-balance", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if recorder.Body.String() != `{"totalBalance":2500}` {
		t.Fatalf("expected camelCase total balance response, got %s", recorder.Body.String())
	}
}

func TestHandlerDeleteAccountReturnsNoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAccountService{}
	router := authenticatedAccountRouter(service)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, "/accounts/account-id", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
	if service.inactivateInput.AccountID != "account-id" {
		t.Fatalf("expected account id from path, got %s", service.inactivateInput.AccountID)
	}
}

func authenticatedAccountRouter(service *fakeAccountService) *gin.Engine {
	router := gin.New()
	router.Use(func(ctx *gin.Context) {
		ctx.Set("authenticated_user", authdomain.AuthenticatedUser{UserID: "authenticated-user"})
		ctx.Next()
	})
	NewHandler(service).RegisterForTest(router)
	return router
}

func (handler Handler) RegisterForTest(router *gin.Engine) {
	router.GET("/accounts", handler.ListAccounts)
	router.POST("/accounts", handler.CreateAccount)
	router.GET("/accounts/total-balance", handler.GetTotalBalance)
	router.PATCH("/accounts/:accountID", handler.UpdateAccount)
	router.DELETE("/accounts/:accountID", handler.DeleteAccount)
}

type fakeAccountService struct {
	err             error
	total           financedomain.Money
	createInput     ports.CreateAccountInput
	updateInput     ports.UpdateAccountInput
	inactivateInput ports.InactivateAccountInput
}

func (service *fakeAccountService) CreateAccount(ctx context.Context, input ports.CreateAccountInput) (ports.AccountDTO, error) {
	service.createInput = input
	if service.err != nil {
		return ports.AccountDTO{}, service.err
	}
	account := fakeAccountDTO(input.UserID)
	if input.IncludeInDashboardTotal != nil {
		account.IncludeInDashboardTotal = *input.IncludeInDashboardTotal
	}
	return account, nil
}

func (service *fakeAccountService) ListAccounts(ctx context.Context, input ports.ListAccountsInput) ([]ports.AccountDTO, error) {
	if service.err != nil {
		return nil, service.err
	}
	return []ports.AccountDTO{fakeAccountDTO(input.UserID)}, nil
}

func (service *fakeAccountService) FindActiveAccountsByUserID(ctx context.Context, userID userdomain.UserID) ([]ports.AccountDTO, error) {
	return service.ListAccounts(ctx, ports.ListAccountsInput{UserID: userID})
}

func (service *fakeAccountService) UpdateAccount(ctx context.Context, input ports.UpdateAccountInput) (ports.AccountDTO, error) {
	service.updateInput = input
	if service.err != nil {
		return ports.AccountDTO{}, service.err
	}
	return fakeAccountDTO(input.UserID), nil
}

func (service *fakeAccountService) InactivateAccount(ctx context.Context, input ports.InactivateAccountInput) error {
	service.inactivateInput = input
	if errors.Is(service.err, domain.ErrAccountNotFound) {
		return service.err
	}
	return service.err
}

func (service *fakeAccountService) GetTotalBalance(ctx context.Context, input ports.GetTotalBalanceInput) (financedomain.Money, error) {
	if service.err != nil {
		return 0, service.err
	}
	return service.total, nil
}

func fakeAccountDTO(userID userdomain.UserID) ports.AccountDTO {
	now := time.Now()
	return ports.AccountDTO{
		ID:                      "account-id",
		UserID:                  userID,
		Name:                    "Checking",
		Type:                    domain.AccountTypeChecking,
		InitialBalance:          financedomain.NewMoney(1500),
		CurrentBalance:          financedomain.NewMoney(1500),
		BankIconID:              "bank_1",
		IncludeInDashboardTotal: true,
		Status:                  domain.AccountStatusActive,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
}
