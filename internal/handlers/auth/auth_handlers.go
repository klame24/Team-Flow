package handlers

import (
	"net/http"
	"team-flow/internal/handlers/auth/dto"
	"team-flow/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandlers interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authHandlers struct {
	authService services.AuthService
}

func NewAuthHandlers(authService services.AuthService) AuthHandlers {
	return &authHandlers{
		authService: authService,
	}
}

func (authHandlers *authHandlers) Register(c *gin.Context) {
	req := dto.RegisterRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	userID, err := authHandlers.authService.Register(
		c.Request.Context(),
		req.Name,
		req.Surname,
		req.Nickname,
		req.Email,
		req.Password,
	)
	if err != nil {
		// нужно обрабатывать ошибки грамотно, проверять наличие почты такой же и ника

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
	}

	c.JSON(http.StatusCreated, gin.H{"id": userID})
}

func (authHandlers *authHandlers) Login(c *gin.Context) {
	req := dto.LoginRequest{}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := authHandlers.authService.Login(c, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}
