package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"team-flow/core/logger"
	"team-flow/internal/auth/jwt"
	"team-flow/internal/config"
	"team-flow/internal/repositories"
	"team-flow/internal/routes"
	"team-flow/internal/services"
	"team-flow/pkg/database/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}

	logger.LogInfo("Configuration loaded successfully!")

	godotenv.Load("env.env")

	ctx := context.Background()

	postgresDB, err := postgres.Connect(ctx, cfg.Database)
	if err != nil {
		logger.LogError(err)
		os.Exit(1)
	}

	logger.LogInfo("Postgres connected successfully!")

	defer postgresDB.Close()

	jwtManager := jwt.NewManager(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)

	// инициализация репозиториев
	userRepo := repositories.NewUserRepository(postgresDB.GetPool())
	tokenRepo := repositories.NewTokenRepository(postgresDB.GetPool())

	// инициализая сервисов
	authService := services.NewAuthService(userRepo, tokenRepo, jwtManager)

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())

	routes.RegisterAllRoutes(r, authService)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		logger.LogInfo("Starting server on port " + cfg.Server.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError(err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.LogInfo("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError(err)
	}

	logger.LogInfo("Server exited gracefully")

}
