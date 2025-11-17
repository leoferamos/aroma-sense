-- Rollback AI-related fields and embeddings table
DROP TABLE IF EXISTS product_embeddings;

DROP INDEX IF EXISTS idx_products_slug_unique;

ALTER TABLE products
    DROP COLUMN IF EXISTS notes_base,
    DROP COLUMN IF EXISTS notes_heart,
    DROP COLUMN IF EXISTS notes_top,
    DROP COLUMN IF EXISTS price_range,
    DROP COLUMN IF EXISTS gender,
    DROP COLUMN IF EXISTS intensity,
    DROP COLUMN IF EXISTS seasons,
    DROP COLUMN IF EXISTS occasions,
    DROP COLUMN IF EXISTS accords,
    DROP COLUMN IF EXISTS thumbnail_url,
    DROP COLUMN IF EXISTS slug;
