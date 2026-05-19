package httpgateway

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	bookingpb "booking-service/proto/booking"
	"github.com/cinema-booking-system/api-gateway/internal/config"
	"github.com/cinema-booking-system/api-gateway/internal/repository"
	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestListMoviesClampsFrontendPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	movies := &fakeMovieClient{}
	router := testRouter(repository.Clients{User: &fakeUserClient{}, Movie: movies, Booking: &fakeBookingClient{}})

	res := perform(router, http.MethodGet, "/movies?page=2&limit=500", "", nil)
	if res.Code != http.StatusOK {
		t.Fatalf("status = %d", res.Code)
	}
	if movies.lastListMovies.GetLimit() != 100 {
		t.Fatalf("limit = %d", movies.lastListMovies.GetLimit())
	}
	if movies.lastListMovies.GetOffset() != 100 {
		t.Fatalf("offset = %d", movies.lastListMovies.GetOffset())
	}
}

func TestGRPCNotFoundMapsToHTTP404(t *testing.T) {
	gin.SetMode(gin.TestMode)
	movies := &fakeMovieClient{getMovieErr: status.Error(codes.NotFound, "movie not found")}
	router := testRouter(repository.Clients{User: &fakeUserClient{}, Movie: movies, Booking: &fakeBookingClient{}})

	res := perform(router, http.MethodGet, "/movies/missing", "", nil)
	if res.Code != http.StatusNotFound {
		t.Fatalf("status = %d", res.Code)
	}
}

func TestAdminRouteRejectsNonAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := &fakeUserClient{role: userpb.Role_ROLE_USER}
	router := testRouter(repository.Clients{User: users, Movie: &fakeMovieClient{}, Booking: &fakeBookingClient{}})

	res := perform(router, http.MethodGet, "/admin/bookings/stats", "Bearer token", nil)
	if res.Code != http.StatusForbidden {
		t.Fatalf("status = %d", res.Code)
	}
}

func TestOwnedBookingRejectsDifferentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := &fakeUserClient{userID: "user-1"}
	bookings := &fakeBookingClient{
		booking: &bookingpb.Booking{Id: "b1", UserId: "user-2"},
	}
	router := testRouter(repository.Clients{User: users, Movie: &fakeMovieClient{}, Booking: bookings})

	res := perform(router, http.MethodGet, "/bookings/b1", "Bearer token", nil)
	if res.Code != http.StatusForbidden {
		t.Fatalf("status = %d", res.Code)
	}
}

func TestCreateBookingInjectsAuthenticatedUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	users := &fakeUserClient{userID: "user-1", email: "u@example.com"}
	bookings := &fakeBookingClient{}
	router := testRouter(repository.Clients{User: users, Movie: &fakeMovieClient{}, Booking: bookings})

	body := []byte(`{"session_id":"s1","seat_id":"seat-1"}`)
	res := perform(router, http.MethodPost, "/bookings", "Bearer token", body)
	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d body = %s", res.Code, res.Body.String())
	}
	if bookings.lastCreate.GetUserId() != "user-1" {
		t.Fatalf("user_id = %q", bookings.lastCreate.GetUserId())
	}
	if bookings.lastCreate.GetUserEmail() != "u@example.com" {
		t.Fatalf("user_email = %q", bookings.lastCreate.GetUserEmail())
	}
}

func TestRateLimiterRejectsExcessRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testConfig()
	cfg.IPRate = 1
	cfg.IPBurst = 1
	router := NewRouter(cfg, repository.Clients{User: &fakeUserClient{}, Movie: &fakeMovieClient{}, Booking: &fakeBookingClient{}})

	first := perform(router, http.MethodGet, "/healthz", "", nil)
	if first.Code != http.StatusOK {
		t.Fatalf("first status = %d", first.Code)
	}
	second := perform(router, http.MethodGet, "/healthz", "", nil)
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d", second.Code)
	}
}

func testRouter(clients repository.Clients) *gin.Engine {
	return NewRouter(testConfig(), clients)
}

func testConfig() config.Config {
	return config.Config{
		HTTPAddr:       ":0",
		RequestTimeout: time.Second,
		AllowedOrigins: []string{"http://localhost:5173"},
		IPRate:         1000,
		IPBurst:        1000,
		UserRate:       1000,
		UserBurst:      1000,
	}
}

