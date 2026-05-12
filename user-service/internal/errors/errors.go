package apperrors

import (
	"errors"
	"fmt"

	"github.com/cinema-booking-system/user-service/internal/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}

func ToGRPC(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return status.Error(codes.NotFound, domain.ErrUserNotFound.Error())
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, domain.ErrEmailAlreadyExists.Error())
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, domain.ErrInvalidCredentials.Error())
	case errors.Is(err, domain.ErrUserBanned):
		return status.Error(codes.PermissionDenied, domain.ErrUserBanned.Error())
	case errors.Is(err, domain.ErrForbidden):
		return status.Error(codes.PermissionDenied, domain.ErrForbidden.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return status.Error(codes.Unauthenticated, domain.ErrUnauthorized.Error())
	case errors.Is(err, domain.ErrInvalidToken),
		errors.Is(err, domain.ErrTokenBlacklisted):
		return status.Error(codes.Unauthenticated, "invalid or revoked token")
	case errors.Is(err, domain.ErrWeakPassword):
		return status.Error(codes.InvalidArgument, domain.ErrWeakPassword.Error())
	case errors.Is(err, domain.ErrInvalidEmail):
		return status.Error(codes.InvalidArgument, domain.ErrInvalidEmail.Error())
	case errors.Is(err, domain.ErrPasswordMismatch):
		return status.Error(codes.InvalidArgument, domain.ErrPasswordMismatch.Error())
	case errors.Is(err, domain.ErrInvalidUserID):
		return status.Error(codes.InvalidArgument, domain.ErrInvalidUserID.Error())
	case errors.Is(err, domain.ErrSelfBanForbidden):
		return status.Error(codes.InvalidArgument, domain.ErrSelfBanForbidden.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
