package services

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
)

// SessionCookieName is the key used to store the session ID in the browser.
const SessionCookieName = "pfa_session"

// SessionDuration sets how long a cart/session remains active before expiring.
// 7 days is standard for e-commerce guest carts.
const SessionDuration = 7 * 24 * time.Hour

// GenerateSecureID creates a highly secure, random 32-byte hex string.
// We use this instead of standard UUIDs to avoid pulling in third-party dependencies.
func GenerateSecureID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetSessionID retrieves the active session ID from the request's cookies.
// It returns an empty string if the cookie does not exist or has expired.
func GetSessionID(r *http.Request) string {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		// ErrNoCookie is returned when the cookie isn't found, which is normal for new visitors.
		return ""
	}
	return cookie.Value
}

// SetSessionCookie creates a secure HTTP-Only cookie and sends it to the browser.
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(SessionDuration),
		MaxAge:   int(SessionDuration.Seconds()),
		HttpOnly: true, // Crucial: Prevents JavaScript (and XSS attacks) from reading the cookie
		Secure:   true, // Crucial: Ensures cookie is only sent over HTTPS (or localhost in dev)
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// ClearSessionCookie removes the session cookie from the browser.
// We will call this after a successful Stripe payment or when a user logs out.
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Set expiration to the past to immediately delete
		MaxAge:   -1,              // Instructs the browser to delete the cookie
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// GetOrCreateSession seamlessly retrieves an existing session or provisions a new one.
// This is the primary function our cart handler will call.
func GetOrCreateSession(w http.ResponseWriter, r *http.Request) (string, error) {
	sessionID := GetSessionID(r)
	if sessionID == "" {
		newID, err := GenerateSecureID()
		if err != nil {
			return "", err
		}
		sessionID = newID
		SetSessionCookie(w, sessionID)
	}
	return sessionID, nil
}
