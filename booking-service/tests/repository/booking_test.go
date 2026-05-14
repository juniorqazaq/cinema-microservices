package repository_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(ctx context.Context, t *testing.T) (*pgxpool.Pool, func()) {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	require.NoError(t, err)

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)

	// Create schema
	schema := `
	CREATE TABLE IF NOT EXISTS seats (
		id VARCHAR(50) PRIMARY KEY,
		is_available BOOLEAN NOT NULL DEFAULT true
	);
	CREATE TABLE IF NOT EXISTS bookings (
		id VARCHAR(50) PRIMARY KEY,
		user_id VARCHAR(50),
		session_id VARCHAR(50),
		seat_id VARCHAR(50),
		status VARCHAR(20),
		created_at TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS payments (
		id VARCHAR(50) PRIMARY KEY,
		booking_id VARCHAR(50),
		amount NUMERIC,
		status VARCHAR(20),
		paid_at TIMESTAMP
	);
	`
	_, err = pool.Exec(ctx, schema)
	require.NoError(t, err)

	teardown := func() {
		pool.Close()
		if err := postgresContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}

	return pool, teardown
}

func TestConcurrentBooking(t *testing.T) {
	ctx := context.Background()
	pool, teardown := setupTestDB(ctx, t)
	defer teardown()

	repo := repository.NewBookingRepository(pool)

	// Insert test seat
	seatID := "seat-10"
	_, err := pool.Exec(ctx, "INSERT INTO seats (id, is_available) VALUES ($1, true)", seatID)
	require.NoError(t, err)

	var wg sync.WaitGroup
	results := make(chan error, 2)

	// 2 goroutines trying to book the same seat
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			
			tx, err := pool.Begin(ctx)
			if err != nil {
				results <- err
				return
			}
			defer tx.Rollback(ctx)

			booking := &domain.Booking{
				ID:        uuid.NewString(),
				UserID:    fmt.Sprintf("user-%d", idx),
				SessionID: "session-1",
				SeatID:    seatID,
				Status:    domain.StatusPending,
				CreatedAt: time.Now(),
			}

			err = repo.Create(ctx, tx, booking)
			if err == nil {
				err = tx.Commit(ctx)
			}
			results <- err
		}(i)
	}

	wg.Wait()
	close(results)

	var successCount, errTakenCount int
	for err := range results {
		if err == nil {
			successCount++
		} else if err == domain.ErrSeatTaken {
			errTakenCount++
		} else {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	assert.Equal(t, 1, successCount, "Exactly 1 booking should succeed")
	assert.Equal(t, 1, errTakenCount, "Exactly 1 booking should fail with ErrSeatTaken")
}

func TestTransactionRollback(t *testing.T) {
	ctx := context.Background()
	pool, teardown := setupTestDB(ctx, t)
	defer teardown()

	repo := repository.NewBookingRepository(pool)

	seatID := "seat-20"
	_, err := pool.Exec(ctx, "INSERT INTO seats (id, is_available) VALUES ($1, true)", seatID)
	require.NoError(t, err)

	tx, err := pool.Begin(ctx)
	require.NoError(t, err)

	booking := &domain.Booking{
		ID:        uuid.NewString(),
		UserID:    "user-1",
		SessionID: "session-1",
		SeatID:    seatID,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	// repo.Create locks the seat and inserts the booking
	err = repo.Create(ctx, tx, booking)
	require.NoError(t, err)

	// Simulate an error occurring AFTER repo.Create (e.g., NATS failed)
	// We call Rollback explicitly
	err = tx.Rollback(ctx)
	require.NoError(t, err)

	// Verify that the seat is STILL available
	var isAvailable bool
	err = pool.QueryRow(ctx, "SELECT is_available FROM seats WHERE id = $1", seatID).Scan(&isAvailable)
	require.NoError(t, err)

	assert.True(t, isAvailable, "Seat should remain available after transaction rollback")
	
	// Verify booking does not exist
	var count int
	err = pool.QueryRow(ctx, "SELECT count(*) FROM bookings WHERE id = $1", booking.ID).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Booking should not exist after transaction rollback")
}
