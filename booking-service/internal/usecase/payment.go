package usecase

import (
	"context"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/publisher"
	"github.com/google/uuid"
)

type PaymentUseCase struct {
	paymentRepo domain.PaymentRepository
	bookingRepo domain.BookingRepository
	publisher   EventPublisher
}

func NewPaymentUseCase(pRepo domain.PaymentRepository, bRepo domain.BookingRepository, pub EventPublisher) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: pRepo,
		bookingRepo: bRepo,
		publisher:   pub,
	}
}

func (uc *PaymentUseCase) ConfirmPayment(ctx context.Context, bookingID string, amount float64, email string) (*domain.Payment, error) {
	now := time.Now()
	payment := &domain.Payment{
		ID:        uuid.NewString(),
		BookingID: bookingID,
		Amount:    amount,
		Status:    domain.PaymentStatusPaid,
		PaidAt:    &now,
	}

	err := uc.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	err = uc.bookingRepo.UpdateStatus(ctx, bookingID, domain.StatusConfirmed)
	if err != nil {
		return nil, err
	}

	event := &publisher.PaymentConfirmedEvent{
		PaymentID: payment.ID,
		BookingID: payment.BookingID,
		Amount:    payment.Amount,
		Email:     email,
	}

	err = uc.publisher.PublishPaymentConfirmed(ctx, event)
	if err != nil {
		return nil, err
	}

	return payment, nil
}
