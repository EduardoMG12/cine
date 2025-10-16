package domain

import (
	"time"
)

type MovieList struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	IsDefault bool      `json:"is_default" db:"is_default"` // for "Want to Watch", "Watched"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	User    *User             `json:"user,omitempty"`
	Entries []*MovieListEntry `json:"entries,omitempty"`
}

type MovieListEntry struct {
	ID          string    `json:"id" db:"id"`
	MovieListID string    `json:"movie_list_id" db:"movie_list_id"`
	MovieID     string    `json:"movie_id" db:"movie_id"`
	AddedAt     time.Time `json:"added_at" db:"added_at"`

	Movie *Movie `json:"movie,omitempty"`
}

type MovieListRepository interface {
	Create(list *MovieList) error
	GetByID(id string) (*MovieList, error)
	GetByUserID(userID string) ([]*MovieList, error)
	GetDefaultList(userID string, listType string) (*MovieList, error) // "want_to_watch", "watched"
	Update(list *MovieList) error
	Delete(id string) error

	AddMovieToList(listID, movieID string) error
	RemoveMovieFromList(listID, movieID string) error
	GetListEntries(listID string, limit, offset int) ([]*MovieListEntry, error)
	IsMovieInList(listID, movieID string) (bool, error)
}

type MovieListService interface {
	CreateList(userID string, name string) (*MovieList, error)
	GetUserLists(userID string) ([]*MovieList, error)
	GetList(listID string) (*MovieList, error)
	UpdateList(listID, userID string, name string) (*MovieList, error)
	DeleteList(listID, userID string) error

	AddToWantToWatch(userID, movieID string) error
	AddToWatched(userID, movieID string) error
	RemoveFromWantToWatch(userID, movieID string) error
	RemoveFromWatched(userID, movieID string) error
	MoveToWatched(userID, movieID string) error

	AddMovieToList(listID, userID, movieID string) error
	RemoveMovieFromList(listID, userID, movieID string) error
	GetListMovies(listID string, page int) ([]*MovieListEntry, error)
}
