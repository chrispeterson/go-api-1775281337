package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// RequestIDHeader is the HTTP header used to propagate request IDs.
const RequestIDHeader = "X-Request-ID"

// contextKey is an unexported type for context keys in this package
// to avoid collisions with keys from other packages.
type contextKey string

// RequestIDKey is the context key under which the request ID is stored.
const RequestIDKey contextKey = "request_id"

// RequestID is middleware that ensures every request carries a unique ID.
// If the incoming request already has an X-Request-ID header that value is
// reused; otherwise a new UUID v4 is generated. The ID is:
//   - stored in the request context under RequestIDKey
//   - echoed back to the client via the X-Request-ID response header
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}

		// Propagate into context so downstream handlers can read it.
		ctx := context.WithValue(r.Context(), RequestIDKey, id)

		// Set on response before calling next so it is always present,
		// even if the handler panics or writes headers early.
		w.Header().Set(RequestIDHeader, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext retrieves the request ID stored in ctx.
// Returns an empty string if the value is absent.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}
