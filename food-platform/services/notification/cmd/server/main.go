package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/food-platform/services/notification/internal/application"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"
	"net/http"

	"github.com/food-platform/shared/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_notification_service", "port", cfg.HTTPPort)

	// Initialize use cases (with mock implementations for now)
	// In production, wire up real repositories, push senders, SMS, WebSocket
	var sendNotifUC *application.SendNotificationUseCase // = application.NewSendNotificationUseCase(...)
	var getNotifsUC *application.GetNotificationsUseCase // = application.NewGetNotificationsUseCase(...)

	// Setup HTTP router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logging)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS([]string{"*"}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"notification","version":"1.0.0"}`))
	})

	// TODO: Add notification routes when repository is implemented
	_ = sendNotifUC
	_ = getNotifsUC

	srv := server.New(r, server.DefaultConfig(cfg.HTTPPort))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// TODO: Start Kafka consumers for order + payment events
		slog.Info("kafka_consumers_starting")
		_ = ctx // Will be used for graceful shutdown of consumers
	}()

	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server_failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting_down_notification_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("notification_service_stopped")
}

type Config struct {
	HTTPPort     int
	DatabaseURL  string
	KafkaBrokers string
	LogLevel     string
}

func loadConfig() Config {
	return Config{
		HTTPPort:     config.GetEnvInt("HTTP_PORT", 8089),
		DatabaseURL:  config.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/notifications_db?sslmode=disable"),
		KafkaBrokers: config.GetEnv("KAFKA_BROKERS", "localhost:9092"),
		LogLevel:     config.GetEnv("LOG_LEVEL", "info"),
	}
}
