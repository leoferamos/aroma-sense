-- Drop protection trigger and column
DROP TRIGGER IF EXISTS trg_prevent_protected_user_changes ON users;
DROP FUNCTION IF EXISTS prevent_protected_user_changes();
ALTER TABLE users DROP COLUMN IF EXISTS is_protected;
