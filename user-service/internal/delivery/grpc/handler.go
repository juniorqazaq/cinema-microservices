package grpc

import (
	"context"
	"errors"

	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/cinema-booking-system/user-service/internal/domain"
	apperrors "github.com/cinema-booking-system/user-service/internal/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	auth    AuthFacade
	profile ProfileFacade
}

type AuthFacade interface {
	Register(ctx context.Context, input domain.RegisterInput) (*domain.User, *domain.TokenPair, error)
	Login(ctx context.Context, input domain.LoginInput) (*domain.User, *domain.TokenPair, error)
	Logout(ctx context.Context, accessToken, refreshToken string) error
	ValidateToken(ctx context.Context, accessToken string) (*domain.ValidateTokenOutput, error)
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error)
}

type ProfileFacade interface {
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, input domain.UpdateProfileInput) (*domain.User, error)
	DeleteAccount(ctx context.Context, userID, password string) error
	ChangePassword(ctx context.Context, input domain.ChangePasswordInput) error
	GetAllUsers(ctx context.Context, input domain.GetAllUsersInput) (*domain.GetAllUsersOutput, error)
	BanUser(ctx context.Context, adminID, userID string, ban bool, reason string) (*domain.User, error)
	GetUserByEmail(ctx context.Context, input domain.GetUserByEmailInput) (*domain.User, error)
}

func NewHandler(auth AuthFacade, profile ProfileFacade) *Handler {
	return &Handler{auth: auth, profile: profile}
}

func (h *Handler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	user, tokens, err := h.auth.Register(ctx, domain.RegisterInput{
		FullName:        req.FullName,
		Email:           req.Email,
		Phone:           req.Phone,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.RegisterResponse{
		User:         toProtoUser(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Message:      "registration successful",
	}, nil
}

func (h *Handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	user, tokens, err := h.auth.Login(ctx, domain.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.LoginResponse{
		User:         toProtoUser(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Message:      "login successful",
	}, nil
}

func (h *Handler) Logout(ctx context.Context, req *userpb.LogoutRequest) (*userpb.LogoutResponse, error) {
	if err := h.auth.Logout(ctx, req.AccessToken, req.RefreshToken); err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.LogoutResponse{Success: true, Message: "logged out successfully"}, nil
}

func (h *Handler) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.GetProfileResponse, error) {
	user, err := h.profile.GetProfile(ctx, req.UserId)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.GetProfileResponse{User: toProtoUser(user)}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.UpdateProfileResponse, error) {
	user, err := h.profile.UpdateProfile(ctx, domain.UpdateProfileInput{
		UserID:   req.UserId,
		FullName: req.FullName,
		Phone:    req.Phone,
	})
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.UpdateProfileResponse{User: toProtoUser(user), Message: "profile updated"}, nil
}

func (h *Handler) DeleteAccount(ctx context.Context, req *userpb.DeleteAccountRequest) (*userpb.DeleteAccountResponse, error) {
	if err := h.profile.DeleteAccount(ctx, req.UserId, req.Password); err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.DeleteAccountResponse{Success: true, Message: "account deleted"}, nil
}

func (h *Handler) ValidateToken(ctx context.Context, req *userpb.ValidateTokenRequest) (*userpb.ValidateTokenResponse, error) {
	out, err := h.auth.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		if errors.Is(err, domain.ErrTokenBlacklisted) || errors.Is(err, domain.ErrInvalidToken) {
			return &userpb.ValidateTokenResponse{Valid: false, Message: "invalid token"}, nil
		}
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.ValidateTokenResponse{
		Valid:   out.Valid,
		UserId:  out.UserID,
		Email:   out.Email,
		Role:    roleToProto(out.Role),
		Message: "token is valid",
	}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, req *userpb.RefreshTokenRequest) (*userpb.RefreshTokenResponse, error) {
	tokens, err := h.auth.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		Message:      "tokens refreshed",
	}, nil
}

func (h *Handler) ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
	if err := h.profile.ChangePassword(ctx, domain.ChangePasswordInput{
		UserID:          req.UserId,
		OldPassword:     req.OldPassword,
		NewPassword:     req.NewPassword,
		ConfirmPassword: req.ConfirmPassword,
	}); err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.ChangePasswordResponse{Success: true, Message: "password changed"}, nil
}

func (h *Handler) GetAllUsers(ctx context.Context, req *userpb.GetAllUsersRequest) (*userpb.GetAllUsersResponse, error) {
	page, pageSize := paginationFrom(req)
	out, err := h.profile.GetAllUsers(ctx, domain.GetAllUsersInput{
		AdminID:    req.AdminId,
		Page:       page,
		PageSize:   pageSize,
		RoleFilter: req.RoleFilter,
		BannedOnly: req.BannedOnly,
	})
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return buildGetAllUsersResponse(out), nil
}

func (h *Handler) BanUser(ctx context.Context, req *userpb.BanUserRequest) (*userpb.BanUserResponse, error) {
	user, err := h.profile.BanUser(ctx, req.AdminId, req.UserId, req.Ban, req.Reason)
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	msg := "user unbanned"
	if req.Ban {
		msg = "user banned"
	}
	return &userpb.BanUserResponse{User: toProtoUser(user), Message: msg}, nil
}

func (h *Handler) GetUserByEmail(ctx context.Context, req *userpb.GetUserByEmailRequest) (*userpb.GetUserByEmailResponse, error) {
	user, err := h.profile.GetUserByEmail(ctx, domain.GetUserByEmailInput{
		AdminID: req.AdminId,
		Email:   req.Email,
	})
	if err != nil {
		return nil, apperrors.ToGRPC(err)
	}
	return &userpb.GetUserByEmailResponse{User: toProtoUser(user)}, nil
}

func paginationFrom(req *userpb.GetAllUsersRequest) (page, pageSize int) {
	page, pageSize = 1, 20
	if req.Pagination != nil {
		if req.Pagination.Page > 0 {
			page = int(req.Pagination.Page)
		}
		if req.Pagination.PageSize > 0 {
			pageSize = int(req.Pagination.PageSize)
		}
	}
	return page, pageSize
}

func buildGetAllUsersResponse(out *domain.GetAllUsersOutput) *userpb.GetAllUsersResponse {
	protoUsers := make([]*userpb.User, len(out.Users))
	for i, u := range out.Users {
		protoUsers[i] = toProtoUser(u)
	}
	return &userpb.GetAllUsersResponse{
		Users: protoUsers,
		Pagination: &userpb.PaginationResponse{
			Page:       int32(out.Page),
			PageSize:   int32(out.PageSize),
			TotalPages: int32(out.TotalPages),
			TotalItems: out.TotalItems,
		},
	}
}

func toProtoUser(u *domain.User) *userpb.User {
	return &userpb.User{
		UserId:    u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Phone:     u.Phone,
		Role:      roleToProto(u.Role),
		IsBanned:  u.IsBanned,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func roleToProto(r domain.Role) userpb.Role {
	switch r {
	case domain.RoleAdmin:
		return userpb.Role_ROLE_ADMIN
	case domain.RoleModerator:
		return userpb.Role_ROLE_MODERATOR
	case domain.RoleUser:
		return userpb.Role_ROLE_USER
	case domain.RoleGuest:
		return userpb.Role_ROLE_GUEST
	default:
		return userpb.Role_ROLE_UNSPECIFIED
	}
}
