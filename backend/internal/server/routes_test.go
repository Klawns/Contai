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

func TestRegisterRoutesDoesNotRegisterLegacyCommitmentRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	registerRoutes(router, dependencies{})

	legacyRoutes := map[string]struct{}{
		http.MethodGet + " /api/payables":                  {},
		http.MethodPost + " /api/payables":                 {},
		http.MethodPatch + " /api/payables/:id":            {},
		http.MethodPatch + " /api/payables/:id/pay":        {},
		http.MethodPatch + " /api/payables/:id/cancel":     {},
		http.MethodGet + " /api/receivables":               {},
		http.MethodPost + " /api/receivables":              {},
		http.MethodPatch + " /api/receivables/:id":         {},
		http.MethodPatch + " /api/receivables/:id/receive": {},
		http.MethodPatch + " /api/receivables/:id/cancel":  {},
	}

	for _, route := range router.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := legacyRoutes[key]; ok {
			t.Fatalf("expected legacy commitment route to be inactive, got %s", key)
		}
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/payables", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected legacy commitment route to be inactive, got %d", recorder.Code)
	}
}

func TestRegisterRoutesIncludesAuthenticatedCreditCardRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	registerRoutes(router, dependencies{})

	routes := router.Routes()
	expected := map[string]string{
		http.MethodGet + " /api/credit-cards":                               "",
		http.MethodPost + " /api/credit-cards":                              "",
		http.MethodPatch + " /api/credit-cards/:cardID":                     "",
		http.MethodPatch + " /api/credit-cards/:cardID/inactivate":          "",
		http.MethodGet + " /api/credit-cards/:cardID/purchases":             "",
		http.MethodPost + " /api/credit-cards/:cardID/purchases":            "",
		http.MethodPatch + " /api/credit-card-purchases/:purchaseID/cancel": "",
		http.MethodGet + " /api/credit-cards/:cardID/invoices":              "",
		http.MethodGet + " /api/credit-card-invoices/:invoiceID":            "",
		http.MethodPatch + " /api/credit-card-invoices/:invoiceID/close":    "",
		http.MethodPatch + " /api/credit-card-invoices/:invoiceID/pay":      "",
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			delete(expected, key)
		}
	}

	if len(expected) > 0 {
		t.Fatalf("expected credit card routes to be registered, missing %#v", expected)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/credit-cards", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected credit card route to require authentication, got %d", recorder.Code)
	}
}

func TestRegisterRoutesIncludesAuthenticatedDashboardRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	registerRoutes(router, dependencies{})

	routes := router.Routes()
	expected := map[string]string{
		http.MethodGet + " /api/dashboard/monthly":        "",
		http.MethodGet + " /api/dashboard/monthly-series": "",
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			delete(expected, key)
		}
	}

	if len(expected) > 0 {
		t.Fatalf("expected dashboard routes to be registered, missing %#v", expected)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/dashboard/monthly-series", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected dashboard route to require authentication, got %d", recorder.Code)
	}
}

func TestRegisterRoutesIncludesAuthenticatedReportRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	registerRoutes(router, dependencies{})

	routes := router.Routes()
	expected := map[string]string{
		http.MethodGet + " /api/reports/financial":     "",
		http.MethodGet + " /api/reports/financial/pdf": "",
	}
	legacy := map[string]struct{}{
		http.MethodGet + " /api/reports/accounts/pdf":           {},
		http.MethodGet + " /api/reports/transactions/pdf":       {},
		http.MethodGet + " /api/reports/period/pdf":             {},
		http.MethodGet + " /api/reports/monthly/pdf":            {},
		http.MethodGet + " /api/reports/account/:accountID/pdf": {},
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, ok := expected[key]; ok {
			delete(expected, key)
		}
		if _, ok := legacy[key]; ok {
			t.Fatalf("expected legacy report route to be inactive, got %s", key)
		}
	}

	if len(expected) > 0 {
		t.Fatalf("expected report routes to be registered, missing %#v", expected)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/reports/financial", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected report route to require authentication, got %d", recorder.Code)
	}
}
