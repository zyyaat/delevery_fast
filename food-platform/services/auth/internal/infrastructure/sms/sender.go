// Package sms provides SMS sending implementations.
package sms

import (
	"context"
	"log/slog"
)

// ============ Interface ============

// Sender is the interface for sending SMS messages.
type Sender interface {
	SendOTP(ctx context.Context, phone, code string) error
}

// ============ Mock Sender (for dev/testing) ============

// MockSender logs the OTP instead of sending it. Use in development only.
type MockSender struct{}

// NewMockSender creates a new MockSender.
func NewMockSender() *MockSender {
	return &MockSender{}
}

// SendOTP logs the OTP code.
func (s *MockSender) SendOTP(ctx context.Context, phone, code string) error {
	slog.Info("mock_sms_otp_sent",
		"phone", phone,
		"code", code,
		"note", "In production, this would send an actual SMS via Twilio",
	)
	return nil
}

// ============ Twilio Sender (TODO: implement) ============

// TwilioConfig holds Twilio credentials.
type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

// TwilioSender sends OTP via Twilio.
// TODO: Implement actual Twilio API calls.
type TwilioSender struct {
	cfg TwilioConfig
}

// NewTwilioSender creates a new TwilioSender.
func NewTwilioSender(cfg TwilioConfig) *TwilioSender {
	return &TwilioSender{cfg: cfg}
}

// SendOTP sends an OTP via Twilio.
func (s *TwilioSender) SendOTP(ctx context.Context, phone, code string) error {
	// TODO: Implement actual Twilio API call
	// For now, log and return nil
	slog.Info("twilio_otp_send_skipped",
		"phone", phone,
		"reason", "Twilio not yet implemented, using mock",
	)
	return nil
}
