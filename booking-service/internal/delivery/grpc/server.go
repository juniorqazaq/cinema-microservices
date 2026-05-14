package grpc

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "booking-service/proto/booking"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
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
	// Logger setup (simple zap production logger)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Recovery handler
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Printf("Recovered from panic: %v", p)
			return status.Errorf(codes.Internal, "panic triggered: %v", p)
		}),
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(logger),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
	)

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

	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gRPC server...")
	s.grpcServer.GracefulStop()
	log.Println("gRPC server stopped.")

	return nil
}
