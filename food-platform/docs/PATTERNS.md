# Code Patterns — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active — Mandatory  
> **Last Updated**: 2026-07-04

This document defines **mandatory code patterns** for all 7 web apps and 12 Go services. Any code that doesn't follow these patterns MUST be flagged in PR review.

---

## Table of Contents

1. [General Principles](#1-general-principles)
2. [Go Backend Patterns](#2-go-backend-patterns)
3. [React Frontend Patterns](#3-react-frontend-patterns)
4. [TypeScript Patterns](#4-typescript-patterns)
5. [Database Patterns](#5-database-patterns)
6. [Error Handling](#6-error-handling)
7. [Testing Patterns](#7-testing-patterns)
8. [Logging & Observability](#8-logging--observability)
9. [Security Patterns](#9-security-patterns)
10. [Naming Conventions](#10-naming-conventions)

---

## 1. General Principles

### 1.1 The 7 Principles

1. **Contract-First**: API contracts (`openapi.yaml`) before any code.
2. **Type Safety**: Generate types from contracts; never hand-write.
3. **Defense in Depth**: Multiple layers of validation (client, API gateway, service, DB).
4. **Fail Fast**: Validate early; reject invalid input at the edge.
5. **Idempotency**: All state-mutating POSTs must accept `X-Idempotency-Key`.
6. **Observability**: Every request has a `request_id`; every action is logged.
7. **Tests Mandatory**: 80%+ coverage on business logic; no PR merges without tests.

### 1.2 Boy Scout Rule

Leave the codebase cleaner than you found it. If you see a small issue, fix it.

---

## 2. Go Backend Patterns

### 2.1 Project Structure (per service)

```
services/order/
├── cmd/
│   └── server/
│       └── main.go              # Entry point — wire everything
├── internal/
│   ├── domain/                  # Core business logic (no deps)
│   │   ├── order.go             # Entity
│   │   ├── status.go            # State machine
│   │   ├── errors.go            # Domain errors
│   │   └── events.go            # Domain events
│   ├── application/             # Use cases
│   │   ├── create_order.go
│   │   ├── cancel_order.go
│   │   └── update_status.go
│   ├── infrastructure/          # External concerns
│   │   ├── postgres/
│   │   │   ├── order_repository.go
│   │   │   └── models.go
│   │   ├── kafka/
│   │   │   ├── producer.go
│   │   │   └── consumer.go
│   │   └── grpc/
│   │       └── fraud_client.go
│   └── interfaces/              # Entry points (HTTP, gRPC)
│       ├── http/
│       │   ├── handlers/
│       │   ├── middleware/
│       │   └── server.go
│       └── grpc/
│           └── server.go
├── migrations/
│   ├── 001_create_orders.sql
│   └── 002_add_partitioning.sql
├── proto/                       # gRPC definitions
│   └── order.proto
├── openapi/                     # REST definitions
│   └── order.yaml
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### 2.2 Domain Entity Pattern

```go
// internal/domain/order.go

package domain

import (
    "time"
    "github.com/google/uuid"
)

type OrderID uuid.UUID

type Order struct {
    id           OrderID
    customerID   uuid.UUID
    restaurantID uuid.UUID
    driverID     *uuid.UUID  // nullable
    status       OrderStatus
    total        Money
    items        []OrderItem
    deliveryAddr Address
    createdAt    time.Time
    updatedAt    time.Time
}

// NewOrder — constructor with validation
func NewOrder(customerID, restaurantID uuid.UUID, items []OrderItem, addr Address) (*Order, error) {
    if len(items) == 0 {
        return nil, ErrEmptyOrder
    }
    total := calculateTotal(items)
    return &Order{
        id:           OrderID(uuid.New()),
        customerID:   customerID,
        restaurantID: restaurantID,
        status:       StatusPending,
        total:        total,
        items:        items,
        deliveryAddr: addr,
        createdAt:    time.Now().UTC(),
    }, nil
}

// Business logic methods — never expose internal state directly
func (o *Order) Confirm(driverID uuid.UUID) error {
    if o.status != StatusPending {
        return ErrInvalidTransition
    }
    o.driverID = &driverID
    o.status = StatusConfirmed
    o.updatedAt = time.Now().UTC()
    return nil
}
```

**Rules**:
- Domain entities have NO dependencies on infrastructure (no `database/sql`, no `kafka`).
- All fields are private; expose via methods.
- Constructors validate input.
- State transitions are methods (not direct field assignment).

### 2.3 Repository Pattern

```go
// internal/infrastructure/postgres/order_repository.go

package postgres

import (
    "context"
    "database/sql"
    "github.com/google/uuid"
    "food-platform/services/order/internal/domain"
)

type OrderRepository interface {
    Save(ctx context.Context, order *domain.Order) error
    FindByID(ctx context.Context, id domain.OrderID) (*domain.Order, error)
    Update(ctx context.Context, order *domain.Order) error
}

type pgOrderRepository struct {
    db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
    return &pgOrderRepository{db: db}
}

func (r *pgOrderRepository) Save(ctx context.Context, order *domain.Order) error {
    query := `INSERT INTO orders (id, customer_id, restaurant_id, status, total, ...)
              VALUES ($1, $2, $3, $4, $5, ...)`
    _, err := r.db.ExecContext(ctx, query,
        order.ID(), order.CustomerID(), order.RestaurantID(),
        order.Status(), order.Total(),
    )
    return err
}
```

**Rules**:
- Repository interface defined in `domain` package.
- Implementation in `infrastructure`.
- Always accept `context.Context` as first parameter.
- Use `sql.DB` (not ORM) for performance-critical services; GORM acceptable for non-critical.

### 2.4 HTTP Handler Pattern

```go
// internal/interfaces/http/handlers/create_order.go

package handlers

import (
    "encoding/json"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
)

type CreateOrderHandler struct {
    createOrder app.CreateOrderUseCase
    logger      Logger
}

func (h *CreateOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    var req CreateOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
        return
    }

    if err := req.Validate(); err != nil {
        writeError(w, http.StatusUnprocessableEntity, "VALIDATION_FAILED", err.Error())
        return
    }

    // Idempotency check
    idempotencyKey := r.Header.Get("X-Idempotency-Key")
    if idempotencyKey == "" {
        idempotencyKey = uuid.New().String()
    }

    cmd := app.CreateOrderCommand{
        CustomerID:   ctx.Value("user_id").(uuid.UUID),
        RestaurantID: req.RestaurantID,
        Items:        req.Items,
        AddressID:    req.AddressID,
        IdempotencyKey: idempotencyKey,
    }

    order, err := h.createOrder.Execute(ctx, cmd)
    if err != nil {
        h.logger.Error("create_order failed", "error", err, "command", cmd)
        writeDomainError(w, err)
        return
    }

    writeJSON(w, http.StatusCreated, toOrderResponse(order))
}
```

### 2.5 gRPC Service Pattern

```go
// internal/interfaces/grpc/server.go

package grpc

import (
    "context"
    "food-platform/services/order/internal/application"
    pb "food-platform/services/order/proto"
)

type OrderGRPCServer struct {
    pb.UnimplementedOrderServiceServer
    getOrder app.GetOrderUseCase
}

func (s *OrderGRPCServer) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
    order, err := s.getOrder.Execute(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
    }
    return toProto(order), nil
}
```

### 2.6 Kafka Producer Pattern

```go
// internal/infrastructure/kafka/producer.go

package kafka

import (
    "context"
    "encoding/json"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventProducer struct {
    producer *kafka.Producer
    topic    string
}

func (p *EventProducer) Publish(ctx context.Context, event DomainEvent) error {
    payload, err := json.Marshal(event)
    if err != nil {
        return err
    }

    msg := &kafka.Message{
        TopicPartition: kafka.TopicPartition{
            Topic:     &p.topic,
            Partition: kafka.PartitionAny,
        },
        Key:   []byte(event.AggregateID()),  // for partition ordering
        Value: payload,
        Headers: []kafka.Header{
            {Key: "event_type", Value: []byte(event.Type())},
            {Key: "event_id", Value: []byte(event.ID())},
            {Key: "event_time", Value: []byte(event.Time().Format(time.RFC3339))},
        },
    }

    return p.producer.Produce(msg, nil)
}
```

### 2.7 Kafka Consumer Pattern

```go
// internal/infrastructure/kafka/consumer.go

package kafka

import (
    "context"
    "github.com/confluentinc/confluent-kafka-go/kafka"
)

type EventHandler interface {
    Handle(ctx context.Context, event Event) error
}

type EventConsumer struct {
    consumer *kafka.Consumer
    handlers map[string]EventHandler
}

func (c *EventConsumer) Run(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
            if err != nil {
                continue
            }

            eventType := getHeader(msg, "event_type")
            handler, ok := c.handlers[eventType]
            if !ok {
                continue  // no handler for this event
            }

            event := parseEvent(eventType, msg.Value)
            if err := handler.Handle(ctx, event); err != nil {
                log.Error("event_handler_failed", "event", event, "error", err)
                // TODO: dead-letter queue
                continue
            }

            c.consumer.CommitMessage(msg)
        }
    }
}
```

### 2.8 Configuration Pattern

```go
// cmd/server/main.go

type Config struct {
    HTTPPort        int           `env:"HTTP_PORT" default:"8080"`
    GRPCPort        int           `env:"GRPC_PORT" default:"9090"`
    DatabaseURL     string        `env:"DATABASE_URL" required:"true"`
    RedisURL        string        `env:"REDIS_URL" required:"true"`
    KafkaBrokers    []string      `env:"KAFKA_BROKERS" required:"true"`
    LogLevel        string        `env:"LOG_LEVEL" default:"info"`
    ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" default:"30s"`
}

func LoadConfig() (*Config, error) {
    var cfg Config
    if err := envconfig.Process("", &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
```

### 2.9 Graceful Shutdown

```go
func main() {
    cfg := LoadConfig()

    srv := server.New(cfg)

    // Start
    go func() {
        if err := srv.Start(); err != nil && err != http.ErrServerClosed {
            log.Fatal("server_failed", "error", err)
        }
    }()

    // Wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info("shutting_down")

    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("shutdown_failed", "error", err)
    }

    log.Info("shutdown_complete")
}
```

---

## 3. React Frontend Patterns

### 3.1 App Structure (per app)

```
apps/customer-web/
├── src/
│   ├── main.tsx                # Entry point
│   ├── App.tsx                 # Root component + routes
│   ├── routes/                 # Page-level components
│   │   ├── home/
│   │   ├── restaurant-detail/
│   │   ├── cart/
│   │   ├── checkout/
│   │   └── order-tracking/
│   ├── components/             # App-specific components
│   │   ├── RestaurantCard/
│   │   ├── CartItem/
│   │   └── OrderTracker/
│   ├── features/               # Feature modules (slice)
│   │   ├── auth/
│   │   ├── cart/
│   │   └── orders/
│   ├── hooks/                  # App-specific hooks
│   ├── lib/                    # App-specific utilities
│   ├── styles/
│   ├── types/
│   └── test/
├── public/
├── package.json
├── tsconfig.json
├── vite.config.ts
└── README.md
```

### 3.2 Component Pattern (Function Components + Hooks)

```tsx
// src/components/RestaurantCard/RestaurantCard.tsx

import { FC } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, Text, Rating } from '@food-platform/ui'
import { Restaurant } from '@food-platform/types'
import { formatEGP, formatDistance } from '@food-platform/utils'
import { trackEvent } from '../../lib/analytics'

interface RestaurantCardProps {
  restaurant: Restaurant
  distanceKm: number
}

export const RestaurantCard: FC<RestaurantCardProps> = ({ restaurant, distanceKm }) => {
  const navigate = useNavigate()

  const handleClick = () => {
    trackEvent('restaurant_card_clicked', { restaurant_id: restaurant.id })
    navigate(`/restaurants/${restaurant.id}`)
  }

  return (
    <Card onClick={handleClick} className="cursor-pointer hover:shadow-md">
      <img
        src={restaurant.cover_url}
        alt={restaurant.name}
        className="w-full h-32 object-cover rounded-t-lg"
        loading="lazy"
      />
      <div className="p-4">
        <Text variant="h4" className="font-bold">{restaurant.name}</Text>
        <div className="flex items-center gap-2 mt-1">
          <Rating value={restaurant.rating} />
          <Text variant="small" className="text-gray-500">
            ({restaurant.rating_count})
          </Text>
        </div>
        <div className="flex justify-between mt-2">
          <Text variant="small">{formatDistance(distanceKm)}</Text>
          <Text variant="small">{restaurant.eta_minutes_min}-{restaurant.eta_minutes_max} min</Text>
        </div>
        {restaurant.promo && (
          <div className="mt-2 text-primary text-sm font-semibold">
            🎁 {restaurant.promo.title}
          </div>
        )}
      </div>
    </Card>
  )
}
```

**Rules**:
- Always function components (no class components).
- TypeScript types for all props.
- Use Tailwind for styling (no inline styles except dynamic).
- Use `@food-platform/ui` components when possible.
- Track analytics events on user interactions.
- Lazy-load images (`loading="lazy"`).

### 3.3 Custom Hooks Pattern (Business Logic)

```tsx
// packages/shared/hooks/useOrder.ts

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient } from '@food-platform/api-client'
import { useToast } from '@food-platform/ui'
import { Order } from '@food-platform/types'

export function useOrder(orderId: string) {
  return useQuery<Order>({
    queryKey: ['order', orderId],
    queryFn: () => apiClient.getOrder(orderId),
    enabled: !!orderId,
    staleTime: 30_000,  // 30s
  })
}

export function useCreateOrder() {
  const queryClient = useQueryClient()
  const toast = useToast()

  return useMutation({
    mutationFn: apiClient.createOrder,
    onSuccess: (order) => {
      queryClient.setQueryData(['order', order.id], order)
      queryClient.invalidateQueries({ queryKey: ['orders', 'active'] })
      toast.success('تم إنشاء طلبك بنجاح!')
    },
    onError: (error) => {
      toast.error(error.message || 'فشل إنشاء الطلب')
    },
  })
}

export function useCancelOrder() {
  const queryClient = useQueryClient()
  const toast = useToast()

  return useMutation({
    mutationFn: ({ orderId, reason }: { orderId: string; reason: string }) =>
      apiClient.cancelOrder(orderId, { reason }),
    onSuccess: (_, { orderId }) => {
      queryClient.invalidateQueries({ queryKey: ['order', orderId] })
      queryClient.invalidateQueries({ queryKey: ['orders', 'active'] })
      toast.success('تم إلغاء الطلب')
    },
  })
}
```

### 3.4 State Management (Zustand)

```tsx
// src/features/auth/auth-store.ts

import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface AuthState {
  user: User | null
  accessToken: string | null
  refreshToken: string | null
  setAuth: (user: User, access: string, refresh: string) => void
  logout: () => void
  isAuthenticated: () => boolean
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      setAuth: (user, access, refresh) =>
        set({ user, accessToken: access, refreshToken: refresh }),
      logout: () =>
        set({ user: null, accessToken: null, refreshToken: null }),
      isAuthenticated: () => !!get().accessToken,
    }),
    { name: 'auth-storage' }
  )
)
```

### 3.5 API Client Pattern

```tsx
// packages/api-client/src/index.ts

import axios, { AxiosInstance, AxiosError } from 'axios'
import { useAuthStore } from '../stores/auth-store'

class ApiClient {
  private client: AxiosInstance

  constructor(baseURL: string) {
    this.client = axios.create({
      baseURL,
      timeout: 10_000,
      headers: { 'Content-Type': 'application/json' },
    })

    // Request interceptor — attach auth
    this.client.interceptors.request.use((config) => {
      const token = useAuthStore.getState().accessToken
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
      config.headers['X-Request-ID'] = crypto.randomUUID()
      return config
    })

    // Response interceptor — handle 401
    this.client.interceptors.response.use(
      (response) => response.data,
      async (error: AxiosError) => {
        if (error.response?.status === 401) {
          await refreshTokenAndRetry(error)
        }
        return Promise.reject(this.normalizeError(error))
      }
    )
  }

  async createOrder(payload: CreateOrderRequest): Promise<Order> {
    return this.client.post('/orders', payload, {
      headers: { 'X-Idempotency-Key': crypto.randomUUID() },
    })
  }

  async getOrder(id: string): Promise<Order> {
    return this.client.get(`/orders/${id}`)
  }

  // ... etc
}

export const apiClient = new ApiClient(import.meta.env.VITE_API_URL)
```

### 3.6 WebSocket Hook Pattern

```tsx
// packages/shared/hooks/useWebSocket.ts

import { useEffect, useRef, useState } from 'react'
import { useAuthStore } from '../stores/auth-store'

interface UseWebSocketOptions {
  onMessage?: (event: WSMessage) => void
  onOpen?: () => void
  onClose?: () => void
  reconnectInterval?: number
  maxReconnectAttempts?: number
}

export function useWebSocket(
  channel: string,
  options: UseWebSocketOptions = {}
) {
  const [isConnected, setIsConnected] = useState(false)
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectAttempts = useRef(0)
  const token = useAuthStore((s) => s.accessToken)

  useEffect(() => {
    if (!token) return

    const connect = () => {
      const ws = new WebSocket(
        `${import.meta.env.VITE_WS_URL}/ws?channel=${channel}`,
        ['food-platform.v1'],
        // headers not supported in browser; use query param
      )

      wsRef.current = ws

      ws.onopen = () => {
        setIsConnected(true)
        reconnectAttempts.current = 0
        options.onOpen?.()
        // Authenticate
        ws.send(JSON.stringify({ event: 'auth', payload: { token } }))
      }

      ws.onmessage = (e) => {
        const message = JSON.parse(e.data) as WSMessage
        options.onMessage?.(message)
      }

      ws.onclose = () => {
        setIsConnected(false)
        options.onClose?.()
        // Exponential backoff
        const delay = Math.min(
          1000 * 2 ** reconnectAttempts.current,
          options.maxReconnectAttempts ?? 30000
        )
        reconnectAttempts.current += 1
        if (reconnectAttempts.current < (options.maxReconnectAttempts ?? 10)) {
          setTimeout(connect, delay)
        }
      }
    }

    connect()

    return () => {
      wsRef.current?.close()
    }
  }, [channel, token])

  const send = (event: string, payload: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ event, payload }))
    }
  }

  return { isConnected, send }
}
```

### 3.7 Form Pattern (React Hook Form + Zod)

```tsx
// src/features/checkout/CheckoutForm.tsx

