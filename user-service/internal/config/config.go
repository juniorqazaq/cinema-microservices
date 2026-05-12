package config

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"
)

var (
	cfgOnce sync.Once
	cfg     *Config
)

type Config struct {
	DB    DBConfig
	JWT   JWTConfig
	Redis RedisConfig
	GRPC  GRPCConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	MaxConns int
}

func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%d",
		d.Host, d.Port, d.User, d.Password, d.Name, d.MaxConns,
	)
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
	Audience   string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type GRPCConfig struct {
	Port string
}

func MustLoad() *Config {
	cfgOnce.Do(func() {
		log := slog.Default()
		cfg = &Config{
			DB: DBConfig{
				Host:     mustGetEnv("DB_HOST", log),
				Port:     getEnvOrDefault("DB_PORT", "5432"),
				User:     mustGetEnv("DB_USER", log),
				Password: mustGetEnv("DB_PASSWORD", log),
				Name:     mustGetEnv("DB_NAME", log),
				MaxConns: 10,
			},
			JWT: JWTConfig{
				Secret:     mustGetEnv("JWT_SECRET", log),
				AccessTTL:  parseDuration("JWT_ACCESS_TTL", 15*time.Minute, log),
				RefreshTTL: parseDuration("JWT_REFRESH_TTL", 7*24*time.Hour, log),
				Issuer:     getEnvOrDefault("JWT_ISSUER", "cinema-booking-system"),
				Audience:   getEnvOrDefault("JWT_AUDIENCE", "cinema-users"),
			},
			Redis: RedisConfig{
				Addr:     mustGetEnv("REDIS_ADDR", log),
				Password: getEnvOrDefault("REDIS_PASSWORD", ""),
				DB:       0,
			},
			GRPC: GRPCConfig{
				Port: getEnvOrDefault("GRPC_PORT", "50051"),
			},
		}
	})
	return cfg
}

func mustGetEnv(key string, log *slog.Logger) string {
	v := os.Getenv(key)
	if v == "" {
		log.Error("required environment variable is not set", "key", key)
		os.Exit(1)
	}
	return v
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseDuration(key string, def time.Duration, log *slog.Logger) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Error("invalid duration", "key", key, "value", v, "err", err)
		os.Exit(1)
	}
	return d
}
