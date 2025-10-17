package handler

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
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

// Register godoc
// @Summary Register a new user
// @Description Register a new user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequestDTO true "Registration data"
// @Success 201 {object} dto.APIResponse{data=dto.AuthResponseDTO}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	domainReq := domain.RegisterRequest{
		Username:    req.Username,
		Email:       req.Email,
		DisplayName: req.DisplayName,
		Password:    req.Password,
	}

	authResp, err := h.authService.Register(domainReq)
	if err != nil {
		if err.Error() == "email already registered" {
			h.sendErrorResponse(w, http.StatusBadRequest, "EMAIL_EXISTS", "Email already registered")
			return
		}
		if err.Error() == "username already taken" {
			h.sendErrorResponse(w, http.StatusBadRequest, "USERNAME_EXISTS", "Username already taken")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "REGISTRATION_FAILED", "Failed to register user")
		return
	}

	respDTO := dto.AuthResponseDTO{
		Token: authResp.Token,
		User:  h.userToDTO(authResp.User),
	}

	h.sendSuccessResponse(w, http.StatusCreated, "User registered successfully", respDTO)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequestDTO true "Login credentials"
// @Success 200 {object} dto.APIResponse{data=dto.AuthResponseDTO}
// @Failure 400 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	domainReq := domain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	authResp, err := h.authService.Login(domainReq)
	if err != nil {
		if err.Error() == "invalid credentials" {
			h.sendErrorResponse(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "LOGIN_FAILED", "Failed to login")
		return
	}

	respDTO := dto.AuthResponseDTO{
		Token: authResp.Token,
		User:  h.userToDTO(authResp.User),
	}

	h.sendSuccessResponse(w, http.StatusOK, "Login successful", respDTO)
}

// Me godoc
// @Summary Get current user info
// @Description Get authenticated user information
// @Tags Authentication
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.APIResponse{data=dto.UserDTO}
// @Failure 401 {object} dto.APIResponse
// @Router /auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context")
		return
	}

	userDTO := h.userToDTO(*user)
	h.sendSuccessResponse(w, http.StatusOK, "User info retrieved", userDTO)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user by invalidating the current session
// @Tags Authentication
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header required")
		return
	}

	token := authHeader[len("Bearer "):]

	err := h.authService.Logout(token)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "LOGOUT_FAILED", "Failed to logout")
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Logout successful", nil)
}

// LogoutAll godoc
// @Summary Logout from all sessions
// @Description Logout user from all active sessions
// @Tags Authentication
// @Produce json
// @Security Bearer
// @Success 200 {object} dto.APIResponse
// @Failure 401 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutAll(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		h.sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not found in context")
		return
	}

	err := h.authService.LogoutAll(user.ID)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "LOGOUT_ALL_FAILED", "Failed to logout from all sessions")
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Logged out from all sessions", nil)
}

func (h *AuthHandler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/auth", func(r chi.Router) {
		// Public routes
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Get("/me", h.Me)
			r.Post("/logout", h.Logout)
			r.Post("/logout-all", h.LogoutAll)
		})
	})
}

// Helper methods
func (h *AuthHandler) userToDTO(user domain.User) dto.UserDTO {
	return dto.UserDTO{
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
		UpdatedAt:         user.UpdatedAt,
	}
}

func (h *AuthHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	response := dto.APIResponse{
		Success: false,
		Error: &dto.APIError{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