import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

const checkoutSchema = z.object({
  address_id: z.string().uuid(),
  payment_method: z.enum(['vodafone_cash', 'instapay', 'card', 'cod']),
  notes: z.string().max(500).optional(),
})

type CheckoutForm = z.infer<typeof checkoutSchema>

export function CheckoutForm({ onSubmit }: { onSubmit: (data: CheckoutForm) => void }) {
  const { register, handleSubmit, formState: { errors, isSubmitting } } =
    useForm<CheckoutForm>({ resolver: zodResolver(checkoutSchema) })

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* ... */}
    </form>
  )
}
```

---

## 4. TypeScript Patterns

### 4.1 Strict Mode

All `tsconfig.json` MUST have:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noFallthroughCasesInSwitch": true,
    "noImplicitReturns": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true
  }
}
```

### 4.2 Type Imports

```typescript
// Always use `import type` for type-only imports
import type { Order, Customer } from '@food-platform/types'
import { apiClient } from '@food-platform/api-client'
```

### 4.3 Never Use `any`

```typescript
// ❌ BAD
function processData(data: any) { ... }

// ✅ GOOD
function processData(data: unknown) {
  if (isOrderData(data)) {
    // ...
  }
}

// ✅ Or use generics
function processData<T>(data: T): T { ... }
```

