package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	connString := os.Getenv("POSTGRES_URL")

	pool, err := pgxpool.New(ctx, connString)

	return pool, err
}
