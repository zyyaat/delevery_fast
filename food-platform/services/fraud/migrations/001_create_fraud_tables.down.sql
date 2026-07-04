-- +migrate Up
DROP TABLE IF EXISTS trust_scores;
DROP TABLE IF EXISTS fraud_scores;

-- +migrate Down
-- (re-create handled by up migration)
