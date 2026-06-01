package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCookieService_SetAccessCookie(t *testing.T) {
	recorder := httptest.NewRecorder()
	expiresAt := time.Date(2026, 6, 1, 12, 30, 0, 0, time.UTC)

	NewCookieService(true).SetAccessCookie(recorder, "token", expiresAt)

	cookie := findCookie(t, recorder.Result().Cookies(), AccessCookieName)
	if cookie.Value != "token" || !cookie.HttpOnly || !cookie.Secure {
		t.Fatalf("expected secure httponly token cookie, got %#v", cookie)
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
}

func TestCookieService_ClearAccessCookie(t *testing.T) {
	recorder := httptest.NewRecorder()

	NewCookieService(false).ClearAccessCookie(recorder)

	cookie := findCookie(t, recorder.Result().Cookies(), AccessCookieName)
	if cookie.Value != "" || cookie.MaxAge != -1 || !cookie.HttpOnly || cookie.Secure {
		t.Fatalf("expected cleared non-secure httponly cookie, got %#v", cookie)
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("expected SameSite=Lax, got %v", cookie.SameSite)
	}
}

func findCookie(t *testing.T, cookies []*http.Cookie, name string) *http.Cookie {
	t.Helper()

	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	t.Fatalf("expected cookie %s", name)
	return nil
}
