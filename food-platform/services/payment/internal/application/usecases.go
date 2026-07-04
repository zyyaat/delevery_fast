// Package application contains use cases for the Payment Service.
package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/food-platform/services/payment/internal/domain"
	"github.com/google/uuid"
)

// ============ Ports ============

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error)
	FindByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}

type ProviderFactory interface {
	GetProvider(method domain.PaymentMethod) (domain.Provider, error)
}

type EventPublisher interface {
	PublishPaymentCaptured(ctx context.Context, event PaymentCapturedEvent) error
	PublishPaymentFailed(ctx context.Context, event PaymentFailedEvent) error
	PublishPaymentRefunded(ctx context.Context, event PaymentRefundedEvent) error
}

// ============ DTOs ============

type ChargePaymentCommand struct {
	OrderID        uuid.UUID
	CustomerID     uuid.UUID
	Method         domain.PaymentMethod
	Amount         float64
	IdempotencyKey string
}

type PaymentDTO struct {
	ID             uuid.UUID `json:"id"`
	OrderID        uuid.UUID `json:"order_id"`
	Method         string    `json:"method"`
	Amount         float64   `json:"amount"`
	Status         string    `json:"status"`
	ProviderTxnID  string    `json:"provider_txn_id,omitempty"`
	RedirectURL    string    `json:"provider_redirect_url,omitempty"`
	FailureReason  string    `json:"failure_reason,omitempty"`
	CreatedAt      string    `json:"created_at"`
}

type RefundCommand struct {
	PaymentID uuid.UUID
	Amount    float64
	Reason    string
	Type      domain.RefundType
}

// ============ Events ============

type PaymentCapturedEvent struct {
	EventID       uuid.UUID
	PaymentID     uuid.UUID
	OrderID       uuid.UUID
	Amount        float64
	Method        string
	ProviderTxnID string
}

type PaymentFailedEvent struct {
	EventID   uuid.UUID
	PaymentID uuid.UUID
	OrderID   uuid.UUID
	Reason    string
	Retryable bool
}

type PaymentRefundedEvent struct {
	EventID   uuid.UUID
	PaymentID uuid.UUID
	OrderID   uuid.UUID
	Amount    float64
	Reason    string
}

// ============ Use Cases ============

// ChargePaymentUseCase processes a payment charge.
type ChargePaymentUseCase struct {
	repo     PaymentRepository
	provider ProviderFactory
	publisher EventPublisher
}

func NewChargePaymentUseCase(repo PaymentRepository, pf ProviderFactory, pub EventPublisher) *ChargePaymentUseCase {
	return &ChargePaymentUseCase{repo: repo, provider: pf, publisher: pub}
}

