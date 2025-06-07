package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestExtractOrderID(t *testing.T) {
	cases := []struct {
		msg string
		expect string
	}{
		{"Đơn hàng 123 của tôi đâu?", "123"},
		{"order 456 status", "456"},
		{"mã đơn hàng 789", "789"},
		{"không có mã", ""},
	}
	for _, c := range cases {
		assert.Equal(t, c.expect, extractOrderID(c.msg))
	}
}

func TestChatbotQueryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/chatbot/query", ChatbotQueryHandler)
	body, _ := json.Marshal(map[string]string{"message": "Đơn hàng 123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chatbot/query", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestChatbotQueryHandler_OrderServiceIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	// Mock order-service
	orderSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"Delivered"}`))
	}))
	defer orderSrv.Close()
	os.Setenv("ORDER_SERVICE_URL", orderSrv.URL)
	r.POST("/chatbot/query", ChatbotQueryHandler)
	body, _ := json.Marshal(map[string]string{"message": "Đơn hàng 123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chatbot/query", bytes.NewReader(body))
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Delivered")
}
