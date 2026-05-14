package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	"github.com/cinema-booking-system/movie-service/internal/usecase"
	"github.com/cinema-booking-system/movie-service/tests/testutil"
)

func TestCreateMovie_Success(t *testing.T) {
	movies := &testutil.FakeMovieRepo{
		CreateFn: func(ctx context.Context, m *domain.Movie) (*domain.Movie, error) {
			out := *m
			out.ID = "movie-1"
			out.CreatedAt = time.Now().UTC()
			return &out, nil
		},
	}
	cache := &testutil.FakeMovieCache{
		DeleteMovieFn: func(ctx context.Context, id string) error { return nil },
	}
	uc := usecase.NewMovieUseCase(movies, cache)
	ctx := context.Background()
	got, err := uc.CreateMovie(ctx, true, &domain.Movie{Title: "Inception", Description: "x", Genre: "sci-fi", Duration: 120, Rating: 8.8})
	if err != nil {
		t.Fatalf("CreateMovie: %v", err)
	}
	if got.Title != "Inception" || got.ID == "" {
		t.Fatalf("unexpected movie: %+v", got)
	}
}

func TestCreateMovie_EmptyTitle(t *testing.T) {
	uc := usecase.NewMovieUseCase(&testutil.FakeMovieRepo{}, &testutil.FakeMovieCache{})
	_, err := uc.CreateMovie(context.Background(), true, &domain.Movie{Title: "", Duration: 10})
	if !errors.Is(err, domain.ErrInvalidArgument) {
		t.Fatalf("want ErrInvalidArgument, got %v", err)
	}
}

func TestGetMovie_CacheHit(t *testing.T) {
	cached := &domain.Movie{ID: "1", Title: "Hit", Description: "", Genre: "g", Duration: 90, Rating: 7, CreatedAt: time.Now().UTC()}
	cache := &testutil.FakeMovieCache{
		GetMovieFn: func(ctx context.Context, id string) (*domain.Movie, error) {
			if id == "1" {
				return cached, nil
			}
			return nil, nil
		},
	}
	uc := usecase.NewMovieUseCase(&testutil.FakeMovieRepo{}, cache)
	got, err := uc.GetMovie(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Hit" {
		t.Fatalf("want cache hit, got %+v", got)
	}
}

func TestGetMovie_CacheMiss(t *testing.T) {
	dbMovie := &domain.Movie{ID: "2", Title: "DB", Description: "", Genre: "g", Duration: 90, Rating: 7, CreatedAt: time.Now().UTC()}
	movies := &testutil.FakeMovieRepo{
		GetFn: func(ctx context.Context, id string) (*domain.Movie, error) {
			if id == "2" {
				return dbMovie, nil
			}
			return nil, domain.ErrMovieNotFound
		},
	}
	var setCalled bool
	cache := &testutil.FakeMovieCache{
		GetMovieFn: func(ctx context.Context, id string) (*domain.Movie, error) { return nil, nil },
		SetMovieFn: func(ctx context.Context, m *domain.Movie) error {
			setCalled = true
			return nil
		},
	}
	uc := usecase.NewMovieUseCase(movies, cache)
	got, err := uc.GetMovie(context.Background(), "2")
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "DB" || !setCalled {
		t.Fatalf("miss path broken: %+v setCalled=%v", got, setCalled)
	}
}

func TestGetAvailableSeats(t *testing.T) {
	sessions := &testutil.FakeSessionRepo{
		GetFn: func(ctx context.Context, id string) (*domain.Session, error) {
			return &domain.Session{ID: id, HallID: "h1"}, nil
		},
	}
	halls := &testutil.FakeHallRepo{
		GetSeatsFn: func(ctx context.Context, hallID string) ([]*domain.Seat, error) {
			return []*domain.Seat{
				{ID: "s1", HallID: hallID, Row: "R1", Number: 1, IsAvailable: true},
				{ID: "s2", HallID: hallID, Row: "R1", Number: 2, IsAvailable: false},
			}, nil
		},
	}
	uc := usecase.NewSessionUseCase(&testutil.FakeMovieRepo{}, halls, sessions, &testutil.FakeMovieCache{})
	seats, err := uc.GetAvailableSeats(context.Background(), "sess-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(seats) != 1 || seats[0].ID != "s1" {
		t.Fatalf("seats=%v", seats)
	}
}

func TestCreateSession_InvalidMovie(t *testing.T) {
	movies := &testutil.FakeMovieRepo{
		GetFn: func(ctx context.Context, id string) (*domain.Movie, error) {
			return nil, domain.ErrMovieNotFound
		},
	}
	uc := usecase.NewSessionUseCase(movies, &testutil.FakeHallRepo{}, &testutil.FakeSessionRepo{}, &testutil.FakeMovieCache{})
	_, err := uc.CreateSession(context.Background(), &domain.Session{MovieID: "bad", HallID: "h", StartTime: time.Now().UTC(), Price: 10})
	if !errors.Is(err, domain.ErrMovieNotFound) {
		t.Fatalf("want ErrMovieNotFound, got %v", err)
	}
}
