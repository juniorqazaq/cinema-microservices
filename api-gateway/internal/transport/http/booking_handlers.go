package httpgateway

import (
	"net/http"

	bookingpb "booking-service/proto/booking"
	"github.com/gin-gonic/gin"
)

func (s *Server) createBooking(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}
	var req createBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": bookingFromProto(resp.GetBooking())})
}

func (s *Server) listUserBookings(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.ListUserBookings(ctx, &bookingpb.ListUserBookingsRequest{UserId: user.ID})
	if err != nil {
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookingsFromProto(resp.GetBookings())})
}

func (s *Server) getBookingHistory(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.GetBookingHistory(ctx, &bookingpb.GetHistoryRequest{UserId: user.ID})
	if err != nil {
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookingsFromProto(resp.GetBookings())})
}

func (s *Server) getBooking(c *gin.Context) {
	booking, owned := s.requireOwnedBooking(c, c.Param("id"))
	if !owned {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookingFromProto(booking)})
}

func (s *Server) cancelBooking(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
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
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{}})
}

func (s *Server) confirmPayment(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return
	}
	var req confirmPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paymentFromProto(resp.GetPayment())})
}

func (s *Server) getPayment(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	resp, err := s.clients.Booking.GetPayment(ctx, &bookingpb.GetPaymentRequest{PaymentId: c.Param("id")})
	cancel()
	if err != nil {
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return
	}
	payment := resp.GetPayment()
	if payment == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}
	if _, ok := s.requireOwnedBooking(c, payment.GetBookingId()); !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": paymentFromProto(payment)})
}

func (s *Server) requireOwnedBooking(c *gin.Context, bookingID string) (*bookingpb.Booking, bool) {
	user, okUser := currentUser(c)
	if !okUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return nil, false
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.GetBooking(ctx, &bookingpb.GetBookingRequest{BookingId: bookingID})
	if err != nil {
		statusCode, message := grpcError(err)
		c.JSON(statusCode, gin.H{"error": message})
		return nil, false
	}
	booking := resp.GetBooking()
	if booking == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "booking not found"})
		return nil, false
	}
	if booking.GetUserId() != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "booking belongs to another user"})
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
