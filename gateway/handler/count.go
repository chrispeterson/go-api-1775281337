package handler

import (
	"encoding/json"
	"net/http"

	"gateway/internal"
)

// CountHandler handles count-related HTTP requests.
type CountHandler struct {
	counter *internal.Counter
}

// NewCountHandler creates a new CountHandler.
func NewCountHandler(counter *internal.Counter) *CountHandler {
	return &CountHandler{
		counter: counter,
	}
}

// CountResponse is the JSON response for GET /count.
type CountResponse struct {
	Count int64 `json:"count"`
}

// ServeHTTP handles GET /count requests.
func (h *CountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	countValue := h.counter.Get()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := CountResponse{Count: countValue}
	json.NewEncoder(w).Encode(response)
}
