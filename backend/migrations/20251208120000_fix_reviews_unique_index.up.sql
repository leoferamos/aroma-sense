-- Fix unique index on reviews to allow re-reviewing after soft delete
-- Drop the old unique index
DROP INDEX IF EXISTS idx_reviews_unique_product_user;

-- Create new unique index that only considers non-deleted reviews
CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_unique_product_user_active ON reviews(product_id, user_id) WHERE deleted_at IS NULL;