// Package postgres implements the Restaurant Catalog repository using PostgreSQL + PostGIS.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/food-platform/restaurant-catalog/internal/domain"
	"github.com/google/uuid"
)

// RestaurantRepository implements application.RestaurantRepository.
type RestaurantRepository struct {
	db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
	return &RestaurantRepository{db: db}
}

// Create inserts a new restaurant.
func (r *RestaurantRepository) Create(ctx context.Context, rest *domain.Restaurant) error {
	query := `
		INSERT INTO restaurants (
			id, name, slug, cuisine_types, rating, rating_count,
			logo_url, cover_url, location, address, city,
			is_open, status, eta_min_minutes, eta_max_minutes,
			delivery_fee, price_range, commission_rate,
			opens_at, closes_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, ST_SetSRID(ST_MakePoint($9, $10), 4326), $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19,
			$20, $21, $22, $23
		)
	`

	cuisines := make([]string, len(rest.CuisineTypes()))
	for i, c := range rest.CuisineTypes() {
		cuisines[i] = string(c)
	}

	_, err := r.db.ExecContext(ctx, query,
		rest.ID(),
		rest.Name(),
		rest.Slug(),
		arrayToPostgres(cuisines),
		rest.Rating(),
		rest.RatingCount(),
		nullableString(rest.LogoURL()),
		nullableString(rest.CoverURL()),
		rest.Longitude(), // MakePoint(lng, lat) — note the order!
		rest.Latitude(),
		rest.Address(),
		rest.City(),
		rest.IsOpen(),
		string(rest.Status()),
		rest.ETAMinMinutes(),
		rest.ETAMaxMinutes(),
		rest.DeliveryFee(),
		rest.PriceRange(),
		rest.CommissionRate(),
		rest.OpensAt(),
		rest.ClosesAt(),
		rest.CreatedAt(),
		rest.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("restaurant_repo.Create: %w", err)
	}
	return nil
}

// FindByID retrieves a restaurant by ID.
func (r *RestaurantRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	query := `
		SELECT id, name, slug, cuisine_types, rating, rating_count,
		       COALESCE(logo_url, ''), COALESCE(cover_url, ''),
		       ST_Y(location) as lat, ST_X(location) as lng,
		       address, city, is_open, status,
		       eta_min_minutes, eta_max_minutes, delivery_fee,
		       price_range, commission_rate, opens_at, closes_at,
		       created_at, updated_at
		FROM restaurants
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanRestaurant(row)
}

// FindBySlug retrieves a restaurant by slug.
func (r *RestaurantRepository) FindBySlug(ctx context.Context, slug string) (*domain.Restaurant, error) {
	query := `
		SELECT id, name, slug, cuisine_types, rating, rating_count,
		       COALESCE(logo_url, ''), COALESCE(cover_url, ''),
		       ST_Y(location) as lat, ST_X(location) as lng,
		       address, city, is_open, status,
		       eta_min_minutes, eta_max_minutes, delivery_fee,
		       price_range, commission_rate, opens_at, closes_at,
		       created_at, updated_at
		FROM restaurants
		WHERE slug = $1
	`
	row := r.db.QueryRowContext(ctx, query, slug)
	return scanRestaurant(row)
}

// Update updates a restaurant.
func (r *RestaurantRepository) Update(ctx context.Context, rest *domain.Restaurant) error {
	query := `
		UPDATE restaurants SET
			name = $2, slug = $3, cuisine_types = $4,
			logo_url = $5, cover_url = $6,
			location = ST_SetSRID(ST_MakePoint($7, $8), 4326),
			address = $9, city = $10, is_open = $11, status = $12,
			eta_min_minutes = $13, eta_max_minutes = $14,
			delivery_fee = $15, commission_rate = $16,
			opens_at = $17, closes_at = $18, updated_at = $19
		WHERE id = $1
	`

	cuisines := make([]string, len(rest.CuisineTypes()))
	for i, c := range rest.CuisineTypes() {
		cuisines[i] = string(c)
	}

	_, err := r.db.ExecContext(ctx, query,
		rest.ID(),
		rest.Name(),
		rest.Slug(),
		arrayToPostgres(cuisines),
		nullableString(rest.LogoURL()),
		nullableString(rest.CoverURL()),
		rest.Longitude(),
		rest.Latitude(),
		rest.Address(),
		rest.City(),
		rest.IsOpen(),
		string(rest.Status()),
		rest.ETAMinMinutes(),
		rest.ETAMaxMinutes(),
		rest.DeliveryFee(),
		rest.CommissionRate(),
		rest.OpensAt(),
		rest.ClosesAt(),
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("restaurant_repo.Update: %w", err)
	}
	return nil
}

// FindNearby finds restaurants within a radius using PostGIS.
func (r *RestaurantRepository) FindNearby(ctx context.Context, lat, lng, radiusKm float64, limit, offset int) ([]*domain.Restaurant, error) {
	query := `
		SELECT id, name, slug, cuisine_types, rating, rating_count,
		       COALESCE(logo_url, ''), COALESCE(cover_url, ''),
		       ST_Y(location) as lat, ST_X(location) as lng,
		       address, city, is_open, status,
		       eta_min_minutes, eta_max_minutes, delivery_fee,
		       price_range, commission_rate, opens_at, closes_at,
		       created_at, updated_at
		FROM restaurants
		WHERE status = 'active'
		  AND is_open = true
		  AND ST_DWithin(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, $3 * 1000)
		ORDER BY ST_Distance(location, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography)
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.QueryContext(ctx, query, lng, lat, radiusKm, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("restaurant_repo.FindNearby: %w", err)
	}
	defer rows.Close()

	return scanRestaurants(rows)
}

