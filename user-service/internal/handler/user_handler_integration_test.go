//go:build integration
// +build integration

package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/luke/sfconnect-backend/user-service/internal/repository"
	"github.com/luke/sfconnect-backend/user-service/internal/service"
)

var (
	testRouter *gin.Engine
	db         *sql.DB // Make db a package-level variable so it can be used in all test functions
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env.example")
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=user_service_user password=password dbname=user_db sslmode=disable"
	}

	var err error
	db, err = repository.NewPostgresDB()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test db: %v", err))
	}
	db.Exec("DROP TABLE IF EXISTS users;")
	// Enable uuid-ossp extension for UUID generation
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	// Run the same SQL migration as production
	db.Exec(`DROP TABLE IF EXISTS users;`)
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		full_name VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'buyer',
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);`)

	repo := repository.NewUserRepository(db)
	srv := service.NewUserService(repo)
	h := NewUserHandler(srv)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	RegisterRoutes(r, h)
	testRouter = r

	code := m.Run()
	db.Exec("DROP TABLE IF EXISTS users;")
	os.Exit(code)
}

func performRequest(method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader([]byte{})
	}
	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func TestRegisterLoginProfileFlow(t *testing.T) {
	// Register
	registerBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "TestPass123!",
		"full_name": "Test User",
	}
	resp := performRequest("POST", "/register", registerBody, "")
	if resp.Code != 201 {
		t.Logf("Register response: %s", resp.Body.String())
	}
	assert.Equal(t, 201, resp.Code)

	// Login
	loginBody := map[string]string{
		"email":    "testuser@example.com",
		"password": "TestPass123!",
	}
	resp = performRequest("POST", "/login", loginBody, "")
	if resp.Code != 200 {
		t.Logf("Login response: %s", resp.Body.String())
	}
	assert.Equal(t, 200, resp.Code)
	var loginResp map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &loginResp)
	token, _ := loginResp["token"].(string)
	assert.NotEmpty(t, token)

	// Get Profile
	resp = performRequest("GET", "/me", nil, token)
	if resp.Code != 200 {
		t.Logf("Profile response: %s", resp.Body.String())
	}
	assert.Equal(t, 200, resp.Code)
	var profileResp map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &profileResp)
	assert.Equal(t, "testuser@example.com", profileResp["email"])
	assert.Equal(t, "Test User", profileResp["full_name"])
}

func TestAdminRoleUpdate(t *testing.T) {
	// Register admin
	registerBody := map[string]string{
		"email":    "admin@example.com",
		"password": "AdminPass123!",
		"full_name": "Admin User",
	}
	resp := performRequest("POST", "/register", registerBody, "")
	if resp.Code != 201 {
		t.Logf("Admin Register response: %s", resp.Body.String())
	}
	assert.Equal(t, 201, resp.Code)

	// Login as admin
	loginBody := map[string]string{
		"email":    "admin@example.com",
		"password": "AdminPass123!",
	}
	resp = performRequest("POST", "/login", loginBody, "")
	if resp.Code != 200 {
		t.Logf("Admin Login response: %s", resp.Body.String())
	}
	assert.Equal(t, 200, resp.Code)
	var loginResp map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &loginResp)
	adminToken, _ := loginResp["token"].(string)

	// Promote to admin in DB
	var adminID string
	db.QueryRow("SELECT id FROM users WHERE email = $1", "admin@example.com").Scan(&adminID)
	db.Exec("UPDATE users SET role = 'admin' WHERE id = $1", adminID)

	// Re-login as admin to get new token with updated role
	resp = performRequest("POST", "/login", loginBody, "")
	if resp.Code != 200 {
		t.Logf("Admin Re-Login response: %s", resp.Body.String())
	}
	assert.Equal(t, 200, resp.Code)
	_ = json.Unmarshal(resp.Body.Bytes(), &loginResp)
	adminToken, _ = loginResp["token"].(string)

	// Register another user
	registerBody = map[string]string{
		"email":    "user2@example.com",
		"password": "User2Pass123!",
		"full_name": "User Two",
	}
	resp = performRequest("POST", "/register", registerBody, "")
	if resp.Code != 201 {
		t.Logf("User2 Register response: %s", resp.Body.String())
	}
	assert.Equal(t, 201, resp.Code)

	// Update role as admin
	var user2ID string
	db.QueryRow("SELECT id FROM users WHERE email = $1", "user2@example.com").Scan(&user2ID)
	roleUpdate := map[string]string{
		"role": "admin",
	}
	resp = performRequest("PUT", "/users/"+user2ID+"/role", roleUpdate, adminToken)
	if resp.Code != 200 {
		t.Logf("Role Update response: %s", resp.Body.String())
	}
	assert.Equal(t, 200, resp.Code)
}
