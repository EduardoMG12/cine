package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type userSessionService struct {
	userSessionRepo domain.UserSessionRepository
	sessionDuration time.Duration
}

func NewUserSessionService(userSessionRepo domain.UserSessionRepository, sessionDuration time.Duration) domain.UserSessionService {
	return &userSessionService{
		userSessionRepo: userSessionRepo,
		sessionDuration: sessionDuration,
	}
}

func (s *userSessionService) CreateSession(userID string, ipAddress, userAgent string) (*domain.UserSession, error) {
	// Generate secure random token
	token, err := s.generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	now := time.Now()
	session := &domain.UserSession{
		UserID:    userID,
		Token:     token,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		CreatedAt: now,
		ExpiresAt: now.Add(s.sessionDuration),
	}

	if err := s.userSessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (s *userSessionService) ValidateSession(token string) (*domain.UserSession, error) {
	session, err := s.userSessionRepo.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("invalid session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		// Session expired, clean it up
		s.userSessionRepo.DeleteByToken(token)
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

func (s *userSessionService) GetUserSessions(userID string) ([]*domain.UserSession, error) {
	return s.userSessionRepo.GetByUserID(userID)
}

func (s *userSessionService) RevokeSession(userID, sessionID string) error {
	// First verify the session belongs to the user
	sessions, err := s.userSessionRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	found := false
	for _, session := range sessions {
		if session.ID == sessionID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("session not found or does not belong to user")
	}

	return s.userSessionRepo.DeleteByID(sessionID)
}

func (s *userSessionService) RevokeAllSessions(userID string) error {
	return s.userSessionRepo.DeleteByUserID(userID)
}

func (s *userSessionService) CleanupExpiredSessions() error {
	return s.userSessionRepo.DeleteExpiredSessions()
}

func (s *userSessionService) generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 256-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
