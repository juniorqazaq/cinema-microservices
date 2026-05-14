package usecase_test

import (
	"context"
	"errors"
	"testing"

	"booking-service/internal/domain"
	"booking-service/internal/usecase"
	"booking-service/tests/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateBooking_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingRepository(ctrl)
	mockDB := mocks.NewMockDB(ctrl)
	mockPublisher := mocks.NewMockEventPublisher(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := usecase.NewBookingUseCase(mockRepo, mockDB, mockPublisher)
	ctx := context.Background()

	mockDB.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)

	mockRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(nil)
	mockPublisher.EXPECT().PublishBookingCreated(ctx, gomock.Any()).Return(nil)

	booking, err := uc.CreateBooking(ctx, "user-1", "session-1", "seat-1", "")

	assert.NoError(t, err)
	assert.NotNil(t, booking)
	assert.Equal(t, "user-1", booking.UserID)
}

func TestCreateBooking_SeatTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingRepository(ctrl)
	mockDB := mocks.NewMockDB(ctrl)
	mockPublisher := mocks.NewMockEventPublisher(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := usecase.NewBookingUseCase(mockRepo, mockDB, mockPublisher)
	ctx := context.Background()

	mockDB.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)

	// repo returns ErrSeatTaken
	mockRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(domain.ErrSeatTaken)

	// Publisher should NOT be called
	mockPublisher.EXPECT().PublishBookingCreated(gomock.Any(), gomock.Any()).Times(0)

	booking, err := uc.CreateBooking(ctx, "user-1", "session-1", "seat-1", "")

	assert.ErrorIs(t, err, domain.ErrSeatTaken)
	assert.Nil(t, booking)
}

func TestCancelBooking(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingRepository(ctrl)
	mockDB := mocks.NewMockDB(ctrl)
	mockPublisher := mocks.NewMockEventPublisher(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := usecase.NewBookingUseCase(mockRepo, mockDB, mockPublisher)
	ctx := context.Background()

	mockDB.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	mockTx.EXPECT().Commit(ctx).Return(nil)

	mockRepo.EXPECT().Cancel(ctx, mockTx, "booking-1").Return(nil)
	mockPublisher.EXPECT().PublishBookingCancelled(ctx, gomock.Any()).Return(nil).Times(1)

	err := uc.CancelBooking(ctx, "booking-1", "")

	assert.NoError(t, err)
}

func TestConfirmPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPaymentRepo := mocks.NewMockPaymentRepository(ctrl)
	mockBookingRepo := mocks.NewMockBookingRepository(ctrl)
	mockPublisher := mocks.NewMockEventPublisher(ctrl)

	uc := usecase.NewPaymentUseCase(mockPaymentRepo, mockBookingRepo, mockPublisher)
	ctx := context.Background()

	mockPaymentRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
	mockBookingRepo.EXPECT().UpdateStatus(ctx, "booking-1", domain.StatusConfirmed).Return(nil)
	mockPublisher.EXPECT().PublishPaymentConfirmed(ctx, gomock.Any()).Return(nil).Times(1)

	payment, err := uc.ConfirmPayment(ctx, "booking-1", 3000, "user@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, payment)
}

func TestCreateBooking_NATSFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBookingRepository(ctrl)
	mockDB := mocks.NewMockDB(ctrl)
	mockPublisher := mocks.NewMockEventPublisher(ctrl)
	mockTx := mocks.NewMockTx(ctrl)

	uc := usecase.NewBookingUseCase(mockRepo, mockDB, mockPublisher)
	ctx := context.Background()

	mockDB.EXPECT().Begin(ctx).Return(mockTx, nil)
	mockTx.EXPECT().Rollback(ctx).Return(nil)
	// Commit should NOT be called

	mockRepo.EXPECT().Create(ctx, mockTx, gomock.Any()).Return(nil)
	mockPublisher.EXPECT().PublishBookingCreated(ctx, gomock.Any()).Return(errors.New("nats down"))

	booking, err := uc.CreateBooking(ctx, "user-1", "session-1", "seat-1", "")

	assert.Error(t, err)
	assert.Equal(t, "nats down", err.Error())
	assert.Nil(t, booking)
}
