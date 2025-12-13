package dto

// UpdateUserRequest represents request to update user profile
type UpdateUserRequest struct {
	DisplayName       *string `json:"display_name" validate:"omitempty,min=2,max=100"`
	Bio               *string `json:"bio" validate:"omitempty,max=500"`
	ProfilePictureURL *string `json:"profile_picture_url" validate:"omitempty,url"`
	Theme             *string `json:"theme" validate:"omitempty,oneof=light dark"`
	IsPrivate         *bool   `json:"is_private"`
}
