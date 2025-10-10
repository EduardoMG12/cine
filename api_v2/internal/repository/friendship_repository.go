package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/jmoiron/sqlx"
)

type friendshipRepository struct {
	db *sqlx.DB
}

func NewFriendshipRepository(db *sqlx.DB) domain.FriendshipRepository {
	return &friendshipRepository{db: db}
}

func (r *friendshipRepository) Create(friendship *domain.Friendship) error {
	query := `
		INSERT INTO friendships (user_id_1, user_id_2, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)`

	friendship.CreatedAt = time.Now()
	friendship.UpdatedAt = time.Now()

	_, err := r.db.Exec(query,
		friendship.UserID1,
		friendship.UserID2,
		friendship.Status,
		friendship.CreatedAt,
		friendship.UpdatedAt,
	)

	return err
}

func (r *friendshipRepository) GetByUsers(userID1, userID2 string) (*domain.Friendship, error) {
	// Validate UUID formats
	if !utils.IsValidUUID(userID1) {
		return nil, fmt.Errorf("invalid user1 UUID format: %s", userID1)
	}
	if !utils.IsValidUUID(userID2) {
		return nil, fmt.Errorf("invalid user2 UUID format: %s", userID2)
	}

	query := `
		SELECT user_id_1, user_id_2, status, created_at, updated_at
		FROM friendships
		WHERE (user_id_1 = $1 AND user_id_2 = $2) 
		   OR (user_id_1 = $2 AND user_id_2 = $1)`

	var friendship domain.Friendship
	err := r.db.Get(&friendship, query, userID1, userID2)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &friendship, nil
}

func (r *friendshipRepository) GetUserFriends(userID string) ([]*domain.Friendship, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return nil, fmt.Errorf("invalid UUID format: %s", userID)
	}
	query := `
		SELECT 
			f.user_id_1, f.user_id_2, f.status, f.created_at, f.updated_at,
			u1.id, u1.username, u1.display_name, u1.profile_picture_url, u1.is_private,
			u2.id, u2.username, u2.display_name, u2.profile_picture_url, u2.is_private
		FROM friendships f
		JOIN users u1 ON f.user_id_1 = u1.id
		JOIN users u2 ON f.user_id_2 = u2.id
		WHERE (f.user_id_1 = $1 OR f.user_id_2 = $1) 
		AND f.status = 'accepted'
		ORDER BY f.updated_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*domain.Friendship

	for rows.Next() {
		var friendship domain.Friendship
		var user1, user2 domain.User

		err := rows.Scan(
			&friendship.UserID1, &friendship.UserID2, &friendship.Status,
			&friendship.CreatedAt, &friendship.UpdatedAt,
			&user1.ID, &user1.Username, &user1.DisplayName,
			&user1.ProfilePictureURL, &user1.IsPrivate,
			&user2.ID, &user2.Username, &user2.DisplayName,
			&user2.ProfilePictureURL, &user2.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		friendship.User1 = &user1
		friendship.User2 = &user2
		friendships = append(friendships, &friendship)
	}

	return friendships, nil
}

func (r *friendshipRepository) GetFriendRequests(userID string) ([]*domain.Friendship, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return nil, fmt.Errorf("invalid UUID format: %s", userID)
	}
	query := `
		SELECT 
			f.user_id_1, f.user_id_2, f.status, f.created_at, f.updated_at,
			u1.id, u1.username, u1.display_name, u1.profile_picture_url, u1.is_private
		FROM friendships f
		JOIN users u1 ON f.user_id_1 = u1.id
		WHERE f.user_id_2 = $1 
		AND f.status = 'pending'
		ORDER BY f.created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*domain.Friendship

	for rows.Next() {
		var friendship domain.Friendship
		var user1 domain.User

		err := rows.Scan(
			&friendship.UserID1, &friendship.UserID2, &friendship.Status,
			&friendship.CreatedAt, &friendship.UpdatedAt,
			&user1.ID, &user1.Username, &user1.DisplayName,
			&user1.ProfilePictureURL, &user1.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		friendship.User1 = &user1
		friendships = append(friendships, &friendship)
	}

	return friendships, nil
}

func (r *friendshipRepository) UpdateStatus(userID1, userID2 string, status domain.FriendshipStatus) error {
	// Validate UUID formats
	if !utils.IsValidUUID(userID1) {
		return fmt.Errorf("invalid user1 UUID format: %s", userID1)
	}
	if !utils.IsValidUUID(userID2) {
		return fmt.Errorf("invalid user2 UUID format: %s", userID2)
	}
	query := `
		UPDATE friendships 
		SET status = $1, updated_at = $2
		WHERE (user_id_1 = $3 AND user_id_2 = $4) 
		   OR (user_id_1 = $4 AND user_id_2 = $3)`

	result, err := r.db.Exec(query, status, time.Now(), userID1, userID2)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no friendship found between users %d and %d", userID1, userID2)
	}

	return nil
}

func (r *friendshipRepository) Delete(userID1, userID2 string) error {
	// Validate UUID formats
	if !utils.IsValidUUID(userID1) {
		return fmt.Errorf("invalid user1 UUID format: %s", userID1)
	}
	if !utils.IsValidUUID(userID2) {
		return fmt.Errorf("invalid user2 UUID format: %s", userID2)
	}

	query := `
		DELETE FROM friendships
		WHERE (user_id_1 = $1 AND user_id_2 = $2) 
		   OR (user_id_1 = $2 AND user_id_2 = $1)`

	result, err := r.db.Exec(query, userID1, userID2)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no friendship found between users %s and %s", userID1, userID2)
	}

	return nil
}
