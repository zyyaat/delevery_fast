package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/food-platform/driver-management/internal/application"
	"github.com/food-platform/driver-management/internal/infrastructure/postgres"
	httpinterfaces "github.com/food-platform/driver-management/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"

	_ "github.com/lib/pq"
)

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_driver_management_service", "port", cfg.HTTPPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed_to_connect_db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	driverRepo := postgres.NewDriverRepository(db)

	registerUC := application.NewRegisterDriverUseCase(driverRepo)
	getUC := application.NewGetDriverUseCase(driverRepo)
	updateStatusUC := application.NewUpdateStatusUseCase(driverRepo)
	updateLocationUC := application.NewUpdateLocationUseCase(driverRepo)

	handler := httpinterfaces.SetupRouter(registerUC, getUC, updateStatusUC, updateLocationUC)
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
	slog.Info("shutting_down_driver_management_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("driver_management_service_stopped")
}

type Config struct {
	HTTPPort    int
	DatabaseURL string
	LogLevel    string
}

func loadConfig() Config {
	return Config{
		HTTPPort:    config.GetEnvInt("HTTP_PORT", 8087),
		DatabaseURL: config.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/drivers_db?sslmode=disable"),
		LogLevel:    config.GetEnv("LOG_LEVEL", "info"),
	}
}

func connectDB(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return db, nil
}
