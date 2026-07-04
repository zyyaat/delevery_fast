# Repository Structure — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Monorepo Tool**: Turborepo + pnpm workspaces

---

## 1. Top-Level Structure

```
food-platform/
├── apps/                          # 7 React web applications
│   ├── customer-web/
│   ├── driver-web/
│   ├── restaurant-web/
│   ├── support-web/
│   ├── command-center/
│   ├── employee-portal/
│   └── field-supervisor-web/
│
├── packages/                      # Shared packages (TypeScript)
│   ├── api-client/                # HTTP + WebSocket client
│   ├── types/                     # TypeScript types (generated from OpenAPI)
│   ├── hooks/                     # Business logic hooks (TanStack Query)
│   ├── utils/                     # Shared utilities (formatters, validators)
│   ├── theme/                     # Design tokens (colors, typography, spacing)
│   ├── ui/                        # Shared components (Button, Card, Input, etc.)
│   ├── auth/                      # Auth helpers (JWT, refresh, WebAuthn)
│   └── eslint-config/             # Shared ESLint config
│
├── services/                      # 12 Go microservices
│   ├── auth/
│   ├── restaurant-catalog/
│   ├── menu/
│   ├── order/
│   ├── payment/
│   ├── delivery-matching/
│   ├── driver-management/
│   ├── geo/
│   ├── notification/
│   ├── fraud/
│   ├── promo/
│   └── analytics/
│
├── infra/                         # Infrastructure
│   ├── docker-compose.yml         # Dev environment
│   ├── docker-compose.test.yml    # Test environment
│   ├── terraform/                 # AWS IaC
│   ├── k8s/                       # Kubernetes manifests
│   └── helm/                      # Helm charts
│
├── tools/                         # Tooling
│   ├── codegen/                   # OpenAPI → Go/TS codegen
│   ├── scripts/                   # Setup & utility scripts
│   └── seeders/                   # DB seed data
│
├── docs/                          # All documentation
│   ├── ARCHITECTURE.md
│   ├── API-CONTRACTS.md
│   ├── PATTERNS.md
│   ├── REPO-STRUCTURE.md
│   ├── ROADMAP.md
│   ├── REFERENCES.md
│   ├── SESSIONS-LOG.md
│   └── adr/                       # Architecture Decision Records
│       ├── 001-go-for-backend.md
│       ├── 002-postgresql-over-mysql.md
│       └── ...
│
├── .github/                       # CI/CD
│   └── workflows/
│       ├── ci.yml                 # Lint + test + build
│       ├── deploy-staging.yml
│       └── deploy-prod.yml
│
├── package.json                   # Root package.json
├── pnpm-workspace.yaml            # pnpm workspace config
├── turbo.json                     # Turborepo config
├── tsconfig.base.json             # Shared TS config
├── .gitignore
├── .editorconfig
├── Makefile                       # Common commands
└── README.md
```

---

## 2. Root Configuration Files

### 2.1 `package.json` (root)

```json
{
  "name": "food-platform",
  "version": "1.0.0",
  "private": true,
  "scripts": {
    "dev": "turbo dev",
    "build": "turbo build",
    "test": "turbo test",
    "lint": "turbo lint",
    "format": "prettier --write \"**/*.{ts,tsx,md}\"",
    "clean": "turbo clean && rm -rf node_modules",
    "codegen": "turbo codegen",
    "type-check": "turbo type-check"
  },
  "devDependencies": {
    "turbo": "^2.0.0",
    "prettier": "^3.3.0",
    "typescript": "^5.5.0",
    "@food-platform/eslint-config": "workspace:*"
  },
  "engines": {
    "node": ">=20.0.0",
    "pnpm": ">=9.0.0"
  },
  "packageManager": "pnpm@9.5.0"
}
```

### 2.2 `pnpm-workspace.yaml`

```yaml
packages:
  - "apps/*"
  - "packages/*"
  - "tools/*"
```

### 2.3 `turbo.json`

```json
{
  "$schema": "https://turbo.build/schema.json",
  "globalDependencies": ["**/.env.*local"],
  "globalEnv": ["NODE_ENV", "CI"],
  "tasks": {
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", ".next/**", "build/**"]
    },
    "dev": {
      "cache": false,
      "persistent": true
    },
    "lint": {
      "dependsOn": ["^build"]
    },
    "test": {
      "dependsOn": ["^build"],
      "outputs": ["coverage/**"]
    },
    "type-check": {
      "dependsOn": ["^build"]
    },
    "codegen": {
      "dependsOn": ["^codegen"],
      "outputs": ["src/generated/**"]
    },
    "clean": {
      "cache": false
    }
  }
}
```

