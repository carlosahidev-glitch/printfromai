// Package services handles external integrations like Firebase, Stripe, and AWS.
package services

import (
	"fmt"
	"log"
	"os"

	"github.com/printfromai/core/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
)

// InitStripe configures the Stripe SDK using your secret key.
// In Cloud Run, this should be injected securely via Secret Manager.
func InitStripe() {
	key := os.Getenv("STRIPE_SECRET_KEY")
	if key == "" {
		log.Println("⚠️  STRIPE_SECRET_KEY is missing. Payment creation will fail.")
	}
	stripe.Key = key
	log.Println("💳 Stripe Service Initialized")
}

// CreatePaymentIntent calculates the final total exclusively on the server to prevent
// frontend manipulation, then requests a secure payment token from Stripe.
func CreatePaymentIntent(cart *models.Cart, uid string, email string) (string, error) {
	// 1. Validate the cart total
	if cart == nil || cart.Total <= 0 {
		return "", fmt.Errorf("invalid cart total")
	}

	// Stripe requires amounts in the smallest currency unit (cents)
	amountInCents := int64(cart.Total * 100)

	// 2. Define the PaymentIntent parameters
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true), // Automatically enables Apple Pay, Google Pay, Cards
		},
	}

	// If the user provided an email during checkout, attach it for Stripe receipts
	if email != "" {
		params.ReceiptEmail = stripe.String(email)
	}

	// Attach metadata so our Webhook handler knows exactly whose cart to finalize
	params.AddMetadata("firebase_uid", uid)

	// 3. Create the Intent via Stripe API
	pi, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Stripe error: %v", err)
		return "", fmt.Errorf("failed to initialize payment gateway")
	}

	// 4. Return the Client Secret (The frontend needs this to render the card input securely)
	return pi.ClientSecret, nil
}
