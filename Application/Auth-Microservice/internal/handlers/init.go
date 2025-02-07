package handlers

import (
	"AuthMicroService/internal/config"
	"AuthMicroService/internal/database/cache"
	"AuthMicroService/internal/database/postgresql"
	"AuthMicroService/internal/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type Dependencies struct {
	db         *postgresql.DBClient
	rd         *cache.RDClient
	smtpClient *utils.SMTPClient
}

func InitDependencies(cfg *config.Config) *Dependencies {
	db, err := postgresql.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	client, err := cache.ConnectToRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to cache: %v", err)
	}

	smtpClient, err := utils.ConnectToSMTP(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to smtp: %v", err)
	}
	return &Dependencies{
		db:         db,
		rd:         client,
		smtpClient: smtpClient,
	}
}

func InitializeRoutes(r *gin.Engine, d *Dependencies) {
	r.POST("/auth", d.AuthHandler)
	r.POST("/authorization", d.AuthorizationHandler)
	r.POST("/registration", d.SingUpHandler)
	r.POST("/logout", d.LogoutHandler)
	r.POST("/refresh", d.RefreshTokens)
	r.POST("/code/ver", d.VerSingupHandler)
	r.POST("/code/resend")
}
