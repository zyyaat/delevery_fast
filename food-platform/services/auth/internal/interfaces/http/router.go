// Package http sets up the HTTP server with routes and middleware.
package http

import (
	"net/http"

	"github.com/food-platform/services/auth/internal/application"
	"github.com/food-platform/services/auth/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

// SetupRouter creates and configures the HTTP router with all routes.
func SetupRouter(
	sendOTP *application.SendOTPUseCase,
	verifyOTP *application.VerifyOTPUseCase,
	refreshToken *application.RefreshTokenUseCase,
	logout *application.LogoutUseCase,
) http.Handler {
	r := chi.NewRouter()

	// Apply standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	// Health endpoints (no auth)
	r.Get("/health", handlers.NewHealthHandler("1.0.0").ServeHTTP)
	r.Get("/ready", handlers.NewReadyHandler().ServeHTTP)

	// Auth routes (no auth required for these)
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/otp/send", handlers.NewSendOTPHandler(sendOTP).ServeHTTP)
		r.Post("/otp/verify", handlers.NewVerifyOTPHandler(verifyOTP).ServeHTTP)
		r.Post("/refresh", handlers.NewRefreshTokenHandler(refreshToken).ServeHTTP)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Post("/logout", handlers.NewLogoutHandler(logout).ServeHTTP)
		})
	})

	// 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":{"code":"NOT_FOUND","message":"Endpoint not found"}}`))
	})

	// Method not allowed
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":{"code":"METHOD_NOT_ALLOWED","message":"Method not allowed"}}`))
	})

	return r
}
