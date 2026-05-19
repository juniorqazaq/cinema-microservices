package usecase

import (
	"context"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	"github.com/cinema-booking-system/movie-service/internal/repository/redis"
)

type SessionUseCase struct {
	movies   domain.MovieRepository
	halls    domain.HallRepository
	sessions domain.SessionRepository
	cache    domain.MovieCache
}

func NewSessionUseCase(
	movies domain.MovieRepository,
	halls domain.HallRepository,
	sessions domain.SessionRepository,
	cache domain.MovieCache,
) *SessionUseCase {
	return &SessionUseCase{movies: movies, halls: halls, sessions: sessions, cache: cache}
}

func (u *SessionUseCase) CreateHall(ctx context.Context, isAdmin bool, name string, capacity int, city, cinemaName string) (*domain.Hall, error) {
	if !isAdmin {
		return nil, domain.ErrForbidden
	}
	if name == "" || capacity <= 0 {
		return nil, domain.ErrInvalidArgument
	}
	h := &domain.Hall{Name: name, Capacity: capacity, City: city, CinemaName: cinemaName}
	created, err := u.halls.Create(ctx, h)
	if err != nil {
		return nil, err
	}
	if err := u.halls.InsertSeatsForHall(ctx, created.ID, capacity); err != nil {
		return nil, err
	}
	return created, nil
}

func (u *SessionUseCase) GetHall(ctx context.Context, hallID string) (*domain.Hall, []*domain.Seat, error) {
	if hallID == "" {
		return nil, nil, domain.ErrInvalidArgument
	}
	h, err := u.halls.GetByID(ctx, hallID)
	if err != nil {
		return nil, nil, err
	}
	seats, err := u.halls.GetSeats(ctx, hallID)
	if err != nil {
		return nil, nil, err
	}
	return h, seats, nil
}

func (u *SessionUseCase) CreateSession(ctx context.Context, s *domain.Session) (*domain.Session, error) {
	if s.MovieID == "" || s.HallID == "" {
		return nil, domain.ErrInvalidArgument
	}
	if s.Price < 0 {
		return nil, domain.ErrInvalidArgument
	}
	if _, err := u.movies.GetByID(ctx, s.MovieID); err != nil {
		return nil, err
	}
	if _, err := u.halls.GetByID(ctx, s.HallID); err != nil {
		return nil, err
	}
	created, err := u.sessions.Create(ctx, s)
	if err != nil {
		if err == domain.ErrInvalidArgument {
			return nil, domain.ErrMovieNotFound
		}
		return nil, err
	}
	dateKey := redis.SessionDateKey(created.StartTime)
	_ = u.cache.DeleteSessionsByDateCache(ctx, dateKey)
	return created, nil
}

func (u *SessionUseCase) GetSession(ctx context.Context, id string) (*domain.Session, error) {
	if id == "" {
		return nil, domain.ErrInvalidArgument
	}
	if cached, err := u.cache.GetSessionCache(ctx, id); err != nil {
		return nil, err
	} else if cached != nil {
		return cached, nil
	}
	s, err := u.sessions.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = u.cache.SetSessionCache(ctx, s)
	return s, nil
}

func (u *SessionUseCase) ListSessions(ctx context.Context, movieID, city string, day *time.Time, limit, offset int) ([]*domain.Session, int64, error) {
	limit, offset, err := clampPagination(limit, offset)
	if err != nil {
		return nil, 0, err
	}
	// Date+city+movie filters need SQL (cache path ignores city).
	if day != nil && (movieID != "" || city != "") {
		return u.sessions.List(ctx, movieID, city, day, limit, offset)
	}
	if day != nil {
		dateKey := redis.SessionDateKey(*day)
		list, err := u.cache.GetSessionsByDateCache(ctx, dateKey)
		if err != nil {
			return nil, 0, err
		}
		if list == nil {
			list, err = u.sessions.GetByDate(ctx, *day)
			if err != nil {
				return nil, 0, err
			}
			_ = u.cache.SetSessionsByDateCache(ctx, dateKey, list)
		}
		filtered := list
		if movieID != "" {
			next := make([]*domain.Session, 0)
			for _, s := range list {
				if s.MovieID == movieID {
					next = append(next, s)
				}
			}
			filtered = next
		}
		total := int64(len(filtered))
		end := offset + limit
		if offset > len(filtered) {
			return []*domain.Session{}, total, nil
		}
		if end > len(filtered) {
			end = len(filtered)
		}
		return filtered[offset:end], total, nil
	}
	return u.sessions.List(ctx, movieID, city, nil, limit, offset)
}

func (u *SessionUseCase) GetAvailableSeats(ctx context.Context, sessionID string) ([]*domain.Seat, error) {
	if sessionID == "" {
		return nil, domain.ErrInvalidArgument
	}
	s, err := u.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	seats, err := u.halls.GetSeats(ctx, s.HallID)
	if err != nil {
		return nil, err
	}
	return seats, nil
}
