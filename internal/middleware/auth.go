// Package middleware contains HTTP interceptors for authentication, authorization, and logging.
package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/printfromai/core/internal/services"
)

// ContextKey is a custom string type to prevent context key collisions across packages.
type ContextKey string

// UserUIDKey is the context key used to store the verified Firebase UID.
const UserUIDKey ContextKey = "user_uid"

// RequireAuth intercepts requests, securely verifies the Firebase JWT, and injects the UID into the context.
// If the user is missing a token or it is invalid, the request is blocked.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			handleUnauthorized(w, r)
			return
		}

		// 2. Extract the raw JWT token string
		idToken := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Verify the token securely using our Firebase Admin SDK service
		token, err := services.FB.Auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			log.Printf("⚠️ Unauthorized access attempt or expired token: %v", err)
			handleUnauthorized(w, r)
			return
		}

		// 4. Inject the validated Firebase UID into the request context
		ctx := context.WithValue(r.Context(), UserUIDKey, token.UID)

		// 5. Pass the request to the next handler with our newly enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// handleUnauthorized sends the appropriate redirect based on the request type.
func handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		// HTMX Magic: If an HTMX background request hits a protected route without auth,
		// we send the HX-Redirect header. HTMX will instantly redirect the browser to login.
		w.Header().Set("HX-Redirect", "/login")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Standard full-page browser redirect
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
