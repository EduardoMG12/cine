package domain

import (
	"time"
)

type User struct {
	ID                string    `json:"id" db:"id"`
	Username          string    `json:"username" db:"username"`
	Email             string    `json:"email" db:"email"`
	PasswordHash      string    `json:"-" db:"password_hash"`
	DisplayName       string    `json:"display_name" db:"display_name"`
	Bio               *string   `json:"bio,omitempty" db:"bio"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty" db:"profile_picture_url"`
	IsPrivate         bool      `json:"is_private" db:"is_private"`
	EmailVerified     bool      `json:"email_verified" db:"email_verified"`
	Theme             string    `json:"theme" db:"theme"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// EmailVerificationToken represents a token for email verification
type EmailVerificationToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// PasswordResetToken represents a token for password reset
type PasswordResetToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id string) error
	UpdateSettings(userID string, settings map[string]interface{}) error

	// Email verification
	CreateEmailVerificationToken(token *EmailVerificationToken) error
	GetEmailVerificationToken(token string) (*EmailVerificationToken, error)
	DeleteEmailVerificationToken(token string) error
	MarkEmailAsVerified(userID string) error

	// Password reset
	CreatePasswordResetToken(token *PasswordResetToken) error
	GetPasswordResetToken(token string) (*PasswordResetToken, error)
	DeletePasswordResetToken(token string) error
}

type UserService interface {
	Register(username, email, password, displayName string) (*User, error)
	Login(email, password string) (*User, string, error) // returns user, JWT token, error
	ConfirmEmail(token string) error
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserProfile(id string) (*User, error) // respects privacy settings
	UpdateProfile(userID string, updates map[string]interface{}) (*User, error)
	UpdateSettings(userID string, settings map[string]interface{}) error
	DeleteUser(id string) error
	ValidateUser(user *User) error
	CreateUser(user *User) error // for creating user entities directly
}
