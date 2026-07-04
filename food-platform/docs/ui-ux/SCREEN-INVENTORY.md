# Screen Inventory — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04

Complete catalog of all screens across the 7 web apps. Each screen has a unique ID for reference.

---

## Summary

| App | Screens | Priority |
|-----|---------|----------|
| Customer Web | 22 | P0 (MVP) |
| Driver Web | 14 | P0 (MVP) |
| Restaurant Web | 18 | P0 (MVP) |
| Support Web | 12 | P1 (Phase 2) |
| Command Center | 10 | P1 (Phase 2) |
| Employee Portal | 12 | P2 (Phase 3) |
| Field Supervisor Web | 11 | P2 (Phase 3) |
| **Total** | **99 screens** | |

---

## 1. Customer Web App (22 screens)

### 1.1 Auth Flow (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-AUTH-01 | Welcome / Onboarding | First-time intro carousel | P0 |
| CUS-AUTH-02 | Phone Login | Enter phone number | P0 |
| CUS-AUTH-03 | OTP Verification | 6-digit code entry | P0 |
| CUS-AUTH-04 | Profile Setup | Name, email (optional) | P0 |

### 1.2 Home & Discovery (5 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-HOME-01 | Home / Discovery | Main feed with restaurants | P0 |
| CUS-HOME-02 | Address Picker | Select delivery address | P0 |
| CUS-HOME-03 | Search Results | Search restaurants/dishes | P0 |
| CUS-HOME-04 | Category Browse | Browse by cuisine type | P0 |
| CUS-HOME-05 | Filter & Sort | Filter restaurants | P1 |

### 1.3 Restaurant & Menu (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-REST-01 | Restaurant Detail | Menu + info + reviews | P0 |
| CUS-REST-02 | Item Detail / Customize | Modifiers, quantity, notes | P0 |
| CUS-REST-03 | Reviews List | All reviews for restaurant | P1 |

### 1.4 Cart & Checkout (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-CART-01 | Cart | Items + pricing + coupon | P0 |
| CUS-CART-02 | Checkout | Address + payment + place order | P0 |
| CUS-CART-03 | Payment Method | Select Vodafone Cash / InstaPay / Card / COD | P0 |
| CUS-CART-04 | Order Confirmation | Success screen after order | P0 |

### 1.5 Orders (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-ORDER-01 | Order Tracking | Live status + map + ETA | P0 |
| CUS-ORDER-02 | Order History | List of past orders | P0 |
| CUS-ORDER-03 | Order Detail | Past order details + reorder | P0 |

### 1.6 Profile (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-PROF-01 | Profile | User info + stats | P0 |
| CUS-PROF-02 | Addresses | CRUD addresses | P0 |
| CUS-PROF-03 | Payment Methods | Saved cards/wallets | P1 |

### 1.7 Loyalty (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-LOYL-01 | Loyalty & Wallet | Points, tier, cashback, wallet | P1 |

### 1.8 Support (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| CUS-SUPP-01 | Help / Support | Chat with support + FAQ | P0 |

---

## 2. Driver Web App (14 screens)

### 2.1 Auth & Onboarding (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| DRV-AUTH-01 | Phone Login | Enter phone | P0 |
| DRV-AUTH-02 | OTP Verification | 6-digit code | P0 |
| DRV-AUTH-03 | Profile Setup | Vehicle type, license upload | P0 |
| DRV-AUTH-04 | Training (basic) | App usage tutorial | P1 |

### 2.2 Home & Status (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| DRV-HOME-01 | Home (Online) | Earnings, heat map, status | P0 |
| DRV-HOME-02 | Heat Map | Demand zones visualization | P0 |
| DRV-HOME-03 | Profile | Driver info, tier, rating | P0 |

### 2.3 Order Flow (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| DRV-ORDR-01 | Order Offer Modal | 15s countdown, accept/reject | P0 |
| DRV-ORDR-02 | Pickup Screen | Navigation to restaurant + items | P0 |
| DRV-ORDR-03 | Dropoff Screen | Navigation to customer + OTP | P0 |
| DRV-ORDR-04 | Order Completed | Earnings summary + next | P0 |

