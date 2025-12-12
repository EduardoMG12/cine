package movie

import (
	"fmt"
	"log"

	"github.com/EduardoMG12/cine/api_v2/internal/data"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

type GetTrendingMoviesUseCase struct {
	movieRepo    domain.MovieRepository
	movieFetcher infrastructure.MovieFetcher
}

func NewGetTrendingMoviesUseCase(
	movieRepo domain.MovieRepository,
	movieFetcher infrastructure.MovieFetcher,
) *GetTrendingMoviesUseCase {
	return &GetTrendingMoviesUseCase{
		movieRepo:    movieRepo,
		movieFetcher: movieFetcher,
	}
}

func (uc *GetTrendingMoviesUseCase) Execute() ([]*dto.MovieDTO, error) {
	const minMoviesThreshold = 100
	const randomMoviesCount = 300

	// Check database first
	count, err := uc.movieRepo.CountMovies()
	if err != nil {
		log.Printf("❌ Error counting movies: %v", err)
		return nil, fmt.Errorf("failed to count movies: %w", err)
	}

	log.Printf("📊 Database has %d movies", count)

	// If we have more than 100 movies, return 300 random ones
	if count >= minMoviesThreshold {
		log.Printf("✅ Database has enough movies, returning %d random ones", randomMoviesCount)
		movies, err := uc.movieRepo.GetRandomMovies(randomMoviesCount)
		if err != nil {
			return nil, fmt.Errorf("failed to get random movies: %w", err)
		}
		return uc.moviesToDTOs(movies), nil
	}

	// Database has less than 100 movies - fetch ALL from seeds and save them
	log.Printf("⚠️  Database has only %d movies, fetching all from seeds...", count)

	saved := 0
	errors := 0

	// Iterate through ALL seeds
	for category, seeds := range data.MovieSeeds {
		log.Printf("🔍 Processing category: %s (%d movies)", category, len(seeds))

		for _, seed := range seeds {
			log.Printf("   Searching for: %s", seed.Title)

			movies, err := uc.movieFetcher.Search(seed.Title, 1)
			if err != nil {
				log.Printf("      ❌ Search failed: %v", err)
				errors++
				continue
			}

			if len(movies) > 0 {
				// Movie was found and automatically saved by the fetcher
				log.Printf("      ✅ Found and saved: %s", movies[0].Title)
				saved++
			} else {
				log.Printf("      ⚠️  No results")
			}
		}
	}

	log.Printf("🎬 Population complete: %d saved, %d errors", saved, errors)

	// Now return 100 random movies from the database
	movies, err := uc.movieRepo.GetRandomMovies(100)
	if err != nil {
		return nil, fmt.Errorf("failed to get random movies after population: %w", err)
	}

	return uc.moviesToDTOs(movies), nil
}

func (uc *GetTrendingMoviesUseCase) moviesToDTOs(movies []*domain.Movie) []*dto.MovieDTO {
	dtos := make([]*dto.MovieDTO, len(movies))
	for i, movie := range movies {
		dtos[i] = uc.movieToDTO(movie)
	}
	return dtos
}

func (uc *GetTrendingMoviesUseCase) movieToDTO(movie *domain.Movie) *dto.MovieDTO {
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
