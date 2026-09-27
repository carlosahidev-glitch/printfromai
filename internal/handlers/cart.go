package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/printfromai/core/internal/models"
	"github.com/printfromai/core/internal/services"
)

// verifyAuth is a helper to securely validate the Firebase JWT sent by HTMX
func verifyAuth(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Verify the token securely with Google's servers
	token, err := services.FB.Auth.VerifyIDToken(context.Background(), tokenString)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}
	return token.UID, nil
}

// AddToCart handles the HTMX form submission containing GCS object paths.
func AddToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Authenticate the User
	uid, err := verifyAuth(r)
	if err != nil {
		log.Printf("Auth error: %v", err)
		http.Error(w, "Unauthorized: Please wait for session to initialize", http.StatusUnauthorized)
		return
	}

	// 2. Parse Standard Form Data (No multipart parsing needed!)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract Lightweight Data
	productSlug := r.FormValue("product")
	sizeValue := r.FormValue("size")
	variantIndexStr := r.FormValue("variant_index")
	qtyStr := r.FormValue("quantity")
	hasBackDesign := r.FormValue("has_back_design") == "true"
	instructions := r.FormValue("instructions")

	frontGCSPath := r.FormValue("front_file_path")
	backGCSPath := r.FormValue("back_file_path")
	additionalGCSPath := r.FormValue("additional_file_path")

	// 3. SERVER-SIDE SECURITY: Validate and Recalculate Price
	product, ok := models.ProductCatalog[productSlug]
	if !ok {
		http.Error(w, "Product not found in catalog", http.StatusBadRequest)
		return
	}

	variantIdx, _ := strconv.Atoi(variantIndexStr)
	if variantIdx < 0 || variantIdx >= len(product.Variants) {
		http.Error(w, "Invalid variant selected", http.StatusBadRequest)
		return
	}
	variant := product.Variants[variantIdx]

	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	// Find the secure base price for this size and quantity
	var basePrice float64 = 0
	pricingTiers := variant.Pricing[sizeValue]
	for _, tier := range pricingTiers {
		if tier.Qty == qty {
			if hasBackDesign {
				basePrice = tier.FrontBack
			} else {
				basePrice = tier.FrontOnly
			}
			break
		}
	}

	if basePrice == 0 {
		http.Error(w, "Invalid quantity or size combination", http.StatusBadRequest)
		return
	}

	// Parse Dynamic Options (e.g., Corner Cuts) and securely add their specific upcharges
	options := make(map[string]string)
	var optionsUpcharge float64 = 0

	for _, opt := range variant.Options {
		formKey := "option_" + opt.Title
		selectedVal := r.FormValue(formKey)
		if selectedVal != "" {
			options[opt.Title] = selectedVal
			for _, choice := range opt.Choices {
				if choice.Value == selectedVal {
					if charge, exists := choice.Upcharges[qty]; exists {
						optionsUpcharge += charge
					}
					break
				}
			}
		}
	}

	finalItemPrice := basePrice + optionsUpcharge

	// 4. Construct the Cart Item
	item := models.CartItem{
		ID:                fmt.Sprintf("item_%d", time.Now().UnixNano()),
		ProductSlug:       product.Slug,
		ProductTitle:      product.Title,
		Size:              sizeValue,
		VariantIndex:      variantIdx,
		VariantName:       variant.Name,
		Quantity:          qty,
		Options:           options,
		FrontFileURL:      frontGCSPath, // Saving the cloud path instead of a local /tmp URL
		BackFileURL:       backGCSPath,
		AdditionalFileURL: additionalGCSPath,
		Instructions:      instructions,
		ItemTotal:         finalItemPrice,
	}

	// 5. Save to Firestore
	if err := services.SaveCartItem(uid, item); err != nil {
		log.Printf("Failed to save cart to Firestore: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 6. Fetch the updated cart to render the UI
	cart, err := services.GetCart(uid)
	if err != nil {
		log.Printf("Failed to fetch cart: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 7. Render the HTMX Response (Slide-out Cart Drawer)
	renderCartDrawer(w, cart)
}

// GetCart handles GET requests to fetch the current user's cart drawer on page load.
func GetCart(w http.ResponseWriter, r *http.Request) {
	uid, err := verifyAuth(r)
	if err != nil {
		// If not authenticated yet, return an empty shell to prevent UI errors
		renderCartDrawer(w, &models.Cart{})
		return
	}

	cart, err := services.GetCart(uid)
	if err != nil {
		// Document might not exist yet, which is fine
		cart = &models.Cart{}
	}

	renderCartDrawer(w, cart)
}

// RemoveFromCart processes HTMX DELETE requests using Go 1.22 path routing
func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid, err := verifyAuth(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Requires Go 1.22 mux routing: mux.HandleFunc("DELETE /cart/remove/{id}", handlers.RemoveFromCart)
	itemID := r.PathValue("id")
	if itemID == "" {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	if err := services.RemoveCartItem(uid, itemID); err != nil {
		log.Printf("Failed to remove item: %v", err)
		http.Error(w, "Failed to remove item", http.StatusInternalServerError)
		return
	}

	cart, err := services.GetCart(uid)
	if err != nil {
		cart = &models.Cart{}
	}

	renderCartDrawer(w, cart)
}

// renderCartDrawer is a helper function to parse and execute the HTML template
func renderCartDrawer(w http.ResponseWriter, cart *models.Cart) {
	tmpl, err := template.ParseFiles("web/templates/partials/cart_drawer.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Error loading cart UI", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, cart); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}
