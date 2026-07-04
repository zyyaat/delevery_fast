package http

import (
	"net/http"

	"github.com/food-platform/menu/internal/application"
	"github.com/food-platform/menu/internal/interfaces/http/handlers"
	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(
	createItem *application.CreateMenuItemUseCase,
	getMenu *application.GetMenuUseCase,
	toggle *application.ToggleAvailabilityUseCase,
	createCat *application.CreateCategoryUseCase,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", handlers.NewHealthHandler().ServeHTTP)

	r.Route("/api/v1/restaurants/{restaurantId}", func(r chi.Router) {
		r.Get("/menu", handlers.NewGetMenuHandler(getMenu).ServeHTTP)
		r.Post("/categories", handlers.NewCreateCategoryHandler(createCat).ServeHTTP)
		r.Post("/items", handlers.NewCreateMenuItemHandler(createItem).ServeHTTP)
		r.Patch("/items/{itemId}/availability", handlers.NewToggleAvailabilityHandler(toggle).ServeHTTP)
	})

	return r
}
