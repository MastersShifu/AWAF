package handlers

import (
	"AuthMicroService/internal/dto"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type TokenRequest struct {
	Token string `json:"token"`
}

// AuthorizationResponse представляет ответ на запрос авторизации
type AuthorizationResponse struct {
	Message string `json:"message"`
}

// AuthorizationHandler обрабатывает запросы на авторизацию
// @Summary Authorize user
// @Description Authorize user with credentials
// @Tags authorization
// @Accept  json
// @Produce  json
// @Param   credentials  body  AuthRequest  true  "User credentials"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Router /authorization [post]
// AuthorizationHandler choose another answers
func (d *Dependencies) AuthorizationHandler(c *gin.Context) {
	var user dto.UserDTO

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println("Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	authUser, err := d.db.AuthenticateUser(user.UsernameOrEmail, user.Password)
	if err != nil {
		log.Println("Authentication failed:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	jwt, err := utils.CreateJWT(authUser.ID, authUser.Version)
	if err != nil {
		log.Println("Creating token error: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creating token failed"})
		return
	}

	rt, err := utils.CreateRT(authUser.ID)
	if err != nil {
		log.Println("Creating token error: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Creating token failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"jwt":     jwt,
		"rt":      rt,
	})
}
