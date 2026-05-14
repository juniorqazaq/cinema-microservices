package testutil

import (
	"context"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
)

type FakeMovieRepo struct {
	CreateFn func(ctx context.Context, m *domain.Movie) (*domain.Movie, error)
	GetFn    func(ctx context.Context, id string) (*domain.Movie, error)
	UpdateFn func(ctx context.Context, m *domain.Movie) (*domain.Movie, error)
	DeleteFn func(ctx context.Context, id string) error
	ListFn   func(ctx context.Context, limit, offset int, genre string, ca, cb *time.Time) ([]*domain.Movie, int64, error)
	SearchFn func(ctx context.Context, title, genre string, limit, offset int) ([]*domain.Movie, int64, error)
}

func (f *FakeMovieRepo) Create(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	if f.CreateFn != nil {
		return f.CreateFn(ctx, m)
	}
	return nil, domain.ErrInvalidArgument
}

func (f *FakeMovieRepo) GetByID(ctx context.Context, id string) (*domain.Movie, error) {
	if f.GetFn != nil {
		return f.GetFn(ctx, id)
	}
	return nil, domain.ErrMovieNotFound
}

func (f *FakeMovieRepo) Update(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
	if f.UpdateFn != nil {
		return f.UpdateFn(ctx, m)
	}
	return nil, domain.ErrMovieNotFound
}

func (f *FakeMovieRepo) Delete(ctx context.Context, id string) error {
	if f.DeleteFn != nil {
		return f.DeleteFn(ctx, id)
	}
	return domain.ErrMovieNotFound
}

func (f *FakeMovieRepo) List(ctx context.Context, limit, offset int, genre string, createdAfter, createdBefore *time.Time) ([]*domain.Movie, int64, error) {
	if f.ListFn != nil {
		return f.ListFn(ctx, limit, offset, genre, createdAfter, createdBefore)
	}
	return nil, 0, nil
}

func (f *FakeMovieRepo) Search(ctx context.Context, title, genre string, limit, offset int) ([]*domain.Movie, int64, error) {
	if f.SearchFn != nil {
		return f.SearchFn(ctx, title, genre, limit, offset)
	}
	return nil, 0, nil
}

type FakeMovieCache struct {
	SetMovieFn            func(ctx context.Context, m *domain.Movie) error
	GetMovieFn            func(ctx context.Context, id string) (*domain.Movie, error)
	DeleteMovieFn         func(ctx context.Context, id string) error
	SetSessionsByDateFn   func(ctx context.Context, date string, sessions []*domain.Session) error
	GetSessionsByDateFn   func(ctx context.Context, date string) ([]*domain.Session, error)
	DeleteSessionsByDateFn func(ctx context.Context, date string) error
	SetSessionFn          func(ctx context.Context, s *domain.Session) error
	GetSessionFn          func(ctx context.Context, id string) (*domain.Session, error)
	DeleteSessionFn       func(ctx context.Context, id string) error
}

func (f *FakeMovieCache) SetMovieCache(ctx context.Context, movie *domain.Movie) error {
	if f.SetMovieFn != nil {
		return f.SetMovieFn(ctx, movie)
	}
	return nil
}

