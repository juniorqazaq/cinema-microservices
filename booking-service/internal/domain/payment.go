package domain

import (
	"context"
	"time"
)

// Константы статусов платежа
const (
	PaymentStatusPending  = "pending"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)

// Структура Payment
type Payment struct {
	ID        string
	BookingID string
	Amount    float64
	Status    string     // pending | paid | refunded
	PaidAt    *time.Time // nil если ещё не оплачено
}

// Интерфейс PaymentRepository
type PaymentRepository interface {
	// Create создаёт новый платёж.
	Create(ctx context.Context, payment *Payment) error

	// GetByID возвращает платёж по ID.
	// Возвращает ErrPaymentNotFound если не найдено.
	GetByID(ctx context.Context, id string) (*Payment, error)

	// GetByBookingID возвращает платёж по ID бронирования.
	// Возвращает ErrPaymentNotFound если не найдено.
	GetByBookingID(ctx context.Context, bookingID string) (*Payment, error)

	// List возвращает все платежи с пагинацией (только для admin).
	List(ctx context.Context, limit, offset int) ([]*Payment, error)

	// UpdateStatus обновляет статус платежа.
	UpdateStatus(ctx context.Context, id, status string) error
}
