# Employee Portal — UI/UX Specification

> **Version**: 1.0  
> **Status**: Active  
> **Last Updated**: 2026-07-04  
> **Based on**: WebAuthn/FIDO2 standards, enterprise security patterns

---

## 1. App Overview

### Purpose

Employee Portal enables internal staff to:
- Login with strong auth (password + TOTP + WebAuthn biometric)
- Process sensitive actions (refunds, approvals) with biometric verification
- View immutable audit log (hash-chained)
- Manage access (quarterly review)
- View anomaly alerts (UEBA)

### Design Inspiration

| Source | What We Take |
|--------|--------------|
| Auth0 | WebAuthn implementation patterns |
| Splunk UEBA | Anomaly detection UI |
| HashiCorp Vault | Secrets management UI |
| Enterprise admin tools | Clean, professional, security-focused |

### Design Principles

1. **Security Over Convenience**: Biometric for every sensitive action
2. **Audit Everything**: Every action logged immutably
3. **Least Privilege**: Only show what user has permission for
4. **Dual Approval**: High-value actions need 2 people
5. **Tamper-Evident**: Hash chain verifies integrity nightly

---

## 2. Screen Specifications

### EMP-AUTH-01: Login (Password + TOTP)

```
┌─────────────────────────────────────────────────────┐
│                                                     │
│           🔐 Employee Portal                        │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │  Username                                     │   │
│  │  [__________________________]                 │   │
│  │                                                │   │
│  │  Password                                     │   │
│  │  [__________________________]                 │   │
│  │                                                │   │
│  │  TOTP Code (6 digits)                         │   │
│  │  [______]                                     │   │
│  │                                                │   │
│  │  [🔐 Login with Biometric]                    │   │
│  │                                                │   │
│  │  Forgot password? | Need help?                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ⚠️ جميع الإجراءات مسجلة ومراقبة                      │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### EMP-AUTH-03: Sensitive Action Verification ⭐

When user attempts sensitive action (e.g., refund >EGP 200):

```
┌─────────────────────────────────────────────────────┐
│                                                     │
│  ⚠️ إجراء حساس - تأكيد الهوية                       │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ تفاصيل الإجراء:                                │   │
│  │                                                │   │
│  │ Action: Process Refund                         │   │
│  │ Order ID: #A7X92F                              │   │
│  │ Customer: Ahmed Mohamed                        │   │
│  │ Amount: EGP 300                                │   │
│  │ Reason: Order delayed 45+ minutes              │   │
│  │                                                │   │
│  │ سيتم تسجيل هذا الإجراء في الـ audit log         │   │
│  │ مع بصمتك الحيوية وتفاصيل الجلسة.                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │                                                │   │
│  │       [🔑 أكد ببصمتك]                          │   │
│  │                                                │   │
│  │   ضع إصبعك على مستشعر البصمة                   │   │
│  │                                                │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  [إلغاء]                                           │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### EMP-DASH-01: Main Dashboard

```
┌─────────────────────────────────────────────────────┐
│  👤 Ahmed K. | Support L2 | Today: 4 يوليو 2026    │
│  🟢 Online | Session: 2h 15m                        │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌────────┐ │
│  │ 📋 طلبات │ │ 💰 refunds│ │ 📊 أداء │ │ ⚠️ أنباء│ │
│  │  اليوم   │ │  اليوم   │ │   90%    │ │   0    │ │
│  │   28     │ │   8      │ │   CSAT   │ │        │ │
│  └──────────┘ └──────────┘ └──────────┘ └────────┘ │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Pending Approvals (2):                             │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🔴 Refund Request - EGP 750                  │   │
│  │ From: Sarah M. (Support L2)                  │   │
│  │ Order: #A7X92F | Customer: Ahmed M.          │   │
│  │ Reason: Order never delivered                │   │
│  │ [✅ Approve] [❌ Reject] [👁️ Review]          │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🟡 Restaurant Activation Request              │   │
│  │ From: Omar T. (Ops Manager)                  │   │
│  │ Restaurant: KFC - Tagamoa                    │   │
│  │ Commission: 18%                              │   │
│  │ [✅ Approve] [❌ Reject]                      │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Recent Activity (Audit Trail):                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 14:32 | refund.issued | EGP 300 | #A7X92F   │   │
│  │ 14:25 | customer.viewed | uuid-cust-123      │   │
│  │ 14:18 | order.cancelled | #B3K45L            │   │
│  │ 14:05 | login | biometric | IP: 192.168.1.100│   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### EMP-AUD-01: Audit Log Viewer ⭐

```
┌─────────────────────────────────────────────────────┐
│  📜 Audit Log Viewer                                │
├─────────────────────────────────────────────────────┤
│                                                     │
│  Filter:                                            │
│  Actor: [All ▾] | Action: [All ▾] | Date: [Today]  │
│  Category: [All ▾] | Search: [___________]         │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ Time       | Actor        | Action           │   │
│  │────────────|──────────────|──────────────────│   │
│  │ 14:32:15   | ahmed.k      | refund.issued    │   │
│  │ 14:25:03   | ahmed.k      | customer.viewed  │   │
│  │ 14:18:42   | ahmed.k      | order.cancelled  │   │
│  │ 14:05:00   | ahmed.k      | auth.login       │   │
│  │ 13:45:22   | sarah.m      | restaurant.approved│  │
│  │ ...                                            │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  [📄 Export] [🔍 Advanced Search]                  │
│                                                     │
│  Chain Status: ✅ Verified (last check: 3:00 AM)   │
│  Total Records: 1,247,392                            │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### EMP-AUD-02: Anomaly Alerts (UEBA)

