package infrastructure

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/lib/pq"
)

// MovieFetcher defines the interface for fetching movies from different sources
type MovieFetcher interface {
	FetchByExternalID(externalID string) (*domain.Movie, error)
	FetchByTitle(title string, year string) (*domain.Movie, error)
	Search(query string, page int) ([]*domain.Movie, error)
	GetProviderName() string
	SetNext(fetcher MovieFetcher)
}

// BaseMovieFetcher provides common chain of responsibility functionality
type BaseMovieFetcher struct {
	next MovieFetcher
}

func (b *BaseMovieFetcher) SetNext(fetcher MovieFetcher) {
	b.next = fetcher
}

func (b *BaseMovieFetcher) tryNext(operation string, args ...interface{}) (*domain.Movie, error) {
	if b.next == nil {
		return nil, fmt.Errorf("no more providers available")
	}

	switch operation {
	case "fetchByExternalID":
		return b.next.FetchByExternalID(args[0].(string))
	case "fetchByTitle":
		return b.next.FetchByTitle(args[0].(string), args[1].(string))
	default:
		return nil, fmt.Errorf("unknown operation: %s", operation)
	}
}

func (b *BaseMovieFetcher) tryNextSearch(query string, page int) ([]*domain.Movie, error) {
	if b.next == nil {
		return nil, fmt.Errorf("no more providers available")
	}
	return b.next.Search(query, page)
}

// OMDbMovieFetcher wraps OMDbService to work with chain of responsibility
type OMDbMovieFetcher struct {
	BaseMovieFetcher
	omdbService *OMDbService
	movieRepo   domain.MovieRepository
	autoSave    bool
}

func NewOMDbMovieFetcher(omdbService *OMDbService, movieRepo domain.MovieRepository, autoSave bool) *OMDbMovieFetcher {
	return &OMDbMovieFetcher{
		omdbService: omdbService,
		movieRepo:   movieRepo,
		autoSave:    autoSave,
	}
}

func (o *OMDbMovieFetcher) GetProviderName() string {
	return "OMDb"
}

func (o *OMDbMovieFetcher) FetchByExternalID(externalID string) (*domain.Movie, error) {
	log.Printf("[MovieFetcher] Trying OMDb for external ID: %s", externalID)

	movieDetails, err := o.omdbService.GetMovieByExternalID(externalID)
	if err != nil {
		log.Printf("[MovieFetcher] OMDb failed: %v. Trying next provider...", err)
		return o.tryNext("fetchByExternalID", externalID)
	}

	movie := o.convertToMovie(movieDetails)

	if o.autoSave {
		if err := o.saveToDatabase(movie); err != nil {
			log.Printf("[MovieFetcher] Failed to save to database: %v", err)
		} else {
			log.Printf("[MovieFetcher] Movie saved to database: %s", movie.Title)
		}
	}

	return movie, nil
}

func (o *OMDbMovieFetcher) FetchByTitle(title string, year string) (*domain.Movie, error) {
	log.Printf("[MovieFetcher] Trying OMDb for title: %s (year: %s)", title, year)

	movieDetails, err := o.omdbService.GetMovieByTitle(title, year)
	if err != nil {
		log.Printf("[MovieFetcher] OMDb failed: %v. Trying next provider...", err)
		return o.tryNext("fetchByTitle", title, year)
	}

	movie := o.convertToMovie(movieDetails)

	if o.autoSave {
		if err := o.saveToDatabase(movie); err != nil {
			log.Printf("[MovieFetcher] Failed to save to database: %v", err)
		} else {
			log.Printf("[MovieFetcher] Movie saved to database: %s", movie.Title)
		}
	}

	return movie, nil
}

