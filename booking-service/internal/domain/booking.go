package domain

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
)

var (
	ErrSeatTaken        = errors.New("seat is already taken")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrAlreadyCancelled = errors.New("booking is already cancelled")
	ErrPaymentNotFound  = errors.New("payment not found")
)

type Booking struct {
	ID        string
	UserID    string
	SessionID string
	SeatID    string
	Status    string
	CreatedAt time.Time
}

type BookingRepository interface {
	Create(ctx context.Context, tx pgx.Tx, booking *Booking) error
	GetByID(ctx context.Context, id string) (*Booking, error)
	ListByUserID(ctx context.Context, userID string) ([]*Booking, error)
	Cancel(ctx context.Context, tx pgx.Tx, id string) error
	GetHistory(ctx context.Context, userID string) ([]*Booking, error)
	AdminListAll(ctx context.Context, limit, offset int) ([]*Booking, error)
	GetStats(ctx context.Context) (total, confirmed, cancelled int64, err error)
	CountAll(ctx context.Context) (int64, error)
	IsSeatAvailable(ctx context.Context, seatID string) (bool, error)
	UpdateStatus(ctx context.Context, id, status string) error
}
