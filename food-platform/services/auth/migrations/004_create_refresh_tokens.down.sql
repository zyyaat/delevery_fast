-- +migrate Up
DROP TABLE IF EXISTS refresh_tokens;

-- +migrate Down
-- (re-create is handled by the up migration)
