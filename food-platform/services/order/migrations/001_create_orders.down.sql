-- +migrate Up
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;

-- +migrate Down
-- (re-create handled by up migration)
