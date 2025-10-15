package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type MovieListHandler struct {
	movieListService domain.MovieListService
	movieService     domain.MovieService
	validator        *validator.Validate
}

func NewMovieListHandler(movieListService domain.MovieListService, movieService domain.MovieService, validator *validator.Validate) *MovieListHandler {
	return &MovieListHandler{
		movieListService: movieListService,
		movieService:     movieService,
		validator:        validator,
	}
}

func (h *MovieListHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// Custom movie lists
	r.Post("/", h.CreateList)
	r.Get("/", h.GetUserLists)
	r.Get("/{id}", h.GetList)
	r.Put("/{id}", h.UpdateList)
	r.Delete("/{id}", h.DeleteList)

	// List operations
	r.Get("/{id}/movies", h.GetListMovies)
	r.Post("/{id}/movies", h.AddMovieToList)
	r.Delete("/{id}/movies/{movieId}", h.RemoveMovieFromList)

	// Default lists operations
	r.Post("/want-to-watch", h.AddToWantToWatch)
	r.Delete("/want-to-watch/{movieId}", h.RemoveFromWantToWatch)
	r.Post("/watched", h.AddToWatched)
	r.Delete("/watched/{movieId}", h.RemoveFromWatched)
	r.Post("/move-to-watched", h.MoveToWatched)

	return r
}

// CreateList creates a new custom movie list
func (h *MovieListHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req dto.CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "invalid_request", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	list, err := h.movieListService.CreateList(userID, req.Name)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "create_failed", err.Error())
		return
	}

	response := h.mapToResponse(list)
	utils.WriteJSONResponse(w, r, http.StatusCreated, response)
}

// GetUserLists returns all lists for the authenticated user
func (h *MovieListHandler) GetUserLists(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	lists, err := h.movieListService.GetUserLists(userID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	response := make([]*dto.MovieListResponse, len(lists))
	for i, list := range lists {
		response[i] = h.mapToResponse(list)
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// GetList returns a specific movie list by ID
func (h *MovieListHandler) GetList(w http.ResponseWriter, r *http.Request) {
	listIDStr := chi.URLParam(r, "id")
	if listIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid list ID")
		return
	}

	list, err := h.movieListService.GetList(listIDStr)
	if err != nil {
		if err.Error() == "movie list not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie list not found")
			return
		}
		utils.WriteJSONError(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	response := h.mapToResponse(list)
	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// UpdateList updates a movie list
func (h *MovieListHandler) UpdateList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	listIDStr := chi.URLParam(r, "id")
	if listIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid list ID")
		return
	}

	var req dto.UpdateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	var name string
	if req.Name != nil {
		name = *req.Name
	}
	list, err := h.movieListService.UpdateList(listIDStr, userID, name)
	if err != nil {
		if err.Error() == "movie list not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie list not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this list" {
			utils.WriteJSONError(w, http.StatusForbidden, "error", "Unauthorized")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	response := h.mapToResponse(list)
	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// DeleteList deletes a movie list
func (h *MovieListHandler) DeleteList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	listIDStr := chi.URLParam(r, "id")
	if listIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid list ID")
		return
	}

	err := h.movieListService.DeleteList(listIDStr, userID)
	if err != nil {
		if err.Error() == "movie list not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie list not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this list" {
			utils.WriteJSONError(w, http.StatusForbidden, "error", "Unauthorized")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "List deleted successfully"})
}

// AddToWantToWatch adds a movie to the want-to-watch list
func (h *MovieListHandler) AddToWantToWatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req dto.AddMovieToListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	err := h.movieListService.AddToWantToWatch(userID, req.MovieID)
	if err != nil {
		if err.Error() == "movie not found: movie not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie not found")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie added to want-to-watch list"})
}

// AddToWatched adds a movie to the watched list
func (h *MovieListHandler) AddToWatched(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req dto.AddMovieToListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	err := h.movieListService.AddToWatched(userID, req.MovieID)
	if err != nil {
		if err.Error() == "movie not found: movie not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie not found")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie added to watched list"})
}

// RemoveFromWantToWatch removes a movie from the want-to-watch list
func (h *MovieListHandler) RemoveFromWantToWatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	movieIDStr := chi.URLParam(r, "movieId")
	if movieIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid movie ID")
		return
	}

	err := h.movieListService.RemoveFromWantToWatch(userID, movieIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie removed from want-to-watch list"})
}

// RemoveFromWatched removes a movie from the watched list
func (h *MovieListHandler) RemoveFromWatched(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	movieIDStr := chi.URLParam(r, "movieId")
	if movieIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid movie ID")
		return
	}

	err := h.movieListService.RemoveFromWatched(userID, movieIDStr)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie removed from watched list"})
}

