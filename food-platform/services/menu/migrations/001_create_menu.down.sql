-- +migrate Up
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS menu_categories;

-- +migrate Down
-- (re-create handled by up migration)
