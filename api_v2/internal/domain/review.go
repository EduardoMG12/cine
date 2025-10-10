package domain

import (
	"time"
)

type Review struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	MovieID   string    `json:"movie_id" db:"movie_id"`
	Rating    *int      `json:"rating,omitempty" db:"rating"` // 1-10 scale
	Content   *string   `json:"content,omitempty" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	User  *User  `json:"user,omitempty"`
	Movie *Movie `json:"movie,omitempty"`
}

type ReviewRepository interface {
	Create(review *Review) error
	GetByID(id string) (*Review, error)
	GetByUserAndMovie(userID, movieID string) (*Review, error)
	GetByMovieID(movieID string, limit, offset int) ([]*Review, error)
	GetByUserID(userID string, limit, offset int) ([]*Review, error)
	Update(review *Review) error
	Delete(id string) error
}

type ReviewService interface {
	CreateReview(userID, movieID string, rating *int, content *string) (*Review, error)
	GetReview(id string) (*Review, error)
	GetMovieReviews(movieID string, page int) ([]*Review, error)
	GetUserReviews(userID string, page int) ([]*Review, error)
	UpdateReview(reviewID, userID string, rating *int, content *string) (*Review, error)
	DeleteReview(reviewID, userID string) error
	ValidateReview(review *Review) error
}
