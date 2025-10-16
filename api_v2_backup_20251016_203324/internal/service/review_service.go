package service

import (
	"errors"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type reviewService struct {
	reviewRepo domain.ReviewRepository
	movieRepo  domain.MovieRepository
	userRepo   domain.UserRepository
}

func NewReviewService(reviewRepo domain.ReviewRepository, movieRepo domain.MovieRepository, userRepo domain.UserRepository) domain.ReviewService {
	return &reviewService{
		reviewRepo: reviewRepo,
		movieRepo:  movieRepo,
		userRepo:   userRepo,
	}
}

func (s *reviewService) CreateReview(userID, movieID string, rating *int, content *string) (*domain.Review, error) {
	// Check if user exists
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if movie exists
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// Check if user already has a review for this movie
	existingReview, err := s.reviewRepo.GetByUserAndMovie(userID, movieID)
	if err == nil && existingReview != nil {
		return nil, errors.New("user already has a review for this movie")
	}

	// Validate that at least rating or content is provided
	if rating == nil && content == nil {
		return nil, errors.New("review must have at least a rating or content")
	}

	review := &domain.Review{
		UserID:  userID,
		MovieID: movieID,
		Rating:  rating,
		Content: content,
	}

	// Validate review
	if err := s.ValidateReview(review); err != nil {
		return nil, fmt.Errorf("review validation failed: %w", err)
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Get the created review with user and movie data
	createdReview, err := s.reviewRepo.GetByID(review.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created review: %w", err)
	}

	return createdReview, nil
}

func (s *reviewService) GetReview(id string) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return review, nil
}

func (s *reviewService) GetMovieReviews(movieID string, page int) ([]*domain.Review, error) {
	if page < 1 {
		page = 1
	}

	// Check if movie exists
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	limit := 20
	offset := (page - 1) * limit

	reviews, err := s.reviewRepo.GetByMovieID(movieID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) GetUserReviews(userID string, page int) ([]*domain.Review, error) {
	if page < 1 {
		page = 1
	}

	// Check if user exists
	if _, err := s.userRepo.GetByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	limit := 20
	offset := (page - 1) * limit

	reviews, err := s.reviewRepo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reviews: %w", err)
	}

	return reviews, nil
}

func (s *reviewService) UpdateReview(reviewID, userID string, rating *int, content *string) (*domain.Review, error) {
	// Get existing review
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("review not found: %w", err)
	}

	// Check ownership
	if review.UserID != userID {
		return nil, errors.New("unauthorized: user does not own this review")
	}

	// Validate that at least rating or content is provided
	if rating == nil && content == nil {
		return nil, errors.New("review must have at least a rating or content")
	}

	// Update fields
	review.Rating = rating
	review.Content = content

	// Validate updated review
	if err := s.ValidateReview(review); err != nil {
		return nil, fmt.Errorf("review validation failed: %w", err)
	}

	if err := s.reviewRepo.Update(review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Get the updated review with user and movie data
	updatedReview, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated review: %w", err)
	}

	return updatedReview, nil
}

func (s *reviewService) DeleteReview(reviewID, userID string) error {
	// Get existing review
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return fmt.Errorf("review not found: %w", err)
	}

	// Check ownership
	if review.UserID != userID {
		return errors.New("unauthorized: user does not own this review")
	}

	if err := s.reviewRepo.Delete(reviewID); err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

func (s *reviewService) ValidateReview(review *domain.Review) error {
	if review == nil {
		return errors.New("review cannot be nil")
	}

	if review.UserID == "" {
		return errors.New("invalid user ID")
	}

	if review.MovieID == "" {
		return errors.New("invalid movie ID")
	}

	// Validate rating if provided
	if review.Rating != nil {
		if *review.Rating < 1 || *review.Rating > 10 {
			return errors.New("rating must be between 1 and 10")
		}
	}

	// Validate content if provided
	if review.Content != nil {
		contentLen := len(*review.Content)
		if contentLen > 2000 {
			return errors.New("review content cannot exceed 2000 characters")
		}
		if contentLen < 1 {
			return errors.New("review content cannot be empty")
		}
	}

	// At least one of rating or content must be provided
	if review.Rating == nil && review.Content == nil {
		return errors.New("review must have at least a rating or content")
	}

	return nil
}
