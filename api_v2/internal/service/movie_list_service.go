package service

import (
	"errors"
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type movieListService struct {
	movieListRepo domain.MovieListRepository
	movieRepo     domain.MovieRepository
}

func NewMovieListService(movieListRepo domain.MovieListRepository, movieRepo domain.MovieRepository) domain.MovieListService {
	return &movieListService{
		movieListRepo: movieListRepo,
		movieRepo:     movieRepo,
	}
}

func (s *movieListService) CreateList(userID int, name string) (*domain.MovieList, error) {
	if name == "" {
		return nil, errors.New("list name cannot be empty")
	}

	list := &domain.MovieList{
		UserID:    userID,
		Name:      name,
		IsDefault: false,
	}

	if err := s.movieListRepo.Create(list); err != nil {
		return nil, fmt.Errorf("failed to create movie list: %w", err)
	}

	return list, nil
}

func (s *movieListService) GetUserLists(userID int) ([]*domain.MovieList, error) {
	lists, err := s.movieListRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user lists: %w", err)
	}

	return lists, nil
}

func (s *movieListService) GetList(listID int) (*domain.MovieList, error) {
	list, err := s.movieListRepo.GetByID(listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie list: %w", err)
	}

	return list, nil
}

func (s *movieListService) UpdateList(listID, userID int, name string) (*domain.MovieList, error) {
	if name == "" {
		return nil, errors.New("list name cannot be empty")
	}

	list, err := s.movieListRepo.GetByID(listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie list: %w", err)
	}

	// Check ownership
	if list.UserID != userID {
		return nil, errors.New("unauthorized: user does not own this list")
	}

	// Don't allow updating default lists
	if list.IsDefault {
		return nil, errors.New("cannot update default lists")
	}

	list.Name = name
	if err := s.movieListRepo.Update(list); err != nil {
		return nil, fmt.Errorf("failed to update movie list: %w", err)
	}

	return list, nil
}

func (s *movieListService) DeleteList(listID, userID int) error {
	list, err := s.movieListRepo.GetByID(listID)
	if err != nil {
		return fmt.Errorf("failed to get movie list: %w", err)
	}

	// Check ownership
	if list.UserID != userID {
		return errors.New("unauthorized: user does not own this list")
	}

	// Don't allow deleting default lists
	if list.IsDefault {
		return errors.New("cannot delete default lists")
	}

	if err := s.movieListRepo.Delete(listID); err != nil {
		return fmt.Errorf("failed to delete movie list: %w", err)
	}

	return nil
}

func (s *movieListService) AddToWantToWatch(userID, movieID int) error {
	// Ensure movie exists
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return fmt.Errorf("movie not found: %w", err)
	}

	// Get or create default want-to-watch list
	list, err := s.movieListRepo.GetDefaultList(userID, "want_to_watch")
	if err != nil {
		// Create default list if it doesn't exist
		list = &domain.MovieList{
			UserID:    userID,
			Name:      "Want to Watch",
			IsDefault: true,
		}
		if err := s.movieListRepo.Create(list); err != nil {
			return fmt.Errorf("failed to create want-to-watch list: %w", err)
		}
	}

	// Check if movie is already in list
	exists, err := s.movieListRepo.IsMovieInList(list.ID, movieID)
	if err != nil {
		return fmt.Errorf("failed to check if movie exists in list: %w", err)
	}
	if exists {
		return errors.New("movie already in want-to-watch list")
	}

	if err := s.movieListRepo.AddMovieToList(list.ID, movieID); err != nil {
		return fmt.Errorf("failed to add movie to want-to-watch list: %w", err)
	}

	return nil
}

