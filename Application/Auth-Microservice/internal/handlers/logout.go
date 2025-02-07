package handlers

import (
	"AuthMicroService/internal/models"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// LogoutHandler choose another answers
func (d *Dependencies) LogoutHandler(c *gin.Context) {
	jwtToken := c.GetHeader("jwt")
	rtToken := c.GetHeader("rt")

	tokens := models.Tokens{
		"rt":  rtToken,
		"jwt": strings.TrimPrefix(jwtToken, "Bearer"),
	}

	claims, err := utils.ExtractClaims(tokens["jwt"].(string))
	if err != nil {
		log.Println("Extracting claims error: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Logout error"})
		return
	}

	err = d.rd.AddTokensToBlackList(claims["id"].(string), tokens)
	if err != nil {
		log.Println("Adding tokens to black list error: ", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Logout error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user logout successful",
	})
}
