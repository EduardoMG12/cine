package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/config"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type TMDbService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewTMDbService(cfg *config.Config) *TMDbService {
	return &TMDbService{
		apiKey:  cfg.TMDb.APIKey,
		baseURL: cfg.TMDb.BaseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TMDbService) SearchMovies(query string, page int) (*domain.TMDbSearchResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("query", query)
	params.Set("page", strconv.Itoa(page))

	endpoint := fmt.Sprintf("%s/search/movie?%s", s.baseURL, params.Encode())

	var response domain.TMDbSearchResponse
	if err := s.makeRequest(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	return &response, nil
}

func (s *TMDbService) GetMovieDetails(movieID int) (*domain.TMDbMovie, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)

	endpoint := fmt.Sprintf("%s/movie/%d?%s", s.baseURL, movieID, params.Encode())

	var movie domain.TMDbMovie
	if err := s.makeRequest(endpoint, &movie); err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}

	return &movie, nil
}

func (s *TMDbService) GetPopularMovies(page int) (*domain.TMDbSearchResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("page", strconv.Itoa(page))

	endpoint := fmt.Sprintf("%s/movie/popular?%s", s.baseURL, params.Encode())

	var response domain.TMDbSearchResponse
	if err := s.makeRequest(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}

	return &response, nil
}

func (s *TMDbService) GetGenres() ([]*domain.TMDbGenre, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)

	endpoint := fmt.Sprintf("%s/genre/movie/list?%s", s.baseURL, params.Encode())

	var response struct {
		Genres []*domain.TMDbGenre `json:"genres"`
	}

	if err := s.makeRequest(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}

	return response.Genres, nil
}

func (s *TMDbService) DiscoverMovies(genreID int, page int) (*domain.TMDbSearchResponse, error) {
	params := url.Values{}
	params.Set("api_key", s.apiKey)
	params.Set("with_genres", strconv.Itoa(genreID))
	params.Set("page", strconv.Itoa(page))

	endpoint := fmt.Sprintf("%s/discover/movie?%s", s.baseURL, params.Encode())

	var response domain.TMDbSearchResponse
	if err := s.makeRequest(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to discover movies: %w", err)
	}

	return &response, nil
}

func (s *TMDbService) makeRequest(endpoint string, target interface{}) error {
	resp, err := s.httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("TMDb API returned status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// TMDbMovieToDomain converts TMDb movie data to domain model
func (s *TMDbService) TMDbMovieToDomain(tmdbMovie *domain.TMDbMovie, genres []*domain.TMDbGenre) *domain.Movie {
	genreMap := make(map[int]string)
	for _, genre := range genres {
		genreMap[genre.ID] = genre.Name
	}

	var movieGenres []string
	for _, genreID := range tmdbMovie.GenreIDs {
		if genreName, exists := genreMap[genreID]; exists {
			movieGenres = append(movieGenres, genreName)
		}
	}

	movie := &domain.Movie{
		ExternalAPIID:  strconv.Itoa(tmdbMovie.ID),
		Title:          tmdbMovie.Title,
		Genres:         movieGenres,
		Adult:          tmdbMovie.Adult,
		CacheExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour cache
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if tmdbMovie.Overview != "" {
		movie.Overview = &tmdbMovie.Overview
	}

	if tmdbMovie.ReleaseDate != "" {
		movie.ReleaseDate = &tmdbMovie.ReleaseDate
	}

	if tmdbMovie.PosterPath != nil && *tmdbMovie.PosterPath != "" {
		posterURL := fmt.Sprintf("https://image.tmdb.org/t/p/w500%s", *tmdbMovie.PosterPath)
		movie.PosterURL = &posterURL
	}

	if tmdbMovie.BackdropPath != nil && *tmdbMovie.BackdropPath != "" {
		backdropURL := fmt.Sprintf("https://image.tmdb.org/t/p/w1280%s", *tmdbMovie.BackdropPath)
		movie.BackdropURL = &backdropURL
	}

	if tmdbMovie.Runtime != nil && *tmdbMovie.Runtime > 0 {
		movie.Runtime = tmdbMovie.Runtime
	}

	if tmdbMovie.VoteAverage > 0 {
		movie.VoteAverage = &tmdbMovie.VoteAverage
	}

	if tmdbMovie.VoteCount > 0 {
		movie.VoteCount = &tmdbMovie.VoteCount
	}

	return movie
}
