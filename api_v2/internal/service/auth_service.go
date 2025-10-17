package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
)

type authService struct {
	userRepo        domain.UserRepository
	sessionRepo     domain.SessionRepository
	passwordService domain.PasswordService
	jwtService      domain.JWTService
}

func NewAuthService(
	userRepo domain.UserRepository,
	sessionRepo domain.SessionRepository,
	passwordService domain.PasswordService,
	jwtService domain.JWTService,
) domain.AuthService {
	return &authService{
		userRepo:        userRepo,
		sessionRepo:     sessionRepo,
		passwordService: passwordService,
		jwtService:      jwtService,
	}
}

func (s *authService) Register(req domain.RegisterRequest) (*domain.AuthResponse, error) {

	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	existingUser, _ = s.userRepo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already taken")
	}

	hashedPassword, err := s.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user := &domain.User{
		Username:      strings.ToLower(req.Username),
		Email:         strings.ToLower(req.Email),
		DisplayName:   req.DisplayName,
		PasswordHash:  hashedPassword,
		IsPrivate:     false,
		EmailVerified: false,
		Theme:         "light",
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := &domain.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		UserAgent: nil,
		IPAddress: nil,
	}

	err = s.sessionRepo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) Login(req domain.LoginRequest) (*domain.AuthResponse, error) {

	user, err := s.userRepo.GetUserByEmail(strings.ToLower(req.Email))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !s.passwordService.ComparePassword(user.PasswordHash, req.Password) {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := &domain.UserSession{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.sessionRepo.CreateSession(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *authService) ValidateToken(token string) (*domain.User, error) {

	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	_, err = s.sessionRepo.GetSessionByToken(token)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	user, err := s.userRepo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (s *authService) Logout(token string) error {
	return s.sessionRepo.DeleteSession(token)
}

func (s *authService) LogoutAll(userID uuid.UUID) error {
	return s.sessionRepo.DeleteUserSessions(userID)
}
