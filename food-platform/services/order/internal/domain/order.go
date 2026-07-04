// Package domain contains the core business logic of the Order Service.
package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ============ Errors ============

var (
	ErrOrderNotFound       = errors.New("order not found")
	ErrInvalidTransition   = errors.New("invalid order status transition")
	ErrOrderNotCancellable = errors.New("order cannot be cancelled at this stage")
	ErrEmptyOrder          = errors.New("order must have at least one item")
	ErrInvalidTotal        = errors.New("invalid order total")
	ErrFraudBlocked        = errors.New("order blocked by fraud detection")
	ErrPaymentRequired     = errors.New("payment is required before confirming order")
)

// ============ Enums ============

// OrderStatus represents the lifecycle state of an order.
type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusPreparing OrderStatus = "preparing"
	StatusReady     OrderStatus = "ready"
	StatusPickedUp  OrderStatus = "picked_up"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
	StatusRefunded  OrderStatus = "refunded"
)

// PaymentMethod represents the payment method for an order.
type PaymentMethod string

const (
	PaymentVodafoneCash PaymentMethod = "vodafone_cash"
	PaymentInstaPay     PaymentMethod = "instapay"
	PaymentCard         PaymentMethod = "card"
	PaymentCOD          PaymentMethod = "cod"
)

// PaymentStatus represents the payment state.
type PaymentStatus string

const (
	PaymentPending  PaymentStatus = "pending"
	PaymentCaptured PaymentStatus = "captured"
	PaymentFailed   PaymentStatus = "failed"
	PaymentRefunded PaymentStatus = "refunded"
)

// ============ State Machine ============

// validTransitions defines which status transitions are allowed.
var validTransitions = map[OrderStatus][]OrderStatus{
	StatusPending:   {StatusConfirmed, StatusCancelled, StatusRefunded},
	StatusConfirmed: {StatusPreparing, StatusCancelled, StatusRefunded},
	StatusPreparing: {StatusReady, StatusCancelled},
	StatusReady:     {StatusPickedUp, StatusCancelled},
	StatusPickedUp:  {StatusDelivered, StatusCancelled},
	StatusDelivered: {StatusRefunded},
	StatusCancelled: {StatusRefunded},
	StatusRefunded:  {}, // Terminal state
}

// CanTransitionTo checks if the order can transition to the new status.
func CanTransitionTo(current, new OrderStatus) bool {
	allowed, ok := validTransitions[current]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == new {
			return true
		}
	}
	return false
}

// IsTerminalStatus returns true if the status is terminal (no further transitions).
func IsTerminalStatus(status OrderStatus) bool {
	return status == StatusRefunded
}

// IsCancellableStatus returns true if the order can still be cancelled.
func IsCancellableStatus(status OrderStatus) bool {
	allowed, ok := validTransitions[status]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == StatusCancelled {
			return true
		}
	}
	return false
}

// ============ Entities ============

// OrderItem represents a single item in an order.
type OrderItem struct {
	id          uuid.UUID
	orderID     uuid.UUID
	menuItemID  uuid.UUID
	name        string
	quantity    int
	unitPrice   float64
	modifiers   string // JSON string of selected modifiers
	notes       string
	lineTotal   float64
}

// NewOrderItem creates a new order item.
func NewOrderItem(menuItemID uuid.UUID, name string, quantity int, unitPrice float64, modifiers string, notes string) (*OrderItem, error) {
	if quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}
	if quantity > 50 {
		return nil, errors.New("quantity cannot exceed 50")
	}
	if unitPrice < 0 {
		return nil, errors.New("unit price cannot be negative")
	}
	if name == "" {
		return nil, errors.New("item name is required")
	}

	lineTotal := unitPrice * float64(quantity)

	return &OrderItem{
		id:         uuid.New(),
		menuItemID: menuItemID,
		name:       name,
		quantity:   quantity,
		unitPrice:  unitPrice,
		modifiers:  modifiers,
		notes:      notes,
		lineTotal:  lineTotal,
	}, nil
}

// Getters
func (oi *OrderItem) ID() uuid.UUID      { return oi.id }
func (oi *OrderItem) OrderID() uuid.UUID { return oi.orderID }
func (oi *OrderItem) MenuItemID() uuid.UUID { return oi.menuItemID }
func (oi *OrderItem) Name() string       { return oi.name }
func (oi *OrderItem) Quantity() int      { return oi.quantity }
func (oi *OrderItem) UnitPrice() float64 { return oi.unitPrice }
func (oi *OrderItem) Modifiers() string  { return oi.modifiers }
func (oi *OrderItem) Notes() string      { return oi.notes }
func (oi *OrderItem) LineTotal() float64 { return oi.lineTotal }

