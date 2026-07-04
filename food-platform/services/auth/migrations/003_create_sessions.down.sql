-- +migrate Up
DROP TABLE IF EXISTS sessions;

-- +migrate Down
-- (re-create is handled by the up migration)
