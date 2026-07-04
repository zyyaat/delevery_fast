-- +migrate Up
DROP TABLE IF EXISTS restaurants;

-- +migrate Down
-- (re-create handled by up migration)
