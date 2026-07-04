-- +migrate Up
DROP TABLE IF EXISTS payments;

-- +migrate Down
-- (re-create handled by up migration)
