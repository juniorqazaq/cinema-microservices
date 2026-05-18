package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserBanned         = errors.New("user is banned")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden: admin access required")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrTokenBlacklisted   = errors.New("token has been revoked")
	ErrWeakPassword       = errors.New("password must be at least 8 characters")
	ErrInvalidEmail       = errors.New("invalid email format")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrPasswordMismatch   = errors.New("passwords do not match")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrSelfBanForbidden      = errors.New("admin cannot ban own account")
	ErrInsufficientBalance   = errors.New("insufficient balance")
)

type Role string

const (
	RoleGuest     Role = "GUEST"
	RoleUser      Role = "USER"
	RoleModerator Role = "MODERATOR"
	RoleAdmin     Role = "ADMIN"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FullName     string
	Phone        string
	Role         Role
	IsBanned     bool
	BanReason    string
	Balance      float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	UpdatePassword(ctx context.Context, userID, passwordHash string) error
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context, page, pageSize int, roleFilter string, bannedOnly bool) ([]*User, int64, error)
	Ban(ctx context.Context, userID string, ban bool, reason string) (*User, error)
	TopUpBalance(ctx context.Context, userID string, amount float64) (*User, error)
	DeductBalance(ctx context.Context, userID string, amount float64) (*User, error)
	UpdateRole(ctx context.Context, userID string, role Role) (*User, error)
}

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, userID, token string, ttl time.Duration) error
	FindRefreshToken(ctx context.Context, token string) (string, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserTokens(ctx context.Context, userID string) error
	AddToBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type UserCache interface {
	SetUserCache(ctx context.Context, user *User) error
	GetUserCache(ctx context.Context, id string) (*User, error)
	DeleteUserCache(ctx context.Context, id string) error
}

type RegisterInput struct {
	FullName        string
	Email           string
	Phone           string
	Password        string
	ConfirmPassword string
}

type LoginInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type UpdateProfileInput struct {
	UserID   string
	FullName string
	Phone    string
}

type ChangePasswordInput struct {
	UserID          string
	OldPassword     string
	NewPassword     string
	ConfirmPassword string
}

type GetAllUsersInput struct {
	AdminID    string
	Page       int
	PageSize   int
	RoleFilter string
	BannedOnly bool
}

type GetUserByEmailInput struct {
	AdminID string
	Email   string
}

type GetAllUsersOutput struct {
	Users      []*User
	TotalItems int64
	TotalPages int
	Page       int
	PageSize   int
}

type ValidateTokenOutput struct {
	Valid  bool
	UserID string
	Email  string
	Role   Role
}

type UserUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*User, *TokenPair, error)
	Login(ctx context.Context, input LoginInput) (*User, *TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	GetProfile(ctx context.Context, userID string) (*User, error)
	UpdateProfile(ctx context.Context, input UpdateProfileInput) (*User, error)
	DeleteAccount(ctx context.Context, userID, password string) error
	ValidateToken(ctx context.Context, accessToken string) (*ValidateTokenOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	ChangePassword(ctx context.Context, input ChangePasswordInput) error
	GetAllUsers(ctx context.Context, input GetAllUsersInput) (*GetAllUsersOutput, error)
	BanUser(ctx context.Context, adminID, userID string, ban bool, reason string) (*User, error)
	GetUserByEmail(ctx context.Context, input GetUserByEmailInput) (*User, error)
	TopUpBalance(ctx context.Context, userID string, amount float64) (*User, error)
	DeductBalance(ctx context.Context, userID string, amount float64) (*User, error)
	UpdateUserRole(ctx context.Context, adminID, userID string, role Role) (*User, error)
}
