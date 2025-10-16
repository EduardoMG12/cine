package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/EduardoMG12/cine/api_v2/internal/utils"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type userSessionRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

func NewUserSessionRepository(db *sqlx.DB, redis *redis.Client) domain.UserSessionRepository {
	return &userSessionRepository{
		db:    db,
		redis: redis,
	}
}

func (r *userSessionRepository) Create(session *domain.UserSession) error {
	// Generate UUID for the session if not set
	if session.ID == "" {
		session.ID = utils.GenerateUUID()
	}

	query := `
		INSERT INTO user_sessions (id, user_id, token, ip_address, user_agent, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		query,
		session.ID,
		session.UserID,
		session.Token,
		session.IPAddress,
		session.UserAgent,
		session.CreatedAt,
		session.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user session: %w", err)
	}

	cacheKey := fmt.Sprintf("session:%s", session.Token)
	sessionData := fmt.Sprintf("%s:%s", session.UserID, session.ID)

	ctx := context.Background()
	r.redis.Set(ctx, cacheKey, sessionData, time.Until(session.ExpiresAt))

	return nil
}

func (r *userSessionRepository) GetByToken(token string) (*domain.UserSession, error) {
	query := `
		SELECT id, user_id, token, ip_address, user_agent, created_at, expires_at
		FROM user_sessions
		WHERE token = $1 AND expires_at > NOW()
	`

	var session domain.UserSession
	err := r.db.Get(&session, query, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user session by token: %w", err)
	}

	return &session, nil
}

func (r *userSessionRepository) GetByUserID(userID string) ([]*domain.UserSession, error) {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return nil, fmt.Errorf("invalid UUID format: %s", userID)
	}

	query := `
		SELECT id, user_id, token, ip_address, user_agent, created_at, expires_at
		FROM user_sessions
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`

	var sessions []*domain.UserSession
	err := r.db.Select(&sessions, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}

	return sessions, nil
}

func (r *userSessionRepository) DeleteByID(id string) error {
	// Validate UUID format
	if !utils.IsValidUUID(id) {
		return fmt.Errorf("invalid UUID format: %s", id)
	}

	query := `DELETE FROM user_sessions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *userSessionRepository) DeleteByToken(token string) error {
	query := `DELETE FROM user_sessions WHERE token = $1`

	result, err := r.db.Exec(query, token)
	if err != nil {
		return fmt.Errorf("failed to delete session by token: %w", err)
	}

	cacheKey := fmt.Sprintf("session:%s", token)
	r.redis.Del(context.Background(), cacheKey)

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *userSessionRepository) DeleteByUserID(userID string) error {
	// Validate UUID format
	if !utils.IsValidUUID(userID) {
		return fmt.Errorf("invalid UUID format: %s", userID)
	}

	sessions, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	ctx := context.Background()
	for _, session := range sessions {
		cacheKey := fmt.Sprintf("session:%s", session.Token)
		r.redis.Del(ctx, cacheKey)
	}

	query := `DELETE FROM user_sessions WHERE user_id = $1`
	_, err = r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete sessions by user ID: %w", err)
	}

	return nil
}

func (r *userSessionRepository) DeleteExpiredSessions() error {
	query := `DELETE FROM user_sessions WHERE expires_at <= NOW()`

	result, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		fmt.Printf("Deleted %d expired sessions\n", rowsAffected)
	}

	return nil
}
