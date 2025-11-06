package auth

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type LogoutUseCase struct {
	sessionRepo domain.SessionRepository
}

func NewLogoutUseCase(sessionRepo domain.SessionRepository) *LogoutUseCase {
	return &LogoutUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *LogoutUseCase) Execute(token string) error {
	if err := uc.sessionRepo.DeleteSession(token); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}
