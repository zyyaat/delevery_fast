// Package domain contains the core business logic of the Payment Service.
package domain

import (
        "context"
        "errors"
        "time"

        "github.com/google/uuid"
)

// ============ Errors ============

var (
        ErrPaymentNotFound      = errors.New("payment not found")
        ErrPaymentAlreadyExists = errors.New("payment already exists for this order")
        ErrPaymentDeclined      = errors.New("payment declined by provider")
        ErrPaymentFailed        = errors.New("payment failed")
        ErrRefundExceedsPayment = errors.New("refund amount exceeds payment amount")
        ErrRefundWindowClosed   = errors.New("refund window has closed")
        ErrInvalidAmount        = errors.New("invalid payment amount")
        ErrDuplicateIdempotency = errors.New("duplicate idempotency key")
        ErrProviderUnavailable  = errors.New("payment provider unavailable")
)

// ============ Enums ============

type PaymentMethod string

const (
        PaymentVodafoneCash PaymentMethod = "vodafone_cash"
        PaymentInstaPay     PaymentMethod = "instapay"
        PaymentCard         PaymentMethod = "card"
        PaymentCOD          PaymentMethod = "cod"
)

type PaymentStatus string

const (
        PaymentStatusPending   PaymentStatus = "pending"
        PaymentStatusCaptured  PaymentStatus = "captured"
        PaymentStatusFailed    PaymentStatus = "failed"
        PaymentStatusRefunded  PaymentStatus = "refunded"
        PaymentStatusPartial   PaymentStatus = "partial_refunded"
)

type RefundType string

const (
        RefundFull    RefundType = "full"
        RefundPartial RefundType = "partial"
        RefundGoodwill RefundType = "goodwill"
)

// ============ Entities ============

// Payment represents a payment transaction for an order.
type Payment struct {
        id              uuid.UUID
        orderID         uuid.UUID
        customerID      uuid.UUID
        method          PaymentMethod
        amount          float64
        refundedAmount  float64
        status          PaymentStatus
        providerTxnID   string
        providerName    string
        idempotencyKey  string
        redirectURL     string
        failureReason   string
        refundReason    string
        createdAt       time.Time
        updatedAt       time.Time
        refundedAt      *time.Time
}

// NewPayment creates a new Payment with validation.
func NewPayment(orderID, customerID uuid.UUID, method PaymentMethod, amount float64, idempotencyKey string) (*Payment, error) {
        if amount <= 0 {
                return nil, ErrInvalidAmount
        }
        if idempotencyKey == "" {
                return nil, errors.New("idempotency key is required")
        }

        now := time.Now().UTC()
        return &Payment{
                id:             uuid.New(),
                orderID:        orderID,
                customerID:     customerID,
                method:         method,
                amount:         amount,
                refundedAmount: 0,
                status:         PaymentStatusPending,
                idempotencyKey: idempotencyKey,
                providerName:   string(method),
                createdAt:      now,
                updatedAt:      now,
        }, nil
}

// ============ Getters ============

func (p *Payment) ID() uuid.UUID              { return p.id }
func (p *Payment) OrderID() uuid.UUID          { return p.orderID }
func (p *Payment) CustomerID() uuid.UUID       { return p.customerID }
func (p *Payment) Method() PaymentMethod       { return p.method }
func (p *Payment) Amount() float64             { return p.amount }
func (p *Payment) RefundedAmount() float64     { return p.refundedAmount }
func (p *Payment) Status() PaymentStatus       { return p.status }
func (p *Payment) ProviderTxnID() string       { return p.providerTxnID }
func (p *Payment) ProviderName() string        { return p.providerName }
func (p *Payment) IdempotencyKey() string      { return p.idempotencyKey }
func (p *Payment) RedirectURL() string         { return p.redirectURL }
func (p *Payment) FailureReason() string       { return p.failureReason }
func (p *Payment) RefundReason() string        { return p.refundReason }
func (p *Payment) CreatedAt() time.Time        { return p.createdAt }
func (p *Payment) UpdatedAt() time.Time        { return p.updatedAt }
func (p *Payment) RefundedAt() *time.Time      { return p.refundedAt }

// ============ State Transitions ============

// MarkCaptured marks the payment as successfully captured.
func (p *Payment) MarkCaptured(providerTxnID string) error {
        if p.status != PaymentStatusPending {
                return errors.New("can only capture pending payments")
        }
        p.providerTxnID = providerTxnID
        p.status = PaymentStatusCaptured
        p.updatedAt = time.Now().UTC()
        return nil
}

// MarkFailed marks the payment as failed.
func (p *Payment) MarkFailed(reason string) {
        p.failureReason = reason
        p.status = PaymentStatusFailed
        p.updatedAt = time.Now().UTC()
}

// SetRedirectURL sets the provider redirect URL (for Vodafone Cash / card auth).
func (p *Payment) SetRedirectURL(url string) {
        p.redirectURL = url
        p.updatedAt = time.Now().UTC()
}

// CanRefund returns true if the payment can be refunded.
func (p *Payment) CanRefund() bool {
        return p.status == PaymentStatusCaptured || p.status == PaymentStatusPartial
}

// Refund processes a refund for the given amount.
// For full refunds, status becomes "refunded".
// For partial refunds, status becomes "partial_refunded".
func (p *Payment) Refund(amount float64, reason string) error {
        if !p.CanRefund() {
                return ErrRefundWindowClosed
        }

        if amount <= 0 {
                return ErrInvalidAmount
        }

        totalRefunded := p.refundedAmount + amount
        if totalRefunded > p.amount {
                return ErrRefundExceedsPayment
        }

        now := time.Now().UTC()
        p.refundedAmount = totalRefunded
        p.refundReason = reason
        p.refundedAt = &now

        if totalRefunded >= p.amount {
                p.status = PaymentStatusRefunded
        } else {
                p.status = PaymentStatusPartial
        }

        p.updatedAt = now
        return nil
}

// IsCOD returns true if the payment is cash on delivery.
func (p *Payment) IsCOD() bool {
        return p.method == PaymentCOD
}

// NetAmount returns the amount after refunds (what the platform actually received).
func (p *Payment) NetAmount() float64 {
        return p.amount - p.refundedAmount
}

// ============ Provider Configuration ============

// ProviderConfig holds configuration for a payment provider.
type ProviderConfig struct {
        Name        string
        APIKey      string
        MerchantID  string
        SandboxMode bool
}

// Provider is the interface for payment providers (Vodafone Cash, InstaPay, Paymob, etc.)
type Provider interface {
        Charge(ctx context.Context, payment *Payment) (*ProviderResult, error)
        Refund(ctx context.Context, payment *Payment, amount float64) error
}

// ProviderResult holds the result of a provider charge operation.
type ProviderResult struct {
        TransactionID string
        RedirectURL   string
        Status        PaymentStatus
        ErrorMessage  string
}
