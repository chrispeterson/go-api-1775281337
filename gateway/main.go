package main

import (
	"context"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	helloHandler := handler.NewHelloHandler()
	mux.Handle("/hello", helloHandler)

	srv := NewServer(logger)
	err := srv.Initialize(mux)
	if err != nil {
		logger.Error("failed to initialize server", "error", err)
		os.Exit(1)
	}

	// Start server in background.
	go func() {
		slog.Info("starting server", "addr", ":"+port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "err", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
