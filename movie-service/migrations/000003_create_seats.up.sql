CREATE TABLE seats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hall_id UUID NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    row TEXT NOT NULL,
    number INT NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE (hall_id, row, number)
);

CREATE INDEX idx_seats_hall ON seats (hall_id);
