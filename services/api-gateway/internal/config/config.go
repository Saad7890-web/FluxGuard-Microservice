package config

import "os"

type Config struct {
	HTTPPort     string
	AuthGRPCAddr string
	JWTSecret    string
}

func Load() *Config {
	return &Config{
		HTTPPort:     getEnv("HTTP_PORT", "8080"),
		AuthGRPCAddr: getEnv("AUTH_GRPC_ADDR", "localhost:50051"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}