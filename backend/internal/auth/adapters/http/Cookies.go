package http

import (
	"net/http"
	"time"
)

const AccessCookieName = "contai_access"

type CookieService struct {
	secure bool
}

func NewCookieService(secure bool) CookieService {
	return CookieService{secure: secure}
}

func (service CookieService) SetAccessCookie(w http.ResponseWriter, accessToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessCookieName,
		Value:    accessToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   service.secure,
	})
}

func (service CookieService) ClearAccessCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0).UTC(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   service.secure,
	})
}
