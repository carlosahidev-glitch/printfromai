// Package handlers manages all HTTP request processing and template rendering.
package handlers

import (
	"log"
	"net/http"

	"github.com/printfromai/core/web/templates/pages"
)

// LoginView serves the sleek authentication gateway.
// It detects if HTMX is requesting a partial swap or if the browser needs a full page load.
func LoginView(w http.ResponseWriter, r *http.Request) {
	isHTMX := r.Header.Get("HX-Request") == "true"

	// Render the Auth UI component
	component := pages.AuthIndex(isHTMX)

	w.Header().Set("Content-Type", "text/html")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering login page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