### 4.4 Branded Types for IDs

```typescript
// packages/types/src/branded.ts

export type Brand<T, B> = T & { __brand: B }

export type OrderID = Brand<string, 'OrderID'>
export type CustomerID = Brand<string, 'CustomerID'>
export type RestaurantID = Brand<string, 'RestaurantID'>

// Usage
function getOrder(id: OrderID): Order { ... }

// Compile-time safety: can't pass CustomerID where OrderID expected
```

---

## 5. Database Patterns

### 5.1 Migrations

Use `golang-migrate`. One file per change, sequentially numbered.

```sql
-- services/order/migrations/001_create_orders.sql

CREATE TABLE orders (
    id              UUID PRIMARY KEY,
    customer_id     UUID NOT NULL REFERENCES users(id),
    restaurant_id   UUID NOT NULL REFERENCES restaurants(id),
    driver_id       UUID REFERENCES drivers(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    subtotal        DECIMAL(10,2) NOT NULL,
    delivery_fee    DECIMAL(10,2) NOT NULL,
    service_fee     DECIMAL(10,2) NOT NULL,
    discount        DECIMAL(10,2) NOT NULL DEFAULT 0,
    total           DECIMAL(10,2) NOT NULL,
    payment_method  VARCHAR(20) NOT NULL,
    delivery_address JSONB NOT NULL,
    latitude        DECIMAL(10,7),
    longitude       DECIMAL(10,7),
    eta_minutes     INT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (total >= 0)
);

CREATE INDEX idx_orders_customer ON orders (customer_id, created_at DESC);
CREATE INDEX idx_orders_restaurant ON orders (restaurant_id, status);
CREATE INDEX idx_orders_driver ON orders (driver_id, status);
```

