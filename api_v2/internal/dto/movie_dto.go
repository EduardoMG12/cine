package dto

import (
	"time"

	"github.com/google/uuid"
)

// Movie DTOs
type MovieDTO struct {
	ID            uuid.UUID  `json:"id"`
	ExternalAPIID string     `json:"external_api_id"`
	Title         string     `json:"title"`
	Overview      *string    `json:"overview,omitempty"`
	ReleaseDate   *time.Time `json:"release_date,omitempty"`
	PosterURL     *string    `json:"poster_url,omitempty"`
	BackdropURL   *string    `json:"backdrop_url,omitempty"`
	Genres        []string   `json:"genres"`
	Runtime       *int       `json:"runtime,omitempty"`
	VoteAverage   *float64   `json:"vote_average,omitempty"`
	VoteCount     *int       `json:"vote_count,omitempty"`
	Adult         bool       `json:"adult"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type GenreDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MovieSearchQuery struct {
	Query string `json:"query" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"omitempty,min=1,max=50"`
}

type RandomMovieByGenreQuery struct {
	Genre string `json:"genre" validate:"required"`
}

// TMDb API Response structures
type TMDbMovieResponse struct {
	ID                  int           `json:"id"`
	Title               string        `json:"title"`
	Overview            string        `json:"overview"`
	ReleaseDate         string        `json:"release_date"`
	PosterPath          string        `json:"poster_path"`
	BackdropPath        string        `json:"backdrop_path"`
	Genres              []TMDbGenre   `json:"genres"`
	Runtime             int           `json:"runtime"`
	VoteAverage         float64       `json:"vote_average"`
	VoteCount           int           `json:"vote_count"`
	Adult               bool          `json:"adult"`
	ProductionCompanies []interface{} `json:"production_companies"`
}

type TMDbGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TMDbSearchResponse struct {
	Page         int                     `json:"page"`
	Results      []TMDbMovieSearchResult `json:"results"`
	TotalPages   int                     `json:"total_pages"`
	TotalResults int                     `json:"total_results"`
}

type TMDbMovieSearchResult struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Overview     string  `json:"overview"`
	ReleaseDate  string  `json:"release_date"`
	PosterPath   string  `json:"poster_path"`
	BackdropPath string  `json:"backdrop_path"`
	GenreIDs     []int   `json:"genre_ids"`
	VoteAverage  float64 `json:"vote_average"`
	VoteCount    int     `json:"vote_count"`
	Adult        bool    `json:"adult"`
}

type TMDbGenresResponse struct {
	Genres []TMDbGenre `json:"genres"`
}

type TMDbDiscoverResponse struct {
	Page         int                     `json:"page"`
	Results      []TMDbMovieSearchResult `json:"results"`
	TotalPages   int                     `json:"total_pages"`
	TotalResults int                     `json:"total_results"`
}
