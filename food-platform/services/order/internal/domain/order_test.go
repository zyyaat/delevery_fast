package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// ============ State Machine Tests ============

func TestCanTransitionTo_ValidTransitions(t *testing.T) {
	valid := []struct {
		from, to OrderStatus
	}{
		{StatusPending, StatusConfirmed},
		{StatusPending, StatusCancelled},
		{StatusPending, StatusRefunded},
		{StatusConfirmed, StatusPreparing},
		{StatusConfirmed, StatusCancelled},
		{StatusPreparing, StatusReady},
		{StatusReady, StatusPickedUp},
		{StatusPickedUp, StatusDelivered},
		{StatusDelivered, StatusRefunded},
		{StatusCancelled, StatusRefunded},
	}

	for _, tt := range valid {
		if !CanTransitionTo(tt.from, tt.to) {
			t.Errorf("expected transition %s → %s to be valid", tt.from, tt.to)
		}
	}
}

func TestCanTransitionTo_InvalidTransitions(t *testing.T) {
	invalid := []struct {
		from, to OrderStatus
	}{
		{StatusPending, StatusPreparing},      // Skip confirmed
		{StatusPending, StatusDelivered},      // Skip everything
		{StatusConfirmed, StatusPickedUp},     // Skip preparing + ready
		{StatusPreparing, StatusDelivered},    // Skip ready + picked_up
		{StatusDelivered, StatusPending},      // Can't go back
		{StatusRefunded, StatusPending},       // Terminal state
		{StatusCancelled, StatusConfirmed},    // Can't reactivate
	}

	for _, tt := range invalid {
		if CanTransitionTo(tt.from, tt.to) {
			t.Errorf("expected transition %s → %s to be invalid", tt.from, tt.to)
		}
	}
}

func TestIsTerminalStatus(t *testing.T) {
	if !IsTerminalStatus(StatusRefunded) {
		t.Error("expected refunded to be terminal")
	}
	if IsTerminalStatus(StatusDelivered) {
		t.Error("expected delivered to NOT be terminal (can be refunded)")
	}
	if IsTerminalStatus(StatusPending) {
		t.Error("expected pending to not be terminal")
	}
}

func TestIsCancellableStatus(t *testing.T) {
	cancellable := []OrderStatus{StatusPending, StatusConfirmed, StatusPreparing, StatusReady, StatusPickedUp}
	for _, s := range cancellable {
		if !IsCancellableStatus(s) {
			t.Errorf("expected %s to be cancellable", s)
		}
	}

	notCancellable := []OrderStatus{StatusDelivered, StatusCancelled, StatusRefunded}
	for _, s := range notCancellable {
		if IsCancellableStatus(s) {
			t.Errorf("expected %s to NOT be cancellable", s)
		}
	}
}

// ============ Order Item Tests ============

func TestNewOrderItem_Valid(t *testing.T) {
	item, err := NewOrderItem(uuid.New(), "Big Mac", 2, 85.0, `[]`, "")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if item.Quantity() != 2 {
		t.Errorf("expected quantity 2, got %d", item.Quantity())
	}
	if item.LineTotal() != 170.0 {
		t.Errorf("expected line total 170, got %f", item.LineTotal())
	}
}

func TestNewOrderItem_InvalidQuantity(t *testing.T) {
	_, err := NewOrderItem(uuid.New(), "Item", 0, 10, "", "")
	if err == nil {
		t.Error("expected error for quantity 0")
	}

	_, err = NewOrderItem(uuid.New(), "Item", 51, 10, "", "")
	if err == nil {
		t.Error("expected error for quantity 51")
	}
}

func TestNewOrderItem_NegativePrice(t *testing.T) {
	_, err := NewOrderItem(uuid.New(), "Item", 1, -5, "", "")
	if err == nil {
		t.Error("expected error for negative price")
	}
}

// ============ Order Tests ============

