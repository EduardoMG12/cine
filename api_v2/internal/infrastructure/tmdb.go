package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
)

const (
	TMDbBaseURL      = "https://api.themoviedb.org/3"
	TMDbImageBaseURL = "https://image.tmdb.org/t/p/w500"
)

type TMDbService struct {
	apiKey     string
	httpClient *http.Client
}

func NewTMDbService(apiKey string) *TMDbService {
	return &TMDbService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetMovie busca detalhes de um filme por ID no TMDb
func (s *TMDbService) GetMovie(tmdbID string) (*dto.TMDbMovieResponse, error) {
	url := fmt.Sprintf("%s/movie/%s?api_key=%s&language=pt-BR", TMDbBaseURL, tmdbID, s.apiKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie from TMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("movie not found in TMDb")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDb API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var movie dto.TMDbMovieResponse
	if err := json.Unmarshal(body, &movie); err != nil {
		return nil, fmt.Errorf("failed to parse TMDb response: %w", err)
	}

	return &movie, nil
}

// SearchMovies busca filmes por nome no TMDb
func (s *TMDbService) SearchMovies(query string, page int) (*dto.TMDbSearchResponse, error) {
	encodedQuery := url.QueryEscape(query)
	url := fmt.Sprintf("%s/search/movie?api_key=%s&language=pt-BR&query=%s&page=%d",
		TMDbBaseURL, s.apiKey, encodedQuery, page)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies in TMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDb API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var searchResp dto.TMDbSearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse TMDb response: %w", err)
	}

	return &searchResp, nil
}

// GetPopularMovies busca filmes populares no TMDb
func (s *TMDbService) GetPopularMovies(page int) (*dto.TMDbDiscoverResponse, error) {
	url := fmt.Sprintf("%s/movie/popular?api_key=%s&language=pt-BR&page=%d",
		TMDbBaseURL, s.apiKey, page)

	return s.fetchMovieList(url)
}

// GetTrendingMovies busca filmes em alta no TMDb
func (s *TMDbService) GetTrendingMovies(timeWindow string) (*dto.TMDbDiscoverResponse, error) {
	if timeWindow != "day" && timeWindow != "week" {
		timeWindow = "week"
	}

	url := fmt.Sprintf("%s/trending/movie/%s?api_key=%s&language=pt-BR",
		TMDbBaseURL, timeWindow, s.apiKey)

	return s.fetchMovieList(url)
}

// GetUpcomingMovies busca próximos lançamentos no TMDb
func (s *TMDbService) GetUpcomingMovies(page int) (*dto.TMDbDiscoverResponse, error) {
	url := fmt.Sprintf("%s/movie/upcoming?api_key=%s&language=pt-BR&page=%d",
		TMDbBaseURL, s.apiKey, page)

	return s.fetchMovieList(url)
}

// GetGenres busca lista de gêneros do TMDb
func (s *TMDbService) GetGenres() (*dto.TMDbGenresResponse, error) {
	url := fmt.Sprintf("%s/genre/movie/list?api_key=%s&language=pt-BR", TMDbBaseURL, s.apiKey)

	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch genres from TMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDb API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var genresResp dto.TMDbGenresResponse
	if err := json.Unmarshal(body, &genresResp); err != nil {
		return nil, fmt.Errorf("failed to parse TMDb response: %w", err)
	}

	return &genresResp, nil
}

// DiscoverByGenre busca filmes por gênero no TMDb
func (s *TMDbService) DiscoverByGenre(genreID int, page int) (*dto.TMDbDiscoverResponse, error) {
	url := fmt.Sprintf("%s/discover/movie?api_key=%s&language=pt-BR&with_genres=%d&page=%d&sort_by=popularity.desc",
		TMDbBaseURL, s.apiKey, genreID, page)

	return s.fetchMovieList(url)
}

// Helper function para buscar lista de filmes
func (s *TMDbService) fetchMovieList(url string) (*dto.TMDbDiscoverResponse, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies from TMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDb API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var discoverResp dto.TMDbDiscoverResponse
	if err := json.Unmarshal(body, &discoverResp); err != nil {
		return nil, fmt.Errorf("failed to parse TMDb response: %w", err)
	}

	return &discoverResp, nil
}

// GetImageURL retorna a URL completa para uma imagem do TMDb
func (s *TMDbService) GetImageURL(path string) string {
	if path == "" {
		return ""
	}
	return TMDbImageBaseURL + path
}

// ConvertTMDbIDToString converte ID do TMDb para string
func ConvertTMDbIDToString(id int) string {
	return strconv.Itoa(id)
}

// ConvertStringToTMDbID converte string para ID do TMDb
func ConvertStringToTMDbID(id string) (int, error) {
	return strconv.Atoi(id)
}
