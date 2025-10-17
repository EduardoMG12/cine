package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTService_GenerateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	userID := uuid.New()
	token, err := jwtService.GenerateToken(userID)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Fatal("Expected token to be non-empty")
	}
}

func TestJWTService_ValidateToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	userID := uuid.New()
	token, err := jwtService.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Expected no error when generating token, got %v", err)
	}

	claims, err := jwtService.ValidateToken(token)
	if err != nil {
		t.Fatalf("Expected no error when validating token, got %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, claims.UserID)
	}
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	_, err := jwtService.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error when validating invalid token")
	}
}

func TestJWTService_ValidateToken_WrongSecret(t *testing.T) {
	jwtService1 := NewJWTService("secret-1")
	jwtService2 := NewJWTService("secret-2")

	userID := uuid.New()
	token, err := jwtService1.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Expected no error when generating token, got %v", err)
	}

	_, err = jwtService2.ValidateToken(token)
	if err == nil {
		t.Fatal("Expected error when validating token with wrong secret")
	}
}

func TestJWTService_ParseToken(t *testing.T) {
	jwtService := NewJWTService("test-secret-key")

	userID := uuid.New()
	token, err := jwtService.GenerateToken(userID)
	if err != nil {
		t.Fatalf("Expected no error when generating token, got %v", err)
	}

	claims, err := jwtService.ParseToken(token)
	if err != nil {
		t.Fatalf("Expected no error when parsing token, got %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, claims.UserID)
	}

	if claims.ExpiresAt.Time.Before(time.Now().Add(23 * time.Hour)) {
		t.Fatal("Expected token to expire in about 24 hours")
	}
}
