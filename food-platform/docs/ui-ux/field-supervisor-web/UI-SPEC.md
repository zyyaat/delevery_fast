# Field Supervisor Web App — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: Salesforce Field Service, GoAudits, SafetyCulture research

---

## 1. App Overview

### Purpose

Field Supervisor Web App enables on-ground supervisors to:
- Receive daily tasks (restaurant verification, driver verification, audits, complaints)
- Conduct on-site verification with photo + GPS proof
- Submit detailed reports
- Plan optimal routes
- Investigate customer complaints on-site

### Design Inspiration

| Source | What We Take |
|--------|--------------|
| Salesforce Field Service | Task management, route planning, GPS tracking |
| GoAudits | Inspection checklists, photo capture |
| SafetyCulture | Mobile-first audit, compliance scoring |

### Design Principles

1. **On-Site First**: Designed for use while standing/walking
2. **Photo-Heavy**: Every verification requires photos with GPS
3. **Offline-Tolerant**: Works in low-signal areas
4. **Battery-Aware**: GPS only when needed
5. **Quick Capture**: Photo + checklist in <5 min per visit

---

## 2. Screen Specifications

### FLD-TASK-01: Task List ⭐

```
┌─────────────────────────────────────────────────────┐
│  👤 Omar F. | Field Supervisor                      │
│  🟢 Online | Today: 4 يوليو 2026                    │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────┐ ┌──────────┐ ┌────────┐              │
│  │ 📋 مهام  │ │ ✅ مكتملة│ │ ⏱️ مسافة│              │
│  │  اليوم   │ │  اليوم   │ │        │              │
│  │   6/8    │ │    6     │ │ 45km   │              │
│  └──────────┘ └──────────┘ └────────┘              │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  📍 موقعك الحالي: معادي                              │
│  أقرب مهمة: 1.2km                                   │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  المهام القادمة:                                     │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ ⚠️ Complaint Investigation (URGENT)          │   │
│  │ KFC - التجمع                                  │   │
│  │ المسافة: 12km | ETA: 35 دقيقة                 │   │
│  │ Priority: 🔴 URGENT                          │   │
│  │ [▶️ ابدأ فوراً]                                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🍔 Restaurant Verification                   │   │
│  │ Pizza Hut - المعصرة                           │   │
│  │ المسافة: 8.5km | ETA: 25 دقيقة                │   │
│  │ Priority: 🟡 Normal                          │   │
│  │ [▶️ ابدأ]                                     │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🛵 Driver Verification                       │   │
│  │ Mahmoud S. - مدينة نصر                        │   │
│  │ المسافة: 5.2km | ETA: 18 دقيقة                │   │
│  │ Priority: 🟢 Low                             │   │
│  │ [▶️ ابدأ]                                     │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  [📍 الموقع على الخريطة]                             │
│  [📊 تقرير اليوم]                                    │
│  [⚙️ الإعدادات]                                      │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### FLD-TASK-02: Restaurant Verification ⭐ (50+ point checklist)

```
┌─────────────────────────────────────────────────────┐
│  ←  Restaurant Verification                         │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🍔 Pizza Hut - المعصرة                        │   │
│  │                                              │   │
│  │ 📍 15 شارع المعصرة الرئيسي                    │   │
│  │ 📞 02-1234-5678                              │   │
│  │ 👤 Owner: Mohamed A.                         │   │
│  │ 📅 Applied: 3 يوليو 2026                     │   │
│  │ ⏰ موعد الزيارة: 11:00 AM                    │   │
│  │                                              │   │
│  │ [📍 افتح الخريطة]                              │   │
│  │ [📞 اتصال بالمطعم]                             │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  حالة الوصول:                                        │
│  [🟢 وصلت للموقع]                                    │
│  (سيتم تسجيل GPS تلقائياً)                            │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  الـ Checklist (50 نقطة)                             │
│                                                     │
│  A. هوية المطعم (5 نقاط)                              │
│  □ ✅ اسم المطعم مطابق                                │
│  □ ✅ صاحب المطعم حاضر                                │
│  □ ⚠️ رخصة التشغيل منتهية                            │
│  □ ✅ البطاقة الضريبية صحيحة                         │
│  □ ✅ عقد الإيجار موجود                                │
│                                                     │
│  B. الموقع (4 نقاط)                                  │
│  □ ✅ العنوان مطابق                                   │
│  □ ✅ GPS مطابق                                       │
│  □ ✅ المدخل سهل                                      │
│  □ ❌ مكان الانتظار غير متاح                           │
│                                                     │
│  [عرض باقي الأقسام ▼]                                 │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  الصور المطلوبة (7 صور):                              │
│  ┌──────────────────────────────────────────────┐   │
│  │ ✅ واجهة المطعم           [📷]                 │   │
│  │ ✅ المدخل                [📷]                 │   │
│  │ ✅ صالة الطعام           [📷]                 │   │
│  │ ✅ المطبخ                [📷]                 │   │
│  │ ⚠️ معدات التخزين (إعادة)  [📷]                 │   │
│  │ ✅ رخصة التشغيل          [📷]                 │   │
│  │ ✅ البطاقة الضريبية      [📷]                 │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ملاحظات المشرف:                                      │
│  [_________________________]                        │
│  [_________________________]                        │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  القرار النهائي:                                      │
│  ○ ✅ Approved (تفعيل فوري)                          │
│  ● ⚠️ Conditional (مهلة 7 أيام)                      │
│  ○ ❌ Rejected (رفض)                                 │
│                                                     │
│  سبب القرار (إجباري):                                 │
│  [_________________________]                        │
│                                                     │
│  [🔒 تأكيد القرار - يتطلب بصمة]                       │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### FLD-TASK-03: Driver Verification

