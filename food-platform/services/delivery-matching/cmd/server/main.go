package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/food-platform/delivery-matching/internal/application"
	httpinterfaces "github.com/food-platform/delivery-matching/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"
	"github.com/google/uuid"
)

// Mock implementations for MVP (will be replaced with real gRPC clients)
type mockGeoService struct{}
type mockDriverService struct{}
type mockOrderService struct{}
type mockPublisher struct{}

func (m *mockGeoService) FindNearbyDrivers(ctx context.Context, lat, lng, radiusKm float64, count int) ([]application.DriverLocation, error) {
	return nil, nil
}
func (m *mockDriverService) GetDriverInfo(ctx context.Context, driverID uuid.UUID) (*application.DriverInfo, error) {
	return &application.DriverInfo{Rating: 4.5, AcceptanceRate: 0.85, CompletionRate: 0.95, VehicleType: "motorcycle"}, nil
}
func (m *mockOrderService) UpdateOrderDriver(ctx context.Context, orderID, driverID uuid.UUID) error { return nil }
func (m *mockPublisher) PublishOrderAssigned(ctx context.Context, orderID, driverID uuid.UUID) error { return nil }
func (m *mockPublisher) PublishDispatchFailed(ctx context.Context, orderID uuid.UUID, reason string) error { return nil }

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_delivery_matching_service", "port", cfg.HTTPPort)

	matchUC := application.NewMatchOrderUseCase(
		&mockGeoService{},
		&mockDriverService{},
		&mockOrderService{},
		&mockPublisher{},
	)
	acceptUC := application.NewAcceptOrderUseCase(&mockOrderService{}, &mockPublisher{})

	handler := httpinterfaces.SetupRouter(matchUC, acceptUC)
	srv := server.New(handler, server.DefaultConfig(cfg.HTTPPort))

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server_failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting_down_delivery_matching_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("delivery_matching_service_stopped")
}

type Config struct {
	HTTPPort int
	LogLevel string
}

func loadConfig() Config {
	return Config{
		HTTPPort: config.GetEnvInt("HTTP_PORT", 8086),
		LogLevel: config.GetEnv("LOG_LEVEL", "info"),
	}
}
