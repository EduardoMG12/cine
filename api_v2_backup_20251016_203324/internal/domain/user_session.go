package domain

import (
	"time"
)

type UserSession struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

type UserSessionRepository interface {
	Create(session *UserSession) error
	GetByToken(token string) (*UserSession, error)
	GetByUserID(userID string) ([]*UserSession, error)
	DeleteByID(id string) error
	DeleteByToken(token string) error
	DeleteByUserID(userID string) error
	DeleteExpiredSessions() error
}

type UserSessionService interface {
	CreateSession(userID string, ipAddress, userAgent string) (*UserSession, error)
	ValidateSession(token string) (*UserSession, error)
	GetUserSessions(userID string) ([]*UserSession, error)
	RevokeSession(userID, sessionID string) error
	RevokeAllSessions(userID string) error
	CleanupExpiredSessions() error
}
