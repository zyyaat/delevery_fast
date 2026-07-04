// Package domain contains the core logic of the Geo Service.
package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCoordinates = errors.New("invalid coordinates")
	ErrDriverNotFound     = errors.New("driver location not found")
)

// DriverLocation represents a driver's real-time GPS position.
type DriverLocation struct {
	driverID  uuid.UUID
	latitude  float64
	longitude float64
	heading   float64
	speed     float64
	timestamp time.Time
}

func NewDriverLocation(driverID uuid.UUID, lat, lng, heading, speed float64) (*DriverLocation, error) {
	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		return nil, ErrInvalidCoordinates
	}
	return &DriverLocation{
		driverID:  driverID,
		latitude:  lat,
		longitude: lng,
		heading:   heading,
		speed:     speed,
		timestamp: time.Now().UTC(),
	}, nil
}

func (dl *DriverLocation) DriverID() uuid.UUID { return dl.driverID }
func (dl *DriverLocation) Latitude() float64   { return dl.latitude }
func (dl *DriverLocation) Longitude() float64  { return dl.longitude }
func (dl *DriverLocation) Heading() float64    { return dl.heading }
func (dl *DriverLocation) Speed() float64      { return dl.speed }
func (dl *DriverLocation) Timestamp() time.Time { return dl.timestamp }

// NearbyDriver represents a driver found in a geospatial search.
type NearbyDriver struct {
	driverID  uuid.UUID
	latitude  float64
	longitude float64
	distance  float64 // meters
	timestamp time.Time
}

func NewNearbyDriver(driverID uuid.UUID, lat, lng, distance float64) *NearbyDriver {
	return &NearbyDriver{
		driverID:  driverID,
		latitude:  lat,
		longitude: lng,
		distance:  distance,
		timestamp: time.Now().UTC(),
	}
}

func (nd *NearbyDriver) DriverID() uuid.UUID { return nd.driverID }
func (nd *NearbyDriver) Latitude() float64   { return nd.latitude }
func (nd *NearbyDriver) Longitude() float64  { return nd.longitude }
func (nd *NearbyDriver) Distance() float64   { return nd.distance }
func (nd *NearbyDriver) DistanceKm() float64 { return nd.distance / 1000 }
func (nd *NearbyDriver) Timestamp() time.Time { return nd.timestamp }

// ETARequest represents an ETA calculation request.
type ETARequest struct {
	originLat      float64
	originLng      float64
	destLat        float64
	destLng        float64
	trafficFactor  float64 // 1.0 = normal, 1.5 = heavy traffic
}

func NewETARequest(originLat, originLng, destLat, destLng float64) *ETARequest {
	return &ETARequest{
		originLat:     originLat,
		originLng:     originLng,
		destLat:       destLat,
		destLng:       destLng,
		trafficFactor: 1.0,
	}
}

func (r *ETARequest) SetTrafficFactor(factor float64) {
	if factor < 0.5 {
		factor = 0.5
	} else if factor > 2.0 {
		factor = 2.0
	}
	r.trafficFactor = factor
}

// CalculateETA estimates travel time in minutes using Haversine distance.
// Assumes average urban speed of 25 km/h with traffic factor.
func CalculateETA(req *ETARequest) int {
	distanceKm := haversine(req.originLat, req.originLng, req.destLat, req.destLng)
	avgSpeedKmH := 25.0 / req.trafficFactor
	etaMinutes := (distanceKm / avgSpeedKmH) * 60

	// Minimum 5 minutes, round to nearest minute
	eta := int(math.Round(etaMinutes))
	if eta < 5 {
		return 5
	}
	return eta
}

func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}
