-- +migrate Up
-- Create refresh_tokens table (for JWT refresh token rotation)
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id   UUID NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    token        VARCHAR(128) NOT NULL UNIQUE,
    expires_at   TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    used_at      TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens (token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
CREATE INDEX idx_refresh_tokens_session_id ON refresh_tokens (session_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens (expires_at);
CREATE INDEX idx_refresh_tokens_unused ON refresh_tokens (user_id) WHERE used_at IS NULL;

COMMENT ON TABLE refresh_tokens IS 'Refresh tokens for JWT rotation';
COMMENT ON COLUMN refresh_tokens.token IS 'Hashed refresh token (not the raw token)';
COMMENT ON COLUMN refresh_tokens.used_at IS 'When the token was used (NULL = unused)';

-- +migrate Down
DROP TABLE IF EXISTS refresh_tokens;
