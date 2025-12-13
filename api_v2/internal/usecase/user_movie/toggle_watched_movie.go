package user_movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type ToggleWatchedMovieUseCase struct {
	watchedRepo domain.WatchedMovieRepository
	movieRepo   domain.MovieRepository
}

func NewToggleWatchedMovieUseCase(
	watchedRepo domain.WatchedMovieRepository,
	movieRepo domain.MovieRepository,
) *ToggleWatchedMovieUseCase {
	return &ToggleWatchedMovieUseCase{
		watchedRepo: watchedRepo,
		movieRepo:   movieRepo,
	}
}

func (uc *ToggleWatchedMovieUseCase) Execute(userID, movieID uuid.UUID) (*dto.ToggleResponse, error) {
	_, err := uc.movieRepo.GetMovieByID(movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found")
	}

	isWatched, err := uc.watchedRepo.IsMovieWatched(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to check watched status: %w", err)
	}

	if isWatched {
		err = uc.watchedRepo.RemoveWatchedMovie(userID, movieID)
		if err != nil {
			return nil, fmt.Errorf("failed to remove from watched list: %w", err)
		}
		return &dto.ToggleResponse{
			Added:   false,
			Message: "Movie removed from watched list",
		}, nil
	}

	_, err = uc.watchedRepo.AddWatchedMovie(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to add to watched list: %w", err)
	}

	return &dto.ToggleResponse{
		Added:   true,
		Message: "Movie added to watched list",
	}, nil
}
