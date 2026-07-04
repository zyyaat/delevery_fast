// Package domain contains the core business logic of the Restaurant Catalog Service.
package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// ============ Errors ============

var (
	ErrRestaurantNotFound  = errors.New("restaurant not found")
	ErrRestaurantClosed    = errors.New("restaurant is closed")
	ErrRestaurantSuspended = errors.New("restaurant is suspended")
	ErrInvalidCoordinates  = errors.New("invalid coordinates")
	ErrInvalidName         = errors.New("invalid restaurant name")
)

// ============ Enums ============

// RestaurantStatus represents the operational status of a restaurant.
type RestaurantStatus string

const (
	RestaurantStatusActive             RestaurantStatus = "active"
	RestaurantStatusPaused             RestaurantStatus = "paused"
	RestaurantStatusSuspended          RestaurantStatus = "suspended"
	RestaurantStatusPendingVerification RestaurantStatus = "pending_verification"
	RestaurantStatusRejected           RestaurantStatus = "rejected"
)

// CuisineType represents a type of cuisine.
type CuisineType string

const (
	CuisineEgyptian    CuisineType = "egyptian"
	CuisineItalian     CuisineType = "italian"
	CuisineAsian       CuisineType = "asian"
	CuisineFastFood    CuisineType = "fast_food"
	CuisineHealthy     CuisineType = "healthy"
	CuisineDesserts    CuisineType = "desserts"
	CuisineIndian      CuisineType = "indian"
	CuisineLebanese    CuisineType = "lebanese"
	CuisineBreakfast   CuisineType = "breakfast"
	CuisineCoffee      CuisineType = "coffee"
	CuisineGroceries   CuisineType = "groceries"
)

// ============ Validation ============

// ValidateCoordinates checks if latitude and longitude are valid.
func ValidateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return ErrInvalidCoordinates
	}
	if lng < -180 || lng > 180 {
		return ErrInvalidCoordinates
	}
	return nil
}

// ============ Entities ============

// Restaurant represents a restaurant in the catalog.
type Restaurant struct {
	id              uuid.UUID
	name            string
	slug            string
	cuisineTypes    []CuisineType
	rating          float64
	ratingCount     int
	logoURL         string
	coverURL        string
	latitude        float64
	longitude       float64
	address         string
	city            string
	isOpen          bool
	status          RestaurantStatus
	etaMinMinutes   int
	etaMaxMinutes   int
	deliveryFee     float64
	priceRange      int // 1-4
	commissionRate  float64
	opensAt         string // "10:00"
	closesAt        string // "02:00"
	createdAt       time.Time
	updatedAt       time.Time
}

// NewRestaurant creates a new Restaurant with validation.
func NewRestaurant(
	name string,
	lat, lng float64,
	address, city string,
	cuisineTypes []CuisineType,
) (*Restaurant, error) {
	if name == "" {
		return nil, ErrInvalidName
	}
	if err := ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}
	if address == "" {
		return nil, errors.New("address is required")
	}

	now := time.Now().UTC()
	return &Restaurant{
		id:             uuid.New(),
		name:           name,
		slug:           generateSlug(name),
		cuisineTypes:   cuisineTypes,
		rating:         0,
		ratingCount:    0,
		latitude:       lat,
		longitude:      lng,
		address:        address,
		city:           city,
		isOpen:         true,
		status:         RestaurantStatusPendingVerification,
		etaMinMinutes:  20,
		etaMaxMinutes:  40,
		deliveryFee:    20.0,
		priceRange:     2,
		commissionRate: 0.15,
		opensAt:        "10:00",
		closesAt:       "23:59",
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// ============ Getters ============

func (r *Restaurant) ID() uuid.UUID           { return r.id }
func (r *Restaurant) Name() string            { return r.name }
func (r *Restaurant) Slug() string            { return r.slug }
func (r *Restaurant) CuisineTypes() []CuisineType { return r.cuisineTypes }
func (r *Restaurant) Rating() float64         { return r.rating }
func (r *Restaurant) RatingCount() int        { return r.ratingCount }
func (r *Restaurant) LogoURL() string         { return r.logoURL }
func (r *Restaurant) CoverURL() string        { return r.coverURL }
func (r *Restaurant) Latitude() float64       { return r.latitude }
func (r *Restaurant) Longitude() float64      { return r.longitude }
func (r *Restaurant) Address() string         { return r.address }
func (r *Restaurant) City() string            { return r.city }
func (r *Restaurant) IsOpen() bool            { return r.isOpen }
func (r *Restaurant) Status() RestaurantStatus { return r.status }
func (r *Restaurant) ETAMinMinutes() int      { return r.etaMinMinutes }
func (r *Restaurant) ETAMaxMinutes() int      { return r.etaMaxMinutes }
func (r *Restaurant) DeliveryFee() float64    { return r.deliveryFee }
func (r *Restaurant) PriceRange() int         { return r.priceRange }
func (r *Restaurant) CommissionRate() float64 { return r.commissionRate }
func (r *Restaurant) OpensAt() string         { return r.opensAt }
func (r *Restaurant) ClosesAt() string        { return r.closesAt }
func (r *Restaurant) CreatedAt() time.Time    { return r.createdAt }
func (r *Restaurant) UpdatedAt() time.Time    { return r.updatedAt }

// ============ Setters ============

// SetName updates the restaurant name.
func (r *Restaurant) SetName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	r.name = name
	r.slug = generateSlug(name)
	r.updatedAt = time.Now().UTC()
	return nil
}

