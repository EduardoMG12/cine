package infrastructure

// MovieProvider defines the interface for external movie data providers
// This allows easy switching between OMDb, TMDb, or other APIs
type MovieProvider interface {
	// GetMovieByExternalID fetches movie details by external ID (IMDb ID, TMDb ID, etc.)
	GetMovieByExternalID(id string) (*MovieDetails, error)

	// GetMovieByTitle fetches movie details by title and optional year
	GetMovieByTitle(title string, year string) (*MovieDetails, error)

	// SearchMovies searches for movies by query
	SearchMovies(query string, page int) (*SearchResults, error)

	// GetProviderName returns the name of the provider (e.g., "OMDb", "TMDb")
	GetProviderName() string
}

// MovieDetails represents unified movie details from any provider
type MovieDetails struct {
	// Basic Info
	Title    string `json:"title"`
	Year     string `json:"year"`
	Released string `json:"released"`
	Runtime  string `json:"runtime"`
	Plot     string `json:"plot"`
	Type     string `json:"type"` // movie, series, episode

	// Media
	Poster      string `json:"poster"`
	BackdropURL string `json:"backdrop_url,omitempty"`

	// Classification
	Rated    string `json:"rated"`
	Genre    string `json:"genre"`
	Language string `json:"language"`
	Country  string `json:"country"`

	// Credits
	Director string `json:"director"`
	Writer   string `json:"writer"`
	Actors   string `json:"actors"`

	// Ratings
	IMDbID     string   `json:"imdb_id"`
	IMDbRating string   `json:"imdb_rating"`
	IMDbVotes  string   `json:"imdb_votes"`
	Metascore  string   `json:"metascore"`
	Ratings    []Rating `json:"ratings"`

	// Additional Info
	Awards     string `json:"awards"`
	BoxOffice  string `json:"box_office"`
	Production string `json:"production"`
	Website    string `json:"website"`

	// Provider Info
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}

// Rating represents a rating from a specific source
type Rating struct {
	Source string `json:"source"`
	Value  string `json:"value"`
}

// SearchResults represents unified search results from any provider
type SearchResults struct {
	Results      []SearchItem `json:"results"`
	TotalResults int          `json:"total_results"`
	Page         int          `json:"page"`
	TotalPages   int          `json:"total_pages"`
	Provider     string       `json:"provider"`
}

// SearchItem represents a single search result item
type SearchItem struct {
	Title      string `json:"title"`
	Year       string `json:"year"`
	Type       string `json:"type"`
	Poster     string `json:"poster"`
	IMDbID     string `json:"imdb_id,omitempty"`
	ProviderID string `json:"provider_id"`
}
