package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailNotVerified   = errors.New("email not verified")
	ErrTokenNotFound      = errors.New("token not found")
	ErrTokenExpired       = errors.New("token expired")
)

type AuthService struct {
	userRepo       domain.UserRepository
	sessionRepo    domain.UserSessionRepository
	jwtManager     *auth.JWTManager
	passwordHasher *auth.PasswordHasher
	cfg            *config.Config
}

func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.UserSessionRepository,
	jwtManager *auth.JWTManager,
	passwordHasher *auth.PasswordHasher,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		jwtManager:     jwtManager,
		passwordHasher: passwordHasher,
		cfg:            cfg,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	existingUser, _ = s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.passwordHasher.GenerateHash(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		DisplayName:   req.DisplayName,
		IsPrivate:     false,
		EmailVerified: false,
		Theme:         "light",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate email verification token
	token, err := auth.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	emailToken := &domain.EmailVerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreateEmailVerificationToken(emailToken); err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	// TODO: Send verification email

	return s.createAuthResponse(user, "", "")
}

func (s *AuthService) Login(req *dto.LoginRequest, ipAddress, userAgent string) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := s.passwordHasher.ComparePasswordAndHash(req.Password, user.PasswordHash); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if email is verified (optional based on requirements)
	// if !user.EmailVerified {
	//     return nil, ErrEmailNotVerified
	// }

	return s.createAuthResponse(user, ipAddress, userAgent)
}

func (s *AuthService) ConfirmEmail(req *dto.ConfirmEmailRequest) error {
	token, err := s.userRepo.GetEmailVerificationToken(req.Token)
	if err != nil {
		return ErrTokenNotFound
	}

	if time.Now().After(token.ExpiresAt) {
		return ErrTokenExpired
	}

	// Mark email as verified
	if err := s.userRepo.MarkEmailAsVerified(token.UserID); err != nil {
		return fmt.Errorf("failed to mark email as verified: %w", err)
	}

	// Delete the verification token
	if err := s.userRepo.DeleteEmailVerificationToken(req.Token); err != nil {
		return fmt.Errorf("failed to delete verification token: %w", err)
	}

	return nil
}

func (s *AuthService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists
		return nil
	}

	// Generate password reset token
	token, err := auth.GenerateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour expiry
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreatePasswordResetToken(resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send reset email

	return nil
}

func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	token, err := s.userRepo.GetPasswordResetToken(req.Token)
	if err != nil {
		return ErrTokenNotFound
	}

	if time.Now().After(token.ExpiresAt) {
		return ErrTokenExpired
	}

	// Hash new password
	hashedPassword, err := s.passwordHasher.GenerateHash(req.NewPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user, err := s.userRepo.GetByID(token.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete the reset token
	if err := s.userRepo.DeletePasswordResetToken(req.Token); err != nil {
		return fmt.Errorf("failed to delete reset token: %w", err)
	}

	// Invalidate all user sessions
	if err := s.sessionRepo.DeleteByUserID(user.ID); err != nil {
		return fmt.Errorf("failed to invalidate sessions: %w", err)
	}

	return nil
}

func (s *AuthService) Logout(sessionID int) error {
	return s.sessionRepo.DeleteByID(sessionID)
}

func (s *AuthService) LogoutAllSessions(userID int) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

func (s *AuthService) createAuthResponse(user *domain.User, ipAddress, userAgent string) (*dto.AuthResponse, error) {
	// Create session
	session := &domain.UserSession{
		UserID:    user.ID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.JWT.Expiration) * time.Hour),
	}

	// Generate JWT token
	tokenString, err := s.jwtManager.GenerateToken(user.ID, user.Email, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session.Token = tokenString

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &dto.AuthResponse{
		Token:     tokenString,
		ExpiresAt: session.ExpiresAt,
		User: dto.UserProfile{
			ID:                user.ID,
			Username:          user.Username,
			Email:             user.Email,
			DisplayName:       user.DisplayName,
			Bio:               user.Bio,
			ProfilePictureURL: user.ProfilePictureURL,
			IsPrivate:         user.IsPrivate,
			EmailVerified:     user.EmailVerified,
			Theme:             user.Theme,
			CreatedAt:         user.CreatedAt,
		},
	}, nil
}
