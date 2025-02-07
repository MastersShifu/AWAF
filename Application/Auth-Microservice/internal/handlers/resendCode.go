package handlers

import (
	"AuthMicroService/internal/models"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func (d *Dependencies) ResendVerCode(c *gin.Context) {
	var code models.Code

	if err := c.ShouldBindJSON(&code); err != nil {
		log.Println("Invalid request format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	newCode := utils.EncodeToString(6)

	if err := d.smtpClient.SendEmail(code.Email, newCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := d.rd.UpdateCode(code.Email, newCode); err != nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification code sanded to email",
	})

}
