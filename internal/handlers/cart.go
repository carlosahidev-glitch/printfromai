package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/printfromai/core/internal/models"
	"github.com/printfromai/core/internal/services"
	// "github.com/printfromai/core/web/templates/partials" // We will uncomment this in the next step
)

// MockCartStore is a temporary in-memory database for our carts.
// In a future step, we will swap this out for a single line of Firestore code!
var MockCartStore = make(map[string]*models.Cart)

// AddToCart handles the HTMX form submission from the product page.
func AddToCart(w http.ResponseWriter, r *http.Request) {
	// 1. Ensure it's a POST request
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Parse Multipart Form (100 MB limit for high-res AI art)
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	// 3. Retrieve or Provision a Secure Session
	sessionID, err := services.GetOrCreateSession(w, r)
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// Fetch existing cart from memory or create a new one
	cart, exists := MockCartStore[sessionID]
	if !exists {
		cart = models.NewCart(sessionID)
		MockCartStore[sessionID] = cart
	}

	// 4. Extract raw frontend form values
	productSlug := r.FormValue("product")
	sizeValue := r.FormValue("size")
	variantIndexStr := r.FormValue("variant_index")
	qtyStr := r.FormValue("quantity")
	hasBackDesign := r.FormValue("has_back_design") == "true"
	instructions := r.FormValue("instructions")

	// 5. SERVER-SIDE SECURITY: Validate and Recalculate Price
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

	qty, _ := strconv.Atoi(qtyStr)

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
			// Verify the upcharge for this quantity securely from the backend
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

	// Calculate absolute final price
	finalItemPrice := basePrice + optionsUpcharge

	// 6. Handle File Uploads (Saving locally to /tmp for development)
	frontURL := saveUploadedFile(r, "front_file", sessionID)
	backURL := saveUploadedFile(r, "back_file", sessionID)
	additionalURL := saveUploadedFile(r, "additional_file", sessionID)

	// 7. Construct the Cart Item and save it
	item := models.CartItem{
		ID:                fmt.Sprintf("item_%d", time.Now().UnixNano()),
		ProductSlug:       product.Slug,
		ProductTitle:      product.Title,
		Size:              sizeValue,
		VariantIndex:      variantIdx,
		VariantName:       variant.Name,
		Quantity:          qty,
		Options:           options,
		FrontFileURL:      frontURL,
		BackFileURL:       backURL,
		AdditionalFileURL: additionalURL,
		Instructions:      instructions,
		ItemTotal:         finalItemPrice,
	}

	cart.Items = append(cart.Items, item)

	// Update Master Cart Totals
	cart.Subtotal = 0
	for _, i := range cart.Items {
		cart.Subtotal += i.ItemTotal
	}
	cart.Total = cart.Subtotal // No taxes or shipping added yet
	cart.UpdatedAt = time.Now()

	// 8. Render the HTMX Response (Slide-out Cart Drawer)
	// Parse the beautiful HTML partial we just created
	tmpl, err := template.ParseFiles("web/templates/partials/cart_drawer.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Error loading cart UI", http.StatusInternalServerError)
		return
	}

	// Execute the template, passing in our dynamically calculated 'cart' struct
	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, cart); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Error rendering cart UI", http.StatusInternalServerError)
		return
	}
}

// GetCart handles GET requests to fetch the current cart drawer.
func GetCart(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/partials/cart_drawer.html")
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// TODO: Later, we will fetch the user's actual cart items from the session/Firestore
	// and pass them into the template here.
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

// saveUploadedFile is a helper to stream form files to local disk for development.
func saveUploadedFile(r *http.Request, formKey string, sessionID string) string {
	file, header, err := r.FormFile(formKey)
	if err != nil {
		return "" // No file uploaded for this field
	}
	defer file.Close()

	// Ensure the directory exists
	uploadDir := filepath.Join("tmp", "uploads", sessionID)
	os.MkdirAll(uploadDir, os.ModePerm)

	// Make the filename secure and unique
	safeFilename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(header.Filename))
	filePath := filepath.Join(uploadDir, safeFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return ""
	}
	defer dst.Close()

	// Stream the binary data safely
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("Error saving file: %v", err)
		return ""
	}

	// Return a relative URL path so we can display it later
	return "/" + filePath
}
