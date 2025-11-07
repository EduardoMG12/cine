package movie

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

type SearchMoviesUseCase struct {
	tmdbService *infrastructure.TMDbService
}

func NewSearchMoviesUseCase(tmdbService *infrastructure.TMDbService) *SearchMoviesUseCase {
	return &SearchMoviesUseCase{
		tmdbService: tmdbService,
	}
}

func (uc *SearchMoviesUseCase) Execute(query string, page int) (*dto.TMDbSearchResponse, error) {
	if page < 1 {
		page = 1
	}

	searchResp, err := uc.tmdbService.SearchMovies(query, page)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}

	// Processar URLs das imagens
	for i := range searchResp.Results {
		if searchResp.Results[i].PosterPath != "" {
			searchResp.Results[i].PosterPath = uc.tmdbService.GetImageURL(searchResp.Results[i].PosterPath)
		}
		if searchResp.Results[i].BackdropPath != "" {
			searchResp.Results[i].BackdropPath = uc.tmdbService.GetImageURL(searchResp.Results[i].BackdropPath)
		}
	}

	return searchResp, nil
}
