package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chrispeterson/go-api-1775281337/gateway/handler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create server instance
	server := NewServer(logger)

	// Create HTTP mux and register handlers
	mux := http.NewServeMux()

	helloHandler := handler.NewHelloHandler()
	mux.Handle("GET /hello", helloHandler)

	// Initialize server with the mux
	if err := server.Initialize(mux); err != nil {
		logger.Error("failed to initialize server", "error", err)
		os.Exit(1)
	}

	logger.Info("starting gateway server", "port", server.GetPort())

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Listen()
	}()

	// Wait for either signal or server error
	select {
	case sig := <-sigChan:
		logger.Info("received signal", "signal", sig.String())
		// Graceful shutdown with 10 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", "error", err)
			os.Exit(1)
		}
		logger.Info("server gracefully shut down")
	case err := <-errChan:
		if err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}
}
