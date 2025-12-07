-- +migrate Up
ALTER TABLE order_items
ADD COLUMN product_slug VARCHAR(255),
ADD COLUMN product_name VARCHAR(255),
ADD COLUMN product_image_url VARCHAR(500);