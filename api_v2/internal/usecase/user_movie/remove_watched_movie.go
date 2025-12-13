package user_movie

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
)

type RemoveWatchedMovieUseCase struct {
	watchedRepo domain.WatchedMovieRepository
}

func NewRemoveWatchedMovieUseCase(
	watchedRepo domain.WatchedMovieRepository,
) *RemoveWatchedMovieUseCase {
	return &RemoveWatchedMovieUseCase{
		watchedRepo: watchedRepo,
	}
}

func (uc *RemoveWatchedMovieUseCase) Execute(userID, movieID uuid.UUID) error {
	err := uc.watchedRepo.RemoveWatchedMovie(userID, movieID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("movie not in watched list")
		}
		return fmt.Errorf("failed to remove watched movie: %w", err)
	}

	return nil
}
