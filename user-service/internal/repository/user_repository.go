package repository

import (
	"database/sql"

	"github.com/luke/sfconnect-backend/user-service/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	Update(user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password_hash, full_name, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(query, user.Email, user.PasswordHash, user.FullName, user.Role).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, full_name, role, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password_hash, full_name, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *models.User) error {
	query := `UPDATE users SET full_name = $1, role = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at`
	return r.db.QueryRow(query, user.FullName, user.Role, user.ID).Scan(&user.UpdatedAt)
}
