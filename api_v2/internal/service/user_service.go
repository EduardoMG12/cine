package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(username, email, password, displayName string) (*domain.User, error) {
	if err := s.validateUsername(username); err != nil {
		return nil, err
	}

	if err := s.validateEmail(email); err != nil {
		return nil, err
	}

	if err := s.validatePassword(password); err != nil {
		return nil, err
	}

	if err := s.validateDisplayName(displayName); err != nil {
		return nil, err
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hasher := auth.NewPasswordHasher()
	hashedPassword, err := hasher.GenerateHash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user with unverified email
	user := &domain.User{
		Username:      username,
		Email:         email,
		PasswordHash:  hashedPassword,
		DisplayName:   displayName,
		EmailVerified: false,    // Set to false initially
		Theme:         "system", // Default theme
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate email verification token
	tokenString, err := generateSecureToken(32)
	if err != nil {
		slog.Warn("Failed to generate email verification token", "error", err)
		// Continue without email verification for now
		return user, nil
	}

	// Create email verification token with 24 hour expiration
	verificationToken := &domain.EmailVerificationToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.CreateEmailVerificationToken(verificationToken); err != nil {
		slog.Warn("Failed to create email verification token", "error", err)
		// Continue without email verification for now
	} else {
		slog.Info("Email verification token created", "user_id", user.ID, "token", tokenString)
		// TODO: Send confirmation email when email service is integrated
	}

	return user, nil
}

func (s *userService) GetUser(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if username, ok := updates["username"].(string); ok {
		if err := s.validateUsername(username); err != nil {
			return nil, err
		}
		user.Username = username
	}

	if email, ok := updates["email"].(string); ok {
		if err := s.validateEmail(email); err != nil {
			return nil, err
		}
		user.Email = email
	}

	if displayName, ok := updates["display_name"].(string); ok {
		if err := s.validateDisplayName(displayName); err != nil {
			return nil, err
		}
		user.DisplayName = displayName
	}

	if bio, ok := updates["bio"].(string); ok {
		if len(bio) > 500 {
			return nil, errors.New("bio must be less than 500 characters")
		}
		user.Bio = &bio
	}

	if avatarURL, ok := updates["profile_picture_url"].(string); ok {
		user.ProfilePictureURL = &avatarURL
	}

	if err := s.ValidateUser(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) ValidateUser(user *domain.User) error {
	if err := s.validateUsername(user.Username); err != nil {
		return err
	}

	if err := s.validateEmail(user.Email); err != nil {
		return err
	}

	if err := s.validateDisplayName(user.DisplayName); err != nil {
		return err
	}

	if user.Bio != nil && len(*user.Bio) > 500 {
		return errors.New("bio must be less than 500 characters")
	}

	return nil
}

func (s *userService) validateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers and underscores")
	}

	return nil
}

func (s *userService) validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (s *userService) validateDisplayName(displayName string) error {
	displayName = strings.TrimSpace(displayName)
	if len(displayName) < 1 || len(displayName) > 100 {
		return errors.New("display name must be between 1 and 100 characters")
	}

	return nil
}

func (s *userService) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	// Check for at least one number, one uppercase, one lowercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	if !hasUpper || !hasLower || !hasNumber {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, and one number")
	}

	return nil
}

// Login validates credentials and returns user with JWT token
func (s *userService) Login(email, password string) (*domain.User, string, error) {
	// TODO: Implement login logic with password validation and JWT generation
	return nil, "", errors.New("not implemented")
}

// ConfirmEmail confirms user email with token
func (s *userService) ConfirmEmail(token string) error {
	// Get the verification token from database
	verificationToken, err := s.userRepo.GetEmailVerificationToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	// Check if token has expired
	if time.Now().After(verificationToken.ExpiresAt) {
		// Clean up expired token
		s.userRepo.DeleteEmailVerificationToken(token)
		return errors.New("verification token has expired")
	}

	// Mark email as verified
	if err := s.userRepo.MarkEmailAsVerified(verificationToken.UserID); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	// Delete the used token
	if err := s.userRepo.DeleteEmailVerificationToken(token); err != nil {
		slog.Warn("Failed to delete used verification token", "error", err)
	}

	return nil
}

// RequestPasswordReset initiates password reset flow
func (s *userService) RequestPasswordReset(email string) error {
	// Check if user exists
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil
	}

	// Generate reset token
	tokenString, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Create password reset token with 1 hour expiration
	resetToken := &domain.PasswordResetToken{
		UserID:    user.ID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := s.userRepo.CreatePasswordResetToken(resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// TODO: Send password reset email (will be implemented when email service is integrated)
	slog.Info("Password reset token created for user", "email", email, "token", tokenString)

	return nil
}

// ResetPassword resets password with token
func (s *userService) ResetPassword(token, newPassword string) error {
	// Get the reset token from database
	resetToken, err := s.userRepo.GetPasswordResetToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	// Check if token has expired
	if time.Now().After(resetToken.ExpiresAt) {
		// Clean up expired token
		s.userRepo.DeletePasswordResetToken(token)
		return errors.New("reset token has expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Hash new password
	hasher := auth.NewPasswordHasher()
	hashedPassword, err := hasher.GenerateHash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user.PasswordHash = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete the used reset token
	if err := s.userRepo.DeletePasswordResetToken(token); err != nil {
		slog.Warn("Failed to delete used reset token", "error", err)
	}

	return nil
}

// GetUserProfile returns user profile respecting privacy settings
func (s *userService) GetUserProfile(id int) (*domain.User, error) {
	// TODO: Implement profile retrieval with privacy checks
	return s.GetUser(id)
}

// UpdateProfile updates user profile
func (s *userService) UpdateProfile(userID int, updates map[string]interface{}) (*domain.User, error) {
	return s.UpdateUser(userID, updates)
}

// UpdateSettings updates user settings
func (s *userService) UpdateSettings(userID int, settings map[string]interface{}) error {
	return s.userRepo.UpdateSettings(userID, settings)
}

// GetUserByEmail retrieves user by email
func (s *userService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(email)
}

// CreateUser creates a new user entity
func (s *userService) CreateUser(user *domain.User) error {
	return s.userRepo.Create(user)
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
