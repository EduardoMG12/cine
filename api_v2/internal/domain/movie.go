package domain

import (
	"time"
)

type Movie struct {
	ID             int       `json:"id" db:"id"`
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
	GetByID(id int) (*Movie, error)
	GetByExternalID(externalID string) (*Movie, error)
	Update(movie *Movie) error
	Delete(id int) error
	Search(query string, limit, offset int) ([]*Movie, error)
	GetPopular(limit, offset int) ([]*Movie, error)
	GetByGenre(genre string, limit, offset int) ([]*Movie, error)
	DeleteExpiredCache() error
}

type MovieService interface {
	GetMovie(id int) (*Movie, error)
	GetMovieByExternalID(externalID string) (*Movie, error)
	SearchMovies(query string, page int) ([]*Movie, error)
	GetPopularMovies(page int) ([]*Movie, error)
	GetMoviesByGenre(genre string, page int) ([]*Movie, error)
	RefreshMovieCache(externalID string) (*Movie, error)
	CleanupExpiredCache() error
}
