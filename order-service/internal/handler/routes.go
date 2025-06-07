package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/luke/sfconnect-backend/order-service/internal/service"
)

func RegisterOrderRoutes(r *gin.Engine, s service.OrderService, authMiddleware gin.HandlerFunc) {
	h := NewOrderHandler(s)
	orders := r.Group("/orders")
	orders.Use(authMiddleware)
	{
		orders.POST("", h.CreateOrder)
		orders.GET("/my-orders", h.ListMyOrders)
		orders.GET(":id", h.GetOrderByID)
		orders.PUT(":id/confirm-ready", h.ConfirmReady)
		orders.PUT(":id/confirm-delivery", h.ConfirmDelivery)
	}
	// admin only
	r.GET("/orders", authMiddleware, h.ListAllOrders)
}
