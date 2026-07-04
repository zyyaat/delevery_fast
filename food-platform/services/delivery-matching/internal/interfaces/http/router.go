package http

import (
	"net/http"

	"github.com/food-platform/delivery-matching/internal/application"
	"github.com/food-platform/delivery-matching/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	match *application.MatchOrderUseCase,
	accept *application.AcceptOrderUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/matching", func(r chi.Router) {
		// Internal route (called by Order Service via Kafka or gRPC)
		r.Post("/match", handlers.NewMatchOrderHandler(match).ServeHTTP)

		// Authenticated route (drivers accept orders)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Post("/accept", handlers.NewAcceptOrderHandler(accept).ServeHTTP)
		})
	})

	return r
}
