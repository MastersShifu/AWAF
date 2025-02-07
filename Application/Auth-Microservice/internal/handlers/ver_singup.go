package handlers

import (
	"AuthMicroService/internal/models"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (d *Dependencies) VerSingupHandler(c *gin.Context) {
	var code models.Code

	if err := c.ShouldBindJSON(&code); err != nil {
		log.Println("Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	credentials, err := d.rd.CheckCode(code.Email, code.Code)
	if err != nil {
		log.Println(err)
		return
	}

	if err := d.db.RegisterNewUser(credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.SendMessage(credentials.Name, credentials.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "user singed up",
		"user_id": credentials.ID,
	})
}
