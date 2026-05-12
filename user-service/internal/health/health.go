package health

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

type Server struct {
	db     *pgxpool.Pool
	redis  *redis.Client
	logger *slog.Logger
}

func NewServer(db *pgxpool.Pool, redisClient *redis.Client, logger *slog.Logger) *Server {
	return &Server{db: db, redis: redisClient, logger: logger}
}

func Register(grpcSrv *grpc.Server, h *Server) {
	healthpb.RegisterHealthServer(grpcSrv, h)
}

func (s *Server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	svc := req.GetService()
	if svc != "" && svc != "user.UserService" {
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVICE_UNKNOWN}, nil
	}
	if err := s.db.Ping(ctx); err != nil {
		s.logger.Warn("healthcheck db ping failed", "err", err)
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}, nil
	}
	if err := s.redis.Ping(ctx).Err(); err != nil {
		s.logger.Warn("healthcheck redis ping failed", "err", err)
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}, nil
	}
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *Server) Watch(_ *healthpb.HealthCheckRequest, _ healthpb.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "watch is not supported")
}
