#!/bin/bash
# Script to create 12 Go service skeletons

set -e

SERVICES_ROOT="/home/z/my-project/food-platform/services"

declare -A SERVICE_PORTS=(
  ["auth"]="8081"
  ["restaurant-catalog"]="8082"
  ["menu"]="8083"
  ["order"]="8084"
  ["payment"]="8085"
  ["delivery-matching"]="8086"
  ["driver-management"]="8087"
  ["geo"]="8088"
  ["notification"]="8089"
  ["fraud"]="8090"
  ["promo"]="8091"
  ["analytics"]="8092"
)

declare -A SERVICE_NAMES=(
  ["auth"]="Authentication Service"
  ["restaurant-catalog"]="Restaurant Catalog Service"
  ["menu"]="Menu Service"
  ["order"]="Order Service"
  ["payment"]="Payment Service"
  ["delivery-matching"]="Delivery Matching Service"
  ["driver-management"]="Driver Management Service"
  ["geo"]="Geo/Tracking Service"
  ["notification"]="Notification Service"
  ["fraud"]="Fraud Detection Service"
  ["promo"]="Promo & Loyalty Service"
  ["analytics"]="Analytics Service"
)

for svc in auth restaurant-catalog menu order payment delivery-matching driver-management geo notification fraud promo analytics; do
  NAME="${SERVICE_NAMES[$svc]}"
  PORT="${SERVICE_PORTS[$svc]}"
  echo "=== Creating $svc (port $PORT) ==="

  SVC_DIR="$SERVICES_ROOT/$svc"
  mkdir -p "$SVC_DIR"/{cmd/server,internal/{domain,application,infrastructure/{postgres,kafka,grpc},interfaces/{http/{handlers,middleware},grpc}},migrations,proto,openapi,deployments}

  # go.mod
  cat > "$SVC_DIR/go.mod" << EOF
module github.com/food-platform/$svc

go 1.22

require (
    github.com/go-chi/chi/v5 v5.0.12
    github.com/google/uuid v1.6.0
    github.com/lib/pq v1.10.9
    github.com/confluentinc/confluent-kafka-go v1.4.2
    google.golang.org/grpc v1.64.0
    google.golang.org/protobuf v1.34.1
)

require (
    github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
    golang.org/x/net v0.25.0 // indirect
    golang.org/x/sys v0.20.0 // indirect
    golang.org/x/text v0.15.0 // indirect
    google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-9c9b29c4b60e // indirect
)
EOF

  # main.go
  cat > "$SVC_DIR/cmd/server/main.go" << EOF
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
    HTTPPort        int           \`env:"HTTP_PORT" default:"$PORT"\`
    DatabaseURL     string        \`env:"DATABASE_URL" required:"true"\`
    RedisURL        string        \`env:"REDIS_URL" required:"true"\`
    KafkaBrokers    string        \`env:"KAFKA_BROKERS" required:"true"\`
    LogLevel        string        \`env:"LOG_LEVEL" default:"info"\`
    ShutdownTimeout time.Duration \`env:"SHUTDOWN_TIMEOUT" default:"30s"\`
}

func main() {
    cfg := &Config{
        HTTPPort:        $PORT,
        DatabaseURL:     getEnvOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/${svc}_db?sslmode=disable"),
        RedisURL:        getEnvOrDefault("REDIS_URL", "localhost:6379"),
        KafkaBrokers:    getEnvOrDefault("KAFKA_BROKERS", "localhost:9092"),
        LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
        ShutdownTimeout: 30 * time.Second,
    }

    // Setup structured logger
    setupLogger(cfg.LogLevel)
    slog.Info("starting $NAME",
        "port", cfg.HTTPPort,
        "service", "$svc",
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
    slog.Info("shutting down $svc service...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("server forced to shutdown", "error", err)
        os.Exit(1)
    }

    slog.Info("$svc service stopped")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(\`{"status":"ok","service":"$svc","version":"1.0.0"}\`))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
    // TODO: Check dependencies (DB, Redis, Kafka)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(\`{"status":"ready"}\`))
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
EOF

  # Makefile
  cat > "$SVC_DIR/Makefile" << 'EOF'
.PHONY: build run test lint migrate migrate-down create-migration docker-build clean

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

migrate:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

create-migration:
	@migrate create -ext sql -dir migrations -seq $(NAME)

docker-build:
	docker build -t food-platform/SERVICE_NAME:latest -f deployments/Dockerfile .

clean:
	rm -rf bin/
EOF
  sed -i "s/SERVICE_NAME/$svc/g" "$SVC_DIR/Makefile"

  # Dockerfile
  cat > "$SVC_DIR/deployments/Dockerfile" << 'EOF'
# Multi-stage build
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

# Distroless final image
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /server /server
COPY --from=builder /app/migrations /migrations

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/server"]
EOF

  # README.md
  cat > "$SVC_DIR/README.md" << EOF
# $NAME

Go microservice for the Food Delivery Platform.

## Quick Start

\`\`\`bash
# Run locally
make run

# Build
make build

# Test
make test

# Run migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/${svc}_db?sslmode=disable make migrate
\`\`\`

## Endpoints

| Endpoint | Description |
|----------|-------------|
| \`GET /health\` | Health check |
| \`GET /ready\` | Readiness check |

## Architecture

This service follows Domain-Driven Design (DDD):

\`\`\`
cmd/server/          # Entry point
internal/
  domain/            # Core business logic (no external deps)
  application/       # Use cases
  infrastructure/    # DB, Kafka, external APIs
    postgres/        # PostgreSQL repository
    kafka/           # Kafka producer/consumer
    grpc/            # gRPC clients to other services
  interfaces/        # Entry points
    http/            # REST handlers + middleware
    grpc/            # gRPC server
migrations/          # SQL migrations
proto/               # gRPC definitions
openapi/             # REST API definitions
\`\`\`

## Documentation

- [Architecture](../../docs/ARCHITECTURE.md) — system architecture
- [API Contracts](../../docs/API-CONTRACTS.md) — REST + gRPC + Kafka contracts
- [Patterns](../../docs/PATTERNS.md) — code patterns
EOF

  # Empty placeholder files for git
  touch "$SVC_DIR/internal/domain/.gitkeep"
  touch "$SVC_DIR/internal/application/.gitkeep"
  touch "$SVC_DIR/internal/infrastructure/postgres/.gitkeep"
  touch "$SVC_DIR/internal/infrastructure/kafka/.gitkeep"
  touch "$SVC_DIR/internal/infrastructure/grpc/.gitkeep"
  touch "$SVC_DIR/internal/interfaces/http/handlers/.gitkeep"
  touch "$SVC_DIR/internal/interfaces/http/middleware/.gitkeep"
  touch "$SVC_DIR/internal/interfaces/grpc/.gitkeep"
  touch "$SVC_DIR/migrations/.gitkeep"
  touch "$SVC_DIR/proto/.gitkeep"
  touch "$SVC_DIR/openapi/.gitkeep"

  echo "  ✓ $svc done"
done

echo ""
echo "✅ All 12 Go service skeletons created!"
