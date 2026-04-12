package main

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	s := NewServer(nil)
	require.NotNil(t, s)
	require.NotNil(t, s.logger)
}

func TestNewServerWithCustomLogger(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	s := NewServer(logger)
	require.NotNil(t, s)
	assert.Equal(t, logger, s.logger)
}

func TestInitializeWithDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")

	s := NewServer(nil)
	mux := http.NewServeMux()
	err := s.Initialize(mux)
	require.NoError(t, err)

	assert.Equal(t, "8080", s.GetPort())
	assert.NotNil(t, s.httpServer)
}

func TestInitializeWithEnvironmentPort(t *testing.T) {
	t.Setenv("PORT", "9000")

	s := NewServer(nil)
	mux := http.NewServeMux()
	err := s.Initialize(mux)
	require.NoError(t, err)

	assert.Equal(t, "9000", s.GetPort())
}

func TestInitializeWithCustomPort(t *testing.T) {
	t.Setenv("PORT", "3000")

	s := NewServer(nil)
	mux := http.NewServeMux()
	err := s.Initialize(mux)
	require.NoError(t, err)

	assert.Equal(t, "3000", s.GetPort())
}

func TestUnknownRoutesReturn404(t *testing.T) {
	s := NewServer(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := s.Initialize(mux)
	require.NoError(t, err)

	// Create test server using the mux
	server := httptest.NewServer(s.httpServer.Handler)
	defer server.Close()

	// Test unmapped route returns 404
	resp, err := http.Get(server.URL + "/nonexistent")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestKnownRoutesReturn200(t *testing.T) {
	s := NewServer(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := s.Initialize(mux)
	require.NoError(t, err)

	// Create test server using the mux
	server := httptest.NewServer(s.httpServer.Handler)
	defer server.Close()

	// Test known endpoint returns 200
	resp, err := http.Get(server.URL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGracefulShutdown(t *testing.T) {
	s := NewServer(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := s.Initialize(mux)
	require.NoError(t, err)

	// Use httptest to avoid actual port binding
	server := httptest.NewServer(s.httpServer.Handler)
	defer server.Close()

	// Test that shutdown succeeds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Shutdown(ctx)
	require.NoError(t, err)
}

func TestShutdownWithoutInitialize(t *testing.T) {
	s := NewServer(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Should not error if server was never initialized
	err := s.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestServerHandlerConfigured(t *testing.T) {
	s := NewServer(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	err := s.Initialize(mux)
	require.NoError(t, err)

	assert.NotNil(t, s.httpServer.Handler)
	assert.Equal(t, mux, s.httpServer.Handler)
}

func TestServerTimeoutConfiguration(t *testing.T) {
	s := NewServer(nil)
	mux := http.NewServeMux()

	err := s.Initialize(mux)
	require.NoError(t, err)

	assert.Equal(t, 15*time.Second, s.httpServer.ReadTimeout)
	assert.Equal(t, 15*time.Second, s.httpServer.WriteTimeout)
	assert.Equal(t, 60*time.Second, s.httpServer.IdleTimeout)
}
