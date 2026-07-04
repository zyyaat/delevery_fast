// Package application contains use cases for the Delivery Matching Service.
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/food-platform/delivery-matching/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

type GeoService interface {
	FindNearbyDrivers(ctx context.Context, lat, lng, radiusKm float64, count int) ([]DriverLocation, error)
}

type DriverService interface {
	GetDriverInfo(ctx context.Context, driverID uuid.UUID) (*DriverInfo, error)
}

type OrderService interface {
	UpdateOrderDriver(ctx context.Context, orderID, driverID uuid.UUID) error
}

type EventPublisher interface {
	PublishOrderAssigned(ctx context.Context, orderID, driverID uuid.UUID) error
	PublishDispatchFailed(ctx context.Context, orderID uuid.UUID, reason string) error
}

// ============ DTOs ============

type DriverLocation struct {
	DriverID   uuid.UUID
	Latitude   float64
	Longitude  float64
	DistanceM  float64
}

type DriverInfo struct {
	DriverID       uuid.UUID
	Rating         float64
	AcceptanceRate float64
	CompletionRate float64
	VehicleType    string
	IdleTime       string // duration string
}

type MatchOrderCommand struct {
	OrderID       uuid.UUID
	RestaurantLat float64
	RestaurantLng float64
	CustomerLat   float64
	CustomerLng   float64
}

type MatchResultDTO struct {
	OrderID    string   `json:"order_id"`
	Matched    bool     `json:"matched"`
	DriverIDs  []string `json:"driver_ids,omitempty"`
	Round      int      `json:"round"`
	RadiusKm   float64  `json:"radius_km"`
	TopScores  []ScoredDriverDTO `json:"top_scores,omitempty"`
}

type ScoredDriverDTO struct {
	DriverID  string  `json:"driver_id"`
	Score     float64 `json:"score"`
	Distance  float64 `json:"distance_km"`
	Rating    float64 `json:"rating"`
}

// ============ Use Cases ============

type MatchOrderUseCase struct {
	geoService    GeoService
	driverService DriverService
	orderService  OrderService
	publisher     EventPublisher
	config        domain.DispatchConfig
}

func NewMatchOrderUseCase(
	geo GeoService,
	driver DriverService,
	order OrderService,
	pub EventPublisher,
) *MatchOrderUseCase {
	return &MatchOrderUseCase{
		geoService:    geo,
		driverService: driver,
		orderService:  order,
		publisher:     pub,
		config:        domain.DefaultDispatchConfig(),
	}
}

// Execute runs the full matching algorithm for an order.
// In production, this would be called by a Kafka consumer listening for
// "order.confirmed" events.
func (uc *MatchOrderUseCase) Execute(ctx context.Context, cmd MatchOrderCommand) (*MatchResultDTO, error) {
	orderLoc := domain.OrderLocation{
		OrderID:       cmd.OrderID,
		RestaurantLat: cmd.RestaurantLat,
		RestaurantLng: cmd.RestaurantLng,
		CustomerLat:   cmd.CustomerLat,
		CustomerLng:   cmd.CustomerLng,
	}

	// Run dispatch rounds
	for roundNum := 1; roundNum <= uc.config.MaxRounds; roundNum++ {
		radius := uc.config.InitialRadiusKm
		if roundNum > 1 {
			radius = uc.config.InitialRadiusKm +
				(float64(roundNum-1) * (uc.config.MaxRadiusKm - uc.config.InitialRadiusKm) / float64(uc.config.MaxRounds-1))
		}

		slog.Info("dispatch_round",
			"order_id", cmd.OrderID,
			"round", roundNum,
			"radius_km", radius,
		)

		// Find nearby drivers via Geo Service
		nearbyDrivers, err := uc.geoService.FindNearbyDrivers(ctx, orderLoc.RestaurantLat, orderLoc.RestaurantLng, radius, 40)
		if err != nil {
			slog.Error("find_nearby_failed", "error", err)
			continue
		}

		if len(nearbyDrivers) == 0 {
			slog.Info("no_drivers_in_round", "round", roundNum, "radius", radius)
			continue
		}

		// Get driver info for scoring
		candidates := make([]domain.DriverCandidate, 0, len(nearbyDrivers))
		for _, nd := range nearbyDrivers {
			info, err := uc.driverService.GetDriverInfo(ctx, nd.DriverID)
			if err != nil {
				slog.Warn("get_driver_info_failed", "driver_id", nd.DriverID, "error", err)
				continue
			}

			candidates = append(candidates, domain.DriverCandidate{
				DriverID:       nd.DriverID,
				Latitude:       nd.Latitude,
				Longitude:      nd.Longitude,
				Distance:       nd.DistanceM,
				Rating:         info.Rating,
				AcceptanceRate: info.AcceptanceRate,
				CompletionRate: info.CompletionRate,
				VehicleType:    info.VehicleType,
			})
		}

		if len(candidates) == 0 {
			continue
		}

		// Score and rank drivers
		scored := domain.MatchDriver(orderLoc, candidates, uc.config.BroadcastCount)

		if len(scored) > 0 {
			// Broadcast to top N drivers (they'll receive push notifications)
			driverIDs := make([]string, 0, len(scored))
			scores := make([]ScoredDriverDTO, 0, len(scored))
			for _, s := range scored {
				driverIDs = append(driverIDs, s.Candidate.DriverID.String())
				scores = append(scores, ScoredDriverDTO{
					DriverID: s.Candidate.DriverID.String(),
					Score:    s.Score,
					Distance: s.Candidate.Distance / 1000,
					Rating:   s.Candidate.Rating,
				})
			}

			slog.Info("drivers_broadcast",
				"order_id", cmd.OrderID,
				"round", roundNum,
				"drivers", len(scored),
				"top_score", scored[0].Score,
			)

			return &MatchResultDTO{
				OrderID:   cmd.OrderID.String(),
				Matched:   false, // Will be true when a driver accepts
				DriverIDs: driverIDs,
				Round:     roundNum,
				RadiusKm:  radius,
				TopScores: scores,
			}, nil
		}
	}

	// All rounds exhausted — no driver found
	slog.Warn("dispatch_failed", "order_id", cmd.OrderID)

	_ = uc.publisher.PublishDispatchFailed(ctx, cmd.OrderID, "no drivers available after 3 rounds")

	return &MatchResultDTO{
		OrderID: cmd.OrderID.String(),
		Matched: false,
		Round:   uc.config.MaxRounds,
		RadiusKm: uc.config.MaxRadiusKm,
	}, fmt.Errorf("no drivers available")
}

// ============ Accept Order Use Case ============

type AcceptOrderUseCase struct {
	orderService OrderService
	publisher    EventPublisher
}

func NewAcceptOrderUseCase(order OrderService, pub EventPublisher) *AcceptOrderUseCase {
	return &AcceptOrderUseCase{orderService: order, publisher: pub}
}

type AcceptOrderCommand struct {
	OrderID  uuid.UUID
	DriverID uuid.UUID
}

func (uc *AcceptOrderUseCase) Execute(ctx context.Context, cmd AcceptOrderCommand) error {
	// Assign driver to order
	if err := uc.orderService.UpdateOrderDriver(ctx, cmd.OrderID, cmd.DriverID); err != nil {
		return fmt.Errorf("assign driver: %w", err)
	}

	// Publish event
	if err := uc.publisher.PublishOrderAssigned(ctx, cmd.OrderID, cmd.DriverID); err != nil {
		slog.Error("publish_order_assigned_failed", "error", err)
	}

	slog.Info("order_assigned",
		"order_id", cmd.OrderID,
		"driver_id", cmd.DriverID,
	)

	return nil
}
