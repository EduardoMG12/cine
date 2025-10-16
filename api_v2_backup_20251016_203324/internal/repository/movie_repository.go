package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// Helper functions for UUID operations
func generateUUID() string {
	return utils.GenerateUUID()
}

func isValidUUID(uuid string) bool {
	return utils.IsValidUUID(uuid)
}

type movieRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewMovieRepository(db *sqlx.DB, redis *redis.Client) domain.MovieRepository {
	return &movieRepository{
		db:    db,
		redis: redis,
	}
}

func (r *movieRepository) Create(movie *domain.Movie) error {
	// Generate UUID for the movie if not set
	if movie.ID == "" {
		movie.ID = generateUUID()
	}

	query := `
		INSERT INTO movies (id, external_api_id, title, overview, release_date, poster_url, backdrop_url, 
		                   genres, runtime, vote_average, vote_count, adult, cache_expires_at, 
		                   created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	now := time.Now()
	movie.CreatedAt = now
	movie.UpdatedAt = now

	_, err := r.db.Exec(
		query,
		movie.ID,
		movie.ExternalAPIID,
		movie.Title,
		movie.Overview,
		movie.ReleaseDate,
		movie.PosterURL,
		movie.BackdropURL,
		pq.Array(movie.Genres),
		movie.Runtime,
		movie.VoteAverage,
		movie.VoteCount,
		movie.Adult,
		movie.CacheExpiresAt,
		movie.CreatedAt,
		movie.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create movie: %w", err)
	}

	r.cacheMovie(movie)
	return nil
}

func (r *movieRepository) GetByID(id string) (*domain.Movie, error) {
	// Validate UUID format
	if !isValidUUID(id) {
		return nil, fmt.Errorf("invalid UUID format: %s", id)
	}

	cacheKey := fmt.Sprintf("movie:id:%s", id)
	if movie := r.getMovieFromCache(cacheKey); movie != nil {
		return movie, nil
	}

	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url, backdrop_url,
		       genres, runtime, vote_average, vote_count, adult, cache_expires_at,
		       created_at, updated_at
		FROM movies 
		WHERE id = $1
	`

	movie := &domain.Movie{}
	err := r.db.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.ExternalAPIID,
		&movie.Title,
		&movie.Overview,
		&movie.ReleaseDate,
		&movie.PosterURL,
		&movie.BackdropURL,
		pq.Array(&movie.Genres),
		&movie.Runtime,
		&movie.VoteAverage,
		&movie.VoteCount,
		&movie.Adult,
		&movie.CacheExpiresAt,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get movie by ID: %w", err)
	}

	r.cacheMovie(movie)
	return movie, nil
}

func (r *movieRepository) GetByExternalID(externalID string) (*domain.Movie, error) {
	cacheKey := fmt.Sprintf("movie:external:%s", externalID)
	if movie := r.getMovieFromCache(cacheKey); movie != nil {
		return movie, nil
	}

	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url, backdrop_url,
		       genres, runtime, vote_average, vote_count, adult, cache_expires_at,
		       created_at, updated_at
		FROM movies 
		WHERE external_api_id = $1
	`

	movie := &domain.Movie{}
	err := r.db.QueryRow(query, externalID).Scan(
		&movie.ID,
		&movie.ExternalAPIID,
		&movie.Title,
		&movie.Overview,
		&movie.ReleaseDate,
		&movie.PosterURL,
		&movie.BackdropURL,
		pq.Array(&movie.Genres),
		&movie.Runtime,
		&movie.VoteAverage,
		&movie.VoteCount,
		&movie.Adult,
		&movie.CacheExpiresAt,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get movie by external ID: %w", err)
	}

	r.cacheMovie(movie)
	return movie, nil
}

func (r *movieRepository) Update(movie *domain.Movie) error {
	query := `
		UPDATE movies 
		SET title = $2, overview = $3, release_date = $4, poster_url = $5, backdrop_url = $6,
		    genres = $7, runtime = $8, vote_average = $9, vote_count = $10, adult = $11,
		    cache_expires_at = $12, updated_at = $13
		WHERE id = $1
	`

	movie.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		movie.ID,
		movie.Title,
		movie.Overview,
		movie.ReleaseDate,
		movie.PosterURL,
		movie.BackdropURL,
		pq.Array(movie.Genres),
		movie.Runtime,
		movie.VoteAverage,
		movie.VoteCount,
		movie.Adult,
		movie.CacheExpiresAt,
		movie.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	r.invalidateMovieCache(movie.ID, movie.ExternalAPIID)
	r.cacheMovie(movie)
	return nil
}

func (r *movieRepository) Delete(id string) error {
	// Validate UUID format
	if !isValidUUID(id) {
		return fmt.Errorf("invalid UUID format: %s", id)
	}

	movie, err := r.GetByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM movies WHERE id = $1`
	_, err = r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	r.invalidateMovieCache(movie.ID, movie.ExternalAPIID)
	return nil
}

