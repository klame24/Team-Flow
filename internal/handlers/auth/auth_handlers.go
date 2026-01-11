package handlers

import (
	"net/http"
	"team-flow/core/errors"
	"team-flow/core/logger"
	"team-flow/core/validator"
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

	if valErr := validator.ValidateStruct(&req); valErr != nil { // почему передаем структуру если функция принимает интерфейс?
		c.JSON(errors.ToHTTPCode(valErr.Code), gin.H{
			"error":   valErr.Code,
			"message": valErr.Message,
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
		logger.LogError(err)

		var statusCode int
		var errorCode, message string

		if err.Error() == "duplicate email" {
			statusCode = http.StatusConflict
			errorCode = errors.ErrCodeDuplicate
			message = "Email already exists"
		} else if err.Error() == "duplicate nickname" {
			statusCode = http.StatusConflict
			errorCode = errors.ErrCodeDuplicate
			message = "Nickname already exists"
		} else {
			statusCode = http.StatusInternalServerError
			errorCode = errors.ErrCodeInternal
			message = "Internal server error"
		}

		c.JSON(statusCode, gin.H{
			"error":   errorCode,
			"message": message,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": userID,
		"message": "User registered successfully",
	})
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

	if valErr := validator.ValidateStruct(&req); valErr != nil {
		c.JSON(errors.ToHTTPCode(valErr.Code), gin.H{
			"error":   valErr.Code,
			"message": valErr.Message,
		})

		return
	}

	accessToken, refreshToken, err := authHandlers.authService.Login(c, req.Email, req.Password)
	if err != nil {
		logger.LogError(err)

		var statusCode int
		var errorCode, message string

		switch err.Error() {
		case "user not found":
			statusCode = http.StatusUnauthorized
			errorCode = errors.ErrCodeUnauthorized
			message = "Invalid data"
		
		case "wrong password":
			statusCode = http.StatusUnauthorized
			errorCode = errors.ErrCodeUnauthorized
			message = "Invalid data"

		default:
			statusCode = http.StatusInternalServerError
			errorCode = errors.ErrCodeInternal
			message = "Internal server error"
		}

		c.JSON(statusCode, gin.H{
			"error": errorCode,
			"message": message,
		})

		return
		
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
	})
}
