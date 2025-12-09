-- +migrate Down
-- Drop soft-delete columns if they exist
ALTER TABLE order_items
  DROP COLUMN IF EXISTS deleted_at;

ALTER TABLE orders
  DROP COLUMN IF EXISTS deleted_at;
