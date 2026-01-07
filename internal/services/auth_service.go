package services

import (
	"context"
	domains "team-flow/core/domains/user"
	"team-flow/internal/models"
	"team-flow/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(
		ctx context.Context, name, surname,
		nickname, email, password string,
	) (int, error)
	// Login
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (authService *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (authService *authService) Register(
	ctx context.Context, name, surname,
	nickname, email, password string,
) (int, error) {

	passwordHash, err := authService.hashPassword(password)
	if err != nil {
		return 0, err
	}

	userDomain := domains.User{
		Name:         name,
		Surname:      surname,
		Nickname:     nickname,
		Email:        email,
		PasswordHash: passwordHash,
	}

	userModel := models.User{
		Name:         userDomain.Name,
		Surname:      userDomain.Surname,
		Nickname:     userDomain.Nickname,
		Email:        userDomain.Email,
		PasswordHash: userDomain.PasswordHash,
	}

	userID, err := authService.userRepo.Create(ctx, &userModel)

	return userID, err
}
