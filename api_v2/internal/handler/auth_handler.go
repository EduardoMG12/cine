package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	userService        domain.UserService
	userSessionService domain.UserSessionService
	jwtManager         *auth.JWTManager
	passwordHasher     *auth.PasswordHasher
}

func NewAuthHandler(
	userService domain.UserService,
	userSessionService domain.UserSessionService,
	jwtManager *auth.JWTManager,
	passwordHasher *auth.PasswordHasher,
) *AuthHandler {
	return &AuthHandler{
		userService:        userService,
		userSessionService: userSessionService,
		jwtManager:         jwtManager,
		passwordHasher:     passwordHasher,
	}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Public routes
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Post("/confirm-email", h.ConfirmEmail)
	r.Post("/forgot-password", h.ForgotPassword)
	r.Post("/reset-password", h.ResetPassword)

	return r
}

// Register creates a new user account
// @Summary Register a new user
// @Description Create a new user account with email verification
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse "User created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 409 {object} utils.ErrorResponse "User already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	// Hash the password
	hashedPassword, err := h.passwordHasher.GenerateHash(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		utils.WriteError(w, http.StatusInternalServerError, err, "Internal server error")
		return
	}

	// Create user
	user := &domain.User{
		Username:      req.Username,
		Email:         req.Email,
		PasswordHash:  hashedPassword,
		DisplayName:   req.DisplayName,
		IsPrivate:     false,
		EmailVerified: false, // Will be verified via email
		Theme:         "light",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = h.userService.ValidateUser(user)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid user data")
		return
	}

	// Save user
	if err := h.userService.CreateUser(user); err != nil {
		slog.Error("Failed to create user", "error", err, "email", req.Email)
		utils.WriteError(w, http.StatusConflict, err, "Failed to create user")
		return
	}

	// TODO: Send confirmation email (implement email service)
	slog.Info("User registered successfully", "user_id", user.ID, "email", user.Email)

	// Return user profile (without sensitive data)
	userProfile := dto.UserProfile{
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
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"user":    userProfile,
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

// Login authenticates a user and returns JWT token
// @Summary User login
// @Description Authenticate user credentials and return access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse "Login successful"
// @Failure 400 {object} utils.ErrorResponse "Invalid credentials"
// @Failure 401 {object} utils.ErrorResponse "Authentication failed"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		slog.Warn("Login attempt with invalid email", "email", req.Email)
		utils.WriteError(w, http.StatusUnauthorized, err, "Invalid credentials")
		return
	}

	// Verify password
	if err := h.passwordHasher.ComparePasswordAndHash(req.Password, user.PasswordHash); err != nil {
		slog.Warn("Login attempt with invalid password", "user_id", user.ID, "email", req.Email)
		utils.WriteError(w, http.StatusUnauthorized, err, "Invalid credentials")
		return
	}

	// Create session
	clientIP := utils.GetClientIP(r)
	userAgent := r.UserAgent()

	session, err := h.userSessionService.CreateSession(user.ID, clientIP, userAgent)
	if err != nil {
		slog.Error("Failed to create session", "error", err, "user_id", user.ID)
		utils.WriteError(w, http.StatusInternalServerError, err, "Failed to create session")
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, session.ID)
	if err != nil {
		slog.Error("Failed to generate token", "error", err, "user_id", user.ID)
		utils.WriteError(w, http.StatusInternalServerError, err, "Failed to generate token")
		return
	}

	slog.Info("User logged in successfully", "user_id", user.ID, "session_id", session.ID)

	// Return auth response
	userProfile := dto.UserProfile{
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
	}

	authResponse := dto.AuthResponse{
		Token:     token,
		ExpiresAt: session.ExpiresAt,
		User:      userProfile,
	}

	utils.WriteJSON(w, http.StatusOK, authResponse)
}

// ConfirmEmail verifies user's email address
// @Summary Confirm email address
// @Description Verify user's email address using confirmation token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ConfirmEmailRequest true "Email confirmation token"
// @Success 200 {object} dto.MessageResponse "Email confirmed successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/confirm-email [post]
func (h *AuthHandler) ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	var req dto.ConfirmEmailRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	err := h.userService.ConfirmEmail(req.Token)
	if err != nil {
		slog.Warn("Email confirmation failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid or expired confirmation token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Message: "Email confirmed successfully",
	})
}

// ForgotPassword initiates password reset process
// @Summary Request password reset
// @Description Send password reset link to user's email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordRequest true "User email for password reset"
// @Success 200 {object} dto.MessageResponse "Password reset email sent"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	err := h.userService.RequestPasswordReset(req.Email)
	if err != nil {
		// Don't reveal whether email exists for security
		slog.Warn("Password reset request failed", "error", err, "email", req.Email)
	}

	// Always return success to prevent email enumeration
	utils.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Message: "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword resets user password with token
// @Summary Reset user password
// @Description Reset user password using reset token from email
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Password reset token and new password"
// @Success 200 {object} dto.MessageResponse "Password reset successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid or expired token"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	err := h.userService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		slog.Warn("Password reset failed", "error", err)
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid or expired reset token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Message: "Password reset successfully",
	})
}

// Helper method to get authenticated user from context
func (h *AuthHandler) GetAuthenticatedUser(r *http.Request) (*auth.Claims, bool) {
	return middleware.GetUserClaims(r.Context())
}
