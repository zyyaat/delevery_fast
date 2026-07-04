// Package kafka implements the event publisher for the Payment Service.
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/food-platform/services/payment/internal/application"
	"github.com/google/uuid"
)

// Publisher implements application.EventPublisher.
type Publisher struct {
	producer interface{ Produce(msg []byte, key string) error }
}

func NewPublisher(producer interface{ Produce(msg []byte, key string) error }) *Publisher {
	return &Publisher{producer: producer}
}

func (p *Publisher) PublishPaymentCaptured(ctx context.Context, event application.PaymentCapturedEvent) error {
	return p.publish(ctx, "captured", event.PaymentID.String(), map[string]interface{}{
		"event_id":       event.EventID.String(),
		"payment_id":     event.PaymentID.String(),
		"order_id":       event.OrderID.String(),
		"amount":         event.Amount,
		"method":         event.Method,
		"provider_txn_id": event.ProviderTxnID,
	})
}

func (p *Publisher) PublishPaymentFailed(ctx context.Context, event application.PaymentFailedEvent) error {
	return p.publish(ctx, "failed", event.PaymentID.String(), map[string]interface{}{
		"event_id":   event.EventID.String(),
		"payment_id": event.PaymentID.String(),
		"order_id":   event.OrderID.String(),
		"reason":     event.Reason,
		"retryable":  event.Retryable,
	})
}

func (p *Publisher) PublishPaymentRefunded(ctx context.Context, event application.PaymentRefundedEvent) error {
	return p.publish(ctx, "refunded", event.PaymentID.String(), map[string]interface{}{
		"event_id":   event.EventID.String(),
		"payment_id": event.PaymentID.String(),
		"order_id":   event.OrderID.String(),
		"amount":     event.Amount,
		"reason":     event.Reason,
	})
}

func (p *Publisher) publish(ctx context.Context, suffix, key string, payload map[string]interface{}) error {
	if p.producer == nil {
		slog.Debug("kafka_event_skipped", "topic", "payment."+suffix, "key", key)
		return nil
	}

	topic := "payment." + suffix
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

// ============ Mock Publisher ============

type MockPublisher struct{}

func NewMockPublisher() *MockPublisher { return &MockPublisher{} }

func (p *MockPublisher) PublishPaymentCaptured(ctx context.Context, event application.PaymentCapturedEvent) error {
	slog.Info("mock_event: payment.captured",
		"payment_id", event.PaymentID,
		"order_id", event.OrderID,
		"amount", event.Amount,
	)
	return nil
}

func (p *MockPublisher) PublishPaymentFailed(ctx context.Context, event application.PaymentFailedEvent) error {
	slog.Info("mock_event: payment.failed",
		"payment_id", event.PaymentID,
		"order_id", event.OrderID,
		"reason", event.Reason,
	)
	return nil
}

func (p *MockPublisher) PublishPaymentRefunded(ctx context.Context, event application.PaymentRefundedEvent) error {
	slog.Info("mock_event: payment.refunded",
		"payment_id", event.PaymentID,
		"order_id", event.OrderID,
		"amount", event.Amount,
	)
	return nil
}

var _ application.EventPublisher = (*MockPublisher)(nil)

// generateEventID helper
func generateEventID() uuid.UUID { return uuid.New() }
