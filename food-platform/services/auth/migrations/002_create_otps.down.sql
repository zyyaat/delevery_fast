-- +migrate Up
DROP TABLE IF EXISTS otps;

-- +migrate Down
-- (re-create is handled by the up migration)
