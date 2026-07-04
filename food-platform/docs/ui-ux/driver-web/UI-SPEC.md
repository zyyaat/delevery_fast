# Driver Web App — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: Uber Eats Driver App, Talabat Rider App research

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

The Driver Web App enables delivery drivers to:
- Go online/offline to receive orders
- View and accept/reject incoming orders
- Navigate to restaurants and customers
- Confirm pickups and dropoffs
- Track earnings and request payouts
- View delivery history

### 1.2 Design Inspiration

| Source | What We Take |
|--------|--------------|
| Uber Eats Driver App | Clean earnings display, 15s order timer, heat map |
| Talabat Rider App | Arabic-first, Vodafone Cash payout, simple flow |
| Uber Eats Base Design | Bottom sheet for order offers, clear CTAs |

### 1.3 Design Principles (App-Specific)

1. **Glanceable**: Driver looks at screen while riding — must understand in <2s
2. **High Contrast**: Readable in bright sunlight
3. **Large Touch Targets**: 56px+ for gloves/wet hands
4. **Audio-First**: Important alerts use sound (driver may not see screen)
5. **Battery-Aware**: Minimize GPS updates when idle
6. **Offline-Tolerant**: Last order info cached for dead zones

---

## 2. Information Architecture

### 2.1 App Shell (Mobile-Optimized for Web)

```
┌─────────────────────────────────┐
│ Top Bar (56px)                  │
│ 🟢 Online    EGP 245  [⚙️]      │
├─────────────────────────────────┤
│                                 │
│                                 │
│   Content (scrollable)          │
│                                 │
│                                 │
│                                 │
├─────────────────────────────────┤
│ Bottom Nav (64px)               │
│ [🏠 الرئيسية] [💰 الأرباح] [📦 السجل] [👤 حسابي] │
└─────────────────────────────────┘
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
              ┌──────▼───────┐
              │ Profile Setup│
              │ (KYC basic)  │
              └──────┬───────┘
                     │
        ┌────────────▼────────────┐
        │       HOME (Online)     │◄──────┐
        └─┬───────┬───────┬───────┘       │
          │       │       │               │
   ┌──────▼──┐ ┌──▼───┐ ┌─▼────────┐     │
   │ Order   │ │Heat  │ │Profile   │     │
   │ Offer   │ │Map   │ │(Tier)    │     │
   │ Modal   │ └──────┘ └──────────┘     │
   └────┬────┘                            │
        │ accept                          │
   ┌────▼─────┐                           │
   │ Pickup   │                           │
   │ Screen   │                           │
   └────┬─────┘                           │
        │ confirm                         │
   ┌────▼─────┐                           │
   │ Dropoff  │                           │
   │ Screen   │                           │
   └────┬─────┘                           │
        │ confirm                         │
   ┌────▼──────────┐                      │
   │   Order       │                      │
   │  Completed    │──────────────────────┘
   └───────────────┘
```

---

## 3. Screen Specifications

---

### DRV-AUTH-01: Phone Login

**Purpose**: Driver enters phone number.

