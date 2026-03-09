package service

import (
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetProfile(userID string) (*model.User, error)
	UpdateProfile(userID string, nickname, avatarURL, statusMessage *string) (*model.User, error)
	SearchUsers(query string) ([]model.User, error)
	DeleteAccount(userID, password string) error
}

type userService struct {
	userRepo   repository.UserRepository
	friendRepo repository.FriendRepository
	chatRepo   repository.ChatRepository
}

func NewUserService(userRepo repository.UserRepository, friendRepo repository.FriendRepository, chatRepo repository.ChatRepository) UserService {
	return &userService{userRepo: userRepo, friendRepo: friendRepo, chatRepo: chatRepo}
}

func (s *userService) GetProfile(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *userService) UpdateProfile(userID string, nickname, avatarURL, statusMessage *string) (*model.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if nickname != nil {
		user.Nickname = *nickname
	}
	if avatarURL != nil {
		user.AvatarURL = *avatarURL
	}
	if statusMessage != nil {
		user.StatusMessage = *statusMessage
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) SearchUsers(query string) ([]model.User, error) {
	return s.userRepo.Search(query)
}

func (s *userService) DeleteAccount(userID, password string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrWrongPassword
	}

	// Clean up friend relationships
	friends, err := s.friendRepo.FindByUserID(userID)
	if err == nil {
		for _, f := range friends {
			_ = s.friendRepo.Delete(f.ID)
		}
	}

	// Remove from chat room memberships
	rooms, err := s.chatRepo.FindRoomsByUserID(userID)
	if err == nil {
		for _, room := range rooms {
			_ = s.chatRepo.RemoveMember(room.ID, userID)
		}
	}

	return s.userRepo.Delete(userID)
}
