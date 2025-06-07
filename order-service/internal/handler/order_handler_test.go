package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luke/sfconnect-backend/order-service/internal/models"
	"github.com/luke/sfconnect-backend/order-service/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mockRepo{
		CreateFn: func(order *models.Order) error {
			order.ID = "order1"
			return nil
		},
	}
	svc := service.NewOrderService(repo)
	r := gin.Default()
	h := NewOrderHandler(svc)
	r.POST("/orders", func(c *gin.Context) {
		c.Set("user_id", "user1")
		h.CreateOrder(c)
	})
	body, _ := json.Marshal(models.OrderCreateRequest{
		Items: []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}{
			{ProductID: "p1", Quantity: 2},
		},
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/orders", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

type mockRepo struct {
	CreateFn      func(order *models.Order) error
	GetByIDFn     func(id string) (*models.Order, error)
	ListByUserFn  func(userID string) ([]*models.Order, error)
	ListAllFn     func() ([]*models.Order, error)
	UpdateStatusFn func(orderID, status string) error
}

func (m *mockRepo) Create(order *models.Order) error                 { return m.CreateFn(order) }
func (m *mockRepo) GetByID(id string) (*models.Order, error)         { return m.GetByIDFn(id) }
func (m *mockRepo) ListByUser(userID string) ([]*models.Order, error) { return m.ListByUserFn(userID) }
func (m *mockRepo) ListAll() ([]*models.Order, error) {
	if m.ListAllFn != nil {
		return m.ListAllFn()
	}
	return nil, nil
}
func (m *mockRepo) UpdateStatus(orderID, status string) error {
	if m.UpdateStatusFn != nil {
		return m.UpdateStatusFn(orderID, status)
	}
	return nil
}