// MoveToWatched moves a movie from want-to-watch to watched
func (h *MovieListHandler) MoveToWatched(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req dto.MoveToWatchedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	// First, find the movie by external ID
	movie, err := h.movieService.GetMovieByExternalID(req.MovieExternalID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie not found")
		return
	}

	err = h.movieListService.MoveToWatched(userID, movie.ID)
	if err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie moved to watched list"})
}

// AddMovieToList adds a movie to a custom list
func (h *MovieListHandler) AddMovieToList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	listIDStr := chi.URLParam(r, "id")
	if listIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid list ID")
		return
	}

	var req dto.AddMovieToListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid request body")
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		utils.WriteValidationError(w, err)
		return
	}

	err := h.movieListService.AddMovieToList(listIDStr, userID, req.MovieID)
	if err != nil {
		if err.Error() == "movie list not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie list not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this list" {
			utils.WriteJSONError(w, http.StatusForbidden, "error", "Unauthorized")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie added to list"})
}

// RemoveMovieFromList removes a movie from a custom list
func (h *MovieListHandler) RemoveMovieFromList(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	listIDStr := chi.URLParam(r, "id")
	if listIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid list ID")
		return
	}

	movieIDStr := chi.URLParam(r, "movieId")
	if movieIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid movie ID")
		return
	}

	err := h.movieListService.RemoveMovieFromList(listIDStr, userID, movieIDStr)
	if err != nil {
		if err.Error() == "movie list not found" {
			utils.WriteJSONError(w, http.StatusNotFound, "error", "Movie list not found")
			return
		}
		if err.Error() == "unauthorized: user does not own this list" {
			utils.WriteJSONError(w, http.StatusForbidden, "error", "Unauthorized")
			return
		}
		utils.WriteJSONError(w, http.StatusBadRequest, "error", err.Error())
		return
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, map[string]string{"message": "Movie removed from list"})
}

// GetListMovies returns movies in a specific list
func (h *MovieListHandler) GetListMovies(w http.ResponseWriter, r *http.Request) {
	listIDStr := chi.URLParam(r, "id")
	if listIDStr == "" {
		utils.WriteJSONError(w, http.StatusBadRequest, "error", "Invalid list ID")
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

	entries, err := h.movieListService.GetListMovies(listIDStr, page)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "error", err.Error())
		return
	}

	response := make([]*dto.MovieListEntryResponse, len(entries))
	for i, entry := range entries {
		response[i] = &dto.MovieListEntryResponse{
			ID:      entry.ID,
			AddedAt: entry.AddedAt,
		}
		if entry.Movie != nil {
			response[i].Movie = dto.MovieResponse{
				ID:          entry.Movie.ID,
				ExternalID:  entry.Movie.ExternalAPIID,
				Title:       entry.Movie.Title,
				Overview:    entry.Movie.Overview,
				PosterURL:   entry.Movie.PosterURL,
				BackdropURL: entry.Movie.BackdropURL,
				ReleaseDate: entry.Movie.ReleaseDate,
				Genres:      entry.Movie.Genres,
				VoteAverage: entry.Movie.VoteAverage,
				Runtime:     entry.Movie.Runtime,
				VoteCount:   entry.Movie.VoteCount,
				Adult:       entry.Movie.Adult,
			}
		}
	}

	utils.WriteJSONResponse(w, r, http.StatusOK, response)
}

// Helper function to map domain to response
func (h *MovieListHandler) mapToResponse(list *domain.MovieList) *dto.MovieListResponse {
	response := &dto.MovieListResponse{
		ID:          list.ID,
		UserID:      list.UserID,
		Name:        list.Name,
		Description: nil,   // Domain doesn't have Description field
		IsPublic:    false, // Domain doesn't have IsPublic field, defaulting to false
		MovieCount:  len(list.Entries),
		CreatedAt:   list.CreatedAt,
		UpdatedAt:   list.UpdatedAt,
	}

	if list.Entries != nil {
		response.Movies = make([]dto.MovieListEntryResponse, len(list.Entries))
		for i, entry := range list.Entries {
			response.Movies[i] = dto.MovieListEntryResponse{
				ID:      entry.ID,
				AddedAt: entry.AddedAt,
			}
			if entry.Movie != nil {
				movieResponse := dto.MovieResponse{
					ID:          entry.Movie.ID,
					ExternalID:  entry.Movie.ExternalAPIID,
					Title:       entry.Movie.Title,
					Overview:    entry.Movie.Overview,
					PosterURL:   entry.Movie.PosterURL,
					BackdropURL: entry.Movie.BackdropURL,
					ReleaseDate: entry.Movie.ReleaseDate,
					Genres:      entry.Movie.Genres,
					VoteAverage: entry.Movie.VoteAverage,
					Runtime:     entry.Movie.Runtime,
					VoteCount:   entry.Movie.VoteCount,
					Adult:       entry.Movie.Adult,
				}
				response.Movies[i].Movie = movieResponse
			}
		}
	}

	return response
}
