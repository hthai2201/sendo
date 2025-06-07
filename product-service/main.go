package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luke/sfconnect-backend/product-service/internal/cache"
	. "github.com/luke/sfconnect-backend/product-service/internal/handler"
	"github.com/luke/sfconnect-backend/product-service/internal/repository"
	"github.com/luke/sfconnect-backend/product-service/internal/service"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SFConnect Product Service API
// @version 1.0
// @description API documentation for the Product Service.
// @host localhost:8082
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @security Bearer
// @schemes http https
// @contact.name SFConnect Team

func main() {
	_ = godotenv.Load()
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	repo := repository.NewProductRepository(db)
	cache := cache.NewProductCache(rdb)
	svc := service.NewProductService(repo, cache)
	h := NewProductHandler(svc)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/products", AuthMiddleware(), AuthorizeRole("partner"), h.Create)
	r.GET("/products/:id", h.GetByID)
	r.GET("/products", h.List)
	r.PUT("/products/:id", AuthMiddleware(), AuthorizeRole("partner"), h.Update)
	r.DELETE("/products/:id", AuthMiddleware(), AuthorizeRole("partner"), h.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	r.Run(":" + port)
}
