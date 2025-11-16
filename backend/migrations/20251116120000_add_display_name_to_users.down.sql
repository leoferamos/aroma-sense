-- Remove display_name from users
ALTER TABLE users DROP COLUMN IF EXISTS display_name;
