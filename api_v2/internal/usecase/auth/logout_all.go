package auth

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
)

type LogoutAllUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewLogoutAllUseCase(sessionRepo domain.SessionRepository) *LogoutAllUseCase {
	return &LogoutAllUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutAllUseCase) Execute(userID uuid.UUID) error {
	if err := uc.sessionRepo.DeleteUserSessions(userID); err != nil {
		return fmt.Errorf("failed to logout from all sessions: %w", err)
	}
	return nil
}
