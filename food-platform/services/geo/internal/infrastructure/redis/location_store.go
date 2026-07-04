package redis

import (
	"context"
	"fmt"

	"github.com/food-platform/geo/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	driverLocationKey = "drivers:locations"
	driverSetKey      = "drivers:online"
)

type LocationStore struct {
	client *redis.Client
}

func NewLocationStore(client *redis.Client) *LocationStore {
	return &LocationStore{client: client}
}

func (s *LocationStore) UpdateLocation(ctx context.Context, driverID uuid.UUID, lat, lng float64) error {
	member := driverID.String()
	err := s.client.GeoAdd(ctx, driverLocationKey, &redis.GeoLocation{
		Name:      member,
		Longitude: lng,
		Latitude:  lat,
	}).Err()
	if err != nil {
		return fmt.Errorf("redis.GeoAdd: %w", err)
	}
	s.client.SAdd(ctx, driverSetKey, member)
	return nil
}

func (s *LocationStore) GetLocation(ctx context.Context, driverID uuid.UUID) (lat, lng float64, err error) {
	member := driverID.String()
	positions, err := s.client.GeoPos(ctx, driverLocationKey, member).Result()
	if err != nil {
		return 0, 0, fmt.Errorf("redis.GeoPos: %w", err)
	}
	if len(positions) == 0 || positions[0] == nil {
		return 0, 0, domain.ErrDriverNotFound
	}
	return positions[0].Latitude, positions[0].Longitude, nil
}

func (s *LocationStore) FindNearby(ctx context.Context, lat, lng, radiusMeters float64, count int) ([]*domain.NearbyDriver, error) {
	// Use GeoSearch with WithCoord and WithDist to get full location info
	results, err := s.client.GeoSearchLocation(ctx, driverLocationKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  lng,
			Latitude:   lat,
			Radius:     radiusMeters,
			RadiusUnit: "m",
			Sort:       "ASC",
			Count:      count,
		},
		WithCoord: true,
		WithDist:  true,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("redis.GeoSearchLocation: %w", err)
	}

	drivers := make([]*domain.NearbyDriver, 0, len(results))
	for _, result := range results {
		driverID, err := uuid.Parse(result.Name)
		if err != nil {
			continue
		}
		drivers = append(drivers, domain.NewNearbyDriver(
			driverID,
			result.Latitude,
			result.Longitude,
			result.Dist,
		))
	}
	return drivers, nil
}

func (s *LocationStore) RemoveDriver(ctx context.Context, driverID uuid.UUID) error {
	member := driverID.String()
	s.client.ZRem(ctx, driverLocationKey, member)
	s.client.SRem(ctx, driverSetKey, member)
	return nil
}
