package dto

import "time"

// Movie DTOs
type MovieResponse struct {
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

type MovieSearchRequest struct {
	Query string `json:"query" validate:"required,min=2"`
	Page  *int   `json:"page,omitempty" validate:"omitempty,min=1"`
}

type MovieSearchResponse struct {
	Movies     []MovieResponse `json:"movies"`
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
	TotalCount int             `json:"total_count"`
}

// Movie List DTOs
type CreateListRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}

type UpdateListRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	IsPublic    *bool   `json:"is_public,omitempty"`
}

type AddMovieToListRequest struct {
	MovieID int `json:"movie_id" validate:"required,min=1"`
}

type MoveMovieRequest struct {
	FromListID int `json:"from_list_id" validate:"required,min=1"`
	ToListID   int `json:"to_list_id" validate:"required,min=1"`
	MovieID    int `json:"movie_id" validate:"required,min=1"`
}

type MovieListResponse struct {
	ID          int                      `json:"id"`
	UserID      int                      `json:"user_id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description,omitempty"`
	IsPublic    bool                     `json:"is_public"`
	MovieCount  int                      `json:"movie_count"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Movies      []MovieListEntryResponse `json:"movies,omitempty"`
}

type MovieListEntryResponse struct {
	ID       int           `json:"id"`
	Position *int          `json:"position,omitempty"`
	AddedAt  time.Time     `json:"added_at"`
	Movie    MovieResponse `json:"movie"`
}

// Quick action DTOs for common operations
type AddToWantToWatchRequest struct {
	MovieExternalID string `json:"movie_external_id" validate:"required"`
}

type AddToWatchedRequest struct {
	MovieExternalID string `json:"movie_external_id" validate:"required"`
}

type MoveToWatchedRequest struct {
	MovieExternalID string `json:"movie_external_id" validate:"required"`
}
