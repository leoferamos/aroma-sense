-- +migrate Up
-- Ensure soft-delete columns exist for orders and order_items
ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;

ALTER TABLE order_items
  ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP WITH TIME ZONE;
