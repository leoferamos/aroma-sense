-- Drop the trigger, function, index and column added in the up migration
DROP TRIGGER IF EXISTS tsvectorupdate ON products;
DROP FUNCTION IF EXISTS products_search_vector_trigger();
DROP INDEX IF EXISTS idx_products_search_vector;
ALTER TABLE products DROP COLUMN IF EXISTS search_vector;
