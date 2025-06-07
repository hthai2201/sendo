package service

import (
	"errors"
	"testing"

	"github.com/luke/sfconnect-backend/user-service/internal/models"
)

type mockRepo struct {
	users map[string]*models.User
}

func (m *mockRepo) Create(user *models.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("duplicate")
	}
	user.ID = "mock-id"
	m.users[user.Email] = user
	return nil
}
func (m *mockRepo) GetByEmail(email string) (*models.User, error) {
	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}
func (m *mockRepo) GetByID(id string) (*models.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}
func (m *mockRepo) Update(user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func newMockRepo() *mockRepo {
	return &mockRepo{users: make(map[string]*models.User)}
}

func TestRegisterAndLogin(t *testing.T) {
	repo := newMockRepo()
	svc := NewUserService(repo)
	// Register
	user, err := svc.Register("a@b.com", "pw", "Test User")
	if err != nil {
		t.Fatalf("Register error: %v", err)
	}
	if user.Email != "a@b.com" {
		t.Errorf("expected email a@b.com, got %s", user.Email)
	}
	// Duplicate
	_, err = svc.Register("a@b.com", "pw", "Test User")
	if err == nil {
		t.Error("expected error for duplicate email")
	}
	// Login
	token, err := svc.Login("a@b.com", "pw")
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	if token == "" {
		t.Error("expected JWT token, got empty string")
	}
	// Wrong password
	_, err = svc.Login("a@b.com", "wrong")
	if err == nil {
		t.Error("expected error for wrong password")
	}
}
