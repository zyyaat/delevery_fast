package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/food-platform/geo/internal/application"
	redisinfra "github.com/food-platform/geo/internal/infrastructure/redis"
	httpinterfaces "github.com/food-platform/geo/internal/interfaces/http"
	"github.com/food-platform/shared/config"
	"github.com/food-platform/shared/logging"
	"github.com/food-platform/shared/server"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := loadConfig()
	logging.Setup(cfg.LogLevel)
	slog.Info("starting_geo_service", "port", cfg.HTTPPort)

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("failed_to_connect_redis", "error", err)
		os.Exit(1)
	}
	slog.Info("redis_connected")

	locationStore := redisinfra.NewLocationStore(rdb)

	updateLocUC := application.NewUpdateLocationUseCase(locationStore)
	getLocUC := application.NewGetLocationUseCase(locationStore)
	findNearbyUC := application.NewFindNearbyUseCase(locationStore)
	calcETAUC := application.NewCalculateETAUseCase()

	handler := httpinterfaces.SetupRouter(updateLocUC, getLocUC, findNearbyUC, calcETAUC)
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
	slog.Info("shutting_down_geo_service")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("geo_service_stopped")
}

type Config struct {
	HTTPPort  int
	RedisURL  string
	LogLevel  string
}

func loadConfig() Config {
	return Config{
		HTTPPort: config.GetEnvInt("HTTP_PORT", 8088),
		RedisURL: config.GetEnv("REDIS_URL", "localhost:6379"),
		LogLevel: config.GetEnv("LOG_LEVEL", "info"),
	}
}
