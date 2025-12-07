-- +migrate Down
ALTER TABLE order_items
DROP COLUMN product_slug,
DROP COLUMN product_name,
DROP COLUMN product_image_url;