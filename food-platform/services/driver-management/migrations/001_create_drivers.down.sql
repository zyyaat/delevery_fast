-- +migrate Up
DROP TABLE IF EXISTS drivers;

-- +migrate Down
-- (re-create handled by up migration)
