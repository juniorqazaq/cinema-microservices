-- Таблица платежей
CREATE TABLE IF NOT EXISTS payments (
    id         UUID           PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_id UUID           NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    amount     DECIMAL(10, 2) NOT NULL,
    status     VARCHAR(20)    NOT NULL DEFAULT 'pending'
                              CHECK (status IN ('pending', 'paid', 'refunded')),
    paid_at    TIMESTAMPTZ    -- NULL пока платёж не совершён
);

-- Индекс для быстрого поиска платежа по бронированию
CREATE INDEX idx_payments_booking ON payments(booking_id);
