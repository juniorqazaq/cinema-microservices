package gateway

import (
	"context"
	"time"

	bookingpb "booking-service/proto/booking"
	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	User    UserClient
	Movie   MovieClient
	Booking BookingClient
}

func DialClients(cfg Config) (Clients, func(), error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userConn, err := grpc.DialContext(ctx, cfg.UserGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return Clients{}, func() {}, err
	}
	movieConn, err := grpc.DialContext(ctx, cfg.MovieGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		_ = userConn.Close()
		return Clients{}, func() {}, err
	}
	bookingConn, err := grpc.DialContext(ctx, cfg.BookingGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		_ = userConn.Close()
		_ = movieConn.Close()
		return Clients{}, func() {}, err
	}

	closeFn := func() {
		_ = userConn.Close()
		_ = movieConn.Close()
		_ = bookingConn.Close()
	}

	return Clients{
		User:    userpb.NewUserServiceClient(userConn),
		Movie:   moviepb.NewMovieServiceClient(movieConn),
		Booking: bookingpb.NewBookingServiceClient(bookingConn),
	}, closeFn, nil
}
