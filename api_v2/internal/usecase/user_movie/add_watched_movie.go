package user_movie

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type AddWatchedMovieUseCase struct {
	watchedRepo domain.WatchedMovieRepository
	movieRepo   domain.MovieRepository
}

func NewAddWatchedMovieUseCase(
	watchedRepo domain.WatchedMovieRepository,
	movieRepo domain.MovieRepository,
) *AddWatchedMovieUseCase {
	return &AddWatchedMovieUseCase{
		watchedRepo: watchedRepo,
		movieRepo:   movieRepo,
	}
}

func (uc *AddWatchedMovieUseCase) Execute(userID, movieID uuid.UUID) (*dto.WatchedMovieDTO, error) {
	_, err := uc.movieRepo.GetMovieByID(movieID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to verify movie: %w", err)
	}

	watched, err := uc.watchedRepo.AddWatchedMovie(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to add watched movie: %w", err)
	}

	return &dto.WatchedMovieDTO{
		ID:        watched.ID,
		UserID:    watched.UserID,
		MovieID:   watched.MovieID,
		WatchedAt: watched.WatchedAt,
		CreatedAt: watched.CreatedAt,
	}, nil
}
