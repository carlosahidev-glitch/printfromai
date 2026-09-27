package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/printfromai/core/internal/app"
)

// main is the entry point for PrintFromAI.com
func main() {
	// 1. Configuration & Environment Setup
	// In Cloud Run, PORT is automatically provided. For local dev, we default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// TODO: Initialize Firebase App, Stripe Client, and AWS SES here in future steps.

	// 2. Initialize the main application router
	// We delegate all routing to our internal/app package to keep main.go clean.
	mux := app.NewRouter()

	// 3. Configure the HTTP Server
	// Adding timeouts is critical for a public-facing server to prevent slow-loris attacks
	// and ensure connections don't hang indefinitely during large file uploads.
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  120 * time.Second, // Increased to allow large file uploads
		WriteTimeout: 120 * time.Second, // Increased for large uploads
		IdleTimeout:  120 * time.Second,
	}

	// 4. Start the Server
	fmt.Printf("🚀 PrintFromAI Server starting on http://localhost:%s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Fatal server error: %v", err)
	}
}
