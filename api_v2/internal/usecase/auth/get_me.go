package auth

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/google/uuid"
)

type GetMeUseCase struct {
	userRepo domain.UserRepository
}

func NewGetMeUseCase(userRepo domain.UserRepository) *GetMeUseCase {
	return &GetMeUseCase{
		userRepo: userRepo,
	}
}

func (uc *GetMeUseCase) Execute(userID uuid.UUID) (*dto.UserDTO, error) {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	userDTO := uc.userToDTO(user)
	return &userDTO, nil
}

func (uc *GetMeUseCase) userToDTO(user *domain.User) dto.UserDTO {
	return dto.UserDTO{
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
	}
}
