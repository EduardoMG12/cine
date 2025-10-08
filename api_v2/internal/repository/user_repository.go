package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type userRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewUserRepository(db *sqlx.DB, redis *redis.Client) domain.UserRepository {
	return &userRepository{
		db:    db,
		redis: redis,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, display_name, bio, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.invalidateUserCache(user.ID)

	return nil
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	cached, err := r.redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var user domain.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}

	query := `
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err = r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if userData, err := json.Marshal(user); err == nil {
		r.redis.Set(context.Background(), cacheKey, userData, time.Hour)
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*domain.User, error) {
	query := `
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	var user domain.User
	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, display_name = $4, bio = $5, avatar_url = $6, updated_at = $7
		WHERE id = $1
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		user.ID,
		user.Username,
		user.Email,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.invalidateUserCache(user.ID)

	return nil
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.invalidateUserCache(id)

	return nil
}

func (r *userRepository) invalidateUserCache(id int) {
	cacheKey := fmt.Sprintf("user:%d", id)
	r.redis.Del(context.Background(), cacheKey)
}
