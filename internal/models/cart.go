package models

import (
	"time"
)

// CartItem represents a single configured product in the user's shopping cart.
// We store the human-readable names (like ProductTitle and VariantName) alongside
// the exact indices so we can easily render the cart drawer UI without having
// to constantly look up the catalog in the HTML template.
type CartItem struct {
	ID                string            `json:"id" firestore:"id"`                               // Unique ID for this cart item
	ProductSlug       string            `json:"productSlug" firestore:"productSlug"`             // e.g., "business-cards"
	ProductTitle      string            `json:"productTitle" firestore:"productTitle"`           // e.g., "Business Cards"
	Size              string            `json:"size" firestore:"size"`                           // e.g., "2x3.5"
	VariantIndex      int               `json:"variantIndex" firestore:"variantIndex"`           // The index of the selected material
	VariantName       string            `json:"variantName" firestore:"variantName"`             // e.g., "Standard (16pt)"
	Quantity          int               `json:"quantity" firestore:"quantity"`                   // Number of prints requested
	Options           map[string]string `json:"options" firestore:"options"`                     // Chosen dynamic options (e.g., "Corner Cut" -> "round-1-8")
	FrontFileURL      string            `json:"frontFileUrl" firestore:"frontFileUrl"`           // Temp or permanent GCS link to front art
	BackFileURL       string            `json:"backFileUrl" firestore:"backFileUrl"`             // Temp or permanent GCS link to back art (empty if 4/0)
	AdditionalFileURL string            `json:"additionalFileUrl" firestore:"additionalFileUrl"` // Link to zip/pdf of extras
	Instructions      string            `json:"instructions" firestore:"instructions"`           // User's custom design notes
	RequiresHumanEdit bool              `json:"requiresHumanEdit" firestore:"requiresHumanEdit"` // True if user selected the $25 prepress upsell
	ItemTotal         float64           `json:"itemTotal" firestore:"itemTotal"`                 // Server-calculated true price for this item
}

// Cart represents the entire shopping session for a user.
// This struct will be stored in Firestore under the "carts" collection, keyed by the SessionID.
type Cart struct {
	SessionID string     `json:"sessionId" firestore:"sessionId"` // The UUID stored in the user's HTTPOnly cookie
	UserID    string     `json:"userId" firestore:"userId"`       // Firebase UID. Empty if this is a Guest cart.
	Items     []CartItem `json:"items" firestore:"items"`         // List of products in the cart
	Subtotal  float64    `json:"subtotal" firestore:"subtotal"`   // Sum of all ItemTotals
	Tax       float64    `json:"tax" firestore:"tax"`             // Calculated tax (if applicable)
	Shipping  float64    `json:"shipping" firestore:"shipping"`   // Shipping cost (0.00 for our Free Ground Shipping)
	Total     float64    `json:"total" firestore:"total"`         // Grand total to be sent to Stripe
	UpdatedAt time.Time  `json:"updatedAt" firestore:"updatedAt"` // Used to clean up abandoned guest carts after 7 days
}

// NewCart is a helper function to instantiate a fresh, empty cart.
func NewCart(sessionID string) *Cart {
	return &Cart{
		SessionID: sessionID,
		UserID:    "", // Defaults to empty (Guest)
		Items:     []CartItem{},
		Subtotal:  0.0,
		Tax:       0.0,
		Shipping:  0.0,
		Total:     0.0,
		UpdatedAt: time.Now(),
	}
}