### 5.2 Query Patterns

```go
// ✅ GOOD — use parameterized queries
db.QueryRowContext(ctx, `SELECT id, status FROM orders WHERE id = $1`, orderID)

// ❌ BAD — SQL injection risk
db.QueryRowContext(ctx, fmt.Sprintf(`SELECT id FROM orders WHERE id = '%s'`, orderID))
```

### 5.3 Transaction Pattern

```go
func (r *pgOrderRepository) CreateWithItems(ctx context.Context, order *domain.Order) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()  // safe: no-op if committed

    // Insert order
    if _, err := tx.ExecContext(ctx, insertOrderQuery, ...); err != nil {
        return fmt.Errorf("insert order: %w", err)
    }

    // Insert items
    for _, item := range order.Items() {
        if _, err := tx.ExecContext(ctx, insertItemQuery, ...); err != nil {
            return fmt.Errorf("insert item: %w", err)
        }
    }

    return tx.Commit()
}
```

---

## 6. Error Handling

### 6.1 Go Error Pattern

```go
// internal/domain/errors.go

package domain

import "errors"

var (
    ErrOrderNotFound       = errors.New("order not found")
    ErrInvalidTransition   = errors.New("invalid order status transition")
    ErrEmptyOrder          = errors.New("order must have at least one item")
    ErrPaymentDeclined     = errors.New("payment declined by provider")
)

// Wrapped errors for context
type OrderError struct {
    OrderID uuid.UUID
    Op      string  // operation
    Err     error
}

func (e *OrderError) Error() string {
    return fmt.Sprintf("%s: order %s: %v", e.Op, e.OrderID, e.Err)
}

func (e *OrderError) Unwrap() error { return e.Err }
```

