package config

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
)

var (
	cfgOnce sync.Once
	cfg     *Config
)

type Config struct {
	DB    DBConfig
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
			Redis: RedisConfig{
				Addr:     mustGetEnv("REDIS_ADDR", log),
				Password: getEnvOrDefault("REDIS_PASSWORD", ""),
				DB:       0,
			},
			GRPC: GRPCConfig{
				Port: getEnvOrDefault("GRPC_PORT", "50052"),
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
