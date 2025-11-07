package movie

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

const (
	RedisCacheDuration = 30 * 24 * time.Hour // 30 days
	RedisMoviePrefix   = "movie:"
)

type GetMovieByIDUseCase struct {
	movieRepo    domain.MovieRepository
	tmdbService  *infrastructure.TMDbService
	redisService *infrastructure.RedisService
}

func NewGetMovieByIDUseCase(
	movieRepo domain.MovieRepository,
	tmdbService *infrastructure.TMDbService,
	redisService *infrastructure.RedisService,
) *GetMovieByIDUseCase {
	return &GetMovieByIDUseCase{
		movieRepo:    movieRepo,
		tmdbService:  tmdbService,
		redisService: redisService,
	}
}

func (uc *GetMovieByIDUseCase) Execute(tmdbID string) (*dto.MovieDTO, error) {
	ctx := context.Background()
	redisKey := RedisMoviePrefix + tmdbID

	// 1. Tentar buscar no Redis
	var cachedMovie dto.MovieDTO
	if err := uc.redisService.Get(ctx, redisKey, &cachedMovie); err == nil {
		return &cachedMovie, nil
	}

	// 2. Tentar buscar no PostgreSQL
	dbMovie, err := uc.movieRepo.GetMovieByExternalID(tmdbID)
	if err == nil && dbMovie.CacheExpiresAt.After(time.Now()) {
		movieDTO := uc.movieToDTO(dbMovie)

		// Salvar no Redis para próximas consultas
		_ = uc.redisService.Set(ctx, redisKey, movieDTO, RedisCacheDuration)

		return movieDTO, nil
	}

	// 3. Buscar no TMDb
	tmdbMovie, err := uc.tmdbService.GetMovie(tmdbID)
	if err != nil {
		return nil, fmt.Errorf("movie not found: %w", err)
	}

	// 4. Converter e salvar no PostgreSQL
	movie := uc.tmdbMovieToMovie(tmdbMovie)
	movie.CacheExpiresAt = time.Now().Add(RedisCacheDuration)

	// Se o filme já existe no DB (cache expirado), atualizar
	if dbMovie != nil {
		movie.ID = dbMovie.ID
		if err := uc.movieRepo.UpdateMovie(movie); err != nil {
			return nil, fmt.Errorf("failed to update movie: %w", err)
		}
	} else {
		if err := uc.movieRepo.CreateMovie(movie); err != nil {
			return nil, fmt.Errorf("failed to create movie: %w", err)
		}
	}

	movieDTO := uc.movieToDTO(movie)

	// 5. Salvar no Redis
	_ = uc.redisService.Set(ctx, redisKey, movieDTO, RedisCacheDuration)

	return movieDTO, nil
}

func (uc *GetMovieByIDUseCase) tmdbMovieToMovie(tmdbMovie *dto.TMDbMovieResponse) *domain.Movie {
	movie := &domain.Movie{
		ExternalAPIID: strconv.Itoa(tmdbMovie.ID),
		Title:         tmdbMovie.Title,
		Adult:         tmdbMovie.Adult,
	}

	if tmdbMovie.Overview != "" {
		movie.Overview = &tmdbMovie.Overview
	}

	if tmdbMovie.ReleaseDate != "" {
		if releaseDate, err := time.Parse("2006-01-02", tmdbMovie.ReleaseDate); err == nil {
			movie.ReleaseDate = &releaseDate
		}
	}

	if tmdbMovie.PosterPath != "" {
		posterURL := uc.tmdbService.GetImageURL(tmdbMovie.PosterPath)
		movie.PosterURL = &posterURL
	}

	if tmdbMovie.BackdropPath != "" {
		backdropURL := uc.tmdbService.GetImageURL(tmdbMovie.BackdropPath)
		movie.BackdropURL = &backdropURL
	}

	if len(tmdbMovie.Genres) > 0 {
		genres := make([]string, len(tmdbMovie.Genres))
		for i, genre := range tmdbMovie.Genres {
			genres[i] = genre.Name
		}
		movie.Genres = genres
	}

	if tmdbMovie.Runtime > 0 {
		movie.Runtime = &tmdbMovie.Runtime
	}

	if tmdbMovie.VoteAverage > 0 {
		movie.VoteAverage = &tmdbMovie.VoteAverage
	}

	if tmdbMovie.VoteCount > 0 {
		movie.VoteCount = &tmdbMovie.VoteCount
	}

	return movie
}

func (uc *GetMovieByIDUseCase) movieToDTO(movie *domain.Movie) *dto.MovieDTO {
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
