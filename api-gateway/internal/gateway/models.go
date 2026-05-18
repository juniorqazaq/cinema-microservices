package gateway

import (
	"net/http"
	"time"

	bookingpb "booking-service/proto/booking"

	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type userJSON struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FullName  string  `json:"full_name,omitempty"`
	Role      string  `json:"role"`
	IsBanned  bool    `json:"is_banned"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type movieJSON struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Genre       string  `json:"genre"`
	Duration    int32   `json:"duration"`
	Rating      float64 `json:"rating"`
	AgeRating   int32   `json:"age_rating"`
	CreatedAt   string  `json:"created_at"`
}

type hallJSON struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Capacity   int32  `json:"capacity"`
	City       string `json:"city"`
	CinemaName string `json:"cinema_name"`
}

type seatJSON struct {
	ID          string `json:"id"`
	HallID      string `json:"hall_id"`
	Row         string `json:"row"`
	Number      int32  `json:"number"`
	IsAvailable bool   `json:"is_available"`
}

type sessionJSON struct {
	ID             string  `json:"id"`
	MovieID        string  `json:"movie_id"`
	HallID         string  `json:"hall_id"`
	StartTime      string  `json:"start_time"`
	Price          float64 `json:"price"`
	City           string  `json:"city,omitempty"`
	CinemaName     string  `json:"cinema_name,omitempty"`
	HallName       string  `json:"hall_name,omitempty"`
	AvailableSeats int32   `json:"available_seats,omitempty"`
	AgeRating      int32   `json:"age_rating,omitempty"`
}

type bookingJSON struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	SessionID      string  `json:"session_id"`
	SeatID         string  `json:"seat_id"`
	Status         string  `json:"status"`
	TicketCategory string  `json:"ticket_category,omitempty"`
	AmountPaid     float64 `json:"amount_paid,omitempty"`
	CreatedAt      string  `json:"created_at"`
}

type paymentJSON struct {
	ID        string  `json:"id"`
	BookingID string  `json:"booking_id"`
	Amount    float64 `json:"amount"`
	Status    string  `json:"status"`
	PaidAt    *string `json:"paid_at"`
}

type statsJSON struct {
	Total     int64 `json:"total"`
	Confirmed int64 `json:"confirmed"`
	Cancelled int64 `json:"cancelled"`
}

func ok(c *gin.Context, payload any) {
	if payload == nil {
		payload = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{"data": payload})
}

func created(c *gin.Context, payload any) {
	c.JSON(http.StatusCreated, gin.H{"data": payload})
}

func fail(c *gin.Context, statusCode int, message string) {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	c.JSON(statusCode, gin.H{"error": message})
}

func failGRPC(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		fail(c, http.StatusInternalServerError, "upstream service error")
		return
	}
	switch st.Code() {
	case codes.InvalidArgument:
		fail(c, http.StatusBadRequest, st.Message())
	case codes.Unauthenticated:
		fail(c, http.StatusUnauthorized, st.Message())
	case codes.PermissionDenied:
		fail(c, http.StatusForbidden, st.Message())
	case codes.NotFound:
		fail(c, http.StatusNotFound, st.Message())
	case codes.AlreadyExists, codes.FailedPrecondition:
		fail(c, http.StatusConflict, st.Message())
	case codes.ResourceExhausted:
		fail(c, http.StatusTooManyRequests, st.Message())
	case codes.DeadlineExceeded:
		fail(c, http.StatusGatewayTimeout, st.Message())
	case codes.Unavailable:
		fail(c, http.StatusServiceUnavailable, st.Message())
	default:
		fail(c, http.StatusInternalServerError, "upstream service error")
	}
}

func userFromProto(u *userpb.User) userJSON {
	if u == nil {
		return userJSON{}
	}
	return userJSON{
		ID:        u.GetUserId(),
		Email:     u.GetEmail(),
		FullName:  u.GetFullName(),
		Role:      roleForFrontend(u.GetRole()),
		IsBanned:  u.GetIsBanned(),
		Balance:   u.GetBalance(),
		CreatedAt: tsString(u.GetCreatedAt()),
		UpdatedAt: tsString(u.GetUpdatedAt()),
	}
}

func movieFromProto(m *moviepb.Movie) movieJSON {
	if m == nil {
		return movieJSON{}
	}
	return movieJSON{
		ID:          m.GetId(),
		Title:       m.GetTitle(),
		Description: m.GetDescription(),
		Genre:       m.GetGenre(),
		Duration:    m.GetDuration(),
		Rating:      m.GetRating(),
		AgeRating:   m.GetAgeRating(),
		CreatedAt:   tsString(m.GetCreatedAt()),
	}
}

func hallFromProto(h *moviepb.Hall) hallJSON {
	if h == nil {
		return hallJSON{}
	}
	return hallJSON{
		ID:         h.GetId(),
		Name:       h.GetName(),
		Capacity:   h.GetCapacity(),
		City:       h.GetCity(),
		CinemaName: h.GetCinemaName(),
	}
}

func seatFromProto(s *moviepb.Seat) seatJSON {
	if s == nil {
		return seatJSON{}
	}
	return seatJSON{
		ID:          s.GetId(),
		HallID:      s.GetHallId(),
		Row:         s.GetRow(),
		Number:      s.GetNumber(),
		IsAvailable: s.GetIsAvailable(),
	}
}

func sessionFromProto(s *moviepb.Session) sessionJSON {
	if s == nil {
		return sessionJSON{}
	}
	return sessionJSON{
		ID:        s.GetId(),
		MovieID:   s.GetMovieId(),
		HallID:    s.GetHallId(),
		StartTime: tsString(s.GetStartTime()),
		Price:     s.GetPrice(),
	}
}

func bookingFromProto(b *bookingpb.Booking) bookingJSON {
	if b == nil {
		return bookingJSON{}
	}
	return bookingJSON{
		ID:             b.GetId(),
		UserID:         b.GetUserId(),
		SessionID:      b.GetSessionId(),
		SeatID:         b.GetSeatId(),
		Status:         b.GetStatus(),
		TicketCategory: b.GetTicketCategory(),
		AmountPaid:     b.GetAmountPaid(),
		CreatedAt:      unixString(b.GetCreatedAt()),
	}
}

func paymentFromProto(p *bookingpb.Payment) paymentJSON {
	if p == nil {
		return paymentJSON{}
	}
	var paidAt *string
	if p.GetPaidAt() > 0 {
		v := unixString(p.GetPaidAt())
		paidAt = &v
	}
	return paymentJSON{
		ID:        p.GetId(),
		BookingID: p.GetBookingId(),
		Amount:    p.GetAmount(),
		Status:    p.GetStatus(),
		PaidAt:    paidAt,
	}
}

func roleForFrontend(role userpb.Role) string {
	if role == userpb.Role_ROLE_ADMIN {
		return "admin"
	}
	return "user"
}

func movieRole(role userpb.Role) moviepb.Role {
	if role == userpb.Role_ROLE_ADMIN {
		return moviepb.Role_ROLE_ADMIN
	}
	return moviepb.Role_ROLE_USER
}

func tsString(ts *timestamppb.Timestamp) string {
	if ts == nil {
		return ""
	}
	return ts.AsTime().UTC().Format(time.RFC3339)
}

func unixString(seconds int64) string {
	if seconds <= 0 {
		return ""
	}
	return time.Unix(seconds, 0).UTC().Format(time.RFC3339)
}
