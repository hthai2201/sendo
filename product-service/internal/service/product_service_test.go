package service

import (
	"testing"

	"github.com/luke/sfconnect-backend/product-service/internal/models"
)

type mockRepo struct {
	products map[string]*models.Product
}

func (m *mockRepo) Create(product *models.Product) error {
	product.ID = "mock-id"
	m.products[product.ID] = product
	return nil
}
func (m *mockRepo) GetByID(id string) (*models.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, nil
	}
	return p, nil
}
func (m *mockRepo) List() ([]*models.Product, error) {
	var out []*models.Product
	for _, p := range m.products {
		out = append(out, p)
	}
	return out, nil
}
func (m *mockRepo) Update(product *models.Product) error {
	m.products[product.ID] = product
	return nil
}
func (m *mockRepo) Delete(id string) error {
	delete(m.products, id)
	return nil
}

type mockCache struct {
	cache map[string]*models.Product
}
func (m *mockCache) GetProduct(id string) (*models.Product, error) {
	return m.cache[id], nil
}
func (m *mockCache) SetProduct(product *models.Product) error {
	m.cache[product.ID] = product
	return nil
}
func (m *mockCache) DeleteProduct(id string) error {
	delete(m.cache, id)
	return nil
}

func newMockService() *ProductService {
	repo := &mockRepo{products: make(map[string]*models.Product)}
	cache := &mockCache{cache: make(map[string]*models.Product)}
	return NewProductService(repo, cache)
}

func TestCreateAndGetProduct(t *testing.T) {
	svc := newMockService()
	p := &models.Product{
		Name: "Test",
		Description: "Desc",
		Price: 10.5,
		ImageURL: "img",
		PartnerID: "partner",
		Stock: 5,
	}
	err := svc.Create(p)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	got, err := svc.GetByID(p.ID)
	if err != nil || got == nil {
		t.Fatalf("GetByID error: %v", err)
	}
	if got.Name != "Test" {
		t.Errorf("expected name Test, got %s", got.Name)
	}
}
