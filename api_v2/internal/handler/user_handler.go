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
	userService domain.UserService
	validator   *validator.Validate
}

type CreateUserRequest struct {
	Username    string `json:"username" validate:"required,min=3,max=30"`
	Email       string `json:"email" validate:"required,email"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	DisplayName *string `json:"display_name,omitempty" validate:"omitempty,min=1,max=100"`
	Bio         *string `json:"bio,omitempty" validate:"omitempty,max=500"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *UserHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateUser)
	r.Get("/{id}", h.GetUser)
	r.Put("/{id}", h.UpdateUser)
	r.Delete("/{id}", h.DeleteUser)

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

	user, err := h.userService.CreateUser(req.Username, req.Email, req.DisplayName)
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
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	user, err := h.userService.UpdateUser(id, updates)
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
