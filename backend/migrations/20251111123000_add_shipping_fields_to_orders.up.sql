-- Add shipping-related fields to orders
ALTER TABLE orders
  ADD COLUMN IF NOT EXISTS shipping_price DECIMAL(10,2) NOT NULL DEFAULT 0,
  ADD COLUMN IF NOT EXISTS shipping_carrier VARCHAR(100),
  ADD COLUMN IF NOT EXISTS shipping_service_code VARCHAR(100),
  ADD COLUMN IF NOT EXISTS shipping_estimated_delivery TIMESTAMP WITH TIME ZONE,
  ADD COLUMN IF NOT EXISTS shipping_tracking VARCHAR(255),
  ADD COLUMN IF NOT EXISTS shipping_status VARCHAR(50);

-- Helpful indexes for lookups/filters
CREATE INDEX IF NOT EXISTS idx_orders_shipping_status ON orders(shipping_status);
CREATE INDEX IF NOT EXISTS idx_orders_shipping_estimated_delivery ON orders(shipping_estimated_delivery);
