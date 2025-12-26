package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() (pool *pgxpool.Pool, err error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(dbURL)

	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create database pool: %w", err)
	}
	return pool, err
}