**Wireframe**:
```
┌─────────────────────────────────┐
│                                 │
│                                 │
│         [App Logo]              │
│                                 │
│         مندوب توصيل              │
│                                 │
│   سجّل رقم موبايلك عشان          │
│   نبعتبلك كود تفعيل              │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🇪🇬 +20                   │   │
│   │ [01 2345 6789        ] │   │
│   └─────────────────────────┘   │
│                                 │
│   [أرسل الكود]                  │
│                                 │
│   ─── أو ───                    │
│                                 │
│   عندك حساب؟ [سجّل دخول]         │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Same as Customer Login but branded "مندوب توصيل"
- Phone validation: Egyptian mobile (01XXXXXXXXX)
- No social login (drivers must be verified individually)

---

### DRV-AUTH-02: OTP Verification

**Purpose**: Verify driver's phone with OTP.

Same as Customer OTP screen (CUS-AUTH-03) with:
- Title: "أدخل الكود"
- 6-digit input boxes
- 2-minute timer
- Resend option

---

### DRV-AUTH-03: Profile Setup (KYC Basic)

**Purpose**: Driver enters basic info for KYC.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  إكمال الملف                  │
├─────────────────────────────────┤
│                                 │
│   خطوة 1 من 3                   │
│   ━━━━━━━━━━━━━━━━ (33%)        │
│                                 │
│   المعلومات الأساسية              │
│                                 │
│   ┌─────────────────────────┐   │
│   │ الاسم الكامل              │   │
│   │ [محمود سعيد            ] │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ الرقم القومي              │   │
│   │ [29XXXXXXXXXX          ] │   │
│   └─────────────────────────┘   │
│                                 │
│   نوع المركبة                    │
│   ● 🛵 موتوسيكل                  │
│   ○ 🚗 عربية                     │
│   ○ 🚲 دراجة                     │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📷 صورة البطاقة (أمامي)   │   │
│   │ [ارفع الصورة]             │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📷 صورة البطاقة (خلفي)    │   │
│   │ [ارفع الصورة]             │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📷 صورة شخصية مع البطاقة │   │
│   │ [ارفع الصورة]             │   │
│   └─────────────────────────┘   │
│                                 │
│   [التالي →]                    │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Progress indicator: "خطوة 1 من 3" + bar
- Step 1: Basic info (name, national ID, vehicle type, ID photos)
- Step 2: Vehicle info (license, registration, insurance)
- Step 3: Bank/wallet info (for payouts)

**Step 2 - Vehicle**:
```
│   ┌─────────────────────────┐   │
│   │ 📷 رخصة القيادة           │   │
│   │ [ارفع الصورة]             │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📷 لوحات المركبة          │   │
│   │ [ارفع الصورة]             │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📷 تأمين المركبة          │   │
│   │ [ارفع الصورة]             │   │
│   └─────────────────────────┘   │
```

**Step 3 - Payout**:
```
│   طريقة استلام الأرباح:           │
│                                 │
│   ● 💚 Vodafone Cash            │
│     رقم المحفظة: [01XXXXXXXXX] │
│                                 │
│   ○ 🟣 InstaPay                │
│                                 │
│   ○ 🏦 تحويل بنكي                │
│     [اسم البنك ▾]               │
│     [رقم الحساب]                │
```

**Validation**:
- All fields required
- National ID: 14 digits, valid Egyptian format
- Photos: max 5MB each, JPG/PNG
- Phone for Vodafone Cash: must match login phone

**States**:
- Uploading: progress bar on photo
- Upload failed: retry button
- Submitted: "جاري المراجعة - هنتواصل معاك في 48 ساعة"

---

### DRV-AUTH-04: Training (Basic)

**Purpose**: Brief training for new drivers.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  تدريب سريع                  │
├─────────────────────────────────┤
│                                 │
│   خطوة 2 من 6                   │
│   ━━━━━━━━━━━━━━━━━ (33%)       │
│                                 │
│   ┌─────────────────────────┐   │
│   │                         │   │
│   │   [Video/GIF 16:9]      │   │
│   │                         │   │
│   └─────────────────────────┘   │
│                                 │
│   استلام الطلب                   │
│                                 │
│   لما يوصلك طلب:                 │
│   1. بص على التفاصيل             │
│   2. اضغط "قبول" خلال 15 ثانية   │
│   3. روح للمطعم                  │
│   4. استلم الطلب                 │
│   5. أكّد الاستلام                │
│                                 │
│   [التالي →]                    │
│                                 │
└─────────────────────────────────┘
```

**Training Modules** (6):
1. استلام الطلب (Order Acceptance)
2. الاستلام من المطعم (Pickup)
3. التوصيل للعميل (Dropoff)
4. التعامل مع العملاء (Customer Service)
5. السلامة المرورية (Road Safety)
6. استلام الأرباح (Earnings)

