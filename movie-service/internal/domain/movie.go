package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrMovieNotFound      = errors.New("movie not found")
	ErrForbidden          = errors.New("forbidden: admin access required")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrInvalidPagination  = errors.New("invalid pagination")
)

type Movie struct {
	ID          string
	Title       string
	Description string
	Genre       string
	Duration    int
	Rating      float64
	AgeRating   int
	CreatedAt   time.Time
}

type MovieRepository interface {
	Create(ctx context.Context, m *Movie) (*Movie, error)
	GetByID(ctx context.Context, id string) (*Movie, error)
	Update(ctx context.Context, m *Movie) (*Movie, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int, genre string, createdAfter, createdBefore *time.Time) ([]*Movie, int64, error)
	Search(ctx context.Context, title, genre string, limit, offset int) ([]*Movie, int64, error)
}
