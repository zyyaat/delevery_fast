package http

import (
	"net/http"

	"github.com/food-platform/analytics/internal/application"
	"github.com/food-platform/analytics/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	dashboard *application.GetDashboardStatsUseCase,
	zones *application.GetZoneMetricsUseCase,
	incidents *application.GetIncidentsUseCase,
	forecast *application.GetForecastUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/analytics", func(r chi.Router) {
		r.Get("/dashboard", handlers.NewDashboardHandler(dashboard).ServeHTTP)
		r.Get("/zones", handlers.NewZonesHandler(zones).ServeHTTP)
		r.Get("/incidents", handlers.NewIncidentsHandler(incidents).ServeHTTP)
		r.Get("/forecast", handlers.NewForecastHandler(forecast).ServeHTTP)
	})

	return r
}
