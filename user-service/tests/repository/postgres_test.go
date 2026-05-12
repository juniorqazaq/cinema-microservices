//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cinema-booking-system/user-service/internal/domain"
	pgRepo "github.com/cinema-booking-system/user-service/internal/repository/postgres"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgres(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "cinema",
			"POSTGRES_PASSWORD": "cinema_secret",
			"POSTGRES_DB":       "cinema_test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	t.Cleanup(func() { _ = container.Terminate(ctx) })

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	return fmt.Sprintf(
		"host=%s port=%s user=cinema password=cinema_secret dbname=cinema_test sslmode=disable pool_max_conns=5",
		host, port.Port(),
	)
}

func applyMigrations(t *testing.T, dsn string) {
	t.Helper()
	ctx := context.Background()
	pool, err := pgRepo.NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("connect for migration: %v", err)
	}
	defer pool.Close()

	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

		CREATE TABLE IF NOT EXISTS users (
			id            UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
			email         TEXT        UNIQUE NOT NULL,
			password_hash TEXT        NOT NULL,
			full_name     TEXT        NOT NULL DEFAULT '',
			phone         TEXT        NOT NULL DEFAULT '',
			role          TEXT        NOT NULL DEFAULT 'USER',
			is_banned     BOOLEAN     NOT NULL DEFAULT false,
			ban_reason    TEXT        NOT NULL DEFAULT '',
			created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);
	`)
	if err != nil {
		t.Fatalf("apply migration: %v", err)
	}
}

func TestCreate_GetByID(t *testing.T) {
	dsn := startPostgres(t)
	applyMigrations(t, dsn)

	ctx := context.Background()
	pool, err := pgRepo.NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	repo := pgRepo.New(pool)

	user := &domain.User{
		Email:        "integration@test.com",
		PasswordHash: "$2a$12$hashedpassword",
		FullName:     "Integration User",
		Phone:        "+9998887766",
		Role:         domain.RoleUser,
	}

	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID after Create")
	}
	if created.Email != user.Email {
		t.Errorf("email mismatch: got %s want %s", created.Email, user.Email)
	}

	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s want %s", fetched.ID, created.ID)
	}
	if fetched.FullName != user.FullName {
		t.Errorf("FullName mismatch: got %s want %s", fetched.FullName, user.FullName)
	}
}

func TestBanUser(t *testing.T) {
	dsn := startPostgres(t)
	applyMigrations(t, dsn)

	ctx := context.Background()
	pool, _ := pgRepo.NewPool(ctx, dsn)
	defer pool.Close()

	repo := pgRepo.New(pool)

	user, _ := repo.Create(ctx, &domain.User{
		Email:        "banme@test.com",
		PasswordHash: "hash",
		Role:         domain.RoleUser,
	})

	banned, err := repo.Ban(ctx, user.ID, true, "violation of terms")
	if err != nil {
		t.Fatalf("Ban: %v", err)
	}
	if !banned.IsBanned {
		t.Error("expected is_banned=true")
	}
	if banned.BanReason != "violation of terms" {
		t.Errorf("ban reason mismatch: got %q", banned.BanReason)
	}
	unbanned, err := repo.Ban(ctx, user.ID, false, "")
	if err != nil {
		t.Fatalf("Unban: %v", err)
	}
	if unbanned.IsBanned {
		t.Error("expected is_banned=false after unban")
	}
}

func TestGetByEmail_NotFound(t *testing.T) {
	dsn := startPostgres(t)
	applyMigrations(t, dsn)

	ctx := context.Background()
	pool, _ := pgRepo.NewPool(ctx, dsn)
	defer pool.Close()

	repo := pgRepo.New(pool)

	_, err := repo.GetByEmail(ctx, "ghost@nowhere.com")
	if err == nil {
		t.Fatal("expected error for missing email, got nil")
	}
	if err != domain.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestCreate_DuplicateEmail(t *testing.T) {
	dsn := startPostgres(t)
	applyMigrations(t, dsn)

	ctx := context.Background()
	pool, _ := pgRepo.NewPool(ctx, dsn)
	defer pool.Close()

	repo := pgRepo.New(pool)

	u := &domain.User{Email: "dup@test.com", PasswordHash: "h", Role: domain.RoleUser}
	if _, err := repo.Create(ctx, u); err != nil {
		t.Fatalf("first Create: %v", err)
	}
	_, err := repo.Create(ctx, u)
	if err != domain.ErrEmailAlreadyExists {
		t.Errorf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestGetAll_Pagination(t *testing.T) {
	dsn := startPostgres(t)
	applyMigrations(t, dsn)

	ctx := context.Background()
	pool, _ := pgRepo.NewPool(ctx, dsn)
	defer pool.Close()

	repo := pgRepo.New(pool)
	for i := 0; i < 5; i++ {
		_, _ = repo.Create(ctx, &domain.User{
			Email:        fmt.Sprintf("user%d@test.com", i),
			PasswordHash: "hash",
			Role:         domain.RoleUser,
		})
	}

	users, total, err := repo.GetAll(ctx, 1, 2, "", false)
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if total < 5 {
		t.Errorf("expected at least 5 total, got %d", total)
	}
	if len(users) != 2 {
		t.Errorf("expected 2 users on page 1 (pageSize=2), got %d", len(users))
	}
}
