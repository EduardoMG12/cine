package user_movie

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
)

type RemoveFavoriteMovieUseCase struct {
	favoriteRepo domain.FavoriteMovieRepository
}

func NewRemoveFavoriteMovieUseCase(
	favoriteRepo domain.FavoriteMovieRepository,
) *RemoveFavoriteMovieUseCase {
	return &RemoveFavoriteMovieUseCase{
		favoriteRepo: favoriteRepo,
	}
}

func (uc *RemoveFavoriteMovieUseCase) Execute(userID, movieID uuid.UUID) error {
	err := uc.favoriteRepo.RemoveFavoriteMovie(userID, movieID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("movie not in favorites list")
		}
		return fmt.Errorf("failed to remove favorite movie: %w", err)
	}

	return nil
}
