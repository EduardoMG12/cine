package domain

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=30"`
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	jwt.RegisteredClaims
}

type AuthService interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
	ValidateToken(token string) (*User, error)
	Logout(token string) error
	LogoutAll(userID uuid.UUID) error
}

type JWTService interface {
	GenerateToken(userID uuid.UUID) (string, error)
	ValidateToken(token string) (*JWTClaims, error)
	ParseToken(token string) (*JWTClaims, error)
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, password string) bool
}
