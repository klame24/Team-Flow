package routes

import (
	"team-flow/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterHealthRoutes(r *gin.Engine) {
	h := &handlers.PingHandler{}
	r.GET("/healthz", h.HealthZ)
}
