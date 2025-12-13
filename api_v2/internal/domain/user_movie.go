package domain

import (
	"time"

	"github.com/google/uuid"
)

// WatchedMovie represents a movie that a user has watched
type WatchedMovie struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	MovieID   uuid.UUID `db:"movie_id"`
	WatchedAt time.Time `db:"watched_at"`
	CreatedAt time.Time `db:"created_at"`
}

// FavoriteMovie represents a movie that a user has marked as favorite
type FavoriteMovie struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	MovieID     uuid.UUID `db:"movie_id"`
	FavoritedAt time.Time `db:"favorited_at"`
	CreatedAt   time.Time `db:"created_at"`
}

// WatchedMovieRepository interface for watched movies operations
type WatchedMovieRepository interface {
	AddWatchedMovie(userID, movieID uuid.UUID) (*WatchedMovie, error)
	RemoveWatchedMovie(userID, movieID uuid.UUID) error
	IsMovieWatched(userID, movieID uuid.UUID) (bool, error)
	GetUserWatchedMovies(userID uuid.UUID) ([]WatchedMovie, error)
}

// FavoriteMovieRepository interface for favorite movies operations
type FavoriteMovieRepository interface {
	AddFavoriteMovie(userID, movieID uuid.UUID) (*FavoriteMovie, error)
	RemoveFavoriteMovie(userID, movieID uuid.UUID) error
	IsMovieFavorite(userID, movieID uuid.UUID) (bool, error)
	GetUserFavoriteMovies(userID uuid.UUID) ([]FavoriteMovie, error)
}
