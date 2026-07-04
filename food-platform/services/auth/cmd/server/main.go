// Package main is the entry point for the Auth Service.
// It wires together the domain, application, infrastructure, and interface layers.
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

	"github.com/food-platform/auth/internal/application"
	"github.com/food-platform/auth/internal/domain"
	"github.com/food-platform/auth/internal/infrastructure/jwt"
	"github.com/food-platform/auth/internal/infrastructure/postgres"
	"github.com/food-platform/auth/internal/infrastructure/sms"
	httpinterfaces "github.com/food-platform/auth/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := loadConfig()

	// Setup structured logging
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_auth_service",
		"port", cfg.HTTPPort,
		"log_level", cfg.LogLevel,
	)

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to PostgreSQL
	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed_to_connect_db", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("database_connected")

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	otpRepo := postgres.NewOTPRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)
	refreshRepo := postgres.NewRefreshTokenRepository(db)

	// Initialize SMS sender
	var smsSender sms.Sender
	if cfg.SMSType == "twilio" {
		smsSender = sms.NewTwilioSender(sms.TwilioConfig{
			AccountSID: cfg.TwilioAccountSID,
			AuthToken:  cfg.TwilioAuthToken,
			FromNumber: cfg.TwilioFromNumber,
		})
	} else {
		smsSender = sms.NewMockSender()
	}

	// Initialize JWT generator
	jwtGenerator := jwt.NewGenerator(jwt.Config{
		SecretKey: cfg.JWTSecret,
		Issuer:    "food-platform-auth",
		Audience:  "food-platform",
		AccessTTL: 15 * time.Minute,
	})

	// Initialize use cases
	sendOTPUseCase := application.NewSendOTPUseCase(userRepo, otpRepo, smsSender)
	verifyOTPUseCase := application.NewVerifyOTPUseCase(userRepo, otpRepo, sessionRepo, refreshRepo, jwtGenerator)
	refreshTokenUseCase := application.NewRefreshTokenUseCase(userRepo, sessionRepo, refreshRepo, jwtGenerator)
	logoutUseCase := application.NewLogoutUseCase(sessionRepo, refreshRepo)

	// Setup HTTP router
	handler := httpinterfaces.SetupRouter(
		sendOTPUseCase,
		verifyOTPUseCase,
		refreshTokenUseCase,
		logoutUseCase,
	)

	// Create and start HTTP server
	srv := server.New(handler, server.DefaultConfig(cfg.HTTPPort))

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			slog.Error("server_failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting_down_auth_service")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("auth_service_stopped")
}

// ============ Configuration ============

type Config struct {
	HTTPPort          int
	DatabaseURL       string
	LogLevel          string
	JWTSecret         string
	SMSType           string // "mock" or "twilio"
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioFromNumber  string
}

func loadConfig() Config {
	return Config{
		HTTPPort:         config.GetEnvInt("HTTP_PORT", 8081),
		DatabaseURL:      config.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable"),
		LogLevel:         config.GetEnv("LOG_LEVEL", "info"),
		JWTSecret:        config.GetEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
		SMSType:          config.GetEnv("SMS_TYPE", "mock"),
		TwilioAccountSID: config.GetEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  config.GetEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromNumber: config.GetEnv("TWILIO_FROM_NUMBER", ""),
	}
}

// connectDB connects to PostgreSQL and verifies the connection.
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

// Suppress unused import warnings (will be used in tests)
var _ = domain.ErrUserNotFound
