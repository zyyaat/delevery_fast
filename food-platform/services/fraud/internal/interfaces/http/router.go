package http

import (
	"net/http"

	"github.com/food-platform/fraud/internal/application"
	"github.com/food-platform/fraud/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(scoreOrder *application.ScoreOrderUseCase, getTrust *application.GetTrustScoreUseCase) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/fraud", func(r chi.Router) {
		// Internal routes (called by Order Service)
		r.Post("/score", handlers.NewScoreOrderHandler(scoreOrder).ServeHTTP)
		r.Get("/trust/{customerId}", handlers.NewGetTrustScoreHandler(getTrust).ServeHTTP)
	})

	return r
}
