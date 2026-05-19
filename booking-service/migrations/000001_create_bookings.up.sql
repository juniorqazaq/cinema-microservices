-- Расширение для генерации UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица бронирований
CREATE TABLE IF NOT EXISTS bookings (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID        NOT NULL,
    session_id UUID        NOT NULL,
    seat_id    UUID        NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending'
                           CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Индекс для быстрого поиска броней по пользователю
CREATE INDEX IF NOT EXISTS idx_bookings_user    ON bookings(user_id);

-- Индекс для быстрого поиска броней по сеансу
CREATE INDEX IF NOT EXISTS idx_bookings_session ON bookings(session_id);
