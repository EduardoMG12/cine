package dto

import "time"

// Movie List DTOs
type CreateMovieListRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type UpdateMovieListRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

type MovieListResponse struct {
	ID        int                       `json:"id"`
	Name      string                    `json:"name"`
	IsDefault bool                      `json:"is_default"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
	User      *UserProfile              `json:"user,omitempty"`
	Entries   []*MovieListEntryResponse `json:"entries,omitempty"`
}

type MovieListEntryResponse struct {
	ID      int                     `json:"id"`
	AddedAt time.Time               `json:"added_at"`
	Movie   *MovieListMovieResponse `json:"movie,omitempty"`
}

// Specific movie response for lists to avoid import issues
type MovieListMovieResponse struct {
	ID          int      `json:"id"`
	ExternalID  string   `json:"external_id"`
	Title       string   `json:"title"`
	Overview    *string  `json:"overview,omitempty"`
	ReleaseDate *string  `json:"release_date,omitempty"`
	PosterURL   *string  `json:"poster_url,omitempty"`
	BackdropURL *string  `json:"backdrop_url,omitempty"`
	Genres      []string `json:"genres"`
	Runtime     *int     `json:"runtime,omitempty"`
	VoteAverage *float64 `json:"vote_average,omitempty"`
	VoteCount   *int     `json:"vote_count,omitempty"`
	Adult       bool     `json:"adult"`
}

type AddMovieToListRequest struct {
	MovieID int `json:"movie_id" validate:"required,gt=0"`
}

type RemoveMovieFromListRequest struct {
	MovieID int `json:"movie_id" validate:"required,gt=0"`
}

// Default Lists Operations
type DefaultListRequest struct {
	MovieID int `json:"movie_id" validate:"required,gt=0"`
}

type MoveToWatchedRequest struct {
	MovieID int `json:"movie_id" validate:"required,gt=0"`
}

// List Query Parameters
type MovieListQuery struct {
	Page  int `json:"page" query:"page" validate:"min=1"`
	Limit int `json:"limit" query:"limit" validate:"min=1,max=100"`
}
