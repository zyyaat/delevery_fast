// Package domain contains the core logic of the Delivery Matching Service.
// This is the most critical algorithm in the platform — it finds the best driver
// for each order using geospatial search + multi-factor scoring.
package domain

import (
	"errors"
	"math"
	"sort"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoDriversAvailable = errors.New("no drivers available in area")
	ErrOrderAlreadyMatched = errors.New("order already has a driver assigned")
	ErrInvalidRadius       = errors.New("invalid search radius")
)

// ============ Entities ============

// OrderLocation represents the pickup and dropoff locations for an order.
type OrderLocation struct {
	OrderID       uuid.UUID
	RestaurantLat float64
	RestaurantLng float64
	CustomerLat   float64
	CustomerLng   float64
}

// DriverCandidate represents a driver found in the geospatial search,
// before scoring.
type DriverCandidate struct {
	DriverID    uuid.UUID
	Latitude    float64
	Longitude   float64
	Distance    float64 // meters from restaurant
	Rating      float64
	AcceptanceRate float64 // 0-1
	CompletionRate float64 // 0-1
	IdleTime    time.Duration // time since last delivery
	VehicleType string
}

// ScoredDriver is a driver with a calculated match score.
type ScoredDriver struct {
	Candidate DriverCandidate
	Score     float64 // 0-100
}

// MatchResult represents the outcome of a matching attempt.
type MatchResult struct {
	OrderID     uuid.UUID
	TopDrivers  []ScoredDriver // Top 3 to broadcast to
	Round       int            // 1, 2, or 3
	RadiusKm    float64
	Matched     bool
	MatchedDriverID *uuid.UUID
}

// ============ Scoring Weights ============

const (
	WeightDistance      = 0.30 // 30% — closer is better
	WeightAcceptance    = 0.20 // 20% — higher acceptance rate is better
	WeightRating        = 0.15 // 15% — higher rating is better
	WeightIdleTime      = 0.15 // 15% — longer idle time is better (fairness)
	WeightCompletion    = 0.10 // 10% — higher completion rate is better
	WeightVehicleMatch  = 0.10 // 10% — motorcycle preferred for short distances
)

// ============ Matching Algorithm ============

// MatchDriver scores all candidates and returns the top N.
// This is the core algorithm inspired by Uber's dispatch system.
func MatchDriver(order OrderLocation, candidates []DriverCandidate, topN int) []ScoredDriver {
	if len(candidates) == 0 {
		return nil
	}

	scored := make([]ScoredDriver, 0, len(candidates))

	for _, c := range candidates {
		score := calculateScore(order, c)
		scored = append(scored, ScoredDriver{
			Candidate: c,
			Score:     score,
		})
	}

	// Sort by score descending
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Return top N
	if len(scored) > topN {
		scored = scored[:topN]
	}

	return scored
}

// calculateScore computes a 0-100 score for a driver candidate.
func calculateScore(order OrderLocation, c DriverCandidate) float64 {
	// 1. Distance score (30%) — closer is better, linear decay after 1.5km
	distanceScore := 100.0
	distanceKm := c.Distance / 1000
	if distanceKm > 1.5 {
		// Linear decay: at 3km, score = 0
		distanceScore = math.Max(0, 100*(1-(distanceKm-1.5)/1.5))
	}

	// 2. Acceptance rate score (20%) — 0-1 scaled to 0-100
	acceptanceScore := c.AcceptanceRate * 100

	// 3. Rating score (15%) — 0-5 scaled to 0-100
	ratingScore := (c.Rating / 5.0) * 100

	// 4. Idle time score (15%) — longer idle = higher score (fairness)
	// Cap at 2 hours (120 min) for max score
	idleMinutes := c.IdleTime.Minutes()
	idleScore := math.Min(100, (idleMinutes/120)*100)

	// 5. Completion rate score (10%) — 0-1 scaled to 0-100
	completionScore := c.CompletionRate * 100

	// 6. Vehicle match score (10%) — motorcycle preferred for short distances
	vehicleScore := 50.0 // default
	totalDistance := distanceKm + haversineDistance(
		order.RestaurantLat, order.RestaurantLng,
		order.CustomerLat, order.CustomerLng,
	)
	if c.VehicleType == "motorcycle" && totalDistance < 10 {
		vehicleScore = 100
	} else if c.VehicleType == "car" && totalDistance >= 10 {
		vehicleScore = 100
	} else if c.VehicleType == "bicycle" && totalDistance < 3 {
		vehicleScore = 100
	}

	// Weighted sum
	totalScore := WeightDistance*distanceScore +
		WeightAcceptance*acceptanceScore +
		WeightRating*ratingScore +
		WeightIdleTime*idleScore +
		WeightCompletion*completionScore +
		WeightVehicleMatch*vehicleScore

	return math.Round(totalScore*100) / 100
}

// haversineDistance calculates distance between two points in km.
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// ============ Dispatch Strategy ============

// DispatchConfig holds the configuration for the dispatch algorithm.
type DispatchConfig struct {
	InitialRadiusKm     float64 // 3km
	MaxRadiusKm         float64 // 5km
	MaxRounds           int     // 3
	InitialTimeoutSec   int     // 15s (per round)
	BroadcastCount      int     // 3 (top 3 drivers)
}

func DefaultDispatchConfig() DispatchConfig {
	return DispatchConfig{
		InitialRadiusKm:   3.0,
		MaxRadiusKm:       5.0,
		MaxRounds:         3,
		InitialTimeoutSec: 15,
		BroadcastCount:    3,
	}
}

// DispatchRound represents one round of the dispatch process.
type DispatchRound struct {
	RoundNumber int
	RadiusKm    float64
	DriverCount int
	TimeoutSec  int
	Drivers     []ScoredDriver
}

// NextRound calculates the parameters for the next dispatch round.
func (r *DispatchRound) NextRound(cfg DispatchConfig) *DispatchRound {
	if r.RoundNumber >= cfg.MaxRounds {
		return nil // No more rounds
	}

	next := &DispatchRound{
		RoundNumber: r.RoundNumber + 1,
		RadiusKm:    r.RadiusKm + (cfg.MaxRadiusKm-cfg.InitialRadiusKm)/float64(cfg.MaxRounds-1),
		TimeoutSec:  r.TimeoutSec + 15, // Add 15s each round
	}

	if next.RadiusKm > cfg.MaxRadiusKm {
		next.RadiusKm = cfg.MaxRadiusKm
	}

	return next
}
