package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cinema-booking-system/user-service/internal/config"
	grpcdelivery "github.com/cinema-booking-system/user-service/internal/delivery/grpc"
	"github.com/cinema-booking-system/user-service/internal/repository/postgres"
	redisrepo "github.com/cinema-booking-system/user-service/internal/repository/redis"
	"github.com/cinema-booking-system/user-service/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)
	logger.Info("user service starting")
	cfg := config.MustLoad()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	db, err := connectPostgres(ctx, cfg.DB.DSN(), logger)
	if err != nil {
		logger.Error("postgres connection failed after retries", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("postgres connected")
	redisClient, err := connectRedis(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, logger)
	if err != nil {
		logger.Error("redis connection failed after retries", "err", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	logger.Info("redis connected")
	userRepo := postgres.New(db)
	tokenRepo := redisrepo.New(redisClient)
	authUC := usecase.NewAuthUseCase(userRepo, tokenRepo, cfg.JWT, logger)
	profileUC := usecase.NewProfileUseCase(userRepo, tokenRepo, tokenRepo)
	handler := grpcdelivery.NewHandler(authUC, profileUC)
	srv := grpcdelivery.NewServer(handler, cfg.GRPC.Port, logger, db, redisClient)
	logger.Info("user service ready", "grpc_port", cfg.GRPC.Port)
	if err := srv.Run(ctx); err != nil {
		logger.Error("grpc server stopped with error", "err", err)
		os.Exit(1)
	}
	logger.Info("user service shutdown complete")
}

func connectPostgres(ctx context.Context, dsn string, logger *slog.Logger) (*pgxpool.Pool, error) {
	delays := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for attempt, d := range delays {
		db, err := postgres.NewPool(ctx, dsn)
		if err == nil {
			return db, nil
		}
		lastErr = err
		logger.Warn("postgres not ready", "attempt", attempt+1, "retry_in", d.String(), "err", err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(d):
		}
	}
	return nil, lastErr
}

func connectRedis(ctx context.Context, addr, password string, db int, logger *slog.Logger) (*redis.Client, error) {
	delays := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second}
	var lastErr error
	for attempt, d := range delays {
		client, err := redisrepo.NewClient(addr, password, db)
		if err == nil {
			return client, nil
		}
		lastErr = err
		logger.Warn("redis not ready", "attempt", attempt+1, "retry_in", d.String(), "err", err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(d):
		}
	}
	return nil, lastErr
}
