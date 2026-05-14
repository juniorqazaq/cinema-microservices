package config

import "os"

type Config struct {
	GRPCPort string
	DBURL    string
	NATSURL  string
}

func New() *Config {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053"
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/booking?sslmode=disable"
	}
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	return &Config{
		GRPCPort: port,
		DBURL:    dbURL,
		NATSURL:  natsURL,
	}
}
