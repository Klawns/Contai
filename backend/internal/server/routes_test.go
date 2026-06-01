package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutesIncludesAuthenticatedAccountRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	registerRoutes(router, dependencies{})

	routes := router.Routes()
	expected := map[string]string{
		http.MethodGet + " /api/accounts":               "",
		http.MethodPost + " /api/accounts":              "",
		http.MethodGet + " /api/accounts/total-balance": "",
		http.MethodPatch + " /api/accounts/:accountID":  "",
		http.MethodDelete + " /api/accounts/:accountID": "",
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			delete(expected, key)
		}
	}

	if len(expected) > 0 {
		t.Fatalf("expected account routes to be registered, missing %#v", expected)
	}

	unauthenticatedPaths := []string{"/api/accounts", "/api/accounts/total-balance"}
	for _, path := range unauthenticatedPaths {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, path, nil)

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected %s to require authentication, got %d", path, recorder.Code)
		}
	}
}

func TestRegisterRoutesIncludesAuthenticatedTransactionRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	registerRoutes(router, dependencies{})

	routes := router.Routes()
	expected := map[string]string{
		http.MethodGet + " /api/transactions":                   "",
		http.MethodPost + " /api/transactions/income":           "",
		http.MethodPost + " /api/transactions/expense":          "",
		http.MethodPost + " /api/transactions/transfer":         "",
		http.MethodPatch + " /api/transactions/:transactionID":  "",
		http.MethodDelete + " /api/transactions/:transactionID": "",
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			delete(expected, key)
		}
	}

	if len(expected) > 0 {
		t.Fatalf("expected transaction routes to be registered, missing %#v", expected)
	}

	unauthenticatedPaths := []string{"/api/transactions", "/api/transactions/income"}
	for _, path := range unauthenticatedPaths {
		recorder := httptest.NewRecorder()
		method := http.MethodGet
		if path == "/api/transactions/income" {
			method = http.MethodPost
		}
		request := httptest.NewRequest(method, path, nil)

		router.ServeHTTP(recorder, request)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected %s to require authentication, got %d", path, recorder.Code)
		}
	}
}