```go
// Always wrap with context
if err := r.Save(ctx, order); err != nil {
    return fmt.Errorf("order_repository.save: %w", err)
}
```

### 6.2 React Error Pattern

```tsx
// ErrorBoundary component
class ErrorBoundary extends React.Component {
  state = { hasError: false, error: null }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    Sentry.captureException(error, { extra: { info } })
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} onRetry={() => location.reload()} />
    }
    return this.props.children
  }
}
```

```tsx
// Hook-level error handling
const { data, error, isLoading } = useOrder(orderId)

if (isLoading) return <Spinner />
if (error) return <ErrorMessage error={error} />
if (!data) return <NotFound />

return <OrderDetail order={data} />
```

---

## 7. Testing Patterns

### 7.1 Go Test Pyramid

```
        /\
       /e2e\        5% — full system tests (Docker)
      /------\
     /integ \     15% — service-to-service
    /----------\
   /   unit    \   80% — pure logic
  /--------------\
```

### 7.2 Go Unit Test

```go
// internal/domain/order_test.go

package domain

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewOrder_WithEmptyItems_ReturnsError(t *testing.T) {
    customerID := uuid.New()
    restaurantID := uuid.New()

    order, err := NewOrder(customerID, restaurantID, []OrderItem{}, Address{})

    require.Error(t, err)
    assert.Equal(t, ErrEmptyOrder, err)
    assert.Nil(t, order)
}

func TestOrder_Confirm_FromPending_Succeeds(t *testing.T) {
    order := newTestOrder(t, StatusPending)
    driverID := uuid.New()

    err := order.Confirm(driverID)

    require.NoError(t, err)
    assert.Equal(t, StatusConfirmed, order.Status())
    assert.Equal(t, driverID, *order.DriverID())
}

func TestOrder_Confirm_FromDelivered_ReturnsError(t *testing.T) {
    order := newTestOrder(t, StatusDelivered)
    driverID := uuid.New()

    err := order.Confirm(driverID)

    require.Error(t, err)
    assert.Equal(t, ErrInvalidTransition, err)
}
```

