// Common types shared across the platform

export type Money = string // decimal as string: "245.50"

export interface Coordinates {
  lat: number
  lng: number
}

export interface Address {
  id: string
  label: string // "Home" | "Work" | "Other"
  street: string
  city: string
  lat: number
  lng: number
  apartment?: string
  building?: string
  notes?: string
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  has_more: boolean
  next_offset?: number
}

export interface Timestamps {
  created_at: string // ISO 8601
  updated_at: string
}

export interface ErrorResponse {
  error: {
    code: string
    message: string
    details?: Record<string, unknown>
    request_id: string
    documentation_url?: string
  }
}

export interface HealthResponse {
  status: 'ok' | 'degraded' | 'down'
  version: string
  uptime_seconds: number
  dependencies?: Record<string, 'ok' | 'degraded' | 'down'>
}

export interface RateLimitInfo {
  limit: number
  remaining: number
  reset: number // Unix timestamp
}

export interface PaginationParams {
  limit?: number
  offset?: number
}

export interface SearchParams extends PaginationParams {
  q?: string
  sort?: string
  order?: 'asc' | 'desc'
}

export interface DateRangeParams {
  from?: string
  to?: string
}

export interface GeoSearchParams {
  lat: number
  lng: number
  radius?: number // km
}

export interface IdempotencyHeaders {
  'X-Idempotency-Key': string
}

export interface AuthHeaders {
  Authorization: `Bearer ${string}`
}

export interface ApiHeaders {
  'X-Request-ID': string
}
