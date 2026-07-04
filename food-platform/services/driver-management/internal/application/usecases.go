// Package application contains use cases for the Driver Management Service.
package application

import (
	"context"
	"fmt"

	"github.com/food-platform/driver-management/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

type DriverRepository interface {
	Create(ctx context.Context, d *domain.Driver) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Driver, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) (*domain.Driver, error)
	Update(ctx context.Context, d *domain.Driver) error
	FindAvailable(ctx context.Context, lat, lng float64, radiusKm float64, limit int) ([]*domain.Driver, error)
	UpdateLocation(ctx context.Context, id uuid.UUID, lat, lng, heading, speed float64) error
}

// ============ DTOs ============

type DriverDTO struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Phone           string    `json:"phone"`
	VehicleType     string    `json:"vehicle_type"`
	VehiclePlate    string    `json:"vehicle_plate,omitempty"`
	LicenseNumber   string    `json:"license_number,omitempty"`
	KYCStatus       string    `json:"kyc_status"`
	Status          string    `json:"status"`
	Tier            string    `json:"tier"`
	Rating          float64   `json:"rating"`
	RatingCount     int       `json:"rating_count"`
	AcceptanceRate  float64   `json:"acceptance_rate"`
	CompletionRate  float64   `json:"completion_rate"`
	TrustScore      int       `json:"trust_score"`
	TotalEarnings   float64   `json:"total_earnings"`
	TotalDeliveries int       `json:"total_deliveries"`
	PhotoURL        string    `json:"photo_url,omitempty"`
}

type EarningsDTO struct {
	TodayEarnings   float64 `json:"today_earnings"`
	WeekEarnings    float64 `json:"week_earnings"`
	TodayDeliveries int     `json:"today_deliveries"`
	HourlyRate      float64 `json:"hourly_rate"`
	PendingPayout   float64 `json:"pending_payout"`
}

// ============ Commands ============

type RegisterDriverCommand struct {
	UserID      uuid.UUID
	Name        string
	Phone       string
	VehicleType string
}

type UpdateStatusCommand struct {
	DriverID uuid.UUID
	Status   string
}

type UpdateLocationCommand struct {
	DriverID uuid.UUID
	Lat      float64
	Lng      float64
	Heading  float64
	Speed    float64
}

type CompleteDeliveryCommand struct {
	DriverID uuid.UUID
	Earnings float64
	Rating   float64
}

// ============ Use Cases ============

type RegisterDriverUseCase struct {
	repo DriverRepository
}

func NewRegisterDriverUseCase(repo DriverRepository) *RegisterDriverUseCase {
	return &RegisterDriverUseCase{repo: repo}
}

func (uc *RegisterDriverUseCase) Execute(ctx context.Context, cmd RegisterDriverCommand) (*DriverDTO, error) {
	driver, err := domain.NewDriver(cmd.UserID, cmd.Name, cmd.Phone, domain.VehicleType(cmd.VehicleType))
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, driver); err != nil {
		return nil, fmt.Errorf("create driver: %w", err)
	}

	return toDriverDTO(driver), nil
}

type GetDriverUseCase struct {
	repo DriverRepository
}

func NewGetDriverUseCase(repo DriverRepository) *GetDriverUseCase {
	return &GetDriverUseCase{repo: repo}
}

func (uc *GetDriverUseCase) Execute(ctx context.Context, id uuid.UUID) (*DriverDTO, error) {
	driver, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDriverDTO(driver), nil
}

type GetDriverByUserUseCase struct {
	repo DriverRepository
}

func NewGetDriverByUserUseCase(repo DriverRepository) *GetDriverByUserUseCase {
	return &GetDriverByUserUseCase{repo: repo}
}

func (uc *GetDriverByUserUseCase) Execute(ctx context.Context, userID uuid.UUID) (*DriverDTO, error) {
	driver, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toDriverDTO(driver), nil
}

type UpdateStatusUseCase struct {
	repo DriverRepository
}

