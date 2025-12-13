package user_movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type GetWatchedMoviesUseCase struct {
	watchedRepo domain.WatchedMovieRepository
	movieRepo   domain.MovieRepository
}

func NewGetWatchedMoviesUseCase(
	watchedRepo domain.WatchedMovieRepository,
	movieRepo domain.MovieRepository,
) *GetWatchedMoviesUseCase {
	return &GetWatchedMoviesUseCase{
		watchedRepo: watchedRepo,
		movieRepo:   movieRepo,
	}
}

func (uc *GetWatchedMoviesUseCase) Execute(userID uuid.UUID) ([]dto.WatchedMovieWithDetailsDTO, error) {
	watchedMovies, err := uc.watchedRepo.GetUserWatchedMovies(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get watched movies: %w", err)
	}

	var result []dto.WatchedMovieWithDetailsDTO
	for _, watched := range watchedMovies {
		movie, err := uc.movieRepo.GetMovieByID(watched.MovieID)
		if err != nil {
			continue
		}

		result = append(result, dto.WatchedMovieWithDetailsDTO{
			WatchedMovie: dto.WatchedMovieDTO{
				ID:        watched.ID,
				UserID:    watched.UserID,
				MovieID:   watched.MovieID,
				WatchedAt: watched.WatchedAt,
				CreatedAt: watched.CreatedAt,
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
