-- Add LGPD compliance fields to users table
-- deleted_at: Soft delete timestamp
-- last_login_at: Track last login for metrics
-- deactivated_by: Admin who deactivated the user (references public_id)
-- deactivated_at: When user was deactivated

ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMP NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deactivated_by UUID NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS deactivated_at TIMESTAMP NULL;

-- Add index for soft delete queries
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add index for deactivated users
CREATE INDEX IF NOT EXISTS idx_users_deactivated_at ON users(deactivated_at);