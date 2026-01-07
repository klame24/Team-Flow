package postgres

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDB(ctx context.Context) (*pgx.Conn, error) {
	connString := os.Getenv("POSTGRES_URL")

	conn, err := pgx.Connect(ctx, connString)

	return conn, err
}
