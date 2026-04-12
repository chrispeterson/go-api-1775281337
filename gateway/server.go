package main

import (
	"context"
	"net/http"
	"time"
)

// Server wraps http.Server with additional configuration and lifecycle methods.
type Server struct {
	http *http.Server
}

// NewServer creates a new Server instance with the given port.
// The port should be formatted as ":8080" or similar.
func NewServer(port string) *Server {
	mux := http.NewServeMux()

	// Register default handler for root path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// ServeMux will handle 404 automatically for unregistered routes
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	httpServer := &http.Server{
		Addr:         port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		http: httpServer,
	}
}

// Listen starts the HTTP server and blocks until an error occurs.
// It returns any error from ListenAndServe.
func (s *Server) Listen() error {
	return s.http.ListenAndServe()
}

// Shutdown gracefully shuts down the server with the given context.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

// GetPort returns the port the server is listening on.
func (s *Server) GetPort() string {
	return s.http.Addr
}
