package grpc

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "booking-service/proto/booking"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	grpcServer *grpc.Server
	handler    *Handler
}

func NewServer(handler *Handler) *Server {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Printf("Recovered from panic: %v", p)
			return status.Errorf(codes.Internal, "panic triggered: %v", p)
		}),
	}

	grpcprom.EnableHandlingTimeHistogram()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpcprom.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
		grpc.StreamInterceptor(grpcprom.StreamServerInterceptor),
	)

	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Println("booking-service metrics on :9102")
		if err := http.ListenAndServe(":9102", mux); err != nil {
			log.Printf("metrics server error: %v", err)
		}
	}()

	return &Server{
		grpcServer: grpcServer,
		handler:    handler,
	}
}

func (s *Server) Run(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	pb.RegisterBookingServiceServer(s.grpcServer, s.handler)
	grpcprom.Register(s.grpcServer)

	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.grpcServer.GracefulStop()
	return nil
}
