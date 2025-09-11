-- Remove public_id column from users table
ALTER TABLE users DROP COLUMN IF EXISTS public_id;
