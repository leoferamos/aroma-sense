-- Add soft-delete support to carts and cart_items
ALTER TABLE carts
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_carts_deleted_at ON carts(deleted_at);

ALTER TABLE cart_items
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_cart_items_deleted_at ON cart_items(deleted_at);