**Specs**:
- Each module: 30-60s video + bullet points
- Progress bar at top
- Auto-advance optional
- Final quiz: 10 questions, 8+ to pass

---

### DRV-HOME-01: Home (Online) ⭐

**Purpose**: Driver's main screen when online.

**Wireframe**:
```
┌─────────────────────────────────┐
│ 🟢 Online  EGP 245  [⚙️]        │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │ 💰 اليوم                  │   │
│   │ EGP 245.50               │   │
│   │ 7 توصيلات • 4.2 ساعات    │   │
│   │ المعدل: EGP 58.5/ساعة     │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📅 هذا الأسبوع             │   │
│   │ EGP 1,820                │   │
│   │ 52 توصيلة                  │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🔥 مناطق الطلب العالي      │   │
│   │                          │   │
│   │  [Heat Map Mini]         │   │
│   │                          │   │
│   │  🔴 معادي - 3 طلبات نشطة  │   │
│   │  🔴 الزمالك - 5 طلبات    │   │
│   │  🟡 مدينة نصر - 2 طلبات  │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🏆 المستوى                │   │
│   │ 🥇 Platinum               │   │
│   │ ⭐ 4.92 (148 تقييم)        │   │
│   │ قبول: 91% | إكمال: 96%   │   │
│   └─────────────────────────┘   │
│                                 │
│   [🛑 إنهاء العمل]               │
│                                 │
├─────────────────────────────────┤
│ [🏠] [💰] [📦] [👤]              │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: status indicator (🟢 Online), today's earnings, settings
- Today card: earnings, deliveries, hours, rate
- Week card: weekly total + count
- Heat map card: mini map + text list of hot zones
- Tier card: badge, rating, acceptance/completion rates
- CTA: "🛑 إنهاء العمل" (large, prominent)

**Status Indicator**:
- 🟢 Online (green) — receiving orders
- 🟡 On Break (yellow) — paused
- 🔴 Offline (red) — not working
- Tapping toggles state with confirmation

**Interactions**:
- Tap status: toggle online/offline (with confirmation)
- Tap heat map: navigate to full Heat Map screen
- Tap tier card: navigate to Profile
- Tap settings (⚙️): open settings menu

**States**:
- Online + waiting: "بنبحث لك على طلبات..."
- Online + order incoming: Order Offer modal appears
- On active delivery: replaced by Active Delivery screen
- Offline: gray status + "روح Online عشان تستقبل طلبات"

**Audio**:
- Order incoming: distinctive sound (loud, 2s)
- Order accepted: confirmation sound (short, 0.5s)
- Order timeout: warning sound (if auto-declined)

---

### DRV-HOME-02: Heat Map

**Purpose**: Full-screen demand visualization.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  مناطق الطلب العالي            │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │                         │   │
│   │   [Full Screen Map]     │   │
│   │                         │   │
│   │   🔴🔴🔴 معادي            │   │
│   │   🔴🔴🔴🔴                │   │
│   │   🔴🔴🔴🔴🔴 الزمالك      │   │
│   │   🟡🟡🟡🔴🔴🔴 مدينة نصر   │   │
│   │   🟡🟡🟡🟡🟢🟢🟢           │   │
│   │   🟢🟢🟢🟢🟢🟢🟢 التحرير   │   │
│   │                         │   │
│   │   🛵 موقعك               │   │
│   │                         │   │
│   └─────────────────────────┘   │
│                                 │
│   🔴 طلب عالي جداً               │
│   🟡 طلب متوسط                   │
│   🟢 طلب منخفض                   │
│                                 │
│   آخر تحديث: 14:32              │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Full-screen Google Maps
- Heat layer: red/yellow/green based on demand/supply ratio
- Driver location: blue dot (current)
- Auto-refresh: every 2 minutes
- Tap zone: show details (active orders, drivers needed)

**Interactions**:
- Tap zone: popup with details
- Pinch: zoom
- Pan: drag map

---

### DRV-HOME-03: Profile

**Purpose**: Driver profile with tier, rating, stats.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  حسابي                        │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │ [Avatar] Mahmoud S.      │   │
│   │ 🥇 Platinum Member        │   │
│   │ 📱 01012345678           │   │
│   │ 🛵 موتوسيكل - لوحة 1234   │   │
│   │ [تعديل الملف]             │   │
│   └─────────────────────────┘   │
│                                 │
│   📊 إحصائياتي                    │
│   ┌─────────────────────────┐   │
│   │ إجمالي التوصيلات: 247     │   │
│   │ إجمالي الأرباح: EGP 8,420 │   │
│   │ متوسط التقييم: 4.92       │   │
│   │ معدل القبول: 91%          │   │
│   │ معدل الإكمال: 96%         │   │
│   └─────────────────────────┘   │
│                                 │
│   🏆 المستوى                      │
│   ┌─────────────────────────┐   │
│   │ 🥇 Platinum               │   │
│   │ تقييم: 4.92 (148)         │   │
│   │ قبول: 91% | إكمال: 96%   │   │
│   │                          │   │
│   │ المزايا:                  │   │
│   │ • أولوية في الطلبات        │   │
│   │ • +10% bonus              │   │
│   │ • دعم VIP                 │   │
│   └─────────────────────────┘   │
│                                 │
│   ⚙️ الإعدادات                    │
│   • الملاحة (Google Maps ▾)     │   │
│   • الإشعارات                   │   │
│   • اللغة                       │   │
│   • تسجيل الخروج                │   │
│                                 │
└─────────────────────────────────┘
```

