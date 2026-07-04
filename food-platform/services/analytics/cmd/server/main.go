package main

import (
        "context"
        "log/slog"
        "os"
        "os/signal"
        "syscall"
        "time"

        "github.com/food-platform/analytics/internal/application"
        httpinterfaces "github.com/food-platform/analytics/internal/interfaces/http"
        "github.com/food-platform/shared/config"
        "github.com/food-platform/shared/logging"
        "github.com/food-platform/shared/server"
)

func main() {
        cfg := loadConfig()
        logging.Setup(cfg.LogLevel)
        slog.Info("starting_analytics_service", "port", cfg.HTTPPort)

        dashboardUC := application.NewGetDashboardStatsUseCase()
        zonesUC := application.NewGetZoneMetricsUseCase()
        incidentsUC := application.NewGetIncidentsUseCase()
        forecastUC := application.NewGetForecastUseCase()

        handler := httpinterfaces.SetupRouter(dashboardUC, zonesUC, incidentsUC, forecastUC)
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
        slog.Info("shutting_down_analytics_service")

        shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer shutdownCancel()

        if err := srv.Stop(shutdownCtx); err != nil {
                slog.Error("server_forced_to_shutdown", "error", err)
                os.Exit(1)
        }

        slog.Info("analytics_service_stopped")
}

type Config struct {
        HTTPPort int
        LogLevel string
}

func loadConfig() Config {
        return Config{
                HTTPPort: config.GetEnvInt("HTTP_PORT", 8092),
                LogLevel: config.GetEnv("LOG_LEVEL", "info"),
        }
}
