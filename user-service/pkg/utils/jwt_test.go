package utils

import (
	"testing"
)

func TestGenerateAndParseJWT(t *testing.T) {
	userID := "test-user"
	role := "buyer"
	token, err := GenerateJWT(userID, role)
	if err != nil {
		t.Fatalf("GenerateJWT error: %v", err)
	}
	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("ParseJWT error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}
