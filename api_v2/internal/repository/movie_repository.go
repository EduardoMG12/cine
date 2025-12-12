package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type movieRepository struct {
	db *sqlx.DB
}

func NewMovieRepository(db *sqlx.DB) domain.MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) CreateMovie(movie *domain.Movie) error {
	movie.ID = uuid.New()
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	query := `
		INSERT INTO movies (
			id, external_api_id, title, overview, release_date, poster_url, 
			backdrop_url, genres, runtime, vote_average, vote_count, adult, 
			cache_expires_at, created_at, updated_at
		) VALUES (
			:id, :external_api_id, :title, :overview, :release_date, :poster_url,
			:backdrop_url, :genres, :runtime, :vote_average, :vote_count, :adult,
			:cache_expires_at, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExec(query, movie)
	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}

	return nil
}

func (r *movieRepository) GetMovieByID(id uuid.UUID) (*domain.Movie, error) {
	var movie domain.Movie
	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url,
			   backdrop_url, genres, runtime, vote_average, vote_count, adult,
			   cache_expires_at, created_at, updated_at
		FROM movies 
		WHERE id = $1
	`

	err := r.db.Get(&movie, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to get movie by id: %w", err)
	}

	return &movie, nil
}

func (r *movieRepository) GetMovieByExternalID(externalID string) (*domain.Movie, error) {
	var movie domain.Movie
	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url,
			   backdrop_url, genres, runtime, vote_average, vote_count, adult,
			   cache_expires_at, created_at, updated_at
		FROM movies 
		WHERE external_api_id = $1
	`

	err := r.db.Get(&movie, query, externalID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to get movie by external id: %w", err)
	}

	return &movie, nil
}

func (r *movieRepository) UpdateMovie(movie *domain.Movie) error {
	movie.UpdatedAt = time.Now()

	query := `
		UPDATE movies SET
			title = :title,
			overview = :overview,
			release_date = :release_date,
			poster_url = :poster_url,
			backdrop_url = :backdrop_url,
			genres = :genres,
			runtime = :runtime,
			vote_average = :vote_average,
			vote_count = :vote_count,
			adult = :adult,
			cache_expires_at = :cache_expires_at,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExec(query, movie)
	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	return nil
}

func (r *movieRepository) DeleteMovie(id uuid.UUID) error {
	query := `DELETE FROM movies WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("movie not found")
	}

	return nil
}

func (r *movieRepository) GetRandomMovie() (*domain.Movie, error) {
	var movie domain.Movie
	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url,
			   backdrop_url, genres, runtime, vote_average, vote_count, adult,
			   cache_expires_at, created_at, updated_at
		FROM movies 
		WHERE cache_expires_at > NOW()
		ORDER BY RANDOM()
		LIMIT 1
	`

	err := r.db.Get(&movie, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no movies found")
		}
		return nil, fmt.Errorf("failed to get random movie: %w", err)
	}

	return &movie, nil
}

func (r *movieRepository) GetRandomMovieByGenre(genre string) (*domain.Movie, error) {
	var movie domain.Movie
	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url,
			   backdrop_url, genres, runtime, vote_average, vote_count, adult,
			   cache_expires_at, created_at, updated_at
		FROM movies 
		WHERE cache_expires_at > NOW()
		  AND $1 = ANY(genres)
		ORDER BY RANDOM()
		LIMIT 1
	`

	err := r.db.Get(&movie, query, genre)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no movies found for genre: %s", genre)
		}
		return nil, fmt.Errorf("failed to get random movie by genre: %w", err)
	}

	return &movie, nil
}

func (r *movieRepository) SearchMovies(queryText string, limit int) ([]*domain.Movie, error) {
	var movies []*domain.Movie

	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url,
			   backdrop_url, genres, runtime, vote_average, vote_count, adult,
			   cache_expires_at, created_at, updated_at
		FROM movies 
		WHERE cache_expires_at > NOW()
		  AND (
			  title ILIKE $1
			  OR overview ILIKE $1
		  )
		ORDER BY vote_average DESC NULLS LAST, vote_count DESC NULLS LAST
		LIMIT $2
	`

	searchPattern := "%" + queryText + "%"
	err := r.db.Select(&movies, query, searchPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	return movies, nil
}

// GetRandomMovies returns N random movies from the database
func (r *movieRepository) GetRandomMovies(limit int) ([]*domain.Movie, error) {
	var movies []*domain.Movie

	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url,
			   backdrop_url, genres, runtime, vote_average, vote_count, adult,
			   cache_expires_at, created_at, updated_at
		FROM movies
		ORDER BY RANDOM()
		LIMIT $1
	`

	err := r.db.Select(&movies, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get random movies: %w", err)
	}

	return movies, nil
}

// CountMovies returns the total number of valid (non-expired) movies in the database
func (r *movieRepository) CountMovies() (int, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM movies
	`

	err := r.db.Get(&count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count movies: %w", err)
	}

	return count, nil
}

// Helper function to convert []string to pq.StringArray for database operations
func StringSliceToArray(s []string) pq.StringArray {
	return pq.StringArray(s)
}
