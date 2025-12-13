package dto

import (
	"time"

	"github.com/google/uuid"
)

// WatchedMovieDTO represents a watched movie response
type WatchedMovieDTO struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	MovieID   uuid.UUID `json:"movie_id"`
	WatchedAt time.Time `json:"watched_at"`
	CreatedAt time.Time `json:"created_at"`
}

// FavoriteMovieDTO represents a favorite movie response
type FavoriteMovieDTO struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	MovieID     uuid.UUID `json:"movie_id"`
	FavoritedAt time.Time `json:"favorited_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// AddWatchedMovieRequest represents request to add a movie to watched list
type AddWatchedMovieRequest struct {
	MovieID uuid.UUID `json:"movie_id" validate:"required"`
}

// RemoveWatchedMovieRequest represents request to remove a movie from watched list
type RemoveWatchedMovieRequest struct {
	MovieID uuid.UUID `json:"movie_id" validate:"required"`
}

// AddFavoriteMovieRequest represents request to add a movie to favorites
type AddFavoriteMovieRequest struct {
	MovieID uuid.UUID `json:"movie_id" validate:"required"`
}

// RemoveFavoriteMovieRequest represents request to remove a movie from favorites
type RemoveFavoriteMovieRequest struct {
	MovieID uuid.UUID `json:"movie_id" validate:"required"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message"`
}

// WatchedMovieWithDetailsDTO represents a watched movie with full movie details
type WatchedMovieWithDetailsDTO struct {
	WatchedMovie WatchedMovieDTO `json:"watched_movie"`
	Movie        MovieDTO        `json:"movie"`
}

// FavoriteMovieWithDetailsDTO represents a favorite movie with full movie details
type FavoriteMovieWithDetailsDTO struct {
	FavoriteMovie FavoriteMovieDTO `json:"favorite_movie"`
	Movie         MovieDTO         `json:"movie"`
}

// ToggleResponse represents the response for toggle operations
type ToggleResponse struct {
	Added   bool   `json:"added"`
	Message string `json:"message"`
}
