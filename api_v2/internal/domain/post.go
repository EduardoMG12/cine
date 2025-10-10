package domain

import (
	"time"
)

type PostVisibility string

const (
	PostVisibilityPublic  PostVisibility = "public"
	PostVisibilityPrivate PostVisibility = "private"
	PostVisibilityFriends PostVisibility = "friends"
)

type Post struct {
	ID         int            `json:"id" db:"id"`
	UserID     int            `json:"user_id" db:"user_id"`
	Content    string         `json:"content" db:"content"`
	Visibility PostVisibility `json:"visibility" db:"visibility"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`

	// Populated by joins
	User *User `json:"user,omitempty"`
}

type PostRepository interface {
	Create(post *Post) error
	GetByID(id int) (*Post, error)
	GetByUserID(userID int, limit, offset int) ([]*Post, error)
	GetPublicFeed(limit, offset int) ([]*Post, error)
	GetUserFeed(userID int, limit, offset int) ([]*Post, error) // includes friends/following posts
	Update(post *Post) error
	Delete(id int) error
}

type PostService interface {
	CreatePost(userID int, content string, visibility PostVisibility) (*Post, error)
	GetPost(postID int, requesterID int) (*Post, error) // respects visibility
	GetUserPosts(userID int, requesterID int, page int) ([]*Post, error)
	GetPublicFeed(page int) ([]*Post, error)
	GetUserFeed(userID int, page int) ([]*Post, error)
	UpdatePost(postID, userID int, content string, visibility PostVisibility) (*Post, error)
	DeletePost(postID, userID int) error
	ValidatePost(post *Post) error
}
