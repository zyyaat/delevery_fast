# Restaurant Web App — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: Uber Eats Manager, Talabat Partner Portal research

---

## Table of Contents

1. [App Overview](#1-app-overview)
2. [Information Architecture](#2-information-architecture)
3. [Screen Specifications](#3-screen-specifications)
4. [UX Flows](#4-ux-flows)
5. [Component Library](#5-component-library)

---

## 1. App Overview

### 1.1 Purpose

The Restaurant Web App enables restaurants to:
- Receive and accept/reject incoming orders
- Manage menu items, modifiers, and availability (86'ing)
- View kitchen display system (KDS) for order preparation
- Track sales analytics and performance
- Create and manage promotions
- Manage operating hours and schedule
- Respond to customer reviews

### 1.2 Design Inspiration

| Source | What We Take |
|--------|--------------|
| Uber Eats Manager | Clean dashboard, order cards with timer, KDS |
| Talabat Partner | Arabic-first, dense info, promo management |
| Square KDS | Color-coded timers, clear statuses |
| Toast POS | Analytics dashboard layout |

### 1.3 Design Principles (App-Specific)

1. **Speed Critical**: 90-second order timer — every second counts
2. **Audio-First**: Loud alerts for new orders (kitchen noise)
3. **Glanceable**: Chef sees order status in <1s from across kitchen
4. **Tablet-First**: Optimized for 10-inch tablet (primary device)
5. **Always-On**: Screen never sleeps during operating hours
6. **Touch-Friendly**: Large buttons for wet/gloved hands

---

## 2. Information Architecture

### 2.1 App Shell (Tablet/Desktop)

```
┌─────────────────────────────────────────────────────┐
│ Top Bar (56px)                                       │
│ 🍔 McDonald's   🟢 مفتوح   [🔊] [⚙️] [👤]            │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ Sidebar │                                           │
│ (240px) │           Content                         │
│         │                                           │
│ 📊 لوحة │                                           │
│ 📦 الطلبات│                                         │
│ 🍽️ المنيو │                                         │
│ 📅 المواعيد│                                        │
│ 🎁 العروض │                                         │
│ 📈 التقارير│                                        │
│ ⭐ المراجعات│                                       │
│ ⚙️ الإعدادات│                                       │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

### 2.2 Navigation Map

```
              ┌─────────────┐
              │ Phone Login  │
              └──────┬───────┘
                     │
              ┌──────▼───────┐
              │ OTP Verify   │
              └──────┬───────┘
                     │
        ┌────────────▼────────────┐
        │       DASHBOARD          │◄──────┐
        └─┬───────┬───────┬────────┘       │
          │       │       │               │
   ┌──────▼──┐ ┌──▼───┐ ┌─▼────────┐     │
   │ Orders  │ │Menu  │ │Schedule  │     │
   │ Active  │ │Manage│ │          │     │
   └────┬────┘ └──┬───┘ └──────────┘     │
        │         │                       │
   ┌────▼────┐ ┌──▼────────┐              │
   │ Inbound │ │ Item      │              │
   │ Modal   │ │ Editor    │              │
   │ (90s)   │ └───────────┘              │
   └────┬────┘                            │
        │ accept                          │
   ┌────▼─────┐                           │
   │ KDS      │                           │
   │ (orders) │                           │
   └──────────┘                           │
                                         │
   ┌──────────┐ ┌──────────┐ ┌──────────┐│
   │ Promos   │ │Analytics │ │ Reviews  ││
   └──────────┘ └──────────┘ └──────────┘│
                                         │
   ┌─────────────────────────────────────┘
   │
   └─→ Back to Dashboard
```

---

## 3. Screen Specifications

---

### REST-AUTH-01: Phone Login

**Purpose**: Restaurant staff login.

**Wireframe**:
```
┌─────────────────────────────────┐
│                                 │
│         [App Logo]              │
│                                 │
│         مطعم / تاجر              │
│                                 │
│   سجّل رقم موبايلك                │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🇪🇬 +20                   │   │
│   │ [01 2345 6789        ] │   │
│   └─────────────────────────┘   │
│                                 │
│   [أرسل الكود]                  │
│                                 │
└─────────────────────────────────┘
```

Same as Driver login but branded "مطعم / تاجر".

---

### REST-AUTH-02: OTP Verification

Same as Customer/Driver OTP with 6-digit input + 2-min timer.

---

### REST-DASH-01: Main Dashboard ⭐

**Purpose**: Restaurant's home screen with today's stats and active orders.

**Wireframe (Tablet/Desktop)**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's - معادي   🟢 مفتوح   [🔊] [⚙️] [👤]    │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊 لوحة │  ┌────────┐ ┌────────┐ ┌────────┐ ┌────┐ │
│ 📦 الطلبات│  │ 📊 اليوم│ │ ⏱️ متوسط│ │ ⭐ تقييم│ │ 📦 │ │
│ 🍽️ المنيو │  │ EGP    │ │ التحضير │ │ 4.6    │ │ 28 │ │
│ 📅 المواعيد│  │ 4,250  │ │ 12 د    │ │ (124)  │ │ طلب│ │
│ 🎁 العروض │  │ ↑12%   │ │         │ │        │ │ نشط│ │
│ 📈 التقارير│  └────────┘ └────────┘ └────────┘ └────┘ │
│ ⭐ المراجعات│                                           │
│ ⚙️ الإعدادات│  🔴 طلبات نشطة (3)                          │
│         │                                           │
│         │  ┌──────────────────────────────────────┐ │
│         │  │ ⏱️ #A7X92F  ⏰ 8:14 دقيقة             │ │
│         │  │                                       │ │
│         │  │ 🟡 يتحضّر                              │ │
│         │  │                                       │ │
│         │  │ • Big Mac Meal × 1                    │ │
│         │  │ • McChicken × 2                       │ │
│         │  │ • Apple Pie × 1                       │ │
│         │  │                                       │ │
│         │  │ EGP 290 | Vodafone Cash               │ │
│         │  │                                       │ │
│         │  │ [✅ جاهز للاستلام] [⚠️ تأخير]          │ │
│         │  └──────────────────────────────────────┘ │
│         │                                           │
│         │  ┌──────────────────────────────────────┐ │
│         │  │ ⏱️ #B3K45L  ⏰ 2:30 دقيقة             │ │
│         │  │ 🔴 جديد                                 │ │
│         │  │ • Pizza Margherita × 1                │ │
│         │  │ • Coca Cola × 2                       │ │
│         │  │ EGP 145 | Card                        │ │
│         │  │ [▶️ ابدأ التحضير]                      │ │
│         │  └──────────────────────────────────────┘ │
│         │                                           │
│         │  ┌──────────────────────────────────────┐ │
│         │  │ ⏱️ #C8M72N  ⏰ 5:45 دقيقة             │ │
│         │  │ 🟢 جاهز                                 │ │
│         │  │ • Sushi Combo × 2                     │ │
│         │  │ • Miso Soup × 2                       │ │
│         │  │ EGP 320 | COD                         │ │
│         │  │ 🛵 المندوب: Mahmoud S. - ETA: 3 دقايق   │ │
│         │  └──────────────────────────────────────┘ │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

**Specs**:
- Top bar: restaurant name + branch, open/closed toggle, sound toggle, settings, profile
- Sidebar: 8 menu items (always visible on desktop, drawer on mobile)
- Stats cards: 4 cards (today's revenue, avg prep time, rating, active orders)
- Active orders: list of order cards sorted by urgency

**Order Card Component**:
```
┌──────────────────────────────────────┐
│ ⏱️ #A7X92F  ⏰ 8:14 دقيقة              │
│                                       │
│ 🟡 يتحضّر                              │
│                                       │
│ • Big Mac Meal × 1                    │
│ • McChicken × 2                       │
│ • Apple Pie × 1                       │
│                                       │
│ EGP 290 | Vodafone Cash               │
│                                       │
│ [✅ جاهز للاستلام] [⚠️ تأخير]          │
└──────────────────────────────────────┘
```

**Timer Color Coding** (Square KDS inspired):
- 🟢 Green: within prep time (e.g., <8 min for 10-min item)
- 🟡 Yellow: approaching limit (8-10 min)
- 🔴 Red: over limit (>10 min, pulsing)

**Status Icons**:
- 🔴 جديد (new, just received)
- 🟡 يتحضّر (preparing)
- 🟢 جاهز (ready, waiting for driver)
- 🔵 استلمه المندوب (picked up)

**Interactions**:
- Tap order card: expand for full details
- Tap "ابدأ التحضير": change status to preparing
- Tap "جاهز للاستلام": change status to ready → notify driver
- Tap "تأخير": open delay dialog (5/10/15 min)
- Tap stats card: navigate to analytics

**Audio**:
- New order: loud notification sound (kitchen-friendly)
- Cannot disable during operating hours
- Volume control in settings

**States**:
- No active orders: empty state "مفيش طلبات نشطة دلوقتي 🎉"
- Restaurant closed: banner "المطعم مقفل. الطلبات مش هتوصلك."

---

### REST-DASH-02: Schedule / Hours

**Purpose**: Manage operating hours.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  ساعات العمل                                │
│ 📦      │                                           │
│ 🍽️      │  ┌─────────────────────────────────────┐  │
│ 📅 ●    │  │ السبت                                │  │
│ 🎁      │  │ ● مفتوح                              │  │
│ 📈      │  │ من [10:00 ص] إلى [02:00 ص            │  │
│ ⭐      │  └─────────────────────────────────────┘  │
│ ⚙️      │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ الأحد                                │  │
│         │  │ ● مفتوح                              │  │
│         │  │ من [10:00 ص] إلى [02:00 ص            │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ الإتنين                              │  │
│         │  │ ○ مقفل                              │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ... (7 days)                             │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 📅 عطلة خاصة                          │  │
│         │  │ [اختر التاريخ] - [سبب العطلة]          │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  [💾 حفظ التغييرات]                        │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

**Specs**:
- 7-day schedule with open/closed toggle per day
- Time pickers for open hours
- Holiday mode: special date + reason
- Save button at bottom

**Interactions**:
- Toggle open/closed per day
- Tap time: open time picker
- Add holiday: date picker + reason
- Save: persist changes

---

### REST-ORDR-01: Active Orders

**Purpose**: List of all current orders.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  الطلبات النشطة (3)                          │
│ 📦 ●    │                                           │
│ 🍽️      │  [الكل] [جديد] [يتحضّر] [جاهز]              │
│ 📅      │                                           │
│ 🎁      │  ┌─────────────────────────────────────┐  │
│ 📈      │  │ 🔴 #A7X92F  ⏰ 8:14  جديدة            │  │
│ ⭐      │  │ • Big Mac Meal × 1                   │  │
│ ⚙️      │  │ • McChicken × 2                      │  │
│         │  │ EGP 290 | VF Cash                    │  │
│         │  │ [▶️ ابدأ التحضير] [❌ رفض]             │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🟡 #B3K45L  ⏰ 12:30  يتحضّر         │  │
│         │  │ • Pizza × 1                          │  │
│         │  │ • Cola × 2                           │  │
│         │  │ EGP 145 | Card                       │  │
│         │  │ [✅ جاهز] [⚠️ تأخير]                  │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🟢 #C8M72N  ⏰ 5:45  جاهز            │  │
│         │  │ • Sushi × 2                          │  │
│         │  │ • Miso Soup × 2                      │  │
│         │  │ EGP 320 | COD                        │  │
│         │  │ 🛵 Mahmoud S. - 3 دقايق               │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

**Specs**:
- Filter tabs: All / New / Preparing / Ready
- Each order card shows: ID, timer, status, items, total, payment, actions
- Sort: by time (oldest first)

---

### REST-ORDR-02: Inbound Order Modal ⭐ (90s Timer)

**Purpose**: New order arrives, 90 seconds to accept/reject.

**Wireframe (Full-screen modal)**:
```
┌─────────────────────────────────────────────────────┐
│                                                     │
│   ⏱️ 75 ثانية متبقية                                 │
│   ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ (83%)    │
│                                                     │
│   🔔 طلب جديد #A7X92F                                │
│                                                     │
│   ┌───────────────────────────────────────────────┐ │
│   │ 🛒 تفاصيل الطلب                                  │ │
│   │                                                 │ │
│   │ 1× Big Mac Meal                                 │ │
│   │   - بدون بصل                                     │ │
│   │   - كوكا صغير                                    │ │
│   │   EGP 85                                        │ │
│   │                                                 │ │
│   │ 2× McChicken                                    │ │
│   │   - كبير                                         │ │
│   │   EGP 90 × 2 = 180                              │ │
│   │                                                 │ │
│   │ 1× Apple Pie                                    │ │
│   │   EGP 25                                        │ │
│   │                                                 │ │
│   │ ─────────────────────────────────              │ │
│   │ الإجمالي: EGP 290                                │ │
│   │                                                 │ │
│   │ الدفع: Vodafone Cash                            │ │
│   │ التوصيل: 4.5km - 22د                            │ │
│   │ العميل: Ahmed M.                                 │ │
│   │                                                 │ │
│   │ 📝 ملاحظات:                                      │ │
│   │ "صلصة إضافية لو سمحت"                            │ │
│   └───────────────────────────────────────────────┘ │
│                                                     │
│   ⏱️ وقت التحضير المتوقع:                              │
│   [-]  15 دقيقة  [+]                                 │
│                                                     │
│   [❌ رفض]              [✅ قبول]                     │
│                                                     │
│   [⚠️ تعديل وقت التحضير]                              │
│                                                     │
└─────────────────────────────────────────────────────┘
```

**Specs**:
- Full-screen modal (can't be dismissed without action)
- Timer: 90 seconds, prominently displayed at top
- Progress bar: visual countdown (red at <30s)
- Order details: items, modifiers, totals, payment, customer, notes
- Prep time selector: - / + buttons (default = estimated)
- Two CTAs: Reject (left, secondary) + Accept (right, primary, large)

**Audio**:
- New order: loud alarm sound (kitchen-penetrating)
- 30s warning: beep
- 10s warning: urgent beep
- Timeout: "missed" sound
- Accepted: success chime

**Timer Behavior**:
- 90s countdown
- Color: green (>60s), yellow (30-60s), red (<30s, pulsing)
- At 0s: auto-reject (counts against restaurant's acceptance rate)
- Persistent: can't close without action

**Interactions**:
- Tap accept: confirm → order moves to Active Orders (preparing)
- Tap reject: prompt for reason → order cancelled → customer notified
- Tap +/- prep time: adjust (max 60 min)
- Tap "تعديل وقت التحضير": custom time input

**Reject Reasons**:
- عنصر غير متوفر (item unavailable)
- المطعم مشغول جداً (too busy)
- خارج منطقة التوصيل (out of area)
- مشكلة في الطلب (order issue)

**States**:
- Incoming: modal + sound
- Accepted: success → navigate to Active Orders
- Rejected: customer notified → modal closes
- Timeout: auto-reject → notification to restaurant

---

### REST-ORDR-03: Order Detail

**Purpose**: Full details of a specific order.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  ← #A7X92F                                 │
│ 📦 ●    │                                           │
│ 🍽️      │  ┌─────────────────────────────────────┐  │
│ 📅      │  │ الحالة: 🟡 يتحضّر                      │  │
│ 🎁      │  │ الوقت: 8:14 دقيقة                      │  │
│ 📈      │  │ ETA التسليم: 8:35 PM                  │  │
│ ⭐      │  └─────────────────────────────────────┘  │
│ ⚙️      │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🛒 العناصر                              │  │
│         │  │                                       │  │
│         │  │ 1× Big Mac Meal        EGP 85         │  │
│         │  │   - بدون بصل                           │  │
│         │  │   - كوكا صغير                          │  │
│         │  │                                       │  │
│         │  │ 2× McChicken           EGP 90 × 2     │  │
│         │  │   - كبير                               │  │
│         │  │                                       │  │
│         │  │ 1× Apple Pie           EGP 25         │  │
│         │  │                                       │  │
│         │  │ ─────────────────────────              │  │
│         │  │ الإجمالي: EGP 290                       │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 👤 العميل                              │  │
│         │  │ Ahmed Mohamed                         │  │
│         │  │ 📞 01012345678 (مخفي)                  │  │
│         │  │ 📍 الزمالك، 26 يوليو، شقة 5            │  │
│         │  │ 📝 "صلصة إضافية لو سمحت"               │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 💳 الدفع                               │  │
│         │  │ Vodafone Cash - مدفوع                 │  │
│         │  │ EGP 290                               │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🛵 المندوب                              │  │
│         │  │ Mahmoud S. - ⭐ 4.8                    │  │
│         │  │ ETA الاستلام: 3 دقايق                  │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  [✅ جاهز للاستلام]                         │
│         │  [⚠️ تأخير 5 دقايق]                        │
│         │  [📞 اتصال بالعميل]                         │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

---

### REST-ORDR-04: Order History

**Purpose**: Past orders with filters.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  سجل الطلبات                                 │
│ 📦 ●    │                                           │
│ 🍽️      │  [📅 اليوم ▾]  [فلترة]  [📊 تصدير]         │
│ 📅      │                                           │
│ 🎁      │  ─── اليوم (28 طلب) ───                     │
│ 📈      │                                           │
│ ⭐      │  ┌─────────────────────────────────────┐  │
│ ⚙️      │  │ 14:32  #A7X92F  ✅ تم التوصيل         │  │
│         │  │ EGP 290 | VF Cash | 3 أصناف          │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 13:15  #B3K45L  ✅ تم التوصيل         │  │
│         │  │ EGP 145 | Card | 2 أصناف             │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 12:00  #C8M72N  ❌ ملغي               │  │
│         │  │ EGP 320 | COD | سبب: عميل رفض         │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ─── أمس (45 طلب) ───                       │
│         │  ...                                     │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

**Specs**:
- Date filter dropdown
- Filter button (status, payment, amount)
- Export button (CSV)
- Grouped by date
- Order cards: time, #, status, total, payment, items count

---

### REST-MENU-01: Menu Overview

**Purpose**: Manage menu categories and items.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  🍽️ المنيو             [➕ إضافة صنف]      │
│ 📦      │                                           │
│ 🍽️ ●    │  ┌─────────────────────────────────────┐  │
│ 📅      │  │ 🍔 البرجر                 12 صنف      │  │
│ 🎁      │  │                                       │  │
│ 📈      │  │ ┌─────┐ ┌─────┐ ┌─────┐              │  │
│ ⭐      │  │ │[img]│ │[img]│ │[img]│              │  │
│ ⚙️      │  │ │Big  │ │McCh │ │QrtP │              │  │
│         │  │ │Mac  │ │icken│ │ounder│             │  │
│         │  │ │85   │ │90   │ │110  │              │  │
│         │  │ │🟢   │ │🟢   │ │🔴 86│              │  │
│         │  │ └─────┘ └─────┘ └─────┘              │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🍕 البيتزا                  8 أصناف   │  │
│         │  │ ┌─────┐ ┌─────┐ ┌─────┐              │  │
│         │  │ │[img]│ │[img]│ │[img]│              │  │
│         │  │ │Marg │ │Pepp │ │Veggie│             │  │
│         │  │ │145  │ │165  │ │155  │              │  │
│         │  │ │🟢   │ │🟡 3 │ │🟢   │              │  │
│         │  │ └─────┘ └─────┘ └─────┘              │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🥤 المشروبات                6 أصناف   │  │
│         │  │ ...                                   │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

**Specs**:
- Categories as collapsible sections
- Items as cards in grid (3-4 per row)
- Each card: image, name, price, availability indicator
- Add item button at top

**Item Card Component**:
```
┌─────────┐
│ [image] │
│         │
│ Big Mac │
│ EGP 85  │
│ 🟢 متاح │
└─────────┘
```
- Image: 1:1, lazy-loaded
- Name: h4
- Price: mono, bold
- Availability: 🟢 available, 🔴 86'd, 🟡 low stock
- Tap: open Item Editor
- Long press: quick toggle availability

---

### REST-MENU-02: Item Editor

**Purpose**: Add/edit menu item.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  ← تعديل: Big Mac                          │
│ 📦      │                                           │
│ 🍽️ ●    │  ┌─────────────────────────────────────┐  │
│ 📅      │  │ 📷 الصور                                │  │
│ 🎁      │  │ [📷 صورة حالية]                          │  │
│ 📈      │  │ [📤 رفع صورة جديدة]                      │  │
│ ⭐      │  └─────────────────────────────────────┘  │
│ ⚙️      │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ الاسم: [Big Mac                  ]    │  │
│         │  │ الوصف: [برجر لحم بقري طازج...    ]   │  │
│         │  │ السعر: [EGP 85                   ]    │  │
│         │  │ التصنيف: [برجر ▾]                       │  │
│         │  │ وقت التحضير: [8 دقايق]                   │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ الإضافات (Modifiers)                    │  │
│         │  │                                       │  │
│         │  │ ☑ بدون بصل (+0)                       │  │
│         │  │ ☑ بدون مخلل (+0)                      │  │
│         │  │ ☑ جبن إضافي (+15)                     │  │
│         │  │ ☑ بيكون (+20)                         │  │
│         │  │ ☐ صلصة حارة (+5)                       │  │
│         │  │ [➕ إضافة modifier]                     │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ التوفر                                  │  │
│         │  │ ● متاح دائماً                           │  │
│         │  │ ○ متاح في أوقات معينة                   │  │
│         │  │ ○ نفد مؤقتاً (86)                       │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  [💾 حفظ]  [❌ إلغاء]                       │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

---

### REST-MENU-03: Category Editor

**Purpose**: Manage menu categories.

Similar to Item Editor but for categories (name, order, items).

---

### REST-MENU-04: Bulk Availability (86'ing)

**Purpose**: Quick toggle item availability.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  التوفر السريع                              │
│ 📦      │                                           │
│ 🍽️ ●    │  ┌─────────────────────────────────────┐  │
│ 📅      │  │ 🍔 البرجر                                │  │
│ 🎁      │  │                                       │  │
│ 📈      │  │ ☑ Big Mac         🟢 متاح              │  │
│ ⭐      │  │ ☑ McChicken       🟢 متاح              │  │
│ ⚙️      │  │ ☐ Quarter Pounder 🔴 نفد (86)         │  │
│         │  │ ☑ Filet-O-Fish    🟢 متاح              │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🍕 البيتزا                               │  │
│         │  │ ☑ Margherita      🟢 متاح              │  │
│         │  │ ☑ Pepperoni       🟡 آخر 3 قطع        │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  [إخفاء قسم كامل ▾]                        │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

**Specs**:
- Toggle availability per item (instant)
- "إخفاء قسم كامل" - hide entire category
- Changes propagate to customer app in <60s

---

### REST-KDS-01: KDS (Kitchen Display System)

**Purpose**: Full-screen order board for kitchen.

**Wireframe (Full-screen, no sidebar)**:
```
┌─────────────────────────────────────────────────────────────┐
│  المطبخ - 10:45 AM                          [إعدادات] [خروج] │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐         │
│  │ #A7X92F  │ │ #B3K45L  │ │ #C8M72N  │ │ #D4N91P  │         │
│  │ ⏰ 8:14   │ │ ⏰ 2:30   │ │ ⏰ 5:45   │ │ ⏰ 0:45   │         │
│  │ 🟡       │ │ 🔴       │ │ 🟢       │ │ 🔴       │         │
│  │          │ │          │ │          │ │          │         │
│  │ 1× Big   │ │ 1× Pizza │ │ 2× Sushi │ │ 3× Burger│         │
│  │   Mac    │ │   Margh  │ │   Combo  │ │   Meal   │         │
│  │   - بصل   │ │          │ │          │ │          │         │
│  │          │ │ 2× Coca  │ │ 2× Miso  │ │ 1× Fries │         │
│  │ 2× McCh  │ │   Cola   │ │   Soup   │ │          │         │
│  │   - كبير  │ │          │ │          │ │ 2× Cola  │         │
│  │          │ │          │ │          │ │          │         │
│  │ 1× Apple │ │          │ │          │ │          │         │
│  │   Pie    │ │          │ │          │ │          │         │
│  │          │ │          │ │          │ │          │         │
│  │ [▶️ بدأ]  │ │ [▶️ بدأ]  │ │ [✅ خلص] │ │ [▶️ بدأ]  │         │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘         │
│                                                             │
│  ┌──────────┐ ┌──────────┐                                   │
│  │ #E5P22Q  │ │ #F6Q33R  │  ...                              │
│  │ ⏰ 1:20   │ │ ⏰ 0:15   │                                   │
│  └──────────┘ └──────────┘                                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Specs**:
- Full-screen mode (browser fullscreen API)
- No sidebar, no top bar (just minimal header)
- Order cards in grid (3-4 per row, scrollable)
- Each card: order #, timer (color-coded), items, action button

**Card Color Coding** (Square KDS inspired):
- 🔴 Red border: NEW (just received, not started)
- 🟡 Yellow border: PREPARING (started)
- 🟢 Green border: READY (done, waiting for driver)

**Timer Color**:
- 🟢 Green: <50% of prep time
- 🟡 Yellow: 50-90% of prep time
- 🔴 Red: >90% of prep time (pulsing)

**Card States**:
```
New (red border):
┌──────────┐
│ #A7X92F  │
│ ⏰ 0:30   │
│ 🔴 جديد   │
│ ...      │
│ [▶️ بدأ]  │
└──────────┘

Preparing (yellow border):
┌──────────┐
│ #B3K45L  │
│ ⏰ 8:14   │
│ 🟡 يتحضّر │
│ ...      │
│ [✅ خلص]  │
└──────────┘

Ready (green border):
┌──────────┐
│ #C8M72N  │
│ ⏰ 5:45   │
│ 🟢 جاهز   │
│ ...      │
│ 🛵 في الطريق│
└──────────┘
```

**Interactions**:
- Tap "بدأ": move to Preparing (yellow)
- Tap "خلص": move to Ready (green) → notify driver
- Tap card: expand for full details
- Auto-scroll: new orders appear on left

**Settings**:
- Number of columns (2/3/4)
- Sound on new order (toggle)
- Auto-bump after pickup (toggle)
- Font size (small/medium/large)

**Display Modes**:
- Default: All orders
- By Station: filter by prep station (grill, fryer, etc.)
- By Priority: sort by urgency

---

### REST-ANAL-01: Sales Analytics

**Purpose**: Charts and metrics.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  📈 التقارير                                 │
│ 📦      │                                           │
│ 🍽️      │  [📅 اليوم ▾]                              │
│ 📅      │                                           │
│ 🎁      │  ┌─────────────────────────────────────┐  │
│ 📈 ●    │  │ مبيعات اليوم                           │  │
│ ⭐      │  │                                       │  │
│ ⚙️      │  │ EGP 4,250                            │  │
│         │  │ 28 طلب                                 │  │
│         │  │ ↑ 12% عن الأمس                         │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 📊 مبيعات الأسبوع                       │  │
│         │  │                                       │  │
│         │  │  [Bar chart - 7 days]                 │  │
│         │  │                                       │  │
│         │  │ السبت:   EGP 3,200                     │  │
│         │  │ الأحد:   EGP 2,800                     │  │
│         │  │ الإتنين: EGP 4,500                     │  │
│         │  │ ...                                   │  │
│         │  │ الإجمالي: EGP 28,400                    │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🍔 الأصناف الأكثر مبيعاً                  │  │
│         │  │                                       │  │
│         │  │ 1. Big Mac      × 42                  │  │
│         │  │ 2. McChicken    × 38                  │  │
│         │  │ 3. Apple Pie    × 35                  │  │
│         │  │ 4. Coca Cola    × 28                  │  │
│         │  │ 5. Quarter Pndr × 22                  │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ ⏰ أوقات الذروة                          │  │
│         │  │                                       │  │
│         │  │  [Heat map by hour]                   │  │
│         │  │                                       │  │
│         │  │ 12pm-3pm: 45% من الطلبات              │  │
│         │  │ 7pm-11pm: 35% من الطلبات              │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

---

### REST-ANAL-02: Reviews

**Purpose**: Customer reviews with reply capability.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  ⭐ المراجعات                                │
│ 📦      │                                           │
│ 🍽️      │  ┌─────────────────────────────────────┐  │
│ 📅      │  │ ⭐ 4.6 (124 تقييم)                      │  │
│ 🎁      │  │                                       │  │
│ 📈      │  │ ⭐⭐⭐⭐⭐ 78 (63%)                       │  │
│ ⭐ ●    │  │ ⭐⭐⭐⭐  28 (23%)                       │  │
│ ⚙️      │  │ ⭐⭐⭐    12 (10%)                       │  │
│         │  │ ⭐⭐      4 (3%)                        │  │
│         │  │ ⭐       2 (1%)                        │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ─── آخر المراجعات ───                       │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ ⭐⭐⭐⭐⭐  Ahmed M.  -  أمس              │  │
│         │  │ #A7X92F - EGP 290                     │  │
│         │  │ "الأكل لذيذ والتوصيل سريع"              │  │
│         │  │                                       │  │
│         │  │ [💬 رد]                                 │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ ⭐⭐⭐  Fatma A.  -  قبل 2 يوم           │  │
│         │  │ #B3K45L - EGP 145                     │  │
│         │  │ "التوصيل اتأخر كتير"                   │  │
│         │  │                                       │  │
│         │  │ 💬 رد المطعم: "آسفين للتأخير..."       │  │
│         │  │ [تعديل الرد]                            │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

---

### REST-ANAL-03: Peak Hours

**Purpose**: Heatmap of orders by hour.

Shows hour-by-hour heatmap + day-of-week heatmap.

---

### REST-PROMO-01: Promotions List

**Purpose**: Manage active and scheduled promotions.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  🎁 العروض              [➕ حملة جديدة]     │
│ 📦      │                                           │
│ 🍽️      │  [نشطة (2)] [مجدولة (1)] [منتهية (5)]      │
│ 📅      │                                           │
│ 🎁 ●    │  ┌─────────────────────────────────────┐  │
│ 📈      │  │ 🔥 "خصم 20% على البرجر"                 │  │
│ ⭐      │  │ 1-15 يوليو                            │  │
│ ⚙️      │  │ 47 طلب - EGP 2,300                    │  │
│         │  │ [إيقاف] [تعديل] [تفاصيل]                │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 🎁 "وجبة عائلية"                       │  │
│         │  │ دائم                                  │  │
│         │  │ 23 طلب - EGP 4,100                    │  │
│         │  │ [إيقاف] [تعديل] [تفاصيل]                │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

---

### REST-PROMO-02: Create Promotion

**Purpose**: Create new promotion.

**Wireframe**:
```
┌─────────────────────────────────────────────────────┐
│ 🍔 McDonald's   🟢 مفتوح            [⚙️] [👤]        │
├─────────┬───────────────────────────────────────────┤
│         │                                           │
│ 📊      │  ➕ حملة جديدة                               │
│ 📦      │                                           │
│ 🍽️      │  النوع:                                     │
│ 📅      │  ● خصم على صنف                              │
│ 🎁 ●    │  ○ خصم على الطلب كله                        │
│ 📈      │  ○ وجبة组合 (combo)                          │
│ ⭐      │  ○ اشتري 1 واحصل على 1                       │
│ ⚙️      │                                           │
│         │  الصنف: [Big Mac ▾]                         │
│         │                                           │
│         │  القيمة: [20 %]                              │
│         │                                           │
│         │  الفترة:                                     │
│         │  من [1 يوليو] إلى [15 يوليو]                  │
│         │                                           │
│         │  الحد الأقصى للاستخدام: [100]                  │
│         │  الحد الأقصى لكل عميل: [1]                    │
│         │                                           │
│         │  الحد الأدنى للطلب: [EGP 200]                 │
│         │                                           │
│         │  ┌─────────────────────────────────────┐  │
│         │  │ 💸 تكلفة الحملة: EGP 500              │  │
│         │  │ (المنصة تدفع 50% = EGP 250)            │  │
│         │  │ (أنت تدفع 50% = EGP 250)               │  │
│         │  └─────────────────────────────────────┘  │
│         │                                           │
│         │  [🚀 إطلاق الحملة]                          │
│         │                                           │
└─────────┴───────────────────────────────────────────┘
```

---

## 4. UX Flows

### 4.1 Receive Order Flow

```
1. New order arrives
   ↓ (loud sound)
2. Inbound Order Modal appears (90s timer)
   ↓
3. Restaurant reviews:
   - Items + modifiers
   - Customer info + notes
   - Payment method
   ↓
4. Adjust prep time if needed
   ↓
5. Tap "✅ قبول"
   ↓
6. Order moves to Active Orders (preparing status)
   ↓
7. Restaurant prepares food
   ↓
8. Tap "✅ جاهز للاستلام"
   ↓
9. Driver notified
   ↓
10. Driver arrives, picks up
    ↓
11. Order status: picked up
    ↓
12. Order moves to history after delivery
```

### 4.2 Reject Order Flow

```
1. Inbound Order Modal
   ↓
2. Tap "❌ رفض"
   ↓
3. Reason prompt:
   - عنصر غير متوفر
   - المطعم مشغول
   - خارج منطقة التوصيل
   - مشكلة في الطلب
   ↓
4. Confirm rejection
   ↓
5. Customer notified (auto)
   ↓
6. Order cancelled
   ↓
7. Modal closes
```

### 4.3 86'ing Item Flow

```
1. Item runs out in kitchen
   ↓
2. Restaurant opens Menu > Bulk Availability
   ↓
3. Toggle item off (instant)
   ↓
4. Customer app updates in <60s
   ↓
5. Item shows "نفد" in customer app
   ↓
6. No new orders with that item
```

### 4.4 Create Promotion Flow

```
1. Promotions > Create New
   ↓
2. Select type (4 types)
   ↓
3. Configure:
   - Item/scope
   - Discount value
   - Date range
   - Usage limits
   ↓
4. See cost split (50/50 with platform)
   ↓
5. Launch
   ↓
6. Promotion active on customer app
```

---

## 5. Component Library

### App-Specific Components

#### InboundOrderModal
```typescript
interface InboundOrderModalProps {
  order: RestaurantOrder
  onAccept: (prepTime: number) => void
  onReject: (reason: string) => void
  onTimeout: () => void
}
```

#### OrderCard
```typescript
interface OrderCardProps {
  order: RestaurantOrder
  onStatusChange: (status: OrderStatus) => void
  onDelay: (minutes: number) => void
}
```

#### KDSCard
```typescript
interface KDSCardProps {
  order: Order
  status: 'new' | 'preparing' | 'ready'
  onStatusChange: (newStatus) => void
}
```

#### Menu Item Card
```typescript
interface MenuItemCardProps {
  item: MenuItem
  onEdit: () => void
  onToggleAvailability: () => void
}
```

#### Timer
```typescript
interface TimerProps {
  startTime: Date
  duration: number  // seconds
  variant: 'order' | 'kds'
}
```

---

## 6. Critical UX Considerations

### 6.1 Kitchen Environment

- **Loud Audio**: Must cut through kitchen noise
- **High Contrast**: Readable in bright/dim light
- **Large Touch Targets**: Wet/gloved hands
- **Always-On Screen**: Prevent sleep during hours

### 6.2 Speed Critical

- 90s order timer: every second counts
- Tap targets: minimum 56px
- No confirmation dialogs for common actions
- Undo for last action (5s window)

### 6.3 Tablet Optimization

- Primary device: 10-inch tablet (Samsung Galaxy Tab)
- Landscape orientation default
- Touch-optimized (no hover states)
- Split view on larger screens

### 6.4 Multi-User

- Multiple staff can be logged in
- Different roles (manager, cashier, chef)
- Chef sees KDS only
- Manager sees everything

---

> **Next**: Read `support-web/UI-SPEC.md` for Support Web App specifications.
