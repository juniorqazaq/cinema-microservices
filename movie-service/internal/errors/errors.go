package apperrors

import (
	"errors"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrMovieNotFound):
		return status.Error(codes.NotFound, domain.ErrMovieNotFound.Error())
	case errors.Is(err, domain.ErrSessionNotFound):
		return status.Error(codes.NotFound, domain.ErrSessionNotFound.Error())
	case errors.Is(err, domain.ErrHallNotFound):
		return status.Error(codes.NotFound, domain.ErrHallNotFound.Error())
	case errors.Is(err, domain.ErrForbidden):
		return status.Error(codes.PermissionDenied, domain.ErrForbidden.Error())
	case errors.Is(err, domain.ErrInvalidArgument),
		errors.Is(err, domain.ErrInvalidPagination):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
