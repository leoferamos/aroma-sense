-- Add 'role' column to users table
ALTER TABLE users ADD COLUMN role VARCHAR(16) NOT NULL DEFAULT 'client';
