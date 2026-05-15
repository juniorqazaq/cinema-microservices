package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr        string
	UserGRPCAddr    string
	MovieGRPCAddr   string
	BookingGRPCAddr string
	RequestTimeout  time.Duration
	AllowedOrigins  []string
	IPRate          float64
	IPBurst         int
	UserRate        float64
	UserBurst       int
}

func LoadConfig() Config {
	return Config{
		HTTPAddr:        envString("HTTP_ADDR", ":8080"),
		UserGRPCAddr:    envString("USER_GRPC_ADDR", "localhost:50051"),
		MovieGRPCAddr:   envString("MOVIE_GRPC_ADDR", "localhost:50052"),
		BookingGRPCAddr: envString("BOOKING_GRPC_ADDR", "localhost:50053"),
		RequestTimeout:  envDuration("REQUEST_TIMEOUT", 5*time.Second),
		AllowedOrigins:  envCSV("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
		IPRate:          envFloat("IP_RATE_PER_SEC", 20),
		IPBurst:         envInt("IP_RATE_BURST", 40),
		UserRate:        envFloat("USER_RATE_PER_SEC", 10),
		UserBurst:       envInt("USER_RATE_BURST", 20),
	}
}

func envString(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func envCSV(key, fallback string) []string {
	raw := envString(key, fallback)
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func envDuration(key string, fallback time.Duration) time.Duration {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func envFloat(key string, fallback float64) float64 {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if n, err := strconv.ParseFloat(v, 64); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}
