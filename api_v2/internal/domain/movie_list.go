package domain

import (
	"time"
)

type MovieList struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	IsDefault bool      `json:"is_default" db:"is_default"` // for "Want to Watch", "Watched"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	User    *User             `json:"user,omitempty"`
	Entries []*MovieListEntry `json:"entries,omitempty"`
}

type MovieListEntry struct {
	ID          int       `json:"id" db:"id"`
	MovieListID int       `json:"movie_list_id" db:"movie_list_id"`
	MovieID     int       `json:"movie_id" db:"movie_id"`
	AddedAt     time.Time `json:"added_at" db:"added_at"`

	Movie *Movie `json:"movie,omitempty"`
}

type MovieListRepository interface {
	Create(list *MovieList) error
	GetByID(id int) (*MovieList, error)
	GetByUserID(userID int) ([]*MovieList, error)
	GetDefaultList(userID int, listType string) (*MovieList, error) // "want_to_watch", "watched"
	Update(list *MovieList) error
	Delete(id int) error

	AddMovieToList(listID, movieID int) error
	RemoveMovieFromList(listID, movieID int) error
	GetListEntries(listID int, limit, offset int) ([]*MovieListEntry, error)
	IsMovieInList(listID, movieID int) (bool, error)
}

type MovieListService interface {
	CreateList(userID int, name string) (*MovieList, error)
	GetUserLists(userID int) ([]*MovieList, error)
	GetList(listID int) (*MovieList, error)
	UpdateList(listID, userID int, name string) (*MovieList, error)
	DeleteList(listID, userID int) error

	AddToWantToWatch(userID, movieID int) error
	AddToWatched(userID, movieID int) error
	RemoveFromWantToWatch(userID, movieID int) error
	RemoveFromWatched(userID, movieID int) error
	MoveToWatched(userID, movieID int) error

	AddMovieToList(listID, userID, movieID int) error
	RemoveMovieFromList(listID, userID, movieID int) error
	GetListMovies(listID int, page int) ([]*MovieListEntry, error)
}
