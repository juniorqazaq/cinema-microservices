package usecase

import (
	"context"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
)

const maxPageSize = 100

type MovieUseCase struct {
	movies domain.MovieRepository
	cache  domain.MovieCache
}

func NewMovieUseCase(movies domain.MovieRepository, cache domain.MovieCache) *MovieUseCase {
	return &MovieUseCase{movies: movies, cache: cache}
}

func clampPagination(limit, offset int) (int, int, error) {
	if offset < 0 {
		return 0, 0, domain.ErrInvalidPagination
	}
	if limit <= 0 || limit > maxPageSize {
		return 0, 0, domain.ErrInvalidPagination
	}
	return limit, offset, nil
}

func (u *MovieUseCase) CreateMovie(ctx context.Context, isAdmin bool, m *domain.Movie) (*domain.Movie, error) {
	if !isAdmin {
		return nil, domain.ErrForbidden
	}
	if m.Title == "" {
		return nil, domain.ErrInvalidArgument
	}
	if m.Duration <= 0 {
		return nil, domain.ErrInvalidArgument
	}
	created, err := u.movies.Create(ctx, m)
	if err != nil {
		return nil, err
	}
	_ = u.cache.DeleteMovieCache(ctx, created.ID)
	return created, nil
}

func (u *MovieUseCase) GetMovie(ctx context.Context, id string) (*domain.Movie, error) {
	if id == "" {
		return nil, domain.ErrInvalidArgument
	}
	if cached, err := u.cache.GetMovieCache(ctx, id); err != nil {
		return nil, err
	} else if cached != nil {
		return cached, nil
	}
	m, err := u.movies.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := u.cache.SetMovieCache(ctx, m); err != nil {
		return m, nil
	}
	return m, nil
}

func (u *MovieUseCase) UpdateMovie(ctx context.Context, isAdmin bool, m *domain.Movie) (*domain.Movie, error) {
	if !isAdmin {
		return nil, domain.ErrForbidden
	}
	if m.ID == "" {
		return nil, domain.ErrInvalidArgument
	}
	if m.Title == "" || m.Duration <= 0 {
		return nil, domain.ErrInvalidArgument
	}
	updated, err := u.movies.Update(ctx, m)
	if err != nil {
		return nil, err
	}
	_ = u.cache.DeleteMovieCache(ctx, m.ID)
	return updated, nil
}

func (u *MovieUseCase) DeleteMovie(ctx context.Context, isAdmin bool, id string) error {
	if !isAdmin {
		return domain.ErrForbidden
	}
	if id == "" {
		return domain.ErrInvalidArgument
	}
	if err := u.movies.Delete(ctx, id); err != nil {
		return err
	}
	_ = u.cache.DeleteMovieCache(ctx, id)
	return nil
}

func (u *MovieUseCase) ListMovies(ctx context.Context, limit, offset int, genre string, createdAfter, createdBefore *time.Time) ([]*domain.Movie, int64, error) {
	limit, offset, err := clampPagination(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return u.movies.List(ctx, limit, offset, genre, createdAfter, createdBefore)
}

func (u *MovieUseCase) SearchMovies(ctx context.Context, title, genre string, limit, offset int) ([]*domain.Movie, int64, error) {
	limit, offset, err := clampPagination(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if title == "" && genre == "" {
		return nil, 0, domain.ErrInvalidArgument
	}
	return u.movies.Search(ctx, title, genre, limit, offset)
}
