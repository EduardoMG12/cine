package infrastructure

import "golang.org/x/crypto/bcrypt"

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *PasswordService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