func TestNewOrder_Valid(t *testing.T) {
	items := []*OrderItem{}
	item, _ := NewOrderItem(uuid.New(), "Pizza", 1, 145.0, "", "")
	items = append(items, item)

	order, err := NewOrder(
		uuid.New(), uuid.New(), items,
		"Maadi, Cairo", 30.0444, 31.2357,
		PaymentVodafoneCash,
		20.0, 0.05, 0.14, 0,
	)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if order.Status() != StatusPending {
		t.Errorf("expected pending, got %s", order.Status())
	}
	if order.Subtotal() != 145.0 {
		t.Errorf("expected subtotal 145, got %f", order.Subtotal())
	}
	if order.DeliveryFee() != 20.0 {
		t.Errorf("expected delivery fee 20, got %f", order.DeliveryFee())
	}
	// service fee = 145 * 0.05 = 7.25
	if order.ServiceFee() != 7.25 {
		t.Errorf("expected service fee 7.25, got %f", order.ServiceFee())
	}
	// vat = (145 + 20 + 7.25 - 0) * 0.14 = 24.115 → 24.12
	if order.VAT() < 24.0 || order.VAT() > 24.5 {
		t.Errorf("expected vat ~24.12, got %f", order.VAT())
	}
	// total = 145 + 20 + 7.25 + 24.12 = 196.37
	if order.Total() < 196.0 || order.Total() > 197.0 {
		t.Errorf("expected total ~196.37, got %f", order.Total())
	}
	if order.OrderNumber() == "" {
		t.Error("expected non-empty order number")
	}
}

func TestNewOrder_EmptyItems(t *testing.T) {
	_, err := NewOrder(
		uuid.New(), uuid.New(), []*OrderItem{},
		"addr", 30, 31, PaymentVodafoneCash,
		20, 0.05, 0.14, 0,
	)
	if err != ErrEmptyOrder {
		t.Errorf("expected ErrEmptyOrder, got %v", err)
	}
}

func TestNewOrder_InvalidCoordinates(t *testing.T) {
	item, _ := NewOrderItem(uuid.New(), "Item", 1, 10, "", "")
	_, err := NewOrder(
		uuid.New(), uuid.New(), []*OrderItem{item},
		"addr", 91, 31, PaymentVodafoneCash,
		20, 0.05, 0.14, 0,
	)
	if err == nil {
		t.Error("expected error for invalid coordinates")
	}
}

func TestOrder_TransitionTo_ValidFlow(t *testing.T) {
	order := createTestOrder(t)

	// pending → confirmed (need payment captured first)
	order.MarkPaymentCaptured()
	if err := order.TransitionTo(StatusConfirmed); err != nil {
		t.Errorf("expected pending → confirmed, got %v", err)
	}

	// confirmed → preparing
	if err := order.TransitionTo(StatusPreparing); err != nil {
		t.Errorf("expected confirmed → preparing, got %v", err)
	}
	if order.PrepStartedAt() == nil {
		t.Error("expected prep started time to be set")
	}

	// preparing → ready
	if err := order.TransitionTo(StatusReady); err != nil {
		t.Errorf("expected preparing → ready, got %v", err)
	}

	// ready → picked_up
	if err := order.TransitionTo(StatusPickedUp); err != nil {
		t.Errorf("expected ready → picked_up, got %v", err)
	}
	if order.PickedUpAt() == nil {
		t.Error("expected picked up time to be set")
	}

	// picked_up → delivered
	if err := order.TransitionTo(StatusDelivered); err != nil {
		t.Errorf("expected picked_up → delivered, got %v", err)
	}
	if order.DeliveredAt() == nil {
		t.Error("expected delivered time to be set")
	}
}

func TestOrder_TransitionTo_InvalidFlow(t *testing.T) {
	order := createTestOrder(t)

	// pending → preparing (skip confirmed) should fail
	err := order.TransitionTo(StatusPreparing)
	if err == nil {
		t.Error("expected error for pending → preparing")
	}

	// pending → delivered (skip everything) should fail
	err = order.TransitionTo(StatusDelivered)
	if err == nil {
		t.Error("expected error for pending → delivered")
	}
}

func TestOrder_TransitionTo_ConfirmedRequiresPayment(t *testing.T) {
	order := createTestOrder(t)

	// Payment not captured → can't confirm
	err := order.TransitionTo(StatusConfirmed)
	if err != ErrPaymentRequired {
		t.Errorf("expected ErrPaymentRequired, got %v", err)
	}

	// COD doesn't require captured payment
	order.paymentMethod = PaymentCOD
	err = order.TransitionTo(StatusConfirmed)
	if err != nil {
		t.Errorf("expected COD to not require payment, got %v", err)
	}
}

func TestOrder_Cancel_Valid(t *testing.T) {
	order := createTestOrder(t)

	err := order.Cancel("changed_mind")
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if order.Status() != StatusCancelled {
		t.Errorf("expected cancelled, got %s", order.Status())
	}
	if order.CancelReason() != "changed_mind" {
		t.Errorf("expected reason, got %s", order.CancelReason())
	}
}