func (s *movieListService) AddToWatched(userID, movieID int) error {
	// Ensure movie exists
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return fmt.Errorf("movie not found: %w", err)
	}

	// Get or create default watched list
	list, err := s.movieListRepo.GetDefaultList(userID, "watched")
	if err != nil {
		// Create default list if it doesn't exist
		list = &domain.MovieList{
			UserID:    userID,
			Name:      "Watched",
			IsDefault: true,
		}
		if err := s.movieListRepo.Create(list); err != nil {
			return fmt.Errorf("failed to create watched list: %w", err)
		}
	}

	// Check if movie is already in list
	exists, err := s.movieListRepo.IsMovieInList(list.ID, movieID)
	if err != nil {
		return fmt.Errorf("failed to check if movie exists in list: %w", err)
	}
	if exists {
		return errors.New("movie already in watched list")
	}

	if err := s.movieListRepo.AddMovieToList(list.ID, movieID); err != nil {
		return fmt.Errorf("failed to add movie to watched list: %w", err)
	}

	return nil
}

func (s *movieListService) RemoveFromWantToWatch(userID, movieID int) error {
	list, err := s.movieListRepo.GetDefaultList(userID, "want_to_watch")
	if err != nil {
		return fmt.Errorf("want-to-watch list not found: %w", err)
	}

	if err := s.movieListRepo.RemoveMovieFromList(list.ID, movieID); err != nil {
		return fmt.Errorf("failed to remove movie from want-to-watch list: %w", err)
	}

	return nil
}

func (s *movieListService) RemoveFromWatched(userID, movieID int) error {
	list, err := s.movieListRepo.GetDefaultList(userID, "watched")
	if err != nil {
		return fmt.Errorf("watched list not found: %w", err)
	}

	if err := s.movieListRepo.RemoveMovieFromList(list.ID, movieID); err != nil {
		return fmt.Errorf("failed to remove movie from watched list: %w", err)
	}

	return nil
}

func (s *movieListService) MoveToWatched(userID, movieID int) error {
	// Remove from want-to-watch list (ignore error if not in list)
	wantToWatchList, err := s.movieListRepo.GetDefaultList(userID, "want_to_watch")
	if err == nil {
		s.movieListRepo.RemoveMovieFromList(wantToWatchList.ID, movieID)
	}

	// Add to watched list
	return s.AddToWatched(userID, movieID)
}

func (s *movieListService) AddMovieToList(listID, userID, movieID int) error {
	// Check if list exists and user owns it
	list, err := s.movieListRepo.GetByID(listID)
	if err != nil {
		return fmt.Errorf("movie list not found: %w", err)
	}

	if list.UserID != userID {
		return errors.New("unauthorized: user does not own this list")
	}

	// Ensure movie exists
	if _, err := s.movieRepo.GetByID(movieID); err != nil {
		return fmt.Errorf("movie not found: %w", err)
	}

	// Check if movie is already in list
	exists, err := s.movieListRepo.IsMovieInList(listID, movieID)
	if err != nil {
		return fmt.Errorf("failed to check if movie exists in list: %w", err)
	}
	if exists {
		return errors.New("movie already in list")
	}

	if err := s.movieListRepo.AddMovieToList(listID, movieID); err != nil {
		return fmt.Errorf("failed to add movie to list: %w", err)
	}

	return nil
}

func (s *movieListService) RemoveMovieFromList(listID, userID, movieID int) error {
	// Check if list exists and user owns it
	list, err := s.movieListRepo.GetByID(listID)
	if err != nil {
		return fmt.Errorf("movie list not found: %w", err)
	}

	if list.UserID != userID {
		return errors.New("unauthorized: user does not own this list")
	}

	if err := s.movieListRepo.RemoveMovieFromList(listID, movieID); err != nil {
		return fmt.Errorf("failed to remove movie from list: %w", err)
	}

	return nil
}

func (s *movieListService) GetListMovies(listID int, page int) ([]*domain.MovieListEntry, error) {
	if page < 1 {
		page = 1
	}

	limit := 20
	offset := (page - 1) * limit

	entries, err := s.movieListRepo.GetListEntries(listID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get list movies: %w", err)
	}

	return entries, nil
}
