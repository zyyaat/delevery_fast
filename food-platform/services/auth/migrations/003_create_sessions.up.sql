-- +migrate Up
-- Create sessions table (active user sessions)
CREATE TABLE IF NOT EXISTS sessions (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token       VARCHAR(128) NOT NULL UNIQUE,
    device_fingerprint  VARCHAR(255),
    user_agent          TEXT,
    ip_address          VARCHAR(45),
    expires_at          TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at          TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_sessions_user_id ON sessions (user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions (refresh_token);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);
CREATE INDEX idx_sessions_revoked_at ON sessions (revoked_at);
CREATE INDEX idx_sessions_user_active ON sessions (user_id) WHERE revoked_at IS NULL;

COMMENT ON TABLE sessions IS 'Active and revoked user sessions';
COMMENT ON COLUMN sessions.refresh_token IS 'Hashed refresh token (not the raw token)';

-- +migrate Down
DROP TABLE IF EXISTS sessions;
