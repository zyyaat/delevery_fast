module github.com/food-platform/fraud

go 1.22

require (
	github.com/food-platform/shared v0.0.0
	github.com/go-chi/chi/v5 v5.0.12
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

replace github.com/food-platform/shared => ../shared
