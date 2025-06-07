package models

// User represents the structure of an individual user in the system
type User struct {
	ID          string `db:"id" json:"id"`
	Email       string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	FullName    string `db:"full_name" json:"full_name"`
	Role        string `db:"role" json:"role"`
	CreatedAt   string `db:"created_at" json:"created_at"`
	UpdatedAt   string `db:"updated_at" json:"updated_at"`
}

// Registration request for Swagger
// @Description Registration request body
// @name RegisterRequest
// @Param email body string true "Email"
// @Param password body string true "Password"
// @Param full_name body string true "Full Name"
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// Login request for Swagger
// @Description Login request body
// @name LoginRequest
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Profile update request for Swagger
// @Description Profile update request body
// @name UpdateProfileRequest
type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
}

// Role update request for Swagger
// @Description Role update request body
// @name UpdateRoleRequest
type UpdateRoleRequest struct {
	Role string `json:"role"`
}
