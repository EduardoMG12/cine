package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

type LoginUseCase struct {
	userRepo        domain.UserRepository
	sessionRepo     domain.SessionRepository
	passwordService *infrastructure.PasswordService
	jwtService      *infrastructure.JWTService
}

func NewLoginUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	passwordService *infrastructure.PasswordService,
	jwtService *infrastructure.JWTService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uc *LoginUseCase) Execute(input dto.LoginRequestDTO) (*dto.AuthResponseDTO, error) {
	user, err := uc.userRepo.GetUserByEmail(strings.ToLower(input.Email))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !uc.passwordService.ComparePassword(user.PasswordHash, input.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := uc.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := &domain.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := uc.sessionRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &dto.AuthResponseDTO{
		Token: token,
		User:  uc.userToDTO(user),
	}, nil
}

func (uc *LoginUseCase) userToDTO(user *domain.User) dto.UserDTO {
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