// SetOrderID sets the parent order ID.
func (oi *OrderItem) SetOrderID(id uuid.UUID) {
	oi.orderID = id
}

// Order represents a customer's food order.
type Order struct {
	id              uuid.UUID
	orderNumber     string
	customerID      uuid.UUID
	restaurantID    uuid.UUID
	driverID        *uuid.UUID
	status          OrderStatus
	items           []*OrderItem
	subtotal        float64
	deliveryFee     float64
	serviceFee      float64
	vat             float64
	discount        float64
	total           float64
	paymentMethod   PaymentMethod
	paymentStatus   PaymentStatus
	deliveryAddress string
	latitude        float64
	longitude       float64
	etaMinutes      int
	scheduledFor    *time.Time
	prepStartedAt   *time.Time
	pickedUpAt      *time.Time
	deliveredAt     *time.Time
	cancelReason    string
	notes           string
	cashbackEarned  float64
	createdAt       time.Time
	updatedAt       time.Time
}

// NewOrder creates a new order with the given parameters.
func NewOrder(
	customerID, restaurantID uuid.UUID,
	items []*OrderItem,
	deliveryAddress string,
	lat, lng float64,
	paymentMethod PaymentMethod,
	deliveryFee, serviceFeeRate, vatRate, discount float64,
) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}
	if err := ValidateCoordinates(lat, lng); err != nil {
		return nil, err
	}
	if deliveryAddress == "" {
		return nil, errors.New("delivery address is required")
	}

	// Calculate subtotal
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.LineTotal()
	}

	// Calculate fees
	serviceFee := subtotal * serviceFeeRate
	vat := (subtotal + deliveryFee + serviceFee - discount) * vatRate
	total := subtotal + deliveryFee + serviceFee + vat - discount

	if total < 0 {
		return nil, ErrInvalidTotal
	}

	now := time.Now().UTC()
	order := &Order{
		id:              uuid.New(),
		orderNumber:     generateOrderNumber(),
		customerID:      customerID,
		restaurantID:    restaurantID,
		status:          StatusPending,
		items:           items,
		subtotal:        roundTo2(subtotal),
		deliveryFee:     roundTo2(deliveryFee),
		serviceFee:      roundTo2(serviceFee),
		vat:             roundTo2(vat),
		discount:        roundTo2(discount),
		total:           roundTo2(total),
		paymentMethod:   paymentMethod,
		paymentStatus:   PaymentPending,
		deliveryAddress: deliveryAddress,
		latitude:        lat,
		longitude:       lng,
		etaMinutes:      35, // Default ETA
		notes:           "",
		cashbackEarned:  0,
		createdAt:       now,
		updatedAt:       now,
	}

	// Set order ID on items
	for _, item := range items {
		item.SetOrderID(order.id)
	}

	return order, nil
}

// ============ Getters ============

func (o *Order) ID() uuid.UUID           { return o.id }
func (o *Order) OrderNumber() string     { return o.orderNumber }
func (o *Order) CustomerID() uuid.UUID   { return o.customerID }
func (o *Order) RestaurantID() uuid.UUID { return o.restaurantID }
func (o *Order) DriverID() *uuid.UUID    { return o.driverID }
func (o *Order) Status() OrderStatus     { return o.status }
func (o *Order) Items() []*OrderItem     { return o.items }
func (o *Order) Subtotal() float64       { return o.subtotal }
func (o *Order) DeliveryFee() float64    { return o.deliveryFee }
func (o *Order) ServiceFee() float64     { return o.serviceFee }
func (o *Order) VAT() float64            { return o.vat }
func (o *Order) Discount() float64       { return o.discount }
func (o *Order) Total() float64          { return o.total }
func (o *Order) PaymentMethod() PaymentMethod { return o.paymentMethod }
func (o *Order) PaymentStatus() PaymentStatus { return o.paymentStatus }
func (o *Order) DeliveryAddress() string { return o.deliveryAddress }
func (o *Order) Latitude() float64       { return o.latitude }
func (o *Order) Longitude() float64      { return o.longitude }
func (o *Order) ETAMinutes() int         { return o.etaMinutes }
func (o *Order) ScheduledFor() *time.Time { return o.scheduledFor }
func (o *Order) PrepStartedAt() *time.Time { return o.prepStartedAt }
func (o *Order) PickedUpAt() *time.Time  { return o.pickedUpAt }
func (o *Order) DeliveredAt() *time.Time { return o.deliveredAt }
func (o *Order) CancelReason() string    { return o.cancelReason }
func (o *Order) Notes() string           { return o.notes }
func (o *Order) CashbackEarned() float64 { return o.cashbackEarned }
func (o *Order) CreatedAt() time.Time    { return o.createdAt }
func (o *Order) UpdatedAt() time.Time    { return o.updatedAt }

