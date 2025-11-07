package movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

// GetPopularMoviesUseCase busca filmes populares
type GetPopularMoviesUseCase struct {
	tmdbService *infrastructure.TMDbService
}

func NewGetPopularMoviesUseCase(tmdbService *infrastructure.TMDbService) *GetPopularMoviesUseCase {
	return &GetPopularMoviesUseCase{
		tmdbService: tmdbService,
	}
}

func (uc *GetPopularMoviesUseCase) Execute(page int) (*dto.TMDbDiscoverResponse, error) {
	if page < 1 {
		page = 1
	}

	response, err := uc.tmdbService.GetPopularMovies(page)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular movies: %w", err)
	}

	// Processar URLs das imagens
	for i := range response.Results {
		if response.Results[i].PosterPath != "" {
			response.Results[i].PosterPath = uc.tmdbService.GetImageURL(response.Results[i].PosterPath)
		}
		if response.Results[i].BackdropPath != "" {
			response.Results[i].BackdropPath = uc.tmdbService.GetImageURL(response.Results[i].BackdropPath)
		}
	}

	return response, nil
}

// GetTrendingMoviesUseCase busca filmes em alta
type GetTrendingMoviesUseCase struct {
	tmdbService *infrastructure.TMDbService
}

func NewGetTrendingMoviesUseCase(tmdbService *infrastructure.TMDbService) *GetTrendingMoviesUseCase {
	return &GetTrendingMoviesUseCase{
		tmdbService: tmdbService,
	}
}

func (uc *GetTrendingMoviesUseCase) Execute(timeWindow string) (*dto.TMDbDiscoverResponse, error) {
	if timeWindow != "day" && timeWindow != "week" {
		timeWindow = "week"
	}

	response, err := uc.tmdbService.GetTrendingMovies(timeWindow)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}

	// Processar URLs das imagens
	for i := range response.Results {
		if response.Results[i].PosterPath != "" {
			response.Results[i].PosterPath = uc.tmdbService.GetImageURL(response.Results[i].PosterPath)
		}
		if response.Results[i].BackdropPath != "" {
			response.Results[i].BackdropPath = uc.tmdbService.GetImageURL(response.Results[i].BackdropPath)
		}
	}

	return response, nil
}

// GetGenresUseCase busca lista de gêneros
type GetGenresUseCase struct {
	tmdbService *infrastructure.TMDbService
}

func NewGetGenresUseCase(tmdbService *infrastructure.TMDbService) *GetGenresUseCase {
	return &GetGenresUseCase{
		tmdbService: tmdbService,
	}
}

func (uc *GetGenresUseCase) Execute() (*dto.TMDbGenresResponse, error) {
	response, err := uc.tmdbService.GetGenres()
	if err != nil {
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}

	return response, nil
}
