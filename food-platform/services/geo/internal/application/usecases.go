// Package application contains use cases for the Geo Service.
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/food-platform/geo/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

// LocationStore is the interface for storing/retrieving driver locations (Redis GEO).
type LocationStore interface {
	UpdateLocation(ctx context.Context, driverID uuid.UUID, lat, lng float64) error
	GetLocation(ctx context.Context, driverID uuid.UUID) (lat, lng float64, err error)
	FindNearby(ctx context.Context, lat, lng, radiusMeters float64, count int) ([]*domain.NearbyDriver, error)
	RemoveDriver(ctx context.Context, driverID uuid.UUID) error
}

// ============ DTOs ============

type LocationDTO struct {
	DriverID  string `json:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Speed     float64 `json:"speed"`
}

type NearbyDriverDTO struct {
	DriverID   string  `json:"driver_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	DistanceKm float64 `json:"distance_km"`
}

type ETADTO struct {
	OriginLat float64 `json:"origin_lat"`
	OriginLng float64 `json:"origin_lng"`
	DestLat   float64 `json:"dest_lat"`
	DestLng   float64 `json:"dest_lng"`
	ETAMinutes int    `json:"eta_minutes"`
	DistanceKm float64 `json:"distance_km"`
}

// ============ Commands ============

type UpdateLocationCommand struct {
	DriverID uuid.UUID
	Lat      float64
	Lng      float64
	Heading  float64
	Speed    float64
}

type FindNearbyCommand struct {
	Lat        float64
	Lng        float64
	RadiusKm   float64
	Count      int
}

type CalculateETACommand struct {
	OriginLat float64
	OriginLng float64
	DestLat   float64
	DestLng   float64
}

// ============ Use Cases ============

type UpdateLocationUseCase struct {
	store LocationStore
}

func NewUpdateLocationUseCase(store LocationStore) *UpdateLocationUseCase {
	return &UpdateLocationUseCase{store: store}
}

func (uc *UpdateLocationUseCase) Execute(ctx context.Context, cmd UpdateLocationCommand) error {
	loc, err := domain.NewDriverLocation(cmd.DriverID, cmd.Lat, cmd.Lng, cmd.Heading, cmd.Speed)
	if err != nil {
		return err
	}

	if err := uc.store.UpdateLocation(ctx, loc.DriverID(), loc.Latitude(), loc.Longitude()); err != nil {
		slog.Error("update_location_failed", "driver_id", cmd.DriverID, "error", err)
		return fmt.Errorf("update location: %w", err)
	}

	return nil
}

type GetLocationUseCase struct {
	store LocationStore
}

func NewGetLocationUseCase(store LocationStore) *GetLocationUseCase {
	return &GetLocationUseCase{store: store}
}

func (uc *GetLocationUseCase) Execute(ctx context.Context, driverID uuid.UUID) (*LocationDTO, error) {
	lat, lng, err := uc.store.GetLocation(ctx, driverID)
	if err != nil {
		return nil, err
	}

	return &LocationDTO{
		DriverID:  driverID.String(),
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

type FindNearbyUseCase struct {
	store LocationStore
}

func NewFindNearbyUseCase(store LocationStore) *FindNearbyUseCase {
	return &FindNearbyUseCase{store: store}
}

func (uc *FindNearbyUseCase) Execute(ctx context.Context, cmd FindNearbyCommand) ([]*NearbyDriverDTO, error) {
	if cmd.Count == 0 {
		cmd.Count = 40
	}
	if cmd.RadiusKm == 0 {
		cmd.RadiusKm = 3
	}

	drivers, err := uc.store.FindNearby(ctx, cmd.Lat, cmd.Lng, cmd.RadiusKm*1000, cmd.Count)
	if err != nil {
		return nil, fmt.Errorf("find nearby: %w", err)
	}

	dtos := make([]*NearbyDriverDTO, 0, len(drivers))
	for _, d := range drivers {
		dtos = append(dtos, &NearbyDriverDTO{
			DriverID:   d.DriverID().String(),
			Latitude:   d.Latitude(),
			Longitude:  d.Longitude(),
			DistanceKm: d.DistanceKm(),
		})
	}
	return dtos, nil
}

type CalculateETAUseCase struct{}

func NewCalculateETAUseCase() *CalculateETAUseCase {
	return &CalculateETAUseCase{}
}

func (uc *CalculateETAUseCase) Execute(ctx context.Context, cmd CalculateETACommand) (*ETADTO, error) {
	req := domain.NewETARequest(cmd.OriginLat, cmd.OriginLng, cmd.DestLat, cmd.DestLng)
	eta := domain.CalculateETA(req)

	return &ETADTO{
		OriginLat:  cmd.OriginLat,
		OriginLng:  cmd.OriginLng,
		DestLat:    cmd.DestLat,
		DestLng:    cmd.DestLng,
		ETAMinutes: eta,
		DistanceKm: 0, // Calculated inside domain
	}, nil
}
