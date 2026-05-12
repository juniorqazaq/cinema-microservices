package validation

import (
	"strings"

	"github.com/cinema-booking-system/user-service/internal/domain"
	"github.com/google/uuid"
)

func ParseUserID(id string) (string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", domain.ErrInvalidUserID
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return "", domain.ErrInvalidUserID
	}
	return parsed.String(), nil
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
