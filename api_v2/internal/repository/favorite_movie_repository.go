package repository

import (
	"database/sql"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type favoriteMovieRepository struct {
	db *sqlx.DB
}

func NewFavoriteMovieRepository(db *sqlx.DB) domain.FavoriteMovieRepository {
	return &favoriteMovieRepository{db: db}
}

func (r *favoriteMovieRepository) AddFavoriteMovie(userID, movieID uuid.UUID) (*domain.FavoriteMovie, error) {
	var favorite domain.FavoriteMovie

	query := `
		INSERT INTO favorite_movies (user_id, movie_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, movie_id) 
		DO UPDATE SET favorited_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, movie_id, favorited_at, created_at
	`

	err := r.db.QueryRowx(query, userID, movieID).StructScan(&favorite)
	if err != nil {
		return nil, fmt.Errorf("failed to add favorite movie: %w", err)
	}

	return &favorite, nil
}

func (r *favoriteMovieRepository) RemoveFavoriteMovie(userID, movieID uuid.UUID) error {
	query := `
		DELETE FROM favorite_movies
		WHERE user_id = $1 AND movie_id = $2
	`

	result, err := r.db.Exec(query, userID, movieID)
	if err != nil {
		return fmt.Errorf("failed to remove favorite movie: %w", err)
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

func (r *favoriteMovieRepository) IsMovieFavorite(userID, movieID uuid.UUID) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM favorite_movies
			WHERE user_id = $1 AND movie_id = $2
		)
	`

	err := r.db.QueryRow(query, userID, movieID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if movie is favorite: %w", err)
	}

	return exists, nil
}

func (r *favoriteMovieRepository) GetUserFavoriteMovies(userID uuid.UUID) ([]domain.FavoriteMovie, error) {
	var favorites []domain.FavoriteMovie

	query := `
		SELECT id, user_id, movie_id, favorited_at, created_at
		FROM favorite_movies
		WHERE user_id = $1
		ORDER BY favorited_at DESC
	`

	err := r.db.Select(&favorites, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user favorite movies: %w", err)
	}

	return favorites, nil
}
