CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    movie_id UUID NOT NULL REFERENCES movies (id) ON DELETE CASCADE,
    hall_id UUID NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    start_time TIMESTAMPTZ NOT NULL,
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0)
);

CREATE INDEX idx_sessions_movie ON sessions (movie_id);
CREATE INDEX idx_sessions_hall ON sessions (hall_id);
CREATE INDEX idx_sessions_start ON sessions (start_time);
