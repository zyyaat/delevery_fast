package main

import (
    "context"
    "fmt"
    "log"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// Config holds service configuration
type Config struct {
    HTTPPort        int           `env:"HTTP_PORT" default:"8086"`
    DatabaseURL     string        `env:"DATABASE_URL" required:"true"`
    RedisURL        string        `env:"REDIS_URL" required:"true"`
    KafkaBrokers    string        `env:"KAFKA_BROKERS" required:"true"`
    LogLevel        string        `env:"LOG_LEVEL" default:"info"`
    ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func main() {
    cfg := &Config{
        HTTPPort:        8086,
        DatabaseURL:     getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/delivery-matching_db?sslmode=disable"),
        RedisURL:        getEnvOrDefault("REDIS_URL", "localhost:6379"),
        KafkaBrokers:    getEnvOrDefault("KAFKA_BROKERS", "localhost:9092"),
        LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
        ShutdownTimeout: 30 * time.Second,
    }

    // Setup structured logger
    setupLogger(cfg.LogLevel)
    slog.Info("starting Delivery Matching Service",
        "port", cfg.HTTPPort,
        "service", "delivery-matching",
    )

    // Create HTTP server
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)
    mux.HandleFunc("/ready", readyHandler)

    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.HTTPPort),
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // Start server in goroutine
    go func() {
        slog.Info("HTTP server listening", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server failed", "error", err)
            os.Exit(1)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    slog.Info("shutting down delivery-matching service...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("server forced to shutdown", "error", err)
        os.Exit(1)
    }

    slog.Info("delivery-matching service stopped")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok","service":"delivery-matching","version":"1.0.0"}`))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Check dependencies (DB, Redis, Kafka)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ready"}`))
}

func setupLogger(level string) {
    var lvl slog.Level
    switch level {
    case "debug":
        lvl = slog.LevelDebug
    case "info":
        lvl = slog.LevelInfo
    case "warn":
        lvl = slog.LevelWarn
    case "error":
        lvl = slog.LevelError
    default:
        lvl = slog.LevelInfo
    }
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
    slog.SetDefault(slog.New(handler))
}

func getEnvOrDefault(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}
