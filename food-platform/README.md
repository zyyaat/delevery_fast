# Food Delivery Platform

> **Status**: In Development  
> **Started**: 2026-07-04  
> **Methodology**: Contract-First + Vertical Slices + Incremental  
> **Target**: Production-grade food delivery platform for Egyptian market

---

## 🎯 Overview

A complete food delivery platform (Uber Eats / Talabat-class) consisting of:

- **7 web applications** (React 18 + Vite)
- **12 Go microservices** (Domain-oriented, event-driven)
- **PostgreSQL + Redis + Kafka** backend
- **Kubernetes** infrastructure on AWS

Targeting **50,000+ orders/day** at full scale in the Egyptian market.

---

## 📁 Project Structure

```
food-platform/
├── apps/                          # 7 React web applications
│   ├── customer-web/              # Customer ordering app
│   ├── driver-web/                # Driver delivery app
│   ├── restaurant-web/            # Restaurant management app
│   ├── support-web/               # Customer support app
│   ├── command-center/            # Operations dashboard
│   ├── employee-portal/           # Internal employee portal
│   └── field-supervisor-web/      # Field supervisor app
│
├── packages/                      # Shared TypeScript packages
│   ├── api-client/                # HTTP + WebSocket client
│   ├── types/                     # TypeScript types (generated)
│   ├── hooks/                     # Business logic hooks
│   ├── utils/                     # Shared utilities
│   ├── theme/                     # Design tokens
│   ├── ui/                        # Shared components
│   └── auth/                      # Auth helpers
│
├── services/                      # 12 Go microservices
│   ├── auth/                      # Authentication
│   ├── restaurant-catalog/        # Restaurant data
│   ├── menu/                      # Menu management
│   ├── order/                     # Order lifecycle
│   ├── payment/                   # Payment processing
│   ├── delivery-matching/         # Driver matching
│   ├── driver-management/         # Driver profiles
│   ├── geo/                       # GPS tracking
│   ├── notification/              # Push/SMS/email
│   ├── fraud/                     # Fraud detection
│   ├── promo/                     # Promotions & loyalty
│   └── analytics/                 # Analytics & reporting
│
├── infra/                         # Infrastructure
│   ├── docker-compose.yml         # Dev environment
│   ├── terraform/                 # AWS IaC
│   └── k8s/                       # Kubernetes manifests
│
├── tools/                         # Tooling & scripts
├── docs/                          # All documentation
└── README.md                      # This file
```

---

## 📚 Documentation

**READ THESE FIRST** (in order):

1. **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** — Full system architecture
2. **[docs/API-CONTRACTS.md](docs/API-CONTRACTS.md)** — REST, WebSocket, Kafka contracts
3. **[docs/PATTERNS.md](docs/PATTERNS.md)** — Code patterns (Go + React)
4. **[docs/REPO-STRUCTURE.md](docs/REPO-STRUCTURE.md)** — Monorepo layout
5. **[docs/ROADMAP.md](docs/ROADMAP.md)** — 12-week execution plan
6. **[docs/REFERENCES.md](docs/REFERENCES.md)** — 133 research sources
7. **[docs/SESSIONS-LOG.md](docs/SESSIONS-LOG.md)** — Session memory (update every session)

---

## 🚀 Quick Start (When Implemented)

### Prerequisites
- Node.js 20+
- pnpm 9+
- Go 1.22+
- Docker + Docker Compose
- Make

### Setup
```bash
# Clone repository
git clone <repo-url>
cd food-platform

# Install dependencies
pnpm install

# Start dev infrastructure (PostgreSQL, Redis, Kafka, etc.)
make docker-up

# Generate types from OpenAPI
make codegen

# Run database migrations
make migrate

# Start all services (dev mode)
make dev
```

### Access Points (when running)
- Customer Web: http://localhost:5173
- Driver Web: http://localhost:5174
- Restaurant Web: http://localhost:5175
- Support Web: http://localhost:5176
- Command Center: http://localhost:5177
- Employee Portal: http://localhost:5178
- Field Supervisor Web: http://localhost:5179
- API Gateway: http://localhost:8080
- WebSocket Gateway: ws://localhost:8081
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- Kafka: localhost:9092

