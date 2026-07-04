// Package application contains use cases for the Notification Service.
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/food-platform/services/notification/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

type NotificationRepository interface {
	Save(ctx context.Context, n *domain.Notification) error
	FindByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Notification, error)
	MarkDelivered(ctx context.Context, id uuid.UUID) error
}

type PushSender interface {
	Send(ctx context.Context, userID uuid.UUID, title, body string, data map[string]interface{}) error
}

type SMSSender interface {
	Send(ctx context.Context, phone, message string) error
}

type WebSocketBroadcaster interface {
	BroadcastToUser(ctx context.Context, userID uuid.UUID, message map[string]interface{}) error
}

// ============ DTOs ============

type NotificationDTO struct {
	ID        uuid.UUID              `json:"id"`
	Channel   string                 `json:"channel"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Status    string                 `json:"status"`
	CreatedAt string                 `json:"created_at"`
}

// ============ Use Cases ============

// SendNotificationUseCase sends a notification via the appropriate channel.
type SendNotificationUseCase struct {
	repo       NotificationRepository
	pushSender PushSender
	smsSender  SMSSender
	wsBroadcaster WebSocketBroadcaster
}

func NewSendNotificationUseCase(
	repo NotificationRepository,
	push PushSender,
	sms SMSSender,
	ws WebSocketBroadcaster,
) *SendNotificationUseCase {
	return &SendNotificationUseCase{
		repo:          repo,
		pushSender:    push,
		smsSender:     sms,
		wsBroadcaster: ws,
	}
}

func (uc *SendNotificationUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
	userRole string,
	channel domain.NotificationChannel,
	title, body string,
	data map[string]interface{},
) error {
	notif := domain.NewNotification(userID, userRole, channel, title, body, data)

	// Save notification
	if err := uc.repo.Save(ctx, notif); err != nil {
		return fmt.Errorf("save notification: %w", err)
	}

	// Send via appropriate channel
	var sendErr error
	switch channel {
	case domain.ChannelPush:
		sendErr = uc.pushSender.Send(ctx, userID, title, body, data)
	case domain.ChannelSMS:
		// In production, fetch user's phone from user service
		sendErr = uc.smsSender.Send(ctx, "", body)
	case domain.ChannelWebSocket:
		wsData := map[string]interface{}{
			"event":     "notification",
			"title":     title,
			"body":      body,
			"data":      data,
			"timestamp": notif.CreatedAt(),
		}
		sendErr = uc.wsBroadcaster.BroadcastToUser(ctx, userID, wsData)
	case domain.ChannelInApp:
		// In-app notifications are just stored in DB; client polls or gets via WebSocket
		notif.MarkSent()
		_ = uc.repo.MarkDelivered(ctx, notif.ID())
		return nil
	}

	if sendErr != nil {
		slog.Error("notification_send_failed",
			"channel", channel,
			"user_id", userID,
			"error", sendErr,
		)
		notif.MarkFailed()
		return sendErr
	}

	notif.MarkSent()
	slog.Info("notification_sent",
		"channel", channel,
		"user_id", userID,
		"notification_id", notif.ID(),
	)
	return nil
}

// GetNotificationsUseCase retrieves notifications for a user.
type GetNotificationsUseCase struct {
	repo NotificationRepository
}

func NewGetNotificationsUseCase(repo NotificationRepository) *GetNotificationsUseCase {
	return &GetNotificationsUseCase{repo: repo}
}

func (uc *GetNotificationsUseCase) Execute(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*NotificationDTO, error) {
	if limit == 0 {
		limit = 20
	}

	notifications, err := uc.repo.FindByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	dtos := make([]*NotificationDTO, 0, len(notifications))
	for _, n := range notifications {
		dtos = append(dtos, &NotificationDTO{
			ID:        n.ID(),
			Channel:   string(n.Channel()),
			Title:     n.Title(),
			Body:      n.Body(),
			Data:      n.Data(),
			Status:    string(n.Status()),
			CreatedAt: n.CreatedAt().Format("2006-01-02T15:04:05Z"),
		})
	}
	return dtos, nil
}
