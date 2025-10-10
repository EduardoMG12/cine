package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/jmoiron/sqlx"
)

type reviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) domain.ReviewRepository {
	return &reviewRepository{
		db: db,
	}
}

func (r *reviewRepository) Create(review *domain.Review) error {
	query := `
		INSERT INTO reviews (user_id, movie_id, rating, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		review.UserID,
		review.MovieID,
		review.Rating,
		review.Content,
		review.CreatedAt,
		review.UpdatedAt,
	).Scan(&review.ID)

	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func (r *reviewRepository) GetByID(id int) (*domain.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.movie_id, r.rating, r.content, r.created_at, r.updated_at,
		       u.id, u.username, u.email, u.display_name, u.bio, u.profile_picture_url, u.is_private, u.created_at,
		       m.id, m.external_api_id, m.title, m.overview, m.poster_url, m.backdrop_url,
		       m.release_date, m.genres, m.vote_average, m.runtime, m.vote_count, m.adult,
		       m.created_at, m.updated_at
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		LEFT JOIN movies m ON r.movie_id = m.id
		WHERE r.id = $1`

	var review domain.Review
	var user domain.User
	var movie domain.Movie

	err := r.db.QueryRow(query, id).Scan(
		&review.ID, &review.UserID, &review.MovieID, &review.Rating, &review.Content, &review.CreatedAt, &review.UpdatedAt,
		&user.ID, &user.Username, &user.Email, &user.DisplayName, &user.Bio, &user.ProfilePictureURL, &user.IsPrivate, &user.CreatedAt,
		&movie.ID, &movie.ExternalAPIID, &movie.Title, &movie.Overview, &movie.PosterURL, &movie.BackdropURL,
		&movie.ReleaseDate, &movie.Genres, &movie.VoteAverage, &movie.Runtime, &movie.VoteCount, &movie.Adult,
		&movie.CreatedAt, &movie.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	review.User = &user
	review.Movie = &movie
	return &review, nil
}

func (r *reviewRepository) GetByUserAndMovie(userID, movieID int) (*domain.Review, error) {
	query := `
		SELECT id, user_id, movie_id, rating, content, created_at, updated_at
		FROM reviews
		WHERE user_id = $1 AND movie_id = $2`

	var review domain.Review
	err := r.db.QueryRow(query, userID, movieID).Scan(
		&review.ID,
		&review.UserID,
		&review.MovieID,
		&review.Rating,
		&review.Content,
		&review.CreatedAt,
		&review.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

func (r *reviewRepository) GetByMovieID(movieID int, limit, offset int) ([]*domain.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.movie_id, r.rating, r.content, r.created_at, r.updated_at,
		       u.id, u.username, u.display_name, u.bio, u.profile_picture_url, u.is_private, u.created_at
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.movie_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, movieID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var review domain.Review
		var user domain.User

		err := rows.Scan(
			&review.ID, &review.UserID, &review.MovieID, &review.Rating, &review.Content, &review.CreatedAt, &review.UpdatedAt,
			&user.ID, &user.Username, &user.DisplayName, &user.Bio, &user.ProfilePictureURL, &user.IsPrivate, &user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}

		review.User = &user
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *reviewRepository) GetByUserID(userID int, limit, offset int) ([]*domain.Review, error) {
	query := `
		SELECT r.id, r.user_id, r.movie_id, r.rating, r.content, r.created_at, r.updated_at,
		       m.id, m.external_api_id, m.title, m.overview, m.poster_url, m.backdrop_url,
		       m.release_date, m.genres, m.vote_average, m.runtime, m.vote_count, m.adult,
		       m.created_at, m.updated_at
		FROM reviews r
		LEFT JOIN movies m ON r.movie_id = m.id
		WHERE r.user_id = $1
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		var review domain.Review
		var movie domain.Movie

		err := rows.Scan(
			&review.ID, &review.UserID, &review.MovieID, &review.Rating, &review.Content, &review.CreatedAt, &review.UpdatedAt,
			&movie.ID, &movie.ExternalAPIID, &movie.Title, &movie.Overview, &movie.PosterURL, &movie.BackdropURL,
			&movie.ReleaseDate, &movie.Genres, &movie.VoteAverage, &movie.Runtime, &movie.VoteCount, &movie.Adult,
			&movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}

		review.Movie = &movie
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *reviewRepository) Update(review *domain.Review) error {
	query := `
		UPDATE reviews
		SET rating = $1, content = $2, updated_at = $3
		WHERE id = $4`

	review.UpdatedAt = time.Now()

	_, err := r.db.Exec(query, review.Rating, review.Content, review.UpdatedAt, review.ID)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

func (r *reviewRepository) Delete(id int) error {
	query := `DELETE FROM reviews WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("review not found")
	}

	return nil
}
