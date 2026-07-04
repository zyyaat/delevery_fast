// API request/response types

import type { Money, Coordinates, SearchParams as CommonSearchParams } from './common'
import type {
  OrderStatus,
  PaymentMethod,
  UserRole,
  DriverStatus,
  PromoType,
  TicketPriority,
  IncidentSeverity,
  IncidentStatus,
  FieldTaskType,
  FieldTaskPriority,
} from './enums'
import type {
  User,
  Restaurant,
  MenuItem,
} from './entities'

// Re-export with alias to avoid name collision
export type { CommonSearchParams as SearchParams }

// ============ Auth ============

export interface SendOtpRequest {
  phone: string
  role: UserRole
}

export interface SendOtpResponse {
  request_id: string
  expires_in: number
  attempts_remaining: number
}

export interface VerifyOtpRequest {
  request_id: string
  code: string
  device_fingerprint: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: 'Bearer'
  user: User
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface EmployeeLoginRequest {
  username: string
  password: string
  totp_code: string
}

// ============ Customer ============

export interface UpdateCustomerRequest {
  name?: string
  email?: string
}

export interface CreateAddressRequest {
  label: string
  street: string
  city: string
  lat: number
  lng: number
  apartment?: string
  building?: string
  notes?: string
}

// ============ Discover ============

export interface NearbyRestaurantsParams {
  lat: number
  lng: number
  radius?: number
  limit?: number
  offset?: number
}

export interface NearbyRestaurantsResponse {
  restaurants: Restaurant[]
  total: number
  has_more: boolean
}

export interface RestaurantSearchParams {
  q: string
  lat: number
  lng: number
}

export interface RestaurantDetail extends Restaurant {
  menu_categories: MenuCategoryWithItems[]
}

export interface MenuCategoryWithItems {
  id: string
  name: string
  display_order: number
  items: MenuItem[]
}

// ============ Cart ============

export interface AddToCartRequest {
  restaurant_id: string
  menu_item_id: string
  quantity: number
  modifiers: {
    modifier_id: string
    option_id: string
  }[]
  notes?: string
}

export interface ApplyCouponRequest {
  code: string
}

// ============ Orders ============

export interface CreateOrderRequest {
  address_id: string
  payment_method: PaymentMethod
  scheduled_for?: string
  notes?: string
}

export interface CreateOrderResponse {
  id: string
  order_number: string
  status: OrderStatus
  total: Money
  eta_minutes: number
  payment: {
    method: PaymentMethod
    status: string
    provider_redirect_url?: string
  }
}

export interface CancelOrderRequest {
  reason: string
}

export interface RateOrderRequest {
  overall_rating: number
  restaurant_rating?: number
  driver_rating?: number
  comment?: string
}

// ============ Driver ============

export interface UpdateDriverStatusRequest {
  status: DriverStatus
}

export interface UpdateDriverLocationRequest {
  lat: number
  lng: number
  heading: number
  speed: number
}

export interface AcceptOrderResponse {
  order_id: string
  status: OrderStatus
  pickup_code: string
}

export interface PayoutRequest {
  amount: Money
  method: PaymentMethod
}

// ============ Restaurant ============

export interface CreateMenuItemRequest {
  category_id: string
  name: string
  description?: string
  price: Money
  image_url?: string
  prep_time_minutes: number
  modifiers?: MenuItem['modifiers']
}

export interface UpdateMenuItemRequest extends Partial<CreateMenuItemRequest> {
  is_available?: boolean
}

export interface CreatePromoRequest {
  type: PromoType
  value: number
  valid_from: string
  valid_to: string
  usage_limit: number
  per_user_limit: number
  min_order_value?: Money
  item_id?: string
}

// ============ Support ============

export interface CreateTicketRequest {
  customer_id?: string
  driver_id?: string
  restaurant_id?: string
  order_id?: string
  priority: TicketPriority
  subject: string
  initial_message: string
  channel: 'chat' | 'email' | 'phone' | 'whatsapp' | 'social'
}

export interface SendMessageRequest {
  text: string
  attachments?: string[]
}

export interface RefundRequest {
  order_id: string
  amount: Money
  reason: string
  type: 'full' | 'partial' | 'coupon'
}

// ============ Command Center ============

export interface SurgeOverrideRequest {
  multiplier: number
  duration_minutes: number
  reason: string
}

export interface CreateIncidentRequest {
  severity: IncidentSeverity
  title: string
  description: string
  affected_services: string[]
}

export interface UpdateIncidentRequest {
  status?: IncidentStatus
  notes?: string
}

// ============ Field Supervisor ============

export interface CompleteTaskRequest {
  result: 'approved' | 'conditional' | 'rejected'
  checklist: Record<string, boolean>
  photos: {
    url: string
    type: string
    lat: number
    lng: number
    timestamp: string
  }[]
  notes: string
  conditions?: string[]
  reason?: string
}

export interface FieldTask {
  id: string
  type: FieldTaskType
  priority: FieldTaskPriority
  title: string
  address: string
  coordinates: Coordinates
  eta_minutes: number
  distance_km: number
  scheduled_at: string
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled'
}

// ============ Audit Log ============

export interface AuditLogEntry {
  id: string
  actor_id: string
  actor_email: string
  actor_role: UserRole
  action: string
  action_category: string
  entity_type: string
  entity_id?: string
  metadata: Record<string, unknown>
  biometric_verified: boolean
  dual_approval: boolean
  approver_id?: string
  ip_address: string
  user_agent: string
  created_at: string
}

// ============ Analytics ============

export interface AnalyticsPoint {
  timestamp: string
  value: number
}

export interface RestaurantAnalytics {
  sales_today: Money
  orders_today: number
  avg_prep_time: number
  rating: number
  rating_count: number
  sales_chart: AnalyticsPoint[]
  top_items: {
    id: string
    name: string
    count: number
    revenue: Money
  }[]
  peak_hours: {
    hour: number
    order_count: number
  }[]
}
