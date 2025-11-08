package movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

type SearchMoviesUseCase struct {
	movieFetcher infrastructure.MovieFetcher
}

func NewSearchMoviesUseCase(movieFetcher infrastructure.MovieFetcher) *SearchMoviesUseCase {
	return &SearchMoviesUseCase{
		movieFetcher: movieFetcher,
	}
}

func (uc *SearchMoviesUseCase) Execute(query string, page int) ([]*dto.MovieDTO, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	movies, err := uc.movieFetcher.Search(query, page)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	dtos := make([]*dto.MovieDTO, len(movies))
	for i, movie := range movies {
		dtos[i] = uc.movieToDTO(movie)
	}

	return dtos, nil
}

func (uc *SearchMoviesUseCase) movieToDTO(movie *domain.Movie) *dto.MovieDTO {
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
