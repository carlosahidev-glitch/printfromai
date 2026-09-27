package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/printfromai/core/internal/models"
	"github.com/printfromai/core/web/templates/pages"
)

// ProductPage handles requests to /product/{slug}
func ProductPage(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	// 1. Fetch from our in-memory Catalog (0ms latency lookup)
	product, exists := models.ProductCatalog[slug]
	if !exists {
		http.NotFound(w, r)
		return
	}

	// 2. Convert Variants to JSON for our instant client-side price calculator
	variantsJSON, err := json.Marshal(product.Variants)
	if err != nil {
		log.Printf("Error marshaling variants to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 3. Check for HTMX request
	isHTMX := r.Header.Get("HX-Request") == "true"

	// 4. Render the strongly-typed templ component, passing our full struct and JSON
	pages.ProductIndex(product, string(variantsJSON), isHTMX).Render(r.Context(), w)
}
