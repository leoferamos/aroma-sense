-- Add refresh token fields to users table
ALTER TABLE users 
ADD COLUMN refresh_token_hash VARCHAR(255),
ADD COLUMN refresh_token_expires_at TIMESTAMP;

-- Add index for faster lookups
CREATE INDEX idx_users_refresh_token_hash ON users(refresh_token_hash);
