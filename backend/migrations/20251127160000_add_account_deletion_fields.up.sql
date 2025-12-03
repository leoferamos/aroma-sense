-- Add account deletion fields for LGPD compliance
-- These fields implement proper account deletion workflow with retention period

ALTER TABLE users ADD COLUMN deletion_requested_at TIMESTAMP NULL;
ALTER TABLE users ADD COLUMN deletion_confirmed_at TIMESTAMP NULL;

-- Add index for efficient queries on deletion status
CREATE INDEX idx_users_deletion_requested_at ON users(deletion_requested_at);
CREATE INDEX idx_users_deletion_confirmed_at ON users(deletion_confirmed_at);