### 7.3 React Component Test

```tsx
// src/components/RestaurantCard/RestaurantCard.test.tsx

import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { RestaurantCard } from './RestaurantCard'
import { mockRestaurant } from './__mocks__'

test('renders restaurant name and rating', () => {
  render(
    <MemoryRouter>
      <RestaurantCard restaurant={mockRestaurant} distanceKm={1.2} />
    </MemoryRouter>
  )

  expect(screen.getByText('Pizza Hut Maadi')).toBeInTheDocument()
  expect(screen.getByText(/4\.6/)).toBeInTheDocument()
  expect(screen.getByText(/1\.2 km/)).toBeInTheDocument()
})

test('calls navigate on click', async () => {
  const { user } = renderWithRouter(
    <RestaurantCard restaurant={mockRestaurant} distanceKm={1.2} />
  )

  await user.click(screen.getByRole('article'))

  expect(screen.location.pathname).toBe(`/restaurants/${mockRestaurant.id}`)
})
```

### 7.4 E2E Test (Playwright)

```ts
// tests/e2e/order-flow.spec.ts

import { test, expect } from '@playwright/test'

test('customer can place an order', async ({ page }) => {
  await page.goto('http://localhost:5173')

  // Login
  await page.fill('[data-testid=phone-input]', '01012345678')
  await page.click('[data-testid=send-otp]')
  await page.fill('[data-testid=otp-input]', '123456')
  await page.click('[data-testid=verify-otp]')

  // Browse
  await page.click('text=Pizza Hut')
  await page.click('text=Margherita')
  await page.click('[data-testid=add-to-cart]')

  // Checkout
  await page.click('[data-testid=checkout]')
  await page.click('[data-testid=payment-vodafone-cash]')
  await page.click('[data-testid=place-order]')

  // Verify
  await expect(page.locator('[data-testid=order-confirmed]')).toBeVisible()
})
```

