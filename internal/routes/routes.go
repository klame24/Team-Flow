package routes

import (
	"team-flow/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterAllRoutes(r *gin.Engine, authService services.AuthService) {
	RegisterAuthRoutes(r, authService)
	RegisterHealthRoutes(r)
}