### 2.4 `tsconfig.base.json`

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "lib": ["ES2022", "DOM", "DOM.Iterable"],
    "module": "ESNext",
    "moduleResolution": "Bundler",
    "jsx": "react-jsx",
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noFallthroughCasesInSwitch": true,
    "noImplicitReturns": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "isolatedModules": true,
    "baseUrl": ".",
    "paths": {
      "@food-platform/api-client": ["./packages/api-client/src"],
      "@food-platform/types": ["./packages/types/src"],
      "@food-platform/hooks": ["./packages/hooks/src"],
      "@food-platform/utils": ["./packages/utils/src"],
      "@food-platform/theme": ["./packages/theme/src"],
      "@food-platform/ui": ["./packages/ui/src"],
      "@food-platform/auth": ["./packages/auth/src"]
    }
  }
}
```

### 2.5 `Makefile`

```makefile
.PHONY: dev build test lint codegen clean docker-up docker-down

dev:
\tturbo dev

build:
\tturbo build

test:
\tturbo test

lint:
\tturbo lint

codegen:
\tturbo codegen

type-check:
\tturbo type-check

docker-up:
\tdocker compose -f infra/docker-compose.yml up -d

docker-down:
\tdocker compose -f infra/docker-compose.yml down

migrate:
\t@for svc in $(SERVICES); do \\
\t\techo "Migrating $$svc..."; \\
\t\tcd services/$$svc && make migrate && cd ../..; \\
\tdone

clean:
\tturbo clean
\trm -rf node_modules

