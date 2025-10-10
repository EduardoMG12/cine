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
	ID         int                `json:"id" db:"id"`
	HostUserID int                `json:"host_user_id" db:"host_user_id"`
	Status     MatchSessionStatus `json:"status" db:"status"`
	CreatedAt  time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" db:"updated_at"`

	// Populated by joins
	Host         *User                      `json:"host,omitempty"`
	Participants []*MatchSessionParticipant `json:"participants,omitempty"`
	Interactions []*MatchInteraction        `json:"interactions,omitempty"`
}

type MatchSessionParticipant struct {
	SessionID int       `json:"session_id" db:"session_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	JoinedAt  time.Time `json:"joined_at" db:"joined_at"`

	// Populated by joins
	User *User `json:"user,omitempty"`
}

type MatchInteraction struct {
	SessionID    int       `json:"session_id" db:"session_id"`
	UserID       int       `json:"user_id" db:"user_id"`
	MovieID      int       `json:"movie_id" db:"movie_id"`
	Liked        bool      `json:"liked" db:"liked"`
	InteractedAt time.Time `json:"interacted_at" db:"interacted_at"`

	// Populated by joins
	User  *User  `json:"user,omitempty"`
	Movie *Movie `json:"movie,omitempty"`
}

type MatchResult struct {
	SessionID int    `json:"session_id"`
	MovieID   int    `json:"movie_id"`
	Movie     *Movie `json:"movie,omitempty"`
}

type MatchSessionRepository interface {
	Create(session *MatchSession) error
	GetByID(id int) (*MatchSession, error)
	GetByHostID(hostID int, limit, offset int) ([]*MatchSession, error)
	GetByParticipantID(userID int, limit, offset int) ([]*MatchSession, error)
	UpdateStatus(sessionID int, status MatchSessionStatus) error
	Delete(id int) error

	// Participant operations
	AddParticipant(sessionID, userID int) error
	RemoveParticipant(sessionID, userID int) error
	GetParticipants(sessionID int) ([]*MatchSessionParticipant, error)
	IsParticipant(sessionID, userID int) (bool, error)

	// Interaction operations
	CreateInteraction(interaction *MatchInteraction) error
	GetInteractions(sessionID int) ([]*MatchInteraction, error)
	GetUserInteractions(sessionID, userID int) ([]*MatchInteraction, error)
	HasUserInteracted(sessionID, userID, movieID int) (bool, error)
}

type MatchService interface {
	CreateSession(hostUserID int, participantIDs []int) (*MatchSession, error)
	JoinSession(sessionID, userID int) error
	LeaveSession(sessionID, userID int) error
	GetSession(sessionID int) (*MatchSession, error)
	GetUserSessions(userID int, page int) ([]*MatchSession, error)

	// Movie suggestions based on participants' preferences
	GetSessionSuggestions(sessionID int, limit int) ([]*Movie, error)

	// Interaction operations
	RecordInteraction(sessionID, userID, movieID int, liked bool) error
	CheckForMatches(sessionID int) ([]*MatchResult, error)

	// Session management
	FinishSession(sessionID, userID int) error
	CancelSession(sessionID, userID int) error
}
