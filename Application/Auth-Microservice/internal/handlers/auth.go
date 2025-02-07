package handlers

import (
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse представляет ответ на запрос аутентификации
type AuthResponse struct {
	Token string `json:"token"`
}

// ErrorResponse представляет ошибку
type ErrorResponse struct {
	Message string `json:"message"`
}

// AuthHandler обрабатывает запросы на аутентификацию
// @Summary Authenticate user
// @Description Authenticate user with token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   Authorization header string true "Bearer token"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /auth [post]

// AuthHandler choose another answers, check token blacklist
func (d *Dependencies) AuthHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		log.Println("Invalid token format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		return
	}

	if err := utils.VerifyToken(tokenString); err != nil {
		log.Println("Invalid token:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User authenticated",
	})
}