setup: docker-up codegen migrate
\t@echo "✅ Setup complete! Run 'make dev' to start."
```

---

## 3. Web Apps (`apps/`)

Each web app has identical structure. Here's `customer-web` as example:

```
apps/customer-web/
├── src/
│   ├── main.tsx                       # Entry point
│   ├── App.tsx                        # Root + routes
│   ├── routes/                        # Page-level components
│   │   ├── home/
│   │   │   ├── HomePage.tsx
│   │   │   ├── HomePage.test.tsx
│   │   │   └── index.ts
│   │   ├── restaurant-detail/
│   │   ├── cart/
│   │   ├── checkout/
│   │   ├── order-tracking/
│   │   ├── orders-history/
│   │   ├── profile/
│   │   └── auth/
│   │       ├── LoginPage.tsx
│   │       └── OtpPage.tsx
│   ├── components/                    # App-specific
│   │   ├── RestaurantCard/
│   │   │   ├── RestaurantCard.tsx
│   │   │   ├── RestaurantCard.test.tsx
│   │   │   └── index.ts
│   │   ├── CartItem/
│   │   ├── OrderTracker/
│   │   ├── AddressPicker/
│   │   └── PaymentSelector/
│   ├── features/                      # Feature modules (slice)
│   │   ├── auth/
│   │   │   ├── auth-store.ts          # Zustand store
│   │   │   ├── auth-api.ts
│   │   │   └── auth-hooks.ts
│   │   ├── cart/
│   │   │   ├── cart-store.ts
│   │   │   └── cart-hooks.ts
│   │   └── orders/
│   │       ├── orders-hooks.ts
│   │       └── orders-utils.ts
│   ├── hooks/                         # App-specific hooks
│   ├── lib/                           # App-specific utilities
│   │   ├── analytics.ts
│   │   ├── sentry.ts
│   │   └── constants.ts
│   ├── styles/
│   │   ├── globals.css
│   │   └── tailwind.config.ts
│   ├── types/
│   │   └── index.ts                   # App-specific types
│   └── test/
│       ├── setup.ts
│       └── utils.tsx
├── public/
│   ├── favicon.ico
│   └── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
├── playwright.config.ts               # E2E config
├── .eslintrc.cjs
├── .env.example
└── README.md
```

### 3.1 Web App `package.json` (template)

```json
{
  "name": "@food-platform/customer-web",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "test": "vitest",
    "test:e2e": "playwright test",
    "lint": "eslint src --ext .ts,.tsx",
    "type-check": "tsc --noEmit",
    "clean": "rm -rf dist node_modules/.vite"
  },
  "dependencies": {
    "@food-platform/api-client": "workspace:*",
    "@food-platform/types": "workspace:*",
    "@food-platform/hooks": "workspace:*",
    "@food-platform/utils": "workspace:*",
    "@food-platform/theme": "workspace:*",
    "@food-platform/ui": "workspace:*",
    "@food-platform/auth": "workspace:*",
    "react": "^18.3.0",
    "react-dom": "^18.3.0",
    "react-router-dom": "^6.24.0",
    "@tanstack/react-query": "^5.51.0",
    "zustand": "^4.5.0",
    "axios": "^1.7.0",
    "zod": "^3.23.0",
    "react-hook-form": "^7.52.0",
    "@hookform/resolvers": "^3.9.0",
    "tailwindcss": "^3.4.0",
    "clsx": "^2.1.0",
    "date-fns": "^3.6.0",
    "@sentry/react": "^8.20.0"
  },
  "devDependencies": {
    "@types/react": "^18.3.0",
    "@types/react-dom": "^18.3.0",
    "@vitejs/plugin-react": "^4.3.0",
    "vite": "^5.3.0",
    "vitest": "^2.0.0",
    "@testing-library/react": "^16.0.0",
    "@testing-library/jest-dom": "^6.4.0",
    "@testing-library/user-event": "^14.5.0",
    "@playwright/test": "^1.45.0",
    "eslint": "^8.57.0",
    "@food-platform/eslint-config": "workspace:*"
  }
}
```

---

## 4. Shared Packages (`packages/`)

### 4.1 `packages/api-client/`

```
packages/api-client/
├── src/
│   ├── index.ts                  # Exports
│   ├── client.ts                 # Axios instance + interceptors
│   ├── endpoints/                # Endpoint-specific functions
│   │   ├── auth.ts
│   │   ├── restaurants.ts
│   │   ├── orders.ts
│   │   ├── drivers.ts
│   │   ├── payments.ts
│   │   └── ...
│   ├── websocket/
│   │   ├── client.ts             # WS client with reconnect
│   │   └── events.ts             # Event type definitions
│   ├── errors.ts                 # Error normalization
│   └── types.ts                  # Re-export from @food-platform/types
├── package.json
├── tsconfig.json
└── README.md
```

### 4.2 `packages/types/`

```
packages/types/
├── src/
│   ├── index.ts                  # All exports
│   ├── generated/                # Generated from OpenAPI (DO NOT EDIT)
│   │   └── api.ts
│   ├── branded.ts                # Branded types (OrderID, CustomerID, ...)
│   ├── enums.ts                  # OrderStatus, UserRole, etc.
│   ├── entities.ts               # Order, Restaurant, Driver, etc.
│   ├── api.ts                    # Request/Response types
│   ├── events.ts                 # WebSocket + Kafka event types
│   └── common.ts                 # Money, Coordinates, Address, Pagination
├── package.json
├── tsconfig.json
└── README.md
```

### 4.3 `packages/hooks/`

```
packages/hooks/
├── src/
│   ├── index.ts
│   ├── auth/
│   │   ├── useLogin.ts
│   │   ├── useLogout.ts
│   │   └── useAuth.ts
│   ├── restaurants/
│   │   ├── useNearbyRestaurants.ts
│   │   ├── useRestaurant.ts
│   │   └── useMenu.ts
│   ├── orders/
│   │   ├── useCreateOrder.ts
│   │   ├── useOrder.ts
│   │   ├── useActiveOrders.ts
│   │   └── useOrderTracking.ts
│   ├── cart/
│   │   ├── useCart.ts
│   │   └── useApplyCoupon.ts
│   ├── drivers/
│   │   └── useDriverEarnings.ts
│   └── common/
│       ├── useWebSocket.ts
│       ├── useDebounce.ts
│       └── useMediaQuery.ts
├── package.json
└── tsconfig.json
```

### 4.4 `packages/ui/`

```
packages/ui/
├── src/
│   ├── index.ts
│   ├── components/
│   │   ├── Button/
│   │   │   ├── Button.tsx
│   │   │   ├── Button.test.tsx
│   │   │   └── index.ts
│   │   ├── Card/
│   │   ├── Input/
│   │   ├── Modal/
│   │   ├── Rating/
│   │   ├── Spinner/
│   │   ├── Toast/
│   │   ├── Tabs/
│   │   ├── Badge/
│   │   └── ...
│   ├── theme/
│   │   ├── colors.ts
│   │   ├── typography.ts
│   │   └── index.ts
│   └── utils/
│       └── cn.ts                 # clsx + tailwind-merge
├── package.json
├── tsconfig.json
└── README.md
```

### 4.5 `packages/auth/`

```
packages/auth/
├── src/
│   ├── index.ts
│   ├── stores/
│   │   └── auth-store.ts         # Zustand store (shared)
│   ├── jwt.ts                    # JWT decode (no verify on client)
│   ├── webauthn.ts               # WebAuthn helpers
│   ├── otp.ts                    # OTP request/verify
│   ├── refresh.ts                # Token refresh logic
│   └── guards/
│       ├── require-auth.tsx      # Route guard
│       └── require-role.tsx
├── package.json
└── tsconfig.json
```

### 4.6 `packages/utils/`

```
packages/utils/
├── src/
│   ├── index.ts
│   ├── format.ts                 # formatEGP, formatDistance, formatDate
│   ├── validate.ts               # isPhone, isEmail, isUUID
│   ├── crypto.ts                 # hash, randomString
│   ├── array.ts                  # groupBy, sortBy, chunk
│   ├── object.ts                 # pick, omit, deepEqual
│   └── constants.ts              # EGP_CURRENCY, PHONE_REGEX, etc.
└── package.json
```

### 4.7 `packages/theme/`

```
packages/theme/
├── src/
│   ├── index.ts
│   ├── tokens.ts                 # Design tokens
│   ├── tailwind.config.ts        # Shared Tailwind config
│   └── css/
│       └── globals.css           # CSS reset + Tailwind base
└── package.json
```

**Tokens**:

```typescript
export const tokens = {
  colors: {
    primary: '#00D4FF',      // cyan
    accent: '#B14EFF',       // purple
    success: '#00E58A',
    warning: '#FFB800',
    danger: '#FF4757',
    bg: '#0A0E1A',
    text: '#EAF0FF',
    textDim: '#9BA8C7',
  },
  typography: {
    fontFamily: {
      heading: 'Inter, Cairo, sans-serif',
      body: 'Cairo, Inter, sans-serif',
      numeric: 'JetBrains Mono, Inter, monospace',
    },
    fontSize: {
      display: '56px',
      h1: '40px',
      h2: '30px',
      h3: '22px',
      body: '16px',
      small: '13px',
      micro: '11px',
    },
  },
  spacing: { /* ... */ },
  borderRadius: { /* ... */ },
}
```

---

## 5. Go Services (`services/`)

Each service has identical structure. Here's `order` as example:

```
services/order/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Core business logic (no deps)
│   │   ├── order.go
│   │   ├── status.go            # State machine
│   │   ├── errors.go
│   │   └── events.go
│   ├── application/             # Use cases
│   │   ├── create_order.go
│   │   ├── cancel_order.go
│   │   ├── update_status.go
│   │   └── get_order.go
│   ├── infrastructure/          # External concerns
│   │   ├── postgres/
│   │   │   ├── order_repository.go
│   │   │   └── models.go
│   │   ├── kafka/
│   │   │   ├── producer.go
│   │   │   └── consumer.go
│   │   └── grpc/
│   │       ├── fraud_client.go
│   │       └── payment_client.go
│   └── interfaces/              # Entry points
│       ├── http/
│       │   ├── handlers/
│       │   │   ├── create_order.go
│       │   │   ├── get_order.go
│       │   │   └── cancel_order.go
│       │   ├── middleware/
│       │   │   ├── auth.go
│       │   │   ├── request_id.go
│       │   │   └── logging.go
│       │   ├── server.go
│       │   └── routes.go
│       └── grpc/
│           └── server.go
├── migrations/
│   ├── 001_create_orders.sql
│   ├── 002_add_partitioning.sql
│   └── 003_add_audit_columns.sql
├── proto/
│   └── order.proto
├── openapi/
│   └── order.yaml
├── tests/
│   ├── integration/
│   └── e2e/
├── deployments/
│   └── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 5.1 Go Service `Makefile`

