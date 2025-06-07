package service

import (
	"github.com/luke/sfconnect-backend/product-service/internal/cache"
	"github.com/luke/sfconnect-backend/product-service/internal/models"
	"github.com/luke/sfconnect-backend/product-service/internal/repository"
)

type ProductService struct {
	repo  repository.ProductRepository
	cache cache.ProductCache
}

func NewProductService(repo repository.ProductRepository, cache cache.ProductCache) *ProductService {
	return &ProductService{repo: repo, cache: cache}
}

func (s *ProductService) Create(product *models.Product) error {
	if err := s.repo.Create(product); err != nil {
		return err
	}
	return s.cache.SetProduct(product)
}

func (s *ProductService) GetByID(id string) (*models.Product, error) {
	// Try cache first
	p, err := s.cache.GetProduct(id)
	if err != nil {
		return nil, err
	}
	if p != nil {
		return p, nil
	}
	// Fallback to DB
	p, err = s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	_ = s.cache.SetProduct(p)
	return p, nil
}

func (s *ProductService) List() ([]*models.Product, error) {
	return s.repo.List()
}

func (s *ProductService) Update(product *models.Product) error {
	if err := s.repo.Update(product); err != nil {
		return err
	}
	return s.cache.SetProduct(product)
}

func (s *ProductService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return s.cache.DeleteProduct(id)
}