---

### DRV-ORDR-01: Order Offer Modal ⭐ (Most Critical)

**Purpose**: Display incoming order, 15s to accept.

**Wireframe**:
```
┌─────────────────────────────────┐
│                                 │
│   ⏱️ 14 ثانية متبقية             │
│   ━━━━━━━━━━━━━━━━━━━ (93%)     │
│                                 │
│   🔔 طلب جديد                    │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📍 McDonald's             │   │
│   │ معادي - 1.2km منك         │   │
│   │                          │   │
│   │ 📍 العميل                  │   │
│   │ الدقي - 4.5km             │   │
│   │                          │   │
│   │ ⏱️ ETA التوصيل: 22 دقيقة   │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 💰 أرباحك من الطلب         │   │
│   │                          │   │
│   │ EGP 28.50                │   │
│   │                          │   │
│   │ • رسوم التوصيل: EGP 20.00 │   │
│   │ • مكافأة المسافة: EGP 5.50│   │
│   │ • Peak bonus: EGP 3.00   │   │
│   │                          │   │
│   │ ⚡ Instant payout متاح     │   │
│   └─────────────────────────┘   │
│                                 │
│   [رفض]      [✅ قبول]           │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Full-screen modal (overlay)
- Timer: 15 seconds, prominently displayed
- Progress bar: visual countdown (red when <5s)
- Order info: restaurant, distance, customer, distance, ETA
- Earnings: total + breakdown (base + distance + peak)
- Instant payout badge: if eligible
- Two CTAs: Reject (left, secondary) + Accept (right, primary, large)

**Timer Behavior**:
- 15s countdown
- Sound at: 0s (incoming), 10s (warning beep), 0s (timeout)
- At 0s: auto-reject (counts against acceptance rate)
- Visual: progress bar shrinks, turns red at 5s

**Sound**:
- New order: distinctive loud sound (2s, repeatable)
- 5s warning: beep
- Accepted: success sound
- Timeout: subtle "missed" sound

**Interactions**:
- Tap accept: navigate to Pickup Screen
- Tap reject: close modal, return to Home
- Tap anywhere else: ignored (prevent accidental dismiss)
- Back button: disabled (must choose accept/reject)

**States**:
- Incoming: modal appears + sound
- Accepted: brief success animation → navigate to Pickup
- Rejected: fade out → return to Home
- Timeout: fade out → return to Home (auto-rejected)
- Network issue: "فقدت الاتصال - حاول تاني"

---

### DRV-ORDR-02: Pickup Screen ⭐

**Purpose**: Navigate to restaurant and confirm pickup.

**Wireframe**:
```
┌─────────────────────────────────┐
│  الطلب #A7X92F  [3 من 5]        │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │ 1 ✅ الذهاب للمطعم        │   │
│   │ 2 🟡 استلام الطلب          │   │
│   │ 3 ⬜ التوصيل للعميل         │   │
│   │ 4 ⬜ تسليم الطلب           │   │
│   │ 5 ⬜ إتمام الطلب           │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📍 McDonald's - معادي    │   │
│   │ 1.2km - 5 دقايق           │   │
│   │                          │   │
│   │ [📍 ابدأ الملاحة]          │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📋 تفاصيل الطلب            │   │
│   │                          │   │
│   │ 1. Big Mac Meal          │   │
│   │    - بدون بصل             │   │
│   │    - كوكا صغير            │   │
│   │                          │   │
│   │ 2. McChicken             │   │
│   │    - كبير                 │   │
│   │                          │   │
│   │ 3. Apple Pie × 2         │   │
│   │                          │   │
│   │ رمز الاستلام: 4892        │   │
│   │ (اعرضه للمطعم)            │   │
│   └─────────────────────────┘   │
│                                 │
│   [✅ استلمت الطلب]              │
│                                 │
│   [⚠️ المطعم مش جاهز]            │
│   [⚠️ عنصر غير متوفر]            │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: order number + step indicator (3 of 5)
- Progress: 5-step timeline (current highlighted)
- Restaurant card: name, address, distance, ETA
- Order items: list with notes
- Pickup code: large, prominent (show to restaurant)
- CTA: "✅ استلمت الطلب" (primary, large)
- Issue buttons: "المطعم مش جاهز", "عنصر غير متوفر"

