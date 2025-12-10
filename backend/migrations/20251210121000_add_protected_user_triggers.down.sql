-- Remove trigger and function that protect protected users
DROP TRIGGER IF EXISTS trg_prevent_protected_user_changes ON users;
DROP FUNCTION IF EXISTS prevent_protected_user_changes();