---

## 8. Logging & Observability

### 8.1 Go Logging (Structured)

```go
import "log/slog"

// ✅ GOOD — structured
slog.InfoContext(ctx, "order_created",
    "order_id", order.ID(),
    "customer_id", order.CustomerID(),
    "total", order.Total(),
    "duration_ms", time.Since(start).Milliseconds(),
)

// ❌ BAD — unstructured
log.Printf("Order %s created for customer %s with total %v",
    order.ID(), order.CustomerID(), order.Total())
```

### 8.2 Request ID Propagation

```go
// middleware/request_id.go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 8.3 OpenTelemetry Tracing

```go
// internal/application/create_order.go

func (uc *CreateOrderUseCase) Execute(ctx context.Context, cmd CreateOrderCommand) (*domain.Order, error) {
    ctx, span := tracer.Start(ctx, "create_order")
    defer span.End()

    span.SetAttributes(
        attribute.String("customer_id", cmd.CustomerID.String()),
        attribute.String("restaurant_id", cmd.RestaurantID.String()),
    )

    // ...
}
```

---

## 9. Security Patterns

### 9.1 Input Validation

```go
// Always validate at the edge
type CreateOrderRequest struct {
    RestaurantID uuid.UUID `json:"restaurant_id"`
    Items        []OrderItemRequest `json:"items"`
    AddressID    uuid.UUID `json:"address_id"`
}