// Search searches restaurants by name.
func (r *RestaurantRepository) Search(ctx context.Context, query string, lat, lng float64, limit, offset int) ([]*domain.Restaurant, error) {
	sqlQuery := `
		SELECT id, name, slug, cuisine_types, rating, rating_count,
		       COALESCE(logo_url, ''), COALESCE(cover_url, ''),
		       ST_Y(location) as lat, ST_X(location) as lng,
		       address, city, is_open, status,
		       eta_min_minutes, eta_max_minutes, delivery_fee,
		       price_range, commission_rate, opens_at, closes_at,
		       created_at, updated_at
		FROM restaurants
		WHERE status = 'active'
		  AND (name ILIKE '%' || $1 || '%' OR address ILIKE '%' || $1 || '%')
	`

	args := []interface{}{query}
	argIdx := 2

	if lat != 0 && lng != 0 {
		sqlQuery += fmt.Sprintf(` ORDER BY ST_Distance(location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography)`, argIdx, argIdx+1)
		args = append(args, lng, lat)
		argIdx += 2
	} else {
		sqlQuery += ` ORDER BY rating DESC, rating_count DESC`
	}

	sqlQuery += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("restaurant_repo.Search: %w", err)
	}
	defer rows.Close()

	return scanRestaurants(rows)
}

// FindByCuisine finds restaurants by cuisine type.
func (r *RestaurantRepository) FindByCuisine(ctx context.Context, cuisine domain.CuisineType, lat, lng float64, limit, offset int) ([]*domain.Restaurant, error) {
	query := `
		SELECT id, name, slug, cuisine_types, rating, rating_count,
		       COALESCE(logo_url, ''), COALESCE(cover_url, ''),
		       ST_Y(location) as lat, ST_X(location) as lng,
		       address, city, is_open, status,
		       eta_min_minutes, eta_max_minutes, delivery_fee,
		       price_range, commission_rate, opens_at, closes_at,
		       created_at, updated_at
		FROM restaurants
		WHERE status = 'active'
		  AND is_open = true
		  AND $1 = ANY(cuisine_types)
	`

	args := []interface{}{string(cuisine)}
	argIdx := 2

	if lat != 0 && lng != 0 {
		query += fmt.Sprintf(` ORDER BY ST_Distance(location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography)`, argIdx, argIdx+1)
		args = append(args, lng, lat)
		argIdx += 2
	} else {
		query += ` ORDER BY rating DESC`
	}

	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("restaurant_repo.FindByCuisine: %w", err)
	}
	defer rows.Close()

	return scanRestaurants(rows)
}

