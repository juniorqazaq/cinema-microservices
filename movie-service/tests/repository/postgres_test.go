//go:build integration

package repository_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cinema-booking-system/movie-service/internal/domain"
	pgRepo "github.com/cinema-booking-system/movie-service/internal/repository/postgres"
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
			"POSTGRES_DB":       "movie_test",
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
		t.Fatalf("start postgres: %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })
	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")
	return fmt.Sprintf(
		"host=%s port=%s user=cinema password=cinema_secret dbname=movie_test sslmode=disable pool_max_conns=5",
		host, port.Port(),
	)
}

func applyMovieMigrations(t *testing.T, dsn string) {
	t.Helper()
	ctx := context.Background()
	pool, err := pgRepo.NewPool(ctx, dsn, 5)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()
	root := findRepoRoot(t)
	for _, name := range []string{
		"000001_create_movies.up.sql",
		"000002_create_halls.up.sql",
		"000003_create_seats.up.sql",
		"000004_create_sessions.up.sql",
	} {
		b, err := os.ReadFile(filepath.Join(root, "migrations", name))
		if err != nil {
			t.Fatalf("read migration %s: %v", name, err)
		}
		if _, err := pool.Exec(ctx, string(b)); err != nil {
			t.Fatalf("exec migration %s: %v", name, err)
		}
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found")
		}
		dir = parent
	}
}

func TestMovieCRUD(t *testing.T) {
	dsn := startPostgres(t)
	applyMovieMigrations(t, dsn)
	ctx := context.Background()
	pool, err := pgRepo.NewPool(ctx, dsn, 5)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repo := pgRepo.NewMovieRepository(pool)
	created, err := repo.Create(ctx, &domain.Movie{Title: "T", Description: "d", Genre: "action", Duration: 100, Rating: 8})
	if err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByID(ctx, created.ID)
	if err != nil || got.Title != "T" {
		t.Fatalf("get: %v %+v", err, got)
	}
	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.GetByID(ctx, created.ID); err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestSearchByGenre(t *testing.T) {
	dsn := startPostgres(t)
	applyMovieMigrations(t, dsn)
	ctx := context.Background()
	pool, err := pgRepo.NewPool(ctx, dsn, 5)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	repo := pgRepo.NewMovieRepository(pool)
	_, _ = repo.Create(ctx, &domain.Movie{Title: "A", Description: "", Genre: "comedy", Duration: 90, Rating: 7})
	_, _ = repo.Create(ctx, &domain.Movie{Title: "B", Description: "", Genre: "comedy", Duration: 95, Rating: 6})
	list, total, err := repo.Search(ctx, "", "comedy", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if total < 2 || len(list) < 2 {
		t.Fatalf("search: total=%d len=%d", total, len(list))
	}
}

func TestSessionListByDate(t *testing.T) {
	dsn := startPostgres(t)
	applyMovieMigrations(t, dsn)
	ctx := context.Background()
	pool, err := pgRepo.NewPool(ctx, dsn, 5)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	movies := pgRepo.NewMovieRepository(pool)
	halls := pgRepo.NewHallRepository(pool)
	sessions := pgRepo.NewSessionRepository(pool)
	m, _ := movies.Create(ctx, &domain.Movie{Title: "M", Description: "", Genre: "g", Duration: 100, Rating: 5})
	h, _ := halls.Create(ctx, &domain.Hall{Name: "H1", Capacity: 5})
	_ = halls.InsertSeatsForHall(ctx, h.ID, 5)
	day := time.Date(2026, 5, 10, 18, 0, 0, 0, time.UTC)
	_, err = sessions.Create(ctx, &domain.Session{MovieID: m.ID, HallID: h.ID, StartTime: day, Price: 12.5})
	if err != nil {
		t.Fatal(err)
	}
	list, err := sessions.GetByDate(ctx, day)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("sessions=%d", len(list))
	}
	if list[0].MovieID != m.ID {
		t.Fatalf("movie id mismatch: %s vs %s", list[0].MovieID, m.ID)
	}
}
