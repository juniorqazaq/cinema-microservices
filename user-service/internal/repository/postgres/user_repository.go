package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cinema-booking-system/user-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: parse config: %w", err)
	}
	cfg.MaxConns = 10

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: connect: %w", err)
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}
	return pool, nil
}

const createUserSQL = `
INSERT INTO users (email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at`

func (r *UserRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	now := time.Now().UTC()
	row := r.db.QueryRow(ctx, createUserSQL,
		u.Email, u.PasswordHash, u.FullName, u.Phone,
		string(u.Role), u.IsBanned, u.BanReason, now, now,
	)
	created, err := scanUser(row)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, domain.ErrEmailAlreadyExists
		}
		return nil, fmt.Errorf("postgres: create user: %w", err)
	}
	return created, nil
}

const getUserByIDSQL = `
SELECT id, email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at
FROM users WHERE id = $1`

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, getUserByIDSQL, id)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("postgres: get user by id: %w", err)
	}
	return u, nil
}

const getUserByEmailSQL = `
SELECT id, email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at
FROM users WHERE email = $1`

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, getUserByEmailSQL, email)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("postgres: get user by email: %w", err)
	}
	return u, nil
}

const updateUserSQL = `
UPDATE users
SET full_name = $1, phone = $2, updated_at = $3
WHERE id = $4
RETURNING id, email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at`

func (r *UserRepository) Update(ctx context.Context, u *domain.User) (*domain.User, error) {
	row := r.db.QueryRow(ctx, updateUserSQL,
		u.FullName, u.Phone, time.Now().UTC(), u.ID,
	)
	updated, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("postgres: update user: %w", err)
	}
	return updated, nil
}

const updatePasswordSQL = `
UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`

func (r *UserRepository) UpdatePassword(ctx context.Context, userID, hash string) error {
	cmd, err := r.db.Exec(ctx, updatePasswordSQL, hash, time.Now().UTC(), userID)
	if err != nil {
		return fmt.Errorf("postgres: update password: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

const deleteUserSQL = `DELETE FROM users WHERE id = $1`

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	cmd, err := r.db.Exec(ctx, deleteUserSQL, id)
	if err != nil {
		return fmt.Errorf("postgres: delete user: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) GetAll(
	ctx context.Context,
	page, pageSize int,
	roleFilter string,
	bannedOnly bool,
) ([]*domain.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	where := "WHERE 1=1"
	args := []any{}
	argIdx := 1

	if roleFilter != "" {
		where += fmt.Sprintf(" AND role = $%d", argIdx)
		args = append(args, roleFilter)
		argIdx++
	}
	if bannedOnly {
		where += fmt.Sprintf(" AND is_banned = $%d", argIdx)
		args = append(args, true)
		argIdx++
	}

	var total int64
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM users %s", where)
	if err := r.db.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("postgres: count users: %w", err)
	}

	listSQL := fmt.Sprintf(`
		SELECT id, email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at
		FROM users %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, argIdx, argIdx+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, listSQL, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("postgres: list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u, err := scanUserFromRows(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("postgres: scan user: %w", err)
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

const banUserSQL = `
UPDATE users
SET is_banned = $1, ban_reason = $2, updated_at = $3
WHERE id = $4
RETURNING id, email, password_hash, full_name, phone, role, is_banned, ban_reason, created_at, updated_at`

func (r *UserRepository) Ban(ctx context.Context, userID string, ban bool, reason string) (*domain.User, error) {
	row := r.db.QueryRow(ctx, banUserSQL, ban, reason, time.Now().UTC(), userID)
	u, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("postgres: ban user: %w", err)
	}
	return u, nil
}

type scannable interface {
	Scan(dest ...any) error
}

func scanUser(row scannable) (*domain.User, error) {
	var u domain.User
	var role string
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone,
		&role, &u.IsBanned, &u.BanReason, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Role = domain.Role(role)
	return &u, nil
}

func scanUserFromRows(rows pgx.Rows) (*domain.User, error) {
	var u domain.User
	var role string
	err := rows.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.Phone,
		&role, &u.IsBanned, &u.BanReason, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Role = domain.Role(role)
	return &u, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
