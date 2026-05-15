package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cinema-booking-system/api-gateway/internal/config"
	"github.com/cinema-booking-system/api-gateway/internal/infra/grpcclient"
	gatewayhttp "github.com/cinema-booking-system/api-gateway/internal/transport/http"
)

func main() {
	cfg := config.LoadConfig()

	clients, closeClients, err := grpcclient.DialClients(cfg)
	if err != nil {
		log.Fatalf("dial grpc clients: %v", err)
	}
	defer closeClients()

	router := gatewayhttp.NewRouter(cfg, clients)
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Printf("api gateway listening on %s", cfg.HTTPAddr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("shutdown api gateway: %v", err)
		}
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("run api gateway: %v", err)
		}
	}
}
