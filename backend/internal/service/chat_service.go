package service

import (
	"errors"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrRoomNotFound    = errors.New("chat room not found")
	ErrNotRoomMember   = errors.New("not a member of this chat room")
	ErrDirectRoomExists = errors.New("direct chat room already exists")
)

type ChatService interface {
	CreateDirectChat(userID, targetUserID string) (*model.ChatRoom, error)
	CreateGroupChat(userID string, name string, memberIDs []string) (*model.ChatRoom, error)
	GetUserChats(userID string) ([]model.ChatRoom, error)
	GetChatMessages(userID, roomID string, limit, offset int) ([]model.Message, error)
	SendMessage(userID, roomID, content string, messageType model.MessageType) (*model.Message, error)
}

type chatService struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
}

func NewChatService(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository) ChatService {
	return &chatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

func (s *chatService) CreateDirectChat(userID, targetUserID string) (*model.ChatRoom, error) {
	existing, err := s.chatRepo.FindDirectRoom(userID, targetUserID)
	if err == nil && existing != nil {
		return existing, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	room := &model.ChatRoom{
		Type: model.RoomTypeDirect,
	}
	if err := s.chatRepo.CreateRoom(room); err != nil {
		return nil, err
	}

	for _, uid := range []string{userID, targetUserID} {
		member := &model.ChatRoomMember{
			ChatRoomID: room.ID,
			UserID:     uid,
		}
		if err := s.chatRepo.AddMember(member); err != nil {
			return nil, err
		}
	}

	return s.chatRepo.FindRoomByID(room.ID)
}

func (s *chatService) CreateGroupChat(userID string, name string, memberIDs []string) (*model.ChatRoom, error) {
	room := &model.ChatRoom{
		Name: name,
		Type: model.RoomTypeGroup,
	}
	if err := s.chatRepo.CreateRoom(room); err != nil {
		return nil, err
	}

	allMembers := append([]string{userID}, memberIDs...)
	for _, uid := range allMembers {
		member := &model.ChatRoomMember{
			ChatRoomID: room.ID,
			UserID:     uid,
		}
		if err := s.chatRepo.AddMember(member); err != nil {
			return nil, err
		}
	}

	return s.chatRepo.FindRoomByID(room.ID)
}

func (s *chatService) GetUserChats(userID string) ([]model.ChatRoom, error) {
	return s.chatRepo.FindRoomsByUserID(userID)
}

func (s *chatService) GetChatMessages(userID, roomID string, limit, offset int) ([]model.Message, error) {
	if err := s.verifyMembership(userID, roomID); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 50
	}

	return s.messageRepo.FindByRoomID(roomID, limit, offset)
}

func (s *chatService) SendMessage(userID, roomID, content string, messageType model.MessageType) (*model.Message, error) {
	if err := s.verifyMembership(userID, roomID); err != nil {
		return nil, err
	}

	msg := &model.Message{
		ChatRoomID:  roomID,
		SenderID:    userID,
		Content:     content,
		MessageType: messageType,
	}

	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	return s.messageRepo.FindByID(msg.ID)
}

func (s *chatService) verifyMembership(userID, roomID string) error {
	members, err := s.chatRepo.FindMembers(roomID)
	if err != nil {
		return ErrRoomNotFound
	}

	for _, m := range members {
		if m.UserID == userID {
			return nil
		}
	}

	return ErrNotRoomMember
}
