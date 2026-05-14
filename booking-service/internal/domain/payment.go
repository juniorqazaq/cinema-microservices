package domain

import (
	"context"
	"time"
)

const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)

type Payment struct {
	ID        string
	BookingID string
	Amount    float64
	Status    string
	PaidAt    *time.Time
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	GetByBookingID(ctx context.Context, bookingID string) (*Payment, error)
	List(ctx context.Context, limit, offset int) ([]*Payment, error)
	CountAll(ctx context.Context) (int64, error)
	UpdateStatus(ctx context.Context, id, status string) error
}
