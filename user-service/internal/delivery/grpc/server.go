package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/cinema-booking-system/user-service/internal/health"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcSrv *grpc.Server
	port    string
	logger  *slog.Logger
}

func NewServer(
	handler *Handler,
	port string,
	logger *slog.Logger,
	db *pgxpool.Pool,
	rdb *redis.Client,
) *Server {
	grpcprom.EnableHandlingTimeHistogram()

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcprom.UnaryServerInterceptor,
			UnaryRecoveryInterceptor(logger),
			UnaryTimingInterceptor(logger),
			UnaryRequestLoggingInterceptor(logger),
		),
		grpc.ChainStreamInterceptor(
			grpcprom.StreamServerInterceptor,
		),
	)
	userpb.RegisterUserServiceServer(srv, handler)
	grpcprom.Register(srv)

	if db != nil && rdb != nil {
		health.Register(srv, health.NewServer(db, rdb, logger))
	}
	reflection.Register(srv)

	// Prometheus metrics HTTP server :9101
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		logger.Info("metrics server listening", "addr", ":9101")
		if err := http.ListenAndServe(":9101", mux); err != nil {
			logger.Error("metrics server failed", "err", err)
		}
	}()

	return &Server{grpcSrv: srv, port: port, logger: logger}
}

func (s *Server) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", s.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("grpc listen %s: %w", addr, err)
	}
	s.logger.Info("grpc listening", "addr", addr)
	errCh := make(chan error, 1)
	go func() {
		if err := s.grpcSrv.Serve(lis); err != nil {
			errCh <- err
		}
	}()
	select {
	case <-ctx.Done():
		s.logger.Info("grpc graceful shutdown")
		s.grpcSrv.GracefulStop()
		return nil
	case err := <-errCh:
		return fmt.Errorf("grpc serve: %w", err)
	}
}
