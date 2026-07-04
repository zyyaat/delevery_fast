-- +migrate Up
-- Create OTPs table (one-time passwords for phone verification)
CREATE TABLE IF NOT EXISTS otps (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone          VARCHAR(20) NOT NULL,
    code           VARCHAR(6) NOT NULL,
    status         VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempts_used  INT NOT NULL DEFAULT 0,
    max_attempts   INT NOT NULL DEFAULT 3,
    expires_at     TIMESTAMPTZ NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_otps_phone ON otps (phone);
CREATE INDEX idx_otps_status ON otps (status);
CREATE INDEX idx_otps_created_at ON otps (created_at DESC);

-- Auto-expire old OTPs after 1 day (cleanup)
-- In production, use a cron job or scheduled task to delete expired OTPs

COMMENT ON TABLE otps IS 'One-time passwords for phone verification';

-- +migrate Down
DROP TABLE IF EXISTS otps;
