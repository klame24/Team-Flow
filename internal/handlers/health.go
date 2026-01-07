package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingHandler struct{}

func (h *PingHandler) HealthZ(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "server is healthy"})
}