**Interactions**:
- Tap navigate: open Google Maps with restaurant as destination
- Tap received: confirm pickup → navigate to Dropoff
- Tap not ready: report delay (extends ETA)
- Tap item unavailable: report issue (item removed from order)

**States**:
- Arrived at restaurant: button changes to "أنا في المطعم"
- Confirmed pickup: success animation → Dropoff screen
- Restaurant delay: notification to customer + new ETA

---

### DRV-ORDR-03: Dropoff Screen ⭐

**Purpose**: Navigate to customer and confirm delivery.

**Wireframe**:
```
┌─────────────────────────────────┐
│  الطلب #A7X92F  [4 من 5]        │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │ 1 ✅ الذهاب للمطعم        │   │
│   │ 2 ✅ استلام الطلب          │   │
│   │ 3 ✅ التوصيل للعميل         │   │
│   │ 4 🟡 تسليم الطلب           │   │
│   │ 5 ⬜ إتمام الطلب           │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📍 12 شارع التحرير، الدقي │   │
│   │ 4.5km - 18 دقيقة          │   │
│   │                          │   │
│   │ [📍 ابدأ الملاحة]          │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 👤 العميل                  │   │
│   │ Ahmed Mohamed             │   │
│   │ 📞 01012345678            │   │
│   │                          │   │
│   │ الشقة: 5 - الدور 2         │   │
│   │                          │   │
│   │ ملاحظات:                  │   │
│   │ "اتصل قبل الوصول،          │   │
│   │  الباب لونه أزرق"          │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🔢 رمز التسليم             │   │
│   │                          │   │
│   │      7 2 9 4             │   │
│   │                          │   │
│   │ (اطلبه من العميل)         │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📸 أو التقط صورة للطلب    │   │
│   │ كدليل على التسليم          │   │
│   │ [📷 التقاط]               │   │
│   └─────────────────────────┘   │
│                                 │
│   [✅ تم التسليم]                │
│                                 │
│   [⚠️ العميل مش موجود]           │
│   [📞 اتصال بالعميل]             │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: order number + step (4 of 5)
- Progress: timeline (3 done, 1 current)
- Customer card: name, phone, address, notes
- Delivery OTP: large, prominent (request from customer)
- Photo proof: optional, take photo as backup
- CTA: "✅ تم التسليم" (primary, large)
- Issue buttons: "العميل مش موجود", "اتصال بالعميل"

**Interactions**:
- Tap navigate: open Google Maps with customer as destination
- Tap call: open phone dialer (anonymized)
- Tap delivered: prompt for OTP → enter → confirm
- Tap photo: open camera → capture → upload
- Tap not available: report issue → support contact

**OTP Entry**:
```
┌─────────────────────────────────┐
│   أدخل رمز التسليم               │
│                                 │
│   ┌─┐ ┌─┐ ┌─┐ ┌─┐               │
│   │ │ │ │ │ │ │ │               │
│   └─┘ └─┘ └─┘ └─┘               │
│                                 │
│   [تأكيد]                       │
│   [إلغاء]                       │
│                                 │
└─────────────────────────────────┘
```

**States**:
- Arrived at customer: button changes to "أنا عند العميل"
- OTP correct: success → Order Completed screen
- OTP wrong: "الرمز مش صحيح" + retry
- Customer unavailable: support intervention

---

### DRV-ORDR-04: Order Completed

**Purpose**: Earnings summary after delivery.

**Wireframe**:
```
┌─────────────────────────────────┐
│                                 │
│                                 │
│         ✅                      │
│      (animation)                │
│                                 │
│      تم التوصيل!                 │
│                                 │
│   أرباحك: EGP 28.50             │
│                                 │
│   • رسوم التوصيل: EGP 20.00     │
│   • مكافأة المسافة: EGP 5.50    │
│   • Peak bonus: EGP 3.00        │
│                                 │
│   ⏱️ المدة: 28 دقيقة              │
│   📏 المسافة: 5.7km              │
│                                 │
│   💰 رصيدك: EGP 273.50           │
│   (هيتحول تلقائياً كل يوم)        │
│                                 │
│   [⚡ سحب فوري]                  │
│                                 │
│   [طلب تاني 🛵]                  │
│                                 │
│   [العودة للرئيسية]              │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Full-screen modal (no top bar, no bottom nav)
- Success animation: checkmark + brief confetti
- Title: "تم التوصيل!"
- Earnings: total + breakdown
- Trip stats: duration, distance
- Wallet balance: with note about daily auto-payout
- CTAs: Instant payout (primary), New order (secondary), Home (text)

