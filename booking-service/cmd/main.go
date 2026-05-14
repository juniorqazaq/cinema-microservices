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
	// 1. config.New()
	cfg := config.New()

	ctx := context.Background()

	// 2. pgxpool.New() с max_conns=10, retry 3 раза (1с, 2с, 4с)
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

	// 3. nats.Connect() с retry 3 раза
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

	// 4. repositories
	bookingRepo := repository.NewBookingRepository(pool)
	paymentRepo := repository.NewPaymentRepository(pool)

	// Publisher
	pub := publisher.NewNATSPublisher(nc)

	// 5. usecases (передать pool для транзакций)
	bookingUC := usecase.NewBookingUseCase(bookingRepo, pool, pub)
	paymentUC := usecase.NewPaymentUseCase(paymentRepo, bookingRepo, pub)

	// 6. grpc handler
	handler := deliveryGrpc.NewHandler(bookingUC, paymentUC, bookingRepo, paymentRepo)

	// 7. grpc server -> serve
	grpcServer := deliveryGrpc.NewServer(handler)
	
	if err := grpcServer.Run(cfg.GRPCPort); err != nil {
		log.Fatalf("Failed to run gRPC server: %v", err)
	}
}
