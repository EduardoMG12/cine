package movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
)

type GetRandomMovieByGenreUseCase struct {
	movieRepo domain.MovieRepository
}

func NewGetRandomMovieByGenreUseCase(movieRepo domain.MovieRepository) *GetRandomMovieByGenreUseCase {
	return &GetRandomMovieByGenreUseCase{
		movieRepo: movieRepo,
	}
}

func (uc *GetRandomMovieByGenreUseCase) Execute(genre string) (*dto.MovieDTO, error) {
	movie, err := uc.movieRepo.GetRandomMovieByGenre(genre)
	if err != nil {
		return nil, fmt.Errorf("failed to get random movie by genre: %w", err)
	}

	return uc.movieToDTO(movie), nil
}

func (uc *GetRandomMovieByGenreUseCase) movieToDTO(movie *domain.Movie) *dto.MovieDTO {
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
