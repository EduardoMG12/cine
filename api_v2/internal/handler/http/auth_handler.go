package http

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/auth"
)

type AuthHandler struct {
	registerUC  *auth.RegisterUseCase
	loginUC     *auth.LoginUseCase
	getMeUC     *auth.GetMeUseCase
	logoutUC    *auth.LogoutUseCase
	logoutAllUC *auth.LogoutAllUseCase
}

func NewAuthHandler(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
	getMeUC *auth.GetMeUseCase,
	logoutUC *auth.LogoutUseCase,
	logoutAllUC *auth.LogoutAllUseCase,
) *AuthHandler {
	return &AuthHandler{
		registerUC:  registerUC,
		loginUC:     loginUC,
		getMeUC:     getMeUC,
		logoutUC:    logoutUC,
		logoutAllUC: logoutAllUC,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequestDTO true "Registration data"
// @Success 201 {object} dto.APIResponse{data=dto.AuthResponseDTO}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.registerUC.Execute(req)
	if err != nil {
		if err.Error() == "email already registered" {
			sendErrorResponse(w, http.StatusBadRequest, "EMAIL_EXISTS", "Email already registered")
			return
		}
		if err.Error() == "username already taken" {
			sendErrorResponse(w, http.StatusBadRequest, "USERNAME_EXISTS", "Username already taken")
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register user")
		return
	}

	sendSuccessResponse(w, http.StatusCreated, "User registered successfully", result)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequestDTO true "Login credentials"
// @Success 200 {object} dto.APIResponse{data=dto.AuthResponseDTO}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.loginUC.Execute(req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			sendErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, "LOGIN_FAILED", "Failed to login")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Login successful", result)
}

// GetMe godoc
// @Summary Get current user
// @Description Get authenticated user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.UserDTO}
// @Failure 401 {object} dto.APIResponse
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context")
		return
	}

	result, err := h.getMeUC.Execute(userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "FAILED", "Failed to get user info")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "User info retrieved", result)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user by invalidating current session
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header required")
		return
	}

	token := authHeader[len("Bearer "):]
	if err := h.logoutUC.Execute(token); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Logout successful", nil)
}

// LogoutAll godoc
// @Summary Logout from all sessions
// @Description Logout user from all active sessions
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/auth/logout-all [post]
func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context")
		return
	}

	if err := h.logoutAllUC.Execute(userID); err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "LOGOUT_ALL_FAILED", "Failed to logout from all sessions")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Logged out from all sessions", nil)
}