```
┌─────────────────────────────────────────────────────┐
│  ←  Driver Verification                             │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🛵 Mahmoud S.                                 │   │
│  │ 📱 01012345678                                │   │
│  │ 🛵 موتوسيكل - لوحة 1234                       │   │
│  │ 📅 Applied: 2 يوليو 2026                     │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  Checklist:                                          │
│                                                     │
│  A. هوية المندوب                                      │
│  □ ✅ البطاقة الشخصية أصلية                           │
│  □ ✅ صورة المندوب مطابقة                              │
│  □ ✅ Face match (liveness)                         │
│  □ ✅ الرقم القومي موجود في القاعدة                    │
│                                                     │
│  B. المركبة                                           │
│  □ ✅ المركبة موجودة فعلياً                            │
│  □ ✅ اللوحات مطابقة                                  │
│  □ ⚠️ التأمين ساري (ينتهي شهر 8)                    │
│  □ ✅ المركبة في حالة جيدة                            │
│  □ ✅ خوذة متوفرة                                     │
│                                                     │
│  C. الفحص العملي                                      │
│  □ ✅ يقدر يركب الموتوسيكل بأمان                      │
│  □ ✅ يعرف إشارات المرور                             │
│  □ ✅ يقدر يستخدم GPS                                │
│  □ ✅ يقدر يتواصل بأدب (اختبار قصير)                 │
│                                                     │
│  الصور المطلوبة:                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │ ✅ صورة المندوب مع المركبة   [📷]              │   │
│  │ ✅ صورة اللوحات              [📷]              │   │
│  │ ✅ صورة الرخصة              [📷]              │   │
│  │ ✅ صورة البطاقة             [📷]              │   │
│  │ ✅ selfie مع البطاقة        [📷]              │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  القرار:                                              │
│  ● ✅ Approved (تفعيل)                               │
│  ○ ❌ Rejected (رفض)                                 │
│                                                     │
│  [🔒 تأكيد القرار]                                    │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### FLD-ROUTE-01: Route Planner

```
┌─────────────────────────────────────────────────────┐
│  🗺️ تخطيط المسار اليوم                              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  مهام اليوم (8):                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 1. 09:00 🍔 Pizza Hut معادي                  │   │
│  │      ← موقعك الحالي (2km)                    │   │
│  │                                              │   │
│  │ 2. 10:30 🛵 Driver Verify                    │   │
│  │      مدينة نصر (5km)                         │   │
│  │                                              │   │
│  │ 3. 12:00 🍔 KFC التجمع                        │   │
│  │      (12km)                                  │   │
│  │                                              │   │
│  │ 4. 13:30 ⚠️ Complaint                         │   │
│  │      الزمالك (8km)                           │   │
│  │                                              │   │
│  │ 5. 15:00 🍔 McDonald's                       │   │
│  │      مصر الجديدة (15km)                      │   │
│  │                                              │   │
│  │ ... 3 مهام إضافية                              │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  إجمالي المسافة: 68km                                │
│  إجمالي الوقت: 6 ساعات                                │
│                                                     │
│  [🔄 إعادة ترتيب تلقائي (أمثل)]                       │
│  [📍 ابدأ بالزيارة الأولى]                             │
│                                                     │
│  خيارات التخطيط:                                      │
│  ☑ تجميع المهام القريبة                                │
│  ☑ تجنب ساعات الذروة                                  │
│  ☑ إعطاء أولوية للمهام العاجلة                         │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### FLD-REPT-01: Daily Report

