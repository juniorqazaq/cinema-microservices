package main

import (
	"context"
	"log"
	"time"

	"booking-service/internal/config"
	deliveryGrpc "booking-service/internal/delivery/grpc"
	"booking-service/internal/publisher"
	"booking-service/internal/repository"
	"booking-service/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

func main() {
	cfg := config.New()

	ctx := context.Background()

	var pool *pgxpool.Pool
	var err error
	dbConfig, err := pgxpool.ParseConfig(cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}
	dbConfig.MaxConns = 10

	dbRetries := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	for i, d := range dbRetries {
		pool, err = pgxpool.NewWithConfig(ctx, dbConfig)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				break
			}
		}
		log.Printf("DB connection failed, retrying in %v (attempt %d/3)...", d, i+1)
		time.Sleep(d)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}
	defer pool.Close()

	var nc *nats.Conn
	natsRetries := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}
	for i, d := range natsRetries {
		nc, err = nats.Connect(cfg.NATSURL)
		if err == nil {
			break
		}
		log.Printf("NATS connection failed, retrying in %v (attempt %d/3)...", d, i+1)
		time.Sleep(d)
	}
	if err != nil {
		log.Fatalf("Failed to connect to NATS after retries: %v", err)
	}
	defer nc.Close()

	bookingRepo := repository.NewBookingRepository(pool)
	paymentRepo := repository.NewPaymentRepository(pool)

	pub := publisher.NewNATSPublisher(nc)

	bookingUC := usecase.NewBookingUseCase(bookingRepo, pool, pub)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, bookingRepo, pub)

	handler := deliveryGrpc.NewHandler(bookingUC, paymentUC, bookingRepo, paymentRepo)

	grpcServer := deliveryGrpc.NewServer(handler)

	if err := grpcServer.Run(cfg.GRPCPort); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
