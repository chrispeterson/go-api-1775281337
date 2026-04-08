package middleware

import "net/http"

// Chain combines multiple middleware functions into a single middleware.
// Middlewares are applied in order: first middleware listed is outermost.
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Apply in reverse so the first middleware listed wraps outermost
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
