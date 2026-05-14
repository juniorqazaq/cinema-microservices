CREATE TABLE halls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    capacity INT NOT NULL CHECK (capacity > 0)
);
