package main

import (
	"github.com/gin-gonic/gin"
	"github.com/luke/sfconnect-backend/chatbot-service/internal/handler"
)

func main() {
	r := gin.Default()
	handler.RegisterChatbotRoutes(r)
	r.Run()
}
