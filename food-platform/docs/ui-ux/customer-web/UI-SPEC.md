# Customer Web App — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: Uber Eats, Talabat, elmenus, DoorDash research

---

## Table of Contents

1. [App Overview](#1-app-overview)
2. [Information Architecture](#2-information-architecture)
3. [Screen Specifications](#3-screen-specifications)
4. [UX Flows](#4-ux-flows)
5. [Component Library (App-Specific)](#5-component-library)
6. [Arabic RTL Notes](#6-arabic-rtl-notes)
7. [States Reference](#7-states-reference)

---

## 1. App Overview

### 1.1 Purpose

The Customer Web App is the primary ordering interface for end customers. It enables users to:
- Discover restaurants near them
- Browse menus with rich visuals
- Customize and order food
- Pay via local methods (Vodafone Cash, InstaPay, Card, COD)
- Track orders in real-time
- Manage profile, addresses, payment methods
- Earn and redeem loyalty rewards

### 1.2 Design Inspiration

| Inspiration | What We Take |
|-------------|--------------|
| **Uber Eats** | Clean image-forward cards, skeleton loading, 5-step tracking bar, minimal text philosophy, bottom sheet for cart |
| **Talabat** | Arabic-first RTL, Cairo font, dense info in cards, Vodafone Cash default, local restaurant focus |
| **elmenus** | Egyptian cultural context, restaurant discovery emphasis, photo quality |
| **DoorDash** | Bold CTAs, item-level rating, clear pricing breakdown |

### 1.3 Design Principles (App-Specific)

1. **Egyptian-First**: Arabic primary, RTL default, Vodafone Cash front and center
2. **Image-Forward**: Large food photography (16:9 covers, 1:1 items)
3. **Glanceable**: User decides in <3s whether to engage
4. **Trust Through Transparency**: Full pricing breakdown, real ETAs, restaurant health scores
5. **Forgiving**: Easy cancel/refund within policy, reversible actions
6. **Fast**: Every screen loads <2s, lazy-load images, optimistic UI

---

## 2. Information Architecture

### 2.1 App Shell (Mobile)

```
┌─────────────────────────────────┐
│ Top Bar (56px)                  │
│ 📍 Address    🔍 Search    👤   │
├─────────────────────────────────┤
│                                 │
│                                 │
│   Content (scrollable)          │
│                                 │
│                                 │
│                                 │
├─────────────────────────────────┤
│ Bottom Nav (64px)               │
│ [🏠 الرئيسية] [🔍 بحث] [📦 طلباتي] [👤 حسابي] │
└─────────────────────────────────┘
```

### 2.2 App Shell (Desktop ≥1024px)

```
┌────┬────────────────────────────┐
│ S  │ Top Bar (56px)             │
│ i  │ 📍 Address   🔍 Search  👤 │
│ d  ├────────────────────────────┤
│ e  │                            │
│ b  │   Content (max-width 1280) │
│ a  │                            │
│ r  │                            │
│    │                            │
│ 🏠 │                            │
│ 🔍 │                            │
│ 📦 │                            │
│ 👤 │                            │
└────┴────────────────────────────┘
```

### 2.3 Navigation Map

```
                    ┌─────────────┐
                    │  Welcome    │
                    │  (Onboard)  │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  Phone Login │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │  OTP Verify │
                    └──────┬──────┘
                           │
                    ┌──────▼──────┐
                    │ Profile Setup│
                    └──────┬──────┘
                           │
              ┌────────────▼────────────┐
              │       HOME              │◄──────┐
              └─┬───────┬───────┬───────┘       │
                │       │       │               │
         ┌──────▼──┐ ┌──▼───┐ ┌─▼────────┐     │
         │Restaurant│ │Search│ │Orders     │     │
         │ Detail   │ └──┬───┘ │ History   │     │
         └────┬────┘    │     └─────┬──────┘     │
              │         │           │             │
         ┌────▼───┐ ┌───▼────┐ ┌────▼─────┐      │
         │  Item  │ │ Filter │ │  Order   │      │
         │ Detail │ └────────┘ │ Tracking │      │
         └────┬───┘            └──────────┘      │
              │                                  │
         ┌────▼───┐                              │
         │  Cart  │──────────────────────┐       │
         └────┬───┘                       │       │
              │                           │       │
         ┌────▼─────┐                     │       │
         │ Checkout │                     │       │
         └────┬─────┘                     │       │
              │                           │       │
         ┌────▼──────────┐                │       │
         │   Order       │                │       │
         │ Confirmation  │────────────────┴───────┘
         └───────────────┘
```

---

## 3. Screen Specifications

---

### CUS-AUTH-01: Welcome / Onboarding

**Purpose**: First-time intro to the app.

**Wireframe**:
```
┌─────────────────────────────────┐
│                                 │
│                                 │
│         [App Logo]              │
│                                 │
│    اطلب أكل من مطاعمك           │
│       المفضلة                   │
│                                 │
│    توصيل سريع في 30 دقيقة       │
│                                 │
│  ● ○ ○   (page indicators)      │
│                                 │
│                                 │
│                                 │
│  [ابدأ الآن]                    │
│                                 │
│  عندك حساب؟ [سجّل دخول]          │
│                                 │
└─────────────────────────────────┘
```

**Carousel Slides** (3):

| Slide | Image | Title | Subtitle |
|-------|-------|-------|----------|
| 1 | Delivery scooter illustration | اطلب أكل من مطاعمك المفضلة | توصيل سريع في 30 دقيقة |
| 2 | Phone with Vodafone Cash | ادفع بـ Vodafone Cash أو InstaPay | طرق دفع محلية متنوعة |
| 3 | Loyalty card illustration | اكسب نقاط وكاش باك | مع كل طلب |

**Specs**:
- Full-screen (no top bar, no bottom nav)
- Background: subtle gradient (primary → primary-light)
- Carousel: swipeable + dots indicator
- Auto-advance: 3s per slide (pause on user interaction)
- CTA: Primary button (full-width, 56px height)
- Skip option: top-right (text link)
- Language: Arabic primary

**Interactions**:
- Swipe left/right: change slide
- Tap CTA: navigate to Phone Login
- Tap "سجّل دخول": navigate to Phone Login (skip carousel)
- Auto-advance pauses on touch, resumes after 5s of inactivity

**States**:
- Loading: brand logo splash (300ms)
- Error: N/A (static content)

**Accessibility**:
- Carousel: `role="region"`, `aria-roledescription="carousel"`
- Dots: `aria-label="شريحة X من 3"`
- Skip link: `aria-label="تخطي المقدمة"`

---

### CUS-AUTH-02: Phone Login

**Purpose**: Enter phone number to receive OTP.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←                          ?   │
├─────────────────────────────────┤
│                                 │
│   أهلاً! 👋                      │
│                                 │
│   سجّل رقم موبايلك عشان         │
│   نبعتبلك كود تفعيل             │
│                                 │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🇪🇬 +20                   │   │
│   │ [01 2345 6789        ] │   │
│   └─────────────────────────┘   │
│                                 │
│   رقم الموبايل                  │
│                                 │
│                                 │
│   [أرسل الكود]                  │
│                                 │
│                                 │
│   بتسجيلك، أنت موافق على         │
│   [الشروط والأحكام]             │
│   و[سياسة الخصوصية]             │
│                                 │
│                                 │
│   ─── أو سجّل بـ ───             │
│                                 │
│   [G] Google                    │
│   [🍎] Apple                    │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back arrow (←) + help (?)
- Title: "أهلاً! 👋" (display-sm, bold)
- Subtitle: "سجّل رقم موبايلك..." (body, secondary)
- Phone input: country code dropdown (🇪🇬 +20 default) + phone number
- Input height: 56px (large for mobile)
- Input type: `tel` (triggers numeric keyboard)
- Max length: 11 digits (Egyptian mobile)
- CTA: Primary button (full-width, 56px)
- CTA disabled until valid phone (10-11 digits starting with 01)
- Terms: links to legal pages
- Divider: "─── أو سجّل بـ ───"
- Social login: Google + Apple buttons (48px height)

**Validation**:
- Empty: helper text "أدخل رقم الموبايل"
- Invalid: "رقم الموبايل مش صحيح" (red)
- Valid: green checkmark icon

**Interactions**:
- Tap CTA: validate → if valid, send OTP → navigate to OTP Verification
- Tap Google/Apple: trigger OAuth flow
- Tap back: return to Welcome

**States**:
- Loading (sending OTP): button shows spinner + "جاري الإرسال..."
- Error (rate limited): toast "حاول تاني بعد دقيقة"
- Error (network): toast "فيه مشكلة بالنت، حاول تاني"

---

### CUS-AUTH-03: OTP Verification

**Purpose**: Enter 6-digit OTP sent via SMS.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←                          ?   │
├─────────────────────────────────┤
│                                 │
│   أدخل الكود                    │
│                                 │
│   بعتبنالك كود على              │
│   +20 10 1234 5678              │
│   [تغيير الرقم]                 │
│                                 │
│                                 │
│   ┌─┐ ┌─┐ ┌─┐ ┌─┐ ┌─┐ ┌─┐      │
│   │ │ │ │ │ │ │ │ │ │ │      │
│   └─┘ └─┘ └─┘ └─┘ └─┘ └─┘      │
│                                 │
│   ⏱️ 02:00 متبقي                │
│                                 │
│   ماوصلكش الكود؟                │
│   [إعادة إرسال] (disabled)      │
│                                 │
│                                 │
│   [تأكيد]                       │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Title: "أدخل الكود" (h1, bold)
- Phone display: "+20 10 1234 5678" + "تغيير الرقم" link
- OTP input: 6 separate boxes (48×56px each, 8px gap)
- Input type: numeric only
- Auto-advance: focus moves to next box on input
- Auto-submit: when 6 digits entered, auto-verify (no need to tap button)
- Timer: 2:00 countdown (mm:ss)
- Resend: disabled until timer hits 0:00, then enabled
- CTA: Primary button (visible but optional due to auto-submit)

**Validation**:
- Auto-verify after 6 digits entered
- Wrong code: shake animation on boxes + error message "الكود مش صحيح"
- 3 wrong attempts: 5-minute lockout

**Interactions**:
- Type digit: auto-advance to next box
- Backspace: clear current + go back
- Paste: distribute digits across boxes
- Tap resend: new OTP sent, timer resets
- Tap "تغيير الرقم": return to Phone Login

**States**:
- Loading (verifying): button spinner
- Error (wrong code): red border + shake animation
- Error (expired): "انتهت صلاحية الكود، ابعت كود جديد"
- Success: brief checkmark animation → navigate to Profile Setup (new user) or Home (returning user)

**Accessibility**:
- OTP boxes: `aria-label="رقم 1"`, `aria-label="رقم 2"`, etc.
- Auto-announce: "تم إدخال X أرقام"
- Error: `aria-live="assertive"`

---

### CUS-AUTH-04: Profile Setup

**Purpose**: New user enters basic info.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←                          ?   │
├─────────────────────────────────┤
│                                 │
│   خلينا نتعرف عليك 🙌            │
│                                 │
│   ┌─────────────────────────┐   │
│   │ الاسم                    │   │
│   │ [أحمد محمد           ]  │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ البريد الإلكتروني (اختياري)│   │
│   │ [ahmed@example.com   ]  │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📍 موقعك الحالي           │   │
│   │ [اختر عنوان التوصيل]      │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ إيه أكلك المفضل؟          │   │
│   │                          │   │
│   │ ☑ مصري                   │   │
│   │ ☑ إيطالي                 │   │
│   │ ☐ آسيوي                  │   │
│   │ ☑ وجبات سريعة             │   │
│   │ ☐ صحي                    │   │
│   │ ☐ حلويات                 │   │
│   └─────────────────────────┘   │
│                                 │
│   [خلينا نبدأ 🚀]               │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Title: "خلينا نتعرف عليك 🙌" (h1, bold)
- Name input: required, min 2 chars
- Email input: optional, validate format
- Address: tap to open Address Picker modal
- Cuisine preferences: checkboxes (multi-select), feeds recommendation engine
- CTA: "خلينا نبدأ 🚀" (primary, full-width)

**Interactions**:
- Tap address: open Address Picker (modal)
- Tap cuisines: toggle checkmarks
- Tap CTA: validate → save → navigate to Home

**Validation**:
- Name required (min 2 chars)
- Address required (must pick before continue)
- Email optional but if entered, must be valid format

---

### CUS-HOME-01: Home / Discovery ⭐ (Most Important)

**Purpose**: Main feed showing restaurants, categories, promotions.

**Wireframe (Mobile)**:
```
┌─────────────────────────────────┐
│ 📍 الزمالك، 26 يوليو  🔍  👤   │
├─────────────────────────────────┤
│                                 │
│ ┌─────────────────────────────┐ │
│ │ 🎁 خصم 50% على أول طلب       │ │
│ │ استخدم الكود: WELCOME50      │ │
│ └─────────────────────────────┘ │
│                                 │
├─────────────────────────────────┤
│                                 │
│ 🔥 رائج قريب منك                 │
│                                 │
│ ┌────┐ ┌────┐ ┌────┐ ┌────┐    │
│ │[img]│ │[img]│ │[img]│ │[img]│   │
│ │ McD │ │Pizza│ │Sushi│ │KFC │   │
│ │4.6⭐│ │4.8⭐│ │4.5⭐│ │4.4⭐│   │
│ │25د │ │35د │ │20د │ │30د │   │
│ │EGP85│ │EGP145│ │EGP90│ │EGP120│   │
│ └────┘ └────┘ └────┘ └────┘    │
│  ← horizontal scroll →          │
│                                 │
├─────────────────────────────────┤
│                                 │
│ ⏰ وصل في 30 دقيقة               │
│                                 │
│ ┌────┐ ┌────┐ ┌────┐           │
│ │ ⚡ │ │ ⚡ │ │ ⚡ │           │
│ │McD │ │KFC │ │Subway│          │
│ └────┘ └────┘ └────┘           │
│                                 │
├─────────────────────────────────┤
│                                 │
│ 🍽️ أنواع المطابخ                  │
│                                 │
│ ┌────┐ ┌────┐ ┌────┐ ┌────┐    │
│ │🍔  │ │🍕  │ │🥗  │ │🍗  │    │
│ │مصري│ │إيطالي│ │آسيوي│ │وجبات│   │
│ └────┘ └────┘ └────┘ └────┘    │
│ ┌────┐ ┌────┐ ┌────┐ ┌────┐    │
│ │🍰  │ │🥗  │ │🌯  │ │🍢  │    │
│ │حلوى│ │صحي │ │شاورما│ │كباب│   │
│ └────┘ └────┘ └────┘ └────┘    │
│                                 │
├─────────────────────────────────┤
│                                 │
│ 🏆 الأعلى تقييماً                 │
│                                 │
│ ┌────┐ ┌────┐ ┌────┐           │
│ │4.9⭐│ │4.9⭐│ │4.8⭐│           │
│ └────┘ └────┘ └────┘           │
│                                 │
├─────────────────────────────────┤
│                                 │
│ ⏮️ اطلب تاني                     │
│ (returning users only)          │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ [img] Pizza Hut              │ │
│ │ آخر طلب: Margherita + Coke  │ │
│ │ EGP 285 - 3 أيام              │ │
│ │ [اطلب نفس الطلب]              │ │
│ └─────────────────────────────┘ │
│                                 │
├─────────────────────────────────┤
│ [🏠] [🔍] [📦] [👤]              │
└─────────────────────────────────┘
```

**Sections (in order)**:

| # | Section | Type | Source |
|---|---------|------|--------|
| 1 | Welcome Banner | Promo (new users only) | If user has 0 orders |
| 2 | Reorder (Last Orders) | Horizontal scroll | Returning users only |
| 3 | Trending Near You | Horizontal scroll | Top 8-10 restaurants |
| 4 | Under 30 Minutes | Horizontal scroll | ETA <30min restaurants |
| 5 | Cuisines | Grid (4-col) | All cuisine types |
| 6 | Top Rated | Horizontal scroll | Rating >4.7 |
| 7 | Offers Near You | Horizontal scroll | Active promos |

**Specs**:
- Top bar: 56px height, address (tap to change), search icon, profile avatar
- Promo banner: 16:9 aspect, primary-color background, dismissible (X)
- Section headers: h2, bold, with "عرض الكل →" link if applicable
- Horizontal scroll: snap to card, peek next card (8px)
- Cards: see RestaurantCard spec below

**RestaurantCard Component**:
```
┌─────────────┐
│             │
│   [Image]   │  ← 16:9, lazy-loaded
│             │
├─────────────┤
│ McDonald's  │  ← h3, bold
│ 🍔 برجر     │  ← caption, secondary
│ ⭐ 4.6 (1.2K)│  ← caption
│ 📍 1.2km     │  ← caption
│ ⏱️ 25-35 د   │  ← caption
│ 💰 EGP 85 للشخصين│ ← caption
│ 🎁 خصم 20%  │  ← badge (optional)
└─────────────┘
```

- Card width: 280px (desktop), 70% of viewport (mobile)
- Card padding: 0 (image) + 12px (content)
- Image: 16:9, border-radius top only (12px)
- Border radius: 12px
- Shadow: md (default), lg (hover)
- Hover: translateY(-2px), shadow-lg, 200ms

**CuisineCard Component**:
```
┌─────────┐
│   🍔    │  ← icon (32px)
│         │
│  مصري   │  ← caption, bold
└─────────┘
```
- Card: square (80×80px mobile, 100×100 desktop)
- Background: bg-secondary
- Border radius: 12px
- Icon: Material Symbol, 32px, primary color
- Tap: navigate to Category Browse

**Interactions**:
- Tap address: open Address Picker
- Tap search icon: navigate to Search
- Tap avatar: navigate to Profile
- Tap restaurant card: navigate to Restaurant Detail
- Tap cuisine card: navigate to Category Browse
- Pull to refresh (mobile): reload all sections
- Scroll: lazy-load more content (infinite scroll)

**States**:
- Loading: skeleton cards (5-7 placeholder cards)
- Empty (no restaurants): "مفيش مطاعم في منطقتك دلوقتي. جرّب وسّع البحث."
- Error: "حصل خطأ. اسحب للتحديث."
- Offline: "مفيش نت. اتصل بالإنترنت وحاول تاني."

**Personalization**:
- Section order changes based on user history
- Returning users see "Reorder" section first
- Time-of-day affects: lunch hours → "Quick Lunch" section
- Weather: rain → "Warm Soups" section

**Accessibility**:
- Section headers: `role="heading"`, `aria-level="2"`
- Horizontal scroll: `role="region"`, `aria-label="رائج قريب منك"`
- Cards: `role="button"`, `aria-label="مطعم McDonald's، تقييم 4.6، 25-35 دقيقة"`

---

### CUS-HOME-02: Address Picker

**Purpose**: Select delivery address.

**Wireframe (Bottom Sheet)**:
```
┌─────────────────────────────────┐
│                                 │
│   ━━━ (drag handle)             │
│                                 │
│   اختر عنوان التوصيل             │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📍 موقعك الحالي            │   │
│   │ حدد موقعك تلقائياً         │   │
│   └─────────────────────────┘   │
│                                 │
│   ─── العناوين المحفوظة ───     │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🏠 المنزل                  │   │
│   │ الزمالك، 26 يوليو، شقة 5  │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 💼 الشغل                   │   │
│   │ مدينة نصر، عباس العقاد     │   │
│   └─────────────────────────┘   │
│                                 │
│   [➕ إضافة عنوان جديد]          │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Bottom sheet: slides up from bottom, 70% viewport height
- Drag handle: visible at top
- Title: "اختر عنوان التوصيل"
- Current location: tap to use GPS
- Saved addresses: list with label + street
- Add new: button at bottom

**Interactions**:
- Tap current location: request GPS permission → get location → close sheet
- Tap saved address: select → close sheet → update home feed
- Tap add new: open Address Form (modal)
- Swipe down: dismiss sheet

---

### CUS-HOME-03: Search Results

**Purpose**: Search restaurants and dishes.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  [🔍 بيتزا              ] ✕  │
├─────────────────────────────────┤
│                                 │
│   47 نتيجة لـ "بيتزا"            │
│                                 │
│   [الكل] [مطاعم] [أطباق]        │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Pizza Hut          │   │
│   │ 🍕 إيطالي • ⭐4.6 (1.2K)  │   │
│   │ 📍 1.2km • ⏱️ 30-40 د     │   │
│   │ 💰 EGP 145 للشخصين        │   │
│   │ 🎁 خصم 20%                │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Domino's           │   │
│   │ 🍕 إيطالي • ⭐4.5 (890)   │   │
│   │ 📍 2.1km • ⏱️ 35-45 د     │   │
│   └─────────────────────────┘   │
│                                 │
│   ...                           │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back arrow + search input + clear (✕)
- Search input: auto-focus, persistent keyboard
- Result count: "47 نتيجة لـ 'بيتزا'"
- Tabs: الكل / مطاعم / أطباق
- Results: vertical list of restaurant cards (full-width)
- Filters: floating button (bottom-right) for advanced filters

**Interactions**:
- Type: debounced search (300ms)
- Tap clear (✕): clear input, show recent searches
- Tap tab: switch filter
- Tap card: navigate to Restaurant Detail
- Tap filter button: open Filter & Sort sheet

**States**:
- Loading: skeleton cards
- Empty: "مفيش نتائج لـ 'بيتزا'. جرب كلمة تانية."
- Recent searches (when input empty): last 5 searches

---

### CUS-HOME-04: Category Browse

**Purpose**: Browse restaurants by cuisine.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  إيطالي                     │
├─────────────────────────────────┤
│                                 │
│   [الأكثر شهرة ▾] [⏱️ الأسرع]    │
│   [⭐ الأعلى تقييماً] [💰 الأرخص]  │
│                                 │
│   23 مطعم إيطالي                 │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Pizza Hut          │   │
│   │ ...                      │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Domino's           │   │
│   │ ...                      │   │
│   └─────────────────────────┘   │
│                                 │
│   ...                           │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back arrow + category name ("إيطالي")
- Sort chips: horizontal scroll (Most Popular, Fastest, Top Rated, Cheapest)
- Count: "23 مطعم إيطالي"
- List: vertical restaurant cards

**Interactions**:
- Tap sort chip: re-sort list
- Tap card: navigate to Restaurant Detail

---

### CUS-HOME-05: Filter & Sort

**Purpose**: Advanced filtering of restaurants.

**Wireframe (Bottom Sheet)**:
```
┌─────────────────────────────────┐
│   ━━━                           │
│                                 │
│   فلترة النتائج                  │
│                                 │
│   التقييم                        │
│   ○ الكل                         │
│   ○ 4.5+ ⭐                      │
│   ● 4.0+ ⭐                      │
│                                 │
│   وقت التوصيل                    │
│   ○ الكل                         │
│   ● أقل من 30 دقيقة              │
│   ○ أقل من 45 دقيقة              │
│                                 │
│   السعر للشخصين                   │
│   ○ الكل                         │
│   ○ أقل من EGP 100              │
│   ● EGP 100 - 200               │
│   ○ EGP 200+                    │
│                                 │
│   عروض فقط                       │
│   [Toggle: OFF]                 │
│                                 │
│   [إلغاء]    [تطبيق (23)]       │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Bottom sheet: 80% viewport height
- Filter groups: Rating, Delivery Time, Price, Offers Only
- Apply button: shows count of matching results
- Reset: link at top

---

### CUS-REST-01: Restaurant Detail ⭐

**Purpose**: View restaurant info and menu.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←                       [♡] [↗]│
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │                         │   │
│   │    [Cover Image 16:9]   │   │
│   │                         │   │
│   └─────────────────────────┘   │
│                                 │
│   Pizza Hut                     │
│   ⭐ 4.6 (1,243 تقييم)           │
│   🍕 إيطالي • 📍 1.2km          │
│   ⏱️ 30-40 دقيقة                 │
│   💰 EGP 85 للشخصين             │
│                                 │
│   🎁 خصم 20% على الطلبات >200    │
│                                 │
├─────────────────────────────────┤
│   🔍 ابحث في المنيو              │
├─────────────────────────────────┤
│   [🍕 بيتزا] [🥤 مشروبات]        │
│   [🥗 مقبلات] [🍰 حلويات]        │
│   ← sticky category tabs →      │
├─────────────────────────────────┤
│                                 │
│   🍕 بيتزا                       │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Margherita         │   │
│   │ صلصة طماطم، موتزاريلا     │   │
│   │ ⭐ 4.7 • EGP 145          │   │
│   │                  [➕ إضافة]│   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Pepperoni          │   │
│   │ بسطرمة لحم، جبن إضافي     │   │
│   │ ⭐ 4.8 • EGP 165 🔥 الأكثر │   │
│   │                  [➕ إضافة]│   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Veggie Supreme     │   │
│   │ خضار مشكل                │   │
│   │ ⭐ 4.5 • EGP 155          │   │
│   │                  [➕ إضافة]│   │
│   └─────────────────────────┘   │
│                                 │
│   🥤 مشروبات                     │
│   ...                           │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back arrow + favorite (♡) + share (↗)
- Cover image: 16:9, lazy-loaded, full-width
- Restaurant header: name (h1), rating + count, cuisine + distance, ETA, price range
- Promo badge: if active, primary-color background
- Search menu: sticky input below header
- Category tabs: sticky, horizontal scroll, active tab highlighted
- Menu items: list with image (left) + info (right) + add button

**MenuItemCard Component**:
```
┌─────────────────────────────────┐
│ [img]  Margherita                │
│  80px  صلصة طماطم، موتزاريلا     │
│        ⭐ 4.7 • EGP 145      [➕]│
└─────────────────────────────────┘
```
- Image: 80×80px (1:1), border-radius 8px
- Title: h4, bold
- Description: caption, secondary (max 2 lines)
- Rating + price: caption
- Add button: icon button (44×44px), primary color
- "🔥 الأكثر طلباً" badge: orange, top-right

**Interactions**:
- Tap back: return to previous screen
- Tap favorite (♡): toggle favorite (heart fills)
- Tap share: open share sheet
- Tap search menu: focus input, filter items live
- Tap category tab: scroll to section + highlight tab
- Tap add (➕): open Item Detail modal
- Scroll: category tabs update based on visible section

**States**:
- Restaurant closed: banner "المطعم مقفل دلوقتي. يفتح 10:00 صباحاً."
- Item unavailable: card grayed out + "نفد" badge
- Loading: skeleton menu items

---

### CUS-REST-02: Item Detail / Customize ⭐

**Purpose**: Customize item before adding to cart.

**Wireframe (Bottom Sheet)**:
```
┌─────────────────────────────────┐
│                                 │
│   ┌─────────────────────────┐   │
│   │                         │   │
│   │    [Item Image 16:9]    │   │
│   │                         │   │
│   └─────────────────────────┘   │
│                                 │
│   Margherita Pizza              │
│   ⭐ 4.7 (234 تقييم)             │
│   صلصة طماطم إيطالية، موتزاريلا  │
│   طازجة، ريحان عضوي             │
│                                 │
│   الحجم (إجباري)                 │
│   ○ صغير (EGP 120)              │
│   ● وسط (EGP 145)               │
│   ○ كبير (EGP 175)              │
│                                 │
│   الإضافات                       │
│   ☑ جبن إضافي (+EGP 15)         │
│   ☐ بيكون (+EGP 20)             │
│   ☐ فطر (+EGP 10)               │
│   ☐ زيتون (+EGP 5)              │
│                                 │
│   التعديلات                      │
│   ☐ بدون بصل                     │
│   ☐ بدون ثوم                     │
│   ☐ حار                          │
│                                 │
│   ملاحظات خاصة                   │
│   ┌─────────────────────────┐   │
│   │ (مثلاً: اتصلوا قبل الوصول) │   │
│   └─────────────────────────┘   │
│                                 │
│   ─────────────────────────     │
│                                 │
│   الكمية         [- 1 +]        │
│                                 │
│   [🛒 أضف للسلة - EGP 160]       │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Bottom sheet: 90% viewport height, scrollable
- Image: 16:9, full-width
- Title: h2, bold
- Rating: caption
- Description: body, secondary
- Required modifiers: marked "إجباري", can't add to cart without selection
- Optional modifiers: checkboxes
- Quantity stepper: - 1 + (min 1, max 20)
- CTA: "🛒 أضف للسلة - EGP {total}" (primary, full-width)
- Total: dynamically updates as user selects

**Validation**:
- Required modifier not selected: CTA disabled + helper text "اختر الحجم"
- Quantity > 20: error "حد أقصى 20 وحدة"

**Interactions**:
- Tap radio/checkbox: toggle, update total
- Type notes: max 200 chars
- Tap quantity +/-: update count + total
- Tap CTA: add to cart, close sheet, show toast "اتضاف للسلة"

---

### CUS-CART-01: Cart ⭐

**Purpose**: Review cart before checkout.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  السلة                       │
├─────────────────────────────────┤
│                                 │
│   📍 التوصيل إلى:                │
│   الزمالك، 26 يوليو، شقة 5      │
│   [تغيير]                       │
│                                 │
├─────────────────────────────────┤
│                                 │
│   Pizza Hut                     │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Margherita (وسط)   │   │
│   │ + جبن إضافي              │   │
│   │ + فطر                    │   │
│   │ EGP 170 [- 1 +] [🗑️]     │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Pepperoni (كبير)   │   │
│   │ EGP 165 [- 2 +] [🗑️]     │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ [img] Coca Cola          │   │
│   │ EGP 25 [- 2 +] [🗑️]      │   │
│   └─────────────────────────┘   │
│                                 │
├─────────────────────────────────┤
│                                 │
│   🎁 كوبون خصم                   │
│   ┌──────────────┐ [تطبيق]      │
│   │ WELCOME50    │              │
│   └──────────────┘              │
│   ✅ مطبق - وفّرت EGP 50         │
│                                 │
├─────────────────────────────────┤
│                                 │
│   تفاصيل الفاتورة:               │
│   Subtotal:           EGP 550   │
│   رسوم التوصيل:       EGP 25    │
│   رسوم الخدمة (5%):    EGP 27.5  │
│   ض.ق.م (14%):         EGP 84.1  │
│   الخصم:             -EGP 50    │
│   ─────────────────────────     │
│   الإجمالي:           EGP 636.6 │
│                                 │
│   💰 كاش باك: EGP 31 (5%)       │
│                                 │
├─────────────────────────────────┤
│                                 │
│   🎁 تستحق خصم EGP 30 لو         │
│   طلبت بـ EGP 700               │
│   [➕ اطلب أكتر]                 │
│                                 │
├─────────────────────────────────┤
│                                 │
│   [🛒 اطلب الآن - EGP 636.6]     │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + "السلة"
- Address section: with "تغيير" link
- Restaurant name: h3, bold
- Cart items: image (60×60), name + modifiers, price + quantity + remove
- Quantity stepper: - N + (min 1, max 20)
- Coupon input: text field + apply button
- Price breakdown: clear, line by line
- Cashback: highlighted (success color)
- "Spend more" prompt: warning color, with CTA
- CTA: "🛒 اطلب الآن - EGP {total}" (primary, full-width, sticky bottom)

**Interactions**:
- Tap change address: open Address Picker
- Tap quantity +/-: update cart + totals
- Tap remove (🗑️): confirm dialog → remove item
- Type coupon + apply: validate → if valid, apply discount
- Tap "اطلب أكتر": stay on cart, prompt for items
- Tap CTA: navigate to Checkout

**Validation**:
- Empty cart: CTA disabled
- Restaurant closed: banner "المطعم مقفل. اختر مطعم تاني."
- Minimum order not met: "الحد الأدنى للطلب EGP 50"

**States**:
- Loading (applying coupon): spinner on apply button
- Coupon invalid: error text "الكوبون مش صحيح أو منتهي"
- Cart empty: empty state with "تصفح المطاعم" CTA

---

### CUS-CART-02: Checkout ⭐

**Purpose**: Final review and place order.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  إتمام الطلب                  │
├─────────────────────────────────┤
│                                 │
│   📍 عنوان التوصيل                │
│   ┌─────────────────────────┐   │
│   │ الزمالك، 26 يوليو، شقة 5  │   │
│   │ علامة: باب أزرق           │   │
│   │ [تغيير]                   │   │
│   └─────────────────────────┘   │
│                                 │
│   ⏰ وقت التوصيل                  │
│   ┌─────────────────────────┐   │
│   │ ● في أسرع وقت (35-45 د)   │   │
│   │ ○ مجدول لوقت لاحق         │   │
│   └─────────────────────────┘   │
│                                 │
│   💳 طريقة الدفع                  │
│   ┌─────────────────────────┐   │
│   │ ● 💚 Vodafone Cash         │   │
│   │   الرصيد: EGP 1,250        │   │
│   │ ○ 🟣 InstaPay             │   │
│   │ ○ 💳 Card (**** 4521)      │   │
│   │ ○ 💵 Cash on Delivery     │   │
│   │   (+EGP 5 رسوم)            │   │
│   │ ➕ إضافة طريقة دفع          │   │
│   └─────────────────────────┘   │
│                                 │
│   🎁 كوبون                       │
│   ✅ WELCOME50 مطبق              │
│                                 │
│   📝 ملاحظات للمطعم               │
│   ┌─────────────────────────┐   │
│   │ (مثلاً: صلصة إضافية)       │   │
│   └─────────────────────────┘   │
│                                 │
│   📞 رقم التواصل                  │
│   01012345678 [تغيير]            │
│                                 │
├─────────────────────────────────┤
│                                 │
│   الإجمالي: EGP 636.6            │
│                                 │
│   [🔒 تأكيد الطلب]                │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + "إتمام الطلب"
- Sections: Address, Time, Payment, Coupon, Notes, Phone
- Payment methods: radio list with icons
- Vodafone Cash: shows current balance (if linked)
- COD: shows fee warning
- CTA: "🔒 تأكيد الطلب" (primary, full-width, sticky bottom)
- Lock icon = security/trust signal

**Interactions**:
- Tap change address: open Address Picker
- Tap scheduled time: open time picker
- Tap payment method: select
- Tap add payment: open Payment Method screen
- Tap CTA: process payment → navigate to Order Confirmation

**States**:
- Loading (processing): button spinner + "جاري تأكيد الطلب..."
- Error (payment failed): toast + retry
- Error (fraud detected): modal explaining + contact support
- Error (restaurant closed): "المطعم قفل. ارجع اختار مطعم تاني."

---

### CUS-CART-04: Order Confirmation

**Purpose**: Success screen after order placed.

**Wireframe**:
```
┌─────────────────────────────────┐
│                                 │
│                                 │
│                                 │
│         ✅                      │
│      (animation)                │
│                                 │
│      تم الطلب!                   │
│                                 │
│   طلبك #A7X92F اتأكد             │
│   المطعم هيبدأ تحضيره            │
│                                 │
│   ⏱️ الوقت المتوقع: 8:35 PM       │
│                                 │
│   💰 EGP 636.6 - Vodafone Cash  │
│                                 │
│   💰 كاش باك: EGP 31             │
│   (هيتضاف لمحفظتك بعد التوصيل)    │
│                                 │
│                                 │
│   [📍 تتبع الطلب]                │
│                                 │
│   [العودة للرئيسية]              │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Full-screen modal (no top bar, no bottom nav)
- Success animation: checkmark draw + confetti (1.5s)
- Title: "تم الطلب!" (display-sm, bold)
- Order details: number, ETA, total, payment
- Cashback info: success color
- Primary CTA: "📍 تتبع الطلب"
- Secondary CTA: "العودة للرئيسية"

**Interactions**:
- Auto-redirect to tracking after 5s (if user doesn't tap)
- Tap track: navigate to Order Tracking
- Tap home: navigate to Home

---

### CUS-ORDER-01: Order Tracking ⭐

**Purpose**: Real-time order status with live map.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  الطلب #A7X92F                │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │ حالة الطلب                │   │
│   │                          │   │
│   │ ✅ تم استلام الطلب         │   │
│   │ ✅ المطعم بدأ التحضير      │   │
│   │ 🟡 الأكل جاهز              │   │
│   │ ⬜ المندوب في الطريق للمطعم │   │
│   │ ⬜ المندوب استلم الطلب     │   │
│   │ ⬜ في الطريق إليك           │   │
│   │ ⬜ تم التوصيل              │   │
│   │                          │   │
│   │ ⏱️ 8:35 PM (15 دقيقة)     │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │                         │   │
│   │   [Live Map]            │   │
│   │                         │   │
│   │   🍔 المطعم              │   │
│   │    ↓                    │   │
│   │   🛵 المندوب (متحرك)     │   │
│   │    ↓                    │   │
│   │   🏠 موقعك               │   │
│   │                         │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 🛵 المندوب                 │   │
│   │ [photo] Mahmoud S.       │   │
│   │ ⭐ 4.8 • موتوسيكل          │   │
│   │ [📞 اتصال] [💬 رسالة]      │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ 📦 تفاصيل الطلب            │   │
│   │ Pizza Hut                │   │
│   │ • Margherita (وسط)       │   │
│   │ • Pepperoni (كبير) × 2    │   │
│   │ • Coca Cola × 2          │   │
│   │ الإجمالي: EGP 636.6       │   │
│   │ الدفع: Vodafone Cash     │   │
│   └─────────────────────────┘   │
│                                 │
│   [❓ مساعدة] [🚫 إلغاء الطلب]   │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + order number
- Status timeline: 7 steps, current highlighted (yellow), completed (green check)
- ETA: prominent, updates every 30s
- Map: Google Maps embed, real-time driver position (5s updates)
- Driver card: photo, name, rating, vehicle, call/message buttons
- Order details: collapsible
- Footer: help + cancel (cancel hidden if past cancellation window)

**Status Timeline Component**:
```
✅ ── ✅ ── 🟡 ── ⬜ ── ⬜ ── ⬜ ── ⬜
تم  بدأ  جاهز  في    استلم  في    تم
الطلب التحضير     الطريق الطلب الطريق التوصيل
```
- Horizontal (mobile) or vertical (desktop)
- Completed: green check + green line
- Current: yellow dot + pulse animation
- Future: gray dot + gray line

**Interactions**:
- Tap call: open phone dialer (anonymized number)
- Tap message: open in-app chat with driver
- Tap cancel: confirm dialog → if past window, explain policy
- Map: pan/zoom allowed

**States**:
- Loading (initial): skeleton
- Driver not assigned yet: "بنبحث لك عن مندوب..."
- Delayed (>10min over ETA): warning banner + apology + new ETA
- Delivered: "تم التوصيل 🎉" + rate prompt

**Real-time Updates**:
- Status changes: animate to next step
- Driver location: smooth pan on map
- ETA: countdown updates every 30s

---

### CUS-ORDER-02: Order History

**Purpose**: List of past orders.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  طلباتي                       │
├─────────────────────────────────┤
│                                 │
│   [النشطة (1)] [السابقة (47)]   │
│                                 │
│   🟡 طلب نشط                     │
│   ┌─────────────────────────┐   │
│   │ #A7X92F  Pizza Hut       │   │
│   │ 3 أصناف - EGP 636.6      │   │
│   │ ⏱️ ETA: 8:35 PM (15 د)    │   │
│   │ [تتبع الطلب]              │   │
│   └─────────────────────────┘   │
│                                 │
│   ─── أمس ───                    │
│                                 │
│   ┌─────────────────────────┐   │
│   │ #B3K45L  McDonald's      │   │
│   │ 2 أصناف - EGP 145        │   │
│   │ ⭐⭐⭐⭐⭐ تم التوصيل 9:20   │   │
│   │ [إعادة الطلب] [تقييم]     │   │
│   └─────────────────────────┘   │
│                                 │
│   ┌─────────────────────────┐   │
│   │ #C8M72N  KFC             │   │
│   │ 4 أصناف - EGP 320        │   │
│   │ ⭐⭐⭐⭐ تم التوصيل 7:15     │   │
│   │ [إعادة الطلب] [تقييم]     │   │
│   └─────────────────────────┘   │
│                                 │
│   ─── قبل أمس ───                │
│   ...                           │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + "طلباتي"
- Tabs: Active / Past
- Active orders: highlighted with status
- Past orders: grouped by date ("أمس", "قبل أمس", "هذا الأسبوع", etc.)
- Order card: number, restaurant, items count, total, status, actions

**Order Card Actions**:
- Active: "تتبع الطلب"
- Past delivered: "إعادة الطلب" + "تقييم" (if not rated)
- Past cancelled: "إعادة الطلب" only

**Interactions**:
- Tap order card: navigate to Order Detail
- Tap reorder: add same items to cart → navigate to Cart
- Tap rate: open rating modal

---

### CUS-PROF-01: Profile

**Purpose**: User profile and stats.

**Wireframe**:
```
┌─────────────────────────────────┐
│  ←  حسابي                        │
├─────────────────────────────────┤
│                                 │
│   ┌─────────────────────────┐   │
│   │ [Avatar] Ahmed Mohamed   │   │
│   │ 🥇 Platinum Member        │   │
│   │ 📧 ahmed@example.com     │   │
│   │ 📱 01012345678           │   │
│   │ [تعديل الملف الشخصي]      │   │
│   └─────────────────────────┘   │
│                                 │
│   📊 إحصائياتي                    │
│   ┌─────────────────────────┐   │
│   │ إجمالي الطلبات: 47        │   │
│   │ إجمالي الإنفاق: EGP 8,520 │   │
│   │ النقاط: 8,520 / 15,000    │   │
│   │ المحفظة: EGP 425 (cashback)│   │
│   └─────────────────────────┘   │
│                                 │
│   📍 العناوين                    │
│   • المنزل - الزمالك             │
│   • الشغل - مدينة نصر            │
│   [➕ إضافة عنوان]               │
│                                 │
│   💳 طرق الدفع                    │
│   • Vodafone Cash               │
│   • Card **** 4521              │
│   [➕ إضافة طريقة]               │
│                                 │
│   🎁 كوبوناتي (3)                │
│   • FREEDEL - توصيل مجاني (7 أيام)│
│   • CASHBACK10 - 10% cashback   │
│   • BIRTHDAY50 - EGP 50         │
│                                 │
│   ⚙️ الإعدادات                    │
│   • الإشعارات                   │
│   • اللغة (عربي/English)         │
│   • الخصوصية                    │
│   • مساعدة                      │
│   • تسجيل الخروج                │
│                                 │
└─────────────────────────────────┘
```

**Specs**:
- Top bar: back + "حسابي"
- Profile header: avatar, name, tier badge, contact info
- Stats card: totals (orders, spend, points, wallet)
- Sections: Addresses, Payment Methods, Coupons, Settings
- Each section: list items + add button

**Interactions**:
- Tap edit profile: open edit modal
- Tap address: edit
- Tap payment: edit/remove
- Tap coupon: copy code
- Tap settings item: navigate to respective screen
- Tap logout: confirm dialog → logout → return to Welcome

---

## 4. UX Flows

### 4.1 First-Time User Flow

```
1. App Open
   ↓
2. Welcome Carousel (3 slides)
   ↓
3. Phone Login
   ↓
4. OTP Verification (2 min timer)
   ↓
5. Profile Setup (name, address, cuisines)
   ↓
6. Home (with Welcome Promo banner)
   ↓
7. Browse Restaurants
   ↓
8. Restaurant Detail
   ↓
9. Add Item to Cart (with customization)
   ↓
10. Cart Review
    ↓
11. Checkout (Vodafone Cash default)
    ↓
12. Order Confirmation
    ↓
13. Order Tracking
```

### 4.2 Returning User Flow

```
1. App Open
   ↓
2. Auto-Login (token valid)
   ↓
3. Home (with "Order Again" section first)
   ↓
4. Quick Reorder (1 tap)
   ↓
5. Cart Review
   ↓
6. Checkout (saved payment)
   ↓
7. Order Confirmation
   ↓
8. Order Tracking
```

### 4.3 Reordering Flow

```
1. Tap "اطلب تاني" on past order
   ↓
2. Modal: "إعادة طلب #B3K45L؟"
   - Items list
   - Total
   - [تأكيد] [إلغاء]
   ↓
3. If items still available:
   - Add to cart
   - Navigate to Cart
   ↓
4. If some items unavailable:
   - Modal: "بعض الأصناف مش متاحة"
   - List unavailable items
   - [كمل بغيرها] [إلغاء]
   ↓
5. Cart Review
   ↓
6. Checkout
```

### 4.4 Cancellation Flow

```
1. User taps "إلغاء الطلب" on tracking
   ↓
2. Modal: "إلغاء الطلب #A7X92F؟"
   - Reason dropdown
   - "هل أنت متأكد؟"
   - [تأكيد الإلغاء] [تراجع]
   ↓
3. If within cancellation window:
   - Cancel order
   - Process refund
   - Show confirmation
   ↓
4. If past window:
   - Modal: "مفيش إلغاء بعد ما المطعم يبدأ"
   - [حسناً]
```

---

## 5. Component Library

### App-Specific Components

#### RestaurantCard
```typescript
interface RestaurantCardProps {
  restaurant: Restaurant
  distanceKm: number
  onClick: () => void
}
```

#### MenuItemCard
```typescript
interface MenuItemCardProps {
  item: MenuItem
  onAdd: () => void
}
```

#### CartItem
```typescript
interface CartItemProps {
  item: CartItem
  onQuantityChange: (qty: number) => void
  onRemove: () => void
}
```

#### OrderStatusTimeline
```typescript
interface OrderStatusTimelineProps {
  currentStatus: OrderStatus
  history: StatusChange[]
}
```

#### AddressPicker
```typescript
interface AddressPickerProps {
  addresses: Address[]
  onPick: (address: Address) => void
  onAddNew: () => void
}
```

#### PaymentMethodSelector
```typescript
interface PaymentMethodSelectorProps {
  methods: PaymentMethod[]
  selected: string
  onSelect: (id: string) => void
}
```

---

## 6. Arabic RTL Notes

- All screens default to RTL (`dir="rtl"`)
- Icons flip where direction matters (arrows, sliders)
- Numbers stay LTR (prices, phone numbers)
- Back button: top-right (not top-left)
- Bottom nav: order doesn't change but icons flip
- Carousel: scrolls right-to-left
- Swipe gestures: inverted (swipe right = back)

---

## 7. States Reference

### 7.1 Loading States
- Skeleton screens for all data-driven screens
- Spinners for button actions
- Shimmer for images

### 7.2 Empty States
| Screen | Empty State |
|--------|-------------|
| Home (no restaurants) | "مفيش مطاعم في منطقتك" + CTA to widen search |
| Search (no results) | "مفيش نتائج لـ '{query}'" + suggestions |
| Cart (empty) | "السلة فاضية" + CTA to browse |
| Orders (none) | "مفيش طلبات بعد" + CTA to browse |
| Favorites (none) | "مفيش مطاعم مفضلة" + CTA to browse |

### 7.3 Error States
- Network error: "مفيش نت. اتصل بالإنترنت وحاول تاني."
- Server error: "حصل خطأ من عندنا. حاول تاني."
- Not found: "الصفحة مش موجودة" + CTA to home
- Unauthorized: redirect to login

### 7.4 Offline State
- Banner: "مفيش نت. بيظهر آخر بيانات."
- Cached data shown
- Actions disabled with tooltip "محتاج نت"

---

> **Next**: Read `driver-web/UI-SPEC.md` for Driver Web App specifications.
