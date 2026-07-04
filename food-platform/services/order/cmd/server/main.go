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

	"github.com/food-platform/order/internal/application"
	"github.com/food-platform/order/internal/infrastructure/kafka"
	"github.com/food-platform/order/internal/infrastructure/postgres"
	httpinterfaces "github.com/food-platform/order/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"

	_ "github.com/lib/pq"
)

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_order_service", "port", cfg.HTTPPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed_to_connect_db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	orderRepo := postgres.NewOrderRepository(db)

	// Initialize event publisher (mock for now; switch to real Kafka in production)
	var publisher application.EventPublisher
	if cfg.KafkaBrokers != "" {
		// TODO: Initialize real Kafka producer
		publisher = kafka.NewMockPublisher()
	} else {
		publisher = kafka.NewMockPublisher()
	}

	// Initialize use cases
	createUC := application.NewCreateOrderUseCase(orderRepo, publisher)
	getUC := application.NewGetOrderUseCase(orderRepo)
	getActiveUC := application.NewGetActiveOrdersUseCase(orderRepo)
	getHistoryUC := application.NewGetOrderHistoryUseCase(orderRepo)
	cancelUC := application.NewCancelOrderUseCase(orderRepo, publisher)
	updateStatusUC := application.NewUpdateOrderStatusUseCase(orderRepo, publisher)

	// Setup HTTP router
	handler := httpinterfaces.SetupRouter(
		createUC, getUC, getActiveUC, getHistoryUC, cancelUC, updateStatusUC,
	)

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
	slog.Info("shutting_down_order_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("order_service_stopped")
}

type Config struct {
	HTTPPort     int
	DatabaseURL  string
	KafkaBrokers string
	LogLevel     string
}

func loadConfig() Config {
	return Config{
		HTTPPort:     config.GetEnvInt("HTTP_PORT", 8084),
		DatabaseURL:  config.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/orders_db?sslmode=disable"),
		KafkaBrokers: config.GetEnv("KAFKA_BROKERS", ""),
		LogLevel:     config.GetEnv("LOG_LEVEL", "info"),
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
