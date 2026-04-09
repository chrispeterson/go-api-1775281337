package router

import (
	"net/http"

	"github.com/chrispeterson/go-api-1775281337/gateway/handler"
)

// Register sets up all routes for the gateway service
func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handler.Health)
}
