package service

import (
	"fmt"

	"github.com/EduardoMG12/cine/api_v2/internal/domain"
)

type socialService struct {
	friendshipRepo domain.FriendshipRepository
	followRepo     domain.FollowRepository
	userRepo       domain.UserRepository
}

func NewSocialService(
	friendshipRepo domain.FriendshipRepository,
	followRepo domain.FollowRepository,
	userRepo domain.UserRepository,
) domain.SocialService {
	return &socialService{
		friendshipRepo: friendshipRepo,
		followRepo:     followRepo,
		userRepo:       userRepo,
	}
}

// Friendship operations

func (s *socialService) SendFriendRequest(senderID, receiverID string) error {
	if senderID == receiverID {
		return fmt.Errorf("cannot send friend request to yourself")
	}

	// Check if users exist
	sender, err := s.userRepo.GetByID(senderID)
	if err != nil || sender == nil {
		return fmt.Errorf("sender not found")
	}

	receiver, err := s.userRepo.GetByID(receiverID)
	if err != nil || receiver == nil {
		return fmt.Errorf("receiver not found")
	}

	// Check if friendship already exists
	existing, err := s.friendshipRepo.GetByUsers(senderID, receiverID)
	if err != nil {
		return err
	}

	if existing != nil {
		switch existing.Status {
		case domain.FriendshipStatusAccepted:
			return fmt.Errorf("users are already friends")
		case domain.FriendshipStatusPending:
			return fmt.Errorf("friend request already pending")
		case domain.FriendshipStatusBlocked:
			return fmt.Errorf("cannot send friend request to blocked user")
		}
	}

	// Create new friendship request
	friendship := &domain.Friendship{
		UserID1: senderID,
		UserID2: receiverID,
		Status:  domain.FriendshipStatusPending,
	}

	return s.friendshipRepo.Create(friendship)
}

func (s *socialService) AcceptFriendRequest(userID, requesterID string) error {
	// Check if friend request exists and is pending
	friendship, err := s.friendshipRepo.GetByUsers(requesterID, userID)
	if err != nil {
		return err
	}

	if friendship == nil {
		return fmt.Errorf("no friend request found")
	}

	if friendship.Status != domain.FriendshipStatusPending {
		return fmt.Errorf("friend request is not pending")
	}

	// Ensure the request was sent TO the current user
	if friendship.UserID2 != userID {
		return fmt.Errorf("unauthorized to accept this friend request")
	}

	return s.friendshipRepo.UpdateStatus(requesterID, userID, domain.FriendshipStatusAccepted)
}

func (s *socialService) DeclineFriendRequest(userID, requesterID string) error {
	// Check if friend request exists and is pending
	friendship, err := s.friendshipRepo.GetByUsers(requesterID, userID)
	if err != nil {
		return err
	}

	if friendship == nil {
		return fmt.Errorf("no friend request found")
	}

	if friendship.Status != domain.FriendshipStatusPending {
		return fmt.Errorf("friend request is not pending")
	}

	// Ensure the request was sent TO the current user
	if friendship.UserID2 != userID {
		return fmt.Errorf("unauthorized to decline this friend request")
	}

	// Delete the friendship record instead of updating to declined
	return s.friendshipRepo.Delete(requesterID, userID)
}

func (s *socialService) RemoveFriend(userID, friendID string) error {
	// Check if friendship exists and is accepted
	friendship, err := s.friendshipRepo.GetByUsers(userID, friendID)
	if err != nil {
		return err
	}

	if friendship == nil || friendship.Status != domain.FriendshipStatusAccepted {
		return fmt.Errorf("users are not friends")
	}

	return s.friendshipRepo.Delete(userID, friendID)
}

