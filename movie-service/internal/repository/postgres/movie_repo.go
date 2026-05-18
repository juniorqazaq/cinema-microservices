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

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

type movieScanner interface {
	Scan(dest ...any) error
}

func scanMovie(row movieScanner) (*domain.Movie, error) {
	var m domain.Movie
	age := 12
	if err := row.Scan(&m.ID, &m.Title, &m.Description, &m.Genre, &m.Duration, &m.Rating, &age, &m.CreatedAt); err != nil {
		return nil, err
	}
	m.AgeRating = age
	return &m, nil
}

const createMovieSQL = `
INSERT INTO movies (title, description, genre, duration, rating, age_rating, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title, description, genre, duration, rating, age_rating, created_at`

func (r *MovieRepository) Create(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	now := time.Now().UTC()
	age := m.AgeRating
	if age == 0 {
		age = 12
	}
	row := r.db.QueryRow(ctx, createMovieSQL, m.Title, m.Description, m.Genre, m.Duration, m.Rating, age, now)
	out, err := scanMovie(row)
	if err != nil {
		return nil, fmt.Errorf("postgres: create movie: %w", err)
	}
	return out, nil
}

const getMovieByIDSQL = `
SELECT id, title, description, genre, duration, rating, age_rating, created_at FROM movies WHERE id = $1`

func (r *MovieRepository) GetByID(ctx context.Context, id string) (*domain.Movie, error) {
	row := r.db.QueryRow(ctx, getMovieByIDSQL, id)
	out, err := scanMovie(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMovieNotFound
		}
		return nil, fmt.Errorf("postgres: get movie: %w", err)
	}
	return out, nil
}

const updateMovieSQL = `
UPDATE movies SET title = $2, description = $3, genre = $4, duration = $5, rating = $6, age_rating = $7
WHERE id = $1
RETURNING id, title, description, genre, duration, rating, age_rating, created_at`

func (r *MovieRepository) Update(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	age := m.AgeRating
	if age == 0 {
		age = 12
	}
	row := r.db.QueryRow(ctx, updateMovieSQL, m.ID, m.Title, m.Description, m.Genre, m.Duration, m.Rating, age)
	out, err := scanMovie(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrMovieNotFound
		}
		return nil, fmt.Errorf("postgres: update movie: %w", err)
	}
	return out, nil
}

const deleteMovieSQL = `DELETE FROM movies WHERE id = $1`

func (r *MovieRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, deleteMovieSQL, id)
	if err != nil {
		return fmt.Errorf("postgres: delete movie: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrMovieNotFound
	}
	return nil
}

func (r *MovieRepository) List(ctx context.Context, limit, offset int, genre string, createdAfter, createdBefore *time.Time) ([]*domain.Movie, int64, error) {
	where := "WHERE 1=1"
	args := []any{}
	argPos := 1
	if genre != "" {
		where += fmt.Sprintf(" AND genre ILIKE $%d", argPos)
		args = append(args, "%"+genre+"%")
		argPos++
	}
	if createdAfter != nil {
		where += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *createdAfter)
		argPos++
	}
	if createdBefore != nil {
		where += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *createdBefore)
		argPos++
	}
	countSQL := "SELECT COUNT(*) FROM movies " + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("postgres: count movies: %w", err)
	}
	listSQL := fmt.Sprintf(`
SELECT id, title, description, genre, duration, rating, age_rating, created_at
FROM movies %s
ORDER BY created_at DESC
LIMIT $%d OFFSET $%d`, where, argPos, argPos+1)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: list movies: %w", err)
	}
	defer rows.Close()
	var list []*domain.Movie
	for rows.Next() {
		m, err := scanMovie(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, m)
	}
	return list, total, rows.Err()
}

func (r *MovieRepository) Search(ctx context.Context, title, genre string, limit, offset int) ([]*domain.Movie, int64, error) {
	where := "WHERE 1=1"
	args := []any{}
	argPos := 1
	if title != "" {
		where += fmt.Sprintf(" AND title ILIKE $%d", argPos)
		args = append(args, "%"+title+"%")
		argPos++
	}
	if genre != "" {
		where += fmt.Sprintf(" AND genre ILIKE $%d", argPos)
		args = append(args, "%"+genre+"%")
		argPos++
	}
	countSQL := "SELECT COUNT(*) FROM movies " + where
	var total int64
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("postgres: count search: %w", err)
	}
	listSQL := fmt.Sprintf(`
SELECT id, title, description, genre, duration, rating, age_rating, created_at
FROM movies %s
ORDER BY title ASC
LIMIT $%d OFFSET $%d`, where, argPos, argPos+1)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: search movies: %w", err)
	}
	defer rows.Close()
	var list []*domain.Movie
	for rows.Next() {
		m, err := scanMovie(rows)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, m)
	}
	return list, total, rows.Err()
}
