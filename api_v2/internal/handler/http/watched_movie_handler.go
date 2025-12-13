package http

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
	"github.com/EduardoMG12/cine/api_v2/internal/usecase/user_movie"
)

type WatchedMovieHandler struct {
	toggleWatchedUC *user_movie.ToggleWatchedMovieUseCase
	getWatchedUC    *user_movie.GetWatchedMoviesUseCase
}

func NewWatchedMovieHandler(
	toggleWatchedUC *user_movie.ToggleWatchedMovieUseCase,
	getWatchedUC *user_movie.GetWatchedMoviesUseCase,
) *WatchedMovieHandler {
	return &WatchedMovieHandler{
		toggleWatchedUC: toggleWatchedUC,
		getWatchedUC:    getWatchedUC,
	}
}

// ToggleWatchedMovie godoc
// @Summary Toggle movie in watched list
// @Description Add or remove a movie from the authenticated user's watched list
// @Tags user-movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.AddWatchedMovieRequest true "Movie ID to toggle in watched list"
// @Success 200 {object} dto.APIResponse{data=dto.ToggleResponse} "Movie toggled in watched list successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request body"
// @Failure 401 {object} dto.APIResponse "User not authenticated"
// @Failure 404 {object} dto.APIResponse "Movie not found"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /api/v1/watched [post]
func (h *WatchedMovieHandler) ToggleWatchedMovie(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	var req dto.AddWatchedMovieRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.toggleWatchedUC.Execute(userID, req.MovieID)
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

// GetWatchedMovies godoc
// @Summary Get user's watched movies
// @Description Get all movies in the authenticated user's watched list with full movie details
// @Tags user-movies
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=[]dto.WatchedMovieWithDetailsDTO} "List of watched movies retrieved successfully"
// @Failure 401 {object} dto.APIResponse "User not authenticated"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /api/v1/watched [get]
func (h *WatchedMovieHandler) GetWatchedMovies(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		sendErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "User not authenticated")
		return
	}

	movies, err := h.getWatchedUC.Execute(userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Watched movies retrieved successfully", movies)
}
