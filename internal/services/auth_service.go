package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	domains "team-flow/core/domains/user"
	"team-flow/internal/auth/jwt"
	"team-flow/internal/models"
	"team-flow/internal/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(
		ctx context.Context, name, surname,
		nickname, email, password string,
	) (int32, error)
	Login(ctx context.Context, email, password string) (string, string, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	tokenRepo  repositories.TokenRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo repositories.UserRepository, tokenRepo repositories.TokenRepository, jwtManager *jwt.Manager) AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
	}
}

func (authService *authService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (authService *authService) VerifyPassowrd(password, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	return err == nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (authService *authService) Register(
	ctx context.Context, name, surname,
	nickname, email, password string,
) (int32, error) {

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

func (authService *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := authService.userRepo.GetByEmail(ctx, email)
	if err != nil {
		// user not found
		return "", "", err
	}

	// проверить пароль из бд с паролем от транспорта
	if !authService.VerifyPassowrd(password, user.PasswordHash) {
		// wrong password
		return "", "", err
	}
	// сгенерить токены
	accessToken, err := authService.jwtManager.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		// cant create accessToken
		return "", "", err
	}

	refreshToken, err := authService.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		// cant create refreshToken
		return "", "", err
	}

	hashRefreshToken := hashToken(refreshToken)

	expiresAt := time.Now().Add(authService.jwtManager.RefreshTokenTTL)

	// сохранить токены в refresh_tokens table
	_, err = authService.tokenRepo.SaveRefreshTokens(ctx, user.ID, string(hashRefreshToken), expiresAt)
	if err != nil {
		// cant save refreshtoken
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
