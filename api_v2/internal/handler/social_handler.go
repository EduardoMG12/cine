package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
)

type SocialHandler struct {
	socialService domain.SocialService
}

func NewSocialHandler(socialService domain.SocialService) *SocialHandler {
	return &SocialHandler{
		socialService: socialService,
	}
}

func (h *SocialHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Friendship routes
	r.Post("/friend-request/{userID}", h.SendFriendRequest)
	r.Post("/friend-request/{userID}/accept", h.AcceptFriendRequest)
	r.Delete("/friend-request/{userID}/decline", h.DeclineFriendRequest)
	r.Delete("/friend/{userID}", h.RemoveFriend)
	r.Post("/block/{userID}", h.BlockUser)
	r.Get("/friends", h.GetFriends)
	r.Get("/friend-requests", h.GetFriendRequests)

	// Follow routes
	r.Post("/follow/{userID}", h.FollowUser)
	r.Delete("/follow/{userID}", h.UnfollowUser)
	r.Get("/followers", h.GetFollowers)
	r.Get("/following", h.GetFollowing)
	r.Get("/followers/count", h.GetFollowersCount)
	r.Get("/following/count", h.GetFollowingCount)

	return r
}

// Friendship endpoints

// SendFriendRequest sends a friend request to another user
// @Summary Send friend request
// @Description Send a friend request to another user
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "Target user ID"
// @Success 200 {object} dto.MessageResponse "Friend request sent"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 409 {object} utils.ErrorResponse "Request already exists"
// @Router /social/friend-request/{userID} [post]
func (h *SocialHandler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.SendFriendRequest(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to send friend request", "error", err, "sender", claims.UserID, "receiver", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "FRIEND_REQUEST_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "FRIEND_REQUEST_SENT"),
	})
}

// AcceptFriendRequest accepts a pending friend request
// @Summary Accept friend request
// @Description Accept a pending friend request from another user
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "Requester user ID"
// @Success 200 {object} dto.MessageResponse "Friend request accepted"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/friend-request/{userID}/accept [post]
func (h *SocialHandler) AcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.AcceptFriendRequest(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to accept friend request", "error", err, "user", claims.UserID, "requester", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "FRIEND_REQUEST_ACCEPT_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "FRIEND_REQUEST_ACCEPTED"),
	})
}

// DeclineFriendRequest declines a pending friend request
// @Summary Decline friend request
// @Description Decline a pending friend request from another user
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "Requester user ID"
// @Success 200 {object} dto.MessageResponse "Friend request declined"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/friend-request/{userID}/decline [delete]
func (h *SocialHandler) DeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.DeclineFriendRequest(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to decline friend request", "error", err, "user", claims.UserID, "requester", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "FRIEND_REQUEST_DECLINE_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "FRIEND_REQUEST_DECLINED"),
	})
}

// RemoveFriend removes a friend from user's friend list
// @Summary Remove friend
// @Description Remove a friend from user's friend list
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "Friend user ID"
// @Success 200 {object} dto.MessageResponse "Friend removed"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/friend/{userID} [delete]
func (h *SocialHandler) RemoveFriend(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.RemoveFriend(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to remove friend", "error", err, "user", claims.UserID, "friend", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "FRIEND_REMOVE_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "FRIEND_REMOVED"),
	})
}

// BlockUser blocks another user
// @Summary Block user
// @Description Block another user to prevent interactions
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "User ID to block"
// @Success 200 {object} dto.MessageResponse "User blocked"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/block/{userID} [post]
func (h *SocialHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.BlockUser(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to block user", "error", err, "user", claims.UserID, "blocked", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "USER_BLOCK_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "USER_BLOCKED"),
	})
}

// GetFriends returns user's friends list
// @Summary Get friends
// @Description Get the current user's friends list
// @Tags Social
// @Security BearerAuth
// @Success 200 {array} dto.UserProfile "Friends list"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/friends [get]
func (h *SocialHandler) GetFriends(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	friends, err := h.socialService.GetFriends(claims.UserID)
	if err != nil {
		slog.Error("Failed to get friends", "error", err, "user", claims.UserID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
		return
	}

	// Convert to DTOs
	var friendProfiles []dto.UserProfile
	for _, friend := range friends {
		friendProfiles = append(friendProfiles, dto.UserProfile{
			ID:                friend.ID,
			Username:          friend.Username,
			DisplayName:       friend.DisplayName,
			Bio:               friend.Bio,
			ProfilePictureURL: friend.ProfilePictureURL,
			IsPrivate:         friend.IsPrivate,
			CreatedAt:         friend.CreatedAt,
		})
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, friendProfiles)
}

// GetFriendRequests returns pending friend requests
// @Summary Get friend requests
// @Description Get pending friend requests for the current user
// @Tags Social
// @Security BearerAuth
// @Success 200 {array} dto.UserProfile "Friend requests list"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/friend-requests [get]
func (h *SocialHandler) GetFriendRequests(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	requesters, err := h.socialService.GetFriendRequests(claims.UserID)
	if err != nil {
		slog.Error("Failed to get friend requests", "error", err, "user", claims.UserID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
		return
	}

	// Convert to DTOs
	var requesterProfiles []dto.UserProfile
	for _, requester := range requesters {
		requesterProfiles = append(requesterProfiles, dto.UserProfile{
			ID:                requester.ID,
			Username:          requester.Username,
			DisplayName:       requester.DisplayName,
			Bio:               requester.Bio,
			ProfilePictureURL: requester.ProfilePictureURL,
			IsPrivate:         requester.IsPrivate,
			CreatedAt:         requester.CreatedAt,
		})
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, requesterProfiles)
}

// Follow endpoints

// FollowUser follows another user
// @Summary Follow user
// @Description Follow another user to see their activities
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "User ID to follow"
// @Success 200 {object} dto.MessageResponse "User followed"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/follow/{userID} [post]
func (h *SocialHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.FollowUser(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to follow user", "error", err, "follower", claims.UserID, "following", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "FOLLOW_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "USER_FOLLOWED"),
	})
}

