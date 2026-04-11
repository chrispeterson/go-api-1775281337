package handler

import (
	"encoding/json"
	"net/http"
)

// HelloHandler handles GET /hello requests
type HelloHandler struct{}

// NewHelloHandler returns a new HelloHandler
func NewHelloHandler() *HelloHandler {
	return &HelloHandler{}
}

// greetingResponse represents the JSON response body
type greetingResponse struct {
	Message string `json:"message"`
}

// errorResponse represents an error JSON response body
type errorResponse struct {
	Error string `json:"error"`
}

// ServeHTTP implements http.Handler for the hello endpoint
func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "method not allowed"})
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "name parameter is required"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(greetingResponse{Message: "Hello, " + name + "!"})
}
