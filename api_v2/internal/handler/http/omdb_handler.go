package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
	"github.com/go-chi/chi/v5"
)

type OMDbHandler struct {
	omdbService *infrastructure.OMDbService
}

func NewOMDbHandler(omdbService *infrastructure.OMDbService) *OMDbHandler {
	return &OMDbHandler{
		omdbService: omdbService,
	}
}

// GetMovieByIMDbID godoc
// @Summary Get movie by IMDb ID from OMDb
// @Description Retrieves detailed movie information from OMDb API using IMDb ID
// @Tags OMDb
// @Accept json
// @Produce json
// @Param imdbId path string true "IMDb ID (e.g., tt0133093)"
// @Success 200 {object} infrastructure.MovieDetails "Movie details"
// @Failure 400 {object} map[string]string "Invalid IMDb ID"
// @Failure 404 {object} map[string]string "Movie not found"
// @Router /api/v1/omdb/{imdbId} [get]
func (h *OMDbHandler) GetMovieByIMDbID(w http.ResponseWriter, r *http.Request) {
	imdbID := chi.URLParam(r, "imdbId")
	if imdbID == "" {
		respondWithError(w, http.StatusBadRequest, "IMDb ID is required")
		return
	}

	movie, err := h.omdbService.GetMovieByExternalID(imdbID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

// GetMovieByTitle godoc
// @Summary Get movie by title from OMDb
// @Description Retrieves movie information from OMDb API using title and optional year
// @Tags OMDb
// @Accept json
// @Produce json
// @Param title query string true "Movie title"
// @Param year query string false "Release year (optional)"
// @Success 200 {object} infrastructure.MovieDetails "Movie details"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 404 {object} map[string]string "Movie not found"
// @Router /api/v1/omdb/title [get]
func (h *OMDbHandler) GetMovieByTitle(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	if title == "" {
		respondWithError(w, http.StatusBadRequest, "title parameter is required")
		return
	}

	year := r.URL.Query().Get("year")

	movie, err := h.omdbService.GetMovieByTitle(title, year)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, movie)
}

// SearchMovies godoc
// @Summary Search movies in OMDb
// @Description Search for movies in OMDb API by query string with pagination
// @Tags OMDb
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} infrastructure.SearchResults "Search results with pagination"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/omdb/search [get]
func (h *OMDbHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	page := getPageFromQuery(r)

	results, err := h.omdbService.SearchMovies(query, page)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

// SearchMoviesByType godoc
// @Summary Search movies by type in OMDb
// @Description Search for content in OMDb API by query and filter by type (movie, series, episode)
// @Tags OMDb
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param type query string true "Content type (movie, series, episode)"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} infrastructure.SearchResults "Search results with pagination"
// @Failure 400 {object} map[string]string "Invalid parameters"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/omdb/search-by-type [get]
func (h *OMDbHandler) SearchMoviesByType(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		respondWithError(w, http.StatusBadRequest, "q parameter is required")
		return
	}

	movieType := r.URL.Query().Get("type") // movie, series, episode
	page := getPageFromQuery(r)

	results, err := h.omdbService.SearchMoviesByType(query, movieType, page)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

// TestConnection godoc
// @Summary Test OMDb API connection
// @Description Tests the connection to OMDb API by fetching a known movie (The Matrix)
// @Tags OMDb
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Connection test result"
// @Router /api/v1/omdb/test [get]
func (h *OMDbHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	// Test with a known movie (The Matrix)
	movie, err := h.omdbService.GetMovieByExternalID("tt0133093")
	if err != nil {
		respondWithJSON(w, http.StatusOK, map[string]interface{}{
			"status":   "error",
			"provider": "OMDb",
			"message":  err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "success",
		"provider":    "OMDb",
		"message":     "Connection successful",
		"test_movie":  movie.Title,
		"test_imdbId": movie.IMDbID,
	})
}

// Helper functions

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func getPageFromQuery(r *http.Request) int {
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	return page
}
