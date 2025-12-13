package user_movie

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type AddFavoriteMovieUseCase struct {
	favoriteRepo domain.FavoriteMovieRepository
	movieRepo    domain.MovieRepository
}

func NewAddFavoriteMovieUseCase(
	favoriteRepo domain.FavoriteMovieRepository,
	movieRepo domain.MovieRepository,
) *AddFavoriteMovieUseCase {
	return &AddFavoriteMovieUseCase{
		favoriteRepo: favoriteRepo,
		movieRepo:    movieRepo,
	}
}

func (uc *AddFavoriteMovieUseCase) Execute(userID, movieID uuid.UUID) (*dto.FavoriteMovieDTO, error) {
	_, err := uc.movieRepo.GetMovieByID(movieID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("movie not found")
		}
		return nil, fmt.Errorf("failed to verify movie: %w", err)
	}

	favorite, err := uc.favoriteRepo.AddFavoriteMovie(userID, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to add favorite movie: %w", err)
	}

	return &dto.FavoriteMovieDTO{
		ID:          favorite.ID,
		UserID:      favorite.UserID,
		MovieID:     favorite.MovieID,
		FavoritedAt: favorite.FavoritedAt,
		CreatedAt:   favorite.CreatedAt,
	}, nil
}
