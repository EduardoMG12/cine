package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `db:"id" json:"id"`
	Username          string    `db:"username" json:"username"`
	Email             string    `db:"email" json:"email"`
	DisplayName       string    `db:"display_name" json:"display_name"`
	Bio               *string   `db:"bio" json:"bio,omitempty"`
	ProfilePictureURL *string   `db:"profile_picture_url" json:"profile_picture_url,omitempty"`
	PasswordHash      string    `db:"password_hash" json:"-"`
	IsPrivate         bool      `db:"is_private" json:"is_private"`
	EmailVerified     bool      `db:"email_verified" json:"email_verified"`
	Theme             string    `db:"theme" json:"theme"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type UserSession struct {
	ID        uuid.UUID `db:"id" json:"-"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UserAgent *string   `db:"user_agent" json:"user_agent,omitempty"`
	IPAddress *string   `db:"ip_address" json:"ip_address,omitempty"`
}

// UserRepository defines the repository interface for user operations
type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uuid.UUID) error
}

// SessionRepository defines the repository interface for session operations
type SessionRepository interface {
	CreateSession(session *UserSession) error
	GetSessionByToken(token string) (*UserSession, error)
	DeleteSession(token string) error
	DeleteUserSessions(userID uuid.UUID) error
}

// UserService defines the service interface for user operations
type UserService interface {
	GetProfile(userID uuid.UUID) (*User, error)
	UpdateProfile(userID uuid.UUID, displayName, bio, profilePictureURL *string) (*User, error)
	UpdateSettings(userID uuid.UUID, theme *string, isPrivate *bool) error
	CheckUsernameAvailability(username string, excludeUserID *uuid.UUID) (bool, error)
}
