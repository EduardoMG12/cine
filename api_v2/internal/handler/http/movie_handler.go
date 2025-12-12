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
	getTrendingUC      *movie.GetTrendingMoviesUseCase
}

func NewMovieHandler(
	getMovieByIDUC *movie.GetMovieByIDUseCase,
	getRandomUC *movie.GetRandomMovieUseCase,
	getRandomByGenreUC *movie.GetRandomMovieByGenreUseCase,
	searchMoviesUC *movie.SearchMoviesUseCase,
	getTrendingUC *movie.GetTrendingMoviesUseCase,
) *MovieHandler {
	return &MovieHandler{
		getMovieByIDUC:     getMovieByIDUC,
		getRandomUC:        getRandomUC,
		getRandomByGenreUC: getRandomByGenreUC,
		searchMoviesUC:     searchMoviesUC,
		getTrendingUC:      getTrendingUC,
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

// GetTrendingMovies godoc
// @Summary Get trending movies
// @Description Get a shuffled list of trending/popular movies from the local database
// @Tags movies
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]dto.MovieDTO}
// @Failure 500 {object} dto.APIResponse
// @Router /api/v1/movies/trending [get]
//
// ROUTE LOGIC EXPLANATION:
// The OMDb API does not provide any official endpoint for trending or popular movies.
// Therefore, this route simulates a "trending movies" feature using the local database.
//
// FLOW:
//  1. First, it attempts to retrieve 100 movies directly from the database.
//  2. If the database already contains 100 or more movies, the route returns them
//     in random order so the user never sees the same sequence twice.
//  3. If the database does not contain enough movies (< 100), the handler must:
//     a) Iterate through an internal array containing around 300 initial movie name seeds
//     (e.g. "Matrix", "Transformers", "Avatar"…) organized by genre categories.
//     b) For each movie seed, call the existing SearchMoviesUseCase to fetch movie data
//     from OMDb API.
//     c) Each successful search result is saved into the database with the predefined
//     genres from our seed data (since OMDb search doesn't return genres).
//     d) After populating the database with search results, another fetch is performed
//     to gather 100 movies.
//  4. Finally, the route returns these movies shuffled to simulate a "trending" catalog.
//
// PURPOSE:
// - Build a local movie catalog progressively
// - Provide consistent results without depending on a real "trending" API
// - Ensure the user always sees a varied list of movies on the home screen
// - Pre-populate genre data that OMDb doesn't provide in search results
func (h *MovieHandler) GetTrendingMovies(w http.ResponseWriter, r *http.Request) {
	result, err := h.getTrendingUC.Execute()
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "TRENDING_FAILED", err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK, "Trending movies retrieved", result)
}
