ALTER TABLE user_contestations
ADD COLUMN IF NOT EXISTS deleted_at timestamptz;
