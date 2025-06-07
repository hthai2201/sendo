package repository

import (
	"database/sql"

	"github.com/luke/sfconnect-backend/product-service/internal/models"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id string) (*models.Product, error)
	List() ([]*models.Product, error)
	Update(product *models.Product) error
	Delete(id string) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	query := `INSERT INTO products (name, description, price, image_url, partner_id, stock) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, product.Name, product.Description, product.Price, product.ImageURL, product.PartnerID, product.Stock).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) GetByID(id string) (*models.Product, error) {
	product := &models.Product{}
	query := `SELECT id, name, description, price, image_url, partner_id, stock, created_at, updated_at FROM products WHERE id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.PartnerID, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) List() ([]*models.Product, error) {
	query := `SELECT id, name, description, price, image_url, partner_id, stock, created_at, updated_at FROM products`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []*models.Product
	for rows.Next() {
		p := &models.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.PartnerID, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepository) Update(product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, price = $3, image_url = $4, stock = $5, updated_at = NOW() WHERE id = $6 RETURNING updated_at`
	return r.db.QueryRow(query, product.Name, product.Description, product.Price, product.ImageURL, product.Stock, product.ID).Scan(&product.UpdatedAt)
}

func (r *productRepository) Delete(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
