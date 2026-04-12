package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Server wraps http.Server with additional configuration and lifecycle methods.
type Server struct {
	logger     *slog.Logger
	httpServer *http.Server
	port       string
}

// NewServer creates a new Server instance with an optional logger.
func NewServer(logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return &Server{
		logger: logger,
	}
}

// Initialize sets up the HTTP server with the given mux and port configuration.
// It reads the PORT environment variable (default 8080) and configures timeouts.
func (s *Server) Initialize(mux *http.ServeMux) error {
	// Read port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s.port = port

	// Create HTTP server with timeouts
	s.httpServer = &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return nil
}

// Listen starts the HTTP server and blocks until an error occurs.
// It returns any error from ListenAndServe.
func (s *Server) Listen() error {
	if s.httpServer == nil {
		return fmt.Errorf("server not initialized, call Initialize() first")
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server with the given context.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}
	return s.httpServer.Shutdown(ctx)
}

// GetPort returns the port the server is listening on.
func (s *Server) GetPort() string {
	return s.port
}
