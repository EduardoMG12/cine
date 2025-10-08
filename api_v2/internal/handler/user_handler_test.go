package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of domain.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(username, email, displayName string) (*domain.User, error) {
	args := m.Called(username, email, displayName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetUser(id int) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	args := m.Called(id, updates)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) ValidateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	createRequest := CreateUserRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	expectedUser := &domain.User{
		ID:          1,
		Username:    createRequest.Username,
		Email:       createRequest.Email,
		DisplayName: createRequest.DisplayName,
	}

	mockService.On("CreateUser", createRequest.Username, createRequest.Email, createRequest.DisplayName).
		Return(expectedUser, nil)

	// Create request
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response domain.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Username, response.Username)
	assert.Equal(t, expectedUser.Email, response.Email)

	mockService.AssertExpectations(t)
}

func TestUserHandler_CreateUser_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Invalid JSON")
}

func TestUserHandler_CreateUser_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	createRequest := CreateUserRequest{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	mockService.On("CreateUser", createRequest.Username, createRequest.Email, createRequest.DisplayName).
		Return(nil, errors.New("username already exists"))

	// Create request
	reqBody, _ := json.Marshal(createRequest)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.CreateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusConflict, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "username already exists")

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_Success(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := 1
	expectedUser := &domain.User{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	mockService.On("GetUser", userID).Return(expectedUser, nil)

	// Create request with URL parameter
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	// Setup chi context with URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(userID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, response.ID)
	assert.Equal(t, expectedUser.Username, response.Username)

	mockService.AssertExpectations(t)
}

func TestUserHandler_GetUser_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Create request with invalid ID
	req := httptest.NewRequest(http.MethodGet, "/users/invalid", nil)
	w := httptest.NewRecorder()

	// Setup chi context with invalid ID
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "Invalid user ID")
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := 999

	mockService.On("GetUser", userID).Return(nil, errors.New("user not found"))

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(userID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.GetUser(w, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Error, "User not found")

	mockService.AssertExpectations(t)
}

func TestUserHandler_UpdateUser_Success(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := 1
	updateRequest := UpdateUserRequest{
		Username:    stringPtr("updateduser"),
		Email:       stringPtr("updated@example.com"),
		DisplayName: stringPtr("Updated User"),
	}

	expectedUser := &domain.User{
		ID:          userID,
		Username:    *updateRequest.Username,
		Email:       *updateRequest.Email,
		DisplayName: *updateRequest.DisplayName,
	}

	expectedUpdates := map[string]interface{}{
		"username":     *updateRequest.Username,
		"email":        *updateRequest.Email,
		"display_name": *updateRequest.DisplayName,
	}

	mockService.On("UpdateUser", userID, expectedUpdates).Return(expectedUser, nil)

	// Create request
	reqBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(userID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.UpdateUser(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Username, response.Username)

	mockService.AssertExpectations(t)
}

func TestUserHandler_DeleteUser_Success(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := 1

	mockService.On("DeleteUser", userID).Return(nil)

	// Create request
	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()

	// Setup chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", strconv.Itoa(userID))
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	// Act
	handler.DeleteUser(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())

	mockService.AssertExpectations(t)
}

func TestUserHandler_Routes(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Act
	router := handler.Routes()

	// Assert
	assert.NotNil(t, router)
	// Simply verify that the router was created without errors
}

// Helper function for pointer to string
func stringPtr(s string) *string {
	return &s
}
