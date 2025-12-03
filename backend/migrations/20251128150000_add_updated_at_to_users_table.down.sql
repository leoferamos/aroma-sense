-- +migrate Down
ALTER TABLE users DROP COLUMN updated_at;