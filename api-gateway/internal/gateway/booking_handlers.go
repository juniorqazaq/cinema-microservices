package gateway

import (
	"net/http"

	bookingpb "booking-service/proto/booking"
	"github.com/gin-gonic/gin"
)

func (s *Server) createBooking(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
		SeatID    string `json:"seat_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.CreateBooking(ctx, &bookingpb.CreateBookingRequest{
		UserId:    user.ID,
		SessionId: req.SessionID,
		SeatId:    req.SeatID,
		UserEmail: user.Email,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	created(c, bookingFromProto(resp.GetBooking()))
}

func (s *Server) listUserBookings(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.ListUserBookings(ctx, &bookingpb.ListUserBookingsRequest{UserId: user.ID})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, bookingsFromProto(resp.GetBookings()))
}

func (s *Server) getBookingHistory(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.GetBookingHistory(ctx, &bookingpb.GetHistoryRequest{UserId: user.ID})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, bookingsFromProto(resp.GetBookings()))
}

func (s *Server) getBooking(c *gin.Context) {
	booking, owned := s.requireOwnedBooking(c, c.Param("id"))
	if !owned {
		return
	}
	ok(c, bookingFromProto(booking))
}

func (s *Server) cancelBooking(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	if _, ok := s.requireOwnedBooking(c, c.Param("id")); !ok {
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	_, err := s.clients.Booking.CancelBooking(ctx, &bookingpb.CancelBookingRequest{
		BookingId: c.Param("id"),
		UserEmail: user.Email,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, gin.H{})
}

func (s *Server) confirmPayment(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req struct {
		BookingID string  `json:"booking_id"`
		Amount    float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Amount <= 0 {
		fail(c, http.StatusBadRequest, "amount must be positive")
		return
	}
	if _, ok := s.requireOwnedBooking(c, req.BookingID); !ok {
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.ConfirmPayment(ctx, &bookingpb.ConfirmPaymentRequest{
		BookingId: req.BookingID,
		Amount:    req.Amount,
		UserEmail: user.Email,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, paymentFromProto(resp.GetPayment()))
}

func (s *Server) getPayment(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	resp, err := s.clients.Booking.GetPayment(ctx, &bookingpb.GetPaymentRequest{PaymentId: c.Param("id")})
	cancel()
	if err != nil {
		failGRPC(c, err)
		return
	}
	payment := resp.GetPayment()
	if payment == nil {
		fail(c, http.StatusNotFound, "payment not found")
		return
	}
	if _, ok := s.requireOwnedBooking(c, payment.GetBookingId()); !ok {
		return
	}
	ok(c, paymentFromProto(payment))
}

func (s *Server) requireOwnedBooking(c *gin.Context, bookingID string) (*bookingpb.Booking, bool) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return nil, false
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.GetBooking(ctx, &bookingpb.GetBookingRequest{BookingId: bookingID})
	if err != nil {
		failGRPC(c, err)
		return nil, false
	}
	booking := resp.GetBooking()
	if booking == nil {
		fail(c, http.StatusNotFound, "booking not found")
		return nil, false
	}
	if booking.GetUserId() != user.ID {
		fail(c, http.StatusForbidden, "booking belongs to another user")
		return nil, false
	}
	return booking, true
}

func bookingsFromProto(in []*bookingpb.Booking) []bookingJSON {
	out := make([]bookingJSON, 0, len(in))
	for _, booking := range in {
		out = append(out, bookingFromProto(booking))
	}
	return out
}
