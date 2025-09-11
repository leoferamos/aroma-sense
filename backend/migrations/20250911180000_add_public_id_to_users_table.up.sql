-- Add public_id (UUID) to users table for public reference
ALTER TABLE users ADD COLUMN public_id UUID NOT NULL DEFAULT gen_random_uuid();
CREATE UNIQUE INDEX idx_users_public_id ON users(public_id);