```
┌─────────────────────────────────────────────────────┐
│  📊 تقرير اليوم - Omar F.                           │
│  4 يوليو 2026                                       │
├─────────────────────────────────────────────────────┤
│                                                     │
│  المهام:                                              │
│  • مكتملة: 7/8 (87%)                                 │
│  • ملغية: 1 (مطعم مقفل)                               │
│  • مدة العمل: 7 ساعات                                  │
│  • المسافة: 68km                                      │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  تفصيل الزيارات:                                      │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🍔 Restaurant Verifications                  │   │
│  │ ✅ Approved: 2                               │   │
│  │ ⚠️ Conditional: 1                            │   │
│  │ ❌ Rejected: 0                               │   │
│  │                                              │   │
│  │ 🛵 Driver Verifications                      │   │
│  │ ✅ Approved: 3                               │   │
│  │ ❌ Rejected: 1 (مستندات وهمية)                │   │
│  │                                              │   │
│  │ ⚠️ Complaint Investigation                   │   │
│  │ 🔴 Resolved: 1                               │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Issues Found:                                       │
│  • 1 رخصة تشغيل منتهية                                │
│  • 1 مركبة في حالة سيئة                               │
│  • 1 مستند وهمي (محول للتحقيق)                        │
│  • 1 مطعم فيه مشكلة نظافة                             │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  التوصيات:                                            │
│  • إعادة فحص Pizza Hut بعد 7 أيام                      │
│  • تفعيل فوري لـ 2 مطعم                                │
│  • suspend المندوب بمستندات وهمية                      │
│                                                     │
│  [📤 رفع التقرير النهائي]                              │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 3. UX Flows

### Restaurant Verification Flow

```
1. Supervisor arrives at restaurant
   ↓
2. Tap "🟢 وصلت للموقع" (GPS recorded)
   ↓
3. Go through 50-point checklist:
   - Identity (5 points)
   - Location (4 points)
   - Hygiene (15 points)
   - Operations (10 points)
   - Safety (8 points)
   - Pricing (3 points)
   - Photos (7 required)
   ↓
4. Take photos with GPS metadata
   ↓
5. Add notes
   ↓
6. Decision: Approve / Conditional / Reject
   ↓
7. Biometric confirmation (WebAuthn)
   ↓
8. Report submitted
   ↓
9. If approved: restaurant activated
   ↓
10. If conditional: 7-day follow-up scheduled
    ↓
11. If rejected: applicant notified
```

### Complaint Investigation Flow

```
1. Receive complaint (P0/P1)
   ↓
2. Navigate to location
   ↓
3. Investigate:
   - Meet customer (if willing)
   - Visit restaurant
   - Speak with driver (if applicable)
   - Take photos
   ↓
4. Submit investigation report:
   - Findings
   - Photos
   - Witness statements
   - Recommendation
   ↓
5. Resolution:
   - Restaurant at fault: suspend + investigate
   - Driver at fault: suspend + investigate
   - Customer fraud: flag account
   ↓
6. Customer compensated (if applicable)
```

---

## 4. Critical UX Considerations

### 4.1 GPS Verification

- Every photo includes GPS metadata
- GPS must be within 50m of restaurant location
- Geofence triggers task start/stop
- If GPS mismatch: flag + review

### 4.2 Photo Capture

- Use device camera (not file picker)
- Capture EXIF data (GPS, timestamp, device)
- Compress before upload (WebP)
- Watermark with timestamp + GPS

### 4.3 Offline Support

- Cache tasks for the day locally
- Allow checklist completion offline
- Queue photos for upload
- Sync when connection restored

### 4.4 Anti-Cheating

- GPS must match task location
- Photos must be fresh (timestamp within 1hr)
- Time on site minimum 15 min
- Spot checks (10% re-verified by another supervisor)

### 4.5 Battery Optimization

- GPS polling: 30s when navigating, 5min when idle
- Photo compression before upload
- Dark mode for OLED screens
- Battery warning at 20%

---

## 5. Component Library

### App-Specific Components

#### TaskCard
```typescript
interface TaskCardProps {
  task: FieldTask
  onStart: () => void
}
```

#### Checklist
```typescript
interface ChecklistProps {
  sections: ChecklistSection[]
  onChange: (section, item, value) => void
}
```

#### PhotoCapture
```typescript
interface PhotoCaptureProps {
  type: 'storefront' | 'interior' | 'kitchen' | 'license' | 'vehicle' | 'document'
  onCapture: (photo: PhotoWithMetadata) => void
  required: boolean
}
```

#### RouteOptimizer
```typescript
interface RouteOptimizerProps {
  tasks: FieldTask[]
  currentLocation: Coordinates
  onOptimize: (orderedTasks: FieldTask[]) => void
}
```

#### GPSVerification
```typescript
interface GPSVerificationProps {
  targetLocation: Coordinates
  onArrive: () => void
  radius: 50  // meters
}
```

---

> Field Supervisor Web App is mobile-first optimized (designed for use on phone/tablet in the field).
