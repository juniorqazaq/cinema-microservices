ALTER TABLE halls
    ADD COLUMN IF NOT EXISTS city TEXT NOT NULL DEFAULT 'Astana',
    ADD COLUMN IF NOT EXISTS cinema_name TEXT NOT NULL DEFAULT '';

ALTER TABLE movies
    ADD COLUMN IF NOT EXISTS age_rating INT NOT NULL DEFAULT 12 CHECK (age_rating IN (0, 6, 12, 16, 18));

CREATE INDEX IF NOT EXISTS idx_halls_city ON halls (city);

UPDATE halls SET cinema_name = name WHERE cinema_name = '';