func TestOrder_Cancel_AfterDelivered(t *testing.T) {
	order := createTestOrder(t)
	order.MarkPaymentCaptured()
	order.TransitionTo(StatusConfirmed)
	order.TransitionTo(StatusPreparing)
	order.TransitionTo(StatusReady)
	order.TransitionTo(StatusPickedUp)
	order.TransitionTo(StatusDelivered)

	err := order.Cancel("test")
	if err != ErrOrderNotCancellable {
		t.Errorf("expected ErrOrderNotCancellable, got %v", err)
	}
}

func TestOrder_Cancel_EmptyReason(t *testing.T) {
	order := createTestOrder(t)

	err := order.Cancel("")
	if err == nil {
		t.Error("expected error for empty reason")
	}
}

func TestOrder_AssignDriver(t *testing.T) {
	order := createTestOrder(t)
	order.MarkPaymentCaptured()
	order.TransitionTo(StatusConfirmed)

	driverID := uuid.New()
	err := order.AssignDriver(driverID)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if order.DriverID() == nil || *order.DriverID() != driverID {
		t.Error("expected driver ID to match")
	}
}

func TestOrder_AssignDriver_InvalidStatus(t *testing.T) {
	order := createTestOrder(t)

	// Status is pending, can't assign driver
	err := order.AssignDriver(uuid.New())
	if err == nil {
		t.Error("expected error when assigning driver to pending order")
	}
}

func TestOrder_CanBeCancelledByCustomer(t *testing.T) {
	order := createTestOrder(t)

	// Pending → customer can cancel
	if !order.CanBeCancelledByCustomer() {
		t.Error("expected customer to cancel pending order")
	}

	// Confirmed → customer can cancel
	order.MarkPaymentCaptured()
	order.TransitionTo(StatusConfirmed)
	if !order.CanBeCancelledByCustomer() {
		t.Error("expected customer to cancel confirmed order")
	}

	// Preparing → customer cannot cancel
	order.TransitionTo(StatusPreparing)
	if order.CanBeCancelledByCustomer() {
		t.Error("expected customer to NOT cancel preparing order")
	}
}

func TestOrder_IsActive(t *testing.T) {
	order := createTestOrder(t)
	if !order.IsActive() {
		t.Error("expected pending order to be active")
	}

	order.Cancel("test")
	if order.IsActive() {
		t.Error("expected cancelled order to not be active")
	}
}

func TestOrder_MarkPaymentCaptured(t *testing.T) {
	order := createTestOrder(t)
	order.MarkPaymentCaptured()
	if order.PaymentStatus() != PaymentCaptured {
		t.Errorf("expected captured, got %s", order.PaymentStatus())
	}
}

func TestOrder_MarkPaymentFailed(t *testing.T) {
	order := createTestOrder(t)
	order.MarkPaymentFailed()
	if order.PaymentStatus() != PaymentFailed {
		t.Errorf("expected failed, got %s", order.PaymentStatus())
	}
}

func TestOrder_SetETA(t *testing.T) {
	order := createTestOrder(t)
	order.SetETA(45)
	if order.ETAMinutes() != 45 {
		t.Errorf("expected 45, got %d", order.ETAMinutes())
	}
}

func TestOrder_SetCashbackEarned(t *testing.T) {
	order := createTestOrder(t)
	order.SetCashbackEarned(31.5)
	if order.CashbackEarned() != 31.5 {
		t.Errorf("expected 31.5, got %f", order.CashbackEarned())
	}
}

func TestOrder_SetScheduledFor(t *testing.T) {
	order := createTestOrder(t)
	tm := time.Now().Add(2 * time.Hour)
	order.SetScheduledFor(&tm)
	if order.ScheduledFor() == nil {
		t.Error("expected scheduled time to be set")
	}
}

// ============ Helpers ============

func createTestOrder(t *testing.T) *Order {
	t.Helper()
	item, _ := NewOrderItem(uuid.New(), "Test Item", 1, 100.0, "", "")
	order, err := NewOrder(
		uuid.New(), uuid.New(),
		[]*OrderItem{item},
		"Test Address", 30.0, 31.0,
		PaymentVodafoneCash,
		20.0, 0.05, 0.14, 0,
	)
	if err != nil {
		t.Fatalf("failed to create test order: %v", err)
	}
	return order
}
