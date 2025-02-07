package main

import (
	_ "AuthMicroService/docs" // импортируйте сгенерированные файлы Swagger
	"AuthMicroService/internal/config"
	"AuthMicroService/internal/handlers"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Auth Microservice API
// @version 1.0
// @description This is a sample server for an authentication microservice.
// @host localhost:8081
// @BasePath /

func main() {
	cfg, err := config.LoadConfig("config.yml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := gin.Default()

	dependencies := handlers.InitDependencies(cfg)
	handlers.InitializeRoutes(r, dependencies)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