func (r *movieRepository) Search(query string, limit, offset int) ([]*domain.Movie, error) {
	sqlQuery := `
		SELECT id, external_api_id, title, overview, release_date, poster_url, backdrop_url,
		       genres, runtime, vote_average, vote_count, adult, cache_expires_at,
		       created_at, updated_at
		FROM movies 
		WHERE to_tsvector('english', title) @@ plainto_tsquery('english', $1)
		   OR title ILIKE '%' || $1 || '%'
		ORDER BY vote_average DESC NULLS LAST, vote_count DESC NULLS LAST
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Queryx(sqlQuery, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	return r.scanMovies(rows)
}

func (r *movieRepository) GetPopular(limit, offset int) ([]*domain.Movie, error) {
	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url, backdrop_url,
		       genres, runtime, vote_average, vote_count, adult, cache_expires_at,
		       created_at, updated_at
		FROM movies 
		WHERE vote_count > 100
		ORDER BY vote_average DESC, vote_count DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Queryx(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}
	defer rows.Close()

	return r.scanMovies(rows)
}

func (r *movieRepository) GetByGenre(genre string, limit, offset int) ([]*domain.Movie, error) {
	query := `
		SELECT id, external_api_id, title, overview, release_date, poster_url, backdrop_url,
		       genres, runtime, vote_average, vote_count, adult, cache_expires_at,
		       created_at, updated_at
		FROM movies 
		WHERE $1 = ANY(genres)
		ORDER BY vote_average DESC NULLS LAST, vote_count DESC NULLS LAST
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Queryx(query, genre, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies by genre: %w", err)
	}
	defer rows.Close()

	return r.scanMovies(rows)
}

func (r *movieRepository) DeleteExpiredCache() error {
	query := `DELETE FROM movies WHERE cache_expires_at < NOW()`
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired cache: %w", err)
	}

	// Also clean up Redis cache
	ctx := context.Background()
	keys, err := r.redis.Keys(ctx, "movie:*").Result()
	if err == nil {
		for _, key := range keys {
			r.redis.Del(ctx, key)
		}
	}

	return nil
}

// Helper methods

func (r *movieRepository) scanMovies(rows *sqlx.Rows) ([]*domain.Movie, error) {
	var movies []*domain.Movie

	for rows.Next() {
		movie := &domain.Movie{}
		err := rows.Scan(
			&movie.ID,
			&movie.ExternalAPIID,
			&movie.Title,
			&movie.Overview,
			&movie.ReleaseDate,
			&movie.PosterURL,
			&movie.BackdropURL,
			pq.Array(&movie.Genres),
			&movie.Runtime,
			&movie.VoteAverage,
			&movie.VoteCount,
			&movie.Adult,
			&movie.CacheExpiresAt,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		movies = append(movies, movie)
	}

	return movies, nil
}

func (r *movieRepository) cacheMovie(movie *domain.Movie) {
	ctx := context.Background()

	// Cache by ID
	idKey := fmt.Sprintf("movie:id:%s", movie.ID)
	r.redis.Set(ctx, idKey, movie, 30*time.Minute)

	// Cache by external ID
	externalKey := fmt.Sprintf("movie:external:%s", movie.ExternalAPIID)
	r.redis.Set(ctx, externalKey, movie, 30*time.Minute)
}

func (r *movieRepository) getMovieFromCache(key string) *domain.Movie {
	ctx := context.Background()

	var movie domain.Movie
	err := r.redis.Get(ctx, key).Scan(&movie)
	if err != nil {
		return nil
	}

	return &movie
}

func (r *movieRepository) invalidateMovieCache(movieID string, externalID string) {
	ctx := context.Background()

	idKey := fmt.Sprintf("movie:id:%s", movieID)
	externalKey := fmt.Sprintf("movie:external:%s", externalID)

	r.redis.Del(ctx, idKey, externalKey)
}
