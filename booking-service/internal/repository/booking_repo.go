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
	var taken bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE session_id = $1 AND seat_id = $2 AND status IN ($3, $4)
		)`,
		booking.SessionID, booking.SeatID, domain.StatusPending, domain.StatusConfirmed,
	).Scan(&taken)
	if err != nil {
		return err
	}
	if taken {
		return domain.ErrSeatTaken
	}

	category := booking.TicketCategory
	if category == "" {
		category = "adult"
	}
	_, err = tx.Exec(ctx,
		"INSERT INTO bookings (id, user_id, session_id, seat_id, status, ticket_category, amount_paid, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		booking.ID, booking.UserID, booking.SessionID, booking.SeatID, booking.Status, category, booking.AmountPaid, booking.CreatedAt,
	)
	return err
}

func (r *bookingRepo) Cancel(ctx context.Context, tx pgx.Tx, id string) error {
	var status string
	err := tx.QueryRow(ctx, "SELECT status FROM bookings WHERE id=$1 FOR UPDATE", id).Scan(&status)
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
	return err
}

func (r *bookingRepo) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	var b domain.Booking
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, session_id, seat_id, status, ticket_category, amount_paid, created_at FROM bookings WHERE id=$1",
		id,
	).Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.TicketCategory, &b.AmountPaid, &b.CreatedAt)

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
		"SELECT id, user_id, session_id, seat_id, status, ticket_category, amount_paid, created_at FROM bookings WHERE user_id=$1 AND status != $2",
		userID, domain.StatusCancelled,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.TicketCategory, &b.AmountPaid, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *bookingRepo) GetHistory(ctx context.Context, userID string) ([]*domain.Booking, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, session_id, seat_id, status, ticket_category, amount_paid, created_at FROM bookings WHERE user_id=$1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.TicketCategory, &b.AmountPaid, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *bookingRepo) AdminListAll(ctx context.Context, limit, offset int) ([]*domain.Booking, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, session_id, seat_id, status, ticket_category, amount_paid, created_at FROM bookings ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var b domain.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.SessionID, &b.SeatID, &b.Status, &b.TicketCategory, &b.AmountPaid, &b.CreatedAt); err != nil {
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
	return true, nil
}

func (r *bookingRepo) IsSeatAvailableForSession(ctx context.Context, sessionID, seatID string) (bool, error) {
	var taken bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM bookings
			WHERE session_id = $1 AND seat_id = $2 AND status IN ($3, $4)
		)`,
		sessionID, seatID, domain.StatusPending, domain.StatusConfirmed,
	).Scan(&taken)
	if err != nil {
		return false, err
	}
	return !taken, nil
}

func (r *bookingRepo) ListTakenSeatIDs(ctx context.Context, sessionID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT seat_id FROM bookings
		WHERE session_id = $1 AND status IN ($2, $3)`,
		sessionID, domain.StatusPending, domain.StatusConfirmed,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
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
