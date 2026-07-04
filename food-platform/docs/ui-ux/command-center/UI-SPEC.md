# Command Center — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: Uber Charon, Datadog, Grafana research

---

## 1. App Overview

### Purpose

Command Center enables operations managers to:
- Monitor real-time metrics (orders, drivers, restaurants, GMV)
- View live map with demand heatmap, driver locations, order flows
- Manage zones (surge pricing, driver notifications)
- Handle incidents (P0/P1/P2)
- Manual interventions (assign driver, cancel order, suspend)
- View forecasts and staffing recommendations

### Design Inspiration

| Source | What We Take |
|--------|--------------|
| Uber Charon | Real-time analytics, merchant monitoring |
| Datadog | Dashboard layout, metric visualization |
| Grafana | Time-series charts, alerting |
| PagerDuty | Incident management workflow |

### Design Principles

1. **Real-Time First**: Updates <5s lag
2. **Glanceable**: Key metrics visible without scrolling
3. **Action-Oriented**: Every metric has a related action
4. **Dark Mode Default**: Better for ops rooms (24/7 monitoring)
5. **Multi-Screen Ready**: Designed for ops room with 4 large displays

---

## 2. Screen Specifications

### OPS-DASH-01: Main Dashboard ⭐

```
┌─────────────────────────────────────────────────────────────────────┐
│  🚨 Command Center | Shift: 8AM-4PM | Today: 4 يوليو 2026            │
│  👤 Omar T. (Ops Manager) | 🟢 Online                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐  │
│  │ 📦 طلبات │ │ 🛵 مناديب│ │ 🍔 مطاعم│ │ ⚠️ مشاكل │ │ 💰 GMV  │  │
│  │  نشطة    │ │  نشطة   │ │  نشطة   │ │  نشطة   │ │  اليوم  │  │
│  │   423    │ │   187    │ │   142    │ │    3     │ │ EGP 84K │  │
│  │   ↑12%   │ │   ↑8%    │ │   ↑5%    │ │   🔴     │ │  ↑18%   │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘  │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  📍 Live Map (القاهرة الكبرى)                                        │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                                                              │    │
│  │                    🔴 معادي                                   │    │
│  │                   🔴🔴🔴                                     │    │
│  │                🔴🔴🔴🔴                                     │    │
│  │             🟡🟡🟡🔴🔴🔴     الزمالك                        │    │
│  │            🟡🟡🟡🟡🔴🔴🔴🔴🔴                               │    │
│  │             🟡🟡🟡🟡🟡     مدينة نصر                         │    │
│  │              🟡🟡🟡🟢🟢🟢                                    │    │
│  │               🟢🟢🟢🟢🟢🟢  التحرير                          │    │
│  │                                                              │    │
│  │  🔴 طلب عالي | 🟡 متوسط | 🟢 هادي | 🛵 مندوب | 🍔 مطعم    │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  📊 Live Metrics                                                    │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │  Orders/min: ████████████ 7.2 (target: 8.0)                │    │
│  │  Avg Delivery Time: ████████ 32 min (target: <35)          │    │
│  │  Driver Utilization: ██████████████ 78% (target: 75-85%)   │    │
│  │  Order Completion: ███████████████ 94% (target: >92%)      │    │
│  │  Cancellation Rate: ██ 4% (target: <5%)                    │    │
│  │  Payment Success: █████████████████ 98.5% (target: >98%)   │    │
│  │  Avg Driver Rating: ████████████ 4.7 (target: >4.5)        │    │
│  │  P95 API Latency: ████████ 380ms (target: <500ms)          │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ⚠️ Active Alerts (3)                                               │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │ 🔴 P0: معادي - طلبات عالية + نقص مناديب (gap: 18 طلب)       │    │
│  │    [تفعيل Surge 1.3x] [طلب مناديب] [Ignore]                │    │
│  │                                                              │    │
│  │ 🟡 P1: Pizza Hut معادي - 5 طلبات delayed >15 دقيقة          │    │
│  │    [اتصال بالمطعم] [Auto-reject] [Ignore]                  │    │
│  │                                                              │    │
│  │ 🟢 P2: Vodafone Cash API - latency مرتفعة (3 ثواني)         │    │
│  │    [Switch to InstaPay] [Monitor] [Ignore]                 │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### OPS-DASH-02: Live Map (Full Screen)

```
┌─────────────────────────────────────────────────────────────────────┐
│  📍 Live Map | 4 يوليو 2026, 14:32                  [⚙️ Layers] [✕] │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                                                              │    │
│  │                    🔴🔴🔴 معادي                               │    │
│  │                  🔴🔴🔴🔴🔴                                  │    │
│  │               🟡🟡🟡🔴🔴🔴🔴     الزمالك                     │    │
│  │              🟡🟡🟡🟡🔴🔴🔴🔴🔴🔴                            │    │
│  │               🟡🟡🟡🟡🟡🟡     مدينة نصر                      │    │
│  │                🟡🟡🟡🟢🟢🟢🟢                                 │    │
│  │                 🟢🟢🟢🟢🟢🟢🟢🟢  التحرير                      │    │
│  │                                                              │    │
│  │     🛵 🛵 🛵 🛵 🛵 🛵 🛵 🛵 🛵 🛵 (drivers)                  │    │
│  │                                                              │    │
│  │     Layers: [☑ Demand] [☑ Drivers] [☐ Orders] [☐ Problems]│    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                     │
│  🔴 طلب عالي جداً | 🟡 متوسط | 🟢 هادي | 🛵 مندوب متاح              │
│                                                                     │
│  آخر تحديث: 14:32:15 (5 ثواني)                                      │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Map Layers**:
1. Demand Heat Map (default)
2. Driver Density
3. Order Flow (lines from restaurant to customer)
4. Problem Zones (red circles)
5. Forecast Overlay

