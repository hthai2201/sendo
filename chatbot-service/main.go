package main

import (
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/luke/sfconnect-backend/chatbot-service/docs"
	"github.com/luke/sfconnect-backend/chatbot-service/internal/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SFConnect Chatbot Service API
// @version 1.0
// @description API documentation for the Chatbot Service.
// @host localhost:8084
// @BasePath /

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.RegisterChatbotRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	r.Run(":" + port)
}
