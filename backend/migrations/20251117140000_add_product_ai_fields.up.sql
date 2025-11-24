-- Add AI-related fields to products and a separate table for embeddings
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS slug VARCHAR(128),
    ADD COLUMN IF NOT EXISTS thumbnail_url VARCHAR(256),
    ADD COLUMN IF NOT EXISTS accords TEXT[],
    ADD COLUMN IF NOT EXISTS occasions TEXT[],
    ADD COLUMN IF NOT EXISTS seasons TEXT[],
    ADD COLUMN IF NOT EXISTS intensity VARCHAR(16),
    ADD COLUMN IF NOT EXISTS gender VARCHAR(16),
    ADD COLUMN IF NOT EXISTS price_range VARCHAR(16),
    ADD COLUMN IF NOT EXISTS notes_top TEXT[],
    ADD COLUMN IF NOT EXISTS notes_heart TEXT[],
    ADD COLUMN IF NOT EXISTS notes_base TEXT[];

-- Unique slug when present
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_slug_unique ON products(slug) WHERE slug IS NOT NULL;

-- Store precomputed embeddings as JSONB
CREATE TABLE IF NOT EXISTS product_embeddings (
    product_id INT PRIMARY KEY REFERENCES products(id) ON DELETE CASCADE,
    embedding JSONB
);