// SetLogoURL sets the logo URL.
func (r *Restaurant) SetLogoURL(url string) {
	r.logoURL = url
	r.updatedAt = time.Now().UTC()
}

// SetCoverURL sets the cover image URL.
func (r *Restaurant) SetCoverURL(url string) {
	r.coverURL = url
	r.updatedAt = time.Now().UTC()
}

// SetLocation updates the restaurant location.
func (r *Restaurant) SetLocation(lat, lng float64, address string) error {
	if err := ValidateCoordinates(lat, lng); err != nil {
		return err
	}
	r.latitude = lat
	r.longitude = lng
	r.address = address
	r.updatedAt = time.Now().UTC()
	return nil
}

// SetHours updates the operating hours.
func (r *Restaurant) SetHours(opensAt, closesAt string) {
	r.opensAt = opensAt
	r.closesAt = closesAt
	r.updatedAt = time.Now().UTC()
}

// SetETA updates the estimated delivery time range.
func (r *Restaurant) SetETA(min, max int) {
	r.etaMinMinutes = min
	r.etaMaxMinutes = max
	r.updatedAt = time.Now().UTC()
}

// SetDeliveryFee updates the delivery fee.
func (r *Restaurant) SetDeliveryFee(fee float64) {
	if fee < 0 {
		fee = 0
	}
	r.deliveryFee = fee
	r.updatedAt = time.Now().UTC()
}

// SetCommissionRate updates the commission rate (0.0 to 1.0).
func (r *Restaurant) SetCommissionRate(rate float64) {
	if rate < 0 {
		rate = 0
	} else if rate > 1 {
		rate = 1
	}
	r.commissionRate = rate
	r.updatedAt = time.Now().UTC()
}

// SetOpen toggles the restaurant open/closed status.
func (r *Restaurant) SetOpen(open bool) {
	r.isOpen = open
	r.updatedAt = time.Now().UTC()
}

// SetStatus updates the restaurant status.
func (r *Restaurant) SetStatus(status RestaurantStatus) {
	r.status = status
	r.updatedAt = time.Now().UTC()
}

// UpdateRating recalculates the rating from a new rating.
func (r *Restaurant) UpdateRating(newRating float64) {
	if newRating < 0 {
		newRating = 0
	} else if newRating > 5 {
		newRating = 5
	}

	total := r.rating * float64(r.ratingCount)
	r.ratingCount++
	r.rating = (total + newRating) / float64(r.ratingCount)
	r.updatedAt = time.Now().UTC()
}

// ============ Business Logic ============

// IsAcceptingOrders returns true if the restaurant can accept new orders.
func (r *Restaurant) IsAcceptingOrders() bool {
	return r.status == RestaurantStatusActive && r.isOpen
}

// DistanceTo calculates the distance from the restaurant to a point (in km) using Haversine formula.
func (r *Restaurant) DistanceTo(lat, lng float64) float64 {
	return haversine(r.latitude, r.longitude, lat, lng)
}

// haversine calculates the great-circle distance between two points in km.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0

	dLat := toRadians(lat2 - lat1)
	dLng := toRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// generateSlug creates a URL-friendly slug from a name.
// This is a simple implementation; production would use a proper slug library.
func generateSlug(name string) string {
	result := make([]byte, 0, len(name))
	prevDash := false

	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			result = append(result, c)
			prevDash = false
		} else if c >= 'A' && c <= 'Z' {
			result = append(result, c+32) // to lowercase
			prevDash = false
		} else if c == ' ' || c == '-' || c == '_' {
			if !prevDash && len(result) > 0 {
				result = append(result, '-')
				prevDash = true
			}
		}
	}

	// Trim trailing dash
	if len(result) > 0 && result[len(result)-1] == '-' {
		result = result[:len(result)-1]
	}

	if len(result) == 0 {
		return uuid.New().String()[:8]
	}

	return string(result)
}
