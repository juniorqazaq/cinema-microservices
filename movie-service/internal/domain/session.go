package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrHallNotFound    = errors.New("hall not found")
	ErrSessionNotFound = errors.New("session not found")
)

type Hall struct {
	ID         string
	Name       string
	Capacity   int
	City       string
	CinemaName string
}

type Seat struct {
	ID           string
	HallID       string
	Row          string
	Number       int
	IsAvailable  bool
}

type Session struct {
	ID        string
	MovieID   string
	HallID    string
	StartTime time.Time
	Price     float64
}

type HallRepository interface {
	Create(ctx context.Context, h *Hall) (*Hall, error)
	GetByID(ctx context.Context, id string) (*Hall, error)
	GetSeats(ctx context.Context, hallID string) ([]*Seat, error)
	InsertSeatsForHall(ctx context.Context, hallID string, capacity int) error
}

type SessionRepository interface {
	Create(ctx context.Context, s *Session) (*Session, error)
	GetByID(ctx context.Context, id string) (*Session, error)
	List(ctx context.Context, movieID, city string, date *time.Time, limit, offset int) ([]*Session, int64, error)
	GetByDate(ctx context.Context, day time.Time) ([]*Session, error)
	GetByMovieID(ctx context.Context, movieID string) ([]*Session, error)
}

// MovieCache abstracts Redis caching for movies and sessions.
type MovieCache interface {
	SetMovieCache(ctx context.Context, movie *Movie) error
	GetMovieCache(ctx context.Context, id string) (*Movie, error)
	DeleteMovieCache(ctx context.Context, id string) error

	SetSessionsByDateCache(ctx context.Context, date string, sessions []*Session) error
	GetSessionsByDateCache(ctx context.Context, date string) ([]*Session, error)
	DeleteSessionsByDateCache(ctx context.Context, date string) error

	SetSessionCache(ctx context.Context, session *Session) error
	GetSessionCache(ctx context.Context, id string) (*Session, error)
	DeleteSessionCache(ctx context.Context, id string) error
}
