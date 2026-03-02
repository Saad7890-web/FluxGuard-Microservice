package config

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 20
	cfg.MinConns = 5
	cfg.MaxConnLifetime = time.Hour

	return pgxpool.NewWithConfig(ctx, cfg)
}