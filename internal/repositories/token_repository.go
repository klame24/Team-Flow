package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository interface {
	SaveRefreshTokens(
		ctx context.Context,
		userID int32,
		tokenHash string,
		expiresAt time.Time,
	) (int, error)
}

type tokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return &tokenRepository{
		db: db,
	}
}

func (tokenRepo *tokenRepository) SaveRefreshTokens(
	ctx context.Context,
	userID int32,
	tokenHash string,
	expiresAt time.Time,
) (int, error) {
	var TokenRefreshID int

	sqlQuery := `
		INSERT INTO refresh_tokens(user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	err := tokenRepo.db.QueryRow(
		ctx, sqlQuery,
		userID,
		tokenHash,
		expiresAt,
		time.Now(),
	).Scan(&TokenRefreshID)

	return TokenRefreshID, err
}
