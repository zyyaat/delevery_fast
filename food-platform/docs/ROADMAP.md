# 12-Week Roadmap — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Total Duration**: 12 weeks (3 months)  
> **Methodology**: Contract-First + Vertical Slices + Incremental

---

## Table of Contents

1. [Overview](#1-overview)
2. [Phase 1: Foundation (Weeks 1-2)](#2-phase-1-foundation-weeks-1-2)
3. [Phase 2: Customer Flow (Weeks 3-5)](#3-phase-2-customer-flow-weeks-3-5)
4. [Phase 3: Restaurant Flow (Weeks 6-7)](#4-phase-3-restaurant-flow-weeks-6-7)
5. [Phase 4: Driver Flow (Weeks 8-9)](#5-phase-4-driver-flow-weeks-8-9)
6. [Phase 5: Support & Ops (Weeks 10-11)](#6-phase-5-support--ops-weeks-10-11)
7. [Phase 6: Trust & Field (Week 12)](#7-phase-6-trust--field-week-12)
8. [Milestones & Exit Criteria](#8-milestones--exit-criteria)
9. [Risk Register](#9-risk-register)
10. [Resource Requirements](#10-resource-requirements)

---

## 1. Overview

### 1.1 Methodology

```
1. Contract-First: OpenAPI/AsyncAPI specs before any code
2. Vertical Slices: Each slice delivers end-to-end functionality
3. Incremental: Every week produces a demo-able increment
4. Test-Driven: 80%+ coverage on business logic
5. Documentation as Code: ARCHITECTURE.md stays in sync
```

### 1.2 Weekly Rhythm

```
Day 1:    Planning + specs review
Day 2-4:  Implementation (pair programming with AI)
Day 5:    Code review + tests + demo
Day 6-7:  Bug fixes + documentation
```

### 1.3 Phase Summary

| Phase | Weeks | Goal | Sessions | Deliverable |
|-------|-------|------|----------|-------------|
| 1. Foundation | 1-2 | Infra + contracts + auth | 3 | Dev env ready, login works |
| 2. Customer | 3-5 | Browse + cart + checkout + tracking | 3 | Customer can place order |
| 3. Restaurant | 6-7 | Receive + accept + KDS + menu | 2 | Restaurant can serve orders |
| 4. Driver | 8-9 | Online + accept + deliver + earn | 2 | Driver can deliver orders |
| 5. Support & Ops | 10-11 | Tickets + command center | 2 | Support can resolve issues |
| 6. Trust & Field | 12 | Employee portal + field supervisor | 1 | Full system ready |

**Total: 13 sessions over 12 weeks**

---

## 2. Phase 1: Foundation (Weeks 1-2)

### Goal
Establish dev environment, write contracts, build auth foundation.

### Session 1 — Documentation & Setup (Week 1, Day 1-2)

**Tasks**:
- [x] Write `ARCHITECTURE.md` (DONE)
- [x] Write `API-CONTRACTS.md` (DONE)
- [x] Write `PATTERNS.md` (DONE)
- [x] Write `REPO-STRUCTURE.md` (DONE)
- [x] Write `ROADMAP.md` (this file)
- [ ] Write `REFERENCES.md`
- [ ] Write `SESSIONS-LOG.md`
- [ ] Initialize monorepo (package.json, turbo.json, pnpm-workspace.yaml)
- [ ] Initialize Go module for each service
- [ ] Create 7 web app skeletons (Vite + React + TS)
- [ ] Create 7 shared packages skeletons
- [ ] Write `README.md` for root + each app/package
- [ ] Setup `.gitignore`, `.editorconfig`, `prettierrc`
- [ ] Setup ESLint shared config

**Deliverable**: Monorepo skeleton committed to Git.

### Session 2 — Dev Infrastructure (Week 1, Day 3-5)

**Tasks**:
- [ ] Write `infra/docker-compose.yml` (PostgreSQL, Redis, Kafka, ES, ClickHouse, MinIO)
- [ ] Write `infra/init-db.sql` (12 databases)
- [ ] Test: `docker compose up -d` works
- [ ] Setup GitHub Actions CI skeleton
  - [ ] Frontend pipeline (lint, type-check, test, build)
  - [ ] Backend pipeline (per-service: lint, test, build)
- [ ] Setup pre-commit hooks (Husky + lint-staged)
- [ ] Write `Makefile` (dev, build, test, lint, codegen, migrate, docker-up)

**Deliverable**: `make setup` brings up full dev env in <5 min.

### Session 3 — Auth Service + API Gateway (Week 2)

**Tasks**:
- [ ] **Auth Service (Go)** — full implementation:
  - [ ] Domain: User, Session, RefreshToken entities
  - [ ] Application: LoginUseCase, VerifyOTPUseCase, RefreshTokenUseCase
  - [ ] Infrastructure: PostgreSQL repository, Redis session store, Twilio SMS
  - [ ] Interfaces: HTTP handlers (send-otp, verify-otp, refresh, logout)
  - [ ] WebSocket support for employee auth events
  - [ ] WebAuthn registration + verification (for employees)
  - [ ] Migrations: users, sessions, refresh_tokens, webauthn_credentials
  - [ ] Unit tests (80%+ coverage)
  - [ ] Integration tests (auth flow end-to-end)
  - [ ] OpenAPI spec for `/auth/*` endpoints
  - [ ] Dockerfile + Helm chart
- [ ] **API Gateway**:
  - [ ] Kong or APISIX configuration
  - [ ] JWT validation plugin (RS256, JWKS from Auth Service)
  - [ ] Rate limiting plugin (per-user, per-IP)
  - [ ] Routing to all 12 services
  - [ ] CORS configuration
  - [ ] Health check endpoint
- [ ] **Shared Go library** (`services/shared/`):
  - [ ] Error types + HTTP error mapping
  - [ ] Structured logging (slog)
  - [ ] OpenTelemetry tracing setup
  - [ ] HTTP middleware (auth, request_id, logging, recovery)
  - [ ] Kafka producer/consumer helpers
  - [ ] PostgreSQL connection pool helper
- [ ] **Shared TypeScript packages**:
  - [ ] `@food-platform/types`: Generated from OpenAPI
  - [ ] `@food-platform/api-client`: Axios + interceptors
  - [ ] `@food-platform/auth`: Zustand store + token refresh
  - [ ] `@food-platform/ui`: Button, Card, Input, Spinner, Toast
  - [ ] `@food-platform/theme`: Design tokens + Tailwind config

**Deliverable**: 
- User can `POST /auth/otp/send` + `POST /auth/otp/verify` and get JWT.
- All 7 web apps can login (with mock role switching).

**Exit Criteria**:
- ✅ Auth Service: 80%+ test coverage, p95 < 100ms
- ✅ API Gateway: JWT validation, rate limiting, routing
- ✅ All 7 web apps: Login page works, redirects to dashboard
- ✅ CI: All checks pass on PR

---

## 3. Phase 2: Customer Flow (Weeks 3-5)

### Goal
Customer can browse restaurants, add to cart, checkout, and track order.

### Session 4 — Restaurant Catalog + Menu + Order Services (Week 3)

**Tasks**:
- [ ] **Restaurant Catalog Service (Go)**:
  - [ ] Domain: Restaurant entity
  - [ ] Migrations: restaurants table (with PostGIS POINT)
  - [ ] HTTP handlers: GET /restaurants/nearby, GET /restaurants/{id}, GET /restaurants/search
  - [ ] PostGIS queries for geospatial search
  - [ ] Redis cache (60s TTL for restaurant list)
  - [ ] Seed 50 restaurants (Cairo: Maadi, Zamalek, Nasr City)
  - [ ] OpenAPI spec
  - [ ] Tests
- [ ] **Menu Service (Go)**:
  - [ ] Domain: MenuItem, Modifier, Category
  - [ ] Migrations: menu_items, modifiers, categories tables
  - [ ] HTTP handlers: GET /restaurants/{id}/menu
  - [ ] Redis cache (300s TTL for menus)
  - [ ] Seed menus for 50 restaurants
  - [ ] OpenAPI spec
  - [ ] Tests
- [ ] **Order Service (Go)** — partial (create + status):
  - [ ] Domain: Order entity + state machine (8 states)
  - [ ] Migrations: orders, order_items tables (partitioned by month)
  - [ ] HTTP handlers: POST /orders, GET /orders/{id}, GET /orders/active, POST /orders/{id}/cancel
  - [ ] Kafka producer: order.created, order.cancelled events
  - [ ] Idempotency key handling
  - [ ] Pricing calculation (subtotal, delivery fee, service fee, VAT, discount)
  - [ ] OpenAPI spec
  - [ ] Tests

**Deliverable**: Customer can `POST /orders` and get back an order ID.

### Session 5 — Customer Web App: Browse + Cart + Checkout (Week 4)

**Tasks**:
- [ ] **Customer Web App** (`apps/customer-web/`):
  - [ ] **Auth pages**:
    - [ ] Phone input page
    - [ ] OTP verification page
    - [ ] Token storage (Zustand + persist)
  - [ ] **Home page**:
    - [ ] Address picker (current location + saved addresses)
    - [ ] Restaurant cards grid (Trending Nearby section)
    - [ ] Cuisine categories grid
    - [ ] Search bar
  - [ ] **Restaurant detail page**:
    - [ ] Restaurant header (cover, name, rating, ETA)
    - [ ] Menu categories sidebar
    - [ ] Menu items list
    - [ ] Item detail modal (modifiers, quantity, notes)
    - [ ] Add to cart
  - [ ] **Cart page**:
    - [ ] Cart items list
    - [ ] Quantity adjustment
    - [ ] Coupon input
    - [ ] Price breakdown (subtotal, fees, VAT, discount, total)
    - [ ] "Spend more, save more" prompt
  - [ ] **Checkout page**:
    - [ ] Address confirmation
    - [ ] Payment method selector (Vodafone Cash, InstaPay, Card, COD)
    - [ ] Order notes
    - [ ] Place order button
    - [ ] Loading states
    - [ ] Error handling (payment declined, fraud detected, etc.)
  - [ ] **Order confirmation page**:
    - [ ] Order ID, ETA, status
    - [ ] "Track order" CTA
  - [ ] **Profile page**:
    - [ ] User info
    - [ ] Addresses CRUD
    - [ ] Payment methods
    - [ ] Order history
  - [ ] **Design system**:
    - [ ] Tailwind config with theme tokens
    - [ ] RTL support (Arabic)
    - [ ] Cairo font for Arabic, Inter for English
  - [ ] **Tests**:
    - [ ] Component tests (React Testing Library)
    - [ ] E2E test (Playwright): login → browse → add to cart → checkout

**Deliverable**: Customer can place an order end-to-end (without payment processing yet).

### Session 6 — Payment Service + Order Tracking (Week 5)

**Tasks**:
- [ ] **Payment Service (Go)**:
  - [ ] Domain: Payment entity
  - [ ] Migrations: payments table
  - [ ] Vodafone Cash integration (sandbox)
  - [ ] InstaPay integration (sandbox)
  - [ ] Paymob (card) integration (sandbox)
  - [ ] Idempotency (idempotency_key unique constraint)
  - [ ] Refund workflow
  - [ ] Kafka producer: payment.captured, payment.failed, payment.refunded
  - [ ] HTTP handlers: POST /payments/charge, POST /payments/refund
  - [ ] PCI-DSS compliance (no card storage, tokenization)
  - [ ] Tests (with mocked providers)
- [ ] **Order Service** — update to call Payment:
  - [ ] On order create: call Payment Service (gRPC, sync, idempotent)
  - [ ] On payment captured: publish order.confirmed
  - [ ] On payment failed: cancel order
- [ ] **Notification Service (Go)** — basic:
  - [ ] Kafka consumer for order events
  - [ ] WebSocket push to relevant clients
  - [ ] SMS via Twilio (for critical status)
- [ ] **WebSocket Gateway (Go)**:
  - [ ] Standalone service
  - [ ] JWT auth on connection
  - [ ] Redis pub/sub for cross-pod fan-out
  - [ ] Heartbeat (30s)
  - [ ] Auto-reconnect (client-side)
- [ ] **Customer Web App — Order Tracking page**:
  - [ ] Order status timeline (6 steps)
  - [ ] Live ETA (updates every 30s)
  - [ ] WebSocket subscription to `order.{id}` channel
  - [ ] Driver info (when assigned)
  - [ ] Cancel order button (with policy check)
  - [ ] Rate order (after delivery)

**Deliverable**: Customer can place order → pay → track until delivered (delivery mocked for now).

**Exit Criteria (Phase 2)**:
- ✅ 50 restaurants seeded with menus
- ✅ Customer can browse, cart, checkout, pay (sandbox), track
- ✅ p95 checkout flow < 3s
- ✅ Order state machine works correctly
- ✅ WebSocket updates work in real-time
- ✅ Payment idempotency verified (no double-charge)

---

## 4. Phase 3: Restaurant Flow (Weeks 6-7)

### Goal
Restaurant can receive orders, accept/reject, manage menu, view analytics.

### Session 7 — Restaurant Web App: Order Reception + Menu (Week 6)

**Tasks**:
- [ ] **Restaurant Web App** (`apps/restaurant-web/`):
  - [ ] **Auth**:
    - [ ] Phone + OTP login
    - [ ] Sub-account support (manager, cashier roles)
  - [ ] **Dashboard**:
    - [ ] Today's stats (orders, revenue, avg prep time, rating)
    - [ ] Active orders count
    - [ ] Open/closed toggle
  - [ ] **Active orders page**:
    - [ ] Order cards (sorted by time)
    - [ ] Inbound order modal (90s countdown timer)
    - [ ] Accept/Reject buttons
    - [ ] Order details (items, modifiers, notes, payment)
    - [ ] Status update buttons (start prep, ready)
    - [ ] Delay button (with reason)
  - [ ] **Menu management page**:
    - [ ] Categories CRUD
    - [ ] Items CRUD (name, description, price, image, modifiers)
    - [ ] 86'ing (toggle availability)
    - [ ] Bulk operations (hide category)
    - [ ] Image upload (to MinIO/S3)
  - [ ] **Schedule page**:
    - [ ] Operating hours per day
    - [ ] Holiday mode
  - [ ] **WebSocket integration**:
    - [ ] Listen for `order.new` events
    - [ ] Audio alert (loud, persistent)
    - [ ] Browser notification
  - [ ] **Tests**:
    - [ ] Component tests
    - [ ] E2E: receive order → accept → mark ready
- [ ] **Restaurant Catalog Service** — update:
  - [ ] Restaurant-facing endpoints (GET /restaurants/me, PUT /restaurants/me, etc.)
  - [ ] Schedule management
- [ ] **Menu Service** — update:
  - [ ] Restaurant-facing CRUD endpoints
  - [ ] WebSocket event: menu.updated
- [ ] **Order Service** — update:
  - [ ] Restaurant-facing endpoints (accept, reject, start_prep, ready, delay)
  - [ ] Kafka producer: order.status_changed events
  - [ ] Auto-reject after 90s (cron job)

**Deliverable**: Restaurant can receive, accept, prepare, and mark orders ready.

### Session 8 — KDS + Analytics + Promotions (Week 7)

**Tasks**:
- [ ] **Restaurant Web App — KDS (Kitchen Display System)**:
  - [ ] Full-screen KDS mode (for kitchen display)
  - [ ] Order cards in columns (by status: new, preparing, ready)
  - [ ] Color-coded timers (green <5min, yellow 5-10min, red >10min)
  - [ ] Bump order (mark done)
  - [ ] Recall (undo bump within 30s)
  - [ ] Auto-refresh via WebSocket
- [ ] **Analytics page**:
  - [ ] Sales chart (daily/weekly/monthly)
  - [ ] Top items list
  - [ ] Peak hours heatmap
  - [ ] Customer ratings breakdown
  - [ ] Reviews list (with reply capability)
  - [ ] Export to CSV
- [ ] **Promotions page**:
  - [ ] Create promo (4 types: flat, percentage, free delivery, BOGO)
  - [ ] Promo list (active, scheduled, expired)
  - [ ] Performance metrics (redemptions, revenue)
  - [ ] Stop promo button
- [ ] **Promo Service (Go)**:
  - [ ] Domain: Promo entity
  - [ ] Migrations: promos, promo_redemptions tables
  - [ ] HTTP handlers: CRUD for promos
  - [ ] Validation (date range, usage limits)
  - [ ] Anti-abuse rules (max 1 per order, 3 per month per user)
  - [ ] Kafka consumer: order.created (to track redemptions)
  - [ ] Tests
- [ ] **Analytics Service (Go)** — basic:
  - [ ] Kafka consumer for all events
  - [ ] ClickHouse writes
  - [ ] Aggregation queries
  - [ ] HTTP handlers: GET /analytics/sales, /analytics/top_items, etc.

**Deliverable**: Restaurant has full toolkit: receive, manage, analyze, promote.

**Exit Criteria (Phase 3)**:
- ✅ Restaurant can accept/reject orders in <30s
- ✅ KDS works on tablet (responsive)
- ✅ Menu updates reflect in customer app in <60s
- ✅ Analytics accurate (matches order data)
- ✅ Promotions apply correctly at checkout

---

## 5. Phase 4: Driver Flow (Weeks 8-9)

### Goal
Driver can go online, receive orders, deliver, and earn.

### Session 9 — Driver Management + Geo + Matching Services (Week 8)

**Tasks**:
- [ ] **Driver Management Service (Go)**:
  - [ ] Domain: Driver entity (profile, KYC, vehicle, rating)
  - [ ] Migrations: drivers, driver_documents tables
  - [ ] HTTP handlers: GET /drivers/me, PUT /drivers/me/status, etc.
  - [ ] Tier calculation (Platinum/Gold/Silver/Standard)
  - [ ] Acceptance rate tracking
  - [ ] Tests
- [ ] **Geo Service (Go)**:
  - [ ] Redis GEO integration (GEOADD, GEOSEARCH)
  - [ ] Driver location ingestion (via WebSocket Gateway)
  - [ ] ETA calculation (Google Maps API)
  - [ ] Kafka consumer: driver.location (for archival)
  - [ ] HTTP handlers: POST /drivers/location, GET /drivers/nearby
  - [ ] Tests
- [ ] **Delivery Matching Service (Go)** ⭐:
  - [ ] Kafka consumer: order.confirmed
  - [ ] Algorithm: GEOSEARCH → score → broadcast to top-3
  - [ ] Scoring: distance (30%), acceptance (20%), rating (15%), idle (15%), completion (10%), vehicle (10%)
  - [ ] Retry logic (3 rounds, expanding radius)
  - [ ] Manual dispatch queue (if all rounds fail)
  - [ ] Kafka producer: order.assigned
  - [ ] Tests (algorithm correctness)
- [ ] **Seed 20 drivers** (for testing)

**Deliverable**: When restaurant marks order ready, system finds a driver in <15s.

### Session 10 — Driver Web App (Week 9)

**Tasks**:
- [ ] **Driver Web App** (`apps/driver-web/`):
  - [ ] **Auth + Onboarding**:
    - [ ] Phone + OTP login
    - [ ] Profile setup (KYC basic — for MVP, full KYC via Field Supervisor later)
    - [ ] Vehicle info
  - [ ] **Home page (online)**:
    - [ ] Today's earnings
    - [ ] This week's earnings
    - [ ] Heat map (hot zones)
    - [ ] Tier badge + rating
    - [ ] Acceptance rate
    - [ ] Online/Offline toggle
  - [ ] **Order offer modal**:
    - [ ] 15-second countdown timer
    - [ ] Restaurant info + distance
    - [ ] Customer location + distance
    - [   ] Earnings breakdown (base + distance + peak bonus)
    - [ ] Accept/Reject buttons
    - [ ] Auto-decline on timeout
  - [ ] **Active delivery page**:
    - [ ] 5-step progress (going to pickup, picked up, going to dropoff, delivered)
    - [ ] Map with route
    - [ ] Restaurant info + items list
    - [ ] Customer info + address + notes
    - [ ] Call/Message buttons (anonymized)
    - [ ] Pickup confirmation (QR code from restaurant)
    - [ ] Dropoff confirmation (OTP from customer + photo)
    - [ ] Report problem button
  - [ ] **Earnings page**:
    - [ ] Daily/weekly breakdown
    - [ ] Per-order earnings
    - [ ] Instant payout button (Vodafone Cash)
    - [ ] Payout history
  - [ ] **Order history page**:
    - [ ] Past deliveries
    - [ ] Ratings received
  - [ ] **WebSocket integration**:
    - [ ] Subscribe to `order.new` channel
    - [ ] Audio alert on new order
    - [ ] Background location updates (every 5s)
  - [ ] **Tests**:
    - [ ] Component tests
    - [ ] E2E: accept order → pickup → deliver

**Deliverable**: Driver can accept, pickup, and deliver orders end-to-end.

**Exit Criteria (Phase 4)**:
- ✅ Driver matching p95 < 15s
- ✅ Driver app shows order offers in <2s
- ✅ GPS tracking updates customer app in real-time
- ✅ OTP + photo verification prevents fake deliveries
- ✅ Instant payout works (sandbox)

---

## 6. Phase 5: Support & Ops (Weeks 10-11)

### Goal
Support can resolve customer issues; ops can monitor and intervene.

### Session 11 — Support Web App + AI Chatbot Foundation (Week 10)

**Tasks**:
- [ ] **Support Web App** (`apps/support-web/`):
  - [ ] **Auth**:
    - [ ] Username + password + TOTP (mandatory 2FA)
    - [ ] Role: support_l1, support_l2
  - [ ] **Dashboard**:
    - [ ] Active tickets count
    - [ ] Today's stats (resolved, avg time, CSAT)
    - [ ] Active alerts
  - [ ] **Ticket queue**:
    - [ ] Sortable/filterable list
    - [ ] Auto-assignment
    - [ ] Priority indicators
  - [ ] **Chat view**:
    - [ ] Customer info panel (with trust score)
    - [ ] Order info panel (with link to order)
    - [ ] Chat messages
    - [ ] Quick actions (refund, cancel, extend ETA, escalate)
    - [ ] Macros (canned responses)
    - [   ] Sentiment indicator (real-time)
  - [ ] **Refund workflow**:
    - [ ] Refund dialog (amount, reason)
    - [ ] Biometric verification (WebAuthn) for refunds >EGP 100
    - [ ] Dual approval for refunds >EGP 500
    - [ ] Refund history per customer
  - [ ] **Knowledge base**:
    - [ ] Article search
    - [ ] Article suggestion in chat
  - [ ] **WebSocket integration**:
    - [ ] Live ticket updates
    - [ ] New ticket alerts
- [ ] **Notification Service** — update:
  - [ ] In-app messaging (chat between support and customer)
  - [ ] Ticket events
- [ ] **Fraud Service (Go)** — basic:
  - [ ] Rule-based fraud scoring (ML model in Phase 6)
  - [ ] Trust score per customer
  - [ ] Kafka consumer: order.created
  - [ ] Kafka producer: fraud.score_calculated
  - [ ] HTTP handler: GET /fraud/score/{customer_id}
  - [ ] Tests

**Deliverable**: Support can resolve customer issues with full context.

### Session 12 — Command Center (Week 11)

**Tasks**:
- [ ] **Command Center** (`apps/command-center/`):
  - [ ] **Auth**:
    - [ ] Username + password + TOTP
    - [   ] Role: ops_manager, super_admin
  - [ ] **Main dashboard**:
    - [ ] KPI cards (orders, drivers, restaurants, GMV, issues)
    - [ ] Live metrics bar (orders/min, avg delivery time, etc.)
    - [ ] Active alerts panel
  - [ ] **Live map**:
    - [ ] Demand heatmap (red/yellow/green zones)
    - [ ] Driver pins (real-time)
    - [ ] Order flow lines
    - [   ] Zone detail on click
  - [ ] **Zone management**:
    - [ ] Surge pricing override (multiplier + duration)
    - [ ] Driver notification ("need drivers in Maadi")
    - [ ] Pause orders in zone
  - [ ] **Manual interventions**:
    - [ ] Assign driver to order manually
    - [ ] Cancel order
    - [ ] Suspend driver
    - [ ] Pause/deactivate restaurant
  - [ ] **Incident management**:
    - [ ] Incident list (P0/P1/P2)
    - [ ] Acknowledge/resolve workflow
    - [ ] Postmortem notes
  - [ ] **Analytics**:
    - [ ] Real-time metrics
    - [ ] Today's GMV chart
    - [ ] Driver utilization
    - [   ] Export to CSV
  - [ ] **WebSocket integration**:
    - [ ] Real-time metric updates
    - [ ] Alert notifications
    - [ ] Incident updates
- [ ] **Analytics Service** — update:
    - [ ] Real-time aggregations (orders/min, GMV/sec)
    - [ ] Historical queries (daily/weekly/monthly)
    - [ ] WebSocket push to Command Center
- [ ] **Forecasting** (basic, ML in Phase 6):
    - [ ] Historical average + day-of-week adjustment
    - [   ] Hourly forecast for next 24h

**Deliverable**: Ops manager has full visibility + control.

**Exit Criteria (Phase 5)**:
- ✅ Support first response time < 2 min
- ✅ Command Center updates in real-time (< 5s lag)
- ✅ Manual interventions work (surge, assign, suspend)
- ✅ Incident workflow functional
- ✅ Basic fraud detection blocks repeat offenders

---

## 7. Phase 6: Trust & Field (Week 12)

### Goal
Complete the platform with employee portal + field supervisor app.

### Session 13 — Employee Portal + Field Supervisor + ML (Week 12)

**Tasks**:
- [ ] **Employee Portal** (`apps/employee-portal/`):
  - [ ] **Auth**:
    - [ ] Username + password + TOTP
    - [ ] WebAuthn biometric enrollment
    - [ ] Sensitive action verification (WebAuthn challenge-response)
  - [ ] **Dashboard**:
    - [ ] Pending approvals (dual approval workflow)
    - [ ] Today's stats
    - [ ] Audit log preview
  - [ ] **Sensitive actions**:
    - [ ] Refund (with biometric + dual approval for >EGP 500)
    - [ ] Restaurant approve/deactivate
    - [ ] Driver approve/deactivate
    - [ ] Payout approval
    - [   ] Menu override
    - [ ] Data export
  - [ ] **Audit log viewer**:
    - [ ] Searchable, filterable
    - [   ] Hash chain verification status
    - [ ] Export to CSV/PDF
  - [ ] **Access review**:
    - [ ] Team list with permissions
    - [ ] Quarterly review workflow
  - [ ] **Anomaly alerts**:
    - [ ] View flagged employees
    - [ ] Investigate workflow
- [ ] **Trust Service (Go)** ⭐:
  - [ ] Immutable audit_logs table (REVOKE UPDATE/DELETE)
  - [   ] Hash chain (prev_hash + record_hash)
  - [ ] Nightly chain verification job
  - [ ] WebAuthn endpoints (registration, verification)
  - [ ] Dual approval workflow engine
  - [ ] Conflict of interest checks
  - [ ] Quarterly access review
  - [ ] Tests
- [ ] **Field Supervisor Web App** (`apps/field-supervisor-web/`):
  - [ ] **Auth**:
    - [ ] Phone + OTP
    - [ ] Role: field_supervisor
  - [ ] **Task list**:
    - [ ] Today's tasks (sorted by priority + ETA)
    - [   ] Task detail (restaurant/driver info, checklist)
  - [ ] **Restaurant verification**:
    - [ ] 50+ point checklist
    - [   ] Photo capture (with GPS metadata)
    - [ ] GPS verification (must be within 50m of restaurant)
    - [   ] Decision: approve/conditional/reject
  - [ ] **Driver verification**:
    - [ ] Document check (ID, license, vehicle)
    - [ ] Photo capture
    - [   ] Practical test (basic)
  - [ ] **Complaint investigation**:
    - [ ] Investigation workflow
    - [   ] Photo evidence
    - [ ] Resolution report
  - [ ] **Route planning**:
    - [ ] TSP solver for daily tasks
    - [   ] Google Maps integration
  - [ ] **Reports**:
    - [ ] Daily report submission
    - [   ] Weekly performance
- [ ] **ML Models** (basic, Python sidecars):
  - [ ] Fraud detection: rule-based + simple logistic regression
  - [   ] Driver fraud: GPS trajectory validation
  - [ ] Anomaly detection: Isolation Forest for UEBA
  - [   ] Sentiment analysis: Arabic BERT for chat sentiment
- [ ] **E2E integration tests**:
  - [ ] Full customer flow (login → order → pay → track → rate)
  - [ ] Full restaurant flow (login → receive → accept → prepare → ready)
  - [ ] Full driver flow (login → accept → pickup → deliver → earn)
  - [ ] Support intervention flow
  - [   ] Ops manual override flow
  - [ ] Employee biometric action flow

**Deliverable**: Full platform operational with all 7 apps + 12 services.

**Exit Criteria (Phase 6)**:
- ✅ All 7 web apps functional
- ✅ All 12 services deployed and communicating
- ✅ Audit log hash chain verified
- ✅ Biometric auth works for sensitive actions
- ✅ Field supervisor can verify restaurant/driver
- ✅ E2E tests pass for all critical flows

---

## 8. Milestones & Exit Criteria

### 8.1 Weekly Milestones

| Week | Milestone | Demo |
|------|-----------|------|
| 1 | Monorepo + dev env | `make setup` works |
| 2 | Auth + API Gateway | User can login |
| 3 | Catalog + Menu + Order services | API returns restaurants |
| 4 | Customer web (browse + cart) | Customer can browse |
| 5 | Payment + tracking | Customer can order end-to-end |
| 6 | Restaurant web | Restaurant can receive orders |
| 7 | KDS + analytics + promos | Restaurant has full toolkit |
| 8 | Driver services | System matches drivers |
| 9 | Driver web app | Driver can deliver |
| 10 | Support web + fraud | Support can resolve issues |
| 11 | Command center | Ops can monitor |
| 12 | Employee portal + field + ML | Full system ready |

### 8.2 Phase Gates (Go/No-Go)

#### Gate 1 — End of Week 2 (Foundation Validated)
**Go criteria (ALL)**:
- ✅ Auth Service: 80% coverage, p95 < 100ms
- ✅ API Gateway: routing + JWT + rate limiting
- ✅ All 7 web apps: login works
- ✅ CI: all green

**No-go**: 1-week extension; focus on stability.

#### Gate 2 — End of Week 5 (Customer Flow Validated)
**Go criteria**:
- ✅ 50 restaurants seeded
- ✅ Customer can browse + cart + checkout + pay (sandbox) + track
- ✅ p95 checkout < 3s
- ✅ WebSocket real-time updates work

**No-go**: 2-week extension; no restaurant app until fixed.

#### Gate 3 — End of Week 9 (Three-Sided Marketplace Validated)
**Go criteria**:
- ✅ Customer + Restaurant + Driver apps all functional
- ✅ End-to-end order flow works (customer → restaurant → driver → customer)
- ✅ Driver matching p95 < 15s
- ✅ GPS tracking real-time

**No-go**: 2-week extension; focus on integration.

#### Gate 4 — End of Week 12 (Production-Ready MVP)
**Go criteria**:
- ✅ All 7 web apps deployed
- ✅ All 12 services in production (staging)
- ✅ Audit log hash chain verified
- ✅ E2E tests pass
- ✅ Load test: 1000 concurrent users
- ✅ Security review completed

**No-go**: 2-week extension; production launch delayed.

---

## 9. Risk Register

| # | Risk | Probability | Impact | Mitigation | Owner |
|---|------|-------------|--------|------------|-------|
| 1 | Vodafone Cash API sandbox unstable | Medium | High | Mock provider; integrate InstaPay first | Backend |
| 2 | WebSocket scaling issues | Low | High | Load test early (Week 5); Redis pub/sub | Backend |
| 3 | WebAuthn browser compatibility | Low | Medium | Test on Chrome, Firefox, Safari; provide TOTP fallback | Frontend |
| 4 | Kafka consumer lag | Medium | Medium | KEDA autoscaling; alert if lag >1000 | Backend |
| 5 | Field supervisor GPS spoofing | Medium | High | Photo metadata + geofencing + spot checks | Trust |
| 6 | Audit log hash chain breaks | Low | Critical | Nightly verification; P0 incident if broken | Trust |
| 7 | Scope creep (adding features) | High | Medium | Strict phase gates; defer to v2 | PM |
| 8 | Team availability | Medium | High | Pair programming; docs as memory | All |

---

## 10. Resource Requirements

### 10.1 Team (Minimum Viable)

| Role | Count | Weeks | Responsibilities |
|------|-------|-------|------------------|
| Tech Lead (you + AI) | 1 | 12 | Architecture, code, review |
| Backend Engineer | 1 | 12 | Go services |
| Frontend Engineer | 1 | 12 | React apps |
| DevOps Engineer | 0.5 | 12 | Infra, CI/CD |
| QA Engineer | 0.5 | 8-12 | Tests, E2E |
| Designer | 0.3 | 4 | UI/UX (Weeks 3-7) |

**Total: 4.3 FTE for 12 weeks**

### 10.2 Infrastructure (Monthly Cost)

| Resource | Dev | Staging | Production (Year 1) |
|----------|-----|---------|---------------------|
| AWS EKS | $100 | $300 | $1,500 |
| RDS PostgreSQL | $50 | $200 | $1,200 |
| ElastiCache Redis | $30 | $100 | $500 |
| MSK Kafka | $100 | $200 | $800 |
| S3 + CloudFront | $10 | $30 | $200 |
| Other (ES, ClickHouse) | $50 | $100 | $500 |
| **Total** | **$340** | **$930** | **$4,700** |

### 10.3 Third-Party Services (Monthly)

| Service | Dev | Production |
|---------|-----|------------|
| Twilio (SMS) | $20 | $500 |
| Google Maps API | $0 (free tier) | $300 |
| Sentry | $0 (free tier) | $100 |
| Vodafone Cash | $0 (sandbox) | Variable (per transaction) |
| InstaPay | $0 (sandbox) | Variable |
| Paymob | $0 (sandbox) | 1.5% per transaction |
| Cloudflare | $0 (free tier) | $50 |

**Total third-party (production)**: ~$1,000/month + transaction fees

---

## Appendix A: Daily Session Structure

Each session (Day 1-5 of a week) follows this structure:

```
09:00 - 09:30: Review yesterday's work + plan today
09:30 - 12:30: Implementation (pair programming)
12:30 - 13:30: Lunch
13:30 - 16:30: Implementation continues
16:30 - 17:30: Code review + tests
17:30 - 18:00: Update SESSIONS-LOG.md + plan tomorrow
```

## Appendix B: Session Log Template

```markdown
## Session N - YYYY-MM-DD

**Goal**: <one sentence>

**Completed**:
- [ ] Task 1
- [ ] Task 2

**Decisions**:
- Decision 1 (rationale)

**Blockers**:
- Blocker 1 (needs X)

**Next session**:
- Task A
- Task B

**Files touched**:
- path/to/file1
- path/to/file2
```

---

> **Next**: Read `REFERENCES.md` for all research sources.
