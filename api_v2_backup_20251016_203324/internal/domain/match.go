package domain

import (
	"time"
)

type MatchSessionStatus string

const (
	MatchSessionStatusActive   MatchSessionStatus = "active"
	MatchSessionStatusFinished MatchSessionStatus = "finished"
	MatchSessionStatusCanceled MatchSessionStatus = "canceled"
)

type MatchSession struct {
	ID         string             `json:"id" db:"id"`
	HostUserID string             `json:"host_user_id" db:"host_user_id"`
	Status     MatchSessionStatus `json:"status" db:"status"`
	CreatedAt  time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" db:"updated_at"`

	Host         *User                      `json:"host,omitempty"`
	Participants []*MatchSessionParticipant `json:"participants,omitempty"`
	Interactions []*MatchInteraction        `json:"interactions,omitempty"`
}

type MatchSessionParticipant struct {
	SessionID string    `json:"session_id" db:"session_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`

	User *User `json:"user,omitempty"`
}

type MatchInteraction struct {
	SessionID    string    `json:"session_id" db:"session_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	MovieID      string    `json:"movie_id" db:"movie_id"`
	Liked        bool      `json:"liked" db:"liked"`
	InteractedAt time.Time `json:"interacted_at" db:"interacted_at"`

	User  *User  `json:"user,omitempty"`
	Movie *Movie `json:"movie,omitempty"`
}

type MatchResult struct {
	SessionID string `json:"session_id"`
	MovieID   string `json:"movie_id"`
	Movie     *Movie `json:"movie,omitempty"`
}

type MatchSessionRepository interface {
	Create(session *MatchSession) error
	GetByID(id string) (*MatchSession, error)
	GetByHostID(hostID string, limit, offset int) ([]*MatchSession, error)
	GetByParticipantID(userID string, limit, offset int) ([]*MatchSession, error)
	UpdateStatus(sessionID string, status MatchSessionStatus) error
	Delete(id string) error

	AddParticipant(sessionID, userID string) error
	RemoveParticipant(sessionID, userID string) error
	GetParticipants(sessionID string) ([]*MatchSessionParticipant, error)
	IsParticipant(sessionID, userID string) (bool, error)

	CreateInteraction(interaction *MatchInteraction) error
	GetInteractions(sessionID string) ([]*MatchInteraction, error)
	GetUserInteractions(sessionID, userID string) ([]*MatchInteraction, error)
	HasUserInteracted(sessionID, userID, movieID string) (bool, error)
}

type MatchService interface {
	CreateSession(hostUserID string, participantIDs []string) (*MatchSession, error)
	JoinSession(sessionID, userID string) error
	LeaveSession(sessionID, userID string) error
	GetSession(sessionID string) (*MatchSession, error)
	GetUserSessions(userID string, page int) ([]*MatchSession, error)

	GetSessionSuggestions(sessionID string, limit int) ([]*Movie, error)

	RecordInteraction(sessionID, userID, movieID string, liked bool) error
	CheckForMatches(sessionID string) ([]*MatchResult, error)

	FinishSession(sessionID, userID string) error
	CancelSession(sessionID, userID string) error
}
