package middleware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the HTTP status code
// written by the handler. If WriteHeader is never called the status defaults
// to 200, matching the net/http implicit behaviour.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// WriteHeader captures the status code and delegates to the underlying writer.
// Subsequent calls are no-ops to match the stdlib contract.
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.wroteHeader {
		rw.status = code
		rw.wroteHeader = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// Status returns the captured HTTP status code.
// Returns http.StatusOK (200) when WriteHeader was never explicitly called.
func (rw *responseWriter) Status() int {
	if !rw.wroteHeader {
		return http.StatusOK
	}
	return rw.status
}

// DefaultLogger returns a *slog.Logger that emits JSON records to stdout.
func DefaultLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// Logger returns middleware that logs each completed request as a structured
// JSON record. Fields: time (auto), level, request_id, method, path,
// status (int), duration_ms (float64).
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := newResponseWriter(w)

			next.ServeHTTP(rw, r)

			durationMS := float64(time.Since(start).Microseconds()) / 1000.0

			logger.InfoContext(
				r.Context(),
				"request",
				slog.String("request_id", RequestIDFromContext(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", rw.Status()),
				slog.Float64("duration_ms", durationMS),
			)
		})
	}
}
