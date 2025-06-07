package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luke/sfconnect-backend/user-service/internal/handler"
	"github.com/luke/sfconnect-backend/user-service/internal/repository"
	"github.com/luke/sfconnect-backend/user-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/example/basic/docs"
)

// @title SFConnect User Service API
// @version 1.0
// @description API documentation for the User Service.
// @host localhost:8080
// @BasePath /

func main() {
	_ = godotenv.Load()
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.RegisterRoutes(r, h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
