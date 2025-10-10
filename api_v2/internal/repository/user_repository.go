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
		INSERT INTO users (username, email, password_hash, display_name, bio, profile_picture_url, is_private, email_verified, theme, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Bio,
		user.ProfilePictureURL,
		user.IsPrivate,
		user.EmailVerified,
		user.Theme,
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
		SELECT id, username, email, display_name, bio, profile_picture_url, created_at, updated_at
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
		SELECT id, username, email, display_name, bio, profile_picture_url, created_at, updated_at
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
		SELECT id, username, email, display_name, bio, profile_picture_url, created_at, updated_at
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
		SET username = $2, email = $3, display_name = $4, bio = $5, profile_picture_url = $6, updated_at = $7
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
		user.ProfilePictureURL,
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

func (r *userRepository) UpdateSettings(userID int, settings map[string]interface{}) error {
	if len(settings) == 0 {
		return nil
	}

	query := "UPDATE users SET updated_at = NOW()"
	args := []interface{}{}
	argIndex := 1

	for key, value := range settings {
		switch key {
		case "theme":
			query += fmt.Sprintf(", theme = $%d", argIndex)
			args = append(args, value)
			argIndex++
		}
	}

	query += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, userID)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user settings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.invalidateUserCache(userID)
	return nil
}

func (r *userRepository) invalidateUserCache(id int) {
	cacheKey := fmt.Sprintf("user:%d", id)
	r.redis.Del(context.Background(), cacheKey)
}

// CreateEmailVerificationToken creates a new email verification token
func (r *userRepository) CreateEmailVerificationToken(token *domain.EmailVerificationToken) error {
	query := `
		INSERT INTO email_verification_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return r.db.QueryRow(
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(&token.ID)
}

// GetEmailVerificationToken retrieves an email verification token
func (r *userRepository) GetEmailVerificationToken(tokenStr string) (*domain.EmailVerificationToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM email_verification_tokens
		WHERE token = $1
	`

	var token domain.EmailVerificationToken
	err := r.db.QueryRow(query, tokenStr).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// DeleteEmailVerificationToken deletes an email verification token
func (r *userRepository) DeleteEmailVerificationToken(tokenStr string) error {
	query := `DELETE FROM email_verification_tokens WHERE token = $1`
	_, err := r.db.Exec(query, tokenStr)
	return err
}

// MarkEmailAsVerified marks a user's email as verified
func (r *userRepository) MarkEmailAsVerified(userID int) error {
	query := `UPDATE users SET email_verified = true WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	if err == nil {
		r.invalidateUserCache(userID)
	}
	return err
}

// CreatePasswordResetToken creates a new password reset token
func (r *userRepository) CreatePasswordResetToken(token *domain.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	return r.db.QueryRow(
		query,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	).Scan(&token.ID)
}

// GetPasswordResetToken retrieves a password reset token
func (r *userRepository) GetPasswordResetToken(tokenStr string) (*domain.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM password_reset_tokens
		WHERE token = $1
	`

	var token domain.PasswordResetToken
	err := r.db.QueryRow(query, tokenStr).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// DeletePasswordResetToken deletes a password reset token
func (r *userRepository) DeletePasswordResetToken(tokenStr string) error {
	query := `DELETE FROM password_reset_tokens WHERE token = $1`
	_, err := r.db.Exec(query, tokenStr)
	return err
}