func (o *OMDbMovieFetcher) Search(query string, page int) ([]*domain.Movie, error) {
	log.Printf("[MovieFetcher] Trying OMDb for search: %s (page: %d)", query, page)

	searchResults, err := o.omdbService.SearchMovies(query, page)
	if err != nil {
		log.Printf("[MovieFetcher] OMDb search failed: %v. Trying next provider...", err)
		return o.tryNextSearch(query, page)
	}

	movies := make([]*domain.Movie, 0, len(searchResults.Results))
	for _, item := range searchResults.Results {
		movie := &domain.Movie{
			ExternalAPIID: item.IMDbID,
			Provider:      "omdb",
			Title:         item.Title,
			LastSyncAt:    timePtr(time.Now()),
			Genres:        pq.StringArray{}, // Initialize as empty array
		}

		if item.Year != "" {
			movie.Overview = &item.Year
		}
		if item.Poster != "" && item.Poster != "N/A" {
			movie.PosterURL = &item.Poster
		}

		// Note: OMDb search API doesn't return Genre
		// Genres will be populated when user requests movie details (FetchByExternalID)

		if o.autoSave {
			if err := o.saveToDatabase(movie); err != nil {
				log.Printf("[MovieFetcher] Failed to save search result: %v", err)
			}
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (o *OMDbMovieFetcher) convertToMovie(details *MovieDetails) *domain.Movie {
	movie := &domain.Movie{
		ExternalAPIID:  details.IMDbID,
		Provider:       "omdb",
		Title:          details.Title,
		Adult:          details.Rated == "R" || details.Rated == "NC-17",
		LastSyncAt:     timePtr(time.Now()),
		CacheExpiresAt: time.Now().Add(48 * time.Hour), // 2 days cache
	}

	if details.Plot != "" && details.Plot != "N/A" {
		movie.Overview = &details.Plot
	}
	if details.Released != "" && details.Released != "N/A" {
		if releaseDate, err := time.Parse("02 Jan 2006", details.Released); err == nil {
			movie.ReleaseDate = &releaseDate
		}
	}
	if details.Poster != "" && details.Poster != "N/A" {
		movie.PosterURL = &details.Poster
	}
	if details.Runtime != "" && details.Runtime != "N/A" {
		var runtime int
		if _, err := fmt.Sscanf(details.Runtime, "%d min", &runtime); err == nil {
			movie.Runtime = &runtime
		}
	}
	if details.IMDbRating != "" && details.IMDbRating != "N/A" {
		var rating float64
		if _, err := fmt.Sscanf(details.IMDbRating, "%f", &rating); err == nil {
			movie.VoteAverage = &rating
		}
	}
	if details.IMDbVotes != "" && details.IMDbVotes != "N/A" {
		var votes int
		cleanVotes := ""
		for _, char := range details.IMDbVotes {
			if char >= '0' && char <= '9' {
				cleanVotes += string(char)
			}
		}
		if _, err := fmt.Sscanf(cleanVotes, "%d", &votes); err == nil {
			movie.VoteCount = &votes
		}
	}
	if details.Genre != "" && details.Genre != "N/A" {
		movie.Genres = convertGenreStringToSlice(details.Genre)
	}

	return movie
}

func (o *OMDbMovieFetcher) saveToDatabase(movie *domain.Movie) error {
	existing, err := o.movieRepo.GetMovieByExternalID(movie.ExternalAPIID)
	if err == nil && existing != nil {
		movie.ID = existing.ID
		return o.movieRepo.UpdateMovie(movie)
	}
	return o.movieRepo.CreateMovie(movie)
}

func splitByComma(s string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if char == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else if char != ' ' || current != "" {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// DatabaseMovieFetcher is the final fallback that searches in local database
type DatabaseMovieFetcher struct {
	BaseMovieFetcher
	movieRepo domain.MovieRepository
}

func NewDatabaseMovieFetcher(movieRepo domain.MovieRepository) *DatabaseMovieFetcher {
	return &DatabaseMovieFetcher{
		movieRepo: movieRepo,
	}
}

func (d *DatabaseMovieFetcher) GetProviderName() string {
	return "Database (Local Cache)"
}

func (d *DatabaseMovieFetcher) FetchByExternalID(externalID string) (*domain.Movie, error) {
	log.Printf("[MovieFetcher] Trying local database for external ID: %s", externalID)

	movie, err := d.movieRepo.GetMovieByExternalID(externalID)
	if err != nil {
		log.Printf("[MovieFetcher] Database lookup failed: %v", err)
		return nil, fmt.Errorf("movie not found in any provider or local database")
	}

	log.Printf("[MovieFetcher] Movie found in local database: %s", movie.Title)
	return movie, nil
}

func (d *DatabaseMovieFetcher) FetchByTitle(title string, year string) (*domain.Movie, error) {
	log.Printf("[MovieFetcher] Trying local database for title: %s", title)

	movies, err := d.movieRepo.SearchMovies(title, 1)
	if err != nil || len(movies) == 0 {
		log.Printf("[MovieFetcher] Database search failed or no results")
		return nil, fmt.Errorf("movie not found in any provider or local database")
	}

	// If year is specified, try to match it
	if year != "" {
		for _, movie := range movies {
			if movie.ReleaseDate != nil && movie.ReleaseDate.Format("2006") == year {
				log.Printf("[MovieFetcher] Movie found in local database: %s (%s)", movie.Title, year)
				return movie, nil
			}
		}
	}

	// Return first result if year not specified or no year match
	log.Printf("[MovieFetcher] Movie found in local database: %s", movies[0].Title)
	return movies[0], nil
}

func (d *DatabaseMovieFetcher) Search(query string, page int) ([]*domain.Movie, error) {
	log.Printf("[MovieFetcher] Trying local database for search: %s", query)

	limit := 10
	movies, err := d.movieRepo.SearchMovies(query, limit)
	if err != nil {
		log.Printf("[MovieFetcher] Database search failed: %v", err)
		return nil, fmt.Errorf("no results found in any provider or local database")
	}

	if len(movies) == 0 {
		log.Printf("[MovieFetcher] No results in local database")
		return nil, fmt.Errorf("no results found")
	}

	log.Printf("[MovieFetcher] Found %d movies in local database", len(movies))
	return movies, nil
}

// NewMovieFetcherChain creates a complete chain: OMDb -> Database
func NewMovieFetcherChain(omdbService *OMDbService, movieRepo domain.MovieRepository, autoSave bool) MovieFetcher {
	// Create fetchers
	omdbFetcher := NewOMDbMovieFetcher(omdbService, movieRepo, autoSave)
	dbFetcher := NewDatabaseMovieFetcher(movieRepo)

	omdbFetcher.SetNext(dbFetcher)

	return omdbFetcher
}

func convertGenreStringToSlice(genreStr string) pq.StringArray {
	if genreStr == "" || genreStr == "N/A" {
		return pq.StringArray{}
	}

	// OMDb returns genres as "Action, Adventure, Sci-Fi"
	// Split by comma and trim spaces
	genres := strings.Split(genreStr, ", ")
	result := make([]string, 0, len(genres))
	for _, genre := range genres {
		trimmed := strings.TrimSpace(genre)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return pq.StringArray(result)
}
