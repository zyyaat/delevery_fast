// Core entities

import type { Timestamps, Money, Coordinates, Address } from './common'
import type {
  OrderStatus,
  PaymentMethod,
  PaymentStatus,
  UserRole,
  DriverStatus,
  DriverTier,
  CustomerTier,
  RestaurantStatus,
  VehicleType,
  CuisineType,
} from './enums'

// ============ User / Customer ============

export interface User {
  id: string
  phone: string
  email?: string
  name: string
  role: UserRole
  is_active: boolean
  trust_score: number // 0-100
}

export interface Customer extends User {
  role: 'customer'
  tier: CustomerTier
  loyalty_points: number
  wallet_balance: Money
  total_orders: number
  total_spent: Money
  addresses: Address[]
}

// ============ Restaurant ============

export interface Restaurant {
  id: string
  name: string
  slug: string
  cuisine_types: CuisineType[]
  rating: number
  rating_count: number
  logo_url?: string
  cover_url?: string
  location: Coordinates
  address: string
  is_open: boolean
  status: RestaurantStatus
  eta_minutes_min: number
  eta_minutes_max: number
  delivery_fee: Money
  price_range: 1 | 2 | 3 | 4
  commission_rate: number
  promo?: RestaurantPromo
}

export interface RestaurantPromo {
  id: string
  title: string
  description: string
  type: 'flat' | 'percentage' | 'free_delivery' | 'bogo'
  value: number
}

// ============ Menu ============

export interface MenuCategory {
  id: string
  name: string
  display_order: number
}

export interface MenuItem {
  id: string
  restaurant_id: string
  category_id: string
  name: string
  description?: string
  price: Money
  image_url?: string
  is_available: boolean
  prep_time_minutes: number
  rating?: number
  rating_count?: number
  is_most_ordered?: boolean
  modifiers?: MenuItemModifier[]
}

export interface MenuItemModifier {
  id: string
  name: string
  required: boolean
  multiple_choice: boolean
  options: MenuItemModifierOption[]
}

export interface MenuItemModifierOption {
  id: string
  name: string
  price_delta: number
}

// ============ Cart ============

export interface CartItem {
  id: string
  menu_item_id: string
  name: string
  image_url?: string
  quantity: number
  unit_price: Money
  modifiers: CartItemModifier[]
  notes?: string
  line_total: Money
}

export interface CartItemModifier {
  modifier_id: string
  modifier_name: string
  option_id: string
  option_name: string
  price_delta: number
}

export interface Cart {
  restaurant_id: string
  items: CartItem[]
  subtotal: Money
  delivery_fee: Money
  service_fee: Money
  vat: Money
  discount: Money
  total: Money
  cashback_earned: Money
  applied_coupon?: string
}

// ============ Order ============

export interface Order extends Timestamps {
  id: string
  order_number: string
  customer_id: string
  restaurant_id: string
  driver_id?: string
  status: OrderStatus
  items: OrderItem[]
  subtotal: Money
  delivery_fee: Money
  service_fee: Money
  vat: Money
  discount: Money
  total: Money
  payment_method: PaymentMethod
  payment_status: PaymentStatus
  delivery_address: Address
  coordinates: Coordinates
  eta_minutes: number
  scheduled_for?: string
  prep_started_at?: string
  picked_up_at?: string
  delivered_at?: string
  cancel_reason?: string
  notes?: string
  cashback_earned?: Money
}

export interface OrderItem {
  id: string
  order_id: string
  menu_item_id: string
  name: string
  quantity: number
  unit_price: Money
  modifiers?: CartItemModifier[]
  notes?: string
  line_total: Money
}

// ============ Payment ============

export interface Payment {
  id: string
  order_id: string
  method: PaymentMethod
  amount: Money
  status: PaymentStatus
  provider_txn_id?: string
  refunded_amount?: Money
  refund_reason?: string
  created_at: string
}

// ============ Driver ============

export interface Driver {
  id: string
  user_id: string
  name: string
  phone: string
  vehicle_type: VehicleType
  vehicle_plate?: string
  rating: number
  rating_count: number
  tier: DriverTier
  trust_score: number
  acceptance_rate: number
  completion_rate: number
  total_earnings: Money
  is_online: boolean
  status: DriverStatus
  location?: Coordinates
  photo_url?: string
}

export interface DriverEarnings {
  today: Money
  this_week: Money
  this_month: Money
  today_deliveries: number
  today_hours: number
  hourly_rate: Money
  pending_payout: Money
}

export interface DriverLocation {
  driver_id: string
  lat: number
  lng: number
  heading: number
  speed: number
  timestamp: string
}

// ============ Delivery ============

export interface Delivery {
  id: string
  order_id: string
  driver_id: string
  pickup_at?: string
  picked_up_at?: string
  delivered_at?: string
  distance_km: number
  duration_minutes: number
  earnings: Money
}

// ============ Promo / Loyalty ============

export interface Promo {
  id: string
  code: string
  title: string
  description: string
  type: PromoType
  value: number
  valid_from: string
  valid_to: string
  usage_limit: number
  used_count: number
  per_user_limit: number
  min_order_value?: Money
  restaurant_id?: string
  is_active: boolean
}

export interface LoyaltyAccount {
  customer_id: string
  tier: CustomerTier
  points: number
  total_earned: number
  total_redeemed: number
  wallet_balance: Money
}
