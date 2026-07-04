package http

import (
	"net/http"

	"github.com/food-platform/restaurant-catalog/internal/application"
	"github.com/food-platform/restaurant-catalog/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	create *application.CreateRestaurantUseCase,
	get *application.GetRestaurantUseCase,
	nearby *application.FindNearbyUseCase,
	search *application.SearchRestaurantsUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/restaurants", func(r chi.Router) {
		r.Post("/", handlers.NewCreateRestaurantHandler(create).ServeHTTP)
		r.Get("/nearby", handlers.NewFindNearbyHandler(nearby).ServeHTTP)
		r.Get("/search", handlers.NewSearchHandler(search).ServeHTTP)
		r.Get("/{id}", handlers.NewGetRestaurantHandler(get).ServeHTTP)
	})

	return r
}
