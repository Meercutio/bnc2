package auth

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const SessionCookieName = "bc_session"

func TokenFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}

	if h := strings.TrimSpace(r.Header.Get("Authorization")); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	}

	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookie.Value)
}

func SetSessionCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration) {
	if w == nil {
		return
	}

	maxAge := int(ttl / time.Second)
	if maxAge < 0 {
		maxAge = 0
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   RequestIsSecure(r),
	})
}

func ClearSessionCookie(w http.ResponseWriter, r *http.Request) {
	if w == nil {
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   RequestIsSecure(r),
	})
}

func RequestIsSecure(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	if strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")), "https") {
		return true
	}
	return strings.EqualFold(strings.TrimSpace(r.Header.Get("X-Forwarded-Ssl")), "on")
}

func AllowedWebSocketOrigin(r *http.Request) bool {
	if r == nil {
		return false
	}

	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil || u.Host == "" {
		return false
	}

	allowedHosts := []string{r.Host}
	if forwardedHost := strings.TrimSpace(r.Header.Get("X-Forwarded-Host")); forwardedHost != "" {
		allowedHosts = append(allowedHosts, forwardedHost)
	}

	hostAllowed := false
	for _, host := range allowedHosts {
		if strings.EqualFold(u.Host, host) {
			hostAllowed = true
			break
		}
	}
	if !hostAllowed {
		return false
	}

	if RequestIsSecure(r) {
		return strings.EqualFold(u.Scheme, "https")
	}
	return strings.EqualFold(u.Scheme, "http")
}
