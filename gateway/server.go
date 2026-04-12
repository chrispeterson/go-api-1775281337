package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Server wraps http.Server with configuration and lifecycle management
type Server struct {
	httpServer *http.Server
	port       string
	logger     *slog.Logger
}

// NewServer creates a new Server instance with the given logger
func NewServer(logger *slog.Logger) *Server {
	if logger == nil {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}

	return &Server{
		logger: logger,
	}
}

// Initialize sets up the HTTP server with routes and configuration
func (s *Server) Initialize(mux *http.ServeMux) error {
	// Read port from environment variable, default to 8080
	s.port = os.Getenv("PORT")
	if s.port == "" {
		s.port = "8080"
	}

	// Create HTTP server with sensible defaults
	s.httpServer = &http.Server{
		Addr:         ":" + s.port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.logger.Info("server initialized", "port", s.port)
	return nil
}

// ListenAndServe starts the HTTP server and blocks until it exits
func (s *Server) ListenAndServe() error {
	if s.httpServer == nil {
		return s.Initialize(http.NewServeMux())
	}

	s.logger.Info("starting server", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server with a timeout
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("shutting down server")
	return s.httpServer.Shutdown(ctx)
}

// GetPort returns the port the server is configured to listen on
func (s *Server) GetPort() string {
	return s.port
}
