package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/dto"
	"github.com/EduardoMG12/cine/api_v2/internal/infrastructure"
)

type RegisterUseCase struct {
	userRepo        domain.UserRepository
	sessionRepo     domain.SessionRepository
	passwordService *infrastructure.PasswordService
	jwtService      *infrastructure.JWTService
}

func NewRegisterUseCase(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	passwordService *infrastructure.PasswordService,
	jwtService *infrastructure.JWTService,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (uc *RegisterUseCase) Execute(input dto.RegisterRequestDTO) (*dto.AuthResponseDTO, error) {
	// Verificar se email já existe
	existingUser, _ := uc.userRepo.GetUserByEmail(strings.ToLower(input.Email))
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Verificar se username já existe
	existingUser, _ = uc.userRepo.GetUserByUsername(strings.ToLower(input.Username))
	if existingUser != nil {
		return nil, fmt.Errorf("username already taken")
	}

	// Hash da senha
	hashedPassword, err := uc.passwordService.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Criar usuário
	user := &domain.User{
		Username:      strings.ToLower(input.Username),
		Email:         strings.ToLower(input.Email),
		DisplayName:   input.DisplayName,
		PasswordHash:  hashedPassword,
		IsPrivate:     false,
		EmailVerified: false,
		Theme:         "light",
	}

	if err := uc.userRepo.CreateUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Gerar token JWT
	token, err := uc.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Criar sessão
	session := &domain.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := uc.sessionRepo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Retornar resposta
	return &dto.AuthResponseDTO{
		Token: token,
		User:  uc.userToDTO(user),
	}, nil
}

func (uc *RegisterUseCase) userToDTO(user *domain.User) dto.UserDTO {
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
