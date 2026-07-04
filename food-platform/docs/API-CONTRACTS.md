# API Contracts — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active — Source of Truth  
> **Last Updated**: 2026-07-04  
> **Standard**: OpenAPI 3.1 + AsyncAPI 2.6 + WebSocket Subprotocol

This document defines **all** contracts between the 7 web apps and the 12 backend services. Any change to these contracts MUST go through a PR review and update of this document.

---

## Table of Contents

1. [Conventions](#1-conventions)
2. [Authentication](#2-authentication)
3. [REST API — Customer Web](#3-rest-api--customer-web)
4. [REST API — Driver Web](#4-rest-api--driver-web)
5. [REST API — Restaurant Web](#5-rest-api--restaurant-web)
6. [REST API — Support Web](#6-rest-api--support-web)
7. [REST API — Command Center](#7-rest-api--command-center)
8. [REST API — Employee Portal](#8-rest-api--employee-portal)
9. [REST API — Field Supervisor](#9-rest-api--field-supervisor)
10. [WebSocket Events](#10-websocket-events)
11. [Kafka Events (AsyncAPI)](#11-kafka-events-asyncapi)
12. [Common Schemas](#12-common-schemas)
13. [Error Codes](#13-error-codes)

---

## 1. Conventions

### 1.1 Base URL

```
Production:  https://api.food-platform.com
Staging:     https://api.staging.food-platform.com
Dev (local): http://localhost:8080
```

### 1.2 API Versioning

- Path-based: `/api/v1/...`
- Backwards-compatible changes within v1
- Breaking changes → `/api/v2/...` (with 6-month deprecation overlap)

### 1.3 Request/Response Format

- `Content-Type: application/json` (REST)
- `Content-Type: application/protobuf` (gRPC)
- Dates: ISO 8601 (`2026-07-04T14:32:15Z`)
- IDs: UUID v4
- Money: decimal as string (`"245.50"`) to avoid float precision issues
- Coordinates: `{ "lat": 30.0444, "lng": 31.2357 }` (decimal degrees, 7 places)

### 1.4 Authentication Header

```
Authorization: Bearer <jwt-access-token>
X-Idempotency-Key: <uuid-v4>   # for POST requests that mutate state
X-Request-ID: <uuid-v4>         # for tracing (auto-generated if missing)
```

### 1.5 Rate Limiting

| Endpoint Category | Rate Limit | Burst |
|-------------------|-----------|-------|
| Public (login, OTP) | 10 req/min/IP | 20 |
| Customer authenticated | 100 req/min/user | 200 |
| Driver authenticated | 200 req/min/user | 500 |
| Restaurant authenticated | 500 req/min/store | 1000 |
| Internal employee | 60 req/min/user | 100 |
| API keys (POS integration) | 1000 req/min/key | 2000 |

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 87
X-RateLimit-Reset: 1720000900
```

---

## 2. Authentication

### 2.1 Customer/Driver Login (Phone + OTP)

#### Request OTP

```
POST /api/v1/auth/otp/send
Content-Type: application/json

{
  "phone": "+201012345678",
  "role": "customer"  // customer | driver | restaurant
}
```

#### Response 200

```json
{
  "request_id": "uuid-v4",
  "expires_in": 120,  // seconds
  "attempts_remaining": 3
}
```

#### Verify OTP

```
POST /api/v1/auth/otp/verify
Content-Type: application/json

{
  "request_id": "uuid-v4",
  "code": "123456",
  "device_fingerprint": "hash-of-device"
}
```

#### Response 200

```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 900,
  "token_type": "Bearer",
  "user": {
    "id": "uuid-v4",
    "role": "customer",
    "name": "Ahmed Mohamed",
    "phone": "+201012345678",
    "trust_score": 85
  }
}
```

### 2.2 Refresh Token

```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGc..."
}
```

**Note**: Old refresh token is invalidated (rotation). One-time use.

### 2.3 Employee Login (Password + TOTP)

```
POST /api/v1/auth/employee/login
Content-Type: application/json

{
  "username": "ahmed.k@food-platform.com",
  "password": "••••••••••••",
  "totp_code": "123456"
}
```

### 2.4 WebAuthn Registration (Employees)

#### Begin Registration

```
POST /api/v1/auth/webauthn/register/begin
Authorization: Bearer <jwt>
```

#### Response 200

```json
{
  "challenge": "base64url-challenge",
  "rp": { "name": "Food Platform", "id": "food-platform.com" },
  "user": {
    "id": "base64url-user-id",
    "name": "ahmed.k@food-platform.com",
    "displayName": "Ahmed K."
  },
  "pubKeyCredParams": [
    { "type": "public-key", "alg": -7 },   // ES256
    { "type": "public-key", "alg": -257 }  // RS256
  ],
  "authenticatorSelection": {
    "authenticatorAttachment": "platform",
    "userVerification": "required"
  },
  "timeout": 60000
}
```

#### Complete Registration

```
POST /api/v1/auth/webauthn/register/complete
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "credential": { /* WebAuthn response */ },
  "label": "MacBook Pro Touch ID"
}
```

### 2.5 Sensitive Action Verification

```
POST /api/v1/auth/webauthn/verify
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "credential": { /* WebAuthn assertion */ },
  "action": "refund.issue",
  "context": {
    "order_id": "uuid",
    "amount": 300
  }
}
```

#### Response 200

```json
{
  "verified": true,
  "action_token": "uuid-v4",  // use this in the actual action request
  "expires_in": 60
}
```

---

## 3. REST API — Customer Web

### 3.1 Profile

#### Get Profile

```
GET /api/v1/customers/me
Authorization: Bearer <jwt>
```

#### Response 200

```json
{
  "id": "uuid-v4",
  "name": "Ahmed Mohamed",
  "phone": "+201012345678",
  "email": "ahmed@example.com",
  "tier": "platinum",
  "trust_score": 92,
  "loyalty_points": 8520,
  "wallet_balance": 425.50,
  "addresses": [
    {
      "id": "uuid-v4",
      "label": "Home",
      "street": "12 شارع التحرير",
      "city": "Zamalek",
      "lat": 30.0626,
      "lng": 31.2197,
      "apartment": "5",
      "building": "12",
      "notes": "باب أزرق"
    }
  ]
}
```

#### Update Profile

```
PUT /api/v1/customers/me
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "name": "Ahmed M.",
  "email": "ahmed.m@example.com"
}
```

### 3.2 Addresses

```
GET    /api/v1/customers/me/addresses
POST   /api/v1/customers/me/addresses
PUT    /api/v1/customers/me/addresses/{id}
DELETE /api/v1/customers/me/addresses/{id}
```

### 3.3 Discover Restaurants

#### Get Nearby Restaurants

```
GET /api/v1/restaurants/nearby?lat=30.0626&lng=31.2197&radius=3&limit=20&offset=0
Authorization: Bearer <jwt>
```

#### Response 200

```json
{
  "restaurants": [
    {
      "id": "uuid-v4",
      "name": "Pizza Hut Maadi",
      "slug": "pizza-hut-maadi",
      "cuisine_types": ["italian", "pizza", "fast_food"],
      "rating": 4.6,
      "rating_count": 1243,
      "eta_minutes_min": 30,
      "eta_minutes_max": 40,
      "delivery_fee": 25.0,
      "is_open": true,
      "logo_url": "https://...",
      "cover_url": "https://...",
      "distance_km": 1.2,
      "promo": {
        "title": "خصم 20%",
        "description": "على الطلبات >EGP 200"
      }
    }
  ],
  "total": 47,
  "has_more": true
}
```

#### Search Restaurants

```
GET /api/v1/restaurants/search?q=pizza&lat=30.0626&lng=31.2197
```

#### Get Restaurant Detail

```
GET /api/v1/restaurants/{id}
```

#### Get Menu

```
GET /api/v1/restaurants/{id}/menu
```

#### Response 200

```json
{
  "categories": [
    {
      "id": "uuid-v4",
      "name": "بيتزا",
      "items": [
        {
          "id": "uuid-v4",
          "name": "Margherita",
          "description": "صلصة طماطم، موتزاريلا، ريحان",
          "price": 145.0,
          "image_url": "https://...",
          "is_available": true,
          "prep_time_minutes": 12,
          "modifiers": [
            {
              "id": "uuid-v4",
              "name": "حجم",
              "required": true,
              "options": [
                { "id": "uuid-v4", "name": "صغير", "price_delta": -25 },
                { "id": "uuid-v4", "name": "وسط", "price_delta": 0 },
                { "id": "uuid-v4", "name": "كبير", "price_delta": 30 }
              ]
            }
          ]
        }
      ]
    }
  ]
}
```

### 3.4 Cart

#### Add Item to Cart

```
POST /api/v1/cart/items
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "restaurant_id": "uuid-v4",
  "menu_item_id": "uuid-v4",
  "quantity": 1,
  "modifiers": [
    { "modifier_id": "uuid-v4", "option_id": "uuid-v4" }
  ],
  "notes": "بدون بصل"
}
```

#### Get Cart

```
GET /api/v1/cart
Authorization: Bearer <jwt>
```

#### Response 200

```json
{
  "restaurant_id": "uuid-v4",
  "items": [/* ... */],
  "subtotal": 550.0,
  "delivery_fee": 25.0,
  "service_fee": 27.5,
  "vat": 84.1,
  "discount": -50.0,
  "total": 636.6,
  "cashback_earned": 31.0
}
```

#### Apply Coupon

```
POST /api/v1/cart/coupon
Content-Type: application/json

{ "code": "WELCOME50" }
```

### 3.5 Orders

#### Create Order

```
POST /api/v1/orders
Authorization: Bearer <jwt>
X-Idempotency-Key: <uuid-v4>
Content-Type: application/json

{
  "address_id": "uuid-v4",
  "payment_method": "vodafone_cash",
  "scheduled_for": null,
  "notes": "اتصلوا قبل الوصول"
}
```

#### Response 201

```json
{
  "id": "uuid-v4",
  "order_number": "A7X92F",
  "status": "pending",
  "total": 636.6,
  "eta_minutes": 35,
  "payment": {
    "method": "vodafone_cash",
    "status": "pending",
    "provider_redirect_url": "https://..."
  }
}
```

#### Get Order

```
GET /api/v1/orders/{id}
```

#### Get Active Orders

```
GET /api/v1/customers/me/orders/active
```

#### Get Order History

```
GET /api/v1/customers/me/orders?limit=20&offset=0
```

#### Cancel Order

```
POST /api/v1/orders/{id}/cancel
Content-Type: application/json

{ "reason": "changed_mind" }
```

#### Rate Order

```
POST /api/v1/orders/{id}/rate
Content-Type: application/json

{
  "overall_rating": 5,
  "restaurant_rating": 5,
  "driver_rating": 4,
  "comment": "الأكل لذيذ، التوصيل اتأخر شوية"
}
```

### 3.6 Loyalty

```
GET /api/v1/customers/me/loyalty
GET /api/v1/customers/me/loyalty/transactions
GET /api/v1/customers/me/wallet
POST /api/v1/customers/me/wallet/redeem { "amount": 50 }
```

### 3.7 Payments

```
GET    /api/v1/customers/me/payment_methods
POST   /api/v1/customers/me/payment_methods
DELETE /api/v1/customers/me/payment_methods/{id}
```

---

## 4. REST API — Driver Web

### 4.1 Profile & Status

#### Get Profile

```
GET /api/v1/drivers/me
```

#### Response 200

```json
{
  "id": "uuid-v4",
  "name": "Mahmoud S.",
  "phone": "+201012345678",
  "vehicle_type": "motorcycle",
  "rating": 4.8,
  "tier": "platinum",
  "trust_score": 95,
  "acceptance_rate": 91,
  "completion_rate": 96,
  "total_earnings": 8420.50,
  "is_online": false
}
```

#### Update Status

```
PUT /api/v1/drivers/me/status
Content-Type: application/json

{ "status": "online" }  // online | offline | on_break
```

### 4.2 Orders

#### Get Available Orders (or via WebSocket)

```
GET /api/v1/drivers/orders/available
```

#### Accept Order

```
POST /api/v1/drivers/orders/{order_id}/accept
```

#### Reject Order

```
POST /api/v1/drivers/orders/{order_id}/reject
Content-Type: application/json

{ "reason": "too_far" }
```

#### Update Order Status

```
POST /api/v1/drivers/orders/{order_id}/status
Content-Type: application/json

{ "status": "picked_up", "photo_url": "https://...", "otp": "7294" }
```

Valid transitions: `picked_up` (from `ready`), `delivered` (from `picked_up`).

### 4.3 Location

```
POST /api/v1/drivers/location
Content-Type: application/json

{
  "lat": 30.0626,
  "lng": 31.2197,
  "heading": 180,
  "speed": 25.5
}
```

**Note**: For real-time, prefer WebSocket (`driver.location` event).

### 4.4 Earnings

```
GET /api/v1/drivers/earnings/today
GET /api/v1/drivers/earnings/week
GET /api/v1/drivers/earnings/breakdown?date=2026-07-04
POST /api/v1/drivers/earnings/payout
  Body: { "amount": 500, "method": "vodafone_cash" }
```

### 4.5 Heat Map

```
GET /api/v1/drivers/zones/hot?lat=30.0626&lng=31.2197
```

---

## 5. REST API — Restaurant Web

### 5.1 Profile & Hours

```
GET /api/v1/restaurants/me
PUT /api/v1/restaurants/me
GET /api/v1/restaurants/me/schedule
PUT /api/v1/restaurants/me/schedule
```

### 5.2 Orders

```
GET /api/v1/restaurants/me/orders/active
GET /api/v1/restaurants/me/orders/history?limit=20&offset=0
POST /api/v1/restaurants/me/orders/{order_id}/accept
POST /api/v1/restaurants/me/orders/{order_id}/reject
  Body: { "reason": "item_unavailable" }
POST /api/v1/restaurants/me/orders/{order_id}/start_prep
POST /api/v1/restaurants/me/orders/{order_id}/ready
POST /api/v1/restaurants/me/orders/{order_id}/delay
  Body: { "minutes": 5 }
```

### 5.3 Menu Management

```
GET    /api/v1/restaurants/me/menu
POST   /api/v1/restaurants/me/menu/items
PUT    /api/v1/restaurants/me/menu/items/{id}
DELETE /api/v1/restaurants/me/menu/items/{id}
POST   /api/v1/restaurants/me/menu/items/{id}/toggle_availability
  Body: { "is_available": false }  // 86'ing
GET    /api/v1/restaurants/me/categories
POST   /api/v1/restaurants/me/categories
PUT    /api/v1/restaurants/me/categories/{id}
```

### 5.4 Promotions

```
GET    /api/v1/restaurants/me/promotions
POST   /api/v1/restaurants/me/promotions
PUT    /api/v1/restaurants/me/promotions/{id}
DELETE /api/v1/restaurants/me/promotions/{id}
```

### 5.5 Analytics

```
GET /api/v1/restaurants/me/analytics/sales?from=2026-07-01&to=2026-07-04
GET /api/v1/restaurants/me/analytics/top_items
GET /api/v1/restaurants/me/analytics/peak_hours
GET /api/v1/restaurants/me/analytics/reviews
```

---

## 6. REST API — Support Web

### 6.1 Tickets

```
GET    /api/v1/support/tickets?status=active&assigned_to=me
GET    /api/v1/support/tickets/{id}
POST   /api/v1/support/tickets
PUT    /api/v1/support/tickets/{id}/assign
PUT    /api/v1/support/tickets/{id}/escalate
PUT    /api/v1/support/tickets/{id}/resolve
PUT    /api/v1/support/tickets/{id}/close
```

### 6.2 Messages

```
GET    /api/v1/support/tickets/{id}/messages
POST   /api/v1/support/tickets/{id}/messages
  Body: { "text": "...", "attachments": ["url1", "url2"] }
POST   /api/v1/support/tickets/{id}/messages/{id}/read
```

### 6.3 Order Actions (Scoped)

```
POST /api/v1/support/orders/{id}/cancel
POST /api/v1/support/orders/{id}/refund
  Body: { "amount": 300, "reason": "order_delayed_45min" }
POST /api/v1/support/orders/{id}/extend_eta
  Body: { "minutes": 15 }
POST /api/v1/support/orders/{id}/reassign_driver
```

### 6.4 Customer Lookup

```
GET /api/v1/support/customers/{id}
GET /api/v1/support/customers/{id}/orders
GET /api/v1/support/customers/{id}/refunds_history
GET /api/v1/support/customers/{id}/trust_score
```

### 6.5 Macros

```
GET    /api/v1/support/macros
POST   /api/v1/support/macros
PUT    /api/v1/support/macros/{id}
```

### 6.6 Knowledge Base

```
GET /api/v1/support/kb/articles
GET /api/v1/support/kb/articles/{id}
POST /api/v1/support/kb/articles  // admin only
PUT /api/v1/support/kb/articles/{id}
```

---

## 7. REST API — Command Center

### 7.1 Dashboard

```
GET /api/v1/ops/dashboard/overview
GET /api/v1/ops/dashboard/metrics
GET /api/v1/ops/dashboard/alerts
```

### 7.2 Map

```
GET /api/v1/ops/map/heatmap?layer=demand
GET /api/v1/ops/map/drivers
GET /api/v1/ops/map/orders
GET /api/v1/ops/map/zones
```

### 7.3 Zones

```
GET /api/v1/ops/zones
GET /api/v1/ops/zones/{id}
PUT /api/v1/ops/zones/{id}/surge
  Body: { "multiplier": 1.3, "duration_minutes": 60, "reason": "peak" }
```

### 7.4 Manual Interventions

```
POST /api/v1/ops/orders/{id}/assign_driver
  Body: { "driver_id": "uuid" }
POST /api/v1/ops/orders/{id}/cancel
POST /api/v1/ops/orders/{id}/refund
POST /api/v1/ops/drivers/{id}/suspend
POST /api/v1/ops/drivers/{id}/notify
  Body: { "message": "...", "bonus": 5 }
POST /api/v1/ops/restaurants/{id}/pause
  Body: { "duration_minutes": 30 }
POST /api/v1/ops/restaurants/{id}/deactivate
```

### 7.5 Incidents

```
GET    /api/v1/ops/incidents?status=active
POST   /api/v1/ops/incidents
PUT    /api/v1/ops/incidents/{id}/acknowledge
PUT    /api/v1/ops/incidents/{id}/resolve
POST   /api/v1/ops/incidents/{id}/postmortem
```

### 7.6 Forecasting

```
GET /api/v1/ops/forecast?date=tomorrow
GET /api/v1/ops/forecast/accuracy
GET /api/v1/ops/staffing/today
GET /api/v1/ops/staffing/forecast
```

### 7.7 Analytics

```
GET /api/v1/ops/analytics/realtime
GET /api/v1/ops/analytics/today
GET /api/v1/ops/analytics/export?range=2026-07-01,2026-07-04&format=csv
```

---

## 8. REST API — Employee Portal

### 8.1 Auth (Special)

```
POST /api/v1/auth/employee/login
  Body: { "username": "...", "password": "...", "totp_code": "..." }
POST /api/v1/auth/employee/refresh
POST /api/v1/auth/employee/logout
POST /api/v1/auth/webauthn/register/begin
POST /api/v1/auth/webauthn/register/complete
POST /api/v1/auth/webauthn/verify
```

### 8.2 Profile & Permissions

```
GET /api/v1/employees/me
PUT /api/v1/employees/me/password
GET /api/v1/employees/me/permissions
GET /api/v1/employees/me/audit_log
```

### 8.3 Sensitive Actions (Require Biometric Token)

All requests must include `X-Action-Token` header from WebAuthn verification.

```
POST /api/v1/actions/refund
  Headers: X-Action-Token: <token>
  Body: { "order_id": "uuid", "amount": 300, "reason": "..." }

POST /api/v1/actions/restaurant/approve
POST /api/v1/actions/restaurant/deactivate
POST /api/v1/actions/driver/approve
POST /api/v1/actions/payout
  Body: { "recipient": "uuid", "amount": 5000, "type": "driver" }
POST /api/v1/actions/menu/override
POST /api/v1/actions/data/export
```

### 8.4 Dual Approval

```
GET  /api/v1/approvals/pending
POST /api/v1/approvals/{id}/approve
  Headers: X-Action-Token: <token>
POST /api/v1/approvals/{id}/reject
  Body: { "reason": "..." }
```

### 8.5 Audit Log (Admin Only)

```
GET /api/v1/audit/logs?filters...
GET /api/v1/audit/logs/{id}
GET /api/v1/audit/verify_chain
POST /api/v1/audit/export
  Body: { "date_range": ["2026-07-01", "2026-07-04"], "format": "csv" }
```

### 8.6 Access Review

```
GET  /api/v1/access-review/current
GET  /api/v1/access-review/team
POST /api/v1/access-review/{employee_id}/update
  Body: { "permissions": ["..."] }
```

### 8.7 Anomaly Alerts

```
GET  /api/v1/anomalies/active
POST /api/v1/anomalies/{id}/investigate
POST /api/v1/anomalies/{id}/resolve
```

### 8.8 Admin (Super Admin Only)

```
GET  /api/v1/admin/employees
POST /api/v1/admin/employees
PUT  /api/v1/admin/employees/{id}/role
POST /api/v1/admin/employees/{id}/suspend
POST /api/v1/admin/employees/{id}/deactivate
```

---

## 9. REST API — Field Supervisor

### 9.1 Tasks

```
GET  /api/v1/supervisor/tasks/today
GET  /api/v1/supervisor/tasks/{id}
POST /api/v1/supervisor/tasks/{id}/start
POST /api/v1/supervisor/tasks/{id}/complete
  Body: {
    "result": "approved" | "conditional" | "rejected",
    "checklist": { /* ... */ },
    "photos": ["url1", "url2"],
    "notes": "..."
  }
```

### 9.2 Restaurant Verification

```
POST /api/v1/supervisor/restaurants/{id}/verify
POST /api/v1/supervisor/restaurants/{id}/photos
  Body: {
    "photo_url": "https://...",
    "type": "storefront" | "interior" | "kitchen" | "license",
    "gps": { "lat": 30.0444, "lng": 31.2357 },
    "timestamp": "2026-07-04T11:15:32Z"
  }
POST /api/v1/supervisor/restaurants/{id}/decision
  Body: { "approved": true, "conditions": [...], "reason": "..." }
```

### 9.3 Driver Verification

```
POST /api/v1/supervisor/drivers/{id}/verify
POST /api/v1/supervisor/drivers/{id}/photos
POST /api/v1/supervisor/drivers/{id}/decision
```

### 9.4 Audits

```
POST /api/v1/supervisor/audits
  Body: { "restaurant_id": "uuid", "type": "surprise" }
GET  /api/v1/supervisor/audits/{id}
```

### 9.5 Complaints

```
GET  /api/v1/supervisor/complaints?status=assigned_to_me
POST /api/v1/supervisor/complaints/{id}/investigate
  Body: { "findings": "...", "photos": [...], "interviews": [...] }
POST /api/v1/supervisor/complaints/{id}/resolve
  Body: { "resolution": "...", "actions": [...] }
```

### 9.6 Route Planning

```
GET  /api/v1/supervisor/route/optimize?task_ids=uuid1,uuid2
POST /api/v1/supervisor/route/start
```

### 9.7 Location

```
POST /api/v1/supervisor/location
  Body: { "lat": 30.0444, "lng": 31.2357 }
GET  /api/v1/supervisor/location/history
```

### 9.8 Reports

```
GET  /api/v1/supervisor/reports/daily
GET  /api/v1/supervisor/reports/weekly
POST /api/v1/supervisor/reports/{task_id}/submit
```

### 9.9 Training

```
GET  /api/v1/supervisor/training/modules
POST /api/v1/supervisor/training/{driver_id}/complete
  Body: { "module_id": "uuid", "score": 18 }
```

---

## 10. WebSocket Events

### 10.1 Connection

```
WSS: wss://api.food-platform.com/ws
Headers: Authorization: Bearer <jwt>
Subprotocol: food-platform.v1
```

### 10.2 Message Format

```json
{
  "event": "string",
  "payload": { /* event-specific */ },
  "timestamp": "2026-07-04T14:32:15Z",
  "id": "uuid-v4"
}
```

### 10.3 Customer Events

#### Inbound (server → client)

| Event | Payload | Trigger |
|-------|---------|---------|
| `order.status_changed` | `{ order_id, status, eta_minutes }` | Order status update |
| `driver.location` | `{ order_id, lat, lng, heading }` | Every 5s when order active |
| `driver.assigned` | `{ order_id, driver: {name, photo, rating, vehicle} }` | Driver accepts |
| `order.ready_for_pickup` | `{ order_id }` | Restaurant marks ready |
| `promo.unlocked` | `{ code, discount }` | Customer qualifies |
| `support.message` | `{ ticket_id, message, agent_name }` | Support replies |

#### Outbound (client → server)

| Event | Payload |
|-------|---------|
| `ping` | `{}` (heartbeat every 30s) |
| `subscribe` | `{ channel: "order.{id}" }` |

### 10.4 Driver Events

#### Inbound

| Event | Payload |
|-------|---------|
| `order.new` | `{ order_id, restaurant, distance, earnings, eta }` |
| `order.cancelled` | `{ order_id, reason }` |
| `order.picked_up_confirmed` | `{ order_id }` |
| `support.message` | `{ ticket_id, message }` |
| `payout.completed` | `{ amount, method, transaction_id }` |
| `zone.hot` | `{ zones: [...] }` |

#### Outbound

| Event | Payload |
|-------|---------|
| `driver.location` | `{ lat, lng, heading, speed }` (every 5s) |
| `driver.status_changed` | `{ status }` |

### 10.5 Restaurant Events

#### Inbound

| Event | Payload |
|-------|---------|
| `order.new` | `{ order_id, customer, items, total, payment_method }` |
| `order.cancelled` | `{ order_id, reason }` |
| `order.picked_up` | `{ order_id, driver_name }` |
| `order.delivered` | `{ order_id }` |
| `review.new` | `{ order_id, rating, comment }` |
| `promo.performance_update` | `{ promo_id, orders, revenue }` |

#### Outbound

| Event | Payload |
|-------|---------|
| `order.accepted` | `{ order_id }` |
| `order.rejected` | `{ order_id, reason }` |
| `order.started_prep` | `{ order_id }` |
| `order.ready` | `{ order_id }` |
| `menu.updated` | `{ items: [...] }` |

### 10.6 Support Events

| Event | Payload |
|-------|---------|
| `ticket.new` | `{ ticket_id, customer, priority }` |
| `ticket.message` | `{ ticket_id, message }` |
| `ticket.escalated` | `{ ticket_id, from, to }` |
| `alert.fraud` | `{ customer_id, score, reasons }` |

### 10.7 Command Center Events

| Event | Payload |
|-------|---------|
| `metric.update` | `{ metric, value, timestamp }` |
| `alert.new` | `{ severity, message, zone? }` |
| `incident.new` | `{ incident_id, severity, title }` |
| `incident.update` | `{ incident_id, status }` |
| `zone.update` | `{ zone_id, demand, drivers }` |

### 10.8 Employee Portal Events

| Event | Payload |
|-------|---------|
| `approval.requested` | `{ approval_id, action, requester }` |
| `approval.resolved` | `{ approval_id, decision }` |
| `anomaly.detected` | `{ employee_id, score, reasons }` |
| `audit.tamper_detected` | `{ record_id }` (P0 alert) |

### 10.9 Field Supervisor Events

| Event | Payload |
|-------|---------|
| `task.assigned` | `{ task_id, type, priority, location }` |
| `task.cancelled` | `{ task_id, reason }` |
| `complaint.assigned` | `{ complaint_id, severity, restaurant }` |

---

## 11. Kafka Events (AsyncAPI)

### 11.1 Topic Naming Convention

```
{domain}.{entity}.{action}

Examples:
  order.created
  order.status_changed
  payment.captured
  driver.location
```

### 11.2 Event Schema (Avro)

All events use Avro schemas registered in Confluent Schema Registry.

#### Order Created

```json
{
  "type": "record",
  "name": "OrderCreated",
  "namespace": "com.foodplatform.order",
  "fields": [
    { "name": "event_id", "type": "string" },
    { "name": "event_time", "type": "string" },
    { "name": "order_id", "type": "string" },
    { "name": "user_id", "type": "string" },
    { "name": "restaurant_id", "type": "string" },
    { "name": "total_amount", "type": "double" },
    { "name": "payment_method", "type": "string" },
    { "name": "delivery_address", "type": "string" },
    { "name": "items_count", "type": "int" },
    { "name": "metadata", "type": "string" }
  ]
}
```

### 11.3 Core Topics

| Topic | Producer | Consumers | Partitions | Retention |
|-------|----------|-----------|-----------|-----------|
| `order.created` | Order | Fraud, Notification, Analytics, Loyalty | 12 | 7 days |
| `order.confirmed` | Delivery Matching | Notification, Restaurant, Customer | 12 | 7 days |
| `order.cancelled` | Order, Support | Payment (refund), Driver (unassign) | 12 | 30 days |
| `order.status_changed` | Order, Driver | Notification, Analytics | 12 | 30 days |
| `payment.captured` | Payment | Order, Loyalty, Analytics | 6 | 30 days |
| `payment.failed` | Payment | Order, Notification | 6 | 30 days |
| `payment.refunded` | Payment | Order, Notification, Analytics | 6 | 30 days |
| `driver.location` | Driver App | Geo Service, Command Center | 24 | 3 days |
| `driver.status_changed` | Driver Mgmt | Delivery Matching, Analytics | 6 | 7 days |
| `restaurant.menu_updated` | Restaurant App | Catalog, Search | 6 | 7 days |
| `fraud.score_calculated` | Fraud | Order, Analytics | 6 | 30 days |
| `audit.action_logged` | Trust | Analytics | 6 | 90 days |

### 11.4 Event Order Guarantees

- **Per-order**: All events for the same `order_id` go to the same partition → guaranteed order.
- **Per-driver**: All `driver.location` events for the same driver → same partition.
- **Cross-entity**: No ordering guarantee.

### 11.5 Consumer Patterns

#### At-Least-Once (default)
- Most consumers
- Must be idempotent (use `event_id` for deduplication)

#### Exactly-Once (Payments)
- Payment Service uses Kafka transactions
- Consumer reads committed messages only

---

## 12. Common Schemas

### 12.1 Money

```typescript
type Money = string;  // decimal as string: "245.50"
```

### 12.2 Coordinates

```typescript
interface Coordinates {
  lat: number;  // decimal degrees, 7 places
  lng: number;
}
```

### 12.3 Address

```typescript
interface Address {
  id: string;
  label: string;       // "Home" | "Work" | "Other"
  street: string;
  city: string;
  lat: number;
  lng: number;
  apartment?: string;
  building?: string;
  notes?: string;
}
```

### 12.4 Pagination

```typescript
interface PaginatedResponse<T> {
  data: T[];
  total: number;
  has_more: boolean;
  next_offset?: number;
}
```

### 12.5 Order Status Enum

```typescript
type OrderStatus =
  | "pending"
  | "confirmed"
  | "preparing"
  | "ready"
  | "picked_up"
  | "delivered"
  | "cancelled"
  | "refunded";
```

### 12.6 User Roles

```typescript
type UserRole =
  | "customer"
  | "driver"
  | "restaurant"
  | "support_l1"
  | "support_l2"
  | "ops_manager"
  | "finance"
  | "super_admin"
  | "field_supervisor"
  | "hr"
  | "read_only_analyst";
```

---

## 13. Error Codes

### 13.1 Error Response Format

```json
{
  "error": {
    "code": "ORDER_NOT_FOUND",
    "message": "Order with id 'uuid' not found",
    "details": { "order_id": "uuid" },
    "request_id": "uuid-v4",
    "documentation_url": "https://docs.food-platform.com/errors/ORDER_NOT_FOUND"
  }
}
```

### 13.2 HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | Success | GET resource |
| 201 | Created | POST creates resource |
| 204 | No Content | DELETE success |
| 400 | Bad Request | Invalid input |
| 401 | Unauthorized | Missing/invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate idempotency key |
| 422 | Unprocessable Entity | Validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Backend bug |
| 503 | Service Unavailable | Maintenance / overload |

### 13.3 Business Error Codes

| Code | HTTP | Description |
|------|------|-------------|
| `AUTH_INVALID_OTP` | 401 | OTP incorrect or expired |
| `AUTH_OTP_ATTEMPTS_EXCEEDED` | 429 | Too many OTP attempts |
| `AUTH_TOKEN_EXPIRED` | 401 | Access token expired |
| `AUTH_REFRESH_INVALID` | 401 | Refresh token invalid or used |
| `AUTH_BIOMETRIC_REQUIRED` | 403 | Sensitive action requires biometric |
| `AUTH_ACTION_TOKEN_EXPIRED` | 401 | Action token expired (60s) |
| `ORDER_NOT_FOUND` | 404 | Order ID doesn't exist |
| `ORDER_INVALID_TRANSITION` | 422 | Invalid status change |
| `ORDER_ALREADY_CANCELLED` | 409 | Order already cancelled |
| `ORDER_CANNOT_CANCEL` | 422 | Cancellation window closed |
| `PAYMENT_DECLINED` | 422 | Payment provider declined |
| `PAYMENT_DUPLICATE` | 409 | Idempotency key reused |
| `RESTAURANT_CLOSED` | 422 | Restaurant not accepting orders |
| `MENU_ITEM_UNAVAILABLE` | 422 | Item 86'd |
| `CART_EMPTY` | 422 | Cannot checkout empty cart |
| `CART_DIFFERENT_RESTAURANTS` | 422 | Items from multiple restaurants |
| `PROMO_INVALID` | 422 | Coupon invalid/expired |
| `PROMO_USAGE_LIMIT_REACHED` | 422 | Coupon fully redeemed |
| `DRIVER_NOT_AVAILABLE` | 422 | No drivers in area |
| `DRIVER_OFFLINE` | 422 | Driver went offline |
| `REFUND_AMOUNT_EXCEEDS_ORDER` | 422 | Refund > order total |
| `REFUND_WINDOW_CLOSED` | 422 | Refund period expired |
| `FRAUD_DETECTED` | 403 | Order blocked by fraud service |
| `FRAUD_REVIEW_REQUIRED` | 202 | Order sent for manual review |
| `RATE_LIMIT_EXCEEDED` | 429 | Rate limit hit |
| `PERMISSION_DENIED` | 403 | Role lacks permission |
| `DUAL_APPROVAL_REQUIRED` | 403 | Action needs second approver |
| `AUDIT_CHAIN_BROKEN` | 500 | Tampering detected (P0) |
| `WEBHOOK_DELIVERY_FAILED` | 200 | Background retry scheduled |

### 13.4 Error Handling Strategy (Client Side)

```typescript
// Pseudocode
try {
  await api.createOrder(payload);
} catch (error) {
  if (error.code === 'AUTH_TOKEN_EXPIRED') {
    await refreshToken();
    retry();
  } else if (error.code === 'PAYMENT_DECLINED') {
    showPaymentError(error.details);
  } else if (error.code === 'FRAUD_DETECTED') {
    showFraudWarning();
    logAnalytics('fraud_blocked', error);
  } else if (error.status === 429) {
    await sleep(getRetryAfter(error));
    retry();
  } else {
    showGenericError(error);
  }
}
```

---

## Appendix A: Code Generation

### From OpenAPI to TypeScript

```bash
openapi-typescript ./openapi.yaml -o ./packages/types/src/generated/api.ts
```

### From OpenAPI to Go

```bash
oapi-codegen -package api -generate types,client ./openapi.yaml > ./services/shared/openapi/types.go
```

### From AsyncAPI to Kafka schemas

```bash
asyncapi-generator @asyncapi/avro-schema ./asyncapi.yaml -o ./tools/schemas/
```

---

> **Next**: Read `PATTERNS.md` for code patterns, `REPO-STRUCTURE.md` for monorepo layout, and `REFERENCES.md` for research sources.
