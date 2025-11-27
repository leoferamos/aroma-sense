-- Remove LGPD compliance fields from users table

DROP INDEX IF EXISTS idx_users_deactivated_at;
DROP INDEX IF EXISTS idx_users_deleted_at;

ALTER TABLE users DROP COLUMN IF EXISTS deactivated_at;
ALTER TABLE users DROP COLUMN IF EXISTS deactivated_by;
ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;