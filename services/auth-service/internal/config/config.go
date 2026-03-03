package config

import (
	"log"
	"os"
)

type Config struct {
	DBUrl        string
	JWTSecret    string
	AccessTTL    int64
	RefreshTTL   int64
	GRPCPort     string
	MigrationsDir string
}

func Load() *Config {
	cfg := &Config{
		DBUrl:        os.Getenv("DB_URL"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		GRPCPort:     getEnv("GRPC_PORT", "50051"),
		MigrationsDir: "./migrations",
	}

	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	// 15 min access, 7 days refresh
	cfg.AccessTTL = 900
	cfg.RefreshTTL = 604800

	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}