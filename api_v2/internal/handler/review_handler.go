package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type ReviewHandler struct {
	reviewService domain.ReviewService
	validator     *validator.Validate
}

func NewReviewHandler(reviewService domain.ReviewService, validator *validator.Validate) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		validator:     validator,
	}
}

func (h *ReviewHandler) Routes(authMiddleware *middleware.AuthMiddleware) chi.Router {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuth)

		r.Post("/", h.CreateReview)
		r.Put("/{id}", h.UpdateReview)
		r.Delete("/{id}", h.DeleteReview)
	})

	// Public routes
	r.Get("/{id}", h.GetReview)
	r.Get("/movie/{movieId}", h.GetMovieReviews)
	r.Get("/user/{userId}", h.GetUserReviews)

	return r
}

// CreateReview creates a new review for a movie
// @Summary Create a movie review
// @Description Create a new review for a movie with rating and/or content
// @Tags Reviews
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateReviewRequest true "Review data"
// @Success 201 {object} dto.ReviewResponse "Review created successfully"
// @Failure 400 {object} utils.ErrorResponse "Invalid request data"
// @Failure 401 {object} utils.ErrorResponse "Unauthorized"
// @Failure 404 {object} utils.ErrorResponse "Movie not found"
// @Failure 409 {object} utils.ErrorResponse "Review already exists"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /reviews [post]
func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	var req dto.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	review, err := h.reviewService.CreateReview(claims.UserID, req.MovieID, req.Rating, req.Content)
	if err != nil {
		if err.Error() == "user already has a review for this movie" {
			utils.WriteJSONError(w, http.StatusConflict, "review_exists", err.Error())
			return
		}
		if err.Error() == "movie not found: movie not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "movie_not_found", "Movie not found")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	response := h.mapToResponse(review)
	utils.WriteJSONResponse(w, r, http.StatusCreated, response)
}

// GetReview returns a specific review by ID
// @Summary Get review by ID
// @Description Get detailed information about a specific review
// @Tags Reviews
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} dto.ReviewResponse "Review details"
// @Failure 400 {object} utils.ErrorResponse "Invalid review ID"
// @Failure 404 {object} utils.ErrorResponse "Review not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /reviews/{id} [get]
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := chi.URLParam(r, "id")
	if reviewIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid review ID")
		return
	}

	review, err := h.reviewService.GetReview(reviewIDStr)
	if err != nil {
		if err.Error() == "failed to get review: review not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "review_not_found", "Review not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "get_failed", err.Error())
		return
	}

	response := h.mapToResponse(review)
	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// GetMovieReviews returns reviews for a specific movie
// @Summary Get movie reviews
// @Description Get all reviews for a specific movie with pagination
// @Tags Reviews
// @Produce json
// @Param movieId path int true "Movie ID"
// @Param page query int false "Page number" default(1)
// @Success 200 {array} dto.ReviewResponse "List of movie reviews"
// @Failure 400 {object} utils.ErrorResponse "Invalid movie ID"
// @Failure 404 {object} utils.ErrorResponse "Movie not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Router /reviews/movie/{movieId} [get]
func (h *ReviewHandler) GetMovieReviews(w http.ResponseWriter, r *http.Request) {
	movieIDStr := chi.URLParam(r, "movieId")
	if movieIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid movie ID")
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	reviews, err := h.reviewService.GetMovieReviews(movieIDStr, page)
	if err != nil {
		if err.Error() == "movie not found: movie not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "movie_not_found", "Movie not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "get_failed", err.Error())
		return
	}

	response := make([]*dto.ReviewResponse, len(reviews))
	for i, review := range reviews {
		response[i] = h.mapToResponse(review)
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// GetUserReviews returns reviews by a specific user
func (h *ReviewHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	if userIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid user ID")
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	reviews, err := h.reviewService.GetUserReviews(userIDStr, page)
	if err != nil {
		if err.Error() == "user not found: user not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "user_not_found", "User not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "get_failed", err.Error())
		return
	}

	response := make([]*dto.ReviewResponse, len(reviews))
	for i, review := range reviews {
		response[i] = h.mapToResponse(review)
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// UpdateReview updates an existing review
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	reviewIDStr := chi.URLParam(r, "id")
	if reviewIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid review ID")
		return
	}

	var req dto.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	review, err := h.reviewService.UpdateReview(reviewIDStr, claims.UserID, req.Rating, req.Content)
	if err != nil {
		if err.Error() == "failed to get review: review not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "review_not_found", "Review not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this review" {
			utils.WriteJSONError(w, http.StatusForbidden, "forbidden", "Unauthorized")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "update_failed", err.Error())
		return
	}

	response := h.mapToResponse(review)
	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// DeleteReview deletes a review
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	reviewIDStr := chi.URLParam(r, "id")
	if reviewIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid review ID")
		return
	}

	err := h.reviewService.DeleteReview(reviewIDStr, claims.UserID)
	if err != nil {
		if err.Error() == "review not found: review not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "review_not_found", "Review not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this review" {
			utils.WriteJSONError(w, http.StatusForbidden, "forbidden", "Unauthorized")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "delete_failed", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Review deleted successfully"})
}

// Helper function to map domain to response
func (h *ReviewHandler) mapToResponse(review *domain.Review) *dto.ReviewResponse {
	response := &dto.ReviewResponse{
		ID:        review.ID,
		Rating:    review.Rating,
		Content:   review.Content,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}

	if review.User != nil {
		response.User = &dto.UserProfile{
			ID:                review.User.ID,
			Username:          review.User.Username,
			DisplayName:       review.User.DisplayName,
			Bio:               review.User.Bio,
			ProfilePictureURL: review.User.ProfilePictureURL,
			IsPrivate:         review.User.IsPrivate,
			CreatedAt:         review.User.CreatedAt,
		}
	}

	if review.Movie != nil {
		response.Movie = &dto.ReviewMovieResponse{
			ID:          review.Movie.ID,
			ExternalID:  review.Movie.ExternalAPIID,
			Title:       review.Movie.Title,
			Overview:    review.Movie.Overview,
			ReleaseDate: review.Movie.ReleaseDate,
			PosterURL:   review.Movie.PosterURL,
			BackdropURL: review.Movie.BackdropURL,
			Genres:      review.Movie.Genres,
			Runtime:     review.Movie.Runtime,
			VoteAverage: review.Movie.VoteAverage,
			VoteCount:   review.Movie.VoteCount,
			Adult:       review.Movie.Adult,
		}
	}

	return response
}
