package grpc

import (
	"context"
	"errors"

	"booking-service/internal/domain"
	"booking-service/internal/usecase"
	pb "booking-service/proto/booking"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	pb.UnimplementedBookingServiceServer
	bookingUC   *usecase.BookingUseCase
	paymentUC   *usecase.PaymentUseCase
	bookingRepo domain.BookingRepository
	paymentRepo domain.PaymentRepository
}

func NewHandler(
	bookingUC *usecase.BookingUseCase,
	paymentUC *usecase.PaymentUseCase,
	bookingRepo domain.BookingRepository,
	paymentRepo domain.PaymentRepository,
) *Handler {
	return &Handler{
		bookingUC:   bookingUC,
		paymentUC:   paymentUC,
		bookingRepo: bookingRepo,
		paymentRepo: paymentRepo,
	}
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrSeatTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrBookingNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAlreadyCancelled):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, domain.ErrPaymentNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error: "+err.Error())
	}
}

func mapBookingToPB(b *domain.Booking) *pb.Booking {
	if b == nil {
		return nil
	}
	return &pb.Booking{
		Id:        b.ID,
		UserId:    b.UserID,
		SessionId: b.SessionID,
		SeatId:    b.SeatID,
		Status:    b.Status,
		CreatedAt: b.CreatedAt.Unix(),
	}
}

func mapPaymentToPB(p *domain.Payment) *pb.Payment {
	if p == nil {
		return nil
	}
	var paidAt int64
	if p.PaidAt != nil {
		paidAt = p.PaidAt.Unix()
	}
	return &pb.Payment{
		Id:        p.ID,
		BookingId: p.BookingID,
		Amount:    p.Amount,
		Status:    p.Status,
		PaidAt:    paidAt,
	}
}

func (h *Handler) CreateBooking(ctx context.Context, req *pb.CreateBookingRequest) (*pb.CreateBookingResponse, error) {
	b, err := h.bookingUC.CreateBooking(ctx, req.UserId, req.SessionId, req.SeatId, req.UserEmail)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.CreateBookingResponse{Booking: mapBookingToPB(b)}, nil
}

func (h *Handler) GetBooking(ctx context.Context, req *pb.GetBookingRequest) (*pb.GetBookingResponse, error) {
	b, err := h.bookingRepo.GetByID(ctx, req.BookingId)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.GetBookingResponse{Booking: mapBookingToPB(b)}, nil
}

func (h *Handler) ListUserBookings(ctx context.Context, req *pb.ListUserBookingsRequest) (*pb.ListUserBookingsResponse, error) {
	bookings, err := h.bookingRepo.ListByUserID(ctx, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}
	pbBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = mapBookingToPB(b)
	}
	return &pb.ListUserBookingsResponse{Bookings: pbBookings}, nil
}

func (h *Handler) CancelBooking(ctx context.Context, req *pb.CancelBookingRequest) (*pb.CancelBookingResponse, error) {
	err := h.bookingUC.CancelBooking(ctx, req.BookingId, req.UserEmail)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.CancelBookingResponse{Success: true}, nil
}

func (h *Handler) ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*pb.ConfirmPaymentResponse, error) {
	p, err := h.paymentUC.ConfirmPayment(ctx, req.BookingId, req.Amount, req.UserEmail)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.ConfirmPaymentResponse{Payment: mapPaymentToPB(p)}, nil
}

func (h *Handler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.GetPaymentResponse, error) {
	var p *domain.Payment
	var err error

	switch {
	case req.PaymentId != "":
		p, err = h.paymentRepo.GetByID(ctx, req.PaymentId)
	case req.BookingId != "":
		p, err = h.paymentRepo.GetByBookingID(ctx, req.BookingId)
	default:
		return nil, status.Error(codes.InvalidArgument, "either payment_id or booking_id must be provided")
	}

	if err != nil {
		return nil, mapError(err)
	}
	return &pb.GetPaymentResponse{Payment: mapPaymentToPB(p)}, nil
}

func (h *Handler) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	payments, err := h.paymentRepo.List(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, mapError(err)
	}
	total, err := h.paymentRepo.CountAll(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	pbPayments := make([]*pb.Payment, len(payments))
	for i, p := range payments {
		pbPayments[i] = mapPaymentToPB(p)
	}
	return &pb.ListPaymentsResponse{Payments: pbPayments, Total: int32(total)}, nil
}

func (h *Handler) RefundPayment(ctx context.Context, req *pb.RefundPaymentRequest) (*pb.RefundPaymentResponse, error) {
	err := h.paymentRepo.UpdateStatus(ctx, req.PaymentId, domain.PaymentStatusRefunded)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.RefundPaymentResponse{Success: true}, nil
}

func (h *Handler) CheckSeatAvailability(ctx context.Context, req *pb.CheckSeatRequest) (*pb.CheckSeatResponse, error) {
	ok, err := h.bookingRepo.IsSeatAvailable(ctx, req.SeatId)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.CheckSeatResponse{IsAvailable: ok}, nil
}

func (h *Handler) GetBookingHistory(ctx context.Context, req *pb.GetHistoryRequest) (*pb.GetHistoryResponse, error) {
	bookings, err := h.bookingRepo.GetHistory(ctx, req.UserId)
	if err != nil {
		return nil, mapError(err)
	}
	pbBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = mapBookingToPB(b)
	}
	return &pb.GetHistoryResponse{Bookings: pbBookings}, nil
}

func (h *Handler) AdminListBookings(ctx context.Context, req *pb.AdminListRequest) (*pb.AdminListResponse, error) {
	bookings, err := h.bookingRepo.AdminListAll(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, mapError(err)
	}
	total, err := h.bookingRepo.CountAll(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	pbBookings := make([]*pb.Booking, len(bookings))
	for i, b := range bookings {
		pbBookings[i] = mapBookingToPB(b)
	}
	return &pb.AdminListResponse{Bookings: pbBookings, Total: int32(total)}, nil
}

func (h *Handler) GetBookingStats(ctx context.Context, _ *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	total, confirmed, cancelled, err := h.bookingRepo.GetStats(ctx)
	if err != nil {
		return nil, mapError(err)
	}
	return &pb.GetStatsResponse{
		Total:     total,
		Confirmed: confirmed,
		Cancelled: cancelled,
	}, nil
}
