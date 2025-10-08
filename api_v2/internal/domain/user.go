package domain

import (
	"time"
)

type User struct {
	ID          int       `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Email       string    `json:"email" db:"email"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Bio         *string   `json:"bio,omitempty" db:"bio"`
	AvatarURL   *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id int) error
}

type UserService interface {
	CreateUser(username, email, displayName string) (*User, error)
	GetUser(id int) (*User, error)
	UpdateUser(id int, updates map[string]interface{}) (*User, error)
	DeleteUser(id int) error
	ValidateUser(user *User) error
}
