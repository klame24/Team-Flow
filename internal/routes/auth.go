package routes

import (
	handlers "team-flow/internal/handlers/auth"
	"team-flow/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, authService services.AuthService) {
	authHandlers := handlers.NewAuthHandlers(authService)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandlers.Register)
		authGroup.POST("/login", authHandlers.Login)
	}
}
