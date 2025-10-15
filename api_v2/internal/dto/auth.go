package dto

import "time"

type RegisterRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}

type SessionResponse struct {
	Username    string `json:"username" validate:"required,min=3,max=30,alphanum"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8,max=100"`
	DisplayName string `json:"display_name" validate:"required,min=2,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ConfirmEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=100"`
}

type AuthResponse struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	User      UserProfile `json:"user"`
}

type UserProfile struct {
	ID                string    `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	DisplayName       string    `json:"display_name"`
	Bio               *string   `json:"bio"`
	ProfilePictureURL *string   `json:"profile_picture_url"`
	IsPrivate         bool      `json:"is_private"`
	EmailVerified     bool      `json:"email_verified"`
	Theme             string    `json:"theme"`
	CreatedAt         time.Time `json:"created_at"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type UpdateProfileRequest struct {
	DisplayName       *string `json:"display_name,omitempty" validate:"omitempty,min=2,max=100"`
	Bio               *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty" validate:"omitempty,url"`
	IsPrivate         *bool   `json:"is_private,omitempty"`
}

type UpdateSettingsRequest struct {
	Theme string `json:"theme" validate:"required,oneof=light dark"`
}

// Session DTOs
type UserSessionResponse struct {
	ID        string    `json:"id"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IsCurrent bool      `json:"is_current"`
}
