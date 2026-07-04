// Common constants

export const EGP_CURRENCY = 'EGP'
export const EGP_LOCALE = 'ar-EG'
export const DEFAULT_TIMEZONE = 'Africa/Cairo'

export const PHONE_REGEX = /^01[0-2,5][0-9]{8}$/
export const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
export const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i

export const OTP_LENGTH = 6
export const OTP_EXPIRY_SECONDS = 120
export const OTP_MAX_ATTEMPTS = 3

export const ORDER_TIMER_SECONDS = 90 // Restaurant
export const DRIVER_OFFER_TIMER_SECONDS = 15

export const JWT_ACCESS_TTL = 900 // 15 min
export const JWT_REFRESH_TTL = 2592000 // 30 days

export const PAGINATION_DEFAULT_LIMIT = 20
export const PAGINATION_MAX_LIMIT = 100

export const RATE_LIMITS = {
  PUBLIC: { requests: 10, window: 60 },
  CUSTOMER: { requests: 100, window: 60 },
  DRIVER: { requests: 200, window: 60 },
  RESTAURANT: { requests: 500, window: 60 },
  EMPLOYEE: { requests: 60, window: 60 },
  API_KEY: { requests: 1000, window: 60 },
} as const

export const TIER_THRESHOLDS = {
  PLATINUM: { rating: 4.9, acceptance: 0.95, completion: 0.98 },
  GOLD: { rating: 4.7, acceptance: 0.85, completion: 0.95 },
  SILVER: { rating: 4.5, acceptance: 0.7, completion: 0.9 },
  STANDARD: { rating: 0, acceptance: 0, completion: 0 },
} as const

export const FRAUD_THRESHOLDS = {
  AUTO_APPROVE: 40,
  REVIEW: 70,
  AUTO_BLOCK: 80,
} as const

export const REFUND_THRESHOLDS = {
  SUPPORT_L1_MAX: 100,
  SUPPORT_L2_MAX: 500,
  DUAL_APPROVAL_REQUIRED: 500,
  BIOMETRIC_REQUIRED: 100,
} as const

export const APP_NAMES = {
  CUSTOMER: 'تطبيق العميل',
  DRIVER: 'تطبيق المندوب',
  RESTAURANT: 'تطبيق المطعم',
  SUPPORT: 'تطبيق الدعم',
  COMMAND_CENTER: 'مركز القيادة',
  EMPLOYEE_PORTAL: 'بوابة الموظفين',
  FIELD_SUPERVISOR: 'تطبيق المشرف الميداني',
} as const

export const PAYMENT_METHODS = {
  VODAFONE_CASH: 'vodafone_cash',
  INSTAPAY: 'instapay',
  CARD: 'card',
  COD: 'cod',
} as const

export const PAYMENT_METHOD_LABELS: Record<string, string> = {
  vodafone_cash: 'Vodafone Cash',
  instapay: 'InstaPay',
  card: 'بطاقة بنكية',
  cod: 'الدفع عند الاستلام',
}
