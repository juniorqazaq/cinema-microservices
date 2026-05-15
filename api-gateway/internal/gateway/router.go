package gateway

import (
	"context"
	"net/http"
	"strconv"
	"time"

	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	cfg     Config
	clients Clients
}

func NewRouter(cfg Config, clients Clients) *gin.Engine {
	s := &Server{cfg: cfg, clients: clients}
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(cfg.AllowedOrigins))
	r.Use(rateLimitMiddleware(NewRateLimiter(cfg.IPRate, cfg.IPBurst), func(c *gin.Context) string {
		return c.ClientIP()
	}))

	r.GET("/healthz", func(c *gin.Context) { ok(c, gin.H{"status": "ok"}) })

	r.POST("/auth/register", s.register)
	r.POST("/auth/login", s.login)
	r.POST("/auth/refresh", s.refreshToken)

	r.GET("/movies", s.listMovies)
	r.GET("/movies/:id", s.getMovie)
	r.GET("/sessions", s.listSessions)
	r.GET("/sessions/:id", s.getSession)
	r.GET("/sessions/:id/seats", s.getAvailableSeats)
	r.GET("/halls/:id", s.getHall)

	protected := r.Group("/")
	protected.Use(authMiddleware(cfg, clients.User))
	protected.Use(rateLimitMiddleware(NewRateLimiter(cfg.UserRate, cfg.UserBurst), func(c *gin.Context) string {
		if user, ok := currentUser(c); ok {
			return user.ID
		}
		return c.ClientIP()
	}))
	protected.POST("/auth/logout", s.logout)
	protected.PUT("/users/password", s.changePassword)
	protected.POST("/bookings", s.createBooking)
	protected.GET("/bookings", s.listUserBookings)
	protected.GET("/bookings/history", s.getBookingHistory)
	protected.GET("/bookings/:id", s.getBooking)
	protected.DELETE("/bookings/:id", s.cancelBooking)
	protected.POST("/payments/confirm", s.confirmPayment)
	protected.GET("/payments/:id", s.getPayment)

	admin := protected.Group("/admin")
	admin.Use(adminMiddleware())
	admin.GET("/users", s.adminListUsers)
	admin.POST("/users/:id/ban", s.adminBanUser)
	admin.POST("/movies", s.adminCreateMovie)
	admin.PUT("/movies/:id", s.adminUpdateMovie)
	admin.DELETE("/movies/:id", s.adminDeleteMovie)
	admin.POST("/halls", s.adminCreateHall)
	admin.POST("/sessions", s.adminCreateSession)
	admin.GET("/bookings/stats", s.adminBookingStats)
	admin.GET("/bookings", s.adminListBookings)

	return r
}

func (s *Server) requestContext(c *gin.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.Request.Context(), s.cfg.RequestTimeout)
}

func (s *Server) register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Phone    string `json:"phone"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.Register(ctx, &userpb.RegisterRequest{
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.Password,
		FullName:        req.FullName,
		Phone:           req.Phone,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	created(c, authTokens{AccessToken: resp.GetAccessToken(), RefreshToken: resp.GetRefreshToken()})
}

func (s *Server) login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.Login(ctx, &userpb.LoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, authTokens{AccessToken: resp.GetAccessToken(), RefreshToken: resp.GetRefreshToken()})
}

func (s *Server) logout(c *gin.Context) {
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.AccessToken == "" {
		req.AccessToken = bearerToken(c.GetHeader("Authorization"))
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	if _, err := s.clients.User.Logout(ctx, &userpb.LogoutRequest{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	}); err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, gin.H{})
}

func (s *Server) refreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.User.RefreshToken(ctx, &userpb.RefreshTokenRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, authTokens{AccessToken: resp.GetAccessToken(), RefreshToken: resp.GetRefreshToken()})
}

func (s *Server) changePassword(c *gin.Context) {
	user, okUser := currentUser(c)
	if !okUser {
		fail(c, http.StatusUnauthorized, "missing authenticated user")
		return
	}
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	_, err := s.clients.User.ChangePassword(ctx, &userpb.ChangePasswordRequest{
		UserId:          user.ID,
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.NewPassword,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, gin.H{})
}

func (s *Server) listMovies(c *gin.Context) {
	limit, offset := limitOffset(c)
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.ListMovies(ctx, &moviepb.ListMoviesRequest{
		Limit:  limit,
		Offset: offset,
		Genre:  c.Query("genre"),
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	out := make([]movieJSON, 0, len(resp.GetMovies()))
	for _, m := range resp.GetMovies() {
		out = append(out, movieFromProto(m))
	}
	ok(c, out)
}

func (s *Server) getMovie(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.GetMovie(ctx, &moviepb.GetMovieRequest{Id: c.Param("id")})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, movieFromProto(resp.GetMovie()))
}

func (s *Server) listSessions(c *gin.Context) {
	limit, offset := limitOffset(c)
	var date *timestamppb.Timestamp
	if raw := c.Query("date"); raw != "" {
		day, err := time.Parse("2006-01-02", raw)
		if err != nil {
			fail(c, http.StatusBadRequest, "invalid date")
			return
		}
		date = timestamppb.New(day)
	}
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.ListSessions(ctx, &moviepb.ListSessionsRequest{
		MovieId: c.Query("movie_id"),
		Date:    date,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		failGRPC(c, err)
		return
	}
	out := make([]sessionJSON, 0, len(resp.GetSessions()))
	for _, session := range resp.GetSessions() {
		out = append(out, sessionFromProto(session))
	}
	ok(c, out)
}

func (s *Server) getSession(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.GetSession(ctx, &moviepb.GetSessionRequest{Id: c.Param("id")})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, sessionFromProto(resp.GetSession()))
}

func (s *Server) getAvailableSeats(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.GetAvailableSeats(ctx, &moviepb.GetAvailableSeatsRequest{SessionId: c.Param("id")})
	if err != nil {
		failGRPC(c, err)
		return
	}
	out := make([]seatJSON, 0, len(resp.GetSeats()))
	for _, seat := range resp.GetSeats() {
		out = append(out, seatFromProto(seat))
	}
	ok(c, out)
}

func (s *Server) getHall(c *gin.Context) {
	ctx, cancel := s.requestContext(c)
	defer cancel()
	resp, err := s.clients.Movie.GetHall(ctx, &moviepb.GetHallRequest{HallId: c.Param("id")})
	if err != nil {
		failGRPC(c, err)
		return
	}
	ok(c, hallFromProto(resp.GetHall()))
}

func limitOffset(c *gin.Context) (int32, int32) {
	page := queryInt(c, "page", 1)
	limit := queryInt(c, "limit", 20)
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return int32(limit), int32((page - 1) * limit)
}

func queryInt(c *gin.Context, name string, fallback int) int {
	if raw := c.Query(name); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			return n
		}
	}
	return fallback
}
