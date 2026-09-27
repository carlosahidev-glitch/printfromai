// Package services handles external integrations like Firebase, Stripe, and AWS.
package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/printfromai/core/internal/models"
)

// FirebaseService holds our core Google Cloud/Firebase clients so we can reuse them.
type FirebaseService struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
	Storage   *storage.Client
	Bucket    string
}

// FB is the global instance of our Firebase services.
var FB *FirebaseService

// InitFirebase initializes Auth, Firestore, and Cloud Storage.
// In Google Cloud Run, it automatically picks up the default service account credentials.
func InitFirebase(projectID string, bucketName string) error {
	ctx := context.Background()

	// Initialize the Firebase App
	conf := &firebase.Config{
		ProjectID:     projectID,
		StorageBucket: bucketName,
	}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Initialize Firebase Auth Client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("error getting auth client: %v", err)
	}

	// Initialize Firestore Client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error getting firestore client: %v", err)
	}

	// Initialize Google Cloud Storage Client
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("error getting storage client: %v", err)
	}

	// Assign to global variable
	FB = &FirebaseService{
		App:       app,
		Auth:      authClient,
		Firestore: firestoreClient,
		Storage:   storageClient,
		Bucket:    bucketName,
	}

	log.Println("🔥 Firebase Services Initialized Successfully (Auth, Firestore, Storage)")
	return nil
}

// ==========================================
// GOOGLE CLOUD STORAGE (FILE UPLOADS)
// ==========================================

// GeneratePresignedUploadURL creates a secure, time-limited URL that allows the frontend
// to upload a file directly to Google Cloud Storage, bypassing our Go server's memory.
func (f *FirebaseService) GeneratePresignedUploadURL(uid string, filename string) (string, string, error) {
	// 1. Create a secure path. We quarantine uploads in a specific folder by User ID.
	// We append a timestamp to the filename to prevent overwriting.
	objectName := fmt.Sprintf("uploads/users/%s/%d_%s", uid, time.Now().Unix(), filename)

	// 2. Configure the presigned URL options.
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "PUT",                                              // The frontend must use a PUT request to upload the file
		Expires: time.Now().Add(15 * time.Minute),                   // The URL is only valid for 15 minutes
		Headers: []string{"Content-Type: application/octet-stream"}, // Enforce binary upload
	}

	// 3. Generate the URL using the default credentials
	url, err := f.Storage.Bucket(f.Bucket).SignedURL(objectName, opts)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate signed URL: %v", err)
	}

	// Return BOTH the Presigned URL (for frontend upload) AND the ObjectName (to save in Firestore)
	return url, objectName, nil
}

// ==========================================
// CART & FIRESTORE OPERATIONS
// ==========================================

// SaveCartItem securely adds a new item to the user's persistent cart in Firestore.
func SaveCartItem(uid string, item models.CartItem) error {
	ctx := context.Background()
	cartRef := FB.Firestore.Collection("carts").Doc(uid)

	// ArrayUnion appends the item to the "items" array atomically, preventing race conditions
	// if the user clicks rapidly or uses multiple tabs.
	_, err := cartRef.Update(ctx, []firestore.Update{
		{Path: "items", Value: firestore.ArrayUnion(item)},
		{Path: "updatedAt", Value: time.Now()},
	})

	if err != nil {
		// If Update fails because the document doesn't exist yet, we Set it to create it.
		_, err = cartRef.Set(ctx, map[string]interface{}{
			"items":     []models.CartItem{item},
			"updatedAt": time.Now(),
		})
		if err != nil {
			return fmt.Errorf("failed to create cart: %v", err)
		}
	}
	return nil
}

// GetCart retrieves the user's current cart and calculates the live total.
func GetCart(uid string) (*models.Cart, error) {
	ctx := context.Background()
	doc, err := FB.Firestore.Collection("carts").Doc(uid).Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("cart not found: %v", err)
	}

	var cart models.Cart
	if err := doc.DataTo(&cart); err != nil {
		return nil, fmt.Errorf("failed to parse cart data: %v", err)
	}

	// Dynamically calculate the total on the server side to ensure pricing accuracy
	cart.Subtotal = 0
	for _, item := range cart.Items {
		cart.Subtotal += item.ItemTotal
	}
	// Add tax/shipping logic here later if needed
	cart.Total = cart.Subtotal

	return &cart, nil
}

// RemoveCartItem filters out a specific item ID and overwrites the cart array.
func RemoveCartItem(uid string, itemID string) error {
	ctx := context.Background()
	cart, err := GetCart(uid)
	if err != nil {
		return err
	}

	var updatedItems []models.CartItem
	for _, item := range cart.Items {
		if item.ID != itemID {
			updatedItems = append(updatedItems, item)
		}
	}

	_, err = FB.Firestore.Collection("carts").Doc(uid).Update(ctx, []firestore.Update{
		{Path: "items", Value: updatedItems},
		{Path: "updatedAt", Value: time.Now()},
	})

	return err
}
