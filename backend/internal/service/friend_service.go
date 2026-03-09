package service

import (
	"errors"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrFriendRequestExists = errors.New("friend request already exists")
	ErrFriendNotFound      = errors.New("friend request not found")
	ErrCannotFriendSelf    = errors.New("cannot send friend request to yourself")
)

type FriendService interface {
	SendFriendRequest(userID, friendID string) (*model.Friend, error)
	AcceptFriendRequest(userID, requestID string) (*model.Friend, error)
	RejectFriendRequest(userID, requestID string) error
	GetFriends(userID string) ([]model.Friend, error)
	BlockUser(userID, friendID string) error
}

type friendService struct {
	friendRepo repository.FriendRepository
}

func NewFriendService(friendRepo repository.FriendRepository) FriendService {
	return &friendService{friendRepo: friendRepo}
}

func (s *friendService) SendFriendRequest(userID, friendID string) (*model.Friend, error) {
	if userID == friendID {
		return nil, ErrCannotFriendSelf
	}

	existing, err := s.friendRepo.FindByUserAndFriend(userID, friendID)
	if err == nil && existing != nil {
		return nil, ErrFriendRequestExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	friend := &model.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   model.FriendStatusPending,
	}

	if err := s.friendRepo.Create(friend); err != nil {
		return nil, err
	}

	return friend, nil
}

func (s *friendService) AcceptFriendRequest(userID, requestID string) (*model.Friend, error) {
	request, err := s.friendRepo.FindByID(requestID)
	if err != nil {
		return nil, ErrFriendNotFound
	}

	if request.FriendID != userID {
		return nil, ErrFriendNotFound
	}

	if err := s.friendRepo.UpdateStatus(requestID, model.FriendStatusAccepted); err != nil {
		return nil, err
	}

	// Create reverse friendship
	reverse := &model.Friend{
		UserID:   userID,
		FriendID: request.UserID,
		Status:   model.FriendStatusAccepted,
	}
	_ = s.friendRepo.Create(reverse)

	return s.friendRepo.FindByID(requestID)
}

func (s *friendService) RejectFriendRequest(userID, requestID string) error {
	request, err := s.friendRepo.FindByID(requestID)
	if err != nil {
		return ErrFriendNotFound
	}

	if request.FriendID != userID {
		return ErrFriendNotFound
	}

	return s.friendRepo.Delete(requestID)
}

func (s *friendService) GetFriends(userID string) ([]model.Friend, error) {
	return s.friendRepo.FindByUserID(userID)
}

func (s *friendService) BlockUser(userID, friendID string) error {
	existing, err := s.friendRepo.FindByUserAndFriend(userID, friendID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil {
		return s.friendRepo.UpdateStatus(existing.ID, model.FriendStatusBlocked)
	}

	blocked := &model.Friend{
		UserID:   userID,
		FriendID: friendID,
		Status:   model.FriendStatusBlocked,
	}
	return s.friendRepo.Create(blocked)
}
