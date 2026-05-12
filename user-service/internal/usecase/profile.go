package usecase

import (
	"context"
	"math"

	"github.com/cinema-booking-system/user-service/internal/domain"
	apperrors "github.com/cinema-booking-system/user-service/internal/errors"
	"github.com/cinema-booking-system/user-service/internal/security"
	"github.com/cinema-booking-system/user-service/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

type ProfileUseCase struct {
	userRepo  domain.UserRepository
	cache     domain.UserCache
	tokenRepo domain.TokenRepository
}

func NewProfileUseCase(
	userRepo domain.UserRepository,
	cache domain.UserCache,
	tokenRepo domain.TokenRepository,
) *ProfileUseCase {
	return &ProfileUseCase{
		userRepo:  userRepo,
		cache:     cache,
		tokenRepo: tokenRepo,
	}
}

func (uc *ProfileUseCase) requireAdmin(ctx context.Context, adminID string) (*domain.User, error) {
	id, err := validation.ParseUserID(adminID)
	if err != nil {
		return nil, err
	}
	admin, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if admin.Role != domain.RoleAdmin {
		return nil, domain.ErrForbidden
	}
	return admin, nil
}

func (uc *ProfileUseCase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	id, err := validation.ParseUserID(userID)
	if err != nil {
		return nil, err
	}
	if user, err := uc.cache.GetUserCache(ctx, id); err == nil {
		return user, nil
	}
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = uc.cache.SetUserCache(ctx, user)
	return user, nil
}

func (uc *ProfileUseCase) UpdateProfile(ctx context.Context, input domain.UpdateProfileInput) (*domain.User, error) {
	id, err := validation.ParseUserID(input.UserID)
	if err != nil {
		return nil, err
	}
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.FullName != "" {
		user.FullName = input.FullName
	}
	if input.Phone != "" {
		user.Phone = input.Phone
	}
	updated, err := uc.userRepo.Update(ctx, user)
	if err != nil {
		return nil, apperrors.Wrap("update profile", err)
	}
	_ = uc.cache.DeleteUserCache(ctx, id)
	return updated, nil
}

func (uc *ProfileUseCase) DeleteAccount(ctx context.Context, userID, password string) error {
	id, err := validation.ParseUserID(userID)
	if err != nil {
		return err
	}
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !security.VerifyPassword(user.PasswordHash, password) {
		return domain.ErrInvalidCredentials
	}
	if err = uc.userRepo.Delete(ctx, id); err != nil {
		return apperrors.Wrap("delete account", err)
	}
	_ = uc.cache.DeleteUserCache(ctx, id)
	_ = uc.tokenRepo.DeleteAllUserTokens(ctx, id)
	return nil
}

func (uc *ProfileUseCase) ChangePassword(ctx context.Context, input domain.ChangePasswordInput) error {
	id, err := validation.ParseUserID(input.UserID)
	if err != nil {
		return err
	}
	if input.NewPassword != input.ConfirmPassword {
		return domain.ErrPasswordMismatch
	}
	if len(input.NewPassword) < 8 {
		return domain.ErrWeakPassword
	}
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if !security.VerifyPassword(user.PasswordHash, input.OldPassword) {
		return domain.ErrInvalidCredentials
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcryptCost)
	if err != nil {
		return apperrors.Wrap("change password: hash", err)
	}
	return uc.userRepo.UpdatePassword(ctx, id, string(newHash))
}

func (uc *ProfileUseCase) GetAllUsers(ctx context.Context, input domain.GetAllUsersInput) (*domain.GetAllUsersOutput, error) {
	if _, err := uc.requireAdmin(ctx, input.AdminID); err != nil {
		return nil, err
	}
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 20
	}
	users, total, err := uc.userRepo.GetAll(ctx, input.Page, input.PageSize, input.RoleFilter, input.BannedOnly)
	if err != nil {
		return nil, apperrors.Wrap("get all users", err)
	}
	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))
	return &domain.GetAllUsersOutput{
		Users:      users,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       input.Page,
		PageSize:   input.PageSize,
	}, nil
}

func (uc *ProfileUseCase) BanUser(ctx context.Context, adminID, userID string, ban bool, reason string) (*domain.User, error) {
	admin, err := uc.requireAdmin(ctx, adminID)
	if err != nil {
		return nil, err
	}
	targetID, err := validation.ParseUserID(userID)
	if err != nil {
		return nil, err
	}
	if ban && admin.ID == targetID {
		return nil, domain.ErrSelfBanForbidden
	}
	user, err := uc.userRepo.Ban(ctx, targetID, ban, reason)
	if err != nil {
		return nil, apperrors.Wrap("ban user", err)
	}
	_ = uc.cache.DeleteUserCache(ctx, targetID)
	if ban {
		_ = uc.tokenRepo.DeleteAllUserTokens(ctx, targetID)
	}
	return user, nil
}

func (uc *ProfileUseCase) GetUserByEmail(ctx context.Context, input domain.GetUserByEmailInput) (*domain.User, error) {
	if _, err := uc.requireAdmin(ctx, input.AdminID); err != nil {
		return nil, err
	}
	email := validation.NormalizeEmail(input.Email)
	if !emailRegex.MatchString(email) {
		return nil, domain.ErrInvalidEmail
	}
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
