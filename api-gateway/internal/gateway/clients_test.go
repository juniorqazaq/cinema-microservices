package gateway

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestDialClientsReturnsQuicklyWhenTargetsAreUnavailable(t *testing.T) {
	t.Parallel()

	cfg := Config{
		UserGRPCAddr:    "127.0.0.1:1",
		MovieGRPCAddr:   "127.0.0.1:2",
		BookingGRPCAddr: "127.0.0.1:3",
	}

	start := time.Now()
	clients, closeClients, err := DialClients(cfg)
	if err != nil {
		t.Fatalf("DialClients() error = %v", err)
	}
	defer closeClients()
	if clients.User == nil || clients.Movie == nil || clients.Booking == nil {
		t.Fatal("DialClients() returned nil service clients")
	}

	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Fatalf("DialClients() took too long: %v", elapsed)
	}
}

func TestRouterReturnsServiceUnavailableWhenMovieServiceIsDown(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	clients, closeClients, err := DialClients(Config{
		UserGRPCAddr:    "127.0.0.1:1",
		MovieGRPCAddr:   "127.0.0.1:2",
		BookingGRPCAddr: "127.0.0.1:3",
	})
	if err != nil {
		t.Fatalf("DialClients() error = %v", err)
	}
	defer closeClients()

	cfg := testConfig()
	cfg.RequestTimeout = 250 * time.Millisecond

	res := perform(NewRouter(cfg, clients), http.MethodGet, "/movies", "", nil)
	if res.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body = %s", res.Code, res.Body.String())
	}
}
