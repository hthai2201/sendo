package service

import (
	"errors"

	"github.com/luke/sfconnect-backend/order-service/internal/models"
	"github.com/luke/sfconnect-backend/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(userID string, req models.OrderCreateRequest) (*models.Order, error)
	GetOrderByID(orderID string, userID string, isAdmin bool) (*models.Order, error)
	ListOrdersByUser(userID string) ([]*models.Order, error)
	UpdateOrderStatus(orderID, nextStatus, userID, role string) error
	ListAllOrders() ([]*models.Order, error)
	CalculateAndSetCommission(orderID string) error
}

type orderService struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderService{repo: repo}
}

func (s *orderService) CreateOrder(userID string, req models.OrderCreateRequest) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	order := &models.Order{
		UserID: userID,
		Status: "pending",
	}
	var total float64
	for _, item := range req.Items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}
		// TODO: fetch product price from product-service (mocked as 10.0 for now)
		unitPrice := 10.0
		order.Items = append(order.Items, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: unitPrice,
		})
		total += unitPrice * float64(item.Quantity)
	}
	order.TotalAmount = total
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *orderService) GetOrderByID(orderID string, userID string, isAdmin bool) (*models.Order, error) {
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("forbidden")
	}
	return order, nil
}

func (s *orderService) ListOrdersByUser(userID string) ([]*models.Order, error) {
	return s.repo.ListByUser(userID)
}

func (s *orderService) ListAllOrders() ([]*models.Order, error) {
	return s.repo.ListAll()
}

const commissionRate = 0.1 // 10% commission

func (s *orderService) CalculateAndSetCommission(orderID string) error {
	order, err := s.repo.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.Status != StatusDelivered {
		return errors.New("order not delivered")
	}
	commission := order.TotalAmount * commissionRate
	if err := s.repo.UpdateCommission(orderID, commission); err != nil {
		return err
	}
	return nil
}

var (
	StatusPending         = "Pending"
	StatusProcessing      = "Processing"
	StatusReadyForDelivery = "ReadyForDelivery"
	StatusShipped         = "Shipped"
	StatusDelivered       = "Delivered"
	StatusCanceled        = "Canceled"
)

func validStatusTransition(current, next, role string) bool {
	switch current {
	case StatusPending:
		if next == StatusProcessing && role == "admin" { return true }
		if next == StatusCanceled && role == "admin" { return true }
	case StatusProcessing:
		if next == StatusReadyForDelivery && role == "partner" { return true }
		if next == StatusCanceled && role == "admin" { return true }
	case StatusReadyForDelivery:
		if next == StatusShipped && role == "admin" { return true }
	case StatusShipped:
		if next == StatusDelivered && role == "buyer" { return true }
	}
	return false
}

func (s *orderService) UpdateOrderStatus(orderID, nextStatus, userID, role string) error {
	order, err := s.repo.GetByID(orderID)
	if err != nil { return err }
	if !validStatusTransition(order.Status, nextStatus, role) {
		return errors.New("invalid status transition")
	}
	if nextStatus == StatusDelivered {
		if err := s.repo.UpdateStatus(orderID, nextStatus); err != nil {
			return err
		}
		// Calculate commission after delivery
		if err := s.CalculateAndSetCommission(orderID); err != nil {
			return err
		}
		return nil
	}
	if err := s.repo.UpdateStatus(orderID, nextStatus); err != nil {
		return err
	}
	return nil
}
