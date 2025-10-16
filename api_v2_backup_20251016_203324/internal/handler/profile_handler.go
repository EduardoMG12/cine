package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProfileHandler struct {
	userService        domain.UserService
	userSessionService domain.UserSessionService
}

func NewProfileHandler(
	userService domain.UserService,
	userSessionService domain.UserSessionService,
) *ProfileHandler {
	return &ProfileHandler{
		userService:        userService,
		userSessionService: userSessionService,
	}
}

func (h *ProfileHandler) Routes(authMiddleware *middleware.AuthMiddleware) chi.Router {
	r := chi.NewRouter()

	// Protected routes
	r.Use(authMiddleware.RequireAuth)

	r.Get("/me", h.GetMyProfile)
	r.Put("/me", h.UpdateMyProfile)
	r.Put("/me/settings", h.UpdateMySettings)

	// Session management
	r.Get("/me/sessions", h.GetMySessions)
	r.Delete("/me/sessions/{sessionId}", h.RevokeSession)
	r.Delete("/me/sessions", h.RevokeAllSessions)

	// Public routes (with optional auth)
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.OptionalAuth)
		r.Get("/{userId}", h.GetUserProfile)
	})

	return r
}

func (h *ProfileHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	user, err := h.userService.GetUser(claims.UserID)
	if err != nil {
		slog.Error("Failed to get user profile", "error", err, "user_id", claims.UserID)
		utils.WriteError(w, http.StatusNotFound, err, "User not found")
		return
	}

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

	utils.WriteJSON(w, http.StatusOK, userProfile)
}

func (h *ProfileHandler) UpdateMyProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	var req dto.UpdateProfileRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.ProfilePictureURL != nil {
		updates["profile_picture_url"] = *req.ProfilePictureURL
	}
	if req.IsPrivate != nil {
		updates["is_private"] = *req.IsPrivate
	}

	user, err := h.userService.UpdateProfile(claims.UserID, updates)
	if err != nil {
		slog.Error("Failed to update profile", "error", err, "user_id", claims.UserID)
		utils.WriteError(w, http.StatusBadRequest, err, "Failed to update profile")
		return
	}

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

	utils.WriteJSON(w, http.StatusOK, userProfile)
}

func (h *ProfileHandler) UpdateMySettings(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	var req dto.UpdateSettingsRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	settings := map[string]interface{}{
		"theme": req.Theme,
	}

	err := h.userService.UpdateSettings(claims.UserID, settings)
	if err != nil {
		slog.Error("Failed to update settings", "error", err, "user_id", claims.UserID)
		utils.WriteError(w, http.StatusBadRequest, err, "Failed to update settings")
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Message: "Settings updated successfully",
	})
}

func (h *ProfileHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		utils.WriteError(w, http.StatusBadRequest, nil, "Invalid user ID")
		return
	}

	// Get requesting user info (if authenticated)
	var requestingUserID string
	if claims, ok := middleware.GetUserClaims(r.Context()); ok {
		requestingUserID = claims.UserID
	}

	user, err := h.userService.GetUserProfile(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err, "User not found")
		return
	}

	// Check privacy settings
	if user.IsPrivate && requestingUserID != userIDStr {
		// TODO: Check if users are friends/followers
		// For now, just return basic info for private profiles
		userProfile := dto.UserProfile{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			IsPrivate:   user.IsPrivate,
			CreatedAt:   user.CreatedAt,
		}
		utils.WriteJSON(w, http.StatusOK, userProfile)
		return
	}

	userProfile := dto.UserProfile{
		ID:                user.ID,
		Username:          user.Username,
		DisplayName:       user.DisplayName,
		Bio:               user.Bio,
		ProfilePictureURL: user.ProfilePictureURL,
		IsPrivate:         user.IsPrivate,
		CreatedAt:         user.CreatedAt,
	}

	// Don't include email and email verification for other users
	if requestingUserID == userIDStr {
		userProfile.Email = user.Email
		userProfile.EmailVerified = user.EmailVerified
		userProfile.Theme = user.Theme
	}

	utils.WriteJSON(w, http.StatusOK, userProfile)
}

func (h *ProfileHandler) GetMySessions(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	sessions, err := h.userSessionService.GetUserSessions(claims.UserID)
	if err != nil {
		slog.Error("Failed to get user sessions", "error", err, "user_id", claims.UserID)
		utils.WriteError(w, http.StatusInternalServerError, err, "Failed to get sessions")
		return
	}

	// Convert to response DTOs
	sessionResponses := make([]dto.UserSessionResponse, len(sessions))
	for i, session := range sessions {
		sessionResponses[i] = dto.UserSessionResponse{
			ID:        session.ID,
			IPAddress: session.IPAddress,
			UserAgent: session.UserAgent,
			CreatedAt: session.CreatedAt,
			ExpiresAt: session.ExpiresAt,
			IsCurrent: session.ID == claims.SessionID,
		}
	}

	utils.WriteJSON(w, http.StatusOK, sessionResponses)
}

func (h *ProfileHandler) RevokeSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	sessionIDStr := chi.URLParam(r, "sessionId")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid session ID")
		return
	}

	err = h.userSessionService.RevokeSession(claims.UserID, sessionIDStr)
	if err != nil {
		slog.Error("Failed to revoke session", "error", err, "user_id", claims.UserID, "session_id", sessionID)
		utils.WriteError(w, http.StatusBadRequest, err, "Failed to revoke session")
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Message: "Session revoked successfully",
	})
}

func (h *ProfileHandler) RevokeAllSessions(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, nil, "Unauthorized")
		return
	}

	err := h.userSessionService.RevokeAllSessions(claims.UserID)
	if err != nil {
		slog.Error("Failed to revoke all sessions", "error", err, "user_id", claims.UserID)
		utils.WriteError(w, http.StatusInternalServerError, err, "Failed to revoke sessions")
		return
	}

	utils.WriteJSON(w, http.StatusOK, dto.MessageResponse{
		Message: "All sessions revoked successfully",
	})
}
