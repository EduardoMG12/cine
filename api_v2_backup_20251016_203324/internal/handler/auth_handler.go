package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/auth"
	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/service"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	authService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Deprecated: Use NewAuthHandler instead
func NewAuthHandlerWithDependencies(
	userService domain.UserService,
	userSessionService domain.UserSessionService,
	jwtManager *auth.JWTManager,
	passwordHasher *auth.PasswordHasher,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		authService: nil, // This will cause errors, use NewAuthHandler instead
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
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", "")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf(`{"details": "%s"}`, err.Error()))
		return
	}

	if len(req.Password) < 6 {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "PASSWORD_TOO_SHORT", "")
		return
	}

	authResponse, err := h.authService.Register(&req)
	if err != nil {
		slog.Error("Failed to register user", "error", err, "email", req.Email)

		// Handler responsibility: map service errors to appropriate HTTP status codes
		switch {
		case errors.Is(err, service.ErrUserAlreadyExists):
			utils.WriteErrorResponse(w, r, http.StatusConflict, "USER_ALREADY_EXISTS", "")
		case strings.Contains(err.Error(), "validation failed"):
			utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_USER_DATA", err.Error())
		default:
			utils.WriteErrorResponse(w, r, http.StatusInternalServerError, "REGISTRATION_FAILED", "")
		}
		return
	}

	slog.Info("User registered successfully", "user_id", authResponse.User.ID, "email", authResponse.User.Email)

	// Handler responsibility: format response properly
	utils.WriteJSONResponse(w, r, http.StatusCreated, map[string]interface{}{
		"token":      authResponse.Token,
		"expires_at": authResponse.ExpiresAt,
		"user":       authResponse.User,
		"message":    i18n.T(lang, "REGISTRATION_SUCCESS"),
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
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", "")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf(`{"details": "%s"}`, err.Error()))
		return
	}

	clientIP := utils.GetClientIP(r)
	userAgent := r.UserAgent()

	if len(strings.TrimSpace(req.Email)) == 0 || len(strings.TrimSpace(req.Password)) == 0 {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "EMPTY_CREDENTIALS", "")
		return
	}

	authResponse, err := h.authService.Login(&req, clientIP, userAgent)
	if err != nil {
		slog.Warn("Login attempt failed", "error", err.Error(), "email", req.Email, "ip", clientIP)

		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			utils.WriteErrorResponse(w, r, http.StatusUnauthorized, "INVALID_CREDENTIALS", "")
		case errors.Is(err, service.ErrEmailNotVerified):
			utils.WriteErrorResponse(w, r, http.StatusUnauthorized, "EMAIL_NOT_VERIFIED", "")
		default:
			utils.WriteErrorResponse(w, r, http.StatusInternalServerError, "LOGIN_FAILED", "")
		}
		return
	}

	slog.Info("User logged in successfully", "user_id", authResponse.User.ID, "ip", clientIP)

	utils.WriteJSONResponse(w, r, http.StatusOK, authResponse)
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
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", "")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf(`{"details": "%s"}`, err.Error()))
		return
	}

	err := h.authService.ConfirmEmail(&req)
	if err != nil {
		slog.Warn("Email confirmation failed", "error", err)

		switch err {
		case service.ErrTokenNotFound, service.ErrTokenExpired:
			utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_TOKEN", "")
		default:
			utils.WriteErrorResponse(w, r, http.StatusInternalServerError, "CONFIRMATION_FAILED", "")
		}
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, dto.MessageResponse{
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
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", "")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf(`{"details": "%s"}`, err.Error()))
		return
	}

	err := h.authService.ForgotPassword(&req)
	if err != nil {
		// Don't reveal whether email exists for security
		slog.Warn("Password reset request failed", "error", err, "email", req.Email)
	}

	// Always return success to prevent email enumeration
	utils.WriteJSONResponse(w, r, http.StatusOK, dto.MessageResponse{
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
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_REQUEST_BODY", "")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteErrorResponse(w, r, http.StatusBadRequest, "VALIDATION_ERROR", fmt.Sprintf(`{"details": "%s"}`, err.Error()))
		return
	}

	err := h.authService.ResetPassword(&req)
	if err != nil {
		slog.Warn("Password reset failed", "error", err)

		switch err {
		case service.ErrTokenNotFound, service.ErrTokenExpired:
			utils.WriteErrorResponse(w, r, http.StatusBadRequest, "INVALID_TOKEN", "")
		default:
			utils.WriteErrorResponse(w, r, http.StatusInternalServerError, "RESET_FAILED", "")
		}
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "PASSWORD_RESET_SUCCESS"),
	})
}

// Helper method to get authenticated user from context
func (h *AuthHandler) GetAuthenticatedUser(r *http.Request) (*auth.Claims, bool) {
	return middleware.GetUserClaims(r.Context())
}
