// Package postgres implements the Driver repository.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/driver-management/internal/domain"
	"github.com/google/uuid"
)

type DriverRepository struct {
	db *sql.DB
}

func NewDriverRepository(db *sql.DB) *DriverRepository {
	return &DriverRepository{db: db}
}

func (r *DriverRepository) Create(ctx context.Context, d *domain.Driver) error {
	query := `
		INSERT INTO drivers (
			id, user_id, name, phone, vehicle_type, vehicle_plate, license_number,
			kyc_status, status, tier, rating, rating_count,
			acceptance_rate, completion_rate, trust_score,
			total_earnings, total_deliveries, photo_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`
	_, err := r.db.ExecContext(ctx, query,
		d.ID(), d.UserID(), d.Name(), d.Phone(),
		string(d.VehicleType()), nullableString(d.VehiclePlate()), nullableString(d.LicenseNumber()),
		string(d.KYCStatus()), string(d.Status()), string(d.Tier()),
		d.Rating(), d.RatingCount(),
		d.AcceptanceRate(), d.CompletionRate(), d.TrustScore(),
		d.TotalEarnings(), d.TotalDeliveries(),
		nullableString(d.PhotoURL()),
		d.CreatedAt(), d.UpdatedAt(),
	)
	return err
}

func (r *DriverRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error) {
	query := `
		SELECT id, user_id, name, phone, vehicle_type,
		       COALESCE(vehicle_plate, ''), COALESCE(license_number, ''),
		       kyc_status, status, tier, rating, rating_count,
		       acceptance_rate, completion_rate, trust_score,
		       total_earnings, total_deliveries,
		       COALESCE(latitude, 0), COALESCE(longitude, 0),
		       COALESCE(heading, 0), COALESCE(speed, 0),
		       last_online_at, COALESCE(photo_url, ''),
		       created_at, updated_at
		FROM drivers WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanDriver(row)
}

func (r *DriverRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error) {
	query := `
		SELECT id, user_id, name, phone, vehicle_type,
		       COALESCE(vehicle_plate, ''), COALESCE(license_number, ''),
		       kyc_status, status, tier, rating, rating_count,
		       acceptance_rate, completion_rate, trust_score,
		       total_earnings, total_deliveries,
		       COALESCE(latitude, 0), COALESCE(longitude, 0),
		       COALESCE(heading, 0), COALESCE(speed, 0),
		       last_online_at, COALESCE(photo_url, ''),
		       created_at, updated_at
		FROM drivers WHERE user_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	return scanDriver(row)
}

func (r *DriverRepository) Update(ctx context.Context, d *domain.Driver) error {
	query := `
		UPDATE drivers SET
			name = $2, phone = $3, vehicle_type = $4, vehicle_plate = $5,
			license_number = $6, kyc_status = $7, status = $8, tier = $9,
			rating = $10, rating_count = $11,
			acceptance_rate = $12, completion_rate = $13, trust_score = $14,
			total_earnings = $15, total_deliveries = $16,
			last_online_at = $17, photo_url = $18, updated_at = $19
		WHERE id = $1
	`
	var lastOnlineAt interface{}
	if d.LastOnlineAt() != nil {
		lastOnlineAt = *d.LastOnlineAt()
	}

	_, err := r.db.ExecContext(ctx, query,
		d.ID(), d.Name(), d.Phone(),
		string(d.VehicleType()), nullableString(d.VehiclePlate()),
		nullableString(d.LicenseNumber()),
		string(d.KYCStatus()), string(d.Status()), string(d.Tier()),
		d.Rating(), d.RatingCount(),
		d.AcceptanceRate(), d.CompletionRate(), d.TrustScore(),
		d.TotalEarnings(), d.TotalDeliveries(),
		lastOnlineAt, nullableString(d.PhotoURL()),
		time.Now().UTC(),
	)
	return err
}

func (r *DriverRepository) FindAvailable(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*domain.Driver, error) {
	query := `
		SELECT id, user_id, name, phone, vehicle_type,
		       COALESCE(vehicle_plate, ''), COALESCE(license_number, ''),
		       kyc_status, status, tier, rating, rating_count,
		       acceptance_rate, completion_rate, trust_score,
		       total_earnings, total_deliveries,
		       COALESCE(latitude, 0), COALESCE(longitude, 0),
		       COALESCE(heading, 0), COALESCE(speed, 0),
		       last_online_at, COALESCE(photo_url, ''),
		       created_at, updated_at
		FROM drivers
		WHERE status = 'online' AND kyc_status = 'verified'
		  AND latitude IS NOT NULL AND longitude IS NOT NULL
		  AND ST_DWithin(
		    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography,
		    ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
		    $3 * 1000
		  )
		ORDER BY ST_Distance(
		    ST_SetSRID(ST_MakePoint(longitude, latitude), 4326)::geography,
		    ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
		)
		LIMIT $4
	`
	rows, err := r.db.QueryContext(ctx, query, lng, lat, radiusKm, limit)
	if err != nil {
		return nil, fmt.Errorf("find available drivers: %w", err)
	}
	defer rows.Close()

	var drivers []*domain.Driver
	for rows.Next() {
		d, err := scanDriver(rows)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, rows.Err()
}

func (r *DriverRepository) UpdateLocation(ctx context.Context, id uuid.UUID, lat, lng, heading, speed float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE drivers SET latitude = $2, longitude = $3, heading = $4, speed = $5, updated_at = $6 WHERE id = $1`,
		id, lat, lng, heading, speed, time.Now().UTC(),
	)
	return err
}

// ============ Helpers ============

func nullableString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanDriver(s scanner) (*domain.Driver, error) {
	var (
		id              uuid.UUID
		userID          uuid.UUID
		name            string
		phone           string
		vehicleType     string
		vehiclePlate    string
		licenseNumber   string
		kycStatus       string
		status          string
		tier            string
		rating          float64
		ratingCount     int
		acceptanceRate  float64
		completionRate  float64
		trustScore      int
		totalEarnings   float64
		totalDeliveries int
		lat             float64
		lng             float64
		heading         float64
		speed           float64
		lastOnlineAt    *time.Time
		photoURL        string
		createdAt       time.Time
		updatedAt       time.Time
	)

	err := s.Scan(
		&id, &userID, &name, &phone, &vehicleType,
		&vehiclePlate, &licenseNumber,
		&kycStatus, &status, &tier, &rating, &ratingCount,
		&acceptanceRate, &completionRate, &trustScore,
		&totalEarnings, &totalDeliveries,
		&lat, &lng, &heading, &speed,
		&lastOnlineAt, &photoURL,
		&createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrDriverNotFound
		}
		return nil, fmt.Errorf("scanDriver: %w", err)
	}

	return domain.ReconstructDriver(
		id, userID, name, phone,
		domain.VehicleType(vehicleType), vehiclePlate, licenseNumber,
		domain.KYCStatus(kycStatus), domain.DriverStatus(status), domain.DriverTier(tier),
		rating, ratingCount, acceptanceRate, completionRate,
		trustScore, totalEarnings, totalDeliveries,
		lat, lng, heading, speed,
		lastOnlineAt, photoURL,
		createdAt, updatedAt,
	), nil
}
