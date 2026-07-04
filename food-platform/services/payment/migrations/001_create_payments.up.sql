-- +migrate Up
-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id         UUID NOT NULL,
    customer_id      UUID NOT NULL,
    method           VARCHAR(20) NOT NULL,
    amount           DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    refunded_amount  DECIMAL(10,2) NOT NULL DEFAULT 0,
    status           VARCHAR(20) NOT NULL DEFAULT 'pending',
    provider_txn_id  VARCHAR(100),
    provider_name    VARCHAR(50),
    idempotency_key  VARCHAR(255) NOT NULL UNIQUE,
    redirect_url     TEXT,
    failure_reason   TEXT,
    refund_reason    TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    refunded_at      TIMESTAMPTZ,

    CHECK (method IN ('vodafone_cash', 'instapay', 'card', 'cod')),
    CHECK (status IN ('pending', 'captured', 'failed', 'refunded', 'partial_refunded')),
    CHECK (refunded_amount >= 0 AND refunded_amount <= amount)
);

CREATE INDEX idx_payments_order_id ON payments (order_id);
CREATE INDEX idx_payments_customer ON payments (customer_id);
CREATE INDEX idx_payments_status ON payments (status);
CREATE INDEX idx_payments_idempotency ON payments (idempotency_key);
CREATE INDEX idx_payments_created_at ON payments (created_at DESC);

COMMENT ON TABLE payments IS 'Payment transactions for orders';
COMMENT ON COLUMN payments.idempotency_key IS 'Unique key to prevent duplicate charges (X-Idempotency-Key header)';

-- +migrate Down
DROP TABLE IF EXISTS payments;