func perform(router http.Handler, method, path, auth string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

type fakeUserClient struct {
	userID string
	email  string
	role   userpb.Role
}

func (f *fakeUserClient) Register(context.Context, *userpb.RegisterRequest, ...grpc.CallOption) (*userpb.RegisterResponse, error) {
	return &userpb.RegisterResponse{AccessToken: "access", RefreshToken: "refresh"}, nil
}

func (f *fakeUserClient) Login(context.Context, *userpb.LoginRequest, ...grpc.CallOption) (*userpb.LoginResponse, error) {
	return &userpb.LoginResponse{AccessToken: "access", RefreshToken: "refresh"}, nil
}

func (f *fakeUserClient) Logout(context.Context, *userpb.LogoutRequest, ...grpc.CallOption) (*userpb.LogoutResponse, error) {
	return &userpb.LogoutResponse{Success: true}, nil
}

func (f *fakeUserClient) ValidateToken(context.Context, *userpb.ValidateTokenRequest, ...grpc.CallOption) (*userpb.ValidateTokenResponse, error) {
	userID := f.userID
	if userID == "" {
		userID = "user-1"
	}
	email := f.email
	if email == "" {
		email = "user@example.com"
	}
	role := f.role
	if role == userpb.Role_ROLE_UNSPECIFIED {
		role = userpb.Role_ROLE_ADMIN
	}
	return &userpb.ValidateTokenResponse{Valid: true, UserId: userID, Email: email, Role: role}, nil
}

func (f *fakeUserClient) RefreshToken(context.Context, *userpb.RefreshTokenRequest, ...grpc.CallOption) (*userpb.RefreshTokenResponse, error) {
	return &userpb.RefreshTokenResponse{AccessToken: "access-2", RefreshToken: "refresh-2"}, nil
}

func (f *fakeUserClient) ChangePassword(context.Context, *userpb.ChangePasswordRequest, ...grpc.CallOption) (*userpb.ChangePasswordResponse, error) {
	return &userpb.ChangePasswordResponse{Success: true}, nil
}

func (f *fakeUserClient) GetAllUsers(context.Context, *userpb.GetAllUsersRequest, ...grpc.CallOption) (*userpb.GetAllUsersResponse, error) {
	return &userpb.GetAllUsersResponse{}, nil
}

func (f *fakeUserClient) BanUser(context.Context, *userpb.BanUserRequest, ...grpc.CallOption) (*userpb.BanUserResponse, error) {
	return &userpb.BanUserResponse{}, nil
}

func (f *fakeUserClient) GetProfile(context.Context, *userpb.GetProfileRequest, ...grpc.CallOption) (*userpb.GetProfileResponse, error) {
	return &userpb.GetProfileResponse{User: &userpb.User{UserId: "user-1", Email: "user@example.com"}}, nil
}

func (f *fakeUserClient) TopUpBalance(context.Context, *userpb.TopUpBalanceRequest, ...grpc.CallOption) (*userpb.TopUpBalanceResponse, error) {
	return &userpb.TopUpBalanceResponse{User: &userpb.User{UserId: "user-1", Email: "user@example.com"}}, nil
}

func (f *fakeUserClient) DeductBalance(context.Context, *userpb.DeductBalanceRequest, ...grpc.CallOption) (*userpb.DeductBalanceResponse, error) {
	return &userpb.DeductBalanceResponse{User: &userpb.User{UserId: "user-1", Email: "user@example.com"}}, nil
}

func (f *fakeUserClient) UpdateUserRole(context.Context, *userpb.UpdateUserRoleRequest, ...grpc.CallOption) (*userpb.UpdateUserRoleResponse, error) {
	return &userpb.UpdateUserRoleResponse{User: &userpb.User{UserId: "user-1", Email: "user@example.com"}}, nil
}

type fakeMovieClient struct {
	lastListMovies *moviepb.ListMoviesRequest
	getMovieErr    error
}

func (f *fakeMovieClient) CreateMovie(context.Context, *moviepb.CreateMovieRequest, ...grpc.CallOption) (*moviepb.CreateMovieResponse, error) {
	return &moviepb.CreateMovieResponse{Movie: &moviepb.Movie{Id: "m1"}}, nil
}

func (f *fakeMovieClient) GetMovie(context.Context, *moviepb.GetMovieRequest, ...grpc.CallOption) (*moviepb.GetMovieResponse, error) {
	if f.getMovieErr != nil {
		return nil, f.getMovieErr
	}
	return &moviepb.GetMovieResponse{Movie: &moviepb.Movie{Id: "m1", Title: "Movie", CreatedAt: timestamppb.Now()}}, nil
}

func (f *fakeMovieClient) UpdateMovie(context.Context, *moviepb.UpdateMovieRequest, ...grpc.CallOption) (*moviepb.UpdateMovieResponse, error) {
	return &moviepb.UpdateMovieResponse{Movie: &moviepb.Movie{Id: "m1"}}, nil
}

func (f *fakeMovieClient) DeleteMovie(context.Context, *moviepb.DeleteMovieRequest, ...grpc.CallOption) (*moviepb.DeleteMovieResponse, error) {
	return &moviepb.DeleteMovieResponse{Success: true}, nil
}

func (f *fakeMovieClient) ListMovies(_ context.Context, req *moviepb.ListMoviesRequest, _ ...grpc.CallOption) (*moviepb.ListMoviesResponse, error) {
	f.lastListMovies = req
	return &moviepb.ListMoviesResponse{}, nil
}

func (f *fakeMovieClient) CreateHall(context.Context, *moviepb.CreateHallRequest, ...grpc.CallOption) (*moviepb.CreateHallResponse, error) {
	return &moviepb.CreateHallResponse{Hall: &moviepb.Hall{Id: "h1"}}, nil
}

func (f *fakeMovieClient) GetHall(context.Context, *moviepb.GetHallRequest, ...grpc.CallOption) (*moviepb.GetHallResponse, error) {
	return &moviepb.GetHallResponse{Hall: &moviepb.Hall{Id: "h1"}}, nil
}

func (f *fakeMovieClient) CreateSession(context.Context, *moviepb.CreateSessionRequest, ...grpc.CallOption) (*moviepb.CreateSessionResponse, error) {
	return &moviepb.CreateSessionResponse{Session: &moviepb.Session{Id: "s1"}}, nil
}

func (f *fakeMovieClient) GetSession(context.Context, *moviepb.GetSessionRequest, ...grpc.CallOption) (*moviepb.GetSessionResponse, error) {
	return &moviepb.GetSessionResponse{Session: &moviepb.Session{Id: "s1"}}, nil
}

func (f *fakeMovieClient) ListSessions(context.Context, *moviepb.ListSessionsRequest, ...grpc.CallOption) (*moviepb.ListSessionsResponse, error) {
	return &moviepb.ListSessionsResponse{}, nil
}

func (f *fakeMovieClient) GetAvailableSeats(context.Context, *moviepb.GetAvailableSeatsRequest, ...grpc.CallOption) (*moviepb.GetAvailableSeatsResponse, error) {
	return &moviepb.GetAvailableSeatsResponse{}, nil
}

type fakeBookingClient struct {
	booking    *bookingpb.Booking
	lastCreate *bookingpb.CreateBookingRequest
}

func (f *fakeBookingClient) CreateBooking(_ context.Context, req *bookingpb.CreateBookingRequest, _ ...grpc.CallOption) (*bookingpb.CreateBookingResponse, error) {
	f.lastCreate = req
	return &bookingpb.CreateBookingResponse{Booking: &bookingpb.Booking{Id: "b1", UserId: req.GetUserId()}}, nil
}

func (f *fakeBookingClient) GetBooking(context.Context, *bookingpb.GetBookingRequest, ...grpc.CallOption) (*bookingpb.GetBookingResponse, error) {
	booking := f.booking
	if booking == nil {
		booking = &bookingpb.Booking{Id: "b1", UserId: "user-1"}
	}
	return &bookingpb.GetBookingResponse{Booking: booking}, nil
}

func (f *fakeBookingClient) ListUserBookings(context.Context, *bookingpb.ListUserBookingsRequest, ...grpc.CallOption) (*bookingpb.ListUserBookingsResponse, error) {
	return &bookingpb.ListUserBookingsResponse{}, nil
}

func (f *fakeBookingClient) CancelBooking(context.Context, *bookingpb.CancelBookingRequest, ...grpc.CallOption) (*bookingpb.CancelBookingResponse, error) {
	return &bookingpb.CancelBookingResponse{Success: true}, nil
}

func (f *fakeBookingClient) ConfirmPayment(context.Context, *bookingpb.ConfirmPaymentRequest, ...grpc.CallOption) (*bookingpb.ConfirmPaymentResponse, error) {
	return &bookingpb.ConfirmPaymentResponse{Payment: &bookingpb.Payment{Id: "p1", BookingId: "b1"}}, nil
}

func (f *fakeBookingClient) GetPayment(context.Context, *bookingpb.GetPaymentRequest, ...grpc.CallOption) (*bookingpb.GetPaymentResponse, error) {
	return &bookingpb.GetPaymentResponse{Payment: &bookingpb.Payment{Id: "p1", BookingId: "b1"}}, nil
}

func (f *fakeBookingClient) GetBookingHistory(context.Context, *bookingpb.GetHistoryRequest, ...grpc.CallOption) (*bookingpb.GetHistoryResponse, error) {
	return &bookingpb.GetHistoryResponse{}, nil
}

func (f *fakeBookingClient) AdminListBookings(context.Context, *bookingpb.AdminListRequest, ...grpc.CallOption) (*bookingpb.AdminListResponse, error) {
	return &bookingpb.AdminListResponse{}, nil
}

func (f *fakeBookingClient) GetBookingStats(context.Context, *bookingpb.GetStatsRequest, ...grpc.CallOption) (*bookingpb.GetStatsResponse, error) {
	return &bookingpb.GetStatsResponse{}, nil
}

func (f *fakeBookingClient) GetSessionTakenSeats(context.Context, *bookingpb.GetSessionTakenSeatsRequest, ...grpc.CallOption) (*bookingpb.GetSessionTakenSeatsResponse, error) {
	return &bookingpb.GetSessionTakenSeatsResponse{}, nil
}
