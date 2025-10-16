package domain

import (
	"time"
)

type FriendshipStatus string

const (
	FriendshipStatusPending  FriendshipStatus = "pending"
	FriendshipStatusAccepted FriendshipStatus = "accepted"
	FriendshipStatusDeclined FriendshipStatus = "declined"
	FriendshipStatusBlocked  FriendshipStatus = "blocked"
)

type Friendship struct {
	UserID1   string           `json:"user_id_1" db:"user_id_1"`
	UserID2   string           `json:"user_id_2" db:"user_id_2"`
	Status    FriendshipStatus `json:"status" db:"status"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`

	User1 *User `json:"user_1,omitempty"`
	User2 *User `json:"user_2,omitempty"`
}

type Follow struct {
	FollowerID  string    `json:"follower_id" db:"follower_id"`
	FollowingID string    `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	Follower  *User `json:"follower,omitempty"`
	Following *User `json:"following,omitempty"`
}

type FriendshipRepository interface {
	Create(friendship *Friendship) error
	GetByUsers(userID1, userID2 string) (*Friendship, error)
	GetUserFriends(userID string) ([]*Friendship, error)
	GetFriendRequests(userID string) ([]*Friendship, error)
	UpdateStatus(userID1, userID2 string, status FriendshipStatus) error
	Delete(userID1, userID2 string) error
}

type FollowRepository interface {
	Create(follow *Follow) error
	GetByUsers(followerID, followingID string) (*Follow, error)
	GetFollowers(userID string, limit, offset int) ([]*Follow, error)
	GetFollowing(userID string, limit, offset int) ([]*Follow, error)
	Delete(followerID, followingID string) error
	GetFollowersCount(userID string) (int, error)
	GetFollowingCount(userID string) (int, error)
}

type SocialService interface {
	SendFriendRequest(senderID, receiverID string) error
	AcceptFriendRequest(userID, requesterID string) error
	DeclineFriendRequest(userID, requesterID string) error
	RemoveFriend(userID, friendID string) error
	BlockUser(userID, blockedUserID string) error
	GetFriends(userID string) ([]*User, error)
	GetFriendRequests(userID string) ([]*User, error)
	AreFriends(userID1, userID2 string) (bool, error)

	FollowUser(followerID, followingID string) error
	UnfollowUser(followerID, followingID string) error
	GetFollowers(userID string, page int) ([]*User, error)
	GetFollowing(userID string, page int) ([]*User, error)
	IsFollowing(followerID, followingID string) (bool, error)
	GetFollowersCount(userID string) (int, error)
	GetFollowingCount(userID string) (int, error)
}
