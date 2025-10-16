package domain

import "github.com/EduardoMG12/cine/api_v2/internal/dto"

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req *dto.LoginRequest, ipAddress, userAgent string) (*dto.AuthResponse, error)
	ConfirmEmail(req *dto.ConfirmEmailRequest) error
	ForgotPassword(req *dto.ForgotPasswordRequest) error
	ResetPassword(req *dto.ResetPasswordRequest) error
	Logout(sessionID string) error
	LogoutAllSessions(userID string) error
}
