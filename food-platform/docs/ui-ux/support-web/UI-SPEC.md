# Support Web App — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: Zendesk, Freshdesk, Intercom research

---

## 1. App Overview

### Purpose

Support Web App enables customer support agents to:
- Handle customer/driver/restaurant tickets via omnichannel (chat, email, phone, WhatsApp)
- Process refunds (with biometric verification for >EGP 100)
- Cancel/modify orders
- Escalate to Tier 2 or Ops
- Access knowledge base for quick resolutions

### Design Inspiration

| Source | What We Take |
|--------|--------------|
| Zendesk | Ticket queue, agent dashboard, macros |
| Freshdesk | SLA management, omnichannel |
| Intercom | Modern chat interface, customer 360 |

### Design Principles

1. **Fast First Response**: <30s target for chat
2. **Context-Aware**: All customer info visible during chat
3. **Action-Oriented**: Quick actions inline (refund, cancel)
4. **Quality-Focused**: Sentiment tracking, QA scoring

---

## 2. Screen Specifications (Key Screens)

### SUP-DASH-01: Main Dashboard

```
┌─────────────────────────────────────────────────────┐
│ 👤 Agent: Sarah M. | Tier 1.5 | Online             │
│ 📊 اليوم: 47 تذكرة | 92% CSAT | Avg: 2:34          │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │
│  │ 4 نشطة   │ │ 7 في     │ │ 12       │ │ 3      │ │
│  │ الآن     │ │ الانتظار │ │ اليوم    │ │ معلّقة │ │
│  └──────────┘ └──────────┘ └──────────┘ └────────┘ │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  التذاكر النشطة:                                     │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🔴 #T-8294 | P0 | Ahmed M. | 🥇 Platinum    │   │
│  │ "الطلب اتأخر 45 دقيقة!"                       │   │
│  │ آخر رسالة: 30 ثانية | أنت تكتب...            │   │
│  │ [فتح المحادثة]                                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🟡 #T-8291 | P1 | Mahmoud K. | 🥈 Gold       │   │
│  │ "ناقص كوكا في الطلب"                          │   │
│  │ آخر رسالة: 1:30 دقيقة | بانتظار ردك          │   │
│  │ [فتح المحادثة]                                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### SUP-TICK-02: Ticket Detail (Chat View) ⭐

```
┌─────────────────────────────────────────────────────┐
│  ←  #T-8294 | Ahmed M. | 🥇 Platinum              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 📋 معلومات العميل                              │   │
│  │ Ahmed Mohamed | 01012345678                  │   │
│  │ 🥇 Platinum | 47 طلب | EGP 8,520 إنفاق       │   │
│  │ Trust Score: 92/100                          │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 📦 الطلب الحالي                                │   │
│  │ #A7X92F | Pizza Hut                          │   │
│  │ EGP 636.6 | Vodafone Cash                    │   │
│  │ الحالة: يتحضّر (المطعم)                       │   │
│  │ الوقت: 45 دقيقة (متأخر!)                      │   │
│  │ [👁️ عرض الطلب] [📞 المطعم] [🛵 المندوب]       │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 💬 المحادثة                                    │   │
│  │                                                │   │
│  │ [12:34] Ahmed: الطلب اتأخر 45 دقيقة!           │   │
│  │          🔴 sentiment: غاضب جداً               │   │
│  │                                                │   │
│  │ [12:35] You: أهلاً أستاذ أحمد، أنا سارة...     │   │
│  │                                                │   │
│  │ [12:36] Ahmed: مش محتاجه، ألغيه                │   │
│  │                                                │   │
│  │ [typing...] Ahmed is typing                   │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ ⚡ Quick Actions                               │   │
│  │ [💰 Refund] [❌ Cancel Order] [📞 Call Cust.]  │   │
│  │ [🎁 Promo Code] [⏰ Extend ETA] [🆘 Escalate]  │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 📝 Macros (ردود جاهزة)                         │   │
│  │ [اعتذار عن التأخير] [شرح سياسة الاسترجاع]      │   │
│  │ [طلب صور للمشكلة] [تأكيد الحل]                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ ✏️ اكتب ردك...                                │   │
│  │ [_________________________]                  │   │
│  │ [😊 Emoji] [📎 Attach] [🎤 Voice]   [إرسال]   │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### SUP-TICK-03: Refund Dialog (with Biometric)

```
┌─────────────────────────────────────────────────────┐
│                                                     │
│   💰 معالجة استرجاع                                 │
│                                                     │
│   ┌──────────────────────────────────────────────┐  │
│   │ الطلب: #A7X92F                                 │  │
│   │ العميل: Ahmed Mohamed                          │  │
│   │ المبلغ الأصلي: EGP 636.6                       │  │
│   └──────────────────────────────────────────────┘  │
│                                                     │
│   نوع الاسترجاع:                                      │
│   ● كامل (100%)                                      │
│   ○ جزئي                                             │
│     المبلغ: [300] EGP                                │
│   ○ كوبون اعتذار                                     │
│                                                     │
│   السبب:                                              │
│   [الطلب اتأخر ▾]                                    │
│                                                     │
│   ⚠️ بما إن المبلغ >EGP 100،                            │
│   لازم تأكيد ببصمتك                                   │
│                                                     │
│   [🔑 أكد ببصمتك]                                    │
│                                                     │
│   [إلغاء]                                           │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 3. UX Flows

### Refund Flow

```
1. Agent taps "💰 Refund" in ticket
   ↓
2. Refund dialog opens
   ↓
3. Select refund type (full/partial/coupon)
   ↓
4. If amount >EGP 100: WebAuthn biometric prompt
   ↓
5. If amount >EGP 500: dual approval needed
   ↓
6. Confirm
   ↓
7. Process refund
   ↓
8. Customer notified
```

---

> Support Web App uses the same dark/light theme as other admin apps. Bottom nav replaced by sidebar.
