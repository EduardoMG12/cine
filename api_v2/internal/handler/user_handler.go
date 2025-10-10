package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	userService        domain.UserService
	userSessionService domain.UserSessionService
	validator          *validator.Validate
}

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=30"`
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	Username          *string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email             *string `json:"email,omitempty" validate:"omitempty,email"`
	DisplayName       *string `json:"display_name,omitempty" validate:"omitempty,min=1,max=100"`
	Bio               *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	ProfilePictureURL *string `json:"profile_picture_url,omitempty"`
}

type UserSettingsRequest struct {
	Theme              *string `json:"theme,omitempty" validate:"omitempty,oneof=light dark"`
	EmailNotifications *bool   `json:"email_notifications,omitempty"`
	PushNotifications  *bool   `json:"push_notifications,omitempty"`
	PrivateProfile     *bool   `json:"private_profile,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func NewUserHandler(userService domain.UserService, userSessionService domain.UserSessionService) *UserHandler {
	return &UserHandler{
		userService:        userService,
		userSessionService: userSessionService,
		validator:          validator.New(),
	}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateUser)
	r.Get("/{id}", h.GetUser)
	r.Put("/{id}", h.UpdateUser)
	r.Delete("/{id}", h.DeleteUser)

	// Session management endpoints
	r.Get("/me/sessions", h.GetMySessions)
	r.Delete("/me/sessions/{sessionId}", h.RevokeSession)
	r.Delete("/me/sessions", h.RevokeAllSessions)

	// User settings endpoint
	r.Put("/me/settings", h.UpdateUserSettings)

	return r
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// TODO: This is temporary - will be replaced with proper auth endpoints
	user, err := h.userService.Register(req.Username, req.Email, "temp_password", req.DisplayName)
	if err != nil {
		slog.Error("Failed to create user", "error", err, "username", req.Username)
		h.writeErrorResponse(w, http.StatusConflict, "Failed to create user", err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, user)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "ID must be a number")
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		slog.Error("Failed to get user", "error", err, "id", id)
		h.writeErrorResponse(w, http.StatusNotFound, "User not found", err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "ID must be a number")
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.ProfilePictureURL != nil {
		updates["profile_picture_url"] = *req.ProfilePictureURL
	}

	user, err := h.userService.UpdateProfile(id, updates)
	if err != nil {
		slog.Error("Failed to update user", "error", err, "id", id)
		h.writeErrorResponse(w, http.StatusBadRequest, "Failed to update user", err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "ID must be a number")
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		slog.Error("Failed to delete user", "error", err, "id", id)
		h.writeErrorResponse(w, http.StatusNotFound, "Failed to delete user", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
	}
}

func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, error, message string) {
	response := ErrorResponse{
		Error:   error,
		Message: message,
	}

	h.writeJSONResponse(w, statusCode, response)
}

// GetMySessions returns all active sessions for the authenticated user
// @Summary Get user sessions
// @Description Retrieve all active sessions for the authenticated user
// @Tags users,sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} domain.UserSession
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me/sessions [get]
func (h *UserHandler) GetMySessions(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from JWT token in middleware
	// For now, using a mock user ID
	userID := 1 // This should come from the JWT middleware context

	sessions, err := h.userSessionService.GetUserSessions(userID)
	if err != nil {
		slog.Error("Failed to get user sessions", "error", err, "userID", userID)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve sessions", err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusOK, sessions)
}

// RevokeSession revokes a specific session for the authenticated user
// @Summary Revoke specific session
// @Description Revoke a specific session by session ID
// @Tags users,sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param sessionId path int true "Session ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me/sessions/{sessionId} [delete]
func (h *UserHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from JWT token in middleware
	userID := 1 // This should come from the JWT middleware context

	sessionIDStr := chi.URLParam(r, "sessionId")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid session ID", "Session ID must be a number")
		return
	}

	err = h.userSessionService.RevokeSession(userID, sessionID)
	if err != nil {
		slog.Error("Failed to revoke session", "error", err, "userID", userID, "sessionID", sessionID)
		h.writeErrorResponse(w, http.StatusNotFound, "Session not found", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RevokeAllSessions revokes all sessions for the authenticated user (logout everywhere)
// @Summary Revoke all sessions
// @Description Revoke all sessions for the authenticated user (logout from all devices)
// @Tags users,sessions
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me/sessions [delete]
func (h *UserHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from JWT token in middleware
	userID := 1 // This should come from the JWT middleware context

	err := h.userSessionService.RevokeAllSessions(userID)
	if err != nil {
		slog.Error("Failed to revoke all sessions", "error", err, "userID", userID)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to revoke sessions", err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserSettings updates user preferences and settings
// @Summary Update user settings
// @Description Update user preferences like theme, notifications, privacy settings
// @Tags users,settings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param settings body UserSettingsRequest true "User settings to update"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/me/settings [put]
func (h *UserHandler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	// TODO: Get user ID from JWT token in middleware
	userID := 1 // This should come from the JWT middleware context

	var req UserSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// Build the updates map for settings
	updates := make(map[string]interface{})
	if req.Theme != nil {
		updates["theme"] = *req.Theme
	}
	if req.EmailNotifications != nil {
		updates["email_notifications"] = *req.EmailNotifications
	}
	if req.PushNotifications != nil {
		updates["push_notifications"] = *req.PushNotifications
	}
	if req.PrivateProfile != nil {
		updates["private_profile"] = *req.PrivateProfile
	}

	user, err := h.userService.UpdateProfile(userID, updates)
	if err != nil {
		slog.Error("Failed to update user settings", "error", err, "userID", userID)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to update settings", err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}
