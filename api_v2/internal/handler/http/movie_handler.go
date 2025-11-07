package http

import (
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/usecase/movie"
	"github.com/go-chi/chi/v5"
)

type MovieHandler struct {
	getMovieByIDUC     *movie.GetMovieByIDUseCase
	getRandomUC        *movie.GetRandomMovieUseCase
	getRandomByGenreUC *movie.GetRandomMovieByGenreUseCase
	searchMoviesUC     *movie.SearchMoviesUseCase
	getPopularUC       *movie.GetPopularMoviesUseCase
	getTrendingUC      *movie.GetTrendingMoviesUseCase
	getGenresUC        *movie.GetGenresUseCase
}

func NewMovieHandler(
	getMovieByIDUC *movie.GetMovieByIDUseCase,
	getRandomUC *movie.GetRandomMovieUseCase,
	getRandomByGenreUC *movie.GetRandomMovieByGenreUseCase,
	searchMoviesUC *movie.SearchMoviesUseCase,
	getPopularUC *movie.GetPopularMoviesUseCase,
	getTrendingUC *movie.GetTrendingMoviesUseCase,
	getGenresUC *movie.GetGenresUseCase,
) *MovieHandler {
	return &MovieHandler{
		getMovieByIDUC:     getMovieByIDUC,
		getRandomUC:        getRandomUC,
		getRandomByGenreUC: getRandomByGenreUC,
		searchMoviesUC:     searchMoviesUC,
		getPopularUC:       getPopularUC,
		getTrendingUC:      getTrendingUC,
		getGenresUC:        getGenresUC,
	}
}

// GetMovieByID godoc
// @Summary Get movie by TMDb ID
// @Description Get detailed information about a specific movie
// @Tags movies
// @Produce json
// @Param id path string true "TMDb Movie ID"
// @Success 200 {object} dto.APIResponse{data=dto.MovieDTO}
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/{id} [get]
func (h *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {
	tmdbID := chi.URLParam(r, "id")
	if tmdbID == "" {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Movie ID is required")
		return
	}

	result, err := h.getMovieByIDUC.Execute(tmdbID)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "MOVIE_NOT_FOUND", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Movie found", result)
}

// GetRandomMovie godoc
// @Summary Get random movie
// @Description Get a random movie from the database
// @Tags movies
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.MovieDTO}
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/random [get]
func (h *MovieHandler) GetRandomMovie(w http.ResponseWriter, r *http.Request) {
	result, err := h.getRandomUC.Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "NO_MOVIES", "No movies found in database")
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Random movie retrieved", result)
}

// GetRandomMovieByGenre godoc
// @Summary Get random movie by genre
// @Description Get a random movie filtered by genre
// @Tags movies
// @Produce json
// @Param genre query string true "Genre name"
// @Success 200 {object} dto.APIResponse{data=dto.MovieDTO}
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/random-by-genre [get]
func (h *MovieHandler) GetRandomMovieByGenre(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	if genre == "" {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_GENRE", "Genre is required")
		return
	}

	result, err := h.getRandomByGenreUC.Execute(genre)
	if err != nil {
		sendErrorResponse(w, http.StatusNotFound, "NO_MOVIES", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Random movie by genre retrieved", result)
}

// SearchMovies godoc
// @Summary Search movies
// @Description Search movies by title in TMDb
// @Tags movies
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.APIResponse{data=dto.TMDbSearchResponse}
// @Failure 400 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/search [get]
func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		sendErrorResponse(w, http.StatusBadRequest, "INVALID_QUERY", "Search query is required")
		return
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	result, err := h.searchMoviesUC.Execute(query, page)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "SEARCH_FAILED", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Movies found", result)
}

// GetPopularMovies godoc
// @Summary Get popular movies
// @Description Get popular movies from TMDb
// @Tags movies
// @Produce json
// @Param page query int false "Page number" default(1)
// @Success 200 {object} dto.APIResponse{data=dto.TMDbDiscoverResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/popular [get]
func (h *MovieHandler) GetPopularMovies(w http.ResponseWriter, r *http.Request) {
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	result, err := h.getPopularUC.Execute(page)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Popular movies retrieved", result)
}

// GetTrendingMovies godoc
// @Summary Get trending movies
// @Description Get trending movies from TMDb
// @Tags movies
// @Produce json
// @Param time_window query string false "Time window (day or week)" default(week)
// @Success 200 {object} dto.APIResponse{data=dto.TMDbDiscoverResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/trending [get]
func (h *MovieHandler) GetTrendingMovies(w http.ResponseWriter, r *http.Request) {
	timeWindow := r.URL.Query().Get("time_window")
	if timeWindow == "" {
		timeWindow = "week"
	}

	result, err := h.getTrendingUC.Execute(timeWindow)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Trending movies retrieved", result)
}

// GetGenres godoc
// @Summary Get movie genres
// @Description Get list of all movie genres from TMDb
// @Tags movies
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.TMDbGenresResponse}
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/genres [get]
func (h *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	result, err := h.getGenresUC.Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "FETCH_FAILED", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Genres retrieved", result)
}
