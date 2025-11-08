package infrastructure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// OMDbService implements the MovieProvider interface for OMDb API
type OMDbService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// OMDb API Response structures (internal)
type omdbMovieResponse struct {
	Title      string       `json:"Title"`
	Year       string       `json:"Year"`
	Rated      string       `json:"Rated"`
	Released   string       `json:"Released"`
	Runtime    string       `json:"Runtime"`
	Genre      string       `json:"Genre"`
	Director   string       `json:"Director"`
	Writer     string       `json:"Writer"`
	Actors     string       `json:"Actors"`
	Plot       string       `json:"Plot"`
	Language   string       `json:"Language"`
	Country    string       `json:"Country"`
	Awards     string       `json:"Awards"`
	Poster     string       `json:"Poster"`
	Ratings    []omdbRating `json:"Ratings"`
	Metascore  string       `json:"Metascore"`
	IMDbRating string       `json:"imdbRating"`
	IMDbVotes  string       `json:"imdbVotes"`
	IMDbID     string       `json:"imdbID"`
	Type       string       `json:"Type"`
	DVD        string       `json:"DVD"`
	BoxOffice  string       `json:"BoxOffice"`
	Production string       `json:"Production"`
	Website    string       `json:"Website"`
	Response   string       `json:"Response"`
	Error      string       `json:"Error,omitempty"`
}

type omdbRating struct {
	Source string `json:"Source"`
	Value  string `json:"Value"`
}

type omdbSearchResponse struct {
	Search       []omdbSearchItem `json:"Search"`
	TotalResults string           `json:"totalResults"`
	Response     string           `json:"Response"`
	Error        string           `json:"Error,omitempty"`
}

type omdbSearchItem struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	IMDbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

func NewOMDbService(apiKey string) *OMDbService {
	return &OMDbService{
		apiKey:  apiKey,
		baseURL: "http://www.omdbapi.com/",
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// GetProviderName returns the name of this provider
func (s *OMDbService) GetProviderName() string {
	return "OMDb"
}

// GetMovieByExternalID fetches movie details by IMDb ID (implements MovieProvider)
func (s *OMDbService) GetMovieByExternalID(imdbID string) (*MovieDetails, error) {
	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("i", imdbID)
	params.Add("plot", "full")
	params.Add("r", "json")

	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())

	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie from OMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDb API returned status code: %d", resp.StatusCode)
	}

	var omdbMovie omdbMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&omdbMovie); err != nil {
		return nil, fmt.Errorf("failed to decode OMDb response: %w", err)
	}

	if omdbMovie.Response == "False" {
		return nil, fmt.Errorf("OMDb API error: %s", omdbMovie.Error)
	}

	return s.convertToMovieDetails(&omdbMovie), nil
}

// GetMovieByTitle fetches movie details by title (implements MovieProvider)
func (s *OMDbService) GetMovieByTitle(title string, year string) (*MovieDetails, error) {
	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("t", title)
	if year != "" {
		params.Add("y", year)
	}
	params.Add("plot", "full")
	params.Add("r", "json")

	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())

	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie from OMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDb API returned status code: %d", resp.StatusCode)
	}

	var omdbMovie omdbMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&omdbMovie); err != nil {
		return nil, fmt.Errorf("failed to decode OMDb response: %w", err)
	}

	if omdbMovie.Response == "False" {
		return nil, fmt.Errorf("OMDb API error: %s", omdbMovie.Error)
	}

	return s.convertToMovieDetails(&omdbMovie), nil
}

// SearchMovies searches for movies by query (implements MovieProvider)
func (s *OMDbService) SearchMovies(query string, page int) (*SearchResults, error) {
	if page < 1 {
		page = 1
	}
	if page > 100 {
		page = 100
	}

	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("s", query)
	params.Add("page", strconv.Itoa(page))
	params.Add("r", "json")

	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())

	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies from OMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDb API returned status code: %d", resp.StatusCode)
	}

	var omdbSearch omdbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&omdbSearch); err != nil {
		return nil, fmt.Errorf("failed to decode OMDb search response: %w", err)
	}

	if omdbSearch.Response == "False" {
		return nil, fmt.Errorf("OMDb API error: %s", omdbSearch.Error)
	}

	return s.convertToSearchResults(&omdbSearch, page), nil
}

