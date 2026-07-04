-- +migrate Up
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone        VARCHAR(20) NOT NULL UNIQUE,
    email        VARCHAR(255) UNIQUE,
    name         VARCHAR(255) NOT NULL,
    role         VARCHAR(50) NOT NULL DEFAULT 'customer',
    status       VARCHAR(20) NOT NULL DEFAULT 'active',
    trust_score  INT NOT NULL DEFAULT 50 CHECK (trust_score >= 0 AND trust_score <= 100),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common queries
CREATE INDEX idx_users_phone ON users (phone);
CREATE INDEX idx_users_role ON users (role);
CREATE INDEX idx_users_status ON users (status);
CREATE INDEX idx_users_created_at ON users (created_at DESC);

-- Add comment
COMMENT ON TABLE users IS 'Authenticated users (customers, drivers, restaurants, employees)';
COMMENT ON COLUMN users.phone IS 'Normalized Egyptian phone number (11 digits, starts with 01)';
COMMENT ON COLUMN users.trust_score IS 'Trust score 0-100, used for fraud detection';

-- +migrate Down
DROP TABLE IF EXISTS users;
