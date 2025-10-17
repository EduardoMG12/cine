package service

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type passwordService struct {
	cost int
}

func NewPasswordService() domain.PasswordService {
	return &passwordService{
		cost: 12,
	}
}

func (s *passwordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

func (s *passwordService) ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