```
┌─────────────────────────────────────────────────────┐
│  ⚠️ Anomaly Alerts (3)                              │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🔴 High Risk: ahmed.k (Support L2)            │   │
│  │ Score: -0.92 (critical)                      │   │
│  │ Reasons:                                      │   │
│  │ • 25 refunds today (avg: 8)                   │   │
│  │ • Avg amount EGP 200 (avg: 80)                │   │
│  │ • Activity at 2 AM (unusual)                  │   │
│  │ • Login from new IP                           │   │
│  │                                              │   │
│  │ Actions:                                      │   │
│  │ [🚫 Freeze Account] [📞 Contact] [👁️ Investigate]│   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
│  ┌──────────────────────────────────────────────┐   │
│  │ 🟡 Medium Risk: sarah.m (Support L2)          │   │
│  │ Score: -0.72                                  │   │
│  │ Reasons:                                      │   │
│  │ • 5 refunds to same customer in 7 days        │   │
│  │                                              │   │
│  │ [👁️ Investigate] [📞 Contact]                 │   │
│  └──────────────────────────────────────────────┘   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 3. UX Flows

### Sensitive Action Flow (Refund >EGP 200)

```
1. User clicks "Process Refund"
   ↓
2. Refund dialog opens (amount, reason)
   ↓
3. User fills details
   ↓
4. Clicks "Process"
   ↓
5. Sensitive Action Dialog appears
   - Shows action details
   - "أكد ببصمتك"
   ↓
6. WebAuthn prompt (Touch ID / Face ID)
   ↓
7. Browser verifies with platform authenticator
   ↓
8. Server validates WebAuthn response
   ↓
9. Action authorized (60s token)
   ↓
10. Refund processed
    ↓
11. Audit log entry created (with biometric_verified=true)
    ↓
12. If amount >EGP 500: dual approval needed
    ↓
13. Second approver notified
    ↓
14. Second approver biometric verification
    ↓
15. Refund executed
```

### Dual Approval Flow

```
1. Initiator (Support L2) creates request
   ↓
2. Biometric verification by initiator
   ↓
3. Request goes to Pending Approvals queue
   ↓
4. Approver (Ops Manager) sees request
   ↓
5. Approver reviews details
   ↓
6. Approver taps "Approve"
   ↓
7. Biometric verification by approver
   ↓
8. Action executed
   ↓
9. Both parties logged in audit
```

---

## 4. Critical UX Considerations

### 4.1 WebAuthn Implementation

- Use `@simplewebauthn/browser` library
- Support platform authenticators (Touch ID, Windows Hello, Face ID)
- Support roaming authenticators (YubiKey) for admins
- Fallback: TOTP for devices without biometric

### 4.2 Audit Log Display

- Read-only (no edit/delete)
- Real-time updates (WebSocket)
- Search by actor, action, entity, date
- Export to CSV (with audit trail)
- Hash chain status indicator (green = verified)

### 4.3 Session Management

- 15 min idle timeout
- 8h max session
- Auto-logout with warning at 5 min remaining
- Re-auth required for sensitive actions (even mid-session)

### 4.4 Permission Visibility

- Users only see actions they have permission for
- Disabled actions show tooltip explaining why
- "Request access" link for needed permissions

---

> Employee Portal uses professional light theme by default. Dark mode available for late-night shifts.
