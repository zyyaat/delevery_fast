// Package postgres implements the Payment repository.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/food-platform/services/payment/internal/domain"
	"github.com/google/uuid"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	query := `
		INSERT INTO payments (
			id, order_id, customer_id, method, amount, refunded_amount,
			status, provider_txn_id, provider_name, idempotency_key,
			redirect_url, failure_reason, refund_reason, created_at, updated_at, refunded_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	var refundedAt interface{}
	if p.RefundedAt() != nil {
		refundedAt = *p.RefundedAt()
	}

	_, err := r.db.ExecContext(ctx, query,
		p.ID(), p.OrderID(), p.CustomerID(),
		string(p.Method()), p.Amount(), p.RefundedAmount(),
		string(p.Status()), p.ProviderTxnID(), p.ProviderName(), p.IdempotencyKey(),
		p.RedirectURL(), p.FailureReason(), p.RefundReason(),
		p.CreatedAt(), p.UpdatedAt(), refundedAt,
	)
	if err != nil {
		return fmt.Errorf("payment_repo.Create: %w", err)
	}
	return nil
}

func (r *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, customer_id, method, amount, refunded_amount,
		       status, COALESCE(provider_txn_id, ''), COALESCE(provider_name, ''),
		       idempotency_key, COALESCE(redirect_url, ''),
		       COALESCE(failure_reason, ''), COALESCE(refund_reason, ''),
		       created_at, updated_at, refunded_at
		FROM payments WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	return scanPayment(row)
}

func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, customer_id, method, amount, refunded_amount,
		       status, COALESCE(provider_txn_id, ''), COALESCE(provider_name, ''),
		       idempotency_key, COALESCE(redirect_url, ''),
		       COALESCE(failure_reason, ''), COALESCE(refund_reason, ''),
		       created_at, updated_at, refunded_at
		FROM payments WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, orderID)
	return scanPayment(row)
}

func (r *PaymentRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Payment, error) {
	query := `
		SELECT id, order_id, customer_id, method, amount, refunded_amount,
		       status, COALESCE(provider_txn_id, ''), COALESCE(provider_name, ''),
		       idempotency_key, COALESCE(redirect_url, ''),
		       COALESCE(failure_reason, ''), COALESCE(refund_reason, ''),
		       created_at, updated_at, refunded_at
		FROM payments WHERE idempotency_key = $1
	`
	row := r.db.QueryRowContext(ctx, query, key)
	return scanPayment(row)
}

func (r *PaymentRepository) Update(ctx context.Context, p *domain.Payment) error {
	query := `
		UPDATE payments SET
			refunded_amount = $2, status = $3, provider_txn_id = $4,
			failure_reason = $5, refund_reason = $6, updated_at = $7, refunded_at = $8
		WHERE id = $1
	`
	var refundedAt interface{}
	if p.RefundedAt() != nil {
		refundedAt = *p.RefundedAt()
	}

	_, err := r.db.ExecContext(ctx, query,
		p.ID(), p.RefundedAmount(), string(p.Status()), p.ProviderTxnID(),
		p.FailureReason(), p.RefundReason(), time.Now().UTC(), refundedAt,
	)
	return err
}

// ============ Helpers ============

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanPayment(s scanner) (*domain.Payment, error) {
	var (
		id              uuid.UUID
		orderID         uuid.UUID
		customerID      uuid.UUID
		method          string
		amount          float64
		refundedAmount  float64
		status          string
		providerTxnID   string
		providerName    string
		idempotencyKey  string
		redirectURL     string
		failureReason   string
		refundReason    string
		createdAt       time.Time
		updatedAt       time.Time
		refundedAt      *time.Time
	)

	err := s.Scan(
		&id, &orderID, &customerID, &method, &amount, &refundedAmount,
		&status, &providerTxnID, &providerName, &idempotencyKey,
		&redirectURL, &failureReason, &refundReason,
		&createdAt, &updatedAt, &refundedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, fmt.Errorf("scanPayment: %w", err)
	}

	return &domain.Payment{
		id:             id,
		orderID:        orderID,
		customerID:     customerID,
		method:         domain.PaymentMethod(method),
		amount:         amount,
		refundedAmount: refundedAmount,
		status:         domain.PaymentStatus(status),
		providerTxnID:  providerTxnID,
		providerName:   providerName,
		idempotencyKey: idempotencyKey,
		redirectURL:    redirectURL,
		failureReason:  failureReason,
		refundReason:   refundReason,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
		refundedAt:     refundedAt,
	}, nil
}
