package http

import (
	"net/http"

	"github.com/food-platform/geo/internal/application"
	"github.com/food-platform/geo/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	updateLoc *application.UpdateLocationUseCase,
	getLoc *application.GetLocationUseCase,
	findNearby *application.FindNearbyUseCase,
	calcETA *application.CalculateETAUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/geo", func(r chi.Router) {
		// Authenticated routes (drivers update their location)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)
			r.Post("/location", handlers.NewUpdateLocationHandler(updateLoc).ServeHTTP)
		})

		// Internal routes (other services query driver locations)
		r.Get("/location", handlers.NewGetLocationHandler(getLoc).ServeHTTP)
		r.Get("/nearby", handlers.NewFindNearbyHandler(findNearby).ServeHTTP)
		r.Get("/eta", handlers.NewCalculateETAHandler(calcETA).ServeHTTP)
	})

	return r
}
