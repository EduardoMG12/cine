package service

import (
	"errors"
	"testing"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	username := "testuser"
	email := "test@example.com"
	displayName := "Test User"

	// Mock repository calls - user doesn't exist
	mockRepo.On("GetByUsername", username).Return(nil, errors.New("user not found"))
	mockRepo.On("GetByEmail", email).Return(nil, errors.New("user not found"))
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	// Act
	user, err := service.CreateUser(username, email, displayName)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, displayName, user.DisplayName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_DuplicateUsername(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	username := "testuser"
	email := "test@example.com"
	displayName := "Test User"

	existingUser := &domain.User{
		ID:       1,
		Username: username,
		Email:    "other@example.com",
	}

	// Mock repository calls - username already exists
	mockRepo.On("GetByUsername", username).Return(existingUser, nil)

	// Act
	user, err := service.CreateUser(username, email, displayName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "username already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	username := "testuser"
	email := "test@example.com"
	displayName := "Test User"

	existingUser := &domain.User{
		ID:       1,
		Username: "otheruser",
		Email:    email,
	}

	// Mock repository calls - email already exists
	mockRepo.On("GetByUsername", username).Return(nil, errors.New("user not found"))
	mockRepo.On("GetByEmail", email).Return(existingUser, nil)

	// Act
	user, err := service.CreateUser(username, email, displayName)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "email already exists")
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser_InvalidUsername(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	testCases := []struct {
		name     string
		username string
		email    string
		display  string
	}{
		{"empty username", "", "test@example.com", "Test User"},
		{"username too short", "ab", "test@example.com", "Test User"},
		{"username too long", "thisusernameiswaytoolongandexceedsthemaximumlength", "test@example.com", "Test User"},
		{"username with spaces", "test user", "test@example.com", "Test User"},
		{"username with special chars", "test@user", "test@example.com", "Test User"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			user, err := service.CreateUser(tc.username, tc.email, tc.display)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, user)
		})
	}
}

func TestUserService_CreateUser_InvalidEmail(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	testCases := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"invalid format", "invalid-email"},
		{"missing @", "testexample.com"},
		{"missing domain", "test@"},
		{"missing local part", "@example.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			user, err := service.CreateUser("testuser", tc.email, "Test User")

			// Assert
			assert.Error(t, err)
			assert.Nil(t, user)
		})
	}
}

func TestUserService_GetUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := 1
	expectedUser := &domain.User{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	mockRepo.On("GetByID", userID).Return(expectedUser, nil)

	// Act
	user, err := service.GetUser(userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Username, user.Username)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUser_InvalidID(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	testCases := []int{0, -1, -100}

	for _, invalidID := range testCases {
		t.Run("invalid ID", func(t *testing.T) {
			// Act
			user, err := service.GetUser(invalidID)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, user)
			assert.Contains(t, err.Error(), "invalid user ID")
		})
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := 999

	mockRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

	// Act
	user, err := service.GetUser(userID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "failed to get user")
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := 1
	existingUser := &domain.User{
		ID:          userID,
		Username:    "olduser",
		Email:       "old@example.com",
		DisplayName: "Old User",
	}

	updates := map[string]interface{}{
		"username":     "newuser",
		"email":        "new@example.com",
		"display_name": "New User",
	}

	mockRepo.On("GetByID", userID).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.User")).Return(nil)

	// Act
	user, err := service.UpdateUser(userID, updates)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "New User", user.DisplayName)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := 1

	mockRepo.On("Delete", userID).Return(nil)

	// Act
	err := service.DeleteUser(userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUserService_DeleteUser_InvalidID(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Act
	err := service.DeleteUser(0)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID")
}
