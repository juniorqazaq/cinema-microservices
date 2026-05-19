-- Legacy booking-service images expected a local seats table.
-- Session-scoped availability is enforced via bookings (migration 000004).
CREATE TABLE IF NOT EXISTS seats (
    id UUID PRIMARY KEY,
    is_available BOOLEAN NOT NULL DEFAULT true
);
