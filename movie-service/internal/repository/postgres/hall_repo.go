package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HallRepository struct {
	db *pgxpool.Pool
}

func NewHallRepository(db *pgxpool.Pool) *HallRepository {
	return &HallRepository{db: db}
}

const createHallSQL = `INSERT INTO halls (name, capacity) VALUES ($1, $2) RETURNING id, name, capacity`

func (r *HallRepository) Create(ctx context.Context, h *domain.Hall) (*domain.Hall, error) {
	row := r.db.QueryRow(ctx, createHallSQL, h.Name, h.Capacity)
	var out domain.Hall
	if err := row.Scan(&out.ID, &out.Name, &out.Capacity); err != nil {
		return nil, fmt.Errorf("postgres: create hall: %w", err)
	}
	return &out, nil
}

const getHallByIDSQL = `SELECT id, name, capacity FROM halls WHERE id = $1`

func (r *HallRepository) GetByID(ctx context.Context, id string) (*domain.Hall, error) {
	row := r.db.QueryRow(ctx, getHallByIDSQL, id)
	var h domain.Hall
	if err := row.Scan(&h.ID, &h.Name, &h.Capacity); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrHallNotFound
		}
		return nil, fmt.Errorf("postgres: get hall: %w", err)
	}
	return &h, nil
}

const getSeatsByHallSQL = `
SELECT id, hall_id, row, number, is_available FROM seats WHERE hall_id = $1 ORDER BY row, number`

func (r *HallRepository) GetSeats(ctx context.Context, hallID string) ([]*domain.Seat, error) {
	rows, err := r.db.Query(ctx, getSeatsByHallSQL, hallID)
	if err != nil {
		return nil, fmt.Errorf("postgres: get seats: %w", err)
	}
	defer rows.Close()
	var list []*domain.Seat
	for rows.Next() {
		var s domain.Seat
		if err := rows.Scan(&s.ID, &s.HallID, &s.Row, &s.Number, &s.IsAvailable); err != nil {
			return nil, err
		}
		list = append(list, &s)
	}
	return list, rows.Err()
}

const insertSeatSQL = `INSERT INTO seats (hall_id, row, number, is_available) VALUES ($1, $2, $3, true)`

// InsertSeatsForHall creates a simple row layout: R1..Rk with numbered seats up to capacity.
func (r *HallRepository) InsertSeatsForHall(ctx context.Context, hallID string, capacity int) error {
	perRow := 10
	if capacity < perRow {
		perRow = capacity
	}
	if perRow == 0 {
		return nil
	}
	rows := (capacity + perRow - 1) / perRow
	n := 0
	for ri := 0; ri < rows && n < capacity; ri++ {
		rowLabel := fmt.Sprintf("R%d", ri+1)
		for seatNum := 1; seatNum <= perRow && n < capacity; seatNum++ {
			if _, err := r.db.Exec(ctx, insertSeatSQL, hallID, rowLabel, seatNum); err != nil {
				return fmt.Errorf("postgres: insert seat: %w", err)
			}
			n++
		}
	}
	return nil
}
