-- Remove legacy 'notes' column and update search trigger to use structured notes arrays
ALTER TABLE products DROP COLUMN IF EXISTS notes;

-- Recompute search_vector for existing rows using arrays
UPDATE products SET search_vector =
  setweight(to_tsvector('portuguese', coalesce(unaccent(name), '')), 'A') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(brand), '')), 'B') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(
    coalesce(array_to_string(notes_top, ', '), '') || ' ' ||
    coalesce(array_to_string(notes_heart, ', '), '') || ' ' ||
    coalesce(array_to_string(notes_base, ', '), '')
  ), '')), 'B') ||
  setweight(to_tsvector('portuguese', coalesce(unaccent(description), '')), 'C');

-- Replace trigger function to use arrays
CREATE OR REPLACE FUNCTION products_search_vector_trigger() RETURNS trigger AS $$
begin
  new.search_vector :=
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.name), '')), 'A') ||
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.brand), '')), 'B') ||
    setweight(to_tsvector('portuguese', coalesce(unaccent(
      coalesce(array_to_string(new.notes_top, ', '), '') || ' ' ||
      coalesce(array_to_string(new.notes_heart, ', '), '') || ' ' ||
      coalesce(array_to_string(new.notes_base, ', '), '')
    ), '')), 'B') ||
    setweight(to_tsvector('portuguese', coalesce(unaccent(new.description), '')), 'C');
  return new;
end
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tsvectorupdate ON products;
CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
  ON products FOR EACH ROW EXECUTE PROCEDURE products_search_vector_trigger();
