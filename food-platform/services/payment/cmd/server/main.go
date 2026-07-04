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

	"github.com/food-platform/services/payment/internal/application"
	"github.com/food-platform/services/payment/internal/infrastructure/kafka"
	"github.com/food-platform/services/payment/internal/infrastructure/postgres"
	"github.com/food-platform/services/payment/internal/infrastructure/providers"
	httpinterfaces "github.com/food-platform/services/payment/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"

	_ "github.com/lib/pq"
)

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_payment_service", "port", cfg.HTTPPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed_to_connect_db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	paymentRepo := postgres.NewPaymentRepository(db)

	// Initialize providers
	vodafoneProvider := providers.NewVodafoneCashProvider(cfg.VodafoneMerchantID, cfg.VodafoneAPIKey, cfg.SandboxMode)
	instaPayProvider := providers.NewInstaPayProvider(cfg.InstaPayAPIKey, cfg.SandboxMode)
	paymobProvider := providers.NewPaymobProvider(cfg.PaymobAPIKey, cfg.PaymobMerchantID, cfg.SandboxMode)

	providerFactory := providers.NewFactory(vodafoneProvider, instaPayProvider, paymobProvider)

	// Initialize event publisher
	var publisher application.EventPublisher
	publisher = kafka.NewMockPublisher()

	// Initialize use cases
	chargeUC := application.NewChargePaymentUseCase(paymentRepo, providerFactory, publisher)
	getUC := application.NewGetPaymentUseCase(paymentRepo)
	getByOrderUC := application.NewGetPaymentByOrderUseCase(paymentRepo)
	refundUC := application.NewRefundPaymentUseCase(paymentRepo, providerFactory, publisher)

	// Setup HTTP router
	handler := httpinterfaces.SetupRouter(chargeUC, getUC, getByOrderUC, refundUC)
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
	slog.Info("shutting_down_payment_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("payment_service_stopped")
}

type Config struct {
	HTTPPort          int
	DatabaseURL       string
	LogLevel          string
	SandboxMode       bool
	VodafoneMerchantID string
	VodafoneAPIKey    string
	InstaPayAPIKey    string
	PaymobAPIKey      string
	PaymobMerchantID  string
}

func loadConfig() Config {
	return Config{
		HTTPPort:           config.GetEnvInt("HTTP_PORT", 8085),
		DatabaseURL:        config.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/payments_db?sslmode=disable"),
		LogLevel:           config.GetEnv("LOG_LEVEL", "info"),
		SandboxMode:        config.GetEnvBool("SANDBOX_MODE", true),
		VodafoneMerchantID: config.GetEnv("VODAFONE_MERCHANT_ID", ""),
		VodafoneAPIKey:     config.GetEnv("VODAFONE_API_KEY", ""),
		InstaPayAPIKey:     config.GetEnv("INSTAPAY_API_KEY", ""),
		PaymobAPIKey:       config.GetEnv("PAYMOB_API_KEY", ""),
		PaymobMerchantID:   config.GetEnv("PAYMOB_MERCHANT_ID", ""),
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
