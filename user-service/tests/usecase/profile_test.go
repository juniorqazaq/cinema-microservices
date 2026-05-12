package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cinema-booking-system/user-service/internal/domain"
	"github.com/cinema-booking-system/user-service/internal/repository/mocks"
	"github.com/cinema-booking-system/user-service/internal/usecase"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

const (
	testAdminUUID  = "650e8400-e29b-41d4-a716-446655440001"
	testTargetUUID = "650e8400-e29b-41d4-a716-446655440002"
)

func newProfileUC(t *testing.T) (*usecase.ProfileUseCase, *mocks.MockUserRepository, *mocks.MockUserCache, *mocks.MockTokenRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mocks.NewMockUserRepository(ctrl)
	cache := mocks.NewMockUserCache(ctrl)
	tokenRepo := mocks.NewMockTokenRepository(ctrl)
	uc := usecase.NewProfileUseCase(userRepo, cache, tokenRepo)
	return uc, userRepo, cache, tokenRepo
}

func TestGetProfile_FromCache(t *testing.T) {
	uc, _, cache, _ := newProfileUC(t)
	u := sampleUser()
	cache.EXPECT().GetUserCache(gomock.Any(), u.ID).Return(u, nil)
	got, err := uc.GetProfile(context.Background(), u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != u.Email {
		t.Fatalf("email: got %q", got.Email)
	}
}

func TestGetProfile_FromDB(t *testing.T) {
	uc, userRepo, cache, _ := newProfileUC(t)
	u := sampleUser()
	cache.EXPECT().GetUserCache(gomock.Any(), u.ID).Return(nil, domain.ErrUserNotFound)
	userRepo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	cache.EXPECT().SetUserCache(gomock.Any(), u).Return(nil)
	got, err := uc.GetProfile(context.Background(), u.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != u.ID {
		t.Fatal("id mismatch")
	}
}

func TestChangePassword_Success(t *testing.T) {
	uc, userRepo, _, _ := newProfileUC(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("oldSecret12"), 12)
	if err != nil {
		t.Fatal(err)
	}
	u := sampleUser()
	u.PasswordHash = string(hash)
	userRepo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	userRepo.EXPECT().UpdatePassword(gomock.Any(), u.ID, gomock.Any()).Return(nil)
	if err := uc.ChangePassword(context.Background(), domain.ChangePasswordInput{
		UserID:          u.ID,
		OldPassword:     "oldSecret12",
		NewPassword:     "newSecret12",
		ConfirmPassword: "newSecret12",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestBanUser_ForbiddenNonAdmin(t *testing.T) {
	uc, userRepo, _, _ := newProfileUC(t)
	admin := sampleUser()
	admin.ID = testAdminUUID
	admin.Role = domain.RoleUser
	userRepo.EXPECT().GetByID(gomock.Any(), testAdminUUID).Return(admin, nil)
	_, err := uc.BanUser(context.Background(), testAdminUUID, testTargetUUID, true, "r")
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("want ErrForbidden, got %v", err)
	}
}

func TestGetAllUsers_DefaultPaging(t *testing.T) {
	uc, userRepo, _, _ := newProfileUC(t)
	admin := sampleUser()
	admin.ID = testAdminUUID
	admin.Role = domain.RoleAdmin
	userRepo.EXPECT().GetByID(gomock.Any(), testAdminUUID).Return(admin, nil)
	userRepo.EXPECT().GetAll(gomock.Any(), 1, 20, "", false).Return([]*domain.User{sampleUser()}, int64(1), nil)
	out, err := uc.GetAllUsers(context.Background(), domain.GetAllUsersInput{
		AdminID:  testAdminUUID,
		Page:     0,
		PageSize: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.TotalItems != 1 {
		t.Fatalf("total: %d", out.TotalItems)
	}
}

func TestGetUserByEmail(t *testing.T) {
	uc, userRepo, _, _ := newProfileUC(t)
	admin := sampleUser()
	admin.ID = testAdminUUID
	admin.Role = domain.RoleAdmin
	u := sampleUser()
	userRepo.EXPECT().GetByID(gomock.Any(), testAdminUUID).Return(admin, nil)
	userRepo.EXPECT().GetByEmail(gomock.Any(), u.Email).Return(u, nil)
	got, err := uc.GetUserByEmail(context.Background(), domain.GetUserByEmailInput{
		AdminID: testAdminUUID,
		Email:   u.Email,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Email != u.Email {
		t.Fatal("email mismatch")
	}
}

func TestUpdateProfile(t *testing.T) {
	uc, userRepo, cache, _ := newProfileUC(t)
	u := sampleUser()
	u.FullName = "Old"
	userRepo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	userRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, updated *domain.User) (*domain.User, error) {
			if updated.FullName != "New" {
				t.Fatalf("full name: %q", updated.FullName)
			}
			return updated, nil
		})
	cache.EXPECT().DeleteUserCache(gomock.Any(), u.ID).Return(nil)
	out, err := uc.UpdateProfile(context.Background(), domain.UpdateProfileInput{
		UserID:   u.ID,
		FullName: "New",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.FullName != "New" {
		t.Fatal("want updated name")
	}
}

func TestDeleteAccount(t *testing.T) {
	uc, userRepo, cache, tokenRepo := newProfileUC(t)
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret1234"), 12)
	u := sampleUser()
	u.PasswordHash = string(hash)
	userRepo.EXPECT().GetByID(gomock.Any(), u.ID).Return(u, nil)
	userRepo.EXPECT().Delete(gomock.Any(), u.ID).Return(nil)
	cache.EXPECT().DeleteUserCache(gomock.Any(), u.ID).Return(nil)
	tokenRepo.EXPECT().DeleteAllUserTokens(gomock.Any(), u.ID).Return(nil)
	if err := uc.DeleteAccount(context.Background(), u.ID, "secret1234"); err != nil {
		t.Fatal(err)
	}
}

func TestBanUser_AsAdmin(t *testing.T) {
	uc, userRepo, cache, tokenRepo := newProfileUC(t)
	admin := sampleUser()
	admin.ID = testAdminUUID
	admin.Role = domain.RoleAdmin
	target := sampleUser()
	target.ID = testTargetUUID
	userRepo.EXPECT().GetByID(gomock.Any(), testAdminUUID).Return(admin, nil)
	userRepo.EXPECT().Ban(gomock.Any(), testTargetUUID, true, "spam").Return(target, nil)
	cache.EXPECT().DeleteUserCache(gomock.Any(), testTargetUUID).Return(nil)
	tokenRepo.EXPECT().DeleteAllUserTokens(gomock.Any(), testTargetUUID).Return(nil)
	out, err := uc.BanUser(context.Background(), testAdminUUID, testTargetUUID, true, "spam")
	if err != nil {
		t.Fatal(err)
	}
	if out.ID != target.ID {
		t.Fatal("user id")
	}
}

func TestBanUser_SelfBan(t *testing.T) {
	uc, userRepo, _, _ := newProfileUC(t)
	admin := sampleUser()
	admin.ID = testAdminUUID
	admin.Role = domain.RoleAdmin
	userRepo.EXPECT().GetByID(gomock.Any(), testAdminUUID).Return(admin, nil)
	_, err := uc.BanUser(context.Background(), testAdminUUID, testAdminUUID, true, "x")
	if !errors.Is(err, domain.ErrSelfBanForbidden) {
		t.Fatalf("want ErrSelfBanForbidden, got %v", err)
	}
}
