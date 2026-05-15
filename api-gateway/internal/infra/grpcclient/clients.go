package grpcclient

import (
	bookingpb "booking-service/proto/booking"

	"github.com/cinema-booking-system/api-gateway/internal/config"
	"github.com/cinema-booking-system/api-gateway/internal/repository"
	moviepb "github.com/cinema-booking-system/movie-service/gen/go/movie"
	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DialClients(cfg config.Config) (repository.Clients, func(), error) {
	userConn, err := grpc.NewClient(cfg.UserGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return repository.Clients{}, func() {}, err
	}
	movieConn, err := grpc.NewClient(cfg.MovieGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = userConn.Close()
		return repository.Clients{}, func() {}, err
	}
	bookingConn, err := grpc.NewClient(cfg.BookingGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = userConn.Close()
		_ = movieConn.Close()
		return repository.Clients{}, func() {}, err
	}

	closeFn := func() {
		_ = userConn.Close()
		_ = movieConn.Close()
		_ = bookingConn.Close()
	}

	return repository.Clients{
		User:    userpb.NewUserServiceClient(userConn),
		Movie:   moviepb.NewMovieServiceClient(movieConn),
		Booking: bookingpb.NewBookingServiceClient(bookingConn),
	}, closeFn, nil
}