// UpdateRating updates the rating and rating count.
func (r *RestaurantRepository) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	query := `
		UPDATE restaurants
		SET rating = (rating * rating_count + $2) / (rating_count + 1),
		    rating_count = rating_count + 1,
		    updated_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, rating)
	return err
}

// ============ Helpers ============

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanRestaurant(s scanner) (*domain.Restaurant, error) {
	var (
		id              uuid.UUID
		name            string
		slug            string
		cuisineTypesStr string
		rating          float64
		ratingCount     int
		logoURL         string
		coverURL        string
		lat             float64
		lng             float64
		address         string
		city            string
		isOpen          bool
		status          string
		etaMin          int
		etaMax          int
		deliveryFee     float64
		priceRange      int
		commissionRate  float64
		opensAt         string
		closesAt        string
		createdAt       time.Time
		updatedAt       time.Time
	)

	err := s.Scan(
		&id, &name, &slug, &cuisineTypesStr, &rating, &ratingCount,
		&logoURL, &coverURL, &lat, &lng,
		&address, &city, &isOpen, &status,
		&etaMin, &etaMax, &deliveryFee,
		&priceRange, &commissionRate, &opensAt, &closesAt,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrRestaurantNotFound
		}
		return nil, fmt.Errorf("scanRestaurant: %w", err)
	}

	cuisines := postgresToArray(cuisineTypesStr)

	r := &domain.Restaurant{
		id:             id,
		name:           name,
		slug:           slug,
		cuisineTypes:   cuisines,
		rating:         rating,
		ratingCount:    ratingCount,
		logoURL:        logoURL,
		coverURL:       coverURL,
		latitude:       lat,
		longitude:      lng,
		address:        address,
		city:           city,
		isOpen:         isOpen,
		status:         domain.RestaurantStatus(status),
		etaMinMinutes:  etaMin,
		etaMaxMinutes:  etaMax,
		deliveryFee:    deliveryFee,
		priceRange:     priceRange,
		commissionRate: commissionRate,
		opensAt:        opensAt,
		closesAt:       closesAt,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}

	return r, nil
}

func scanRestaurants(rows *sql.Rows) ([]*domain.Restaurant, error) {
	var restaurants []*domain.Restaurant
	for rows.Next() {
		r, err := scanRestaurant(rows)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, r)
	}
	return restaurants, rows.Err()
}

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// arrayToPostgres converts a string slice to PostgreSQL array format: {a,b,c}
func arrayToPostgres(arr []string) string {
	if len(arr) == 0 {
		return "{}"
	}
	escaped := make([]string, len(arr))
	for i, s := range arr {
		// Escape quotes and backslashes
		s = strings.ReplaceAll(s, `\`, `\\`)
		s = strings.ReplaceAll(s, `"`, `\"`)
		escaped[i] = `"` + s + `"`
	}
	return "{" + strings.Join(escaped, ",") + "}"
}

// postgresToArray converts a PostgreSQL array string back to a slice.
func postgresToArray(s string) []domain.CuisineType {
	if s == "" || s == "{}" {
		return nil
	}

	// Remove surrounding braces
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	cuisines := make([]domain.CuisineType, 0, len(parts))
	for _, p := range parts {
		p = strings.Trim(p, `"`)
		if p != "" {
			cuisines = append(cuisines, domain.CuisineType(p))
		}
	}
	return cuisines
}
