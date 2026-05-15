package gateway

import (
	"context"

	bookingpb "booking-service/proto/booking"
	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"google.golang.org/grpc"
)

type UserClient interface {
	Register(context.Context, *userpb.RegisterRequest, ...grpc.CallOption) (*userpb.RegisterResponse, error)
	Login(context.Context, *userpb.LoginRequest, ...grpc.CallOption) (*userpb.LoginResponse, error)
	Logout(context.Context, *userpb.LogoutRequest, ...grpc.CallOption) (*userpb.LogoutResponse, error)
	ValidateToken(context.Context, *userpb.ValidateTokenRequest, ...grpc.CallOption) (*userpb.ValidateTokenResponse, error)
	RefreshToken(context.Context, *userpb.RefreshTokenRequest, ...grpc.CallOption) (*userpb.RefreshTokenResponse, error)
	ChangePassword(context.Context, *userpb.ChangePasswordRequest, ...grpc.CallOption) (*userpb.ChangePasswordResponse, error)
	GetAllUsers(context.Context, *userpb.GetAllUsersRequest, ...grpc.CallOption) (*userpb.GetAllUsersResponse, error)
	BanUser(context.Context, *userpb.BanUserRequest, ...grpc.CallOption) (*userpb.BanUserResponse, error)
}

type MovieClient interface {
	CreateMovie(context.Context, *moviepb.CreateMovieRequest, ...grpc.CallOption) (*moviepb.CreateMovieResponse, error)
	GetMovie(context.Context, *moviepb.GetMovieRequest, ...grpc.CallOption) (*moviepb.GetMovieResponse, error)
	UpdateMovie(context.Context, *moviepb.UpdateMovieRequest, ...grpc.CallOption) (*moviepb.UpdateMovieResponse, error)
	DeleteMovie(context.Context, *moviepb.DeleteMovieRequest, ...grpc.CallOption) (*moviepb.DeleteMovieResponse, error)
	ListMovies(context.Context, *moviepb.ListMoviesRequest, ...grpc.CallOption) (*moviepb.ListMoviesResponse, error)
	CreateHall(context.Context, *moviepb.CreateHallRequest, ...grpc.CallOption) (*moviepb.CreateHallResponse, error)
	GetHall(context.Context, *moviepb.GetHallRequest, ...grpc.CallOption) (*moviepb.GetHallResponse, error)
	CreateSession(context.Context, *moviepb.CreateSessionRequest, ...grpc.CallOption) (*moviepb.CreateSessionResponse, error)
	GetSession(context.Context, *moviepb.GetSessionRequest, ...grpc.CallOption) (*moviepb.GetSessionResponse, error)
	ListSessions(context.Context, *moviepb.ListSessionsRequest, ...grpc.CallOption) (*moviepb.ListSessionsResponse, error)
	GetAvailableSeats(context.Context, *moviepb.GetAvailableSeatsRequest, ...grpc.CallOption) (*moviepb.GetAvailableSeatsResponse, error)
}

type BookingClient interface {
	CreateBooking(context.Context, *bookingpb.CreateBookingRequest, ...grpc.CallOption) (*bookingpb.CreateBookingResponse, error)
	GetBooking(context.Context, *bookingpb.GetBookingRequest, ...grpc.CallOption) (*bookingpb.GetBookingResponse, error)
	ListUserBookings(context.Context, *bookingpb.ListUserBookingsRequest, ...grpc.CallOption) (*bookingpb.ListUserBookingsResponse, error)
	CancelBooking(context.Context, *bookingpb.CancelBookingRequest, ...grpc.CallOption) (*bookingpb.CancelBookingResponse, error)
	ConfirmPayment(context.Context, *bookingpb.ConfirmPaymentRequest, ...grpc.CallOption) (*bookingpb.ConfirmPaymentResponse, error)
	GetPayment(context.Context, *bookingpb.GetPaymentRequest, ...grpc.CallOption) (*bookingpb.GetPaymentResponse, error)
	GetBookingHistory(context.Context, *bookingpb.GetHistoryRequest, ...grpc.CallOption) (*bookingpb.GetHistoryResponse, error)
	AdminListBookings(context.Context, *bookingpb.AdminListRequest, ...grpc.CallOption) (*bookingpb.AdminListResponse, error)
	GetBookingStats(context.Context, *bookingpb.GetStatsRequest, ...grpc.CallOption) (*bookingpb.GetStatsResponse, error)
}
