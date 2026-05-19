package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{db: db}
}

type sessionScanner interface {
	Scan(dest ...any) error
}

func scanSession(row sessionScanner) (*domain.Session, error) {
	var s domain.Session
	if err := row.Scan(&s.ID, &s.MovieID, &s.HallID, &s.StartTime, &s.Price); err != nil {
		return nil, err
	}
	return &s, nil
}

const createSessionSQL = `
INSERT INTO sessions (movie_id, hall_id, start_time, price)
VALUES ($1, $2, $3, $4)
RETURNING id, movie_id, hall_id, start_time, price::float8`

func (r *SessionRepository) Create(ctx context.Context, s *domain.Session) (*domain.Session, error) {
	row := r.db.QueryRow(ctx, createSessionSQL, s.MovieID, s.HallID, s.StartTime.UTC(), s.Price)
	out, err := scanSession(row)
	if err != nil {
		if isForeignKeyViolation(err) {
			return nil, domain.ErrInvalidArgument
		}
		return nil, fmt.Errorf("postgres: create session: %w", err)
	}
	return out, nil
}

const getSessionByIDSQL = `
SELECT id, movie_id, hall_id, start_time, price::float8 FROM sessions WHERE id = $1`

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	row := r.db.QueryRow(ctx, getSessionByIDSQL, id)
	out, err := scanSession(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("postgres: get session: %w", err)
	}
	return out, nil
}

func (r *SessionRepository) List(ctx context.Context, movieID, city string, date *time.Time, limit, offset int) ([]*domain.Session, int64, error) {
	from := "FROM sessions s"
	where := "WHERE 1=1"
	args := []any{}
	argPos := 1
	if city != "" {
		from += " JOIN halls h ON h.id = s.hall_id"
		where += fmt.Sprintf(" AND h.city = $%d", argPos)
		args = append(args, city)
		argPos++
	}
	if movieID != "" {
		where += fmt.Sprintf(" AND s.movie_id = $%d::uuid", argPos)
		args = append(args, movieID)
		argPos++
	}
	if date != nil {
		start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		where += fmt.Sprintf(" AND s.start_time >= $%d AND s.start_time < $%d", argPos, argPos+1)
		args = append(args, start, end)
		argPos += 2
	}
	countSQL := "SELECT COUNT(*) " + from + " " + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("postgres: count sessions: %w", err)
	}
	listSQL := fmt.Sprintf(`
SELECT s.id, s.movie_id, s.hall_id, s.start_time, s.price::float8 %s %s
ORDER BY s.start_time ASC
LIMIT $%d OFFSET $%d`, from, where, argPos, argPos+1)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: list sessions: %w", err)
	}
	defer rows.Close()
	var list []*domain.Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, s)
	}
	return list, total, rows.Err()
}

const getSessionsByDateSQL = `
SELECT id, movie_id, hall_id, start_time, price::float8 FROM sessions
WHERE start_time >= $1 AND start_time < $2
ORDER BY start_time ASC`

func (r *SessionRepository) GetByDate(ctx context.Context, day time.Time) ([]*domain.Session, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	rows, err := r.db.Query(ctx, getSessionsByDateSQL, start, end)
	if err != nil {
		return nil, fmt.Errorf("postgres: sessions by date: %w", err)
	}
	defer rows.Close()
	var list []*domain.Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

const getSessionsByMovieSQL = `
SELECT id, movie_id, hall_id, start_time, price::float8 FROM sessions WHERE movie_id = $1 ORDER BY start_time ASC`

func (r *SessionRepository) GetByMovieID(ctx context.Context, movieID string) ([]*domain.Session, error) {
	rows, err := r.db.Query(ctx, getSessionsByMovieSQL, movieID)
	if err != nil {
		return nil, fmt.Errorf("postgres: sessions by movie: %w", err)
	}
	defer rows.Close()
	var list []*domain.Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
