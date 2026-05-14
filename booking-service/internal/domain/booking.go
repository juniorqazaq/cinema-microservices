package domain

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

// Константы статусов
const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
)

// Кастомные ошибки
var (
	ErrSeatTaken        = errors.New("seat is already taken")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrAlreadyCancelled = errors.New("booking is already cancelled")
	ErrPaymentNotFound  = errors.New("payment not found")
)

// Структура Booking
type Booking struct {
	ID        string
	UserID    string
	SessionID string
	SeatID    string
	Status    string // pending | confirmed | cancelled
	CreatedAt time.Time
}

// Интерфейс BookingRepository
type BookingRepository interface {
	// Create создаёт бронирование внутри переданной транзакции.
	// Транзакцию открывает UseCase, не Repository.
	Create(ctx context.Context, tx pgx.Tx, booking *Booking) error

	// GetByID возвращает бронирование по ID.
	// Возвращает ErrBookingNotFound если не найдено.
	GetByID(ctx context.Context, id string) (*Booking, error)

	// ListByUserID возвращает активные бронирования пользователя.
	ListByUserID(ctx context.Context, userID string) ([]*Booking, error)

	// Cancel отменяет бронирование внутри переданной транзакции.
	// Возвращает ErrAlreadyCancelled если уже отменено.
	Cancel(ctx context.Context, tx pgx.Tx, id string) error

	// GetHistory возвращает всю историю бронирований пользователя.
	GetHistory(ctx context.Context, userID string) ([]*Booking, error)

	// AdminListAll возвращает все бронирования с пагинацией (только для admin).
	AdminListAll(ctx context.Context, limit, offset int) ([]*Booking, error)

	// GetStats возвращает статистику: всего, подтверждённых, отменённых.
	GetStats(ctx context.Context) (total, confirmed, cancelled int64, err error)

	// UpdateStatus обновляет статус бронирования (например, при подтверждении оплаты).
	UpdateStatus(ctx context.Context, id, status string) error
}
