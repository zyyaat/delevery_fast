# Design System — Food Delivery Platform

> **Version**: 1.0  
> **Status**: Active — Source of Truth for All UI Work  
> **Last Updated**: 2026-07-04  
> **Based on**: Uber Eats Base Design System, Talabat Egypt, elmenus, DoorDash research

---

## Table of Contents

1. [Design Principles](#1-design-principles)
2. [Brand Identity](#2-brand-identity)
3. [Color Palette](#3-color-palette)
4. [Typography](#4-typography)
5. [Spacing & Layout](#5-spacing--layout)
6. [Border Radius & Elevation](#6-border-radius--elevation)
7. [Component Library](#7-component-library)
8. [Iconography](#8-iconography)
9. [Imagery](#9-imagery)
10. [Motion & Animation](#10-motion--animation)
11. [States (Empty / Loading / Error / Success)](#11-states)
12. [Arabic RTL Considerations](#12-arabic-rtl-considerations)
13. [Responsive Breakpoints](#13-responsive-breakpoints)
14. [Accessibility](#14-accessibility)
15. [Competitive Analysis](#15-competitive-analysis)

---

## 1. Design Principles

### 1.1 Our 7 Design Principles

Inspired by Uber's Base design system ([base.uber.com](https://base.uber.com)) and adapted for the Egyptian market:

1. **Egyptian-First, Global-Standard**
   - Arabic is the primary language (RTL by default)
   - Local payment methods front and center (Vodafone Cash)
   - Egyptian cultural references in copy
   - But quality matches global apps (Uber Eats, DoorDash)

2. **Speed Over Decoration**
   - Every screen loads in <2s
   - Images optimized (WebP, lazy-load)
   - No unnecessary animations
   - Inspired by Uber Eats: minimal UI, fast interactions

3. **Glanceable Information**
   - Users decide in <3s whether to engage
   - Card-based layouts (like Uber Eats)
   - Clear hierarchy: image → name → meta → CTA
   - Inspired by Talabat: dense info, but organized

4. **Trust Through Transparency**
   - Show full pricing breakdown (no hidden fees)
   - Real-time tracking with honesty about delays
   - Restaurant health scores visible
   - Driver earnings transparent

5. **Forgiving Interactions**
   - Every action is reversible (within policy)
   - Clear cancel/refund flows
   - Confirmation for destructive actions
   - "Are you sure?" dialogs for irreversible

6. **Inclusive by Default**
   - WCAG 2.1 AA compliance
   - Large touch targets (44×44px min)
   - Color contrast >4.5:1 for text
   - Screen reader support (ARIA labels in Arabic)

7. **Delightful Micro-interactions**
   - Subtle animations on key actions
   - Haptic feedback (mobile)
   - Sound on order confirmation
   - Inspired by Talabat's energetic feel + Uber's polish

### 1.2 Decision Framework

When in doubt, ask:
- Will this slow the user down? → Remove it
- Will an Egyptian user understand this? → Test with locals
- Does this match Uber Eats quality? → If not, improve
- Is this accessible? → Verify contrast, touch target, ARIA

---

## 2. Brand Identity

### 2.1 Brand Personality

| Trait | Description | Inspiration |
|-------|-------------|-------------|
| **Trustworthy** | Reliable, transparent, safe | Uber Eats (track record) |
| **Local** | Egyptian, Arabic-first | Talabat (local presence) |
| **Energetic** | Fast, alive, hungry | Talabat (orange = energy) |
| **Premium-feel** | Polished, modern | Uber Eats (Base design) |
| **Affordable** | Not luxury, accessible | elmenus (broad audience) |

### 2.2 Logo Treatment

- **Primary**: Wordmark "طلبات" + English "TALABAT" lockup (placeholder — to be designed)
- **Symbol**: Stylized fork/utensil forming the Arabic letter "ط"
- **Safe area**: 2x logo height on all sides
- **Min size**: 24px width (digital), 8mm (print)
- **Clear space**: No other elements within safe area

### 2.3 Voice & Tone

| Context | Tone | Example |
|---------|------|---------|
| Marketing | Energetic, friendly | "جوعان؟ احنا جينالك!" |
| Onboarding | Warm, guiding | "أهلاً! خلينا نضبط موقعك الأول" |
| Errors | Apologetic, helpful | "حصل خطأ صغير، حاول تاني" |
| Success | Celebratory, brief | "تم! طلبك في الطريق 🛵" |
| Cancellations | Understanding | "تمام، اتلغي. نرجع في أي وقت" |
| Notifications | Concise, actionable | "طلبك جاهز للاستلام 📦" |

---

## 3. Color Palette

### 3.1 Primary Brand Colors

Inspired by Talabat's orange (energy) + Uber Eats green (freshness), we chose a **vibrant orange-red** as primary (Egyptian appetite-stimulating) + **deep teal** as secondary (trust).

| Role | Name | Hex | RGB | Usage |
|------|------|-----|-----|-------|
| **Primary** | Tandoor Orange | `#FF5722` | (255, 87, 34) | CTAs, active states, brand |
| **Primary Dark** | Ember | `#E64A19` | (230, 74, 25) | Hover, pressed |
| **Primary Light** | Sandstone | `#FFAB91` | (255, 171, 145) | Backgrounds, badges |
| **Secondary** | Nile Teal | `#00897B` | (0, 137, 123) | Trust elements, links |
| **Secondary Dark** | Deep Sea | `#00695C` | (0, 105, 92) | Hover |
| **Secondary Light** | Mist | `#B2DFDB` | (178, 223, 219) | Backgrounds |

### 3.2 Semantic Colors

| Role | Name | Hex | Usage |
|------|------|-----|-------|
| **Success** | Oasis Green | `#2E7D32` | Order confirmed, delivered |
| **Success Light** | Mint | `#C8E6C9` | Success backgrounds |
| **Warning** | Sun Gold | `#F9A825` | Delays, attention needed |
| **Warning Light** | Honey | `#FFF9C4` | Warning backgrounds |
| **Error** | Hibiscus | `#D32F2F` | Errors, cancellations |
| **Error Light** | Rose | `#FFCDD2` | Error backgrounds |
| **Info** | Sky Blue | `#0288D1` | Information, neutral alerts |
| **Info Light** | Ice | `#B3E5FC` | Info backgrounds |

### 3.3 Neutrals

Following Material Design 3 neutral scale, adapted for warm Egyptian aesthetic:

| Role | Name | Hex | Usage |
|------|------|-----|-------|
| **BG Primary** | Papyrus | `#FAFAFA` | App background |
| **BG Secondary** | Linen | `#F5F5F5` | Cards, sections |
| **BG Tertiary** | Stone | `#EEEEEE` | Hover, dividers |
| **Surface** | White | `#FFFFFF` | Modals, elevated cards |
| **Border** | Smoke | `#E0E0E0` | Default borders |
| **Border Strong** | Ash | `#BDBDBD` | Focused borders |
| **Text Primary** | Charcoal | `#212121` | Headings, body |
| **Text Secondary** | Slate | `#616161` | Meta, captions |
| **Text Tertiary** | Fog | `#9E9E9E` | Placeholders |
| **Text Disabled** | Haze | `#BDBDBD` | Disabled states |

### 3.4 Dark Mode Palette (for Employee Portal & Command Center)

| Role | Name | Hex | Usage |
|------|------|-----|-------|
| **BG Primary** | Midnight | `#0A0E1A` | App background |
| **BG Secondary** | Carbon | `#131A2E` | Cards |
| **Surface** | Graphite | `#1A2240` | Elevated cards |
| **Border** | Slate | `#2A3656` | Default borders |
| **Text Primary** | Snow | `#EAF0FF` | Body text |
| **Text Secondary** | Sky | `#9BA8C7` | Meta |
| **Accent** | Cyan Neon | `#00D4FF` | Primary in dark |
| **Accent 2** | Purple Neon | `#B14EFF` | Secondary in dark |

### 3.5 Data Visualization Colors

For Command Center charts:

```
Series 1: #FF5722 (Tandoor Orange)
Series 2: #00897B (Nile Teal)
Series 3: #F9A825 (Sun Gold)
Series 4: #0288D1 (Sky Blue)
Series 5: #7B1FA2 (Royal Purple)
Series 6: #C2185B (Rose Pink)
Series 7: #388E3C (Forest Green)
Series 8: #5D4037 (Earth Brown)
```

### 3.6 Competitive Color Analysis

| Brand | Primary | RGB | Personality |
|-------|---------|-----|-------------|
| Uber Eats | `#06C167` | (6, 193, 103) | Fresh, healthy, global |
| Talabat | `#FF5A00` | (255, 90, 0) | Energetic, urgent |
| elmenus | `#FF4F00` | (255, 79, 0) | Food-focused |
| DoorDash | `#FF3008` | (255, 48, 8) | Bold, fast |
| **Ours** | `#FF5722` | (255, 87, 34) | Warm, energetic, Egyptian |

We chose orange-red because:
1. Appetite-stimulating (color psychology)
2. Differentiates from Uber Eats green
3. Aligns with Talabat/DoorDash (proven in delivery)
4. Warm tone fits Egyptian aesthetic
5. Stands out in app store screenshots

---

## 4. Typography

### 4.1 Font Stack

| Language | Font | Weights | Why |
|----------|------|---------|-----|
| **Arabic** | Cairo | 400, 500, 600, 700, 800 | Modern, clean, supports Egyptian Arabic. Free Google Font. Used by Talabat. |
| **English** | Inter | 400, 500, 600, 700, 800 | Geometric, neutral, pairs with Cairo. Used by Uber-style apps. |
| **Numbers** | JetBrains Mono | 400, 500, 600, 700 | Tabular figures for prices/timers. Inspired by Uber Eats. |

### 4.2 Font Loading

```css
/* Google Fonts */
@import url('https://fonts.googleapis.com/css2?family=Cairo:wght@400;500;600;700;800&family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500;600;700&display=swap');
```

### 4.3 Type Scale

Based on 1.250 (Major Third) modular scale, optimized for mobile-first:

| Token | Size | Line Height | Weight | Usage |
|-------|------|-------------|--------|-------|
| `display-lg` | 48px | 56px | 800 | Hero marketing |
| `display-md` | 40px | 48px | 800 | Page titles |
| `display-sm` | 32px | 40px | 700 | Section headers |
| `h1` | 28px | 36px | 700 | Screen titles |
| `h2` | 24px | 32px | 700 | Section titles |
| `h3` | 20px | 28px | 600 | Card titles |
| `h4` | 18px | 24px | 600 | Subtitles |
| `body-lg` | 16px | 24px | 400 | Lead body |
| `body` | 14px | 20px | 400 | Default body |
| `body-sm` | 13px | 18px | 400 | Secondary text |
| `caption` | 12px | 16px | 500 | Captions, labels |
| `overline` | 11px | 16px | 600 | Overlines, tags |
| `mono-lg` | 16px | 24px | 600 | Prices (large) |
| `mono` | 14px | 20px | 500 | Prices, timers |
| `mono-sm` | 12px | 16px | 500 | Timestamps |

### 4.4 Type Usage Examples

```html
<!-- Display -->
<h1 class="text-display-md font-extrabold">اطلب أكل من مطاعمك المفضلة</h1>

<!-- Heading -->
<h2 class="text-h2 font-bold">رائج قريب منك</h2>

<!-- Body -->
<p class="text-body text-text-primary">بيتزا طازجة من فرن حجر، تصل في 30 دقيقة</p>

<!-- Price (mono) -->
<span class="text-mono-lg font-semibold">EGP 145</span>

<!-- Caption -->
<span class="text-caption text-text-secondary">آخر تحديث: 2 دقيقة</span>
```

### 4.5 Arabic Typography Rules

- **Always RTL** for Arabic content (`dir="rtl"`)
- Use Cairo font for Arabic (designed for screens)
- Numbers stay LTR even in RTL context (`direction: ltr` on numeric spans)
- Line height +2px for Arabic (Arabic glyphs need more vertical space)
- Letter spacing: 0 for Arabic (default), -0.01em for Latin headings
- Avoid ALL CAPS for Arabic (no equivalent)

### 4.6 Competitive Typography Analysis

| Brand | Arabic Font | English Font | Style |
|-------|-------------|--------------|-------|
| Uber Eats | Noto Sans Arabic | Uber Move | Geometric, clean |
| Talabat | Cairo | Cairo | Modern, friendly |
| elmenus | Cairo | Open Sans | Friendly, readable |
| DoorDash | Noto Sans Arabic | DDOak | Bold, distinctive |
| **Ours** | **Cairo** | **Inter** | Modern, Egyptian |

---

## 5. Spacing & Layout

### 5.1 Spacing Scale (4px base)

| Token | Value | Usage |
|-------|-------|-------|
| `space-0` | 0 | No spacing |
| `space-1` | 4px | Tight inline |
| `space-2` | 8px | Inline elements |
| `space-3` | 12px | Compact padding |
| `space-4` | 16px | Default padding |
| `space-5` | 20px | Card padding |
| `space-6` | 24px | Section spacing |
| `space-8` | 32px | Large section |
| `space-10` | 40px | Page sections |
| `space-12` | 48px | Hero spacing |
| `space-16` | 64px | Page padding (desktop) |

### 5.2 Grid System

**Mobile (default)**:
- 4-column grid
- 16px gutters
- 16px page padding

**Tablet (≥768px)**:
- 8-column grid
- 24px gutters
- 32px page padding

**Desktop (≥1024px)**:
- 12-column grid
- 32px gutters
- 64px max page padding (centered, max-width 1440px)

### 5.3 Container Widths

| Name | Max Width | Usage |
|------|-----------|-------|
| `container-sm` | 640px | Modals, forms |
| `container-md` | 768px | Article, long form |
| `container-lg` | 1024px | Dashboard |
| `container-xl` | 1280px | Marketing pages |
| `container-full` | 1440px | App shell |

### 5.4 Touch Targets

Minimum touch target sizes (WCAG 2.5.5):

| Element | Min Size |
|---------|----------|
| Buttons | 44×44px |
| List items | 48×48px |
| Icon buttons | 44×44px (with padding) |
| Form inputs | 48px height |
| Bottom nav | 56px height |

---

## 6. Border Radius & Elevation

### 6.1 Border Radius

| Token | Value | Usage |
|-------|-------|-------|
| `radius-none` | 0 | Full-bleed images |
| `radius-sm` | 4px | Small badges, tags |
| `radius-md` | 8px | Buttons, inputs |
| `radius-lg` | 12px | Cards |
| `radius-xl` | 16px | Modals, sheets |
| `radius-2xl` | 24px | Hero cards |
| `radius-full` | 9999px | Pills, avatars |

### 6.2 Elevation (Shadows)

Material Design 3 inspired, but softer:

| Token | Value | Usage |
|-------|-------|-------|
| `shadow-none` | none | Flat |
| `shadow-sm` | `0 1px 2px rgba(0,0,0,0.05)` | Subtle (cards in list) |
| `shadow-md` | `0 2px 8px rgba(0,0,0,0.08)` | Default cards |
| `shadow-lg` | `0 8px 24px rgba(0,0,0,0.12)` | Modals, dropdowns |
| `shadow-xl` | `0 16px 48px rgba(0,0,0,0.16)` | Overlays |
| `shadow-2xl` | `0 24px 64px rgba(0,0,0,0.20)` | Dialogs |

### 6.3 Border

| Token | Value | Usage |
|-------|-------|-------|
| `border-default` | 1px solid #E0E0E0 | Default borders |
| `border-strong` | 2px solid #BDBDBD | Focused |
| `border-error` | 1px solid #D32F2F | Error states |
| `divider` | 1px solid #EEEEEE | Section dividers |

---

## 7. Component Library

### 7.1 Component Inventory

| Category | Components |
|----------|-----------|
| **Buttons** | Primary, Secondary, Ghost, Destructive, Icon, FAB |
| **Forms** | Input, Textarea, Select, Checkbox, Radio, Toggle, Slider, OTP Input |
| **Navigation** | Bottom Nav, Tabs, Breadcrumb, Pagination, Drawer |
| **Layout** | Card, List, Grid, Stack, Divider, Spacer |
| **Feedback** | Toast, Snackbar, Alert, Banner, Progress, Spinner, Skeleton |
| **Overlays** | Modal, Bottom Sheet, Popover, Tooltip, Dialog |
| **Data Display** | Table, Avatar, Badge, Chip, Tag, Rating, Stat, Timeline |
| **Media** | Image, Avatar, Carousel, Gallery |
| **Restaurant-Specific** | RestaurantCard, MenuItemCard, CartItem, OrderCard |
| **Map** | MapView, Marker, Route, HeatLayer |

### 7.2 Button Specifications

```typescript
// Button variants
type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'destructive' | 'outline'
type ButtonSize = 'sm' | 'md' | 'lg'
type ButtonIcon = 'none' | 'left' | 'right' | 'icon-only'

// Specs
const buttonSpecs = {
  sm: { height: 32, paddingX: 12, fontSize: 13 },
  md: { height: 44, paddingX: 16, fontSize: 14 },
  lg: { height: 56, paddingX: 24, fontSize: 16 },
  radius: 8,
  fontWeight: 600,
  iconSize: { sm: 16, md: 20, lg: 24 },
  transition: '150ms ease',
}

// Variants
const buttonVariants = {
  primary: { bg: '#FF5722', color: '#FFFFFF', hover: '#E64A19' },
  secondary: { bg: '#00897B', color: '#FFFFFF', hover: '#00695C' },
  ghost: { bg: 'transparent', color: '#FF5722', hover: 'rgba(255,87,34,0.08)' },
  destructive: { bg: '#D32F2F', color: '#FFFFFF', hover: '#B71C1C' },
  outline: { bg: 'transparent', color: '#212121', border: '1px solid #BDBDBD' },
}
```

### 7.3 Card Specifications

```typescript
const cardSpecs = {
  padding: 20,
  radius: 12,
  shadow: '0 2px 8px rgba(0,0,0,0.08)',
  hoverShadow: '0 8px 24px rgba(0,0,0,0.12)',
  hoverTransform: 'translateY(-2px)',
  transition: '200ms ease',
}
```

### 7.4 Input Specifications

```typescript
const inputSpecs = {
  height: 48,
  paddingX: 16,
  radius: 8,
  fontSize: 14,
  borderColor: { default: '#E0E0E0', focus: '#FF5722', error: '#D32F2F' },
  borderWidth: 1,
  focusBorderWidth: 2,
  labelColor: '#616161',
  placeholderColor: '#9E9E9E',
  helperColor: '#616161',
  errorColor: '#D32F2F',
}
```

---

## 8. Iconography

### 8.1 Icon Library

Use **Material Symbols** (Google) — variable font, supports RTL flipping:

```html
<link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Rounded:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200" rel="stylesheet">

<span class="material-symbols-rounded">restaurant</span>
```

### 8.2 Icon Sizes

| Token | Size | Usage |
|-------|------|-------|
| `icon-sm` | 16px | Inline, badges |
| `icon-md` | 20px | Buttons, list items |
| `icon-lg` | 24px | Navigation, default |
| `icon-xl` | 32px | Empty states |
| `icon-2xl` | 48px | Hero, onboarding |

### 8.3 Custom Icons

For brand-specific concepts:
- Delivery scooter (custom)
- Egyptian food (koshari, foul — custom illustrations)
- Vodafone Cash logo (third-party brand)
- InstaPay logo (third-party brand)

### 8.4 Icon Rules

- Always use rounded variant (matches our radius)
- Weight 400 (regular), except 600 for active states
- FILL 0 (outline) default, FILL 1 for active
- Ensure RTL-flippable (e.g., arrow_back becomes arrow_forward in RTL)

---

## 9. Imagery

### 9.1 Restaurant Imagery

- **Aspect ratio**: 16:9 (cover), 1:1 (logo)
- **Resolution**: 2x minimum (display density)
- **Format**: WebP (smaller), fallback JPEG
- **Quality**: 80% (balance size/quality)
- **Lazy load**: Always (with skeleton placeholder)
- **Hover**: Subtle zoom (1.05x, 300ms)

### 9.2 Food Item Imagery

- **Aspect ratio**: 1:1 (square) or 4:3
- **Background**: Clean, neutral (white or light gray)
- **Style**: Appetizing, well-lit
- **Fallback**: Generic food icon if no image

### 9.3 Avatar/Profile Images

- **Shape**: Circle
- **Sizes**: 24px (inline), 32px (list), 40px (card), 80px (profile), 120px (hero)
- **Fallback**: Initials on brand-color background
- **Border**: 2px white (for stacking)

### 9.4 Empty State Illustrations

- Style: Flat, friendly, Egyptian-contextual
- Colors: Use brand palette
- Examples: Empty cart, no restaurants, no orders, error

### 9.5 Illustration References

- unDraw (free, customizable)
- Storyset (free with attribution)
- Custom illustrations for hero sections (commission if budget allows)

---

## 10. Motion & Animation

### 10.1 Motion Principles

Inspired by Uber Eats Base design system:

1. **Quick**: 150-250ms for UI feedback
2. **Smooth**: Use `cubic-bezier(0.4, 0, 0.2, 1)` (Material standard)
3. **Purposeful**: Every animation communicates state change
4. **Respectful**: Honor `prefers-reduced-motion`

### 10.2 Easing Curves

```css
--ease-standard: cubic-bezier(0.4, 0, 0.2, 1);    /* Default */
--ease-decelerate: cubic-bezier(0, 0, 0.2, 1);    /* Entering */
--ease-accelerate: cubic-bezier(0.4, 0, 1, 1);    /* Exiting */
--ease-spring: cubic-bezier(0.5, 1.5, 0.5, 1);    /* Playful */
```

### 10.3 Duration Tokens

| Token | Duration | Usage |
|-------|----------|-------|
| `duration-fast` | 100ms | Hover, focus |
| `duration-normal` | 200ms | Default transitions |
| `duration-slow` | 300ms | Page transitions, modals |
| `duration-slower` | 450ms | Complex animations |
| `duration-slowest` | 600ms | Onboarding sequences |

### 10.4 Common Animations

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Button hover | background-color change | 150ms | standard |
| Card hover | translateY(-2px) + shadow | 200ms | standard |
| Modal open | opacity + scale(0.95 → 1) | 200ms | decelerate |
| Modal close | opacity + scale(1 → 0.95) | 150ms | accelerate |
| Bottom sheet | translateY(100% → 0) | 300ms | decelerate |
| Toast enter | translateY(-20px) + opacity | 200ms | decelerate |
| Toast exit | opacity | 150ms | accelerate |
| Skeleton shimmer | gradient sweep | 1500ms | linear (loop) |
| List item enter | opacity + translateY(10px) | 200ms | decelerate |
| Pull to refresh | rotation | 1000ms | linear |

### 10.5 Page Transitions

- **Forward navigation**: Slide in from right (RTL: from left)
- **Back navigation**: Slide out to right (RTL: to left)
- **Modal**: Scale + fade
- **Tab switch**: Cross-fade

### 10.6 Reduced Motion

```css
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 11. States

### 11.1 Empty States

Every list/grid MUST have an empty state. Components:

1. **Illustration** (icon or graphic, 80-120px)
2. **Title** (h3, bold)
3. **Description** (body, secondary color)
4. **CTA** (primary button, optional)

Example (no orders):
```
┌──────────────────────┐
│                      │
│        📦            │
│                      │
│   مفيش طلبات بعد     │
│                      │
│ لما تطلب، هتلاقي     │
│ طلباتك هنا           │
│                      │
│  [تصفح المطاعم]      │
│                      │
└──────────────────────┘
```

### 11.2 Loading States

**Skeleton screens** (preferred over spinners):

```html
<div class="skeleton-card">
  <div class="skeleton-image"></div>  <!-- shimmer -->
  <div class="skeleton-line w-3/4"></div>
  <div class="skeleton-line w-1/2"></div>
</div>
```

- Use shimmer animation (gradient sweep, 1.5s loop)
- Match the layout of the actual content
- Show for 300ms+ (avoid flash on fast loads)

**Spinners** for short actions (<300ms):
- Size: 24px (default), 40px (large)
- Color: Primary brand color
- Use for button loading, form submission

### 11.3 Error States

Three types:

1. **Inline errors** (form validation)
   - Red text below input
   - Red border on input
   - Error icon

2. **Inline error blocks** (API errors)
   - Red-tinted card
   - Error icon + message
   - Retry button

3. **Full-screen errors** (page load failures)
   - Large illustration
   - Error message
   - "Try again" button
   - "Contact support" link

### 11.4 Success States

- **Toast notifications** (default): Green-tinted, auto-dismiss 3s
- **Inline confirmation**: Checkmark + message, fades after 2s
- **Full-screen success**: For major actions (order placed) — illustration + confirmation + CTA

---

## 12. Arabic RTL Considerations

### 12.1 RTL by Default

All Arabic screens use `dir="rtl"`. Tailwind config:

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      // Logical properties (Tailwind 3.3+)
    }
  },
  // Enable RTL plugin
  plugins: [
    require('tailwindcss-rtl'),
  ],
}
```

### 12.2 RTL-Specific Adjustments

| Element | LTR | RTL |
|---------|-----|-----|
| Text direction | left-to-right | right-to-left |
| Icons (arrows) | → points right | → points left (flip) |
| Margins (ml-) | margin-left | margin-right (use ms-/me-) |
| Padding (pl-) | padding-left | padding-right (use ps-/pe-) |
| Float | left | right |
| Text align | left | right (default) |

### 12.3 Use Logical Properties

```css
/* ❌ BAD — physical */
margin-left: 16px;
padding-right: 8px;

/* ✅ GOOD — logical */
margin-inline-start: 16px;  /* mis-4 in Tailwind */
padding-inline-end: 8px;    /* pie-2 in Tailwind */
```

### 12.4 Numbers in RTL

Numbers stay LTR even in RTL context:

```html
<span dir="rtl">
  السعر: <span dir="ltr">EGP 145</span>
</span>
```

### 12.5 Mixed Content

For mixed Arabic/English (e.g., brand names):

```html
<span>اطلب من <span dir="ltr" class="font-en">McDonald's</span> دلوقتي</span>
```

### 12.6 Font Considerations

- Arabic: Cairo (designed for screens, supports Egyptian Arabic)
- English: Inter (geometric, pairs well)
- Switch automatically based on character detection
- Numbers: JetBrains Mono (tabular figures)

### 12.7 Date/Time Formatting

Use Arabic month names:

```javascript
const formatter = new Intl.DateTimeFormat('ar-EG', {
  day: 'numeric',
  month: 'long',
  year: 'numeric',
})
// "4 يوليو 2026"
```

### 12.8 Currency Formatting

```javascript
const formatter = new Intl.NumberFormat('ar-EG', {
  style: 'currency',
  currency: 'EGP',
  currencyDisplay: 'code',
})
// "EGP 145.00"
```

---

## 13. Responsive Breakpoints

### 13.1 Breakpoints

| Name | Min Width | Usage |
|------|-----------|-------|
| `xs` | 0 | Mobile portrait (default) |
| `sm` | 640px | Mobile landscape, small tablet |
| `md` | 768px | Tablet portrait |
| `lg` | 1024px | Tablet landscape, small desktop |
| `xl` | 1280px | Desktop |
| `2xl` | 1536px | Large desktop |

### 13.2 Mobile-First Approach

All base styles target mobile. Use `sm:`, `md:`, etc. for larger screens.

```css
/* Base = mobile */
.card { padding: 16px; }

/* Tablet+ */
@screen md {
  .card { padding: 24px; }
}

/* Desktop+ */
@screen lg {
  .card { padding: 32px; }
}
```

### 13.3 Layout Patterns per Breakpoint

| Breakpoint | Layout |
|------------|--------|
| Mobile (<640px) | Single column, bottom nav, stacked cards |
| Tablet (640-1024px) | 2-column grids, side nav option |
| Desktop (1024px+) | Multi-column, sidebar nav, dense data |

### 13.4 App Shell per Breakpoint

**Mobile**:
```
┌─────────────────┐
│ Top Bar (56px)  │
├─────────────────┤
│                 │
│   Content       │
│                 │
├─────────────────┤
│ Bottom Nav (64) │
└─────────────────┘
```

**Desktop**:
```
┌────┬────────────┐
│ S  │ Top Bar    │
│ i  ├────────────┤
│ d  │            │
│ e  │  Content   │
│ b  │            │
│ a  │            │
│ r  │            │
└────┴────────────┘
```

---

## 14. Accessibility

### 14.1 WCAG 2.1 AA Compliance

Mandatory for all screens:

- **Color contrast**: 4.5:1 for normal text, 3:1 for large text
- **Touch targets**: 44×44px minimum
- **Focus indicators**: Visible focus ring (2px solid primary)
- **Keyboard nav**: All interactive elements reachable via Tab
- **Screen reader**: ARIA labels in Arabic, semantic HTML
- **Color independence**: Information not conveyed by color alone

### 14.2 ARIA Patterns

```html
<!-- Button -->
<button aria-label="إضافة للسلة">أضف</button>

<!-- Loading -->
<div role="status" aria-live="polite">
  <span class="sr-only">جاري التحميل...</span>
  <Spinner />
</div>

<!-- Toast -->
<div role="alert" aria-live="assertive">
  تم إنشاء طلبك بنجاح
</div>

<!-- Modal -->
<div role="dialog" aria-modal="true" aria-labelledby="title">
  <h2 id="title">تأكيد الطلب</h2>
</div>
```

### 14.3 Screen Reader Considerations

- Arabic ARIA labels (not English)
- Announce dynamic changes (order status, cart updates)
- Hide decorative icons (`aria-hidden="true"`)
- Provide text alternatives for images

### 14.4 Keyboard Navigation

- Tab order follows visual order
- Skip-to-content link on every page
- Escape closes modals/dropdowns
- Enter/Space activates buttons
- Arrow keys for lists/tabs

---

## 15. Competitive Analysis

### 15.1 Feature Comparison Matrix

| Feature | Uber Eats | Talabat | elmenus | DoorDash | **Ours** |
|---------|-----------|---------|---------|----------|----------|
| **Visual Style** | Clean, image-heavy | Dense, info-rich | Friendly, food-focused | Bold, fast | Warm, Egyptian, polished |
| **Primary Color** | Green (#06C167) | Orange (#FF5A00) | Orange (#FF4F00) | Red (#FF3008) | Orange-Red (#FF5722) |
| **Typography** | Uber Move | Cairo | Cairo | DDOak | Cairo + Inter |
| **Card Style** | Large images, minimal text | Compact, dense info | Medium, food imagery | Bold, photo-forward | Medium-large, balanced |
| **Bottom Nav** | 4 tabs (Home, Browse, Activity, Account) | 5 tabs (Home, Search, Offers, Orders, Account) | 4 tabs | 4 tabs | 4 tabs (Home, Search, Orders, Profile) |
| **Search** | Top bar, prominent | Top bar, prominent | Hero search | Top bar | Top bar + hero on home |
| **Address Picker** | Top bar, tap to change | Top bar, tap to change | Optional | Top bar | Top bar, prominent |
| **Restaurant Card** | Image, name, cuisine, ETA, delivery fee | Image, name, rating, cuisine, ETA, promo | Image, name, rating, cuisine | Image, name, ETA, fee | Image, name, rating, cuisine, ETA, promo |
| **Menu Layout** | Categories sidebar, items list | Categories tabs, items list | Categories sidebar, items list | Categories sidebar, items list | Categories sidebar, items list |
| **Cart** | Bottom sheet | Bottom bar with total | Side panel (desktop) | Bottom sheet | Bottom sheet + dedicated page |
| **Checkout** | Single page, payment selection | Multi-step (Talabat Pay issues) | Single page | Single page | Single page, 3 payment methods |
| **Order Tracking** | 5-step bar, live map, ETA | Status updates, map | Basic status | 5-step, live map | 5-step bar, live map, ETA range |
| **Rating** | Star + thumbs per item | Star + text | Star + text | Star + thumbs | Star + text + photos |
| **Loyalty** | Uber Rewards (points) | Talabat Pro (subscription) | None | DashPass (subscription) | Tiered (Silver/Gold/Platinum) + cashback |
| **Push Notifications** | Order updates, promos | Order updates, promos | Order updates | Order updates, promos | Order updates, smart re-engagement |
| **Arabic Support** | Yes (RTL) | Yes (RTL, primary) | Yes (RTL, primary) | Limited | Yes (RTL, primary) |
| **Payment Methods** | Card, PayPal, Cash | Card, Cash, VF Cash | Card, Cash | Card, Cash | Vodafone Cash, InstaPay, Card, COD |
| **Empty States** | Friendly illustrations | Basic text | Basic text | Friendly illustrations | Custom Egyptian illustrations |
| **Loading States** | Skeleton screens | Spinners | Spinners | Skeleton screens | Skeleton screens |
| **Error States** | Friendly with retry | Basic | Basic | Friendly with retry | Friendly with retry + support link |

### 15.2 What We Adopt from Each

**From Uber Eats**:
- ✅ Clean, image-forward design
- ✅ Skeleton loading states
- ✅ 5-step order tracking bar
- ✅ Minimal text philosophy
- ✅ Bottom sheet for cart

**From Talabat**:
- ✅ Arabic-first, RTL by default
- ✅ Cairo font
- ✅ Dense info in cards (Egyptian users expect this)
- ✅ Vodafone Cash as default payment
- ✅ Local restaurant focus

**From elmenus**:
- ✅ Egyptian cultural context in copy
- ✅ Restaurant discovery focus
- ✅ Photo quality emphasis

**From DoorDash**:
- ✅ Bold CTAs
- ✅ DashPass-style subscription (planned Phase 4)
- ✅ Item-level rating (thumbs up/down)

### 15.3 What We Do Differently

1. **Lower commission** → visible "why us" badge on restaurants
2. **Driver earnings transparency** → unique to us
3. **Trust scores** → visible on every restaurant
4. **Egyptian payment methods** → InstaPay native (Talabat lacks this)
5. **Anti-fraud visibility** → customer sees their own trust score
6. **Field supervisor verification** → "Verified" badge on restaurants

### 15.4 Screens We Must Match Quality Of

Priority screens where we MUST match Uber Eats/Talabat quality:

1. **Home page** (first impression)
2. **Restaurant detail** (conversion critical)
3. **Cart + checkout** (money moment)
4. **Order tracking** (satisfaction moment)
5. **Order history** (retention)
6. **Profile** (account management)

Secondary screens (can be simpler):
- Search results
- Reviews
- Settings
- Help/Support

---

## Appendix A: Tailwind Config

```javascript
// tailwind.config.js
module.exports = {
  content: ['./src/**/*.{ts,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#FF5722',
          dark: '#E64A19',
          light: '#FFAB91',
        },
        secondary: {
          DEFAULT: '#00897B',
          dark: '#00695C',
          light: '#B2DFDB',
        },
        success: { DEFAULT: '#2E7D32', light: '#C8E6C9' },
        warning: { DEFAULT: '#F9A825', light: '#FFF9C4' },
        error: { DEFAULT: '#D32F2F', light: '#FFCDD2' },
        info: { DEFAULT: '#0288D1', light: '#B3E5FC' },
        // Neutrals
        bg: { primary: '#FAFAFA', secondary: '#F5F5F5', tertiary: '#EEEEEE' },
        text: { primary: '#212121', secondary: '#616161', tertiary: '#9E9E9E' },
        border: { DEFAULT: '#E0E0E0', strong: '#BDBDBD' },
      },
      fontFamily: {
        sans: ['Cairo', 'Inter', 'sans-serif'],
        en: ['Inter', 'sans-serif'],
        ar: ['Cairo', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      fontSize: {
        'display-lg': ['48px', { lineHeight: '56px', fontWeight: '800' }],
        'display-md': ['40px', { lineHeight: '48px', fontWeight: '800' }],
        'display-sm': ['32px', { lineHeight: '40px', fontWeight: '700' }],
        h1: ['28px', { lineHeight: '36px', fontWeight: '700' }],
        h2: ['24px', { lineHeight: '32px', fontWeight: '700' }],
        h3: ['20px', { lineHeight: '28px', fontWeight: '600' }],
        h4: ['18px', { lineHeight: '24px', fontWeight: '600' }],
        'body-lg': ['16px', { lineHeight: '24px' }],
        body: ['14px', { lineHeight: '20px' }],
        'body-sm': ['13px', { lineHeight: '18px' }],
        caption: ['12px', { lineHeight: '16px', fontWeight: '500' }],
        overline: ['11px', { lineHeight: '16px', fontWeight: '600' }],
        'mono-lg': ['16px', { lineHeight: '24px', fontWeight: '600' }],
        mono: ['14px', { lineHeight: '20px', fontWeight: '500' }],
        'mono-sm': ['12px', { lineHeight: '16px', fontWeight: '500' }],
      },
      spacing: {
        1: '4px', 2: '8px', 3: '12px', 4: '16px',
        5: '20px', 6: '24px', 8: '32px', 10: '40px',
        12: '48px', 16: '64px',
      },
      borderRadius: {
        sm: '4px', md: '8px', lg: '12px',
        xl: '16px', '2xl': '24px', full: '9999px',
      },
      boxShadow: {
        sm: '0 1px 2px rgba(0,0,0,0.05)',
        md: '0 2px 8px rgba(0,0,0,0.08)',
        lg: '0 8px 24px rgba(0,0,0,0.12)',
        xl: '0 16px 48px rgba(0,0,0,0.16)',
        '2xl': '0 24px 64px rgba(0,0,0,0.20)',
      },
      transitionDuration: {
        fast: '100ms', normal: '200ms', slow: '300ms',
        slower: '450ms', slowest: '600ms',
      },
      transitionTimingFunction: {
        standard: 'cubic-bezier(0.4, 0, 0.2, 1)',
        decelerate: 'cubic-bezier(0, 0, 0.2, 1)',
        accelerate: 'cubic-bezier(0.4, 0, 1, 1)',
        spring: 'cubic-bezier(0.5, 1.5, 0.5, 1)',
      },
    },
  },
  plugins: [
    require('tailwindcss-rtl'),
  ],
}
```

---

## Appendix B: Design Tokens (TypeScript)

```typescript
// packages/theme/src/tokens.ts

export const tokens = {
  colors: {
    primary: '#FF5722',
    primaryDark: '#E64A19',
    primaryLight: '#FFAB91',
    secondary: '#00897B',
    secondaryDark: '#00695C',
    secondaryLight: '#B2DFDB',
    success: '#2E7D32',
    warning: '#F9A825',
    error: '#D32F2F',
    info: '#0288D1',
    bgPrimary: '#FAFAFA',
    bgSecondary: '#F5F5F5',
    bgTertiary: '#EEEEEE',
    surface: '#FFFFFF',
    border: '#E0E0E0',
    borderStrong: '#BDBDBD',
    textPrimary: '#212121',
    textSecondary: '#616161',
    textTertiary: '#9E9E9E',
  },
  darkColors: {
    bgPrimary: '#0A0E1A',
    bgSecondary: '#131A2E',
    surface: '#1A2240',
    border: '#2A3656',
    textPrimary: '#EAF0FF',
    textSecondary: '#9BA8C7',
    accent: '#00D4FF',
    accent2: '#B14EFF',
  },
  typography: {
    fontFamily: {
      sans: 'Cairo, Inter, sans-serif',
      en: 'Inter, sans-serif',
      ar: 'Cairo, sans-serif',
      mono: 'JetBrains Mono, monospace',
    },
    fontSize: {
      displayLg: '48px',
      displayMd: '40px',
      displaySm: '32px',
      h1: '28px',
      h2: '24px',
      h3: '20px',
      h4: '18px',
      bodyLg: '16px',
      body: '14px',
      bodySm: '13px',
      caption: '12px',
      overline: '11px',
    },
  },
  spacing: {
    1: '4px', 2: '8px', 3: '12px', 4: '16px',
    5: '20px', 6: '24px', 8: '32px', 10: '40px',
    12: '48px', 16: '64px',
  },
  borderRadius: {
    sm: '4px', md: '8px', lg: '12px',
    xl: '16px', '2xl': '24px', full: '9999px',
  },
  shadows: {
    sm: '0 1px 2px rgba(0,0,0,0.05)',
    md: '0 2px 8px rgba(0,0,0,0.08)',
    lg: '0 8px 24px rgba(0,0,0,0.12)',
    xl: '0 16px 48px rgba(0,0,0,0.16)',
  },
  animation: {
    durationFast: '100ms',
    durationNormal: '200ms',
    durationSlow: '300ms',
    easeStandard: 'cubic-bezier(0.4, 0, 0.2, 1)',
    easeDecelerate: 'cubic-bezier(0, 0, 0.2, 1)',
    easeAccelerate: 'cubic-bezier(0.4, 0, 1, 1)',
    easeSpring: 'cubic-bezier(0.5, 1.5, 0.5, 1)',
  },
  breakpoints: {
    sm: '640px',
    md: '768px',
    lg: '1024px',
    xl: '1280px',
    '2xl': '1536px',
  },
} as const
```

---

> **Next**: Read `SCREEN-INVENTORY.md` for the full list of screens, then per-app `UI-SPEC.md` files.
