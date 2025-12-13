package user_movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type ToggleFavoriteMovieUseCase struct {
	favoriteRepo domain.FavoriteMovieRepository
	movieRepo    domain.MovieRepository
}

func NewToggleFavoriteMovieUseCase(
	favoriteRepo domain.FavoriteMovieRepository,
	movieRepo domain.MovieRepository,
) *ToggleFavoriteMovieUseCase {
	return &ToggleFavoriteMovieUseCase{
		favoriteRepo: favoriteRepo,
		movieRepo:    movieRepo,
	}
}

func (uc *ToggleFavoriteMovieUseCase) Execute(userID, movieID uuid.UUID) (*dto.ToggleResponse, error) {
	_, err := uc.movieRepo.GetMovieByID(movieID)
	if err != nil {
		return nil, fmt.Errorf("movie not found")
	}

	isFavorite, err := uc.favoriteRepo.IsMovieFavorite(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to check favorite status: %w", err)
	}

	if isFavorite {
		err = uc.favoriteRepo.RemoveFavoriteMovie(userID, movieID)
		if err != nil {
			return nil, fmt.Errorf("failed to remove from favorites: %w", err)
		}
		return &dto.ToggleResponse{
			Added:   false,
			Message: "Movie removed from favorites",
		}, nil
	}

	_, err = uc.favoriteRepo.AddFavoriteMovie(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to add to favorites: %w", err)
	}

	return &dto.ToggleResponse{
		Added:   true,
		Message: "Movie added to favorites",
	}, nil
}
