-- Remove soft-delete columns from carts and cart_items
DROP INDEX IF EXISTS idx_cart_items_deleted_at;
ALTER TABLE cart_items DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_carts_deleted_at;
ALTER TABLE carts DROP COLUMN IF EXISTS deleted_at;
