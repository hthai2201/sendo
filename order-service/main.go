package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/luke/sfconnect-backend/order-service/docs"
	"github.com/luke/sfconnect-backend/order-service/internal/handler"
	"github.com/luke/sfconnect-backend/order-service/internal/repository"
	"github.com/luke/sfconnect-backend/order-service/internal/service"
	"github.com/luke/sfconnect-backend/order-service/pkg/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SFConnect Order Service API
// @version 1.0
// @description API documentation for the Order Service.
// @host localhost:8083
// @BasePath /

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	handler.RegisterOrderRoutes(r, svc, utils.AuthMiddleware())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
