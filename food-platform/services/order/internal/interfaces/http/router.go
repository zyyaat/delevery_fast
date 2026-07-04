package http

import (
	"net/http"

	"github.com/food-platform/services/order/internal/application"
	"github.com/food-platform/services/order/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	create *application.CreateOrderUseCase,
	get *application.GetOrderUseCase,
	getActive *application.GetActiveOrdersUseCase,
	getHistory *application.GetOrderHistoryUseCase,
	cancel *application.CancelOrderUseCase,
	updateStatus *application.UpdateOrderStatusUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/orders", func(r chi.Router) {
		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)

			r.Post("/", handlers.NewCreateOrderHandler(create).ServeHTTP)
			r.Get("/", handlers.NewGetOrderHistoryHandler(getHistory).ServeHTTP)
			r.Get("/active", handlers.NewGetActiveOrdersHandler(getActive).ServeHTTP)
			r.Get("/{id}", handlers.NewGetOrderHandler(get).ServeHTTP)
			r.Post("/{id}/cancel", handlers.NewCancelOrderHandler(cancel).ServeHTTP)
			r.Patch("/{id}/status", handlers.NewUpdateStatusHandler(updateStatus).ServeHTTP)
		})
	})

	return r
}
