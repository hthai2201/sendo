package repository

import (
	"database/sql"

	"github.com/luke/sfconnect-backend/order-service/internal/models"
)

type OrderRepository interface {
	Create(order *models.Order) error
	GetByID(id string) (*models.Order, error)
	ListByUser(userID string) ([]*models.Order, error)
	UpdateStatus(orderID, status string) error
	UpdateCommission(orderID string, commission float64) error
	ListAll() ([]*models.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *models.Order) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	orderQuery := `INSERT INTO orders (user_id, total_amount, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err = tx.QueryRow(orderQuery, order.UserID, order.TotalAmount, order.Status).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}
	for _, item := range order.Items {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4) RETURNING id`
		err = tx.QueryRow(itemQuery, order.ID, item.ProductID, item.Quantity, item.UnitPrice).Scan(&item.ID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *orderRepository) GetByID(id string) (*models.Order, error) {
	order := &models.Order{}
	orderQuery := `SELECT id, user_id, total_amount, status, created_at, updated_at FROM orders WHERE id = $1`
	err := r.db.QueryRow(orderQuery, id).Scan(&order.ID, &order.UserID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}
	itemQuery := `SELECT id, order_id, product_id, quantity, unit_price FROM order_items WHERE order_id = $1`
	rows, err := r.db.Query(itemQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	return order, nil
}

func (r *orderRepository) ListByUser(userID string) ([]*models.Order, error) {
	query := `SELECT id FROM orders WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*models.Order
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		order, err := r.GetByID(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(orderID, status string) error {
	_, err := r.db.Exec(`UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`, status, orderID)
	return err
}

func (r *orderRepository) UpdateCommission(orderID string, commission float64) error {
	_, err := r.db.Exec(`UPDATE orders SET commission = $1, updated_at = NOW() WHERE id = $2`, commission, orderID)
	return err
}

func (r *orderRepository) ListAll() ([]*models.Order, error) {
	query := `SELECT id FROM orders`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orders []*models.Order
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		order, err := r.GetByID(id)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
