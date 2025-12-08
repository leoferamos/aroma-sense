-- Revert unique index fix
-- Drop the new unique index
DROP INDEX IF EXISTS idx_reviews_unique_product_user_active;

-- Recreate the old unique index
CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_unique_product_user ON reviews(product_id, user_id);