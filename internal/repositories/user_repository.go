package repositories

import (
	"context"
	"team-flow/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int32, error)
	// UPDATE
	// DELETE
	// GET_BY_ID
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (userRepo *userRepository) Create(ctx context.Context, user *models.User) (int32, error) {
	var userID int32

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

func (userRepo *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	sqlQuery := `
		SELECT 
			id, name, surname, nickname, email, password_hash, created_at
		FROM users
		WHERE users.email=$1;
	`

	err := userRepo.db.QueryRow(ctx, sqlQuery, email).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Nickname,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	return user, err
}
