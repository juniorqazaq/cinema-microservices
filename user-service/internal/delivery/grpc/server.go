package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	userpb "github.com/cinema-booking-system/user-service/gen/go/user"
	"github.com/cinema-booking-system/user-service/internal/health"
	"github.com/jackc/pgx/v5/pgxpool"
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
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryRecoveryInterceptor(logger),
			UnaryTimingInterceptor(logger),
			UnaryRequestLoggingInterceptor(logger),
		),
	)
	userpb.RegisterUserServiceServer(srv, handler)
	if db != nil && rdb != nil {
		health.Register(srv, health.NewServer(db, rdb, logger))
	}
	reflection.Register(srv)
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
