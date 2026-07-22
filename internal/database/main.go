package database

import (
	"context"
	"fmt"

	"github.com/RedditUclaista/community-service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}
