package utils

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mysecret"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if !CheckPasswordHash(password, hash) {
		t.Error("CheckPasswordHash failed: password should match hash")
	}
	if CheckPasswordHash("wrong", hash) {
		t.Error("CheckPasswordHash failed: wrong password should not match hash")
	}
}
