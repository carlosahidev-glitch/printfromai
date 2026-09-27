package handlers

import (
	"net/http"

	"github.com/printfromai/core/web/templates/pages"
)

// Home handles requests to the root "/" URL.
func Home(w http.ResponseWriter, r *http.Request) {
	// Ensure we only match exactly "/", otherwise return 404.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Check if request is from HTMX
	isHTMX := r.Header.Get("HX-Request") == "true"

	// Render the templ component (it automatically handles the HTMX conditional logic inside)
	pages.HomeIndex(isHTMX).Render(r.Context(), w)
}
