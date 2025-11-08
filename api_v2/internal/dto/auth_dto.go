package dto

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequestDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequestDTO struct {
	Username    string `json:"username" validate:"required,min=3,max=30"`
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8"`
}

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

type AuthResponseDTO struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
