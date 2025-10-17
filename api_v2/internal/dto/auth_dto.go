package dto

import (
	"time"

	"github.com/google/uuid"
)

// LoginRequestDTO representa a requisição de login
type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequestDTO representa a requisição de registro
type RegisterRequestDTO struct {
	Username    string `json:"username" validate:"required,min=3,max=30"`
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8"`
}

// UserDTO representa um usuário na resposta da API
type UserDTO struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	DisplayName       string    `json:"display_name"`
	Bio               *string   `json:"bio,omitempty"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	IsPrivate         bool      `json:"is_private"`
	EmailVerified     bool      `json:"email_verified"`
	Theme             string    `json:"theme"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AuthResponseDTO representa a resposta de autenticação
type AuthResponseDTO struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

// APIResponse representa a estrutura padrão de resposta da API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError representa um erro na resposta da API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
