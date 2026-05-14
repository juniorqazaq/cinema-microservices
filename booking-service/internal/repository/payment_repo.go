package repository

import (
	"context"
	"errors"

	"booking-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type paymentRepo struct {
	pool *pgxpool.Pool
}

func NewPaymentRepository(pool *pgxpool.Pool) domain.PaymentRepository {
	return &paymentRepo{pool: pool}
}

func (r *paymentRepo) Create(ctx context.Context, p *domain.Payment) error {
	query := `
		INSERT INTO payments (id, booking_id, amount, status, paid_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query, p.ID, p.BookingID, p.Amount, p.Status, p.PaidAt)
	return err
}

func (r *paymentRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	var p domain.Payment
	query := `
		SELECT id, booking_id, amount, status, paid_at
		FROM payments
		WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.BookingID, &p.Amount, &p.Status, &p.PaidAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) GetByBookingID(ctx context.Context, bookingID string) (*domain.Payment, error) {
	var p domain.Payment
	query := `
		SELECT id, booking_id, amount, status, paid_at
		FROM payments
		WHERE booking_id = $1
	`
	err := r.pool.QueryRow(ctx, query, bookingID).Scan(&p.ID, &p.BookingID, &p.Amount, &p.Status, &p.PaidAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) List(ctx context.Context, limit, offset int) ([]*domain.Payment, error) {
	query := `
		SELECT id, booking_id, amount, status, paid_at
		FROM payments
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*domain.Payment
	for rows.Next() {
		var p domain.Payment
		if err := rows.Scan(&p.ID, &p.BookingID, &p.Amount, &p.Status, &p.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, &p)
	}
	return payments, nil
}

func (r *paymentRepo) CountAll(ctx context.Context) (int64, error) {
	var n int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM payments").Scan(&n)
	return n, err
}

func (r *paymentRepo) UpdateStatus(ctx context.Context, id, status string) error {
	var query string
	if status == domain.PaymentStatusPaid {
		query = "UPDATE payments SET status = $1, paid_at = NOW() WHERE id = $2"
	} else {
		query = "UPDATE payments SET status = $1 WHERE id = $2"
	}

	cmd, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return domain.ErrPaymentNotFound
	}

	return nil
}
