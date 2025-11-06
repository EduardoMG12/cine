package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileRequest struct {
	DisplayName       *string `json:"display_name" validate:"omitempty,min=2,max=100"`
	Bio               *string `json:"bio" validate:"omitempty,max=500"`
	ProfilePictureURL *string `json:"profile_picture_url" validate:"omitempty,url"`
}

type UpdateSettingsRequest struct {
	Theme     *string `json:"theme" validate:"omitempty,oneof=light dark"`
	IsPrivate *bool   `json:"is_private"`
}

type PublicProfile struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	DisplayName       string    `json:"display_name"`
	Bio               *string   `json:"bio,omitempty"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	IsPrivate         bool      `json:"is_private"`
	CreatedAt         time.Time `json:"created_at"`
	ReviewCount       int       `json:"review_count,omitempty"`
	MovieListCount    int       `json:"movie_list_count,omitempty"`
	FollowerCount     int       `json:"follower_count,omitempty"`
	FollowingCount    int       `json:"following_count,omitempty"`
}

type UserSearchResult struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	DisplayName       string    `json:"display_name"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	IsPrivate         bool      `json:"is_private"`
}

type UsernameAvailabilityResponse struct {
	Username  string `json:"username"`
	Available bool   `json:"available"`
}

type UserProfileResponse struct {
	ID                uuid.UUID `json:"id"`
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	DisplayName       string    `json:"display_name"`
	Bio               *string   `json:"bio,omitempty"`
	ProfilePictureURL *string   `json:"profile_picture_url,omitempty"`
	IsPrivate         bool      `json:"is_private"`
	EmailVerified     bool      `json:"email_verified"`
	Theme             string    `json:"theme"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
