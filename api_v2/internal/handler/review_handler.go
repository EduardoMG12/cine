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

	review, err := h.reviewService.CreateReview(claims.UserID, req.MovieID, &req.Rating, req.Content)
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
	utils.WriteJSONResponse(w, http.StatusCreated, response)
}

// GetReview returns a specific review by ID
func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	reviewIDStr := chi.URLParam(r, "id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid review ID")
		return
	}

	review, err := h.reviewService.GetReview(reviewID)
	if err != nil {
		if err.Error() == "failed to get review: review not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "review_not_found", "Review not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "get_failed", err.Error())
		return
	}

	response := h.mapToResponse(review)
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetMovieReviews returns reviews for a specific movie
func (h *ReviewHandler) GetMovieReviews(w http.ResponseWriter, r *http.Request) {
	movieIDStr := chi.URLParam(r, "movieId")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid movie ID")
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	reviews, err := h.reviewService.GetMovieReviews(movieID, page)
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

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetUserReviews returns reviews by a specific user
func (h *ReviewHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid user ID")
		return
	}

	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	reviews, err := h.reviewService.GetUserReviews(userID, page)
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

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// UpdateReview updates an existing review
func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	reviewIDStr := chi.URLParam(r, "id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
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

	review, err := h.reviewService.UpdateReview(reviewID, claims.UserID, req.Rating, req.Content)
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
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// DeleteReview deletes a review
func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		utils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
		return
	}

	reviewIDStr := chi.URLParam(r, "id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_id", "Invalid review ID")
		return
	}

	err = h.reviewService.DeleteReview(reviewID, claims.UserID)
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

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Review deleted successfully"})
}

// Helper function to map domain to response
func (h *ReviewHandler) mapToResponse(review *domain.Review) *dto.ReviewResponse {
	var rating int
	if review.Rating != nil {
		rating = *review.Rating
	}

	response := &dto.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		MovieID:   review.MovieID,
		Rating:    rating,
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
		response.Movie = &dto.MovieResponse{
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
