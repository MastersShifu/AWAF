package handlers

import (
	"AuthMicroService/internal/models"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func (d *Dependencies) RefreshTokens(c *gin.Context) {
	jwtToken := c.GetHeader("jwt")
	rtToken := c.GetHeader("rt")

	tokens := models.Tokens{
		"rt":  rtToken,
		"jwt": strings.TrimPrefix(jwtToken, "Bearer"),
	}

	claims, err := utils.ExtractClaims(tokens["jwt"].(string))
	if err != nil {
		log.Println("Extracting claims error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token error"})
		return
	}

	err = d.rd.CheckTokensBlackList(claims["id"].(string), tokens)
	if err != nil {
		log.Println("Tokens in black list: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token error"})
		return
	}

	err = d.db.GetCredentialsVersion(claims["id"].(string), int(claims["version"].(float64)))
	if err != nil {
		log.Println("Getting credentials version: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token error"})
		return
	}

	jwtToken, err = utils.CreateJWT(claims["id"].(string), int(claims["version"].(float64)))
	if err != nil {
		log.Println("Creating token error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token error"})
		return
	}

	rtToken, err = utils.CreateRT(claims["id"].(string))
	if err != nil {
		log.Println("Creating token error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token error"})
		return
	}

	err = d.rd.AddTokensToBlackList(claims["id"].(string), tokens)
	if err != nil {
		log.Println("Adding tokens to black list: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token error"})
		return
	}

	c.JSON(200, gin.H{
		"rt":  rtToken,
		"jwt": jwtToken,
	})
}
