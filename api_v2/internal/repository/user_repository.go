package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (
			id, username, email, display_name, bio, profile_picture_url, 
			password_hash, is_private, email_verified, theme, created_at, updated_at
		) VALUES (
			:id, :username, :email, :display_name, :bio, :profile_picture_url,
			:password_hash, :is_private, :email_verified, :theme, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, display_name, bio, profile_picture_url, 
			   password_hash, is_private, email_verified, theme, created_at, updated_at
		FROM users 
		WHERE id = $1
	`

	err := r.db.Get(&user, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, display_name, bio, profile_picture_url, 
			   password_hash, is_private, email_verified, theme, created_at, updated_at
		FROM users 
		WHERE email = $1
	`

	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, email, display_name, bio, profile_picture_url, 
			   password_hash, is_private, email_verified, theme, created_at, updated_at
		FROM users 
		WHERE username = $1
	`

	err := r.db.Get(&user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *userRepository) UpdateUser(user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users SET 
			username = :username,
			email = :email,
			display_name = :display_name,
			bio = :bio,
			profile_picture_url = :profile_picture_url,
			is_private = :is_private,
			email_verified = :email_verified,
			theme = :theme,
			updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExec(query, user)
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

	return nil
}

func (r *userRepository) DeleteUser(id uuid.UUID) error {
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

	return nil
}
