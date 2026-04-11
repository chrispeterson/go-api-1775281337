package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the response body for the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// Health handles GET /health requests and returns the service health status.
func Health(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set Content-Type header before writing response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Serialize response using encoding/json.Marshal
	resp := HealthResponse{Status: "ok"}
	body, err := json.Marshal(resp)
	if err != nil {
		// If marshaling fails, write error response
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(body)
}
