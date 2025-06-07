package service

import (
	"errors"

	"github.com/luke/sfconnect-backend/user-service/internal/models"
	"github.com/luke/sfconnect-backend/user-service/internal/repository"
	"github.com/luke/sfconnect-backend/user-service/pkg/utils"
)

// UserService implementation

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(email, password, fullName string) (*models.User, error) {
	existing, _ := s.repo.GetByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Email:        email,
		PasswordHash: hash,
		FullName:     fullName,
		Role:         "buyer",
	}
	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *UserService) GetProfile(id string) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) UpdateProfile(id, fullName string) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.FullName = fullName
	return s.repo.Update(user)
}

// UpdateUser allows updating the full user object (admin use only)
func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}
