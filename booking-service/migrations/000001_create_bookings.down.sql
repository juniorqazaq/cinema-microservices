-- Удаляем индексы перед таблицей
DROP INDEX IF EXISTS idx_bookings_session;
DROP INDEX IF EXISTS idx_bookings_user;

-- Удаляем таблицу бронирований
DROP TABLE IF EXISTS bookings;
