package handler

import (
	"encoding/json"
	"net/http"
)

// HealthResponse represents the response from the health endpoint
type HealthResponse struct {
	Status string `json:"status"`
}

// Health handles GET /health requests
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := HealthResponse{
		Status: "ok",
	}

	_ = json.NewEncoder(w).Encode(response)
}
