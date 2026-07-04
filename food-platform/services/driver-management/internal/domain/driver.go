// Package domain contains the core business logic of the Driver Management Service.
package domain

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

// ============ Errors ============

var (
	ErrDriverNotFound      = errors.New("driver not found")
	ErrDriverSuspended     = errors.New("driver is suspended")
	ErrDriverOffline       = errors.New("driver is offline")
	ErrInvalidLicense      = errors.New("invalid license number")
	ErrInvalidVehicle      = errors.New("invalid vehicle info")
	ErrKYCPending          = errors.New("driver KYC not verified")
	ErrTierNotMet          = errors.New("driver does not meet tier requirements")
	ErrInvalidCoordinates  = errors.New("invalid coordinates")
)

// ============ Enums ============

type DriverStatus string

const (
	DriverStatusOffline     DriverStatus = "offline"
	DriverStatusOnline      DriverStatus = "online"
	DriverStatusOnBreak     DriverStatus = "on_break"
	DriverStatusOnDelivery  DriverStatus = "on_delivery"
	DriverStatusSuspended   DriverStatus = "suspended"
)

type DriverTier string

const (
	TierPlatinum DriverTier = "platinum"
	TierGold     DriverTier = "gold"
	TierSilver   DriverTier = "silver"
	TierStandard DriverTier = "standard"
)

type VehicleType string

const (
	VehicleMotorcycle VehicleType = "motorcycle"
	VehicleCar        VehicleType = "car"
	VehicleBicycle    VehicleType = "bicycle"
)

type KYCStatus string

const (
	KYCPending   KYCStatus = "pending"
	KYCVerified  KYCStatus = "verified"
	KYCRejected  KYCStatus = "rejected"
)

// ============ Validation ============

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

// Driver represents a delivery driver.
type Driver struct {
	id              uuid.UUID
	userID          uuid.UUID
	name            string
	phone           string
	vehicleType     VehicleType
	vehiclePlate    string
	licenseNumber   string
	kycStatus       KYCStatus
	status          DriverStatus
	tier            DriverTier
	rating          float64
	ratingCount     int
	acceptanceRate  float64 // 0-1
	completionRate  float64 // 0-1
	trustScore      int     // 0-100
	totalEarnings   float64
	totalDeliveries int
	latitude        float64
	longitude       float64
	heading         float64
	speed           float64
	lastOnlineAt    *time.Time
	photoURL        string
	createdAt       time.Time
	updatedAt       time.Time
}