func (f *FakeMovieCache) GetMovieCache(ctx context.Context, id string) (*domain.Movie, error) {
	if f.GetMovieFn != nil {
		return f.GetMovieFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeMovieCache) DeleteMovieCache(ctx context.Context, id string) error {
	if f.DeleteMovieFn != nil {
		return f.DeleteMovieFn(ctx, id)
	}
	return nil
}

func (f *FakeMovieCache) SetSessionsByDateCache(ctx context.Context, date string, sessions []*domain.Session) error {
	if f.SetSessionsByDateFn != nil {
		return f.SetSessionsByDateFn(ctx, date, sessions)
	}
	return nil
}

func (f *FakeMovieCache) GetSessionsByDateCache(ctx context.Context, date string) ([]*domain.Session, error) {
	if f.GetSessionsByDateFn != nil {
		return f.GetSessionsByDateFn(ctx, date)
	}
	return nil, nil
}

func (f *FakeMovieCache) DeleteSessionsByDateCache(ctx context.Context, date string) error {
	if f.DeleteSessionsByDateFn != nil {
		return f.DeleteSessionsByDateFn(ctx, date)
	}
	return nil
}

func (f *FakeMovieCache) SetSessionCache(ctx context.Context, session *domain.Session) error {
	if f.SetSessionFn != nil {
		return f.SetSessionFn(ctx, session)
	}
	return nil
}

func (f *FakeMovieCache) GetSessionCache(ctx context.Context, id string) (*domain.Session, error) {
	if f.GetSessionFn != nil {
		return f.GetSessionFn(ctx, id)
	}
	return nil, nil
}

func (f *FakeMovieCache) DeleteSessionCache(ctx context.Context, id string) error {
	if f.DeleteSessionFn != nil {
		return f.DeleteSessionFn(ctx, id)
	}
	return nil
}

type FakeHallRepo struct {
	CreateFn           func(ctx context.Context, h *domain.Hall) (*domain.Hall, error)
	GetFn              func(ctx context.Context, id string) (*domain.Hall, error)
	GetSeatsFn         func(ctx context.Context, hallID string) ([]*domain.Seat, error)
	InsertSeatsForHallFn func(ctx context.Context, hallID string, capacity int) error
}

func (f *FakeHallRepo) Create(ctx context.Context, h *domain.Hall) (*domain.Hall, error) {
	if f.CreateFn != nil {
		return f.CreateFn(ctx, h)
	}
	return nil, domain.ErrInvalidArgument
}

func (f *FakeHallRepo) GetByID(ctx context.Context, id string) (*domain.Hall, error) {
	if f.GetFn != nil {
		return f.GetFn(ctx, id)
	}
	return nil, domain.ErrHallNotFound
}

func (f *FakeHallRepo) GetSeats(ctx context.Context, hallID string) ([]*domain.Seat, error) {
	if f.GetSeatsFn != nil {
		return f.GetSeatsFn(ctx, hallID)
	}
	return nil, nil
}

func (f *FakeHallRepo) InsertSeatsForHall(ctx context.Context, hallID string, capacity int) error {
	if f.InsertSeatsForHallFn != nil {
		return f.InsertSeatsForHallFn(ctx, hallID, capacity)
	}
	return nil
}

type FakeSessionRepo struct {
	CreateFn func(ctx context.Context, s *domain.Session) (*domain.Session, error)
	GetFn    func(ctx context.Context, id string) (*domain.Session, error)
	ListFn   func(ctx context.Context, movieID string, date *time.Time, limit, offset int) ([]*domain.Session, int64, error)
	ByDateFn func(ctx context.Context, day time.Time) ([]*domain.Session, error)
	ByMovieFn func(ctx context.Context, movieID string) ([]*domain.Session, error)
}

func (f *FakeSessionRepo) Create(ctx context.Context, s *domain.Session) (*domain.Session, error) {
	if f.CreateFn != nil {
		return f.CreateFn(ctx, s)
	}
	return nil, domain.ErrInvalidArgument
}

func (f *FakeSessionRepo) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	if f.GetFn != nil {
		return f.GetFn(ctx, id)
	}
	return nil, domain.ErrSessionNotFound
}

func (f *FakeSessionRepo) List(ctx context.Context, movieID string, date *time.Time, limit, offset int) ([]*domain.Session, int64, error) {
	if f.ListFn != nil {
		return f.ListFn(ctx, movieID, date, limit, offset)
	}
	return nil, 0, nil
}

func (f *FakeSessionRepo) GetByDate(ctx context.Context, day time.Time) ([]*domain.Session, error) {
	if f.ByDateFn != nil {
		return f.ByDateFn(ctx, day)
	}
	return nil, nil
}

func (f *FakeSessionRepo) GetByMovieID(ctx context.Context, movieID string) ([]*domain.Session, error) {
	if f.ByMovieFn != nil {
		return f.ByMovieFn(ctx, movieID)
	}
	return nil, nil
}