// Additional helper methods specific to OMDb

// SearchMoviesByType searches with a specific type filter
func (s *OMDbService) SearchMoviesByType(query string, movieType string, page int) (*SearchResults, error) {
	if page < 1 {
		page = 1
	}
	if page > 100 {
		page = 100
	}

	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("s", query)
	params.Add("page", strconv.Itoa(page))
	params.Add("r", "json")

	if movieType != "" {
		params.Add("type", movieType) // movie, series, episode
	}

	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())

	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies from OMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDb API returned status code: %d", resp.StatusCode)
	}

	var omdbSearch omdbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&omdbSearch); err != nil {
		return nil, fmt.Errorf("failed to decode OMDb search response: %w", err)
	}

	if omdbSearch.Response == "False" {
		return nil, fmt.Errorf("OMDb API error: %s", omdbSearch.Error)
	}

	return s.convertToSearchResults(&omdbSearch, page), nil
}

// GetMovieByIMDbIDWithPlot allows specifying plot length
func (s *OMDbService) GetMovieByIMDbIDWithPlot(imdbID string, plotType string) (*MovieDetails, error) {
	if plotType != "short" && plotType != "full" {
		plotType = "full"
	}

	params := url.Values{}
	params.Add("apikey", s.apiKey)
	params.Add("i", imdbID)
	params.Add("plot", plotType)
	params.Add("r", "json")

	fullURL := fmt.Sprintf("%s?%s", s.baseURL, params.Encode())

	resp, err := s.httpClient.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movie from OMDb: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OMDb API returned status code: %d", resp.StatusCode)
	}

	var omdbMovie omdbMovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&omdbMovie); err != nil {
		return nil, fmt.Errorf("failed to decode OMDb response: %w", err)
	}

	if omdbMovie.Response == "False" {
		return nil, fmt.Errorf("OMDb API error: %s", omdbMovie.Error)
	}

	return s.convertToMovieDetails(&omdbMovie), nil
}

// Conversion helpers

func (s *OMDbService) convertToMovieDetails(omdb *omdbMovieResponse) *MovieDetails {
	ratings := make([]Rating, len(omdb.Ratings))
	for i, r := range omdb.Ratings {
		ratings[i] = Rating{
			Source: r.Source,
			Value:  r.Value,
		}
	}

	return &MovieDetails{
		Title:      omdb.Title,
		Year:       omdb.Year,
		Released:   omdb.Released,
		Runtime:    omdb.Runtime,
		Plot:       omdb.Plot,
		Type:       omdb.Type,
		Poster:     omdb.Poster,
		Rated:      omdb.Rated,
		Genre:      omdb.Genre,
		Language:   omdb.Language,
		Country:    omdb.Country,
		Director:   omdb.Director,
		Writer:     omdb.Writer,
		Actors:     omdb.Actors,
		IMDbID:     omdb.IMDbID,
		IMDbRating: omdb.IMDbRating,
		IMDbVotes:  omdb.IMDbVotes,
		Metascore:  omdb.Metascore,
		Ratings:    ratings,
		Awards:     omdb.Awards,
		BoxOffice:  omdb.BoxOffice,
		Production: omdb.Production,
		Website:    omdb.Website,
		Provider:   "OMDb",
		ProviderID: omdb.IMDbID,
	}
}

func (s *OMDbService) convertToSearchResults(omdb *omdbSearchResponse, page int) *SearchResults {
	results := make([]SearchItem, len(omdb.Search))
	for i, item := range omdb.Search {
		results[i] = SearchItem{
			Title:      item.Title,
			Year:       item.Year,
			Type:       item.Type,
			Poster:     item.Poster,
			IMDbID:     item.IMDbID,
			ProviderID: item.IMDbID,
			// Genre is not available in search results from OMDb API
			// Only available in detailed movie fetch
		}
	}

	totalResults, _ := strconv.Atoi(omdb.TotalResults)
	totalPages := (totalResults + 9) / 10 // OMDb returns 10 results per page

	return &SearchResults{
		Results:      results,
		TotalResults: totalResults,
		Page:         page,
		TotalPages:   totalPages,
		Provider:     "OMDb",
	}
}
