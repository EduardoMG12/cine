package user

import (
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type UpdateUserUseCase struct {
	userRepo domain.UserRepository
}

func NewUpdateUserUseCase(userRepo domain.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UpdateUserUseCase) Execute(userID uuid.UUID, req *dto.UpdateUserRequest) (*dto.UserProfileResponse, error) {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if req.DisplayName != nil {
		user.DisplayName = *req.DisplayName
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if req.ProfilePictureURL != nil {
		user.ProfilePictureURL = req.ProfilePictureURL
	}

	if req.Theme != nil {
		user.Theme = *req.Theme
	}

	if req.IsPrivate != nil {
		user.IsPrivate = *req.IsPrivate
	}

	user.UpdatedAt = time.Now()

	if err := uc.userRepo.UpdateUser(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &dto.UserProfileResponse{
		ID:                user.ID,
		Username:          user.Username,
		Email:             user.Email,
		DisplayName:       user.DisplayName,
		Bio:               user.Bio,
		ProfilePictureURL: user.ProfilePictureURL,
		IsPrivate:         user.IsPrivate,
		EmailVerified:     user.EmailVerified,
		Theme:             user.Theme,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
	}, nil
}