**Interactions**:
- Tap instant payout: process payout → show confirmation
- Tap new order: navigate to Home (online)
- Tap home: navigate to Home

---

### DRV-EARN-01: Earnings Dashboard

**Purpose**: Detailed earnings view.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  الأرباح                      │
├─────────────────────────────────┤
│                                 │
│   [اليوم] [الأسبوع] [الشهر]     │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 💰 هذا الأسبوع             │   │
│   │                          │   │
│   │ EGP 1,820.50             │   │
│   │ 52 توصيلة                  │   │
│   │ 38 ساعة                   │   │
│   │ المعدل: EGP 47.9/ساعة      │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📊 الرسم البياني           │   │
│   │                          │   │
│   │  [Bar chart - 7 days]    │   │
│   │                          │   │
│   │ السبت:   EGP 320          │   │
│   │ الأحد:   EGP 280          │   │
│   │ الإتنين: EGP 450          │   │
│   │ ...                      │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 💵 رصيدك القابل للسحب      │   │
│   │                          │   │
│   │ EGP 1,820.50             │   │
│   │                          │   │
│   │ [⚡ سحب فوري (Vodafone)]   │   │
│   │ [📅 سحب يومي (بكرة)]       │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📋 آخر التوصيلات           │   │
│   │                          │   │
│   │ 14:30 - #A7X92F          │   │
│   │ EGP 28.50 ✅              │   │
│   │                          │   │
│   │ 13:15 - #B3K45L          │   │
│   │ EGP 22.00 ✅              │   │
│   │                          │   │
│   │ 12:00 - #C8M72N          │   │
│   │ EGP 35.00 ✅              │   │
│   │                          │   │
│   │ [+ 49 توصيلة]             │   │
│   └─────────────────────────┘   │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + "الأرباح"
- Period tabs: Today / Week / Month
- Summary card: total, deliveries, hours, rate
- Chart: bar chart (7 days for week view)
- Wallet card: balance + payout options
- Recent deliveries: list with time, order #, amount

