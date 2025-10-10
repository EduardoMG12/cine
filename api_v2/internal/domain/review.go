package domain

import (
	"time"
)

type Review struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	MovieID   int       `json:"movie_id" db:"movie_id"`
	Rating    *int      `json:"rating,omitempty" db:"rating"` // 1-10 scale
	Content   *string   `json:"content,omitempty" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Populated by joins
	User  *User  `json:"user,omitempty"`
	Movie *Movie `json:"movie,omitempty"`
}

type ReviewRepository interface {
	Create(review *Review) error
	GetByID(id int) (*Review, error)
	GetByUserAndMovie(userID, movieID int) (*Review, error)
	GetByMovieID(movieID int, limit, offset int) ([]*Review, error)
	GetByUserID(userID int, limit, offset int) ([]*Review, error)
	Update(review *Review) error
	Delete(id int) error
}

type ReviewService interface {
	CreateReview(userID, movieID int, rating *int, content *string) (*Review, error)
	GetReview(id int) (*Review, error)
	GetMovieReviews(movieID int, page int) ([]*Review, error)
	GetUserReviews(userID int, page int) ([]*Review, error)
	UpdateReview(reviewID, userID int, rating *int, content *string) (*Review, error)
	DeleteReview(reviewID, userID int) error
	ValidateReview(review *Review) error
}
