package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"contai/internal/auth/app/contracts"
	authdomain "contai/internal/auth/domain"
	userports "contai/internal/users/app/ports"
	userdomain "contai/internal/users/domain"

	"github.com/gin-gonic/gin"
)

func TestHandler_CreateUserSetsCookieWithoutReturningToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(
		fakeHTTPAuthService{},
		fakeHTTPUserService{},
		NewCookieService(false),
	)
	registerTestRoutes(router, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(`{"name":"John Doe","email":"john@example.com","password":"secret"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "access-token") {
		t.Fatalf("expected response without token, got %s", recorder.Body.String())
	}
	cookie := findCookie(t, recorder.Result().Cookies(), AccessCookieName)
	if cookie.Value != "access-token" || !cookie.HttpOnly {
		t.Fatalf("expected httponly access cookie, got %#v", cookie)
	}
}

func TestHandler_LoginSetsCookieWithoutReturningToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(fakeHTTPAuthService{}, fakeHTTPUserService{}, NewCookieService(false))
	registerTestRoutes(router, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"john@example.com","password":"secret"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if strings.Contains(recorder.Body.String(), "access-token") {
		t.Fatalf("expected response without token, got %s", recorder.Body.String())
	}
	cookie := findCookie(t, recorder.Result().Cookies(), AccessCookieName)
	if cookie.Value != "access-token" {
		t.Fatalf("expected access cookie, got %#v", cookie)
	}
}

func TestHandler_LogoutClearsCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(fakeHTTPAuthService{}, fakeHTTPUserService{}, NewCookieService(false))
	registerTestRoutes(router, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", recorder.Code)
	}
	cookie := findCookie(t, recorder.Result().Cookies(), AccessCookieName)
	if cookie.MaxAge != -1 {
		t.Fatalf("expected cleared cookie, got %#v", cookie)
	}
}

func TestHandler_MeRequiresValidCookie(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(fakeHTTPAuthService{}, fakeHTTPUserService{}, NewCookieService(false))
	registerTestRoutes(router, handler)

	unauthorized := httptest.NewRecorder()
	router.ServeHTTP(unauthorized, httptest.NewRequest(http.MethodGet, "/api/auth/me", nil))
	if unauthorized.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 without cookie, got %d", unauthorized.Code)
	}

	authorized := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	request.AddCookie(&http.Cookie{Name: AccessCookieName, Value: "access-token"})

	router.ServeHTTP(authorized, request)

	if authorized.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d with body %s", authorized.Code, authorized.Body.String())
	}
	if !strings.Contains(authorized.Body.String(), "john@example.com") {
		t.Fatalf("expected current user, got %s", authorized.Body.String())
	}
}

func TestHandler_CreateUserReturnsWeakPasswordError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewHandler(
		fakeHTTPAuthService{},
		fakeHTTPUserService{createErr: userdomain.ErrUserPasswordTooWeak},
		NewCookieService(false),
	)
	registerTestRoutes(router, handler)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(`{"name":"John Doe","email":"john@example.com","password":"short"}`))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d with body %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "weak password") {
		t.Fatalf("expected weak password error, got %s", recorder.Body.String())
	}
}

func registerTestRoutes(router *gin.Engine, handler Handler) {
	router.POST("/api/users", handler.CreateUser)
	router.POST("/api/auth/login", handler.Login)
	router.POST("/api/auth/logout", handler.Logout)
	router.GET("/api/auth/me", handler.AuthMiddleware(), handler.Me)
}

type fakeHTTPAuthService struct{}

func (service fakeHTTPAuthService) Login(ctx context.Context, input contracts.LoginInput) (authdomain.AuthenticatedUser, contracts.AuthTokens, error) {
	return authdomain.AuthenticatedUser{
			UserID: userdomain.UserID("user-id"),
			Email:  "john@example.com",
			Status: userdomain.UserStatusActive,
		}, contracts.AuthTokens{
			AccessToken: "access-token",
			AccessClaims: authdomain.AuthClaims{
				Type:      authdomain.AuthTokenTypeAccess,
				UserID:    userdomain.UserID("user-id"),
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(time.Hour),
			},
		}, nil
}

func (service fakeHTTPAuthService) Logout(ctx context.Context) error {
	return nil
}

func (service fakeHTTPAuthService) ValidateAccessToken(ctx context.Context, accessToken string) (authdomain.AuthenticatedUser, error) {
	return authdomain.AuthenticatedUser{
		UserID: userdomain.UserID("user-id"),
		Email:  "john@example.com",
		Status: userdomain.UserStatusActive,
	}, nil
}

type fakeHTTPUserService struct {
	createErr error
}

func (service fakeHTTPUserService) CreateUser(ctx context.Context, input userports.CreateUserInput) (userports.UserDTO, error) {
	if service.createErr != nil {
		return userports.UserDTO{}, service.createErr
	}

	return userports.UserDTO{
		ID:        userdomain.UserID("user-id"),
		Name:      input.Name,
		Email:     input.Email,
		Status:    userdomain.UserStatusActive,
		CreatedAt: time.Now(),
	}, nil
}

func (service fakeHTTPUserService) GetUserByID(ctx context.Context, userID userdomain.UserID) (userports.UserDTO, error) {
	return userports.UserDTO{}, nil
}
