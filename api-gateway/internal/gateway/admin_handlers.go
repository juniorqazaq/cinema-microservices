package gateway

import (
	"encoding/json"
	"net/http"
	"time"

	bookingpb "booking-service/proto/booking"
	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) adminListUsers(c *gin.Context) {
	admin, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	page := int32(queryInt(c, "page", 1))
	limit := int32(queryInt(c, "limit", 50))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.GetAllUsers(ctx, &userpb.GetAllUsersRequest{
		AdminId: admin.ID,
		Pagination: &userpb.PaginationRequest{
			Page:     page,
			PageSize: limit,
		},
		RoleFilter: c.Query("role"),
		BannedOnly: c.Query("banned_only") == "true",
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	out := make([]userJSON, 0, len(resp.GetUsers()))
	for _, user := range resp.GetUsers() {
		out = append(out, userFromProto(user))
	}
	ok(c, out)
}

func (s *Server) adminBanUser(c *gin.Context) {
	admin, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	_, err := s.clients.User.BanUser(ctx, &userpb.BanUserRequest{
		AdminId: admin.ID,
		UserId:  c.Param("id"),
		Ban:     true,
		Reason:  "banned by admin",
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, gin.H{})
}

func (s *Server) adminCreateMovie(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req moviePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.CreateMovie(ctx, &moviepb.CreateMovieRequest{
		Title:         req.Title.Value(),
		Description:   req.Description.Value(),
		Genre:         req.Genre.Value(),
		Duration:      req.Duration.Value(),
		Rating:        req.Rating.Value(),
		RequesterRole: movieRole(user.Role),
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	created(c, movieFromProto(resp.GetMovie()))
}

func (s *Server) adminUpdateMovie(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req moviePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	current, err := s.clients.Movie.GetMovie(ctx, &moviepb.GetMovieRequest{Id: c.Param("id")})
	cancel()
	if err != nil {
		failGRPC(c, err)
		return
	}
	movie := current.GetMovie()
	if movie == nil {
		fail(c, http.StatusNotFound, "movie not found")
		return
	}
	title := req.Title.Or(movie.GetTitle())
	description := req.Description.Or(movie.GetDescription())
	genre := req.Genre.Or(movie.GetGenre())
	duration := req.Duration.Or(movie.GetDuration())
	rating := req.Rating.Or(movie.GetRating())

	ctx, cancel = s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.UpdateMovie(ctx, &moviepb.UpdateMovieRequest{
		Id:            c.Param("id"),
		Title:         title,
		Description:   description,
		Genre:         genre,
		Duration:      duration,
		Rating:        rating,
		RequesterRole: movieRole(user.Role),
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, movieFromProto(resp.GetMovie()))
}

func (s *Server) adminDeleteMovie(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	_, err := s.clients.Movie.DeleteMovie(ctx, &moviepb.DeleteMovieRequest{
		Id:            c.Param("id"),
		RequesterRole: movieRole(user.Role),
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, gin.H{})
}

func (s *Server) adminCreateHall(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req struct {
		Name     string `json:"name"`
		Capacity int32  `json:"capacity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.CreateHall(ctx, &moviepb.CreateHallRequest{
		Name:          req.Name,
		Capacity:      req.Capacity,
		RequesterRole: movieRole(user.Role),
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	created(c, hallFromProto(resp.GetHall()))
}

func (s *Server) adminCreateSession(c *gin.Context) {
	var req struct {
		MovieID   string  `json:"movie_id"`
		HallID    string  `json:"hall_id"`
		StartTime string  `json:"start_time"`
		Price     float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid start_time")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.CreateSession(ctx, &moviepb.CreateSessionRequest{
		MovieId:   req.MovieID,
		HallId:    req.HallID,
		StartTime: timestamppb.New(start.UTC()),
		Price:     req.Price,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	created(c, sessionFromProto(resp.GetSession()))
}

func (s *Server) adminListBookings(c *gin.Context) {
	limit, offset := limitOffset(c)
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.AdminListBookings(ctx, &bookingpb.AdminListRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, bookingsFromProto(resp.GetBookings()))
}

func (s *Server) adminBookingStats(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Booking.GetBookingStats(ctx, &bookingpb.GetStatsRequest{})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, statsJSON{
		Total:     resp.GetTotal(),
		Confirmed: resp.GetConfirmed(),
		Cancelled: resp.GetCancelled(),
	})
}

type moviePayload struct {
	Title       optionalString  `json:"title"`
	Description optionalString  `json:"description"`
	Genre       optionalString  `json:"genre"`
	Duration    optionalInt32   `json:"duration"`
	Rating      optionalFloat64 `json:"rating"`
}

type optionalString struct {
	value string
	set   bool
}

func (o *optionalString) UnmarshalJSON(b []byte) error {
	o.set = true
	if string(b) == "null" {
		o.value = ""
		return nil
	}
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	o.value = value
	return nil
}

func (o optionalString) Value() string {
	return o.value
}

func (o optionalString) Or(fallback string) string {
	if o.set {
		return o.value
	}
	return fallback
}

type optionalInt32 struct {
	value int32
	set   bool
}

func (o *optionalInt32) UnmarshalJSON(b []byte) error {
	o.set = true
	if string(b) == "null" {
		o.value = 0
		return nil
	}
	var value int32
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	o.value = value
	return nil
}

func (o optionalInt32) Value() int32 {
	return o.value
}

func (o optionalInt32) Or(fallback int32) int32 {
	if o.set {
		return o.value
	}
	return fallback
}

type optionalFloat64 struct {
	value float64
	set   bool
}

func (o *optionalFloat64) UnmarshalJSON(b []byte) error {
	o.set = true
	if string(b) == "null" {
		o.value = 0
		return nil
	}
	var value float64
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	o.value = value
	return nil
}

func (o optionalFloat64) Value() float64 {
	return o.value
}

func (o optionalFloat64) Or(fallback float64) float64 {
	if o.set {
		return o.value
	}
	return fallback
}