```makefile
.PHONY: build run test lint migrate create-migration

build:
\tgo build -o bin/server ./cmd/server

run:
\tgo run ./cmd/server

test:
\tgo test -v -race -cover ./...

lint:
\tgolangci-lint run

migrate:
\tmigrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
\tmigrate -path migrations -database "$(DATABASE_URL)" down

create-migration:
\t@migrate create -ext sql -dir migrations -seq $(NAME)

grpc-gen:
\tprotoc --go_out=. --go-grpc_out=. proto/*.proto

docker-build:
\tdocker build -t food-platform/order:latest -f deployments/Dockerfile .

clean:
\trm -rf bin/
```

### 5.2 Go Service `Dockerfile`

```dockerfile
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
EXPOSE 8080 9090

ENTRYPOINT ["/server"]
```

---

## 6. Infrastructure (`infra/`)

### 6.1 `docker-compose.yml` (dev)

```yaml
version: '3.9'

services:
  postgres:
    image: postgis/postgis:16-3.4
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: postgres
    ports:
      - "5432:5432"
    volumes:
      - pg_data:/var/lib/postgresql/data
      - ./init-db.sql:/docker-entrypoint-initdb.d/init-db.sql

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data

  zookeeper:
    image: confluentinc/cp-zookeeper:7.6.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
    ports:
      - "2181:2181"

  kafka:
    image: confluentinc/cp-kafka:7.6.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

  schema-registry:
    image: confluentinc/cp-schema-registry:7.6.0
    depends_on:
      - kafka
    ports:
      - "8081:8081"
    environment:
      SCHEMA_REGISTRY_HOST_NAME: schema-registry
      SCHEMA_REGISTRY_KAFKASTORE_CONNECTION_URL: zookeeper:2181

  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.13.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
    volumes:
      - es_data:/usr/share/elasticsearch/data

  clickhouse:
    image: clickhouse/clickhouse-server:24.3
    ports:
      - "8123:8123"
      - "9000:9000"
    volumes:
      - ch_data:/var/lib/clickhouse

  minio:
    image: minio/minio:latest
    ports:
      - "9001:9000"
      - "9002:9001"
    environment:
      MINIO_ROOT_USER: minio
      MINIO_ROOT_PASSWORD: minio123
    command: server /data --console-address ":9001"
    volumes:
      - minio_data:/data

volumes:
  pg_data:
  redis_data:
  es_data:
  ch_data:
  minio_data:
```

