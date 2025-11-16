-- Add display_name to users (nullable)
ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name VARCHAR(64);
