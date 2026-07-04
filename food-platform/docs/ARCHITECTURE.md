# Architecture Document — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active — Source of Truth  
> **Last Updated**: 2026-07-04  
> **Owner**: Engineering Team

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Business Context](#2-business-context)
3. [High-Level Architecture](#3-high-level-architecture)
4. [Domain-Oriented Microservices](#4-domain-oriented-microservices)
5. [Communication Patterns](#5-communication-patterns)
6. [Data Architecture](#6-data-architecture)
7. [Authentication & Authorization](#7-authentication--authorization)
8. [Real-Time Infrastructure](#8-real-time-infrastructure)
9. [Security & Anti-Fraud](#9-security--anti-fraud)
10. [Observability](#10-observability)
11. [DevOps & Infrastructure](#11-devops--infrastructure)
12. [Scaling Strategy](#12-scaling-strategy)
13. [Decision Records Summary](#13-decision-records-summary)
14. [Research References](#14-research-references)

---

## 1. Executive Summary

This document describes the architecture of a **production-grade food delivery platform** targeting the Egyptian market. The platform is engineered to handle **50,000+ orders/day** at full scale, serving **7 distinct user types** through dedicated web applications, all backed by a **Go-powered event-driven backend** with **PostgreSQL as the source of truth**.

### 1.1 Key Architectural Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Backend language | Go 1.22 | Best concurrency model for I/O-heavy services; fast compile times; small binary size; strong standard library for HTTP/gRPC. Reference: Uber Engineering ([1]), Go benchmarks ([2]). |
| Frontend framework | React 18 + Vite (Web) | Largest talent pool in Egypt; ecosystem maturity; component reuse across all 7 apps. |
| Primary database | PostgreSQL 16 + PostGIS | ACID guarantees for orders/payments; JSONB for flexible metadata; PostGIS for geospatial. Reference: Uber Engineering ([3]). |
| Event bus | Apache Kafka | Industry standard for event-driven; replay capability; exactly-once for payments. Reference: DoorDash Engineering ([4]). |
| Container orchestration | Kubernetes (EKS) | Portability; rich ecosystem; horizontal pod autoscaling for peak hours. |
| Monorepo | Turborepo + pnpm | Shared packages between 7 web apps; atomic commits across packages. |
| API design | Contract-First (OpenAPI 3.1) | Type generation for both Go and TypeScript; documentation as code. |
| Mobile strategy | Web-first (Phase 1) → React Native (Phase 2) | Faster iteration for MVP; 70% code reuse when migrating to mobile. |

### 1.2 Scale Targets

| Metric | Year 1 Target |
|--------|---------------|
| Daily orders | 50,000 |
| Active restaurants | 2,500+ |
| Active drivers | 15,000+ |
| Concurrent users | 50,000 |
| p95 API latency | <500ms |
| Uptime (Tier 0) | 99.95% |
| Order completion rate | >92% |

---

## 2. Business Context

### 2.1 Market

- **Target market**: Egypt (Cairo → Giza → Alexandria → Mansoura)
- **Market size**: $1.7B by 2027, 12% YoY growth
- **Competitors**: Talabat (22-25% commission), Uber Eats (cards only), elmenus (no own fleet)
- **Differentiator**: 12-15% commission (SME-friendly) + 75% of commission to drivers (vs 55% industry) + native anti-fraud layer

### 2.2 User Types (7)

| # | User | App | Platform | Description |
|---|------|-----|----------|-------------|
| 1 | Customer | customer-web | Web (later RN) | Browse, order, track, pay |
| 2 | Driver | driver-web | Web (later RN) | Accept orders, navigate, earn |
| 3 | Restaurant | restaurant-web | Web | Receive orders, manage menu |
| 4 | Support Agent | support-web | Web | Tickets, refunds, escalation |
| 5 | Ops Manager | command-center | Web | Live ops dashboard |
| 6 | Internal Employee | employee-portal | Web | Onboarding, biometric, audit |
| 7 | Field Supervisor | field-supervisor-web | Web (later RN) | On-site verification |

### 2.3 Payment Methods

| Method | Market Share (Target) | Provider |
|--------|----------------------|----------|
| Vodafone Cash | 45% | Vodafone API |
| InstaPay | 25% | InstaPay API |
| Card (Visa/Master) | 20% | Paymob |
| Cash on Delivery | 10% | N/A (driver-collected) |

---

## 3. High-Level Architecture

### 3.1 Layered Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│  Layer 1: Client Layer (7 React Web Apps)                       │
│  customer | driver | restaurant | support | ops | portal | field│
└────────────────────────────┬────────────────────────────────────┘
                             │ HTTPS / WSS
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  Layer 2: Edge Layer                                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐           │
│  │ Kong API     │  │ WebSocket    │  │ Cloudflare   │           │
│  │ Gateway      │  │ Gateway      │  │ CDN + WAF    │           │
│  │ (JWT, RL,    │  │ (real-time)  │  │              │           │
│  │  routing)    │  │              │  │              │           │
│  └──────────────┘  └──────────────┘  └──────────────┘           │
└────────────────────────────┬────────────────────────────────────┘
                             │ gRPC / HTTP
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│  Layer 3: Service Layer (12 Go Microservices)                  │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
│  │  Auth   │ │Catalog  │ │  Menu   │ │  Order  │ │Payment  │    │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘    │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐    │
│  │Matching │ │ Driver  │ │  Geo    │ │ Notif   │ │  Fraud  │    │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘    │
│  ┌─────────┐ ┌─────────┐                                        │
│  │  Promo  │ │Analytics│                                        │
│  └─────────┘ └─────────┘                                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
┌──────────────────┐ ┌──────────────┐ ┌──────────────┐
│  PostgreSQL 16   │ │  Redis 7     │ │  Kafka 3.x   │
│  + PostGIS       │ │  (cache,     │ │  (event bus) │
│  (12 databases)  │ │   geo, sess) │ │              │
└──────────────────┘ └──────────────┘ └──────────────┘
              ▲              ▲              ▲
              │              │              │
┌──────────────────┐ ┌──────────────┐ ┌──────────────┐
│  ElasticSearch   │ │  ClickHouse  │ │  S3 / MinIO  │
│  (search)        │ │  (analytics) │ │  (files)     │
└──────────────────┘ └──────────────┘ └──────────────┘
```

### 3.2 Request Flow Example (Order Creation)

```
1. Customer taps "Place Order" in customer-web
2. customer-web → POST /api/v1/orders (HTTPS)
3. Kong API Gateway:
   - Validates JWT
   - Rate limits (100 RPS/user)
   - Routes to Order Service
4. Order Service (Go):
   - Validates cart + pricing
   - Calls Fraud Service (gRPC, sync, <100ms)
   - Calls Payment Service (gRPC, sync, idempotent)
   - Creates Order (status: PENDING)
   - Publishes OrderCreated event to Kafka
5. Kafka fans out to:
   - Notification Service (push to restaurant)
   - Analytics Service (record)
   - Loyalty Service (points)
6. Restaurant accepts via restaurant-web
7. Order Service publishes OrderConfirmed
8. Delivery Matching Service consumes:
   - Queries Redis GEO for nearby drivers
   - Computes scores
   - Pushes to top-3 drivers (via WebSocket)
9. Driver accepts via driver-web
10. Order Service publishes OrderAssigned
11. Customer sees driver assigned (via WebSocket push)
```

**Total p95 latency target**: <2s for steps 1-5; <10s for steps 6-10.

---

## 4. Domain-Oriented Microservices

Inspired by Uber's **Domain-Oriented Microservice Architecture** ([1]). Each service owns its data, contracts, and scaling policy.

### 4.1 Service Catalog

| # | Service | Responsibility | Database | Tech | Throughput |
|---|---------|----------------|----------|------|------------|
| 1 | **Auth** | JWT, OTP, sessions, WebAuthn | `auth_db` (PG) | Go + Redis | 500 RPS |
| 2 | **Restaurant Catalog** | Restaurant master data, hours, service areas | `restaurants_db` (PG + PostGIS) | Go | 2K RPS read |
| 3 | **Menu** | Items, pricing, modifiers, availability | `menus_db` (PG) | Go | 3K RPS read |
| 4 | **Order** ⭐ | Order lifecycle, state machine, validation | `orders_db` (PG, partitioned) | Go + Kafka | 1K RPS write |
| 5 | **Payment** | Vodafone Cash, InstaPay, card, refunds | `payments_db` (PG) | Go + Vault | 500 RPS |
| 6 | **Delivery Matching** ⭐ | Driver-customer matching, dispatch | `dispatch_db` (PG + Redis GEO) | Go | 200 RPS |
| 7 | **Driver Management** | Profiles, KYC, vehicle, earnings | `drivers_db` (PG) | Go | 300 RPS |
| 8 | **Geo/Tracking** | Real-time GPS ingestion, ETA | Redis only | Go + Redis GEO | 50K writes/sec |
| 9 | **Notification** | Push, SMS, email, in-app | `notifications_db` (PG) | Go + Kafka | 1K RPS |
| 10 | **Fraud Detection** | ML scoring, rule engine, review queue | `fraud_db` (PG) | Go + Python ML | 500 RPS sync |
| 11 | **Promo/Loyalty** | Coupons, campaigns, points, cashback | `promos_db` (PG) | Go | 200 RPS |
| 12 | **Analytics** | Aggregations, dashboards, exports | ClickHouse | Go + Flink | Async |

### 4.2 Service Boundaries (DDD)

Each service follows **Domain-Driven Design** principles:

```
services/order/
├── cmd/
│   └── server/main.go           # Entry point
├── internal/
│   ├── domain/                  # Entities, value objects, domain events
│   │   ├── order.go
│   │   ├── status.go
│   │   └── events.go
│   ├── application/             # Use cases, application services
│   │   ├── create_order.go
│   │   ├── cancel_order.go
│   │   └── update_status.go
│   ├── infrastructure/          # DB, Kafka, external APIs
│   │   ├── postgres/
│   │   ├── kafka/
│   │   └── grpc/
│   └── interfaces/              # HTTP, gRPC handlers
│       ├── http/
│       └── grpc/
├── migrations/
├── proto/                       # gRPC definitions
├── openapi/                     # REST definitions
├── tests/
├── Dockerfile
├── go.mod
└── README.md
```

### 4.3 Database-per-Service

Each of the 12 microservices owns its own PostgreSQL logical database. **No service reads another's database directly** — only via API or event. This enforces domain boundaries and lets each service optimize its schema independently.

Cross-service joins are replaced by:
- **API composition** (for reads): Order Service calls Restaurant Catalog API to enrich order with restaurant name.
- **Event-carried state transfer** (for writes): When Menu Service updates a price, it publishes `MenuPriceUpdated`; Order Service consumes and caches the new price.

Reference: Microservices.io — Database per service pattern ([5]).

---

## 5. Communication Patterns

### 5.1 Three Patterns (Picked Per Use Case)

| Pattern | When | Protocol | Latency | Failure Mode |
|---------|------|----------|---------|--------------|
| **Synchronous** | Caller NEEDS result to proceed | gRPC | <100ms p99 | Caller fails fast |
| **Asynchronous** | Caller announces something happened | Kafka | ~200ms | Consumer lag (alert if >1000) |
| **Real-time** | Clients need push | WebSocket | <50ms | Client auto-reconnects |

### 5.2 Synchronous (gRPC)

**Use cases**:
- Order Service → Fraud Service (need risk score before confirming)
- Order Service → Payment Service (need confirmation before notifying restaurant)
- Customer Web → API Gateway → any read query

**Properties**:
- Protocol Buffers for typed contracts
- Circuit breakers via `go-resiliency`
- Retries with jitter (max 3)
- Timeout at 2s
- gRPC interceptors for auth, logging, tracing

**Decision rule**: If the caller NEEDS the result to proceed → gRPC.

### 5.3 Asynchronous (Kafka)

**Use cases**:
- `OrderCreated` → Notification, Analytics, Loyalty (fan-out)
- `DriverLocation` → Geo Service, Command Center (high-volume stream)
- `MenuUpdated` → Catalog, Search (eventual consistency OK)

**Properties**:
- Throughput: 50K+ msgs/sec
- Retention: 7-30 days
- Avro schemas (Schema Registry)
- Exactly-once semantics for payments
- Consumer group per service

**Core Topics**:

| Topic | Producer | Consumers |
|-------|----------|-----------|
| `order.created` | Order | Fraud, Notification, Analytics, Loyalty |
| `order.confirmed` | Delivery Matching | Notification, Restaurant, Customer |
| `order.cancelled` | Order, Support | Payment (refund), Driver (unassign) |
| `payment.captured` | Payment | Order, Loyalty, Analytics |
| `payment.failed` | Payment | Order, Notification |
| `driver.location` | Driver App | Geo Service, Command Center (streaming) |
| `driver.status_changed` | Driver Mgmt | Delivery Matching, Analytics |
| `restaurant.menu_updated` | Restaurant App | Catalog, Search |

**Decision rule**: If the caller just announces something happened → Kafka.

### 5.4 Real-time (WebSocket)

**Use cases**:
- Customer Web — order tracking (driver location every 5s)
- Driver Web — new order pushes (sub-second)
- Command Center — live ops dashboard
- Support Web — live chat

**Properties**:
- WS Gateway: standalone Go service
- 100K concurrent connections per pod
- Horizontal scaling via Redis pub/sub
- Client auto-reconnect with exponential backoff
- Missed events fetched via REST fallback

**Decision rule**: If clients need push → WebSocket.

---

## 6. Data Architecture

### 6.1 Storage Tiers

```
┌─────────────────────────────────────────────────┐
│ Tier 1: Primary OLTP (PostgreSQL 16)            │
│  • 12 logical databases (one per service)       │
│  • 1 primary + 2 read replicas per critical svc │
│  • Failover via Patroni                         │
│  • Range partitioning for orders (by month)     │
└─────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────┐
│ Tier 2: Geospatial (PostgreSQL + PostGIS)       │
│  • Restaurant Catalog + Delivery Matching       │
│  • For sub-ms queries: Redis GEO (GEOADD)       │
│  • Driver location stream: Redis only           │
└─────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────┐
│ Tier 3: Cache (Redis 7 cluster)                 │
│  • 6 nodes (3 shards + replicas)                │
│  • Sessions, hot menu data, driver availability │
│  • TTL-based eviction + event-driven invalidate │
└─────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────┐
│ Tier 4: Search & Analytics                      │
│  • ElasticSearch 8 (restaurant/menu search)     │
│  • ClickHouse (analytics, 100M+ rows/day)       │
│  • S3/MinIO (images, documents, audit logs)     │
└─────────────────────────────────────────────────┘
```

### 6.2 Database-per-Service Pattern

Reference: Microservices.io ([5]).

**Benefits**:
- Enforces domain boundaries
- Each service can optimize its schema independently
- No cascading failures from shared DB
- Independent scaling

**Trade-offs**:
- Cross-service joins require API composition
- Distributed transactions → use Saga pattern ([6])
- Eventual consistency for cross-service queries

### 6.3 Core Entities (12)

| Entity | Service Owner | Key Fields |
|--------|---------------|------------|
| `users` | Auth | id, phone, email, name, status |
| `restaurants` | Catalog | id, name, slug, location (POINT), commission_rate |
| `menu_items` | Menu | id, restaurant_id, name, price, is_available |
| `orders` ⭐ | Order | id, user_id, restaurant_id, driver_id, status, total_amount, ETA |
| `order_items` | Order | id, order_id, menu_item_id, quantity, unit_price, modifiers |
| `payments` | Payment | id, order_id, provider, amount, status, provider_txn_id |
| `drivers` | Driver Mgmt | id, user_id, vehicle_type, license_number, rating, total_earnings |
| `deliveries` | Delivery Matching | id, order_id, driver_id, pickup_at, delivered_at, distance_km, earnings |
| `driver_locations` (partitioned) | Geo | driver_id, location (POINT), heading, speed, recorded_at |
| `promos` | Promo | id, code, type, value, valid_from, valid_to, usage_limit |
| `audit_logs` (immutable, hash-chained) | Trust | actor_id, action, entity_type, metadata, prev_hash, record_hash |
| `fraud_scores` | Fraud | order_id, score (0-100), reasons, model_version |

### 6.4 Partitioning Strategy

- **orders** table: Range partitioning by `created_at` month. Each partition holds ~1.5M rows at peak. Old partitions (>12 months) moved to cold storage (S3 + Parquet).
- **driver_locations**: Time-series, archived to S3 after 7 days.
- **audit_logs**: Append-only, monthly partitions, 7-year retention (compliance).

### 6.5 Sharding Strategy

- **Stage 1 (MVP, 500 orders/day)**: Single PostgreSQL instance, no sharding.
- **Stage 2 (5K orders/day)**: Critical services (Order, Payment, Driver) on dedicated clusters; read replicas added.
- **Stage 3 (50K orders/day)**: Orders sharded by city (Cairo/Giza/Alex shards); driver_locations moved to Redis-only + S3 archive.

**Decision rule**: Shard only when a single primary can't keep up with write throughput (>5K writes/sec). Premature sharding adds complexity without benefit.

---

## 7. Authentication & Authorization

### 7.1 Auth Stack (3 User Categories)

#### Customer/Driver Auth
- Phone + OTP (SMS via Twilio/local provider)
- Optional password/biometric on device
- JWT access token (15 min TTL) + refresh token (30 days)
- Refresh rotation — old token invalidated on use
- Device binding (refresh token tied to device fingerprint)

#### Restaurant Auth
- Phone + OTP for primary account
- Sub-accounts (manager, cashier) with role scoping
- API keys for POS integration (rotated quarterly)
- IP allowlist (configurable per restaurant)

#### Internal Employee Auth (Most Strict)
- Username + password + TOTP (mandatory 2FA)
- **Biometric (WebAuthn)** for sensitive actions: refunds >EGP 200, restaurant approval, driver payout, menu overrides
- Hardware security key (YubiKey) for Command Center admins
- Session timeout: 15 min idle, 8h max

Reference: FIDO Alliance WebAuthn spec ([7]), Auth0 MFA implementation ([8]).

### 7.2 JWT Structure

```json
{
  "sub": "uuid-user-id",
  "iss": "food-platform-auth",
  "aud": "food-platform",
  "iat": 1720000000,
  "exp": 1720000900,  // 15 min
  "role": "customer",  // customer | driver | restaurant | support_l1 | support_l2 | ops_manager | finance | super_admin | field_supervisor | hr
  "scope": "orders:read orders:create",
  "device_id": "uuid-device",
  "session_id": "uuid-session"
}
```

### 7.3 RBAC Matrix (Summary)

| Action | Super Admin | Ops Mgr | Finance | Support L2 | Support L1 | Field Sup | HR |
|--------|-------------|---------|---------|------------|------------|-----------|-----|
| Refund <EGP 100 | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ |
| Refund 100-500 | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| Refund >EGP 500 | ✅ | ✅ | ✅* | ❌ | ❌ | ❌ | ❌ |
| Restaurant approve | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| Driver approve | ✅ | ✅ | ❌ | ❌ | ❌ | ✅* | ❌ |
| Driver payout | ✅ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| View customer PII | ✅ | ✅ | partial | ✅ | partial | partial | ❌ |
| Manage employees | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ |

`✅*` = requires dual approval

### 7.4 Authorization Flow

```
1. Request arrives at API Gateway
2. Gateway extracts JWT from Authorization header
3. Gateway validates JWT signature (RS256, public key from Auth Service)
4. Gateway injects user_id, role, scope into upstream headers
5. Backend service:
   - Reads role from header
   - Checks permission via Authorization Service (CASL-like)
   - All decisions logged to audit_logs
```

---

## 8. Real-Time Infrastructure

### 8.1 WebSocket Gateway

Standalone Go service for persistent connections.

```
┌─────────────┐    WSS    ┌─────────────┐    Redis Pub/Sub    ┌─────────────┐
│ Client App  │ ←──────→ │ WS Gateway  │ ←─────────────────→ │ Backend Svc │
└─────────────┘           │ (Go)        │                     └─────────────┘
                          │ 100K conn/  │
                          │ pod         │
                          └─────────────┘
```

**Properties**:
- 100K concurrent connections per pod (goroutines)
- Horizontal scaling via Redis pub/sub (cross-pod message fan-out)
- Heartbeat every 30s
- Auto-reconnect with exponential backoff (1s, 2s, 4s, ..., max 30s)
- Missed events fetched via REST fallback (`GET /events?since=timestamp`)

### 8.2 Driver Location Streaming

```
Driver Web (every 5s)
  ↓ GPS update
WebSocket Gateway
  ↓ publish to Kafka "driver.location"
Geo Service
  ├── Redis GEO (GEOADD for matching)
  └── Redis pub/sub "driver.{driver_id}.location"
      ↓
WebSocket Gateway
  ↓ push to Customer Web (only the customer watching this order)
Customer Web
  → map update
```

**Why Redis pub/sub, not Kafka, for the last hop?**
- Redis pub/sub latency: <50ms (excellent for real-time UI)
- Kafka latency: ~200ms (for general events)
- Not every GPS update needs to be a Kafka event

### 8.3 Order Status Notifications

Every order status change → push notification to relevant clients:

```
order.status_changed event (Kafka)
  ↓
Notification Service
  ├── WebSocket push to Customer Web
  ├── WebSocket push to Driver Web
  ├── WebSocket push to Restaurant Web
  ├── Push notification (FCM for Android, APNs for iOS — when web app not focused)
  └── SMS for critical status (delivered, cancelled)
```

---

## 9. Security & Anti-Fraud

### 9.1 Three-Layer Fraud Defense

Reference: Delivery Hero anti-fraud ([9]), Sift Science food delivery ([10]), Incognia fraud detection ([11]).

#### Layer 1: Customer Fraud
**Threats**: Fake accounts (promo abuse), stolen credit cards, friendly fraud (claim non-delivery), COD abuse.

**Detection**:
- Device fingerprinting (FingerprintJS)
- Phone number age + verification
- Order velocity per device/IP/household
- ML model trained on historical chargebacks
- Trust score (0-100) per user, recalculated each order

**Response**: Block auto-refund for users with trust score <30; require OTP for high-value COD orders.

#### Layer 2: Driver Fraud
**Threats**: Fake deliveries, collusion with restaurants, GPS spoofing, multi-accounting, long-hauling.

**Detection**:
- GPS trajectory validation (compare actual vs optimal route)
- Customer delivery confirmation (photo + OTP)
- Restaurant pickup confirmation scan
- Cross-driver restaurant frequency analysis
- Driver earnings anomaly detection

**Response**: Auto-suspend drivers with >3 flagged orders in 7 days.

#### Layer 3: Internal Employee Fraud ⭐ (Hardest)
**Threats**: Manual refunds for friends/family, fake restaurant approvals for kickbacks, customer data leaks, order manipulation.

**Detection**:
- **Biometric auth** (WebAuthn) on every portal login + sensitive action
- **Immutable audit log** (append-only, hash-chained)
- Anomaly detection on refund volume per employee
- **Segregation of duties** (refund >EGP 500 requires 2 people)
- Quarterly access review

**Response**: Immediate account suspension; legal action per Egyptian Cybercrime Law 175/2018 ([12]).

### 9.2 Immutable Audit Log (Hash-Chained)

Reference: Emergent Mind immutable audit log ([13]), cryptographic audit trail patterns ([14]).

```sql
CREATE TABLE audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    actor_id        UUID NOT NULL,
    actor_type      VARCHAR(50) NOT NULL,
    actor_email     TEXT NOT NULL,
    actor_role      VARCHAR(50) NOT NULL,
    action          TEXT NOT NULL,
    action_category VARCHAR(50) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       UUID,
    metadata        JSONB NOT NULL,
    biometric_verified BOOLEAN NOT NULL DEFAULT false,
    dual_approval      BOOLEAN NOT NULL DEFAULT false,
    approver_id        UUID,
    ip_address      INET NOT NULL,
    user_agent      TEXT NOT NULL,
    session_id      UUID NOT NULL,
    prev_hash       BYTEA NOT NULL,
    record_hash     BYTEA NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

REVOKE UPDATE, DELETE ON audit_logs FROM PUBLIC;
```

**Hash chain**:
```
record_hash = SHA256(actor_id + action + metadata + timestamp + prev_hash)
```

**Tamper detection**: Nightly job walks the chain. If any record's hash doesn't match → P0 incident, system lockdown.

### 9.3 Compliance

#### Egyptian Law 175/2018 (Anti-Cybercrime)
- Article 7: Penalizes unauthorized system access
- Article 11: Penalizes data interception
- Article 12: Penalizes modification/destruction
- Article 16: Service provider must protect data
- Article 25: Obligation to report crimes

**Compliance checklist**:
- ✅ 7-year retention for financial records (tax law)
- ✅ 90-day hot retention for audit logs
- ✅ 7-year cold retention for audit logs
- ✅ AES-256 encryption at rest
- ✅ TLS 1.3 in transit
- ✅ PII column-level encryption
- ✅ NTRA breach notification within 72h
- ✅ Maintain evidence chain

Reference: Andersen Egypt Law 175/2018 translation ([12]), ID Egypt legal overview ([15]).

#### PCI-DSS (for card payments)
- Never store full PAN (card number) — only last 4 digits + token
- Payment tokenization via Paymob
- All card input fields in provider's iframe (out of our scope)
- Quarterly ASV scans, annual SAQ-A submission
- Network segmentation: payment subnet isolated

---

## 10. Observability

### 10.1 Three Pillars

| Pillar | Tool | Retention |
|--------|------|-----------|
| Metrics | Prometheus + Grafana | 90 days |
| Logs | Loki (structured JSON) | 30 days |
| Traces | Jaeger + OpenTelemetry | 7 days (10% sampling) |
| Errors | Sentry (client + backend) | 90 days |

### 10.2 Service Tiers & SLOs

| Tier | Services | Uptime SLO | Latency p95 | Error Budget/Month |
|------|----------|-----------|-------------|-------------------|
| Tier 0 (Critical) | Order, Payment, Auth | 99.95% | <500ms | 22 min |
| Tier 1 (Important) | Catalog, Menu, Matching, Geo, Notification | 99.9% | <1s | 43 min |
| Tier 2 (Standard) | Driver Mgmt, Promo, Fraud, Analytics | 99.5% | <2s | 3.6h |

### 10.3 Golden Signals (per service)

1. **Latency**: p50, p95, p99 per endpoint
2. **Traffic**: RPS, concurrent connections, Kafka msgs/sec
3. **Errors**: HTTP 5xx rate, Kafka consumer errors, client errors
4. **Saturation**: CPU, memory, disk, network, DB connection pool, Kafka consumer lag

### 10.4 Alerting

- P0/P1 → PagerDuty (pages on-call)
- P2/P3 → Slack notifications
- All alerts include runbook link in description

---

## 11. DevOps & Infrastructure

### 11.1 Container Orchestration (Kubernetes)

```
EKS Cluster (AWS)
├── Node Group: backend (Go services)
│   ├── m5.large × 3 (min) → m5.2xlarge × 20 (max)
├── Node Group: web (React apps served via Nginx)
│   ├── t3.medium × 2
├── Node Group: data (PostgreSQL, Redis, Kafka)
│   ├── r5.xlarge × 3 (stateful, pinned)
└── Node Group: spot (non-critical workloads)
    ├── m5.large × 5 (analytics, batch jobs)
```

**Autoscaling**:
- HPA on CPU >70% AND RPS >threshold
- KEDA for Kafka-driven scaling (lag-based)
- Cluster autoscaler for node-level scaling
- Pod anti-affinity for HA across AZs

### 11.2 CI/CD Pipeline

```
push → lint → test → build → scan (Trivy) → deploy to staging → smoke test → promote to prod (manual approval)
```

- **CI**: GitHub Actions
- **CD**: ArgoCD (GitOps — cluster state declared in Git)
- **Helm charts** per service (versioned, reusable)
- **Image scanning**: Trivy (block on CRITICAL CVEs)
- **Deploy preview env** per PR (ephemeral, auto-destroy after 24h)

### 11.3 Infrastructure as Code

- **Terraform** for AWS resources (VPC, EKS, RDS, S3, etc.)
- **Terragrunt** for multi-env (dev/staging/prod) DRY
- **State** in S3 + DynamoDB lock
- **Drift detection** nightly

**Policy**: "If it's not in Terraform, it doesn't exist."

### 11.4 Secrets Management

- AWS Secrets Manager for production secrets
- External Secrets Operator syncs to K8s secrets
- Sealed Secrets for GitOps-friendly sensitive values
- Rotation: DB credentials every 90 days, API keys quarterly
- No secrets in code, env files, or chat

### 11.5 Disaster Recovery

- **RPO**: 5 minutes (Postgres WAL streaming + 5-min snapshots)
- **RTO**: 30 minutes (multi-AZ failover, tested quarterly)
- **Cross-region replication** for critical data (us-east-1 → us-west-2)
- **Quarterly DR drill**
- **Chaos engineering** game days every quarter

---

## 12. Scaling Strategy

### 12.1 Peak Profile

Lunch (12-3pm) and dinner (7-11pm) carry **70% of daily volume**. Provisioning for peak = provisioning for ~3x average.

### 12.2 Five Scaling Tactics

1. **Horizontal Pod Autoscaling (HPA)**: CPU >70% AND RPS >threshold; min 3 replicas (HA), max 50 per service.
2. **KEDA Event-Driven Autoscaling**: Kafka lag-based for Notification, Analytics. If lag >1000 → add pods.
3. **Read Replicas**: 80% read traffic to replicas, 20% to primary. PgBouncer pools connections.
4. **Caching Strategy**: Redis for hot reads (restaurants, menus, promos). TTL: 60s/300s/30s. Event-driven invalidation. Hit rate target: >85%.
5. **Pre-Scaled Capacity**: Pre-scale Order + Payment services 15 min before peak (cron-triggered, based on historical data + weather).

---

## 13. Decision Records Summary

Detailed ADRs in `/docs/adr/`. Summary:

| ADR | Decision | Status |
|-----|----------|--------|
| 001 | Go 1.22 for backend | Accepted |
| 002 | PostgreSQL 16 + PostGIS for primary DB | Accepted |
| 003 | Apache Kafka for event bus | Accepted |
| 004 | Kubernetes (EKS) for orchestration | Accepted |
| 005 | React 18 + Vite for all web apps | Accepted |
| 006 | Turborepo + pnpm for monorepo | Accepted |
| 007 | Contract-First (OpenAPI 3.1) | Accepted |
| 008 | WebAuthn for employee biometric auth | Accepted |
| 009 | Hash-chained audit log (not blockchain) | Accepted |
| 010 | Web-first strategy (mobile later via React Native) | Accepted |

---

## 14. Research References

All references are detailed in `/docs/REFERENCES.md`. Key sources:

- Uber Engineering Blog — Domain-Oriented Microservice Architecture ([1])
- Go official benchmarks ([2])
- Uber Engineering — Scaling Uber Eats ([3])
- DoorDash Engineering — Kafka + Flink ([4])
- Microservices.io — Database per service ([5])
- Microservices.io — Saga pattern ([6])
- FIDO Alliance — WebAuthn / FIDO2 ([7])
- Auth0 — MFA with WebAuthn ([8])
- Delivery Hero — Real-time anti-fraud system ([9])
- Sift Science — Food delivery fraud ([10])
- Incognia — Food delivery fraud prevention ([11])
- Andersen Egypt — Law 175/2018 translation ([12])
- Emergent Mind — Immutable audit log ([13])
- GitHub — Cryptographic audit trail ([14])
- ID Egypt — Cybersecurity laws overview ([15])
- Uber Engineering — H3 hexagonal spatial index ([16])
- Uber Engineering — Risk Entity Watch ([17])
- Radar — Real-time delivery tracking ([18])
- Talabat Partner Portal — Direct observation ([19])
- Zendesk — Support tier levels ([20])
- PagerDuty — Incident severity classification ([21])

---

> **Next**: Read `API-CONTRACTS.md` for the detailed REST, WebSocket, and Kafka event specifications.
