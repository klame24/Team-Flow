package repositories

import (
	"context"
	"team-flow/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int, error)
	// UPDATE
	// DELETE
	// GET_BY_ID
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepo *userRepository) Create(ctx context.Context, user *models.User) (int, error) {
	var userID int

	user.CreatedAt = time.Now()

	sqlQuery := `
		INSERT INTO users (name, surname, nickname, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`

	err := userRepo.db.QueryRow(
		ctx,
		sqlQuery,
		user.Name,
		user.Surname,
		user.Nickname,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
	).Scan(&userID)

	return userID, err
}
