package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/jmoiron/sqlx"
)

type movieListRepository struct {
	db *sqlx.DB
}

func NewMovieListRepository(db *sqlx.DB) domain.MovieListRepository {
	return &movieListRepository{
		db: db,
	}
}

func (r *movieListRepository) Create(list *domain.MovieList) error {
	query := `
		INSERT INTO movie_lists (user_id, name, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	now := time.Now()
	list.CreatedAt = now
	list.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		list.UserID,
		list.Name,
		list.IsDefault,
		list.CreatedAt,
		list.UpdatedAt,
	).Scan(&list.ID)

	if err != nil {
		return fmt.Errorf("failed to create movie list: %w", err)
	}

	return nil
}

func (r *movieListRepository) GetByID(id int) (*domain.MovieList, error) {
	query := `
		SELECT ml.id, ml.user_id, ml.name, ml.is_default, ml.created_at, ml.updated_at,
		       u.id, u.username, u.email, u.created_at
		FROM movie_lists ml
		LEFT JOIN users u ON ml.user_id = u.id
		WHERE ml.id = $1`

	var list domain.MovieList
	var user domain.User

	err := r.db.QueryRow(query, id).Scan(
		&list.ID, &list.UserID, &list.Name, &list.IsDefault, &list.CreatedAt, &list.UpdatedAt,
		&user.ID, &user.Username, &user.Email, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie list not found")
		}
		return nil, fmt.Errorf("failed to get movie list: %w", err)
	}

	list.User = &user
	return &list, nil
}

func (r *movieListRepository) GetByUserID(userID int) ([]*domain.MovieList, error) {
	query := `
		SELECT id, user_id, name, is_default, created_at, updated_at
		FROM movie_lists
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user movie lists: %w", err)
	}
	defer rows.Close()

	var lists []*domain.MovieList
	for rows.Next() {
		var list domain.MovieList
		err := rows.Scan(
			&list.ID,
			&list.UserID,
			&list.Name,
			&list.IsDefault,
			&list.CreatedAt,
			&list.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie list: %w", err)
		}
		lists = append(lists, &list)
	}

	return lists, nil
}

func (r *movieListRepository) GetDefaultList(userID int, listType string) (*domain.MovieList, error) {
	var listName string
	switch listType {
	case "want_to_watch":
		listName = "Want to Watch"
	case "watched":
		listName = "Watched"
	default:
		return nil, fmt.Errorf("invalid list type: %s", listType)
	}

	query := `
		SELECT id, user_id, name, is_default, created_at, updated_at
		FROM movie_lists
		WHERE user_id = $1 AND is_default = true AND name = $2`

	var list domain.MovieList
	err := r.db.QueryRow(query, userID, listName).Scan(
		&list.ID,
		&list.UserID,
		&list.Name,
		&list.IsDefault,
		&list.CreatedAt,
		&list.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("default list not found")
		}
		return nil, fmt.Errorf("failed to get default list: %w", err)
	}

	return &list, nil
}

func (r *movieListRepository) Update(list *domain.MovieList) error {
	query := `
		UPDATE movie_lists
		SET name = $1, updated_at = $2
		WHERE id = $3`

	list.UpdatedAt = time.Now()

	_, err := r.db.Exec(query, list.Name, list.UpdatedAt, list.ID)
	if err != nil {
		return fmt.Errorf("failed to update movie list: %w", err)
	}

	return nil
}

func (r *movieListRepository) Delete(id int) error {
	// First delete all entries in the list
	_, err := r.db.Exec("DELETE FROM movie_list_entries WHERE movie_list_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete movie list entries: %w", err)
	}

	// Then delete the list itself
	_, err = r.db.Exec("DELETE FROM movie_lists WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete movie list: %w", err)
	}

	return nil
}

func (r *movieListRepository) AddMovieToList(listID, movieID int) error {
	query := `
		INSERT INTO movie_list_entries (movie_list_id, movie_id, added_at)
		VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, listID, movieID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add movie to list: %w", err)
	}

	return nil
}

func (r *movieListRepository) RemoveMovieFromList(listID, movieID int) error {
	query := `DELETE FROM movie_list_entries WHERE movie_list_id = $1 AND movie_id = $2`

	result, err := r.db.Exec(query, listID, movieID)
	if err != nil {
		return fmt.Errorf("failed to remove movie from list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("movie not found in list")
	}

	return nil
}

func (r *movieListRepository) GetListEntries(listID int, limit, offset int) ([]*domain.MovieListEntry, error) {
	query := `
		SELECT mle.id, mle.movie_list_id, mle.movie_id, mle.added_at,
		       m.id, m.external_api_id, m.title, m.overview, m.poster_url, m.backdrop_url,
		       m.release_date, m.genres, m.vote_average, m.runtime, m.vote_count, m.adult,
		       m.created_at, m.updated_at
		FROM movie_list_entries mle
		LEFT JOIN movies m ON mle.movie_id = m.id
		WHERE mle.movie_list_id = $1
		ORDER BY mle.added_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, listID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get list entries: %w", err)
	}
	defer rows.Close()

	var entries []*domain.MovieListEntry
	for rows.Next() {
		var entry domain.MovieListEntry
		var movie domain.Movie

		err := rows.Scan(
			&entry.ID, &entry.MovieListID, &entry.MovieID, &entry.AddedAt,
			&movie.ID, &movie.ExternalAPIID, &movie.Title, &movie.Overview, &movie.PosterURL, &movie.BackdropURL,
			&movie.ReleaseDate, &movie.Genres, &movie.VoteAverage, &movie.Runtime, &movie.VoteCount, &movie.Adult,
			&movie.CreatedAt, &movie.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan list entry: %w", err)
		}

		entry.Movie = &movie
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (r *movieListRepository) IsMovieInList(listID, movieID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM movie_list_entries WHERE movie_list_id = $1 AND movie_id = $2)`

	var exists bool
	err := r.db.QueryRow(query, listID, movieID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if movie exists in list: %w", err)
	}

	return exists, nil
}
