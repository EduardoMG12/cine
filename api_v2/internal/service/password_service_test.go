package service

import (
	"testing"
)

func TestPasswordService_HashPassword(t *testing.T) {
	passwordService := NewPasswordService()

	password := "testpassword123"
	hash, err := passwordService.HashPassword(password)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatal("Expected hash to be non-empty")
	}

	if hash == password {
		t.Fatal("Expected hash to be different from password")
	}
}

func TestPasswordService_ComparePassword(t *testing.T) {
	passwordService := NewPasswordService()

	password := "testpassword123"
	wrongPassword := "wrongpassword123"

	hash, err := passwordService.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error when hashing, got %v", err)
	}

	if !passwordService.ComparePassword(hash, password) {
		t.Fatal("Expected password comparison to return true for correct password")
	}

	if passwordService.ComparePassword(hash, wrongPassword) {
		t.Fatal("Expected password comparison to return false for wrong password")
	}
}

func TestPasswordService_ComparePassword_EmptyHash(t *testing.T) {
	passwordService := NewPasswordService()

	password := "testpassword123"

	if passwordService.ComparePassword("", password) {
		t.Fatal("Expected password comparison to return false for empty hash")
	}
}

func TestPasswordService_ComparePassword_EmptyPassword(t *testing.T) {
	passwordService := NewPasswordService()

	hash, err := passwordService.HashPassword("testpassword123")
	if err != nil {
		t.Fatalf("Expected no error when hashing, got %v", err)
	}

	if passwordService.ComparePassword(hash, "") {
		t.Fatal("Expected password comparison to return false for empty password")
	}
}
