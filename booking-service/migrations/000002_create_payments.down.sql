-- Удаляем индекс перед таблицей
DROP INDEX IF EXISTS idx_payments_booking;

-- Удаляем таблицу платежей
-- (сначала payments, потому что она ссылается на bookings через FK)
DROP TABLE IF EXISTS payments;
