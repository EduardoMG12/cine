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
	ID         string         `json:"id" db:"id"`
	UserID     string         `json:"user_id" db:"user_id"`
	Content    string         `json:"content" db:"content"`
	Visibility PostVisibility `json:"visibility" db:"visibility"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at" db:"updated_at"`

	User *User `json:"user,omitempty"`
}

type PostRepository interface {
	Create(post *Post) error
	GetByID(id string) (*Post, error)
	GetByUserID(userID string, limit, offset int) ([]*Post, error)
	GetPublicFeed(limit, offset int) ([]*Post, error)
	GetUserFeed(userID string, limit, offset int) ([]*Post, error) // includes friends/following posts
	Update(post *Post) error
	Delete(id string) error
}

type PostService interface {
	CreatePost(userID string, content string, visibility PostVisibility) (*Post, error)
	GetPost(postID string, requesterID string) (*Post, error) // respects visibility
	GetUserPosts(userID string, requesterID string, page int) ([]*Post, error)
	GetPublicFeed(page int) ([]*Post, error)
	GetUserFeed(userID string, page int) ([]*Post, error)
	UpdatePost(postID, userID string, content string, visibility PostVisibility) (*Post, error)
	DeletePost(postID, userID string) error
	ValidatePost(post *Post) error
}
