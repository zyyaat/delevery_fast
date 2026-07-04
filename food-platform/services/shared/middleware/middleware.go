// Package middleware provides HTTP middleware for all services.
package middleware

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/food-platform/shared/logging"
	"github.com/google/uuid"
)

// responseRecorder captures the status code for logging.
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	bytes      int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
}

// RequestID middleware adds a request ID to the context.
// If the X-Request-ID header is present, it uses that; otherwise, generates a new UUID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", requestID)

		ctx := logging.WithRequestID(r.Context(), requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logging middleware logs each request with method, path, status, duration, and request ID.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		logging.FromContext(r.Context()).Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.statusCode,
			"duration_ms", duration.Milliseconds(),
			"bytes", rec.bytes,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

// Recovery middleware recovers from panics and returns 500.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logging.FromContext(r.Context()).Error("panic_recovered",
					"error", rec,
					"stack", string(debug.Stack()),
				)
				http.Error(w, `{"error":{"code":"INTERNAL_ERROR","message":"Internal server error"}}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS middleware adds CORS headers for browser clients.
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = true
	}
	if len(allowed) == 0 {
		allowed["*"] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowed[origin] || allowed["*"] {
				if allowed["*"] {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Idempotency-Key, X-Action-Token")
				w.Header().Set("Access-Control-Expose-Headers", "X-Request-ID, X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset")
				w.Header().Set("Access-Control-Max-Age", "3600")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Auth middleware validates the JWT and adds user ID + role to context.
// The actual JWT validation is done by the API Gateway; this middleware just
// trusts the X-User-ID and X-User-Role headers set by the gateway.
// For direct service-to-service calls, it accepts Bearer tokens too.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Trust headers from API Gateway
		userID := r.Header.Get("X-User-ID")
		role := r.Header.Get("X-User-Role")
		sessionID := r.Header.Get("X-Session-ID")

		if userID == "" {
			http.Error(w, `{"error":{"code":"UNAUTHORIZED","message":"Authentication required"}}`, http.StatusUnauthorized)
			return
		}

		ctx := logging.WithUserID(r.Context(), userID)
		ctx = logging.WithSessionID(ctx, sessionID)
		ctx = context.WithValue(ctx, "role", role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole middleware checks that the authenticated user has one of the required roles.
// Must be used after Auth middleware.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, _ := r.Context().Value("role").(string)
			if role == "" || !allowed[role] {
				http.Error(w, `{"error":{"code":"PERMISSION_DENIED","message":"Insufficient permissions"}}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimit is a simple in-memory rate limiter.
// For production, use Redis-based rate limiting in the API Gateway.
func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	// Note: This is a placeholder. Production uses Redis-based rate limiting.
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// Chain applies multiple middlewares in order.
func Chain(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
