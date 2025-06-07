package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

type ChatbotQuery struct {
	Message string `json:"message"`
}

type ChatbotReply struct {
	Reply string `json:"reply"`
}

func RegisterChatbotRoutes(r *gin.Engine) {
	r.POST("/chatbot/query", ChatbotQueryHandler)
}

func ChatbotQueryHandler(c *gin.Context) {
	var req ChatbotQuery
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	orderID := extractOrderID(req.Message)
	if orderID != "" {
		orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
		if orderServiceURL == "" {
			orderServiceURL = "http://order-service:8080" // default for Docker Compose
		}
		resp, err := http.Get(fmt.Sprintf("%s/orders/%s", orderServiceURL, orderID))
		if err != nil || resp.StatusCode != 200 {
			c.JSON(http.StatusOK, ChatbotReply{Reply: "Xin lỗi, không thể truy vấn trạng thái đơn hàng lúc này."})
			return
		}
		defer resp.Body.Close()
		var order struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
			c.JSON(http.StatusOK, ChatbotReply{Reply: "Xin lỗi, không thể đọc trạng thái đơn hàng."})
			return
		}
		reply := fmt.Sprintf("Đơn hàng #%s hiện đang ở trạng thái: %s.", orderID, order.Status)
		c.JSON(http.StatusOK, ChatbotReply{Reply: reply})
		return
	}
	c.JSON(http.StatusOK, ChatbotReply{Reply: "Xin lỗi, tôi không tìm thấy mã đơn hàng trong câu hỏi của bạn."})
}

func extractOrderID(msg string) string {
	re := regexp.MustCompile(`(?i)(?:order|đơn hàng|mã đơn hàng)[^\d]*(\d+)`)
	match := re.FindStringSubmatch(msg)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
