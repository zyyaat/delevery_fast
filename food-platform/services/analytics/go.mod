module github.com/food-platform/analytics

go 1.22

require (
	github.com/food-platform/shared v0.0.0
	github.com/go-chi/chi/v5 v5.0.12
)

require github.com/google/uuid v1.6.0 // indirect

replace github.com/food-platform/shared => ../shared
