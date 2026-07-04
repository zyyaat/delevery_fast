// Package infrastructure implements Kafka event consumers for the Notification Service.
package infrastructure

import (
        "context"
        "encoding/json"
        "fmt"
        "log/slog"

        "github.com/food-platform/notification/internal/application"
        "github.com/food-platform/notification/internal/domain"
        "github.com/google/uuid"
)

// EventHandler handles events from other services and sends notifications.
type EventHandler struct {
        sendNotif *application.SendNotificationUseCase
}

func NewEventHandler(sendNotif *application.SendNotificationUseCase) *EventHandler {
        return &EventHandler{sendNotif: sendNotif}
}

// ============ Order Events ============

// HandleOrderCreated handles order.created events.
// Notifies the restaurant about a new order.
func (h *EventHandler) HandleOrderCreated(ctx context.Context, data []byte) error {
        var event struct {
                OrderID      string `json:"order_id"`
                UserID       string `json:"user_id"`
                RestaurantID string `json:"restaurant_id"`
                TotalAmount  float64 `json:"total_amount"`
        }

        if err := json.Unmarshal(data, &event); err != nil {
                return err
        }

        restaurantID, _ := uuid.Parse(event.RestaurantID)

        return h.sendNotif.Execute(
                ctx,
                restaurantID,
                "restaurant",
                domain.ChannelWebSocket,
                "🔔 طلب جديد!",
                fmt.Sprintf("طلب جديد بقيمة EGP %.2f — تفقد التطبيق", event.TotalAmount),
                map[string]interface{}{
                        "type":      "order.new",
                        "order_id":  event.OrderID,
                        "expires_in_seconds": 90,
                },
        )
}

// HandleOrderStatusChanged handles order.status_changed events.
// Notifies the customer about status changes.
func (h *EventHandler) HandleOrderStatusChanged(ctx context.Context, data []byte) error {
        var event struct {
                OrderID        string `json:"order_id"`
                PreviousStatus string `json:"previous_status"`
                NewStatus      string `json:"new_status"`
        }

        if err := json.Unmarshal(data, &event); err != nil {
                return err
        }

        // In production, fetch customer_id from Order Service
        // For now, log and skip
        slog.Info("order_status_changed_notification",
                "order_id", event.OrderID,
                "from", event.PreviousStatus,
                "to", event.NewStatus,
        )

        // Map status to notification message
        messages := map[string]string{
                "confirmed":  "✅ تم تأكيد طلبك! المطعم بدأ التحضير",
                "preparing":  "👨‍🍳 المطعم بيحضّر طلبك",
                "ready":      "📦 طلبك جاهز! بنبحث عن مندوب",
                "picked_up":  "🛵 المندوب استلم طلبك — في الطريق!",
                "delivered":  "🎉 طلبك وصل! بالهنا والشفا",
                "cancelled":  "❌ تم إلغاء طلبك",
        }

        body, ok := messages[event.NewStatus]
        if !ok {
                return nil
        }

        // TODO: Get customer_id from Order Service via gRPC
        // For now, skip
        return nil
}

// HandleOrderCancelled handles order.cancelled events.
func (h *EventHandler) HandleOrderCancelled(ctx context.Context, data []byte) error {
        var event struct {
                OrderID string `json:"order_id"`
                Reason  string `json:"reason"`
        }

        if err := json.Unmarshal(data, &event); err != nil {
                return err
        }

        slog.Info("order_cancelled_notification", "order_id", event.OrderID, "reason", event.Reason)
        return nil
}

// ============ Payment Events ============

// HandlePaymentCaptured handles payment.captured events.
func (h *EventHandler) HandlePaymentCaptured(ctx context.Context, data []byte) error {
        var event struct {
                PaymentID string  `json:"payment_id"`
                OrderID   string  `json:"order_id"`
                Amount    float64 `json:"amount"`
                Method    string  `json:"method"`
        }

        if err := json.Unmarshal(data, &event); err != nil {
                return err
        }

        slog.Info("payment_captured_notification",
                "payment_id", event.PaymentID,
                "order_id", event.OrderID,
                "amount", event.Amount,
        )
        return nil
}

// HandlePaymentFailed handles payment.failed events.
func (h *EventHandler) HandlePaymentFailed(ctx context.Context, data []byte) error {
        var event struct {
                PaymentID string `json:"payment_id"`
                OrderID   string `json:"order_id"`
                Reason    string `json:"reason"`
        }

        if err := json.Unmarshal(data, &event); err != nil {
                return err
        }

        slog.Info("payment_failed_notification",
                "payment_id", event.PaymentID,
                "order_id", event.OrderID,
                "reason", event.Reason,
        )
        return nil
}

// HandlePaymentRefunded handles payment.refunded events.
func (h *EventHandler) HandlePaymentRefunded(ctx context.Context, data []byte) error {
        var event struct {
                PaymentID string  `json:"payment_id"`
                OrderID   string  `json:"order_id"`
                Amount    float64 `json:"amount"`
                Reason    string  `json:"reason"`
        }

        if err := json.Unmarshal(data, &event); err != nil {
                return err
        }

        slog.Info("payment_refunded_notification",
                "payment_id", event.PaymentID,
                "order_id", event.OrderID,
                "amount", event.Amount,
                "reason", event.Reason,
        )
        return nil
}
