package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) domain.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) CreateSession(session *domain.UserSession) error {
	session.ID = uuid.New()
	session.CreatedAt = time.Now()

	query := `
		INSERT INTO user_sessions (
			id, user_id, token, expires_at, created_at, user_agent, ip_address
		) VALUES (
			:id, :user_id, :token, :expires_at, :created_at, :user_agent, :ip_address
		)
	`

	_, err := r.db.NamedExec(query, session)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *sessionRepository) GetSessionByToken(token string) (*domain.UserSession, error) {
	var session domain.UserSession
	query := `
		SELECT id, user_id, token, expires_at, created_at, user_agent, ip_address
		FROM user_sessions 
		WHERE token = $1 AND expires_at > NOW()
	`

	err := r.db.Get(&session, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session by token: %w", err)
	}

	return &session, nil
}

func (r *sessionRepository) DeleteSession(token string) error {
	query := `DELETE FROM user_sessions WHERE token = $1`

	result, err := r.db.Exec(query, token)
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

func (r *sessionRepository) DeleteUserSessions(userID uuid.UUID) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user sessions: %w", err)
	}

	return nil
}
