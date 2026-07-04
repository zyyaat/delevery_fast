# Food Delivery Platform — Makefile
# Common commands for development

.PHONY: help setup dev build test lint format type-check codegen clean docker-up docker-down docker-logs docker-reset migrate migrate-down

# Default target
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ============ Setup ============

setup: docker-up install codegen ## Full setup: start services, install deps, generate code
	@echo "✅ Setup complete! Run 'make dev' to start development."

install: ## Install Node.js dependencies
	pnpm install

install:go: ## Install Go dependencies for all services
	@for svc in services/*/; do \
		echo "Installing Go deps for $$svc..."; \
		(cd $$svc && go mod download) || true; \
	done

# ============ Development ============

dev: ## Start all apps in dev mode
	pnpm dev

dev:customer: ## Start customer web only
	cd apps/customer-web && pnpm dev

dev:driver: ## Start driver web only
	cd apps/driver-web && pnpm dev

dev:restaurant: ## Start restaurant web only
	cd apps/restaurant-web && pnpm dev

dev:all-apps: ## Start all 7 web apps
	@echo "Starting all apps in parallel..."
	pnpm -r --parallel dev

# ============ Build ============

build: ## Build all packages and apps
	pnpm build

build:apps: ## Build all web apps only
	pnpm -r --filter="./apps/*" build

build:packages: ## Build all shared packages only
	pnpm -r --filter="./packages/*" build

# ============ Testing ============

test: ## Run all tests
	pnpm test

test:watch: ## Run tests in watch mode
	pnpm test -- --watch

test:e2e: ## Run E2E tests (Playwright)
	pnpm -r --filter="./apps/*" test:e2e

test:coverage: ## Run tests with coverage
	pnpm test -- --coverage

# ============ Code Quality ============

lint: ## Lint all code
	pnpm lint

lint:fix: ## Lint and fix
	pnpm lint -- --fix

format: ## Format all code with Prettier
	pnpm format

format:check: ## Check formatting without writing
	pnpm format:check

type-check: ## TypeScript type checking
	pnpm type-check

# ============ Code Generation ============

codegen: ## Generate types from OpenAPI specs
	@echo "TODO: Add OpenAPI codegen"
	# pnpm codegen

# ============ Docker / Infrastructure ============

docker-up: ## Start all Docker services (PostgreSQL, Redis, Kafka, etc.)
	docker compose -f infra/docker-compose.yml up -d
	@echo ""
	@echo "✅ Services started:"
	@echo "   PostgreSQL: localhost:5432 (postgres/postgres)"
	@echo "   Redis:      localhost:6379"
	@echo "   Kafka:      localhost:9092"
	@echo "   Kafka UI:   http://localhost:9000"
	@echo "   ES:         http://localhost:9200"
	@echo "   ClickHouse: http://localhost:8123"
	@echo "   MinIO:      http://localhost:9002 (minio/minio123)"
	@echo "   Mailhog:    http://localhost:8025"

docker-down: ## Stop all Docker services
	docker compose -f infra/docker-compose.yml down

docker-logs: ## Tail Docker logs
	docker compose -f infra/docker-compose.yml logs -f

docker-reset: docker-down ## Stop services and delete all data (DESTRUCTIVE)
	docker compose -f infra/docker-compose.yml down -v
	@echo "⚠️  All Docker volumes deleted!"

docker-ps: ## Show running Docker services
	docker compose -f infra/docker-compose.yml ps

# ============ Database ============

migrate: ## Run migrations for all services
	@for svc in auth restaurant-catalog menu order payment delivery-matching driver-management geo notification fraud promo analytics; do \
		echo "Migrating $$svc..."; \
		(cd services/$$svc && DATABASE_URL=postgres://postgres:postgres@localhost:5432/$${svc//-/_}_db?sslmode=disable make migrate) || echo "⚠️  $$svc migration skipped (no migrations yet)"; \
	done

migrate:svc: ## Run migrations for a single service (usage: make migrate:svc SVC=auth)
	@if [ -z "$(SVC)" ]; then echo "Usage: make migrate:svc SVC=auth"; exit 1; fi
	cd services/$(SVC) && DATABASE_URL=postgres://postgres:postgres@localhost:5432/$(SVC)//_db?sslmode=disable make migrate

migrate-down: ## Rollback last migration for all services
	@for svc in auth restaurant-catalog menu order payment delivery-matching driver-management geo notification fraud promo analytics; do \
		echo "Rolling back $$svc..."; \
		(cd services/$$svc && DATABASE_URL=postgres://postgres:postgres@localhost:5432/$${svc//-/_}_db?sslmode=disable make migrate-down) || true; \
	done

# ============ Go Services ============

go:run: ## Run a Go service (usage: make go:run SVC=auth)
	@if [ -z "$(SVC)" ]; then echo "Usage: make go:run SVC=auth"; exit 1; fi
	cd services/$(SVC) && go run ./cmd/server

go:test: ## Test all Go services
	@for svc in services/*/; do \
		echo "Testing $$svc..."; \
		(cd $$svc && go test -v -race -cover ./...) || true; \
	done

go:lint: ## Lint all Go services
	@for svc in services/*/; do \
		echo "Linting $$svc..."; \
		(cd $$svc && golangci-lint run ./...) || true; \
	done

go:build: ## Build all Go services
	@for svc in services/*/; do \
		echo "Building $$svc..."; \
		(cd $$svc && go build -o bin/server ./cmd/server) || true; \
	done

# ============ Cleanup ============

clean: ## Clean all build artifacts and node_modules
	pnpm clean
	@for svc in services/*/; do \
		(cd $$svc && rm -rf bin/) || true; \
	done
	rm -rf node_modules

clean:dist: ## Clean build outputs only (keep node_modules)
	pnpm -r exec rm -rf dist .turbo
	@for svc in services/*/; do \
		(cd $$svc && rm -rf bin/) || true; \
	done

# ============ Documentation ============

docs:serve: ## Serve docs locally (if mdserve installed)
	@echo "Open docs/ in your editor or use a Markdown previewer"

# ============ Misc ============

env: ## Copy .env.example files to .env
	@find . -name ".env.example" -exec sh -c 'cp -n "$$1" "$${1%.example}"' _ {} \;
	@echo "✅ .env files created from .env.example"

info: ## Show project info
	@echo "Food Delivery Platform"
	@echo "======================"
	@echo "Apps:       $$(ls apps | wc -l) web apps"
	@echo "Packages:   $$(ls packages | wc -l) shared packages"
	@echo "Services:   $$(ls services | wc -l) Go microservices"
	@echo ""
	@echo "Documentation:"
	@echo "  Architecture:   docs/ARCHITECTURE.md"
	@echo "  API Contracts:  docs/API-CONTRACTS.md"
	@echo "  Patterns:       docs/PATTERNS.md"
	@echo "  Roadmap:        docs/ROADMAP.md"
	@echo "  Design System:  docs/ui-ux/DESIGN-SYSTEM.md"
