package models

type Order struct {
	ID         string  `db:"id" json:"id"`
	UserID     string  `db:"user_id" json:"user_id"`
	TotalAmount float64 `db:"total_amount" json:"total_amount"`
	Status     string  `db:"status" json:"status"`
	Commission float64 `db:"commission" json:"commission"`
	CreatedAt  string  `db:"created_at" json:"created_at"`
	UpdatedAt  string  `db:"updated_at" json:"updated_at"`
	Items      []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        string  `db:"id" json:"id"`
	OrderID   string  `db:"order_id" json:"order_id"`
	ProductID string  `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	UnitPrice float64 `db:"unit_price" json:"unit_price"`
}

// OrderCreateRequest for creating an order
// swagger:model
type OrderCreateRequest struct {
	Items []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}
