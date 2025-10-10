package domain

import (
	"time"
)

type Movie struct {
	ID             string    `json:"id" db:"id"`
	ExternalAPIID  string    `json:"external_api_id" db:"external_api_id"`
	Title          string    `json:"title" db:"title"`
	Overview       *string   `json:"overview,omitempty" db:"overview"`
	ReleaseDate    *string   `json:"release_date,omitempty" db:"release_date"`
	PosterURL      *string   `json:"poster_url,omitempty" db:"poster_url"`
	BackdropURL    *string   `json:"backdrop_url,omitempty" db:"backdrop_url"`
	Genres         []string  `json:"genres" db:"genres"`
	Runtime        *int      `json:"runtime,omitempty" db:"runtime"`
	VoteAverage    *float64  `json:"vote_average,omitempty" db:"vote_average"`
	VoteCount      *int      `json:"vote_count,omitempty" db:"vote_count"`
	Adult          bool      `json:"adult" db:"adult"`
	CacheExpiresAt time.Time `json:"cache_expires_at" db:"cache_expires_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type MovieRepository interface {
	Create(movie *Movie) error
	GetByID(id string) (*Movie, error)
	GetByExternalID(externalID string) (*Movie, error)
	Update(movie *Movie) error
	Delete(id string) error
	Search(query string, limit, offset int) ([]*Movie, error)
	GetPopular(limit, offset int) ([]*Movie, error)
	GetByGenre(genre string, limit, offset int) ([]*Movie, error)
	DeleteExpiredCache() error
}

type MovieService interface {
	GetMovie(id string) (*Movie, error)
	GetMovieByExternalID(externalID string) (*Movie, error)
	SearchMovies(query string, page int) ([]*Movie, error)
	GetPopularMovies(page int) ([]*Movie, error)
	GetMoviesByGenre(genre string, page int) ([]*Movie, error)
	RefreshMovieCache(externalID string) (*Movie, error)
	CleanupExpiredCache() error
}

type TMDbService interface {
	SearchMovies(query string, page int) (*TMDbSearchResponse, error)
	GetMovieDetails(movieID int) (*TMDbMovie, error)
	GetPopularMovies(page int) (*TMDbSearchResponse, error)
	GetGenres() ([]*TMDbGenre, error)
	DiscoverMovies(genreID int, page int) (*TMDbSearchResponse, error)
	TMDbMovieToDomain(tmdbMovie *TMDbMovie, genres []*TMDbGenre) *Movie
}

type TMDbSearchResponse struct {
	Page         int         `json:"page"`
	Results      []TMDbMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

type TMDbMovie struct {
	ID               int         `json:"id"`
	Title            string      `json:"title"`
	OriginalTitle    string      `json:"original_title"`
	Overview         string      `json:"overview"`
	ReleaseDate      string      `json:"release_date"`
	PosterPath       *string     `json:"poster_path"`
	BackdropPath     *string     `json:"backdrop_path"`
	GenreIDs         []int       `json:"genre_ids"`
	Genres           []TMDbGenre `json:"genres,omitempty"` // For detailed movie response
	Runtime          *int        `json:"runtime,omitempty"`
	VoteAverage      float64     `json:"vote_average"`
	VoteCount        int         `json:"vote_count"`
	Adult            bool        `json:"adult"`
	Popularity       float64     `json:"popularity"`
	Video            bool        `json:"video"`
	OriginalLanguage string      `json:"original_language"`
}

type TMDbGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
