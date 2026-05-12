package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cinema-booking-system/user-service/internal/config"
	"github.com/cinema-booking-system/user-service/internal/domain"
	"github.com/cinema-booking-system/user-service/internal/repository/mocks"
	"github.com/cinema-booking-system/user-service/internal/usecase"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func jwtCfg() config.JWTConfig {
	return config.JWTConfig{
		Secret:     "test-secret-key-32-bytes-long-xx",
		AccessTTL:  15 * time.Minute,
		RefreshTTL: 7 * 24 * time.Hour,
		Issuer:     "cinema-booking-system",
		Audience:   "cinema-users",
	}
}

const testUserUUID = "550e8400-e29b-41d4-a716-446655440000"

func newAuthUC(t *testing.T) (*usecase.AuthUseCase, *mocks.MockUserRepository, *mocks.MockTokenRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenRepo := mocks.NewMockTokenRepository(ctrl)
	uc := usecase.NewAuthUseCase(userRepo, tokenRepo, jwtCfg(), nil)
	return uc, userRepo, tokenRepo
}

func sampleUser() *domain.User {
	return &domain.User{
		ID:    testUserUUID,
		Email: "john@example.com",
		Role:  domain.RoleUser,
	}
}

func TestRegister_Success(t *testing.T) {
	uc, userRepo, tokenRepo := newAuthUC(t)
	ctx := context.Background()
	input := domain.RegisterInput{
		FullName:        "John Doe",
		Email:           "john@example.com",
		Phone:           "+1234567890",
		Password:        "securePass123",
		ConfirmPassword: "securePass123",
	}
	created := sampleUser()
	created.FullName = input.FullName
	userRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(created, nil)
	tokenRepo.EXPECT().
		SaveRefreshToken(ctx, created.ID, gomock.Any(), gomock.Any()).
		Return(nil)
	user, tokens, err := uc.Register(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if user.Email != input.Email {
		t.Errorf("email mismatch: got %s want %s", user.Email, input.Email)
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		t.Error("expected non-empty tokens")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	uc, userRepo, _ := newAuthUC(t)
	ctx := context.Background()
	input := domain.RegisterInput{
		Email:           "john@example.com",
		Password:        "securePass123",
		ConfirmPassword: "securePass123",
	}
	userRepo.EXPECT().
		Create(ctx, gomock.Any()).
		Return(nil, domain.ErrEmailAlreadyExists)
	_, _, err := uc.Register(ctx, input)
	if !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	uc, _, _ := newAuthUC(t)
	_, _, err := uc.Register(context.Background(), domain.RegisterInput{
		Email:           "a@b.com",
		Password:        "short",
		ConfirmPassword: "short",
	})
	if !errors.Is(err, domain.ErrWeakPassword) {
		t.Fatalf("expected ErrWeakPassword, got: %v", err)
	}
}

func TestRegister_PasswordMismatch(t *testing.T) {
	uc, _, _ := newAuthUC(t)
	_, _, err := uc.Register(context.Background(), domain.RegisterInput{
		Email:           "a@b.com",
		Password:        "securePass123",
		ConfirmPassword: "differentPass",
	})
	if !errors.Is(err, domain.ErrPasswordMismatch) {
		t.Fatalf("expected ErrPasswordMismatch, got: %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	uc, userRepo, _ := newAuthUC(t)
	ctx := context.Background()
	user := sampleUser()
	user.PasswordHash = "$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQyCrxQ6NQ.P6v6y5K1XLe8Fm"
	userRepo.EXPECT().
		GetByEmail(ctx, "john@example.com").
		Return(user, nil)
	_, _, err := uc.Login(ctx, domain.LoginInput{
		Email:    "john@example.com",
		Password: "wrongPassword",
	})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_BannedUser(t *testing.T) {
	uc, userRepo, _ := newAuthUC(t)
	ctx := context.Background()
	banned := sampleUser()
	banned.IsBanned = true
	userRepo.EXPECT().
		GetByEmail(ctx, "john@example.com").
		Return(banned, nil)
	_, _, err := uc.Login(ctx, domain.LoginInput{
		Email:    "john@example.com",
		Password: "anyPassword",
	})
	if !errors.Is(err, domain.ErrUserBanned) {
		t.Fatalf("expected ErrUserBanned, got: %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	uc, userRepo, _ := newAuthUC(t)
	ctx := context.Background()
	userRepo.EXPECT().
		GetByEmail(ctx, "ghost@example.com").
		Return(nil, domain.ErrUserNotFound)
	_, _, err := uc.Login(ctx, domain.LoginInput{Email: "ghost@example.com", Password: "pass"})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials for unknown email, got: %v", err)
	}
}

func TestValidateToken_Blacklisted(t *testing.T) {
	uc, _, tokenRepo := newAuthUC(t)
	ctx := context.Background()
	tokenRepo.EXPECT().
		IsBlacklisted(ctx, "some.token.here").
		Return(true, nil)
	out, err := uc.ValidateToken(ctx, "some.token.here")
	if !errors.Is(err, domain.ErrTokenBlacklisted) {
		t.Fatalf("expected ErrTokenBlacklisted, got: %v", err)
	}
	if out.Valid {
		t.Error("expected valid=false for blacklisted token")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	uc, _, tokenRepo := newAuthUC(t)
	ctx := context.Background()
	expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMSIsImV4cCI6MX0.invalid"
	tokenRepo.EXPECT().
		IsBlacklisted(ctx, expiredToken).
		Return(false, nil)
	out, err := uc.ValidateToken(ctx, expiredToken)
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken for expired token, got: %v", err)
	}
	if out.Valid {
		t.Error("expected valid=false for expired token")
	}
}

func TestRefreshToken_Expired(t *testing.T) {
	uc, _, tokenRepo := newAuthUC(t)
	ctx := context.Background()
	tokenRepo.EXPECT().
		FindRefreshToken(ctx, "expired-or-missing-refresh").
		Return("", domain.ErrInvalidToken)
	_, err := uc.RefreshToken(ctx, "expired-or-missing-refresh")
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestLogout_BlacklistsAccessToken(t *testing.T) {
	uc, _, tokenRepo := newAuthUC(t)
	ctx := context.Background()
	cfg := jwtCfg()
	now := time.Now()
	accessClaims := &usecase.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   testUserUUID,
			Issuer:    cfg.Issuer,
			Audience:  jwt.ClaimStrings{cfg.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			ID:        uuid.NewString(),
		},
		UserID: testUserUUID,
		Email:  "a@b.com",
		Role:   domain.RoleUser,
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(cfg.Secret))
	if err != nil {
		t.Fatal(err)
	}
	refresh := "refresh-opaque-value"
	tokenRepo.EXPECT().AddToBlacklist(ctx, accessToken, gomock.Any()).Return(nil)
	tokenRepo.EXPECT().DeleteRefreshToken(ctx, refresh).Return(nil)
	if err := uc.Logout(ctx, accessToken, refresh); err != nil {
		t.Fatal(err)
	}
}

func TestRefreshToken_BannedUser(t *testing.T) {
	uc, userRepo, tokenRepo := newAuthUC(t)
	ctx := context.Background()
	tokenRepo.EXPECT().
		FindRefreshToken(ctx, "valid.refresh.token").
		Return(testUserUUID, nil)
	banned := sampleUser()
	banned.IsBanned = true
	userRepo.EXPECT().
		GetByID(ctx, testUserUUID).
		Return(banned, nil)
	_, err := uc.RefreshToken(ctx, "valid.refresh.token")
	if !errors.Is(err, domain.ErrUserBanned) {
		t.Fatalf("expected ErrUserBanned, got: %v", err)
	}
}
