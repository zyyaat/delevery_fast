// Package kafka implements the EventPublisher using Apache Kafka.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/food-platform/services/order/internal/application"
	"github.com/google/uuid"
)

// Publisher implements application.EventPublisher using Kafka.
type Publisher struct {
	producer interface{ Produce(msg []byte, key string) error } // Wraps kafka.Producer
	topicPrefix string
}

// NewPublisher creates a new Kafka event publisher.
// In production, pass a real kafka.Producer.
func NewPublisher(producer interface{ Produce(msg []byte, key string) error }) *Publisher {
	return &Publisher{
		producer:    producer,
		topicPrefix: "order.",
	}
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, event application.OrderCreatedEvent) error {
	return p.publish(ctx, "created", event.OrderID.String(), map[string]interface{}{
		"event_id":       event.EventID.String(),
		"event_time":     ctx.Value("timestamp"),
		"order_id":       event.OrderID.String(),
		"user_id":        event.CustomerID.String(),
		"restaurant_id":  event.RestaurantID.String(),
		"total_amount":   event.TotalAmount,
		"payment_method": event.PaymentMethod,
		"items_count":    event.ItemsCount,
	})
}

func (p *Publisher) PublishOrderConfirmed(ctx context.Context, event application.OrderConfirmedEvent) error {
	driverID := ""
	if event.DriverID != nil {
		driverID = event.DriverID.String()
	}
	return p.publish(ctx, "confirmed", event.OrderID.String(), map[string]interface{}{
		"event_id":    event.EventID.String(),
		"order_id":    event.OrderID.String(),
		"driver_id":   driverID,
		"eta_minutes": event.ETAMinutes,
	})
}

func (p *Publisher) PublishOrderCancelled(ctx context.Context, event application.OrderCancelledEvent) error {
	return p.publish(ctx, "cancelled", event.OrderID.String(), map[string]interface{}{
		"event_id":     event.EventID.String(),
		"order_id":     event.OrderID.String(),
		"reason":       event.Reason,
		"cancelled_by": event.CancelledBy,
	})
}

func (p *Publisher) PublishOrderStatusChanged(ctx context.Context, event application.OrderStatusChangedEvent) error {
	return p.publish(ctx, "status_changed", event.OrderID.String(), map[string]interface{}{
		"event_id":        event.EventID.String(),
		"order_id":        event.OrderID.String(),
		"previous_status": event.PreviousStatus,
		"new_status":      event.NewStatus,
	})
}

func (p *Publisher) publish(ctx context.Context, suffix, key string, payload map[string]interface{}) error {
	if p.producer == nil {
		slog.Debug("kafka_event_skipped", "topic", p.topicPrefix+suffix, "key", key)
		return nil
	}

	topic := p.topicPrefix + suffix
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	if err := p.producer.Produce(data, key); err != nil {
		slog.Error("kafka_publish_failed", "topic", topic, "error", err)
		return err
	}

	slog.Info("kafka_event_published", "topic", topic, "key", key)
	return nil
}

// ============ Mock Publisher (for dev/testing) ============

// MockPublisher logs events instead of publishing to Kafka.
type MockPublisher struct{}

func NewMockPublisher() *MockPublisher {
	return &MockPublisher{}
}

func (p *MockPublisher) PublishOrderCreated(ctx context.Context, event application.OrderCreatedEvent) error {
	slog.Info("mock_event: order.created",
		"order_id", event.OrderID,
		"customer_id", event.CustomerID,
		"total", event.TotalAmount,
	)
	return nil
}

func (p *MockPublisher) PublishOrderConfirmed(ctx context.Context, event application.OrderConfirmedEvent) error {
	slog.Info("mock_event: order.confirmed", "order_id", event.OrderID)
	return nil
}

func (p *MockPublisher) PublishOrderCancelled(ctx context.Context, event application.OrderCancelledEvent) error {
	slog.Info("mock_event: order.cancelled", "order_id", event.OrderID, "reason", event.Reason)
	return nil
}

func (p *MockPublisher) PublishOrderStatusChanged(ctx context.Context, event application.OrderStatusChangedEvent) error {
	slog.Info("mock_event: order.status_changed",
		"order_id", event.OrderID,
		"from", event.PreviousStatus,
		"to", event.NewStatus,
	)
	return nil
}

// Ensure MockPublisher implements EventPublisher
var _ application.EventPublisher = (*MockPublisher)(nil)

// Generate event ID helper
func generateEventID() uuid.UUID {
	return uuid.New()
}
