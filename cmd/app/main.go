package main

import (
	"context"
	"fmt"
	"os"
	"team-flow/internal/auth/jwt"
	"team-flow/internal/repositories"
	"team-flow/internal/routes"
	"team-flow/internal/services"
	"team-flow/pkg/database/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("env.env")

	ctx := context.Background()

	pool, err := postgres.ConnectDB(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connect to PostgresDB")

	jwtManager := jwt.NewManager(
		os.Getenv("SECRET_KEY"),
		15*time.Minute,
		30*24*time.Hour,
	)

	// инициализация репозиториев
	userRepo := repositories.NewUserRepository(pool)
	tokenRepo := repositories.NewTokenRepository(pool)

	// инициализая сервисов
	authService := services.NewAuthService(userRepo, tokenRepo, jwtManager)

	r := gin.Default()

	routes.RegisterAllRoutes(r, authService)

	if err := r.Run(":5050"); err != nil {
		panic(err)
	}

}
