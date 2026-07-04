-- +migrate Up
-- Create webauthn_credentials table (for employee biometric authentication)
-- Reference: FIDO2/WebAuthn standard
CREATE TABLE IF NOT EXISTS webauthn_credentials (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credential_id   TEXT NOT NULL UNIQUE,
    public_key      BYTEA NOT NULL,
    label           VARCHAR(100),
    sign_count      BIGINT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at    TIMESTAMPTZ
);

-- Indexes
CREATE INDEX idx_webauthn_user_id ON webauthn_credentials (user_id);
CREATE INDEX idx_webauthn_credential_id ON webauthn_credentials (credential_id);

COMMENT ON TABLE webauthn_credentials IS 'WebAuthn/FIDO2 credentials for biometric auth';
COMMENT ON COLUMN webauthn_credentials.credential_id IS 'Base64URL-encoded credential ID';
COMMENT ON COLUMN webauthn_credentials.public_key IS 'Public key from authenticator';
COMMENT ON COLUMN webauthn_credentials.sign_count IS 'Signature counter for replay attack prevention';

-- +migrate Down
DROP TABLE IF EXISTS webauthn_credentials;
