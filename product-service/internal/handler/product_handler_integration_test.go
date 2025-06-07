package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/luke/sfconnect-backend/product-service/internal/cache"
	"github.com/luke/sfconnect-backend/product-service/internal/models"
	"github.com/luke/sfconnect-backend/product-service/internal/repository"
	"github.com/luke/sfconnect-backend/product-service/internal/service"
	"github.com/redis/go-redis/v9"
)

var testRouter *gin.Engine
var db *sql.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env.example")
	db, _ = sql.Open("postgres", os.Getenv("PRODUCT_DB_DSN"))
	db.Exec("DROP TABLE IF EXISTS products;")
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";"`)
	db.Exec(`CREATE TABLE IF NOT EXISTS products (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price NUMERIC(12,2) NOT NULL,
		image_url VARCHAR(512),
		partner_id UUID NOT NULL,
		stock INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`)
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	repo := repository.NewProductRepository(db)
	cache := cache.NewProductCache(rdb)
	svc := service.NewProductService(repo, cache)
	h := NewProductHandler(svc)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/products", h.Create)
	r.GET("/products/:id", h.GetByID)
	r.GET("/products", h.List)
	r.PUT("/products/:id", h.Update)
	r.DELETE("/products/:id", h.Delete)
	testRouter = r
	code := m.Run()
	db.Exec("DROP TABLE IF EXISTS products;")
	os.Exit(code)
}

func performRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}
	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func TestProductCRUD(t *testing.T) {
	// Create
	createBody := map[string]interface{}{
		"name": "Apple",
		"description": "Fresh apple",
		"price": 2.5,
		"image_url": "img",
		"stock": 10,
	}
	resp := performRequest("POST", "/products", createBody)
	if resp.Code != 201 {
		t.Fatalf("Create failed: %s", resp.Body.String())
	}
	var prod models.Product
	_ = json.Unmarshal(resp.Body.Bytes(), &prod)
	// Get
	resp = performRequest("GET", "/products/"+prod.ID, nil)
	if resp.Code != 200 {
		t.Fatalf("Get failed: %s", resp.Body.String())
	}
	// List
	resp = performRequest("GET", "/products", nil)
	if resp.Code != 200 {
		t.Fatalf("List failed: %s", resp.Body.String())
	}
	// Update
	updateBody := map[string]interface{}{"name": "Green Apple"}
	resp = performRequest("PUT", "/products/"+prod.ID, updateBody)
	if resp.Code != 200 {
		t.Fatalf("Update failed: %s", resp.Body.String())
	}
	// Delete
	resp = performRequest("DELETE", "/products/"+prod.ID, nil)
	if resp.Code != 200 {
		t.Fatalf("Delete failed: %s", resp.Body.String())
	}
}
