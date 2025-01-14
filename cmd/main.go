package main

import (
	"log"

	"github.com/aminul-i-abid/url-shortener/internal/db"
	"github.com/aminul-i-abid/url-shortener/internal/middlewares"
	"github.com/aminul-i-abid/url-shortener/internal/services"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	router := gin.Default()

	// middleware to enforce rate limiting
	router.Use(middlewares.RateLimiterMiddleware())

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiV1 := router.Group("/api/v1")
	services.Routes(apiV1)

	db.ConnectDB()
	defer db.DisconnectDB()

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