func NewUpdateStatusUseCase(repo DriverRepository) *UpdateStatusUseCase {
	return &UpdateStatusUseCase{repo: repo}
}

func (uc *UpdateStatusUseCase) Execute(ctx context.Context, cmd UpdateStatusCommand) error {
	driver, err := uc.repo.FindByID(ctx, cmd.DriverID)
	if err != nil {
		return err
	}

	switch domain.DriverStatus(cmd.Status) {
	case domain.DriverStatusOnline:
		if err := driver.GoOnline(); err != nil {
			return err
		}
	case domain.DriverStatusOffline:
		if err := driver.GoOffline(); err != nil {
			return err
		}
	case domain.DriverStatusOnBreak:
		if err := driver.GoOnBreak(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid status: %s", cmd.Status)
	}

	return uc.repo.Update(ctx, driver)
}

type UpdateLocationUseCase struct {
	repo DriverRepository
}

func NewUpdateLocationUseCase(repo DriverRepository) *UpdateLocationUseCase {
	return &UpdateLocationUseCase{repo: repo}
}

func (uc *UpdateLocationUseCase) Execute(ctx context.Context, cmd UpdateLocationCommand) error {
	driver, err := uc.repo.FindByID(ctx, cmd.DriverID)
	if err != nil {
		return err
	}

	if err := driver.SetLocation(cmd.Lat, cmd.Lng, cmd.Heading, cmd.Speed); err != nil {
		return err
	}

	// For high-frequency updates, use the optimized location update method
	return uc.repo.UpdateLocation(ctx, cmd.DriverID, cmd.Lat, cmd.Lng, cmd.Heading, cmd.Speed)
}

type CompleteDeliveryUseCase struct {
	repo DriverRepository
}

func NewCompleteDeliveryUseCase(repo DriverRepository) *CompleteDeliveryUseCase {
	return &CompleteDeliveryUseCase{repo: repo}
}

func (uc *CompleteDeliveryUseCase) Execute(ctx context.Context, cmd CompleteDeliveryCommand) error {
	driver, err := uc.repo.FindByID(ctx, cmd.DriverID)
	if err != nil {
		return err
	}

	if err := driver.CompleteDelivery(cmd.Earnings, cmd.Rating); err != nil {
		return err
	}

	return uc.repo.Update(ctx, driver)
}

type FindAvailableDriversUseCase struct {
	repo DriverRepository
}

func NewFindAvailableDriversUseCase(repo DriverRepository) *FindAvailableDriversUseCase {
	return &FindAvailableDriversUseCase{repo: repo}
}

func (uc *FindAvailableDriversUseCase) Execute(ctx context.Context, lat, lng, radiusKm float64, limit int) ([]*DriverDTO, error) {
	if limit == 0 {
		limit = 40
	}

	drivers, err := uc.repo.FindAvailable(ctx, lat, lng, radiusKm, limit)
	if err != nil {
		return nil, err
	}

	dtos := make([]*DriverDTO, 0, len(drivers))
	for _, d := range drivers {
		dtos = append(dtos, toDriverDTO(d))
	}
	return dtos, nil
}

// ============ Helpers ============

func toDriverDTO(d *domain.Driver) *DriverDTO {
	return &DriverDTO{
		ID:              d.ID(),
		Name:            d.Name(),
		Phone:           d.Phone(),
		VehicleType:     string(d.VehicleType()),
		VehiclePlate:    d.VehiclePlate(),
		LicenseNumber:   d.LicenseNumber(),
		KYCStatus:       string(d.KYCStatus()),
		Status:          string(d.Status()),
		Tier:            string(d.Tier()),
		Rating:          d.Rating(),
		RatingCount:     d.RatingCount(),
		AcceptanceRate:  d.AcceptanceRate(),
		CompletionRate:  d.CompletionRate(),
		TrustScore:      d.TrustScore(),
		TotalEarnings:   d.TotalEarnings(),
		TotalDeliveries: d.TotalDeliveries(),
		PhotoURL:        d.PhotoURL(),
	}
}
