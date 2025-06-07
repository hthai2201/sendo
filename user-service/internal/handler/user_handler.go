package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luke/sfconnect-backend/user-service/internal/models"
	"github.com/luke/sfconnect-backend/user-service/internal/service"
	"github.com/luke/sfconnect-backend/user-service/pkg/utils"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func jsonError(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "User registration info"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, 400, "invalid request")
		return
	}
	user, err := h.service.Register(req.Email, req.Password, req.FullName)
	if err != nil {
		jsonError(c, 400, err.Error())
		return
	}
	c.JSON(201, user)
}

// Login godoc
// @Summary Login and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "User login info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, 400, "invalid request")
		return
	}
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		jsonError(c, 401, err.Error())
		return
	}
	c.JSON(200, gin.H{"token": token})
}

// AuthMiddleware validates JWT and injects user_id and role into context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		}
		claims, err := utils.ParseJWT(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AuthorizeRole checks if the user has the required role
func AuthorizeRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.service.GetProfile(userID.(string))
	if err != nil {
		jsonError(c, 404, "user not found")
		return
	}
	c.JSON(200, user)
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags users
// @Accept json
// @Produce json
// @Param body body models.UpdateProfileRequest true "Profile update info"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /me [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, 400, "invalid request")
		return
	}
	err := h.service.UpdateProfile(userID.(string), req.FullName)
	if err != nil {
		jsonError(c, 500, err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "profile updated"})
}

// UpdateUserRole godoc
// @Summary Update user role
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param body body models.UpdateRoleRequest true "Role update info"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /users/{id}/role [put]
func (h *UserHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonError(c, 400, "invalid request")
		return
	}
	if req.Role != "buyer" && req.Role != "partner" && req.Role != "admin" {
		jsonError(c, 400, "invalid role")
		return
	}
	user, err := h.service.GetProfile(userID)
	if err != nil {
		jsonError(c, 404, "user not found")
		return
	}
	user.Role = req.Role
	err = h.service.UpdateUser(user)
	if err != nil {
		jsonError(c, 500, err.Error())
		return
	}
	c.JSON(200, gin.H{"message": "role updated"})
}

func RegisterRoutes(r *gin.Engine, h *UserHandler) {
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/me", AuthMiddleware(), h.GetProfile)
	r.PUT("/me", AuthMiddleware(), h.UpdateProfile)
	r.PUT("/users/:id/role", AuthMiddleware(), AuthorizeRole("admin"), h.UpdateUserRole)
}
