package dto

import "time"

// Review DTOs
type CreateReviewRequest struct {
	MovieID int     `json:"movie_id" validate:"required,gt=0"`
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=10"`
	Content *string `json:"content,omitempty" validate:"omitempty,min=1,max=2000"`
}

type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=10"`
	Content *string `json:"content,omitempty" validate:"omitempty,min=1,max=2000"`
}

type ReviewResponse struct {
	ID        int                  `json:"id"`
	Rating    *int                 `json:"rating,omitempty"`
	Content   *string              `json:"content,omitempty"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	User      *UserProfile         `json:"user,omitempty"`
	Movie     *ReviewMovieResponse `json:"movie,omitempty"`
}

// Specific movie response for reviews to avoid import cycles
type ReviewMovieResponse struct {
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

// Query Parameters
type ReviewQuery struct {
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}
