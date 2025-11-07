package domain

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	ExternalAPIID  string     `db:"external_api_id" json:"external_api_id"` // TMDb ID
	Title          string     `db:"title" json:"title"`
	Overview       *string    `db:"overview" json:"overview,omitempty"`
	ReleaseDate    *time.Time `db:"release_date" json:"release_date,omitempty"`
	PosterURL      *string    `db:"poster_url" json:"poster_url,omitempty"`
	BackdropURL    *string    `db:"backdrop_url" json:"backdrop_url,omitempty"`
	Genres         []string   `db:"genres" json:"genres"`
	Runtime        *int       `db:"runtime" json:"runtime,omitempty"` // minutes
	VoteAverage    *float64   `db:"vote_average" json:"vote_average,omitempty"`
	VoteCount      *int       `db:"vote_count" json:"vote_count,omitempty"`
	Adult          bool       `db:"adult" json:"adult"`
	CacheExpiresAt time.Time  `db:"cache_expires_at" json:"-"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MovieRepository defines the repository interface for movie operations
type MovieRepository interface {
	CreateMovie(movie *Movie) error
	GetMovieByID(id uuid.UUID) (*Movie, error)
	GetMovieByExternalID(externalID string) (*Movie, error)
	UpdateMovie(movie *Movie) error
	DeleteMovie(id uuid.UUID) error
	GetRandomMovie() (*Movie, error)
	GetRandomMovieByGenre(genre string) (*Movie, error)
	SearchMovies(query string, limit int) ([]*Movie, error)
}
