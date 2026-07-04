package http

import (
	"net/http"

	"github.com/food-platform/services/payment/internal/application"
	"github.com/food-platform/services/payment/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	charge *application.ChargePaymentUseCase,
	get *application.GetPaymentUseCase,
	getByOrder *application.GetPaymentByOrderUseCase,
	refund *application.RefundPaymentUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/payments", func(r chi.Router) {
		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth)

			r.Post("/", handlers.NewChargePaymentHandler(charge).ServeHTTP)
			r.Get("/{id}", handlers.NewGetPaymentHandler(get).ServeHTTP)
			r.Get("/order/{orderId}", handlers.NewGetPaymentByOrderHandler(getByOrder).ServeHTTP)

			// Support/admin routes (require elevated roles)
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("support_l1", "support_l2", "ops_manager", "finance", "super_admin"))
				r.Post("/{id}/refund", handlers.NewRefundHandler(refund).ServeHTTP)
			})
		})
	})

	return r
}
