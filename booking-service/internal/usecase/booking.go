package usecase

import (
	"context"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/publisher"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EventPublisher interface {
	PublishBookingCreated(ctx context.Context, event *publisher.BookingCreatedEvent) error
	PublishBookingCancelled(ctx context.Context, event *publisher.BookingCancelledEvent) error
	PublishPaymentConfirmed(ctx context.Context, event *publisher.PaymentConfirmedEvent) error
}

type DB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type BookingUseCase struct {
	repo      domain.BookingRepository
	pool      DB
	publisher EventPublisher
}

func NewBookingUseCase(repo domain.BookingRepository, pool DB, pub EventPublisher) *BookingUseCase {
	return &BookingUseCase{
		repo:      repo,
		pool:      pool,
		publisher: pub,
	}
}

func (uc *BookingUseCase) CreateBooking(ctx context.Context, userID, sessionID, seatID, userEmail, ticketCategory string, amount float64) (*domain.Booking, error) {
	if ticketCategory == "" {
		ticketCategory = "adult"
	}
	booking := &domain.Booking{
		ID:             uuid.NewString(),
		UserID:         userID,
		SessionID:      sessionID,
		SeatID:         seatID,
		Status:         domain.StatusPending,
		TicketCategory: ticketCategory,
		AmountPaid:     amount,
		CreatedAt:      time.Now(),
	}

	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = uc.repo.Create(ctx, tx, booking)
	if err != nil {
		return nil, err
	}

	email := userEmail
	if email == "" {
		email = "user@example.com"
	}
	event := &publisher.BookingCreatedEvent{
		BookingID: booking.ID,
		UserID:    booking.UserID,
		Email:     email,
		Movie:     "Unknown Movie",
		Seat:      seatID,
		Time:      booking.CreatedAt.Format(time.RFC3339),
	}

	err = uc.publisher.PublishBookingCreated(ctx, event)
	if err != nil {
		return nil, err // NATS упал — транзакция откатится через defer
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return booking, nil
}

func (uc *BookingUseCase) CancelBooking(ctx context.Context, id, userEmail string) error {
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = uc.repo.Cancel(ctx, tx, id)
	if err != nil {
		return err
	}

	event := &publisher.BookingCancelledEvent{
		BookingID: id,
		Email:     userEmail,
		Date:      time.Now().UTC().Format(time.RFC3339),
	}

	err = uc.publisher.PublishBookingCancelled(ctx, event)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
