-- Re-add 'notes' column and restore original search trigger using notes
ALTER TABLE products ADD COLUMN IF NOT EXISTS notes TEXT;

-- Recompute search_vector for existing rows using 'notes'
UPDATE products SET search_vector =
  setweight(to_tsvector('portuguese', coalesce(unaccent(name), '' )), 'A') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(brand), '' )), 'B') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(notes), '' )), 'B') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(description), '' )), 'C');

-- Restore trigger function using notes field
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

DROP TRIGGER IF EXISTS tsvectorupdate ON products;
CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
  ON products FOR EACH ROW EXECUTE PROCEDURE products_search_vector_trigger();
