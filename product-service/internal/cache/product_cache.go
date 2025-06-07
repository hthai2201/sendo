package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/luke/sfconnect-backend/product-service/internal/models"
	"github.com/redis/go-redis/v9"
)

type ProductCache interface {
	GetProduct(id string) (*models.Product, error)
	SetProduct(product *models.Product) error
	DeleteProduct(id string) error
}

type productCache struct {
	rdb *redis.Client
	ctx context.Context
}

func NewProductCache(rdb *redis.Client) ProductCache {
	return &productCache{rdb: rdb, ctx: context.Background()}
}

func (c *productCache) GetProduct(id string) (*models.Product, error) {
	val, err := c.rdb.Get(c.ctx, "product:"+id).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var p models.Product
	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (c *productCache) SetProduct(product *models.Product) error {
	b, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return c.rdb.Set(c.ctx, "product:"+product.ID, b, 10*time.Minute).Err()
}

func (c *productCache) DeleteProduct(id string) error {
	return c.rdb.Del(c.ctx, "product:"+id).Err()
}