**Interactions**:
- Tap period: switch view
- Tap instant payout: confirm → process → success
- Tap daily payout: schedule for next day
- Tap delivery: view details

---

### DRV-EARN-02: Payout

**Purpose**: Process instant or scheduled payout.

**Wireframe (Instant Payout Modal)**:
```
┌─────────────────────────────────┐
│                                 │
│   ⚡ سحب فوري                    │
│                                 │
│   ┌─────────────────────────┐   │
│   │ رصيدك القابل للسحب         │   │
│   │                          │   │
│   │ EGP 1,820.50             │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ المبلغ                    │   │
│   │ [1820.50              ] │   │
│   └─────────────────────────┘   │
│                                 │
│   طريقة الاستلام:                │
│   ● 💚 Vodafone Cash            │
│   │   01012345678             │   │
│   ○ 🟣 InstaPay                │   │
│                                 │
│   رسوم الخدمة: EGP 2            │
│   هتستلم: EGP 1818.50           │
│                                 │
│   [⚡ تأكيد السحب]               │
│                                 │
│   [إلغاء]                       │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Modal overlay
- Available balance displayed
- Amount input (default = full balance)
- Payout method: Vodafone Cash (default) or InstaPay
- Fee: EGP 2 (transparent)
- Net amount: clearly shown
- CTA: "⚡ تأكيد السحب"

**States**:
- Processing: spinner + "جاري التحويل..."
- Success: "تم! EGP 1818.50 في محفظتك" + checkmark
- Failed: "حصل خطأ، حاول تاني"

---

### DRV-HIST-01: Order History

**Purpose**: Past deliveries list.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  سجل التوصيلات                │
├─────────────────────────────────┤
│                                 │
│   [📅 اليوم ▾]  [فلترة]          │
│                                 │
│   ─── اليوم (7) ───              │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 14:30  #A7X92F          │   │
│   │ McDonald's → Ahmed M.    │   │
│   │ 5.7km • 28 دقيقة          │   │
│   │ EGP 28.50 ✅              │   │
│   │ ⭐⭐⭐⭐⭐ 4.9 (تقييم العميل)│   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 13:15  #B3K45L          │   │
│   │ KFC → Fatma A.           │   │
│   │ 3.2km • 22 دقيقة          │   │
│   │ EGP 22.00 ✅              │   │
│   │ ⭐⭐⭐⭐ 4.5                │   │
│   └─────────────────────────┘   │
│                                 │
│   ─── أمس (12) ───               │
│   ...                           │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + "سجل التوصيلات"
- Date filter: dropdown (Today, Yesterday, This Week, This Month, Custom)
- Filter button: advanced filters
- Grouped by date
- Order card: time, #, route (restaurant → customer), distance, duration, earnings, rating

**Interactions**:
- Tap date filter: open dropdown
- Tap filter: open advanced filters
- Tap order card: view details (modal)

---

## 4. UX Flows

### 4.1 Driver Goes Online Flow

```
1. App Open
   ↓
2. Auto-Login (token valid)
   ↓
3. Home (Offline status)
   ↓
4. Tap "ابدأ العمل"
   ↓
5. Status changes to 🟢 Online
   ↓
6. Home shows earnings (today: 0), heat map
   ↓
7. Waiting for orders...
   ↓
8. Order Offer Modal appears (with sound)
```

### 4.2 Accept Order Flow

```
1. Order Offer Modal appears
   ↓ (15s timer)
