package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db    *sqlx.DB
	mock  sqlmock.Sqlmock
	redis *redis.Client
	repo  domain.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	// Create mock database
	mockDB, mock, err := sqlmock.New()
	suite.Require().NoError(err)

	suite.db = sqlx.NewDb(mockDB, "postgres")
	suite.mock = mock

	// Create mock Redis client (for testing, we'll use a real Redis client but with a test DB)
	suite.redis = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use test database
	})

	// Create repository instance
	suite.repo = NewUserRepository(suite.db, suite.redis)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	suite.db.Close()
	suite.redis.Close()
}

func (suite *UserRepositoryTestSuite) TestCreate_Success() {
	user := &domain.User{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Bio:         nil,
		AvatarURL:   nil,
	}

	// Expected SQL query
	suite.mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (username, email, display_name, bio, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`)).WithArgs(
		user.Username,
		user.Email,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
		sqlmock.AnyArg(), // created_at
		sqlmock.AnyArg(), // updated_at
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Execute
	err := suite.repo.Create(user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, user.ID)
	assert.False(suite.T(), user.CreatedAt.IsZero())
	assert.False(suite.T(), user.UpdatedAt.IsZero())

	// Verify all expectations were met
	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestCreate_DatabaseError() {
	user := &domain.User{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	// Mock database error
	suite.mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (username, email, display_name, bio, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`)).WithArgs(
		user.Username,
		user.Email,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
		sqlmock.AnyArg(),
		sqlmock.AnyArg(),
	).WillReturnError(sql.ErrConnDone)

	// Execute
	err := suite.repo.Create(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "failed to create user")

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetByID_Success() {
	userID := 1
	expectedUser := &domain.User{
		ID:          userID,
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Bio:         nil,
		AvatarURL:   nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Clear Redis cache for this test
	cacheKey := "user:1"
	suite.redis.Del(context.Background(), cacheKey)

	// Mock database query
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "display_name", "bio", "avatar_url", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID,
		expectedUser.Username,
		expectedUser.Email,
		expectedUser.DisplayName,
		expectedUser.Bio,
		expectedUser.AvatarURL,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
	)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`)).WithArgs(userID).WillReturnRows(rows)

	// Execute
	user, err := suite.repo.GetByID(userID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.ID, user.ID)
	assert.Equal(suite.T(), expectedUser.Username, user.Username)
	assert.Equal(suite.T(), expectedUser.Email, user.Email)
	assert.Equal(suite.T(), expectedUser.DisplayName, user.DisplayName)

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetByID_NotFound() {
	userID := 999

	// Clear Redis cache
	cacheKey := "user:999"
	suite.redis.Del(context.Background(), cacheKey)

	// Mock database query returning no rows
	suite.mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`)).WithArgs(userID).WillReturnError(sql.ErrNoRows)

	// Execute
	user, err := suite.repo.GetByID(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Contains(suite.T(), err.Error(), "user not found")

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestGetByEmail_Success() {
	email := "test@example.com"
	expectedUser := &domain.User{
		ID:          1,
		Username:    "testuser",
		Email:       email,
		DisplayName: "Test User",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock database query
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "display_name", "bio", "avatar_url", "created_at", "updated_at",
	}).AddRow(
		expectedUser.ID,
		expectedUser.Username,
		expectedUser.Email,
		expectedUser.DisplayName,
		expectedUser.Bio,
		expectedUser.AvatarURL,
		expectedUser.CreatedAt,
		expectedUser.UpdatedAt,
	)

	suite.mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, username, email, display_name, bio, avatar_url, created_at, updated_at
		FROM users
		WHERE email = $1
	`)).WithArgs(email).WillReturnRows(rows)

	// Execute
	user, err := suite.repo.GetByEmail(email)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), expectedUser.Email, user.Email)

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdate_Success() {
	user := &domain.User{
		ID:          1,
		Username:    "updateduser",
		Email:       "updated@example.com",
		DisplayName: "Updated User",
	}

	// Mock successful update
	suite.mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users 
		SET username = $2, email = $3, display_name = $4, bio = $5, avatar_url = $6, updated_at = $7
		WHERE id = $1
	`)).WithArgs(
		user.ID,
		user.Username,
		user.Email,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
		sqlmock.AnyArg(), // updated_at
	).WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Execute
	err := suite.repo.Update(user)

	// Assert
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), user.UpdatedAt.IsZero())

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestUpdate_NotFound() {
	user := &domain.User{
		ID:          999,
		Username:    "updateduser",
		Email:       "updated@example.com",
		DisplayName: "Updated User",
	}

	// Mock update with no rows affected
	suite.mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users 
		SET username = $2, email = $3, display_name = $4, bio = $5, avatar_url = $6, updated_at = $7
		WHERE id = $1
	`)).WithArgs(
		user.ID,
		user.Username,
		user.Email,
		user.DisplayName,
		user.Bio,
		user.AvatarURL,
		sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	// Execute
	err := suite.repo.Update(user)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestDelete_Success() {
	userID := 1

	// Mock successful deletion
	suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Execute
	err := suite.repo.Delete(userID)

	// Assert
	assert.NoError(suite.T(), err)

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

func (suite *UserRepositoryTestSuite) TestDelete_NotFound() {
	userID := 999

	// Mock deletion with no rows affected
	suite.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM users WHERE id = $1`)).
		WithArgs(userID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	// Execute
	err := suite.repo.Delete(userID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "user not found")

	suite.Require().NoError(suite.mock.ExpectationsWereMet())
}

// Run the test suite
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
