-- Remove refresh token fields from users table
DROP INDEX IF EXISTS idx_users_refresh_token_hash;

ALTER TABLE users 
DROP COLUMN IF EXISTS refresh_token_hash,
DROP COLUMN IF EXISTS refresh_token_expires_at;
