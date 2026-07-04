-- +migrate Up
DROP TABLE IF EXISTS webauthn_credentials;

-- +migrate Down
-- (re-create is handled by the up migration)
