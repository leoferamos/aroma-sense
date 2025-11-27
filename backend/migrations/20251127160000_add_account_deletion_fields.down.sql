-- Remove account deletion fields

DROP INDEX IF EXISTS idx_users_deletion_confirmed_at;
DROP INDEX IF EXISTS idx_users_deletion_requested_at;

ALTER TABLE users DROP COLUMN deletion_confirmed_at;
ALTER TABLE users DROP COLUMN deletion_requested_at;
