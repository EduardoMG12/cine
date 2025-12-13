package http

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/user_movie"
)

type FavoriteMovieHandler struct {
	toggleFavoriteUC *user_movie.ToggleFavoriteMovieUseCase
	getFavoriteUC    *user_movie.GetFavoriteMoviesUseCase
}

func NewFavoriteMovieHandler(
	toggleFavoriteUC *user_movie.ToggleFavoriteMovieUseCase,
	getFavoriteUC *user_movie.GetFavoriteMoviesUseCase,
) *FavoriteMovieHandler {
	return &FavoriteMovieHandler{
		toggleFavoriteUC: toggleFavoriteUC,
		getFavoriteUC:    getFavoriteUC,
	}
}

// ToggleFavoriteMovie godoc
// @Summary Toggle movie in favorites
// @Description Add or remove a movie from the authenticated user's favorites list
// @Tags user-movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddFavoriteMovieRequest true "Movie ID to toggle in favorites"
// @Success 200 {object} dto.APIResponse{data=dto.ToggleResponse} "Movie toggled in favorites successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request body"
// @Failure 401 {object} dto.APIResponse "User not authenticated"
// @Failure 404 {object} dto.APIResponse "Movie not found"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /api/v1/favorites [post]
func (h *FavoriteMovieHandler) ToggleFavoriteMovie(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req dto.AddFavoriteMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.toggleFavoriteUC.Execute(userID, req.MovieID)
	if err != nil {
		if err.Error() == "movie not found" {
			sendErrorResponse(w, http.StatusNotFound, "MOVIE_NOT_FOUND", err.Error())
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, result.Message, result)
}

// GetFavoriteMovies godoc
// @Summary Get user's favorite movies
// @Description Get all movies in the authenticated user's favorites list with full movie details
// @Tags user-movies
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.FavoriteMovieWithDetailsDTO} "List of favorite movies retrieved successfully"
// @Failure 401 {object} dto.APIResponse "User not authenticated"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /api/v1/favorites [get]
func (h *FavoriteMovieHandler) GetFavoriteMovies(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	movies, err := h.getFavoriteUC.Execute(userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Favorite movies retrieved successfully", movies)
}