### OPS-ZONE-01: Zone Detail (Maadi)

```
┌─────────────────────────────────────────────────────┐
│  📍 معادي - تفاصيل المنطقة                          │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │
│  │ 📦 طلبات │ │ 🛵 مناديب│ │ ⏱️ Avg   │ │ ⚠️ مشاكل│ │
│  │   72     │ │   31     │ │ 28 دقيقة │ │   1    │ │
│  └──────────┘ └──────────┘ └──────────┘ └────────┘ │
│                                                     │
│  Demand vs Supply:                                  │
│  ┌──────────────────────────────────────────────┐   │
│  │ الطلب:    ████████████████ 72 طلب            │   │
│  │ المناديب: ████████ 31 مندوب (نقص 8)          │   │
│  │ Gap:      -8 (تحتاج 8 مناديب إضافيين)         │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  الإجراءات المتاحة:                                   │
│  [⚡ تفعيل Surge 1.3x]                              │
│  [📢 إرسال مناديب من مناطق قريبة]                   │
│  [💰 Bonus للمناديب اللي ييجوا للمنطقة]             │
│  [🚫 تعليق الطلبات الجديدة مؤقتاً]                  │
│                                                     │
│  آخر ساعة:                                           │
│  ┌──────────────────────────────────────────────┐   │
│  │  12pm  1pm  2pm  3pm  4pm                    │   │
│  │  ────────────────────────────                │   │
│  │  الطلب  ████████████                         │   │
│  │  المناديب ████████                           │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### OPS-INC-01: Incident List

```
┌─────────────────────────────────────────────────────┐
│  🚨 Incidents (آخر 24 ساعة)                          │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Active (1):                                        │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🔴 P0: Vodafone Cash API Down               │   │
│  │ Started: 1:00 PM | Duration: 15 min         │   │
│  │ Owner: Ahmed K. (On-call)                   │   │
│  │ Status: Mitigating                          │   │
│  │ [فتح Incident]                              │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  Recent (5):                                        │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🟡 P1: معادي - driver shortage              │   │
│  │ أمس 8:15 PM | Duration: 22 min | Resolved   │   │
│  └──────────────────────────────────────────────┘   │
│  ...                                                │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 3. UX Flows

### Surge Activation Flow

```
1. Alert: "معادي - نقص مناديب"
   ↓
2. Tap "تفعيل Surge"
   ↓
3. Surge dialog:
   - Multiplier: 1.3x (slider 1.0-1.5)
   - Duration: 60 min (dropdown)
   - Reason: peak + driver shortage
   ↓
4. Tap "تطبيق"
   ↓
5. Surge active in zone
   ↓
6. Customer app shows "أسعار مرتفعة" in that zone
   ↓
7. Drivers see surge indicator → move to zone
```

### Incident Response Flow

```
1. Auto-detection (e.g., payment success rate drops)
   ↓
2. Alert appears in dashboard
   ↓
3. Ops Manager acknowledges (assigns to on-call)
   ↓
4. On-call investigates
   ↓
5. Identifies root cause
   ↓
6. Applies mitigation (e.g., switch payment provider)
   ↓
7. Monitors recovery
   ↓
8. Resolves incident
   ↓
9. Postmortem (next day)
```

---

## 4. Critical UX Considerations

### 4.1 Real-Time Updates

- WebSocket connection for live metrics
- Map updates every 5s
- Alerts appear instantly
- No page refreshes needed

### 4.2 Multi-Screen Setup

For ops rooms with 4 displays:
- Screen 1: Live Map (always visible)
- Screen 2: KPIs Dashboard
- Screen 3: Active Alerts
- Screen 4: Forecast

### 4.3 Dark Mode

- Default (better for 24/7 monitoring)
- Reduced eye strain
- Higher contrast for charts

### 4.4 Mobile Companion (Phase 3+)

- Critical alerts on mobile
- Quick actions (acknowledge, escalate)
- Not full dashboard (read-only summaries)

---

> Command Center uses dark theme by default (--bg: #0A0E1A, --text: #EAF0FF).