func (uc *ChargePaymentUseCase) Execute(ctx context.Context, cmd ChargePaymentCommand) (*PaymentDTO, error) {
	// Check idempotency — if a payment with this key exists, return it
	existing, err := uc.repo.FindByIdempotencyKey(ctx, cmd.IdempotencyKey)
	if err == nil && existing != nil {
		// Return the existing payment (idempotent response)
		return toDTO(existing), nil
	}

	// Create new payment
	payment, err := domain.NewPayment(cmd.OrderID, cmd.CustomerID, cmd.Method, cmd.Amount, cmd.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	// For COD, mark as captured immediately (no provider call)
	if payment.IsCOD() {
		payment.MarkCaptured("COD-" + payment.ID().String()[:8])
		if err := uc.repo.Create(ctx, payment); err != nil {
			return nil, fmt.Errorf("create payment: %w", err)
		}

		_ = uc.publisher.PublishPaymentCaptured(ctx, PaymentCapturedEvent{
			EventID:       uuid.New(),
			PaymentID:     payment.ID(),
			OrderID:       payment.OrderID(),
			Amount:        payment.Amount(),
			Method:        string(payment.Method()),
			ProviderTxnID: payment.ProviderTxnID(),
		})

		return toDTO(payment), nil
	}

	// Get provider
	provider, err := uc.provider.GetProvider(cmd.Method)
	if err != nil {
		return nil, err
	}

	// Charge via provider
	result, err := provider.Charge(ctx, payment)
	if err != nil {
		payment.MarkFailed(err.Error())
		_ = uc.repo.Create(ctx, payment)
		_ = uc.publisher.PublishPaymentFailed(ctx, PaymentFailedEvent{
			EventID:   uuid.New(),
			PaymentID: payment.ID(),
			OrderID:   payment.OrderID(),
			Reason:    err.Error(),
			Retryable: true,
		})
		return nil, domain.ErrPaymentFailed
	}

	// Update payment with result
	if result.RedirectURL != "" {
		payment.SetRedirectURL(result.RedirectURL)
	}

	if result.Status == domain.PaymentStatusCaptured {
		payment.MarkCaptured(result.TransactionID)
	} else if result.Status == domain.PaymentStatusFailed {
		payment.MarkFailed(result.ErrorMessage)
	}

	// Save payment
	if err := uc.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	// Publish event
	if payment.Status() == domain.PaymentStatusCaptured {
		_ = uc.publisher.PublishPaymentCaptured(ctx, PaymentCapturedEvent{
			EventID:       uuid.New(),
			PaymentID:     payment.ID(),
			OrderID:       payment.OrderID(),
			Amount:        payment.Amount(),
			Method:        string(payment.Method()),
			ProviderTxnID: payment.ProviderTxnID(),
		})
	} else if payment.Status() == domain.PaymentStatusFailed {
		_ = uc.publisher.PublishPaymentFailed(ctx, PaymentFailedEvent{
			EventID:   uuid.New(),
			PaymentID: payment.ID(),
			OrderID:   payment.OrderID(),
			Reason:    payment.FailureReason(),
			Retryable: true,
		})
	}

	return toDTO(payment), nil
}

// GetPaymentUseCase retrieves a payment by ID.
type GetPaymentUseCase struct {
	repo PaymentRepository
}

func NewGetPaymentUseCase(repo PaymentRepository) *GetPaymentUseCase {
	return &GetPaymentUseCase{repo: repo}
}

func (uc *GetPaymentUseCase) Execute(ctx context.Context, id uuid.UUID) (*PaymentDTO, error) {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toDTO(p), nil
}

// GetPaymentByOrderUseCase retrieves a payment by order ID.
type GetPaymentByOrderUseCase struct {
	repo PaymentRepository
}

func NewGetPaymentByOrderUseCase(repo PaymentRepository) *GetPaymentByOrderUseCase {
	return &GetPaymentByOrderUseCase{repo: repo}
}

func (uc *GetPaymentByOrderUseCase) Execute(ctx context.Context, orderID uuid.UUID) (*PaymentDTO, error) {
	p, err := uc.repo.FindByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return toDTO(p), nil
}

// RefundPaymentUseCase processes a refund.
type RefundPaymentUseCase struct {
	repo      PaymentRepository
	provider  ProviderFactory
	publisher EventPublisher
}

func NewRefundPaymentUseCase(repo PaymentRepository, pf ProviderFactory, pub EventPublisher) *RefundPaymentUseCase {
	return &RefundPaymentUseCase{repo: repo, provider: pf, publisher: pub}
}

func (uc *RefundPaymentUseCase) Execute(ctx context.Context, cmd RefundCommand) (*PaymentDTO, error) {
	payment, err := uc.repo.FindByID(ctx, cmd.PaymentID)
	if err != nil {
		return nil, err
	}

	// Calculate refund amount
	refundAmount := cmd.Amount
	if cmd.Type == domain.RefundFull {
		refundAmount = payment.Amount() - payment.RefundedAmount()
	}

	// For non-COD, call provider to process refund
	if !payment.IsCOD() {
		provider, err := uc.provider.GetProvider(payment.Method())
		if err != nil {
			slog.Error("refund_provider_error", "error", err)
			// Continue anyway — mark as refunded in our system
		} else {
			if err := provider.Refund(ctx, payment, refundAmount); err != nil {
				slog.Error("refund_failed", "payment_id", payment.ID(), "error", err)
				// Continue — we'll mark as refunded in our system and handle reconciliation later
			}
		}
	}

	// Apply refund to payment
	if err := payment.Refund(refundAmount, cmd.Reason); err != nil {
		return nil, err
	}

	// Save
	if err := uc.repo.Update(ctx, payment); err != nil {
		return nil, fmt.Errorf("update payment: %w", err)
	}

	// Publish event
	_ = uc.publisher.PublishPaymentRefunded(ctx, PaymentRefundedEvent{
		EventID:   uuid.New(),
		PaymentID: payment.ID(),
		OrderID:   payment.OrderID(),
		Amount:    refundAmount,
		Reason:    cmd.Reason,
	})

	return toDTO(payment), nil
}

// ============ Helpers ============

func toDTO(p *domain.Payment) *PaymentDTO {
	return &PaymentDTO{
		ID:            p.ID(),
		OrderID:       p.OrderID(),
		Method:        string(p.Method()),
		Amount:        p.Amount(),
		Status:        string(p.Status()),
		ProviderTxnID: p.ProviderTxnID(),
		RedirectURL:   p.RedirectURL(),
		FailureReason: p.FailureReason(),
		CreatedAt:     p.CreatedAt().Format("2006-01-02T15:04:05Z"),
	}
}
