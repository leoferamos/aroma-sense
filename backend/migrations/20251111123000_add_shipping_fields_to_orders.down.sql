-- Revert shipping-related fields from orders
DROP INDEX IF EXISTS idx_orders_shipping_status;
DROP INDEX IF EXISTS idx_orders_shipping_estimated_delivery;

ALTER TABLE orders
  DROP COLUMN IF EXISTS shipping_status,
  DROP COLUMN IF EXISTS shipping_tracking,
  DROP COLUMN IF EXISTS shipping_estimated_delivery,
  DROP COLUMN IF EXISTS shipping_service_code,
  DROP COLUMN IF EXISTS shipping_carrier,
  DROP COLUMN IF EXISTS shipping_price;
