// Enums used across the platform

export type OrderStatus =
  | 'pending'
  | 'confirmed'
  | 'preparing'
  | 'ready'
  | 'picked_up'
  | 'delivered'
  | 'cancelled'
  | 'refunded'

export type PaymentMethod = 'vodafone_cash' | 'instapay' | 'card' | 'cod'

export type PaymentStatus = 'pending' | 'captured' | 'failed' | 'refunded' | 'partial_refunded'

export type UserRole =
  | 'customer'
  | 'driver'
  | 'restaurant'
  | 'support_l1'
  | 'support_l2'
  | 'ops_manager'
  | 'finance'
  | 'super_admin'
  | 'field_supervisor'
  | 'hr'
  | 'read_only_analyst'

export type DriverStatus = 'offline' | 'online' | 'on_break' | 'on_delivery'

export type DriverTier = 'platinum' | 'gold' | 'silver' | 'standard'

export type CustomerTier = 'platinum' | 'gold' | 'silver' | 'standard'

export type RestaurantStatus = 'active' | 'paused' | 'suspended' | 'pending_verification' | 'rejected'

export type VehicleType = 'motorcycle' | 'car' | 'bicycle'

export type PromoType = 'flat' | 'percentage' | 'free_delivery' | 'bogo'

export type CuisineType =
  | 'egyptian'
  | 'italian'
  | 'asian'
  | 'fast_food'
  | 'healthy'
  | 'desserts'
  | 'indian'
  | 'lebanese'
  | 'breakfast'
  | 'coffee'
  | 'groceries'

export type IncidentSeverity = 'P0' | 'P1' | 'P2' | 'P3'

export type IncidentStatus = 'detected' | 'triaged' | 'acknowledged' | 'investigating' | 'identified' | 'mitigating' | 'resolved' | 'postmortem'

export type TicketPriority = 'P0' | 'P1' | 'P2' | 'P3'

export type TicketStatus = 'created' | 'queued' | 'assigned' | 'in_progress' | 'resolved' | 'escalated' | 'closed'

export type AuditActionCategory = 'financial' | 'operational' | 'data_access' | 'auth'

export type FieldTaskType = 'restaurant_verification' | 'driver_verification' | 'audit' | 'complaint_investigation' | 'training'

export type FieldTaskPriority = 'urgent' | 'normal' | 'low'

export type FieldTaskResult = 'approved' | 'conditional' | 'rejected'

export type Sentiment = 'very_angry' | 'angry' | 'neutral' | 'satisfied' | 'very_satisfied'
