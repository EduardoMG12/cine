package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/jmoiron/sqlx"
)

type followRepository struct {
	db *sqlx.DB
}

func NewFollowRepository(db *sqlx.DB) domain.FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(follow *domain.Follow) error {
	query := `
		INSERT INTO follows (follower_id, following_id, created_at)
		VALUES ($1, $2, $3)`

	follow.CreatedAt = time.Now()

	_, err := r.db.Exec(query,
		follow.FollowerID,
		follow.FollowingID,
		follow.CreatedAt,
	)

	return err
}

func (r *followRepository) GetByUsers(followerID, followingID string) (*domain.Follow, error) {
	// Validate UUID formats
	if !utils.IsValidUUID(followerID) {
		return nil, fmt.Errorf("invalid follower UUID format: %s", followerID)
	}
	if !utils.IsValidUUID(followingID) {
		return nil, fmt.Errorf("invalid following UUID format: %s", followingID)
	}

	query := `
		SELECT follower_id, following_id, created_at
		FROM follows
		WHERE follower_id = $1 AND following_id = $2`

	var follow domain.Follow
	err := r.db.Get(&follow, query, followerID, followingID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &follow, nil
}

func (r *followRepository) GetFollowers(userID string, limit, offset int) ([]*domain.Follow, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return nil, fmt.Errorf("invalid UUID format: %s", userID)
	}
	query := `
		SELECT 
			f.follower_id, f.following_id, f.created_at,
			u.id, u.username, u.display_name, u.profile_picture_url, u.is_private
		FROM follows f
		JOIN users u ON f.follower_id = u.id
		WHERE f.following_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*domain.Follow

	for rows.Next() {
		var follow domain.Follow
		var follower domain.User

		err := rows.Scan(
			&follow.FollowerID, &follow.FollowingID, &follow.CreatedAt,
			&follower.ID, &follower.Username, &follower.DisplayName,
			&follower.ProfilePictureURL, &follower.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		follow.Follower = &follower
		follows = append(follows, &follow)
	}

	return follows, nil
}

func (r *followRepository) GetFollowing(userID string, limit, offset int) ([]*domain.Follow, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return nil, fmt.Errorf("invalid UUID format: %s", userID)
	}
	query := `
		SELECT 
			f.follower_id, f.following_id, f.created_at,
			u.id, u.username, u.display_name, u.profile_picture_url, u.is_private
		FROM follows f
		JOIN users u ON f.following_id = u.id
		WHERE f.follower_id = $1
		ORDER BY f.created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var follows []*domain.Follow

	for rows.Next() {
		var follow domain.Follow
		var following domain.User

		err := rows.Scan(
			&follow.FollowerID, &follow.FollowingID, &follow.CreatedAt,
			&following.ID, &following.Username, &following.DisplayName,
			&following.ProfilePictureURL, &following.IsPrivate,
		)
		if err != nil {
			return nil, err
		}

		follow.Following = &following
		follows = append(follows, &follow)
	}

	return follows, nil
}

func (r *followRepository) Delete(followerID, followingID string) error {
	// Validate UUID formats
	if !utils.IsValidUUID(followerID) {
		return fmt.Errorf("invalid follower UUID format: %s", followerID)
	}
	if !utils.IsValidUUID(followingID) {
		return fmt.Errorf("invalid following UUID format: %s", followingID)
	}

	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND following_id = $2`

	result, err := r.db.Exec(query, followerID, followingID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no follow relationship found from user %s to user %s", followerID, followingID)
	}

	return nil
}

func (r *followRepository) GetFollowersCount(userID string) (int, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return 0, fmt.Errorf("invalid UUID format: %s", userID)
	}
	query := `SELECT COUNT(*) FROM follows WHERE following_id = $1`

	var count int
	err := r.db.Get(&count, query, userID)
	return count, err
}

func (r *followRepository) GetFollowingCount(userID string) (int, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return 0, fmt.Errorf("invalid UUID format: %s", userID)
	}
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`

	var count int
	err := r.db.Get(&count, query, userID)
	return count, err
}
