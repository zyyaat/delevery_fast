-- +migrate Up
CREATE TABLE IF NOT EXISTS fraud_scores (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id      UUID NOT NULL UNIQUE,
    customer_id   UUID NOT NULL,
    score         INT NOT NULL CHECK (score >= 0 AND score <= 100),
    decision      VARCHAR(20) NOT NULL,
    reasons       JSONB NOT NULL DEFAULT '[]',
    model_version VARCHAR(20) NOT NULL DEFAULT 'v1.0',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CHECK (decision IN ('approve', 'review', 'block'))
);

CREATE INDEX idx_fraud_scores_order ON fraud_scores (order_id);
CREATE INDEX idx_fraud_scores_customer ON fraud_scores (customer_id);
CREATE INDEX idx_fraud_scores_decision ON fraud_scores (decision);

CREATE TABLE IF NOT EXISTS trust_scores (
    customer_id      UUID PRIMARY KEY,
    score            INT NOT NULL DEFAULT 50 CHECK (score >= 0 AND score <= 100),
    total_orders     INT NOT NULL DEFAULT 0,
    refund_count     INT NOT NULL DEFAULT 0,
    chargeback_count INT NOT NULL DEFAULT 0,
    fraud_flags      INT NOT NULL DEFAULT 0,
    last_updated     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trust_scores_score ON trust_scores (score);

COMMENT ON TABLE fraud_scores IS 'Fraud risk scores per order';
COMMENT ON TABLE trust_scores IS 'Customer trust scores (0-100, higher = more trustworthy)';

-- +migrate Down
DROP TABLE IF EXISTS trust_scores;
DROP TABLE IF EXISTS fraud_scores;
