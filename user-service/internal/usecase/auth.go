package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"time"

	"github.com/cinema-booking-system/user-service/internal/config"
	"github.com/cinema-booking-system/user-service/internal/domain"
	apperrors "github.com/cinema-booking-system/user-service/internal/errors"
	"github.com/cinema-booking-system/user-service/internal/security"
	"github.com/cinema-booking-system/user-service/internal/validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

type Claims struct {
	jwt.RegisteredClaims
	UserID string      `json:"user_id"`
	Email  string      `json:"email"`
	Role   domain.Role `json:"role"`
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type AuthUseCase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	cfg       config.JWTConfig
	log       *slog.Logger
}

func NewAuthUseCase(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	cfg config.JWTConfig,
	log *slog.Logger,
) *AuthUseCase {
	if log == nil {
		log = slog.Default()
	}
	return &AuthUseCase{userRepo: userRepo, tokenRepo: tokenRepo, cfg: cfg, log: log}
}

func (uc *AuthUseCase) Register(ctx context.Context, input domain.RegisterInput) (*domain.User, *domain.TokenPair, error) {
	email := validation.NormalizeEmail(input.Email)
	if !emailRegex.MatchString(email) {
		return nil, nil, domain.ErrInvalidEmail
	}
	input.Email = email
	if len(input.Password) < 8 {
		return nil, nil, domain.ErrWeakPassword
	}
	if input.Password != input.ConfirmPassword {
		return nil, nil, domain.ErrPasswordMismatch
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcryptCost)
	if err != nil {
		return nil, nil, apperrors.Wrap("register: hash password", err)
	}
	user := &domain.User{
		Email:        email,
		PasswordHash: string(hash),
		FullName:     input.FullName,
		Phone:        input.Phone,
		Role:         domain.RoleUser,
	}
	created, err := uc.userRepo.Create(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	tokens, err := uc.generateTokenPair(ctx, created)
	if err != nil {
		return nil, nil, err
	}
	return created, tokens, nil
}

func (uc *AuthUseCase) Login(ctx context.Context, input domain.LoginInput) (*domain.User, *domain.TokenPair, error) {
	email := validation.NormalizeEmail(input.Email)
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		_ = security.VerifyPassword(security.DummyPasswordHash(), input.Password)
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, apperrors.Wrap("login: get user", err)
	}
	if user.IsBanned {
		return nil, nil, domain.ErrUserBanned
	}
	if !security.VerifyPassword(user.PasswordHash, input.Password) {
		return nil, nil, domain.ErrInvalidCredentials
	}
	tokens, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return user, tokens, nil
}

func (uc *AuthUseCase) Logout(ctx context.Context, accessToken, refreshToken string) error {
	if accessToken == "" {
		return nil
	}
	claims, err := uc.parseToken(accessToken)
	if err != nil {
		return nil
	}
	remaining := time.Until(claims.ExpiresAt.Time)
	if remaining > 0 {
		if err = uc.tokenRepo.AddToBlacklist(ctx, accessToken, remaining); err != nil {
			return apperrors.Wrap("logout: blacklist token", err)
		}
	}
	if refreshToken != "" {
		_ = uc.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	}
	return nil
}

func (uc *AuthUseCase) ValidateToken(ctx context.Context, accessToken string) (*domain.ValidateTokenOutput, error) {
	blacklisted, err := uc.tokenRepo.IsBlacklisted(ctx, accessToken)
	if err != nil {
		return nil, apperrors.Wrap("validate token: check blacklist", err)
	}
	if blacklisted {
		uc.log.Warn("jwt validation failed", "reason", "blacklisted")
		return &domain.ValidateTokenOutput{Valid: false}, domain.ErrTokenBlacklisted
	}
	claims, err := uc.parseToken(accessToken)
	if err != nil {
		uc.log.Warn("jwt validation failed", "reason", "parse_or_claims")
		return &domain.ValidateTokenOutput{Valid: false}, domain.ErrInvalidToken
	}
	return &domain.ValidateTokenOutput{
		Valid:  true,
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	userID, err := uc.tokenRepo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperrors.Wrap("refresh token: get user", err)
	}
	if user.IsBanned {
		return nil, domain.ErrUserBanned
	}
	_ = uc.tokenRepo.DeleteRefreshToken(ctx, refreshToken)
	tokens, err := uc.generateTokenPair(ctx, user)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *AuthUseCase) generateTokenPair(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	now := time.Now()
	accessClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    uc.cfg.Issuer,
			Audience:  jwt.ClaimStrings{uc.cfg.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.cfg.AccessTTL)),
			ID:        uuid.NewString(),
		},
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(uc.cfg.Secret))
	if err != nil {
		return nil, apperrors.Wrap("generate access token", err)
	}
	refreshClaims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    uc.cfg.Issuer,
			Audience:  jwt.ClaimStrings{uc.cfg.Audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(uc.cfg.RefreshTTL)),
			ID:        uuid.NewString(),
		},
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(uc.cfg.Secret))
	if err != nil {
		return nil, apperrors.Wrap("generate refresh token", err)
	}
	if err = uc.tokenRepo.SaveRefreshToken(ctx, user.ID, refreshToken, uc.cfg.RefreshTTL); err != nil {
		return nil, apperrors.Wrap("save refresh token", err)
	}
	return &domain.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (uc *AuthUseCase) parseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(uc.cfg.Secret), nil
		},
		jwt.WithIssuer(uc.cfg.Issuer),
		jwt.WithAudience(uc.cfg.Audience),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}
	return claims, nil
}
