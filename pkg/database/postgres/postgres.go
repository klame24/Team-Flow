package postgres

import (
	"context"
	"fmt"
	"team-flow/core/logger"
	"team-flow/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool *pgxpool.Pool
}

func Connect(ctx context.Context, cfg config.DatabaseConfig) (*Database, error) {
	start := time.Now()

	pool, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	elapsed := time.Since(start)
	logger.LogInfo(fmt.Sprintf("Connected to database in %v", elapsed))

	return &Database{Pool: pool}, nil
}

func (db *Database) GetPool() *pgxpool.Pool {
	return db.Pool
}

func (db *Database) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		logger.LogInfo("Database connection closed")
	}
}
