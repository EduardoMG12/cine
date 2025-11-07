package movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
)

type GetRandomMovieUseCase struct {
	movieRepo domain.MovieRepository
}

func NewGetRandomMovieUseCase(movieRepo domain.MovieRepository) *GetRandomMovieUseCase {
	return &GetRandomMovieUseCase{
		movieRepo: movieRepo,
	}
}

func (uc *GetRandomMovieUseCase) Execute() (*dto.MovieDTO, error) {
	movie, err := uc.movieRepo.GetRandomMovie()
	if err != nil {
		return nil, fmt.Errorf("failed to get random movie: %w", err)
	}

	return uc.movieToDTO(movie), nil
}

func (uc *GetRandomMovieUseCase) movieToDTO(movie *domain.Movie) *dto.MovieDTO {
	return &dto.MovieDTO{
		ID:            movie.ID,
		ExternalAPIID: movie.ExternalAPIID,
		Title:         movie.Title,
		Overview:      movie.Overview,
		ReleaseDate:   movie.ReleaseDate,
		PosterURL:     movie.PosterURL,
		BackdropURL:   movie.BackdropURL,
		Genres:        movie.Genres,
		Runtime:       movie.Runtime,
		VoteAverage:   movie.VoteAverage,
		VoteCount:     movie.VoteCount,
		Adult:         movie.Adult,
		CreatedAt:     movie.CreatedAt,
		UpdatedAt:     movie.UpdatedAt,
	}
}
