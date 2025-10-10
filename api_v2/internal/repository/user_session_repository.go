package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
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
	query := `
		INSERT INTO user_sessions (user_id, token, ip_address, user_agent, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		session.UserID,
		session.Token,
		session.IPAddress,
		session.UserAgent,
		session.CreatedAt,
		session.ExpiresAt,
	).Scan(&session.ID)

	if err != nil {
		return fmt.Errorf("failed to create user session: %w", err)
	}

	// Cache session in Redis for fast lookups
	cacheKey := fmt.Sprintf("session:%s", session.Token)
	sessionData := fmt.Sprintf("%d:%d", session.UserID, session.ID)

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

func (r *userSessionRepository) GetByUserID(userID int) ([]*domain.UserSession, error) {
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

func (r *userSessionRepository) DeleteByID(id int) error {
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

	// Remove from cache
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

func (r *userSessionRepository) DeleteByUserID(userID int) error {
	// First get all tokens to remove from cache
	sessions, err := r.GetByUserID(userID)
	if err != nil {
		return err
	}

	// Remove from cache
	ctx := context.Background()
	for _, session := range sessions {
		cacheKey := fmt.Sprintf("session:%s", session.Token)
		r.redis.Del(ctx, cacheKey)
	}

	// Delete from database
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
