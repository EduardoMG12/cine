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
	UserID1   int              `json:"user_id_1" db:"user_id_1"`
	UserID2   int              `json:"user_id_2" db:"user_id_2"`
	Status    FriendshipStatus `json:"status" db:"status"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`

	User1 *User `json:"user_1,omitempty"`
	User2 *User `json:"user_2,omitempty"`
}

type Follow struct {
	FollowerID  int       `json:"follower_id" db:"follower_id"`
	FollowingID int       `json:"following_id" db:"following_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`

	Follower  *User `json:"follower,omitempty"`
	Following *User `json:"following,omitempty"`
}

type FriendshipRepository interface {
	Create(friendship *Friendship) error
	GetByUsers(userID1, userID2 int) (*Friendship, error)
	GetUserFriends(userID int) ([]*Friendship, error)
	GetFriendRequests(userID int) ([]*Friendship, error)
	UpdateStatus(userID1, userID2 int, status FriendshipStatus) error
	Delete(userID1, userID2 int) error
}

type FollowRepository interface {
	Create(follow *Follow) error
	GetByUsers(followerID, followingID int) (*Follow, error)
	GetFollowers(userID int, limit, offset int) ([]*Follow, error)
	GetFollowing(userID int, limit, offset int) ([]*Follow, error)
	Delete(followerID, followingID int) error
	GetFollowersCount(userID int) (int, error)
	GetFollowingCount(userID int) (int, error)
}

type SocialService interface {
	SendFriendRequest(senderID, receiverID int) error
	AcceptFriendRequest(userID, requesterID int) error
	DeclineFriendRequest(userID, requesterID int) error
	RemoveFriend(userID, friendID int) error
	BlockUser(userID, blockedUserID int) error
	GetFriends(userID int) ([]*User, error)
	GetFriendRequests(userID int) ([]*User, error)
	AreFriends(userID1, userID2 int) (bool, error)

	FollowUser(followerID, followingID int) error
	UnfollowUser(followerID, followingID int) error
	GetFollowers(userID int, page int) ([]*User, error)
	GetFollowing(userID int, page int) ([]*User, error)
	IsFollowing(followerID, followingID int) (bool, error)
	GetFollowersCount(userID int) (int, error)
	GetFollowingCount(userID int) (int, error)
}
