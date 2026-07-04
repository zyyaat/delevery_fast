// Package domain contains the core logic of the Notification Service.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
)

type NotificationChannel string

const (
	ChannelPush     NotificationChannel = "push"
	ChannelSMS      NotificationChannel = "sms"
	ChannelEmail    NotificationChannel = "email"
	ChannelInApp    NotificationChannel = "in_app"
	ChannelWebSocket NotificationChannel = "websocket"
)

type NotificationStatus string

const (
	NotificationPending  NotificationStatus = "pending"
	NotificationSent     NotificationStatus = "sent"
	NotificationFailed   NotificationStatus = "failed"
	NotificationDelivered NotificationStatus = "delivered"
)

type Notification struct {
	id         uuid.UUID
	userID     uuid.UUID
	userRole   string
	channel    NotificationChannel
	title      string
	body       string
	data       map[string]interface{}
	status     NotificationStatus
	createdAt  time.Time
	sentAt     *time.Time
}

func NewNotification(
	userID uuid.UUID,
	userRole string,
	channel NotificationChannel,
	title, body string,
	data map[string]interface{},
) *Notification {
	return &Notification{
		id:        uuid.New(),
		userID:    userID,
		userRole:  userRole,
		channel:   channel,
		title:     title,
		body:      body,
		data:      data,
		status:    NotificationPending,
		createdAt: time.Now().UTC(),
	}
}

func (n *Notification) ID() uuid.UUID { return n.id }
func (n *Notification) UserID() uuid.UUID { return n.userID }
func (n *Notification) UserRole() string { return n.userRole }
func (n *Notification) Channel() NotificationChannel { return n.channel }
func (n *Notification) Title() string { return n.title }
func (n *Notification) Body() string { return n.body }
func (n *Notification) Data() map[string]interface{} { return n.data }
func (n *Notification) Status() NotificationStatus { return n.status }
func (n *Notification) CreatedAt() time.Time { return n.createdAt }
func (n *Notification) SentAt() *time.Time { return n.sentAt }

func (n *Notification) MarkSent() {
	now := time.Now().UTC()
	n.sentAt = &now
	n.status = NotificationSent
}

func (n *Notification) MarkFailed() {
	n.status = NotificationFailed
}
