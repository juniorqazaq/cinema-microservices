CREATE UNIQUE INDEX IF NOT EXISTS idx_bookings_session_seat_active
    ON bookings (session_id, seat_id)
    WHERE status IN ('pending', 'confirmed');