### 6.2 `init-db.sql`

```sql
-- Create databases for each service
CREATE DATABASE auth_db;
CREATE DATABASE restaurants_db;
CREATE DATABASE menus_db;
CREATE DATABASE orders_db;
CREATE DATABASE payments_db;
CREATE DATABASE dispatch_db;
CREATE DATABASE drivers_db;
CREATE DATABASE notifications_db;
CREATE DATABASE fraud_db;
CREATE DATABASE promos_db;
CREATE DATABASE audit_db;
CREATE DATABASE analytics_db;
```

### 6.3 `terraform/` (structure)

```
infra/terraform/
├── modules/
│   ├── vpc/
│   ├── eks/
│   ├── rds/
│   ├── elasticache/
│   ├── msk/
│   └── s3/
├── environments/
│   ├── dev/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── terraform.tfvars
│   ├── staging/
│   └── prod/
└── README.md
```

---

## 7. CI/CD (`.github/workflows/`)

### 7.1 `ci.yml`

```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: pnpm/action-setup@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'pnpm'
      - run: pnpm install --frozen-lockfile
      - run: pnpm lint
      - run: pnpm type-check
      - run: pnpm test
      - run: pnpm build

  backend:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        service: [auth, restaurant-catalog, menu, order, payment,
                  delivery-matching, driver-management, geo,
                  notification, fraud, promo, analytics]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Lint
        working-directory: services/${{ matrix.service }}
        run: |
          go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
          golangci-lint run
      - name: Test
        working-directory: services/${{ matrix.service }}
        run: go test -v -race -cover ./...
      - name: Build
        working-directory: services/${{ matrix.service }}
        run: docker build -t food-platform/${{ matrix.service }}:${{ github.sha }} .
```

---

## 8. Workspace Dependency Graph

```
                    ┌────────────────────────────────────────┐
                    │              packages/types              │
                    └────────────────────┬────────────────────┘
                                         │
            ┌────────────────────────────┼────────────────────┐
            │                            │                    │
   ┌────────▼────────┐       ┌──────────▼─────────┐  ┌───────▼────────┐
   │  packages/api-   │       │  packages/hooks    │  │  packages/ui   │
   │     client       │       │                    │  │                 │
   └────────┬─────────┘       └──────────┬─────────┘  └───────┬────────┘
            │                            │                    │
            └─────────────┬──────────────┴────────────────────┘
                          │
                ┌─────────▼──────────┐
                │  packages/auth     │
                └─────────┬──────────┘
                          │
        ┌─────────────────┼─────────────────┐
        │                 │                 │
   ┌────▼─────┐    ┌──────▼──────┐   ┌──────▼──────┐
   │ customer │    │   driver    │   │ restaurant  │   ... (7 apps)
   │   web    │    │     web     │   │     web     │
   └──────────┘    └─────────────┘   └─────────────┘
```

All 7 apps depend on:
- `@food-platform/api-client`
- `@food-platform/types`
- `@food-platform/hooks`
- `@food-platform/utils`
- `@food-platform/theme`
- `@food-platform/ui`
- `@food-platform/auth`

---

> **Next**: Read `ROADMAP.md` for the 12-week execution plan.
