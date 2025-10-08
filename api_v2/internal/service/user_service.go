package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type userService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) domain.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(username, email, displayName string) (*domain.User, error) {
	if err := s.validateUsername(username); err != nil {
		return nil, err
	}

	if err := s.validateEmail(email); err != nil {
		return nil, err
	}

	if err := s.validateDisplayName(displayName); err != nil {
		return nil, err
	}

	existingUser, err := s.userRepo.GetByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, err = s.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user := &domain.User{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if username, ok := updates["username"].(string); ok {
		if err := s.validateUsername(username); err != nil {
			return nil, err
		}
		user.Username = username
	}

	if email, ok := updates["email"].(string); ok {
		if err := s.validateEmail(email); err != nil {
			return nil, err
		}
		user.Email = email
	}

	if displayName, ok := updates["display_name"].(string); ok {
		if err := s.validateDisplayName(displayName); err != nil {
			return nil, err
		}
		user.DisplayName = displayName
	}

	if bio, ok := updates["bio"].(string); ok {
		if len(bio) > 500 {
			return nil, errors.New("bio must be less than 500 characters")
		}
		user.Bio = &bio
	}

	if avatarURL, ok := updates["avatar_url"].(string); ok {
		user.AvatarURL = &avatarURL
	}

	if err := s.ValidateUser(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUser(id int) error {
	if id <= 0 {
		return errors.New("invalid user ID")
	}

	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *userService) ValidateUser(user *domain.User) error {
	if err := s.validateUsername(user.Username); err != nil {
		return err
	}

	if err := s.validateEmail(user.Email); err != nil {
		return err
	}

	if err := s.validateDisplayName(user.DisplayName); err != nil {
		return err
	}

	if user.Bio != nil && len(*user.Bio) > 500 {
		return errors.New("bio must be less than 500 characters")
	}

	return nil
}

func (s *userService) validateUsername(username string) error {
	if len(username) < 3 || len(username) > 30 {
		return errors.New("username must be between 3 and 30 characters")
	}

	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	if !matched {
		return errors.New("username can only contain letters, numbers and underscores")
	}

	return nil
}

func (s *userService) validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func (s *userService) validateDisplayName(displayName string) error {
	displayName = strings.TrimSpace(displayName)
	if len(displayName) < 1 || len(displayName) > 100 {
		return errors.New("display name must be between 1 and 100 characters")
	}

	return nil
}
