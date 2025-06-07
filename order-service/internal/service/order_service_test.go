package service

import (
	"testing"

	"github.com/luke/sfconnect-backend/order-service/internal/models"
	"github.com/stretchr/testify/assert"
)

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

func TestCreateOrder(t *testing.T) {
	repo := &mockRepo{
		CreateFn: func(order *models.Order) error {
			order.ID = "order1"
			return nil
		},
	}
	svc := NewOrderService(repo)
	order, err := svc.CreateOrder("user1", models.OrderCreateRequest{
		Items: []struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}{
			{ProductID: "p1", Quantity: 2},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, "user1", order.UserID)
	assert.Equal(t, 20.0, order.TotalAmount)
	assert.Equal(t, "pending", order.Status)
}

func TestListOrdersByUser(t *testing.T) {
	repo := &mockRepo{
		ListByUserFn: func(userID string) ([]*models.Order, error) {
			return []*models.Order{{ID: "o1", UserID: userID}}, nil
		},
	}
	svc := NewOrderService(repo)
	orders, err := svc.ListOrdersByUser("user1")
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, "user1", orders[0].UserID)
}
