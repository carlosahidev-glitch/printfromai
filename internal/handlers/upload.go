package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/printfromai/core/internal/middleware"
	"github.com/printfromai/core/internal/services"
)

// UploadURLResponse defines the JSON structure sent back to the frontend.
type UploadURLResponse struct {
	URL        string `json:"url"`
	ObjectName string `json:"objectName"`
	Error      string `json:"error,omitempty"`
}

// GetUploadURL generates a secure, presigned Google Cloud Storage URL.
// The frontend uses this to upload massive print files directly to the cloud,
// bypassing the Go server entirely to prevent memory spikes.
func GetUploadURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Get the securely verified Firebase UID from the context middleware
	uid := middleware.GetUID(r)
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(UploadURLResponse{Error: "Unauthorized"})
		return
	}

	// 2. Extract the requested filename from the query string
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(UploadURLResponse{Error: "Filename is required"})
		return
	}

	// 3. Generate the Presigned URL via our Firebase/GCS service
	url, objectName, err := services.FB.GeneratePresignedUploadURL(uid, filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(UploadURLResponse{Error: "Failed to generate upload URL"})
		return
	}

	// 4. Return the secure URL and the storage path back to the client
	json.NewEncoder(w).Encode(UploadURLResponse{
		URL:        url,
		ObjectName: objectName,
	})
}
