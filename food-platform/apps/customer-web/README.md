# Customer Web App

The customer-facing web app for browsing restaurants, ordering food, tracking deliveries, and managing account.

## Quick Start

```bash
# Copy env file
cp .env.example .env

# Install dependencies (from monorepo root)
pnpm install

# Start dev server
pnpm dev
```

App runs at http://localhost:5173

## Routes

| Path | Description |
|------|-------------|
| `/login` | Phone login |
| `/otp` | OTP verification |
| `/` | Home (restaurants feed) |
| `/restaurants/:id` | Restaurant detail + menu |
| `/cart` | Shopping cart |
| `/checkout` | Checkout + payment |
| `/orders/:id` | Order tracking |
| `/orders` | Order history |
| `/profile` | User profile |

## Tech Stack

- React 18 + Vite
- React Router 6
- TanStack Query (data fetching)
- Zustand (state)
- Tailwind CSS (styled with @food-platform/theme)
- React Hook Form + Zod (forms)

## Documentation

- [UI Spec](../../docs/ui-ux/customer-web/UI-SPEC.md) — detailed screen specs
- [Architecture](../../docs/ARCHITECTURE.md) — system architecture
- [API Contracts](../../docs/API-CONTRACTS.md) — REST + WebSocket + Kafka
