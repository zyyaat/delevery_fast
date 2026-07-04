-- +migrate Up
-- (down migration for users table - drops the table)
DROP TABLE IF EXISTS users;

-- +migrate Down
-- (re-create is handled by the up migration)
