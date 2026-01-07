package main

import (
	"context"
	"fmt"
	"team-flow/internal/repositories"
	"team-flow/internal/routes"
	"team-flow/internal/services"
	"team-flow/pkg/database/postgres"

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

	// инициализация репозиториев
	userRepo := repositories.NewUserRepository(pool)

	// инициализая сервисов
	authService := services.NewAuthService(userRepo)

	r := gin.Default()

	routes.RegisterAllRoutes(r, authService)

	if err := r.Run(":5050"); err != nil {
		panic(err)
	}

}
