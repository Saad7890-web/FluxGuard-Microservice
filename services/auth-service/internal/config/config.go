package config

import (
	"log"
	"os"
)

type Config struct {
	DBUrl         string
	RedisURL      string
	JWTSecret     string
	AccessTTL     int64
	RefreshTTL    int64
	GRPCPort      string
	MigrationsDir string
}

func Load() *Config {
	cfg := &Config{
		DBUrl:         os.Getenv("DB_URL"),
		RedisURL:      os.Getenv("REDIS_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		GRPCPort:      getEnv("GRPC_PORT", "50051"),
		MigrationsDir: "./migrations",
		AccessTTL:     900,     // 15 min
		RefreshTTL:    604800,  // 7 days
	}

	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatal("DB_URL and JWT_SECRET are required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}