// NewDriver creates a new Driver with validation.
func NewDriver(userID uuid.UUID, name, phone string, vehicleType VehicleType) (*Driver, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if phone == "" {
		return nil, errors.New("phone is required")
	}

	now := time.Now().UTC()
	return &Driver{
		id:             uuid.New(),
		userID:         userID,
		name:           name,
		phone:          phone,
		vehicleType:    vehicleType,
		kycStatus:      KYCPending,
		status:         DriverStatusOffline,
		tier:           TierStandard,
		rating:         0,
		ratingCount:    0,
		acceptanceRate: 0,
		completionRate: 0,
		trustScore:     50,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// ============ Getters ============

func (d *Driver) ID() uuid.UUID        { return d.id }
func (d *Driver) UserID() uuid.UUID    { return d.userID }
func (d *Driver) Name() string         { return d.name }
func (d *Driver) Phone() string        { return d.phone }
func (d *Driver) VehicleType() VehicleType { return d.vehicleType }
func (d *Driver) VehiclePlate() string { return d.vehiclePlate }
func (d *Driver) LicenseNumber() string { return d.licenseNumber }
func (d *Driver) KYCStatus() KYCStatus { return d.kycStatus }
func (d *Driver) Status() DriverStatus { return d.status }
func (d *Driver) Tier() DriverTier     { return d.tier }
func (d *Driver) Rating() float64      { return d.rating }
func (d *Driver) RatingCount() int     { return d.ratingCount }
func (d *Driver) AcceptanceRate() float64 { return d.acceptanceRate }
func (d *Driver) CompletionRate() float64 { return d.completionRate }
func (d *Driver) TrustScore() int      { return d.trustScore }
func (d *Driver) TotalEarnings() float64 { return d.totalEarnings }
func (d *Driver) TotalDeliveries() int { return d.totalDeliveries }
func (d *Driver) Latitude() float64   { return d.latitude }
func (d *Driver) Longitude() float64  { return d.longitude }
func (d *Driver) Heading() float64    { return d.heading }
func (d *Driver) Speed() float64      { return d.speed }
func (d *Driver) LastOnlineAt() *time.Time { return d.lastOnlineAt }
func (d *Driver) PhotoURL() string    { return d.photoURL }
func (d *Driver) CreatedAt() time.Time { return d.createdAt }
func (d *Driver) UpdatedAt() time.Time { return d.updatedAt }

// ============ Setters ============

func (d *Driver) SetVehicleInfo(plate, license string) {
	d.vehiclePlate = plate
	d.licenseNumber = license
	d.updatedAt = time.Now().UTC()
}

func (d *Driver) SetPhotoURL(url string) {
	d.photoURL = url
	d.updatedAt = time.Now().UTC()
}

func (d *Driver) SetKYCVerified() {
	d.kycStatus = KYCVerified
	d.updatedAt = time.Now().UTC()
}

func (d *Driver) SetKYCRejected() {
	d.kycStatus = KYCRejected
	d.updatedAt = time.Now().UTC()
}

func (d *Driver) SetLocation(lat, lng, heading, speed float64) error {
	if err := ValidateCoordinates(lat, lng); err != nil {
		return err
	}
	d.latitude = lat
	d.longitude = lng
	d.heading = heading
	d.speed = speed
	d.updatedAt = time.Now().UTC()
	return nil
}

// ============ Status Management ============

// GoOnline sets the driver status to online.
func (d *Driver) GoOnline() error {
	if d.kycStatus != KYCVerified {
		return ErrKYCPending
	}
	if d.status == DriverStatusSuspended {
		return ErrDriverSuspended
	}
	if d.status == DriverStatusOnDelivery {
		return errors.New("cannot go online while on delivery")
	}
	d.status = DriverStatusOnline
	now := time.Now().UTC()
	d.lastOnlineAt = &now
	d.updatedAt = now
	return nil
}

// GoOffline sets the driver status to offline.
func (d *Driver) GoOffline() error {
	if d.status == DriverStatusOnDelivery {
		return errors.New("cannot go offline while on delivery")
	}
	d.status = DriverStatusOffline
	d.updatedAt = time.Now().UTC()
	return nil
}

// GoOnBreak sets the driver status to on_break.
func (d *Driver) GoOnBreak() error {
	if d.status == DriverStatusOnDelivery {
		return errors.New("cannot go on break while on delivery")
	}
	d.status = DriverStatusOnBreak
	d.updatedAt = time.Now().UTC()
	return nil
}

// StartDelivery sets the driver status to on_delivery.
func (d *Driver) StartDelivery() error {
	if d.status != DriverStatusOnline {
		return ErrDriverOffline
	}
	d.status = DriverStatusOnDelivery
	d.updatedAt = time.Now().UTC()
	return nil
}

// CompleteDelivery sets the driver back to online and records earnings.
func (d *Driver) CompleteDelivery(earnings float64, rating float64) error {
	if d.status != DriverStatusOnDelivery {
		return errors.New("driver is not on delivery")
	}

	d.totalEarnings += earnings
	d.totalDeliveries++
	d.UpdateRating(rating)
	d.RecalculateTier()
	d.status = DriverStatusOnline
	d.updatedAt = time.Now().UTC()
	return nil
}

// Suspend sets the driver status to suspended.
func (d *Driver) Suspend() {
	d.status = DriverStatusSuspended
	d.updatedAt = time.Now().UTC()
}

// Reactivate sets the driver back to offline (from suspended).
func (d *Driver) Reactivate() {
	if d.status == DriverStatusSuspended {
		d.status = DriverStatusOffline
		d.updatedAt = time.Now().UTC()
	}
}

// IsAvailable returns true if the driver can receive new orders.
func (d *Driver) IsAvailable() bool {
	return d.status == DriverStatusOnline && d.kycStatus == KYCVerified
}

// ============ Rating & Tier ============

// UpdateRating recalculates the rating from a new rating.
func (d *Driver) UpdateRating(newRating float64) {
	if newRating < 0 {
		newRating = 0
	} else if newRating > 5 {
		newRating = 5
	}

	total := d.rating * float64(d.ratingCount)
	d.ratingCount++
	d.rating = (total + newRating) / float64(d.ratingCount)
	d.updatedAt = time.Now().UTC()
}

// UpdateAcceptanceRate updates the acceptance rate based on accepted/total orders.
func (d *Driver) UpdateAcceptanceRate(accepted, total int) {
	if total > 0 {
		d.acceptanceRate = float64(accepted) / float64(total)
	}
	d.updatedAt = time.Now().UTC()
}

// UpdateCompletionRate updates the completion rate.
func (d *Driver) UpdateCompletionRate(completed, accepted int) {
	if accepted > 0 {
		d.completionRate = float64(completed) / float64(accepted)
	}
	d.updatedAt = time.Now().UTC()
}

// RecalculateTier checks if the driver qualifies for a higher tier.
func (d *Driver) RecalculateTier() {
	tier := TierStandard

	// Silver: 4.5+ rating, 70%+ acceptance, 90%+ completion
	if d.rating >= 4.5 && d.acceptanceRate >= 0.70 && d.completionRate >= 0.90 {
		tier = TierSilver
	}

	// Gold: 4.7+ rating, 85%+ acceptance, 95%+ completion
	if d.rating >= 4.7 && d.acceptanceRate >= 0.85 && d.completionRate >= 0.95 {
		tier = TierGold
	}

	// Platinum: 4.9+ rating, 95%+ acceptance, 98%+ completion
	if d.rating >= 4.9 && d.acceptanceRate >= 0.95 && d.completionRate >= 0.98 {
		tier = TierPlatinum
	}

	if tier != d.tier {
		d.tier = tier
		d.updatedAt = time.Now().UTC()
	}
}

// SetTrustScore updates the trust score (0-100).
func (d *Driver) SetTrustScore(score int) {
	if score < 0 {
		score = 0
	} else if score > 100 {
		score = 100
	}
	d.trustScore = score
	d.updatedAt = time.Now().UTC()
}

// ============ Helpers ============

// DistanceTo calculates distance from the driver to a point (in km).
func (d *Driver) DistanceTo(lat, lng float64) float64 {
	return haversine(d.latitude, d.longitude, lat, lng)
}

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

// ============ Earnings ============

// Earnings represents a driver's earnings summary.
type Earnings struct {
	driverID        uuid.UUID
	todayEarnings   float64
	weekEarnings    float64
	monthEarnings   float64
	todayDeliveries int
	todayHours      float64
	pendingPayout   float64
}

func NewEarnings(driverID uuid.UUID) *Earnings {
	return &Earnings{driverID: driverID}
}

func (e *Earnings) DriverID() uuid.UUID    { return e.driverID }
func (e *Earnings) TodayEarnings() float64 { return e.todayEarnings }
func (e *Earnings) WeekEarnings() float64  { return e.weekEarnings }
func (e *Earnings) MonthEarnings() float64 { return e.monthEarnings }
func (e *Earnings) TodayDeliveries() int   { return e.todayDeliveries }
func (e *Earnings) TodayHours() float64    { return e.todayHours }
func (e *Earnings) PendingPayout() float64 { return e.pendingPayout }

func (e *Earnings) AddEarnings(amount float64) {
	e.todayEarnings += amount
	e.weekEarnings += amount
	e.monthEarnings += amount
	e.pendingPayout += amount
	e.todayDeliveries++
}

func (e *Earnings) ResetDaily() {
	e.todayEarnings = 0
	e.todayDeliveries = 0
	e.todayHours = 0
}

func (e *Earnings) Payout(amount float64) error {
	if amount > e.pendingPayout {
		return errors.New("payout amount exceeds pending")
	}
	e.pendingPayout -= amount
	return nil
}

func (e *Earnings) HourlyRate() float64 {
	if e.todayHours == 0 {
		return 0
	}
	return e.todayEarnings / e.todayHours
}

// ============ Reconstructor ============

func ReconstructDriver(
	id, userID uuid.UUID, name, phone string,
	vehicleType VehicleType, vehiclePlate, licenseNumber string,
	kycStatus KYCStatus, status DriverStatus, tier DriverTier,
	rating float64, ratingCount int,
	acceptanceRate, completionRate float64,
	trustScore int, totalEarnings float64, totalDeliveries int,
	lat, lng, heading, speed float64,
	lastOnlineAt *time.Time, photoURL string,
	createdAt, updatedAt time.Time,
) *Driver {
	return &Driver{
		id: id, userID: userID, name: name, phone: phone,
		vehicleType: vehicleType, vehiclePlate: vehiclePlate, licenseNumber: licenseNumber,
		kycStatus: kycStatus, status: status, tier: tier,
		rating: rating, ratingCount: ratingCount,
		acceptanceRate: acceptanceRate, completionRate: completionRate,
		trustScore: trustScore, totalEarnings: totalEarnings, totalDeliveries: totalDeliveries,
		latitude: lat, longitude: lng, heading: heading, speed: speed,
		lastOnlineAt: lastOnlineAt, photoURL: photoURL,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}
