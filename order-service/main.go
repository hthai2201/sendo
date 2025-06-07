package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/luke/sfconnect-backend/order-service/internal/handler"
	"github.com/luke/sfconnect-backend/order-service/internal/repository"
	"github.com/luke/sfconnect-backend/order-service/internal/service"
	"github.com/luke/sfconnect-backend/order-service/pkg/utils"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	r := gin.Default()
	handler.RegisterOrderRoutes(r, svc, utils.AuthMiddleware())
	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
