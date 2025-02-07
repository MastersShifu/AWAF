package handlers

import (
	"AuthMicroService/internal/models"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// SingUpHandler choose another answers
func (d *Dependencies) SingUpHandler(c *gin.Context) {
	var credentials models.Credentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if credentials.Name == "" || credentials.Password == "" || credentials.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not all fields"})
		return
	}

	if err := d.db.CheckUserExist(credentials.Name, credentials.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	credentials.Password = string(hashedPassword)

	code := utils.EncodeToString(6)

	if err := d.smtpClient.SendEmail(credentials.Email, code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	if err := d.rd.AddCode(credentials.Email, code, credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Verification code sanded to email",
	})
}