// ============ Setters ============

func (o *Order) SetNotes(notes string) {
	o.notes = notes
	o.updatedAt = time.Now().UTC()
}

func (o *Order) SetETA(minutes int) {
	o.etaMinutes = minutes
	o.updatedAt = time.Now().UTC()
}

func (o *Order) SetScheduledFor(t *time.Time) {
	o.scheduledFor = t
	o.updatedAt = time.Now().UTC()
}

func (o *Order) SetCashbackEarned(amount float64) {
	o.cashbackEarned = roundTo2(amount)
	o.updatedAt = time.Now().UTC()
}

// ============ State Transitions ============

// TransitionTo attempts to change the order status.
// Returns ErrInvalidTransition if the transition is not allowed.
func (o *Order) TransitionTo(newStatus OrderStatus) error {
	if !CanTransitionTo(o.status, newStatus) {
		return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, o.status, newStatus)
	}

	// Additional validation
	switch newStatus {
	case StatusConfirmed:
		if o.paymentStatus != PaymentCaptured && o.paymentMethod != PaymentCOD {
			return ErrPaymentRequired
		}
	case StatusPreparing:
		now := time.Now().UTC()
		o.prepStartedAt = &now
	case StatusPickedUp:
		now := time.Now().UTC()
		o.pickedUpAt = &now
	case StatusDelivered:
		now := time.Now().UTC()
		o.deliveredAt = &now
	}

	o.status = newStatus
	o.updatedAt = time.Now().UTC()
	return nil
}

// AssignDriver assigns a driver to the order.
func (o *Order) AssignDriver(driverID uuid.UUID) error {
	if o.status != StatusConfirmed && o.status != StatusReady {
		return ErrInvalidTransition
	}
	o.driverID = &driverID
	o.updatedAt = time.Now().UTC()
	return nil
}

// Cancel cancels the order with a reason.
func (o *Order) Cancel(reason string) error {
	if !IsCancellableStatus(o.status) {
		return ErrOrderNotCancellable
	}
	if reason == "" {
		return errors.New("cancel reason is required")
	}

	o.cancelReason = reason
	o.status = StatusCancelled
	o.updatedAt = time.Now().UTC()
	return nil
}

// MarkPaymentCaptured updates the payment status to captured.
func (o *Order) MarkPaymentCaptured() {
	o.paymentStatus = PaymentCaptured
	o.updatedAt = time.Now().UTC()
}

// MarkPaymentFailed updates the payment status to failed.
func (o *Order) MarkPaymentFailed() {
	o.paymentStatus = PaymentFailed
	o.updatedAt = time.Now().UTC()
}

// MarkPaymentRefunded updates the payment status to refunded.
func (o *Order) MarkPaymentRefunded() {
	o.paymentStatus = PaymentRefunded
	o.updatedAt = time.Now().UTC()
}

// CanBeCancelledByCustomer returns true if the customer can still cancel.
// Customers can cancel before preparation starts.
func (o *Order) CanBeCancelledByCustomer() bool {
	return o.status == StatusPending || o.status == StatusConfirmed
}

// CanBeCancelledByRestaurant returns true if the restaurant can cancel.
// Restaurants can cancel before the order is picked up.
func (o *Order) CanBeCancelledByRestaurant() bool {
	return o.status == StatusConfirmed || o.status == StatusPreparing || o.status == StatusReady
}

// IsActive returns true if the order is in an active (non-terminal) state.
func (o *Order) IsActive() bool {
	return o.status != StatusCancelled && o.status != StatusRefunded && o.status != StatusDelivered
}

// ============ Helpers ============

// ValidateCoordinates checks if coordinates are valid.
func ValidateCoordinates(lat, lng float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return errors.New("invalid longitude")
	}
	return nil
}

// roundTo2 rounds to 2 decimal places.
func roundTo2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}

// generateOrderNumber generates a short, human-readable order number.
// Format: 6 alphanumeric characters (e.g., "A7X92F")
func generateOrderNumber() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1) // Ensure uniqueness by adding tiny delay
	}
	return string(b)
}