### 2.4 Earnings (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| DRV-EARN-01 | Earnings Dashboard | Today/week + breakdown | P0 |
| DRV-EARN-02 | Payout | Instant/daily/weekly options | P0 |

### 2.5 History (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| DRV-HIST-01 | Order History | Past deliveries + ratings | P0 |

---

## 3. Restaurant Web App (18 screens)

### 3.1 Auth (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-AUTH-01 | Phone Login | Enter phone | P0 |
| REST-AUTH-02 | OTP Verification | 6-digit code | P0 |

### 3.2 Dashboard (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-DASH-01 | Main Dashboard | Today's stats + active orders | P0 |
| REST-DASH-02 | Schedule / Hours | Operating hours management | P0 |

### 3.3 Order Management (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-ORDR-01 | Active Orders | List of current orders | P0 |
| REST-ORDR-02 | Inbound Order Modal | 90s timer, accept/reject | P0 |
| REST-ORDR-03 | Order Detail | Items, customer, payment | P0 |
| REST-ORDR-04 | Order History | Past orders + filter | P0 |

### 3.4 Menu Management (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-MENU-01 | Menu Overview | Categories + items list | P0 |
| REST-MENU-02 | Item Editor | Add/edit item + modifiers | P0 |
| REST-MENU-03 | Category Editor | Manage categories | P0 |
| REST-MENU-04 | Bulk Availability | 86'ing items | P0 |

### 3.5 Kitchen Display (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-KDS-01 | KDS (Kitchen Display) | Full-screen order board | P1 |

### 3.6 Analytics (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-ANAL-01 | Sales Analytics | Charts + top items | P1 |
| REST-ANAL-02 | Reviews | Customer reviews + reply | P1 |
| REST-ANAL-03 | Peak Hours | Heatmap of orders by hour | P1 |

### 3.7 Promotions (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| REST-PROMO-01 | Promotions List | Active/scheduled/expired | P1 |
| REST-PROMO-02 | Create Promotion | 4 promo types | P1 |

---

## 4. Support Web App (12 screens)

### 4.1 Auth (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| SUP-AUTH-01 | Employee Login | Username + password + TOTP | P0 |
| SUP-AUTH-02 | 2FA Setup | TOTP enrollment | P0 |

### 4.2 Dashboard (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| SUP-DASH-01 | Main Dashboard | Today's stats + active tickets | P0 |
| SUP-DASH-02 | Team Status | Other agents online + load | P1 |

### 4.3 Ticket Management (5 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| SUP-TICK-01 | Ticket Queue | Sortable list of tickets | P0 |
| SUP-TICK-02 | Ticket Detail | Chat + customer info + actions | P0 |
| SUP-TICK-03 | Refund Dialog | Biometric + dual approval | P0 |
| SUP-TICK-04 | Escalation | Send to Tier 2 / Ops | P0 |
| SUP-TICK-05 | Macros | Canned response library | P1 |

### 4.4 Knowledge Base (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| SUP-KB-01 | KB Search | Article search | P1 |
| SUP-KB-02 | KB Article | Read article | P1 |

### 4.5 Customer Lookup (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| SUP-CUST-01 | Customer 360 | Customer profile + orders + refunds | P0 |

---

## 5. Command Center (10 screens)

### 5.1 Auth (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| OPS-AUTH-01 | Login | Username + password + TOTP | P0 |

### 5.2 Dashboard (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| OPS-DASH-01 | Main Dashboard | KPIs + alerts | P0 |
| OPS-DASH-02 | Live Map | Heatmap + drivers + orders | P0 |
| OPS-DASH-03 | Real-time Metrics | Orders/min, GMV, etc. | P0 |

### 5.3 Zone Management (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| OPS-ZONE-01 | Zone Detail | Per-zone metrics + actions | P0 |
| OPS-ZONE-02 | Surge Control | Manual surge override | P0 |

