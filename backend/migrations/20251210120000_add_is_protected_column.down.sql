-- Drop is_protected flag from users table
ALTER TABLE users DROP COLUMN IF EXISTS is_protected;
