package service

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

var (
	ErrMovieNotFound = errors.New("movie not found")
)

type MovieService struct {
	movieRepo   domain.MovieRepository
	tmdbService domain.TMDbService
}

func NewMovieService(movieRepo domain.MovieRepository, tmdbService domain.TMDbService) domain.MovieService {
	return &MovieService{
		movieRepo:   movieRepo,
		tmdbService: tmdbService,
	}
}

func (s *MovieService) GetMovie(id string) (*domain.Movie, error) {
	movie, err := s.movieRepo.GetByID(id)
	if err != nil {
		return nil, ErrMovieNotFound
	}

	// Check if cache is expired
	if time.Now().After(movie.CacheExpiresAt) {
		// Try to refresh from TMDb
		if refreshed, err := s.refreshMovieFromTMDb(movie.ExternalAPIID); err == nil {
			return refreshed, nil
		}
		// If refresh fails, return cached version
		slog.Warn("Failed to refresh movie cache, returning cached version", "movie_id", id, "error", err)
	}

	return movie, nil
}

func (s *MovieService) GetMovieByExternalID(externalID string) (*domain.Movie, error) {
	// Try to get from local database first
	movie, err := s.movieRepo.GetByExternalID(externalID)
	if err == nil {
		// Check if cache is expired
		if time.Now().After(movie.CacheExpiresAt) {
			// Try to refresh from TMDb
			if refreshed, err := s.refreshMovieFromTMDb(externalID); err == nil {
				return refreshed, nil
			}
			// If refresh fails, return cached version
			slog.Warn("Failed to refresh movie cache, returning cached version", "external_id", externalID, "error", err)
		}
		return movie, nil
	}

	// Movie not in database, fetch from TMDb
	return s.fetchAndStoreFromTMDb(externalID)
}

func (s *MovieService) SearchMovies(query string, page int) ([]*domain.Movie, error) {
	if page < 1 {
		page = 1
	}

	// Search in local database first
	limit := 20
	offset := (page - 1) * limit
	localMovies, err := s.movieRepo.Search(query, limit, offset)
	if err == nil && len(localMovies) > 0 {
		return localMovies, nil
	}

	// If no local results or error, search TMDb
	tmdbResponse, err := s.tmdbService.SearchMovies(query, page)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	// Convert TMDb movies to domain movies and store them
	var movies []*domain.Movie
	genres, _ := s.tmdbService.GetGenres() // Get genres for conversion

	for _, tmdbMovie := range tmdbResponse.Results {
		movie := s.tmdbService.TMDbMovieToDomain(&tmdbMovie, genres)

		// Try to store in database (ignore errors for existing movies)
		if _, err := s.movieRepo.GetByExternalID(movie.ExternalAPIID); err != nil {
			s.movieRepo.Create(movie)
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *MovieService) GetPopularMovies(page int) ([]*domain.Movie, error) {
	if page < 1 {
		page = 1
	}

	// Try local database first
	limit := 20
	offset := (page - 1) * limit
	localMovies, err := s.movieRepo.GetPopular(limit, offset)
	if err == nil && len(localMovies) >= 10 { // If we have enough local popular movies
		return localMovies, nil
	}

	// Fetch from TMDb
	tmdbResponse, err := s.tmdbService.GetPopularMovies(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}

	// Convert and store
	var movies []*domain.Movie
	genres, _ := s.tmdbService.GetGenres()

	for _, tmdbMovie := range tmdbResponse.Results {
		movie := s.tmdbService.TMDbMovieToDomain(&tmdbMovie, genres)

		// Try to update existing or create new
		if existing, err := s.movieRepo.GetByExternalID(movie.ExternalAPIID); err == nil {
			// Update existing movie
			movie.ID = existing.ID
			movie.CreatedAt = existing.CreatedAt
			s.movieRepo.Update(movie)
		} else {
			// Create new movie
			s.movieRepo.Create(movie)
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *MovieService) GetMoviesByGenre(genre string, page int) ([]*domain.Movie, error) {
	if page < 1 {
		page = 1
	}

	// Try local database first
	limit := 20
	offset := (page - 1) * limit
	localMovies, err := s.movieRepo.GetByGenre(genre, limit, offset)
	if err == nil && len(localMovies) > 0 {
		return localMovies, nil
	}

	// Find genre ID for TMDb API
	genres, err := s.tmdbService.GetGenres()
	if err != nil {
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}

	var genreID int
	for _, g := range genres {
		if g.Name == genre {
			genreID = g.ID
			break
		}
	}

	if genreID == 0 {
		return nil, fmt.Errorf("genre not found: %s", genre)
	}

	// Fetch from TMDb
	tmdbResponse, err := s.tmdbService.DiscoverMovies(genreID, page)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies by genre: %w", err)
	}

	// Convert and store
	var movies []*domain.Movie
	for _, tmdbMovie := range tmdbResponse.Results {
		movie := s.tmdbService.TMDbMovieToDomain(&tmdbMovie, genres)

		// Try to update existing or create new
		if existing, err := s.movieRepo.GetByExternalID(movie.ExternalAPIID); err == nil {
			movie.ID = existing.ID
			movie.CreatedAt = existing.CreatedAt
			s.movieRepo.Update(movie)
		} else {
			s.movieRepo.Create(movie)
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (s *MovieService) RefreshMovieCache(externalID string) (*domain.Movie, error) {
	return s.refreshMovieFromTMDb(externalID)
}

func (s *MovieService) CleanupExpiredCache() error {
	return s.movieRepo.DeleteExpiredCache()
}

// Helper methods

func (s *MovieService) fetchAndStoreFromTMDb(externalID string) (*domain.Movie, error) {
	movieIDInt, err := strconv.Atoi(externalID)
	if err != nil {
		return nil, fmt.Errorf("invalid external ID: %s", externalID)
	}

	tmdbMovie, err := s.tmdbService.GetMovieDetails(movieIDInt)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie from TMDb: %w", err)
	}

	genres, _ := s.tmdbService.GetGenres()
	movie := s.tmdbService.TMDbMovieToDomain(tmdbMovie, genres)

	if err := s.movieRepo.Create(movie); err != nil {
		return nil, fmt.Errorf("failed to store movie: %w", err)
	}

	return movie, nil
}

func (s *MovieService) refreshMovieFromTMDb(externalID string) (*domain.Movie, error) {
	movieIDInt, err := strconv.Atoi(externalID)
	if err != nil {
		return nil, fmt.Errorf("invalid external ID: %s", externalID)
	}

	tmdbMovie, err := s.tmdbService.GetMovieDetails(movieIDInt)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh movie from TMDb: %w", err)
	}

	genres, _ := s.tmdbService.GetGenres()
	movie := s.tmdbService.TMDbMovieToDomain(tmdbMovie, genres)

	// Update existing movie
	existing, err := s.movieRepo.GetByExternalID(externalID)
	if err != nil {
		// Movie doesn't exist, create it
		if err := s.movieRepo.Create(movie); err != nil {
			return nil, fmt.Errorf("failed to create movie: %w", err)
		}
		return movie, nil
	}

	// Update existing
	movie.ID = existing.ID
	movie.CreatedAt = existing.CreatedAt
	if err := s.movieRepo.Update(movie); err != nil {
		return nil, fmt.Errorf("failed to update movie: %w", err)
	}

	return movie, nil
}
