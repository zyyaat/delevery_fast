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

	"github.com/food-platform/menu/internal/application"
	"github.com/food-platform/menu/internal/infrastructure/postgres"
	httpinterfaces "github.com/food-platform/menu/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"

	_ "github.com/lib/pq"
)

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_menu_service", "port", cfg.HTTPPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := connectDB(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed_to_connect_db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	itemRepo := postgres.NewMenuItemRepository(db)
	catRepo := postgres.NewCategoryRepository(db)

	createItemUC := application.NewCreateMenuItemUseCase(itemRepo)
	getMenuUC := application.NewGetMenuUseCase(itemRepo, catRepo)
	toggleUC := application.NewToggleAvailabilityUseCase(itemRepo)
	createCatUC := application.NewCreateCategoryUseCase(catRepo)

	handler := httpinterfaces.SetupRouter(createItemUC, getMenuUC, toggleUC, createCatUC)
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
	slog.Info("shutting_down_menu_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("menu_service_stopped")
}

type Config struct {
	HTTPPort    int
	DatabaseURL string
	LogLevel    string
}

func loadConfig() Config {
	return Config{
		HTTPPort:    config.GetEnvInt("HTTP_PORT", 8083),
		DatabaseURL: config.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/menus_db?sslmode=disable"),
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
