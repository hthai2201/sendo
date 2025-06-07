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

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order for the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Param order body models.OrderCreateRequest true "Order create request"
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Router /orders [post]
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

// GetOrderByID godoc
// @Summary Get order by ID
// @Description Get an order by its ID (user or admin)
// @Tags orders
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} models.Order
// @Failure 403 {object} map[string]string
// @Router /orders/{id} [get]
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

// ListMyOrders godoc
// @Summary List my orders
// @Description List all orders for the authenticated user
// @Tags orders
// @Produce json
// @Success 200 {array} models.Order
// @Router /orders/my-orders [get]
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
