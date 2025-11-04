-- Add search_vector column and GIN index for full-text search on products
ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector tsvector;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Populate existing rows
UPDATE products SET search_vector =
  setweight(to_tsvector('portuguese', coalesce(unaccent(name), '')), 'A') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(brand), '')), 'B') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(notes), '')), 'B') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(description), '')), 'C');

-- Create GIN index for fast search
CREATE INDEX IF NOT EXISTS idx_products_search_vector ON products USING GIN(search_vector);

-- Trigger function to keep search_vector up-to-date
CREATE OR REPLACE FUNCTION products_search_vector_trigger() RETURNS trigger AS $$
begin
  new.search_vector :=
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.name), '')), 'A') ||
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.brand), '')), 'B') ||
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.notes), '')), 'B') ||
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.description), '')), 'C');
  return new;
end
$$ LANGUAGE plpgsql;

-- Ensure trigger is not duplicated when running migration again
DROP TRIGGER IF EXISTS tsvectorupdate ON products;
CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
  ON products FOR EACH ROW EXECUTE PROCEDURE products_search_vector_trigger();
