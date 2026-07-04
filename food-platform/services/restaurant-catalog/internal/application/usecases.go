// Package application contains the use cases for the Restaurant Catalog Service.
package application

import (
	"context"
	"fmt"
	"math"

	"github.com/food-platform/restaurant-catalog/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

// RestaurantRepository is the interface for persisting restaurants.
type RestaurantRepository interface {
	Create(ctx context.Context, r *domain.Restaurant) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Restaurant, error)
	Update(ctx context.Context, r *domain.Restaurant) error
	FindNearby(ctx context.Context, lat, lng float64, radiusKm float64, limit, offset int) ([]*domain.Restaurant, error)
	Search(ctx context.Context, query string, lat, lng float64, limit, offset int) ([]*domain.Restaurant, error)
	FindByCuisine(ctx context.Context, cuisine domain.CuisineType, lat, lng float64, limit, offset int) ([]*domain.Restaurant, error)
	UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error
}

// ============ DTOs ============

// CreateRestaurantCommand creates a new restaurant.
type CreateRestaurantCommand struct {
	Name         string
	Latitude     float64
	Longitude    float64
	Address      string
	City         string
	CuisineTypes []domain.CuisineType
}

// UpdateRestaurantCommand updates an existing restaurant.
type UpdateRestaurantCommand struct {
	ID           uuid.UUID
	Name         *string
	Latitude     *float64
	Longitude    *float64
	Address      *string
	LogoURL      *string
	CoverURL     *string
	OpensAt      *string
	ClosesAt     *string
	ETAMin       *int
	ETAMax       *int
	DeliveryFee  *float64
	IsOpen       *bool
}

// NearbyQuery finds restaurants near a location.
type NearbyQuery struct {
	Latitude  float64
	Longitude float64
	RadiusKm  float64
	Limit     int
	Offset    int
}

// RestaurantDTO is the output representation of a restaurant.
type RestaurantDTO struct {
	ID            uuid.UUID                  `json:"id"`
	Name          string                     `json:"name"`
	Slug          string                     `json:"slug"`
	CuisineTypes  []string                   `json:"cuisine_types"`
	Rating        float64                    `json:"rating"`
	RatingCount   int                        `json:"rating_count"`
	LogoURL       string                     `json:"logo_url,omitempty"`
	CoverURL      string                     `json:"cover_url,omitempty"`
	Latitude      float64                    `json:"latitude"`
	Longitude     float64                    `json:"longitude"`
	Address       string                     `json:"address"`
	City          string                     `json:"city"`
	IsOpen        bool                       `json:"is_open"`
	Status        string                     `json:"status"`
	ETAMinMinutes int                        `json:"eta_minutes_min"`
	ETAMaxMinutes int                        `json:"eta_minutes_max"`
	DeliveryFee   float64                    `json:"delivery_fee"`
	PriceRange    int                        `json:"price_range"`
	DistanceKm    float64                    `json:"distance_km,omitempty"`
}

// ============ Use Cases ============

// CreateRestaurantUseCase creates a new restaurant.
type CreateRestaurantUseCase struct {
	repo RestaurantRepository
}

func NewCreateRestaurantUseCase(repo RestaurantRepository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{repo: repo}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, cmd CreateRestaurantCommand) (*RestaurantDTO, error) {
	restaurant, err := domain.NewRestaurant(
		cmd.Name,
		cmd.Latitude,
		cmd.Longitude,
		cmd.Address,
		cmd.City,
		cmd.CuisineTypes,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, restaurant); err != nil {
		return nil, fmt.Errorf("create restaurant: %w", err)
	}

	return toDTO(restaurant, 0), nil
}

// GetRestaurantUseCase retrieves a restaurant by ID.
type GetRestaurantUseCase struct {
	repo RestaurantRepository
}

func NewGetRestaurantUseCase(repo RestaurantRepository) *GetRestaurantUseCase {
	return &GetRestaurantUseCase{repo: repo}
}

func (uc *GetRestaurantUseCase) Execute(ctx context.Context, id uuid.UUID) (*RestaurantDTO, error) {
	r, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(r, 0), nil
}

// FindNearbyUseCase finds restaurants near a location.
type FindNearbyUseCase struct {
	repo RestaurantRepository
}

func NewFindNearbyUseCase(repo RestaurantRepository) *FindNearbyUseCase {
	return &FindNearbyUseCase{repo: repo}
}

