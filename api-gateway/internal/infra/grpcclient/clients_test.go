package grpcclient

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cinema-booking-system/api-gateway/internal/config"
	gatewayhttp "github.com/cinema-booking-system/api-gateway/internal/transport/http"
	"github.com/gin-gonic/gin"
)

func TestDialClientsReturnsQuicklyWhenTargetsAreUnavailable(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
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

	clients, closeClients, err := DialClients(config.Config{
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

	res := perform(gatewayhttp.NewRouter(cfg, clients), http.MethodGet, "/movies", "", nil)
	if res.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d body = %s", res.Code, res.Body.String())
	}
}

func testConfig() config.Config {
	return config.Config{
		HTTPAddr:       ":0",
		RequestTimeout: time.Second,
		AllowedOrigins: []string{"http://localhost:5173"},
		IPRate:         1000,
		IPBurst:        1000,
		UserRate:       1000,
		UserBurst:      1000,
	}
}

func perform(router http.Handler, method, path, auth string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
