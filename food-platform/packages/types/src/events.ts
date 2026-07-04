// WebSocket event types

import type { Coordinates, Money } from './common'
import type { OrderStatus, DriverStatus, IncidentSeverity, TicketPriority } from './enums'

export type WSMessageType =
  // Customer
  | 'order.status_changed'
  | 'driver.location'
  | 'driver.assigned'
  | 'order.ready_for_pickup'
  | 'promo.unlocked'
  | 'support.message'
  // Driver
  | 'order.new'
  | 'order.cancelled'
  | 'order.picked_up_confirmed'
  | 'payout.completed'
  | 'zone.hot'
  // Restaurant
  | 'order.new_restaurant'
  | 'order.picked_up'
  | 'order.delivered'
  | 'review.new'
  | 'promo.performance_update'
  | 'menu.updated'
  // Support
  | 'ticket.new'
  | 'ticket.message'
  | 'ticket.escalated'
  | 'alert.fraud'
  // Command Center
  | 'metric.update'
  | 'alert.new'
  | 'incident.new'
  | 'incident.update'
  | 'zone.update'
  // Employee Portal
  | 'approval.requested'
  | 'approval.resolved'
  | 'anomaly.detected'
  | 'audit.tamper_detected'
  // Field Supervisor
  | 'task.assigned'
  | 'task.cancelled'
  | 'complaint.assigned'
  // System
  | 'auth'
  | 'ping'
  | 'pong'
  | 'subscribe'
  | 'unsubscribe'
  | 'error'

export interface WSMessage<T = unknown> {
  event: WSMessageType
  payload: T
  timestamp: string
  id: string
}

// ============ Customer Events ============

export interface OrderStatusChangedPayload {
  order_id: string
  status: OrderStatus
  eta_minutes?: number
}

export interface DriverLocationPayload {
  order_id: string
  lat: number
  lng: number
  heading: number
}

export interface DriverAssignedPayload {
  order_id: string
  driver: {
    id: string
    name: string
    photo_url?: string
    rating: number
    vehicle_type: string
    vehicle_plate?: string
  }
}

export interface PromoUnlockedPayload {
  code: string
  discount: Money
}

// ============ Driver Events ============

export interface OrderNewDriverPayload {
  order_id: string
  restaurant: {
    name: string
    address: string
    distance_km: number
  }
  customer: {
    area: string
    distance_km: number
  }
  earnings: Money
  earnings_breakdown: {
    base: Money
    distance_bonus: Money
    peak_bonus: Money
  }
  eta_minutes: number
  pickup_code: string
  expires_in_seconds: number
}

export interface ZoneHotPayload {
  zones: {
    id: string
    name: string
    demand_level: 'low' | 'medium' | 'high'
    bonus?: Money
  }[]
}

// ============ Restaurant Events ============

export interface OrderNewRestaurantPayload {
  order_id: string
  order_number: string
  customer: {
    name: string
    phone: string
  }
  items: {
    name: string
    quantity: number
    notes?: string
  }[]
  total: Money
  payment_method: string
  delivery_address: string
  notes?: string
  expires_in_seconds: number // 90s timer
}

// ============ Support Events ============

export interface TicketNewPayload {
  ticket_id: string
  customer_name: string
  priority: TicketPriority
  subject: string
  channel: string
}

export interface AlertFraudPayload {
  customer_id: string
  score: number
  reasons: string[]
}

// ============ Command Center Events ============

export interface MetricUpdatePayload {
  metric: string
  value: number
  timestamp: string
}

export interface AlertNewPayload {
  severity: IncidentSeverity
  message: string
  zone?: string
  actions?: string[]
}

export interface IncidentNewPayload {
  incident_id: string
  severity: IncidentSeverity
  title: string
  description: string
}

// ============ Employee Portal Events ============

export interface ApprovalRequestedPayload {
  approval_id: string
  action: string
  requester: {
    id: string
    name: string
  }
  details: Record<string, unknown>
}

export interface AnomalyDetectedPayload {
  employee_id: string
  employee_name: string
  score: number // -1 to 0 (lower = more anomalous)
  reasons: string[]
  severity: 'low' | 'medium' | 'high'
}

// ============ Field Supervisor Events ============

export interface TaskAssignedPayload {
  task_id: string
  type: string
  priority: 'urgent' | 'normal' | 'low'
  title: string
  address: string
  coordinates: Coordinates
  eta_minutes: number
}

// ============ Auth ============

export interface AuthPayload {
  token: string
}

export interface SubscribePayload {
  channel: string
}

export interface ErrorPayload {
  code: string
  message: string
}