func (uc *FindNearbyUseCase) Execute(ctx context.Context, q NearbyQuery) ([]*RestaurantDTO, error) {
	if q.Limit == 0 {
		q.Limit = 20
	}
	if q.RadiusKm == 0 {
		q.RadiusKm = 3
	}

	restaurants, err := uc.repo.FindNearby(ctx, q.Latitude, q.Longitude, q.RadiusKm, q.Limit, q.Offset)
	if err != nil {
		return nil, err
	}

	dtos := make([]*RestaurantDTO, 0, len(restaurants))
	for _, r := range restaurants {
		dist := r.DistanceTo(q.Latitude, q.Longitude)
		dtos = append(dtos, toDTO(r, dist))
	}
	return dtos, nil
}

// SearchRestaurantsUseCase searches restaurants by name.
type SearchRestaurantsUseCase struct {
	repo RestaurantRepository
}

func NewSearchRestaurantsUseCase(repo RestaurantRepository) *SearchRestaurantsUseCase {
	return &SearchRestaurantsUseCase{repo: repo}
}

func (uc *SearchRestaurantsUseCase) Execute(ctx context.Context, query string, lat, lng float64, limit, offset int) ([]*RestaurantDTO, error) {
	if limit == 0 {
		limit = 20
	}

	restaurants, err := uc.repo.Search(ctx, query, lat, lng, limit, offset)
	if err != nil {
		return nil, err
	}

	dtos := make([]*RestaurantDTO, 0, len(restaurants))
	for _, r := range restaurants {
		dist := 0.0
		if lat != 0 && lng != 0 {
			dist = r.DistanceTo(lat, lng)
		}
		dtos = append(dtos, toDTO(r, dist))
	}
	return dtos, nil
}

// UpdateRestaurantUseCase updates restaurant details.
type UpdateRestaurantUseCase struct {
	repo RestaurantRepository
}

func NewUpdateRestaurantUseCase(repo RestaurantRepository) *UpdateRestaurantUseCase {
	return &UpdateRestaurantUseCase{repo: repo}
}

func (uc *UpdateRestaurantUseCase) Execute(ctx context.Context, cmd UpdateRestaurantCommand) (*RestaurantDTO, error) {
	r, err := uc.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	if cmd.Name != nil {
		if err := r.SetName(*cmd.Name); err != nil {
			return nil, err
		}
	}
	if cmd.Latitude != nil && cmd.Longitude != nil {
		if err := r.SetLocation(*cmd.Latitude, *cmd.Longitude, r.Address()); err != nil {
			return nil, err
		}
	}
	if cmd.Address != nil {
		if err := r.SetLocation(r.Latitude(), r.Longitude(), *cmd.Address); err != nil {
			return nil, err
		}
	}
	if cmd.LogoURL != nil {
		r.SetLogoURL(*cmd.LogoURL)
	}
	if cmd.CoverURL != nil {
		r.SetCoverURL(*cmd.CoverURL)
	}
	if cmd.OpensAt != nil && cmd.ClosesAt != nil {
		r.SetHours(*cmd.OpensAt, *cmd.ClosesAt)
	}
	if cmd.ETAMin != nil && cmd.ETAMax != nil {
		r.SetETA(*cmd.ETAMin, *cmd.ETAMax)
	}
	if cmd.DeliveryFee != nil {
		r.SetDeliveryFee(*cmd.DeliveryFee)
	}
	if cmd.IsOpen != nil {
		r.SetOpen(*cmd.IsOpen)
	}

	if err := uc.repo.Update(ctx, r); err != nil {
		return nil, fmt.Errorf("update restaurant: %w", err)
	}

	return toDTO(r, 0), nil
}

// ============ Helpers ============

func toDTO(r *domain.Restaurant, distanceKm float64) *RestaurantDTO {
	cuisines := make([]string, len(r.CuisineTypes()))
	for i, c := range r.CuisineTypes() {
		cuisines[i] = string(c)
	}

	dto := &RestaurantDTO{
		ID:            r.ID(),
		Name:          r.Name(),
		Slug:          r.Slug(),
		CuisineTypes:  cuisines,
		Rating:        math.Round(r.Rating()*10) / 10,
		RatingCount:   r.RatingCount(),
		LogoURL:       r.LogoURL(),
		CoverURL:      r.CoverURL(),
		Latitude:      r.Latitude(),
		Longitude:     r.Longitude(),
		Address:       r.Address(),
		City:          r.City(),
		IsOpen:        r.IsOpen(),
		Status:        string(r.Status()),
		ETAMinMinutes: r.ETAMinMinutes(),
		ETAMaxMinutes: r.ETAMaxMinutes(),
		DeliveryFee:   r.DeliveryFee(),
		PriceRange:    r.PriceRange(),
	}

	if distanceKm > 0 {
		dto.DistanceKm = math.Round(distanceKm*10) / 10
	}

	return dto
}
