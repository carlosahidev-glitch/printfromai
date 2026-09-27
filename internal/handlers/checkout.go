// Package handlers manages all HTTP request processing and template rendering.
package handlers

import (
	"log"
	"net/http"

	"github.com/printfromai/core/internal/services"
	"github.com/printfromai/core/web/templates/pages"
)

// CheckoutView renders the secure, two-column checkout page.
// It retrieves the user's active cart from Firestore to display the order summary.
func CheckoutView(w http.ResponseWriter, r *http.Request) {
	isHTMX := r.Header.Get("HX-Request") == "true"

	// TODO: Replace this hardcoded UID with the actual UID from our upcoming Auth Middleware
	// uid := r.Context().Value("user_uid").(string)
	uid := "anonymous_guest_123"

	// Fetch the cart from Firestore
	cart, err := services.GetCart(uid)
	if err != nil || cart == nil || len(cart.Items) == 0 {
		// If cart is empty or missing, bounce them back to the home page
		if isHTMX {
			w.Header().Set("HX-Redirect", "/")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	component := pages.CheckoutIndex(*cart, isHTMX)
	w.Header().Set("Content-Type", "text/html")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering checkout page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ProcessCheckout handles the shipping form submission via HTMX.
// It generates a Stripe PaymentIntent and returns the credit card input UI as a partial.
func ProcessCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// uid := r.Context().Value("user_uid").(string)
	uid := "anonymous_guest_123"

	// Capture the user's email for the Stripe receipt
	email := r.FormValue("email")

	// Fetch the cart to calculate the exact, un-tampered server-side total
	cart, err := services.GetCart(uid)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusBadRequest)
		return
	}

	// Call our Stripe service to generate the secure Client Secret
	clientSecret, err := services.CreatePaymentIntent(cart, uid, email)
	if err != nil {
		log.Printf("Stripe intent error: %v", err)
		http.Error(w, "<div class='text-red-500'>Payment gateway error. Please try again.</div>", http.StatusInternalServerError)
		return
	}

	// Return ONLY the Stripe Elements HTML snippet to HTMX, which will smoothly swap it into the page
	component := pages.StripePaymentPartial(clientSecret)
	w.Header().Set("Content-Type", "text/html")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering Stripe partial: %v", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}
