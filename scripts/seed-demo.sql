-- Demo seed for movie_db: halls, seats (10 rows x 15 seats), sessions for 7 days.
-- Run against movie_db after movies exist: psql $MOVIE_DATABASE_URL -f scripts/seed-demo.sql

-- Astana cinemas
INSERT INTO halls (name, capacity, city, cinema_name)
VALUES
    ('Hall 1', 150, 'Astana', 'Keruen Cineplex'),
    ('Hall 2', 120, 'Astana', 'Keruen Cineplex'),
    ('IMAX', 150, 'Astana', 'Khan Shatyr IMAX'),
    ('Hall 1', 150, 'Almaty', 'Kinopark 8'),
    ('Hall 2', 120, 'Almaty', 'Kinopark 8')
ON CONFLICT DO NOTHING;

-- Seats for halls that have none yet (A–J rows, 15 seats each, VIP = rows A–B)
INSERT INTO seats (hall_id, row, number, is_available)
SELECT h.id, chr(64 + r.row_num), s.seat_num, true
FROM halls h
CROSS JOIN generate_series(1, 10) AS r(row_num)
CROSS JOIN generate_series(1, 15) AS s(seat_num)
WHERE NOT EXISTS (SELECT 1 FROM seats st WHERE st.hall_id = h.id)
ON CONFLICT DO NOTHING;

-- Sessions: today through 7 days, every 3 hours (UTC stored; display in Asia/Almaty on client)
INSERT INTO sessions (movie_id, hall_id, start_time, price)
SELECT m.id, h.id, gs, 2500
FROM movies m
CROSS JOIN halls h
CROSS JOIN generate_series(
    date_trunc('hour', NOW()),
    NOW() + interval '7 days',
    interval '3 hours'
) AS gs
WHERE NOT EXISTS (
    SELECT 1 FROM sessions s
    WHERE s.movie_id = m.id AND s.hall_id = h.id AND s.start_time = gs
);
