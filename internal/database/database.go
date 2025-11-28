package database

import (
	"context"
	"fmt"
	"os"
	"pastebin/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(dbConfig *config.DBConfig) (pool *pgxpool.Pool, err error) {
	dbURL := os.Getenv("DATABASE_URL")
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
