package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestLimitBodyReturnsPayloadTooLarge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/limited", limitBody(4), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/limited", strings.NewReader("too large"))

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func TestRateLimiterReturnsTooManyRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	limiter := newRateLimiter(1, time.Minute)
	router.POST("/limited", limiter.Middleware(), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	first := httptest.NewRecorder()
	router.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/limited", nil))
	if first.Code != http.StatusNoContent {
		t.Fatalf("expected first request status 204, got %d", first.Code)
	}

	second := httptest.NewRecorder()
	router.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/limited", nil))
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("expected second request status 429, got %d with body %s", second.Code, second.Body.String())
	}
}

func TestSecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(securityHeaders(true))
	router.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	headers := recorder.Result().Header
	if headers.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options header, got %q", headers.Get("X-Content-Type-Options"))
	}
	if headers.Get("X-Frame-Options") != "DENY" {
		t.Fatalf("expected X-Frame-Options header, got %q", headers.Get("X-Frame-Options"))
	}
	if headers.Get("Referrer-Policy") == "" {
		t.Fatal("expected Referrer-Policy header")
	}
	if headers.Get("Strict-Transport-Security") == "" {
		t.Fatal("expected HSTS header in production")
	}
}
