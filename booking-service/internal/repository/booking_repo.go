package repository

import (
	"context"
	"errors"

	"booking-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type bookingRepo struct {
	pool *pgxpool.Pool
}

func NewBookingRepository(pool *pgxpool.Pool) domain.BookingRepository {
	return &bookingRepo{pool: pool}
}

func (r *bookingRepo) Create(ctx context.Context, tx pgx.Tx, booking *domain.Booking) error {
	var seatID string
	err := tx.QueryRow(ctx, "SELECT id FROM seats WHERE id=$1 AND is_available=true FOR UPDATE", booking.SeatID).Scan(&seatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrSeatTaken
		}
		return err
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO bookings (id, user_id, session_id, seat_id, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		booking.ID, booking.UserID, booking.SessionID, booking.SeatID, booking.Status, booking.CreatedAt,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE seats SET is_available=false WHERE id=$1", booking.SeatID)
	return err
}

func (r *bookingRepo) Cancel(ctx context.Context, tx pgx.Tx, id string) error {
	var status string
	var seatID string
	err := tx.QueryRow(ctx, "SELECT status, seat_id FROM bookings WHERE id=$1 FOR UPDATE", id).Scan(&status, &seatID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrBookingNotFound
		}
		return err
	}

	if status == domain.StatusCancelled {
		return domain.ErrAlreadyCancelled
	}

	_, err = tx.Exec(ctx, "UPDATE bookings SET status=$1 WHERE id=$2", domain.StatusCancelled, id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE seats SET is_available=true WHERE id=$1", seatID)
	return err
}

func (r *bookingRepo) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	var b domain.Booking
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, session_id, seat_id, status, created_at FROM bookings WHERE id=$1",
		id,
	).Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (r *bookingRepo) ListByUserID(ctx context.Context, userID string) ([]*domain.Booking, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, session_id, seat_id, status, created_at FROM bookings WHERE user_id=$1 AND status != $2",
		userID, domain.StatusCancelled,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *bookingRepo) GetHistory(ctx context.Context, userID string) ([]*domain.Booking, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, session_id, seat_id, status, created_at FROM bookings WHERE user_id=$1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *bookingRepo) AdminListAll(ctx context.Context, limit, offset int) ([]*domain.Booking, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, session_id, seat_id, status, created_at FROM bookings ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *bookingRepo) GetStats(ctx context.Context) (total, confirmed, cancelled int64, err error) {
	query := `
		SELECT 
			COUNT(*),
			COUNT(*) FILTER (WHERE status = $1),
			COUNT(*) FILTER (WHERE status = $2)
		FROM bookings
	`
	err = r.pool.QueryRow(ctx, query, domain.StatusConfirmed, domain.StatusCancelled).Scan(&total, &confirmed, &cancelled)
	return
}

func (r *bookingRepo) CountAll(ctx context.Context) (int64, error) {
	var n int64
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM bookings").Scan(&n)
	return n, err
}

func (r *bookingRepo) IsSeatAvailable(ctx context.Context, seatID string) (bool, error) {
	var avail bool
	err := r.pool.QueryRow(ctx, "SELECT is_available FROM seats WHERE id = $1", seatID).Scan(&avail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return avail, nil
}

func (r *bookingRepo) UpdateStatus(ctx context.Context, id, status string) error {
	cmd, err := r.pool.Exec(ctx, "UPDATE bookings SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return domain.ErrBookingNotFound
	}
	return nil
}