// UnfollowUser unfollows another user
// @Summary Unfollow user
// @Description Unfollow a user to stop seeing their activities
// @Tags Social
// @Security BearerAuth
// @Param userID path int true "User ID to unfollow"
// @Success 200 {object} dto.MessageResponse "User unfollowed"
// @Failure 400 {object} utils.ErrorResponse "Invalid request"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/follow/{userID} [delete]
func (h *SocialHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	lang := middleware.GetLanguageFromContext(r.Context())
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	userIDStr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "INVALID_USER_ID", nil)
		return
	}

	err = h.socialService.UnfollowUser(claims.UserID, userID)
	if err != nil {
		slog.Warn("Failed to unfollow user", "error", err, "follower", claims.UserID, "following", userID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusBadRequest, "UNFOLLOW_FAILED", map[string]interface{}{
			"details": err.Error(),
		})
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, dto.MessageResponse{
		Message: i18n.T(lang, "USER_UNFOLLOWED"),
	})
}

// GetFollowers returns user's followers
// @Summary Get followers
// @Description Get the current user's followers with pagination
// @Tags Social
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Success 200 {array} dto.UserProfile "Followers list"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/followers [get]
func (h *SocialHandler) GetFollowers(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	followers, err := h.socialService.GetFollowers(claims.UserID, page)
	if err != nil {
		slog.Error("Failed to get followers", "error", err, "user", claims.UserID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
		return
	}

	// Convert to DTOs
	var followerProfiles []dto.UserProfile
	for _, follower := range followers {
		followerProfiles = append(followerProfiles, dto.UserProfile{
			ID:                follower.ID,
			Username:          follower.Username,
			DisplayName:       follower.DisplayName,
			Bio:               follower.Bio,
			ProfilePictureURL: follower.ProfilePictureURL,
			IsPrivate:         follower.IsPrivate,
			CreatedAt:         follower.CreatedAt,
		})
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, followerProfiles)
}

// GetFollowing returns users that current user is following
// @Summary Get following
// @Description Get users that the current user is following with pagination
// @Tags Social
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Success 200 {array} dto.UserProfile "Following list"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/following [get]
func (h *SocialHandler) GetFollowing(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	following, err := h.socialService.GetFollowing(claims.UserID, page)
	if err != nil {
		slog.Error("Failed to get following", "error", err, "user", claims.UserID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
		return
	}

	// Convert to DTOs
	var followingProfiles []dto.UserProfile
	for _, user := range following {
		followingProfiles = append(followingProfiles, dto.UserProfile{
			ID:                user.ID,
			Username:          user.Username,
			DisplayName:       user.DisplayName,
			Bio:               user.Bio,
			ProfilePictureURL: user.ProfilePictureURL,
			IsPrivate:         user.IsPrivate,
			CreatedAt:         user.CreatedAt,
		})
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, followingProfiles)
}

// GetFollowersCount returns the count of user's followers
// @Summary Get followers count
// @Description Get the count of the current user's followers
// @Tags Social
// @Security BearerAuth
// @Success 200 {object} map[string]int "Followers count"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/followers/count [get]
func (h *SocialHandler) GetFollowersCount(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	count, err := h.socialService.GetFollowersCount(claims.UserID)
	if err != nil {
		slog.Error("Failed to get followers count", "error", err, "user", claims.UserID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, map[string]int{
		"count": count,
	})
}

// GetFollowingCount returns the count of users that current user is following
// @Summary Get following count
// @Description Get the count of users that the current user is following
// @Tags Social
// @Security BearerAuth
// @Success 200 {object} map[string]int "Following count"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Router /social/following/count [get]
func (h *SocialHandler) GetFollowingCount(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteErrorResponse(w, r.Context(), http.StatusUnauthorized, "UNAUTHORIZED", nil)
		return
	}

	count, err := h.socialService.GetFollowingCount(claims.UserID)
	if err != nil {
		slog.Error("Failed to get following count", "error", err, "user", claims.UserID)
		utils.WriteErrorResponse(w, r.Context(), http.StatusInternalServerError, "INTERNAL_ERROR", nil)
		return
	}

	utils.WriteJSONResponse(w, r.Context(), http.StatusOK, map[string]int{
		"count": count,
	})
}
