package repository

import (
	"database/sql"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type watchedMovieRepository struct {
	db *sqlx.DB
}

func NewWatchedMovieRepository(db *sqlx.DB) domain.WatchedMovieRepository {
	return &watchedMovieRepository{db: db}
}

func (r *watchedMovieRepository) AddWatchedMovie(userID, movieID uuid.UUID) (*domain.WatchedMovie, error) {
	var watched domain.WatchedMovie

	query := `
		INSERT INTO watched_movies (user_id, movie_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, movie_id) 
		DO UPDATE SET watched_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, movie_id, watched_at, created_at
	`

	err := r.db.QueryRowx(query, userID, movieID).StructScan(&watched)
	if err != nil {
		return nil, fmt.Errorf("failed to add watched movie: %w", err)
	}

	return &watched, nil
}

func (r *watchedMovieRepository) RemoveWatchedMovie(userID, movieID uuid.UUID) error {
	query := `
		DELETE FROM watched_movies
		WHERE user_id = $1 AND movie_id = $2
	`

	result, err := r.db.Exec(query, userID, movieID)
	if err != nil {
		return fmt.Errorf("failed to remove watched movie: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *watchedMovieRepository) IsMovieWatched(userID, movieID uuid.UUID) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM watched_movies
			WHERE user_id = $1 AND movie_id = $2
		)
	`

	err := r.db.QueryRow(query, userID, movieID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if movie is watched: %w", err)
	}

	return exists, nil
}

func (r *watchedMovieRepository) GetUserWatchedMovies(userID uuid.UUID) ([]domain.WatchedMovie, error) {
	var watched []domain.WatchedMovie

	query := `
		SELECT id, user_id, movie_id, watched_at, created_at
		FROM watched_movies
		WHERE user_id = $1
		ORDER BY watched_at DESC
	`

	err := r.db.Select(&watched, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user watched movies: %w", err)
	}

	return watched, nil
}