func (r *CreateOrderRequest) Validate() error {
    if len(r.Items) == 0 {
        return ErrEmptyOrder
    }
    if len(r.Items) > 50 {
        return ErrTooManyItems
    }
    for _, item := range r.Items {
        if item.Quantity < 1 || item.Quantity > 20 {
            return ErrInvalidQuantity
        }
    }
    return nil
}
```

### 9.2 SQL Injection Prevention

Always use parameterized queries. **Never** concatenate strings into SQL.

### 9.3 XSS Prevention (React)

React escapes by default. Never use `dangerouslySetInnerHTML` without sanitization.

```tsx
// ❌ BAD
<div dangerouslySetInnerHTML={{ __html: userContent }} />

// ✅ GOOD
import DOMPurify from 'dompurify'
<div dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(userContent) }} />
```

### 9.4 CSRF Prevention

Use `SameSite=Strict` cookies + double-submit cookie for state-mutating requests.

---

## 10. Naming Conventions

### 10.1 Go

- **Packages**: lowercase, single word (`domain`, `application`, `infrastructure`)
- **Files**: lowercase with underscores (`order_repository.go`)
- **Types**: CamelCase (`OrderRepository`, `PgOrderRepository`)
- **Interfaces**: often end in `-er` (`Reader`, `Writer`) or descriptive (`OrderRepository`)
- **Functions**: CamelCase (`CreateOrder`, `FindByID`)
- **Variables**: CamelCase; short names in small scope (`db`, `ctx`)

### 10.2 TypeScript / React

- **Files**: PascalCase for components (`RestaurantCard.tsx`), camelCase for utilities (`formatPrice.ts`)
- **Components**: PascalCase (`RestaurantCard`, `OrderDetail`)
- **Functions**: camelCase (`formatPrice`, `calculateTotal`)
- **Variables**: camelCase
- **Constants**: UPPER_SNAKE (`MAX_RETRY_ATTEMPTS`)
- **Types/Interfaces**: PascalCase (`Order`, `Restaurant`)

### 10.3 Database

- **Tables**: snake_case, plural (`orders`, `order_items`)
- **Columns**: snake_case (`created_at`, `customer_id`)
- **Indexes**: `idx_{table}_{columns}` (`idx_orders_customer_id`)
- **Foreign keys**: `{table}_{column}_fk` (`orders_customer_id_fk`)

---

> **Next**: Read `REPO-STRUCTURE.md` for the monorepo layout.
