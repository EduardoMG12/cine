package user_movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type GetFavoriteMoviesUseCase struct {
	favoriteRepo domain.FavoriteMovieRepository
	movieRepo    domain.MovieRepository
}

func NewGetFavoriteMoviesUseCase(
	favoriteRepo domain.FavoriteMovieRepository,
	movieRepo domain.MovieRepository,
) *GetFavoriteMoviesUseCase {
	return &GetFavoriteMoviesUseCase{
		favoriteRepo: favoriteRepo,
		movieRepo:    movieRepo,
	}
}

func (uc *GetFavoriteMoviesUseCase) Execute(userID uuid.UUID) ([]dto.FavoriteMovieWithDetailsDTO, error) {
	favoriteMovies, err := uc.favoriteRepo.GetUserFavoriteMovies(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get favorite movies: %w", err)
	}

	var result []dto.FavoriteMovieWithDetailsDTO
	for _, favorite := range favoriteMovies {
		movie, err := uc.movieRepo.GetMovieByID(favorite.MovieID)
		if err != nil {
			continue
		}

		result = append(result, dto.FavoriteMovieWithDetailsDTO{
			FavoriteMovie: dto.FavoriteMovieDTO{
				ID:          favorite.ID,
				UserID:      favorite.UserID,
				MovieID:     favorite.MovieID,
				FavoritedAt: favorite.FavoritedAt,
				CreatedAt:   favorite.CreatedAt,
			},
			Movie: dto.MovieDTO{
				ID:            movie.ID,
				ExternalAPIID: movie.ExternalAPIID,
				Title:         movie.Title,
				Overview:      movie.Overview,
				PosterURL:     movie.PosterURL,
				Genres:        movie.Genres,
				Adult:         movie.Adult,
				CreatedAt:     movie.CreatedAt,
				UpdatedAt:     movie.UpdatedAt,
			},
		})
	}

	return result, nil
}