2. Driver reviews:
   - Restaurant + distance
   - Customer + distance
   - Earnings
   ↓
3. Tap "✅ قبول"
   ↓
4. Navigate to Pickup Screen
   ↓
5. Tap "📍 ابدأ الملاحة" → opens Google Maps
   ↓
6. Arrive at restaurant
   ↓
7. Show pickup code (4892) to restaurant
   ↓
8. Receive order
   ↓
9. Tap "✅ استلمت الطلب"
   ↓
10. Navigate to Dropoff Screen
```

### 4.3 Complete Delivery Flow

```
1. Dropoff Screen
   ↓
2. Tap "📍 ابدأ الملاحة" → Google Maps to customer
   ↓
3. Arrive at customer
   ↓
4. Tap "📞 اتصال" (if needed)
   ↓
5. Request OTP from customer
   ↓
6. Enter OTP
   ↓ (validates)
7. Tap "✅ تم التسليم"
   ↓
8. Order Completed screen
   ↓
9. Tap "طلب تاني" → return to Home (online)
```

### 4.4 Instant Payout Flow

```
1. Earnings Dashboard
   ↓
2. Tap "⚡ سحب فوري"
   ↓
3. Payout Modal
   - Amount (default = full)
   - Method (Vodafone Cash default)
   ↓
4. Tap "⚡ تأكيد السحب"
   ↓
5. Processing (spinner)
   ↓
6. Success: "تم! EGP 1818.50 في محفظتك"
   ↓
7. Tap "تم" → return to Earnings
```

### 4.5 Order Rejection Flow

```
1. Order Offer Modal
   ↓
2. Tap "رفض"
   ↓
3. Optional: reason prompt
   - "ليه رفضت؟" (optional)
   - [بعيد] [صغير] [مش هيناسبني]
   ↓
4. Modal closes
   ↓
5. Return to Home (waiting for next order)
```

---

## 5. Component Library

### App-Specific Components

#### OrderOfferModal
```typescript
interface OrderOfferModalProps {
  order: OrderOffer
  onAccept: () => void
  onReject: (reason?: string) => void
  onTimeout: () => void
}
```

#### DeliveryProgress
```typescript
interface DeliveryProgressProps {
  currentStep: 1 | 2 | 3 | 4 | 5
}
```

#### EarningsCard
```typescript
interface EarningsCardProps {
  total: Money
  deliveries: number
  hours: number
  rate: Money
}
```

#### HeatMap
```typescript
interface HeatMapProps {
  zones: HotZone[]
  driverLocation: Coordinates
  onZoneClick?: (zone: HotZone) => void
}
```

#### StatusToggle
```typescript
interface StatusToggleProps {
  status: 'online' | 'offline' | 'on_break'
  onToggle: (newStatus) => void
}
```

#### OTPInput
```typescript
interface OTPInputProps {
  length: 4 | 6
  onComplete: (code: string) => void
}
```

#### TierBadge
```typescript
interface TierBadgeProps {
  tier: 'platinum' | 'gold' | 'silver' | 'standard'
  rating: number
}
```

---

## 6. Critical UX Considerations

### 6.1 Driver Safety

- **One-hand operation**: All critical buttons reachable with thumb
- **Voice navigation**: Default to Google Maps voice nav
- **No typing while moving**: Disable text inputs when speed > 10km/h
- **Emergency button**: SOS in top bar (always accessible)

### 6.2 Battery Optimization

- Reduce GPS polling when idle (30s vs 5s)
- Dim map when not navigating
- Disable non-essential animations

### 6.3 Network Tolerance

- Cache last order info
- Queue location updates when offline
- Sync when connection restored

### 6.4 Audio Alerts

- Distinctive sounds for:
  - New order (urgent, loud)
  - Order accepted (success)
  - Order timeout (subtle)
  - Customer message (gentle)
- Respects device volume
- Vibration on important alerts

---

> **Next**: Read `restaurant-web/UI-SPEC.md` for Restaurant Web App specifications.
