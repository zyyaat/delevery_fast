package http

import (
	"net/http"

	"github.com/food-platform/driver-management/internal/application"
	"github.com/food-platform/driver-management/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	register *application.RegisterDriverUseCase,
	get *application.GetDriverUseCase,
	updateStatus *application.UpdateStatusUseCase,
	updateLocation *application.UpdateLocationUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/drivers", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)

			r.Post("/", handlers.NewRegisterDriverHandler(register).ServeHTTP)
			r.Get("/me", handlers.NewGetDriverHandler(get).ServeHTTP)
			r.Put("/me/status", handlers.NewUpdateStatusHandler(updateStatus).ServeHTTP)
			r.Post("/me/location", handlers.NewUpdateLocationHandler(updateLocation).ServeHTTP)
		})
	})

	return r
}
