package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	userService        domain.UserService
	userSessionService domain.UserSessionService
	jwtManager         *auth.JWTManager
	passwordHasher     *auth.PasswordHasher
	config             *config.Config
}

func NewAuthHandler(
	userService domain.UserService,
	userSessionService domain.UserSessionService,
	jwtManager *auth.JWTManager,
	passwordHasher *auth.PasswordHasher,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		userService:        userService,
		userSessionService: userSessionService,
		jwtManager:         jwtManager,
		passwordHasher:     passwordHasher,
		config:             config,
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
	lang := middleware.GetLanguageFromContext(r.Context())

	var req dto.RegisterRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_REQUEST_BODY", nil)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "VALIDATION_ERROR", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	// Hash the password
	hashedPassword, err := h.passwordHasher.GenerateHash(req.Password)
	if err != nil {
		slog.Error("Failed to hash password", "error", err)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
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
		Theme:         h.config.Application.DefaultTheme,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = h.userService.ValidateUser(user)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_DATA", nil)
		return
	}

	// Save user
	if err := h.userService.CreateUser(user); err != nil {
		slog.Error("Failed to create user", "error", err, "email", req.Email)
		utils.WriteErrorResponse(w, r.Context(), http.StatusConflict, "USER_ALREADY_EXISTS", nil)
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

	utils.WriteJSONResponse(w, r.Context(), http.StatusCreated, map[string]interface{}{
		"user":    userProfile,
		"message": i18n.T(lang, "REGISTRATION_SUCCESS"),
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
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_REQUEST_BODY", nil)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "VALIDATION_ERROR", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	// Get user by email
	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		slog.Warn("Login attempt with invalid email", "email", req.Email)
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "INVALID_CREDENTIALS", nil)
		return
	}

	// Verify password
	if err := h.passwordHasher.ComparePasswordAndHash(req.Password, user.PasswordHash); err != nil {
		slog.Warn("Login attempt with invalid password", "user_id", user.ID, "email", req.Email)
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "INVALID_CREDENTIALS", nil)
		return
	}

	// Create session
	clientIP := utils.GetClientIP(r)
	userAgent := r.UserAgent()

	session, err := h.userSessionService.CreateSession(user.ID, clientIP, userAgent)
	if err != nil {
		slog.Error("Failed to create session", "error", err, "user_id", user.ID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "SESSION_CREATE_FAILED", nil)
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, session.ID)
	if err != nil {
		slog.Error("Failed to generate token", "error", err, "user_id", user.ID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", nil)
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

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, authResponse)
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
	lang := middleware.GetLanguageFromContext(r.Context())

	var req dto.ConfirmEmailRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_REQUEST_BODY", nil)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "VALIDATION_ERROR", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	err := h.userService.ConfirmEmail(req.Token)
	if err != nil {
		slog.Warn("Email confirmation failed", "error", err)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_TOKEN", nil)
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "EMAIL_CONFIRMED"),
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
	lang := middleware.GetLanguageFromContext(r.Context())

	var req dto.ForgotPasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_REQUEST_BODY", nil)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "VALIDATION_ERROR", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	err := h.userService.RequestPasswordReset(req.Email)
	if err != nil {
		// Don't reveal whether email exists for security
		slog.Warn("Password reset request failed", "error", err, "email", req.Email)
	}

	// Always return success to prevent email enumeration
	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "PASSWORD_RESET_SENT"),
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
	lang := middleware.GetLanguageFromContext(r.Context())

	var req dto.ResetPasswordRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_REQUEST_BODY", nil)
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "VALIDATION_ERROR", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	err := h.userService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		slog.Warn("Password reset failed", "error", err)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_TOKEN", nil)
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "PASSWORD_RESET_SUCCESS"),
	})
}

// Helper method to get authenticated user from context
func (h *AuthHandler) GetAuthenticatedUser(r *http.Request) (*auth.Claims, bool) {
	return middleware.GetUserClaims(r.Context())
}
