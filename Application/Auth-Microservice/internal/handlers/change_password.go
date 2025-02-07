package handlers

import (
	"AuthMicroService/internal/models"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (d *Dependencies) ResetPasswordHandler(c *gin.Context) {
	var credentials models.Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := d.db.CheckUserEmailExist(credentials.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code := utils.EncodeToString(6)

	if err := d.smtpClient.SendEmail(credentials.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}
