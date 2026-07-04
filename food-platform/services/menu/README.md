# Menu Service

Go microservice for the Food Delivery Platform.

## Quick Start

```bash
# Run locally
make run

# Build
make build

# Test
make test

# Run migrations
DATABASE_URL=postgres://postgres:postgres@localhost:5432/menu_db?sslmode=disable make migrate
```

## Endpoints

| Endpoint | Description |
|----------|-------------|
| `GET /health` | Health check |
| `GET /ready` | Readiness check |

## Architecture

This service follows Domain-Driven Design (DDD):

```
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
```

## Documentation

- [Architecture](../../docs/ARCHITECTURE.md) — system architecture
- [API Contracts](../../docs/API-CONTRACTS.md) — REST + gRPC + Kafka contracts
- [Patterns](../../docs/PATTERNS.md) — code patterns
