package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luke/sfconnect-backend/order-service/internal/models"
	"github.com/luke/sfconnect-backend/order-service/internal/service"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(s service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

// POST /orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	var req models.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	order, err := h.service.CreateOrder(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// GET /orders/:id
func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")
	isAdmin := c.GetString("role") == "admin"
	order, err := h.service.GetOrderByID(orderID, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// GET /orders/my-orders
func (h *OrderHandler) ListMyOrders(c *gin.Context) {
	userID := c.GetString("user_id")
	orders, err := h.service.ListOrdersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// PUT /orders/:id/confirm-ready
func (h *OrderHandler) ConfirmReady(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if role != "partner" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	err := h.service.UpdateOrderStatus(orderID, "ReadyForDelivery", userID, role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "order marked as ready for delivery"})
}

// PUT /orders/:id/confirm-delivery
func (h *OrderHandler) ConfirmDelivery(c *gin.Context) {
	orderID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")
	if role != "buyer" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	err := h.service.UpdateOrderStatus(orderID, "Delivered", userID, role)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "order marked as delivered"})
}

// GET /orders (admin only)
func (h *OrderHandler) ListAllOrders(c *gin.Context) {
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(403, gin.H{"error": "forbidden"})
		return
	}
	orders, err := h.service.ListAllOrders()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, orders)
}