---

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.22
- **HTTP Router**: chi
- **gRPC**: protobuf + go-grpc
- **Database**: PostgreSQL 16 + PostGIS
- **Cache**: Redis 7
- **Event Bus**: Apache Kafka 3.x (with Schema Registry)
- **Search**: ElasticSearch 8
- **Analytics**: ClickHouse
- **Object Storage**: S3 / MinIO

### Frontend
- **Framework**: React 18 + Vite
- **State**: Zustand
- **Data Fetching**: TanStack Query
- **Routing**: React Router 6
- **Styling**: Tailwind CSS + shadcn/ui-style components
- **Forms**: React Hook Form + Zod
- **WebSocket**: native + reconnect logic

### Infrastructure
- **Container Orchestration**: Kubernetes (EKS)
- **IaC**: Terraform + Terragrunt
- **CI/CD**: GitHub Actions + ArgoCD (GitOps)
- **Observability**: Prometheus, Grafana, Loki, Jaeger, Sentry
- **Secrets**: AWS Secrets Manager + External Secrets Operator

---

## 📊 Architecture Summary

### 7 Web Apps
| App | Purpose | Users |
|-----|---------|-------|
| Customer Web | Browse, order, track | Customers |
| Driver Web | Accept orders, deliver | Drivers |
| Restaurant Web | Receive orders, manage menu | Restaurants |
| Support Web | Handle tickets, refunds | Support Agents |
| Command Center | Live ops monitoring | Ops Managers |
| Employee Portal | Internal admin (biometric) | Internal Employees |
| Field Supervisor | On-site verification | Field Supervisors |

### 12 Microservices
| Service | Responsibility |
|---------|----------------|
| Auth | JWT, OTP, WebAuthn, sessions |
| Restaurant Catalog | Restaurant data, geospatial |
| Menu | Items, modifiers, availability |
| Order | Order lifecycle, state machine |
| Payment | Vodafone Cash, InstaPay, cards |
| Delivery Matching | Driver-customer matching |
| Driver Management | Driver profiles, KYC |
| Geo | GPS tracking, ETA |
| Notification | Push, SMS, email |
| Fraud | ML scoring, rule engine |
| Promo | Coupons, loyalty, cashback |
| Analytics | Aggregations, dashboards |

### Communication Patterns
- **Synchronous (gRPC)**: When caller needs result to proceed
- **Asynchronous (Kafka)**: When announcing events
- **Real-time (WebSocket)**: When clients need push

---

## 🗺️ Roadmap (12 Weeks)

| Phase | Weeks | Goal |
|-------|-------|------|
| 1. Foundation | 1-2 | Infra + contracts + auth |
| 2. Customer Flow | 3-5 | Browse + cart + checkout + tracking |
| 3. Restaurant Flow | 6-7 | Receive + accept + KDS + menu |
| 4. Driver Flow | 8-9 | Online + accept + deliver + earn |
| 5. Support & Ops | 10-11 | Tickets + command center |
| 6. Trust & Field | 12 | Employee portal + field supervisor |

See [docs/ROADMAP.md](docs/ROADMAP.md) for full details.

---

## 🔒 Security Highlights

- **Authentication**: Phone+OTP (customers/drivers), WebAuthn biometric (employees)
- **Authorization**: RBAC with 11 roles
- **Audit Log**: Immutable, hash-chained (tamper-evident)
- **Fraud Detection**: 3-layer defense (customer, driver, internal)
- **Compliance**: Egyptian Law 175/2018, PCI-DSS

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) Section 9 for details.

---

## 📝 Development Methodology

1. **Contract-First**: API contracts before any code
2. **Vertical Slices**: Each slice delivers end-to-end functionality
3. **Test-Driven**: 80%+ coverage on business logic
4. **Documentation as Code**: Docs stay in sync with code
5. **Session Memory**: `SESSIONS-LOG.md` updated every session

---

## 🤝 Contributing

1. Read all docs in `docs/` first
2. Follow patterns in `docs/PATTERNS.md`
3. Update `docs/SESSIONS-LOG.md` after each session
4. Write tests for all business logic
5. Ensure CI passes before merge

---

## 📜 License

Proprietary — All rights reserved.

---

> **Note**: This project is in active development. Refer to `docs/SESSIONS-LOG.md` for current status.
