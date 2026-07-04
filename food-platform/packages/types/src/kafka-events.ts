// Kafka event types (AsyncAPI equivalent)

import type { Money, Coordinates, Address } from './common'
import type { OrderStatus, PaymentMethod, DriverStatus, UserRole } from './enums'

// ============ Order Events ============

export interface OrderCreatedEvent {
  event_id: string
  event_time: string
  order_id: string
  user_id: string
  restaurant_id: string
  total_amount: Money
  payment_method: PaymentMethod
  delivery_address: Address
  items_count: number
  metadata?: Record<string, unknown>
}

export interface OrderConfirmedEvent {
  event_id: string
  event_time: string
  order_id: string
  driver_id: string
  eta_minutes: number
}

export interface OrderCancelledEvent {
  event_id: string
  event_time: string
  order_id: string
  reason: string
  cancelled_by: UserRole
  refund_amount?: Money
}

export interface OrderStatusChangedEvent {
  event_id: string
  event_time: string
  order_id: string
  previous_status: OrderStatus
  new_status: OrderStatus
  actor: {
    type: 'system' | 'customer' | 'driver' | 'restaurant' | 'support' | 'ops'
    id?: string
  }
}

// ============ Payment Events ============

export interface PaymentCapturedEvent {
  event_id: string
  event_time: string
  order_id: string
  payment_id: string
  amount: Money
  method: PaymentMethod
  provider_txn_id: string
}

export interface PaymentFailedEvent {
  event_id: string
  event_time: string
  order_id: string
  payment_id: string
  reason: string
  retryable: boolean
}

export interface PaymentRefundedEvent {
  event_id: string
  event_time: string
  order_id: string
  payment_id: string
  amount: Money
  reason: string
  refunded_by: string
}

// ============ Driver Events ============

export interface DriverLocationEvent {
  event_id: string
  event_time: string
  driver_id: string
  lat: number
  lng: number
  heading: number
  speed: number
}

export interface DriverStatusChangedEvent {
  event_id: string
  event_time: string
  driver_id: string
  previous_status: DriverStatus
  new_status: DriverStatus
  location?: Coordinates
}

// ============ Restaurant Events ============

export interface RestaurantMenuUpdatedEvent {
  event_id: string
  event_time: string
  restaurant_id: string
  changes: {
    type: 'item_added' | 'item_updated' | 'item_removed' | 'availability_changed'
    item_id: string
    is_available?: boolean
  }[]
}

// ============ Fraud Events ============

export interface FraudScoreCalculatedEvent {
  event_id: string
  event_time: string
  order_id: string
  customer_id: string
  score: number // 0-100
  reasons: string[]
  model_version: string
  decision: 'approve' | 'review' | 'block'
}

// ============ Audit Events ============

export interface AuditActionLoggedEvent {
  event_id: string
  event_time: string
  audit_log_id: string
  actor_id: string
  action: string
  entity_type: string
  entity_id?: string
  biometric_verified: boolean
  dual_approval: boolean
  record_hash: string
  prev_hash: string
}

// ============ Topic Definitions ============

export interface KafkaTopics {
  'order.created': OrderCreatedEvent
  'order.confirmed': OrderConfirmedEvent
  'order.cancelled': OrderCancelledEvent
  'order.status_changed': OrderStatusChangedEvent
  'payment.captured': PaymentCapturedEvent
  'payment.failed': PaymentFailedEvent
  'payment.refunded': PaymentRefundedEvent
  'driver.location': DriverLocationEvent
  'driver.status_changed': DriverStatusChangedEvent
  'restaurant.menu_updated': RestaurantMenuUpdatedEvent
  'fraud.score_calculated': FraudScoreCalculatedEvent
  'audit.action_logged': AuditActionLoggedEvent
}

export type KafkaTopicName = keyof KafkaTopics