### 5.4 Incident Management (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| OPS-INC-01 | Incident List | Active + recent incidents | P0 |
| OPS-INC-02 | Incident Detail | Acknowledge/resolve/postmortem | P0 |

### 5.5 Manual Interventions (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| OPS-INTV-01 | Manual Dispatch | Assign driver to order | P0 |
| OPS-INTV-02 | Restaurant/Driver Control | Suspend/pause | P0 |

---

## 6. Employee Portal (12 screens)

### 6.1 Auth (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| EMP-AUTH-01 | Login | Username + password + TOTP | P0 |
| EMP-AUTH-02 | WebAuthn Enrollment | Register biometric | P0 |
| EMP-AUTH-03 | Sensitive Action Verify | Biometric re-auth | P0 |

### 6.2 Dashboard (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| EMP-DASH-01 | Main Dashboard | Pending approvals + stats | P0 |
| EMP-DASH-02 | Personal Activity | Own audit log | P0 |

### 6.3 Sensitive Actions (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| EMP-ACT-01 | Refund Action | Process refund with biometric | P0 |
| EMP-ACT-02 | Restaurant Approve | Approve new restaurant | P0 |
| EMP-ACT-03 | Driver Approve | Approve new driver | P0 |
| EMP-ACT-04 | Payout Approve | Approve driver/restaurant payout | P0 |

### 6.4 Audit & Compliance (3 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| EMP-AUD-01 | Audit Log Viewer | Search/filter all actions | P0 |
| EMP-AUD-02 | Anomaly Alerts | Flagged employees | P1 |
| EMP-AUD-03 | Access Review | Quarterly review | P1 |

---

## 7. Field Supervisor Web App (11 screens)

### 7.1 Auth (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| FLD-AUTH-01 | Phone Login | Enter phone | P0 |
| FLD-AUTH-02 | OTP Verification | 6-digit code | P0 |

### 7.2 Tasks (4 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| FLD-TASK-01 | Task List | Today's tasks sorted by priority | P0 |
| FLD-TASK-02 | Restaurant Verification | 50+ point checklist + photos | P0 |
| FLD-TASK-03 | Driver Verification | Documents + photos + practical test | P0 |
| FLD-TASK-04 | Complaint Investigation | Evidence collection + report | P0 |

### 7.3 Route (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| FLD-ROUTE-01 | Route Planner | TSP optimization for day | P0 |
| FLD-ROUTE-02 | Navigation | Google Maps integration | P0 |

### 7.4 Reports (2 screens)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| FLD-REPT-01 | Daily Report | Submit end-of-day report | P0 |
| FLD-REPT-02 | Performance | Weekly stats + ranking | P1 |

### 7.5 Training (1 screen)

| ID | Screen | Purpose | Priority |
|----|--------|---------|----------|
| FLD-TRNG-01 | Driver Training | Deliver training modules | P1 |

---

## Screen ID Convention

Format: `{APP_PREFIX}-{CATEGORY}-{NUMBER}`

| App | Prefix | Example |
|-----|--------|---------|
| Customer Web | CUS | CUS-HOME-01 |
| Driver Web | DRV | DRV-ORDR-01 |
| Restaurant Web | REST | REST-MENU-02 |
| Support Web | SUP | SUP-TICK-03 |
| Command Center | OPS | OPS-INC-02 |
| Employee Portal | EMP | EMP-AUD-01 |
| Field Supervisor | FLD | FLD-TASK-02 |

---

## Priority Definitions

| Priority | Meaning | Phase |
|----------|---------|-------|
| P0 | MVP — must have for launch | Phase 1-2 (Weeks 1-9) |
| P1 | Important — should have for launch | Phase 3-5 (Weeks 6-11) |
| P2 | Nice to have — post-launch | Phase 6+ (Week 12+) |

---

> **Next**: Read per-app `UI-SPEC.md` files for detailed screen specifications.