func (s *socialService) BlockUser(userID, blockedUserID string) error {
	if userID == blockedUserID {
		return fmt.Errorf("cannot block yourself")
	}

	// Check if friendship exists
	existing, err := s.friendshipRepo.GetByUsers(userID, blockedUserID)
	if err != nil {
		return err
	}

	if existing != nil {
		// Update existing friendship to blocked
		return s.friendshipRepo.UpdateStatus(userID, blockedUserID, domain.FriendshipStatusBlocked)
	} else {
		// Create new blocked relationship
		friendship := &domain.Friendship{
			UserID1: userID,
			UserID2: blockedUserID,
			Status:  domain.FriendshipStatusBlocked,
		}
		return s.friendshipRepo.Create(friendship)
	}
}

func (s *socialService) GetFriends(userID string) ([]*domain.User, error) {
	friendships, err := s.friendshipRepo.GetUserFriends(userID)
	if err != nil {
		return nil, err
	}

	var friends []*domain.User
	for _, friendship := range friendships {
		// Return the friend (not the current user)
		if friendship.User1.ID == userID {
			friends = append(friends, friendship.User2)
		} else {
			friends = append(friends, friendship.User1)
		}
	}

	return friends, nil
}

func (s *socialService) GetFriendRequests(userID string) ([]*domain.User, error) {
	friendships, err := s.friendshipRepo.GetFriendRequests(userID)
	if err != nil {
		return nil, err
	}

	var requesters []*domain.User
	for _, friendship := range friendships {
		requesters = append(requesters, friendship.User1)
	}

	return requesters, nil
}

func (s *socialService) AreFriends(userID1, userID2 string) (bool, error) {
	friendship, err := s.friendshipRepo.GetByUsers(userID1, userID2)
	if err != nil {
		return false, err
	}

	return friendship != nil && friendship.Status == domain.FriendshipStatusAccepted, nil
}

// Follow operations

func (s *socialService) FollowUser(followerID, followingID string) error {
	if followerID == followingID {
		return fmt.Errorf("cannot follow yourself")
	}

	// Check if users exist
	follower, err := s.userRepo.GetByID(followerID)
	if err != nil || follower == nil {
		return fmt.Errorf("follower not found")
	}

	following, err := s.userRepo.GetByID(followingID)
	if err != nil || following == nil {
		return fmt.Errorf("user to follow not found")
	}

	// Check if already following
	existing, err := s.followRepo.GetByUsers(followerID, followingID)
	if err != nil {
		return err
	}

	if existing != nil {
		return fmt.Errorf("already following this user")
	}

	// Create follow relationship
	follow := &domain.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.followRepo.Create(follow)
}

func (s *socialService) UnfollowUser(followerID, followingID string) error {
	// Check if follow relationship exists
	existing, err := s.followRepo.GetByUsers(followerID, followingID)
	if err != nil {
		return err
	}

	if existing == nil {
		return fmt.Errorf("not following this user")
	}

	return s.followRepo.Delete(followerID, followingID)
}

func (s *socialService) GetFollowers(userID string, page int) ([]*domain.User, error) {
	limit := 20 // Fixed page size
	offset := (page - 1) * limit

	follows, err := s.followRepo.GetFollowers(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var followers []*domain.User
	for _, follow := range follows {
		followers = append(followers, follow.Follower)
	}

	return followers, nil
}

func (s *socialService) GetFollowing(userID string, page int) ([]*domain.User, error) {
	limit := 20 // Fixed page size
	offset := (page - 1) * limit

	follows, err := s.followRepo.GetFollowing(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var following []*domain.User
	for _, follow := range follows {
		following = append(following, follow.Following)
	}

	return following, nil
}

func (s *socialService) IsFollowing(followerID, followingID string) (bool, error) {
	follow, err := s.followRepo.GetByUsers(followerID, followingID)
	if err != nil {
		return false, err
	}

	return follow != nil, nil
}

func (s *socialService) GetFollowersCount(userID string) (int, error) {
	return s.followRepo.GetFollowersCount(userID)
}

func (s *socialService) GetFollowingCount(userID string) (int, error) {
	return s.followRepo.GetFollowingCount(userID)
}
