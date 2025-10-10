package handler

import (
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/service"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
)

type MovieHandler struct {
	movieService domain.MovieService
}

func NewMovieHandler(movieService domain.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/search", h.SearchMovies)
	r.Get("/popular", h.GetPopularMovies)
	r.Get("/genre/{genre}", h.GetMoviesByGenre)
	r.Get("/{id}", h.GetMovie)
	r.Get("/external/{externalId}", h.GetMovieByExternalID)

	return r
}

// SearchMovies handles movie search requests
// @Summary Search movies
// @Description Search for movies by title
// @Tags Movies
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.MovieSearchResponse
// @Failure 400 {object} utils.ErrorResponse
// @Router /movies/search [get]
func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "MISSING_QUERY", "Search query is required")
		return
	}

	if len(query) < 2 {
		utils.WriteJSONError(w, http.StatusBadRequest, "QUERY_TOO_SHORT", "Search query must be at least 2 characters")
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	movies, err := h.movieService.SearchMovies(query, page)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "SEARCH_FAILED", "Failed to search movies")
		return
	}

	response := dto.MovieSearchResponse{
		Movies:     h.convertMoviesToDTO(movies),
		Page:       page,
		TotalPages: 1, // This would need to be calculated based on total results
		TotalCount: len(movies),
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetPopularMovies handles requests for popular movies
// @Summary Get popular movies
// @Description Get a list of popular movies
// @Tags Movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.MovieSearchResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /movies/popular [get]
func (h *MovieHandler) GetPopularMovies(w http.ResponseWriter, r *http.Request) {
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	movies, err := h.movieService.GetPopularMovies(page)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get popular movies")
		return
	}

	response := dto.MovieSearchResponse{
		Movies:     h.convertMoviesToDTO(movies),
		Page:       page,
		TotalPages: 1,
		TotalCount: len(movies),
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetMoviesByGenre handles requests for movies by genre
// @Summary Get movies by genre
// @Description Get movies filtered by genre
// @Tags Movies
// @Produce json
// @Param genre path string true "Genre name"
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.MovieSearchResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /movies/genre/{genre} [get]
func (h *MovieHandler) GetMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	genre := chi.URLParam(r, "genre")
	if genre == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "MISSING_GENRE", "Genre is required")
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	movies, err := h.movieService.GetMoviesByGenre(genre, page)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get movies by genre")
		return
	}

	response := dto.MovieSearchResponse{
		Movies:     h.convertMoviesToDTO(movies),
		Page:       page,
		TotalPages: 1,
		TotalCount: len(movies),
	}

	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetMovie handles requests for a specific movie by internal ID
// @Summary Get movie by ID
// @Description Get movie details by internal database ID
// @Tags Movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} dto.MovieResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /movies/{id} [get]
func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "INVALID_ID", "Invalid movie ID")
		return
	}

	movie, err := h.movieService.GetMovie(id)
	if err != nil {
		if err == service.ErrMovieNotFound {
			utils.WriteJSONError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "Movie not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get movie")
		return
	}

	response := h.convertMovieToDTO(movie)
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// GetMovieByExternalID handles requests for a movie by TMDb external ID
// @Summary Get movie by external ID
// @Description Get movie details by TMDb external ID
// @Tags Movies
// @Produce json
// @Param externalId path string true "TMDb Movie ID"
// @Success 200 {object} dto.MovieResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /movies/external/{externalId} [get]
func (h *MovieHandler) GetMovieByExternalID(w http.ResponseWriter, r *http.Request) {
	externalID := chi.URLParam(r, "externalId")
	if externalID == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "INVALID_EXTERNAL_ID", "Invalid external ID")
		return
	}

	movie, err := h.movieService.GetMovieByExternalID(externalID)
	if err != nil {
		if err == service.ErrMovieNotFound {
			utils.WriteJSONError(w, http.StatusNotFound, "MOVIE_NOT_FOUND", "Movie not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "FETCH_FAILED", "Failed to get movie")
		return
	}

	response := h.convertMovieToDTO(movie)
	utils.WriteJSONResponse(w, http.StatusOK, response)
}

// Helper methods

func (h *MovieHandler) convertMovieToDTO(movie *domain.Movie) dto.MovieResponse {
	return dto.MovieResponse{
		ID:          movie.ID,
		ExternalID:  movie.ExternalAPIID,
		Title:       movie.Title,
		Overview:    movie.Overview,
		ReleaseDate: movie.ReleaseDate,
		PosterURL:   movie.PosterURL,
		BackdropURL: movie.BackdropURL,
		Genres:      movie.Genres,
		Runtime:     movie.Runtime,
		VoteAverage: movie.VoteAverage,
		VoteCount:   movie.VoteCount,
		Adult:       movie.Adult,
	}
}

func (h *MovieHandler) convertMoviesToDTO(movies []*domain.Movie) []dto.MovieResponse {
	var response []dto.MovieResponse
	for _, movie := range movies {
		response = append(response, h.convertMovieToDTO(movie))
	}
	return response
}
