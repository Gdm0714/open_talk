package repository

import (
	"github.com/godongmin/open_talk/backend/internal/model"
	"gorm.io/gorm"
)

type ChatRepository interface {
	CreateRoom(room *model.ChatRoom) error
	FindRoomByID(id string) (*model.ChatRoom, error)
	FindRoomsByUserID(userID string) ([]model.ChatRoom, error)
	AddMember(member *model.ChatRoomMember) error
	RemoveMember(roomID, userID string) error
	FindMembers(roomID string) ([]model.ChatRoomMember, error)
	FindDirectRoom(userID1, userID2 string) (*model.ChatRoom, error)
	IsMember(roomID, userID string) bool
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateRoom(room *model.ChatRoom) error {
	return r.db.Create(room).Error
}

func (r *chatRepository) FindRoomByID(id string) (*model.ChatRoom, error) {
	var room model.ChatRoom
	if err := r.db.Preload("Members").Preload("Members.User").Where("id = ?", id).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *chatRepository) FindRoomsByUserID(userID string) ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	if err := r.db.
		Joins("JOIN chat_room_members ON chat_room_members.chat_room_id = chat_rooms.id").
		Where("chat_room_members.user_id = ?", userID).
		Preload("Members").
		Preload("Members.User").
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *chatRepository) AddMember(member *model.ChatRoomMember) error {
	return r.db.Create(member).Error
}

func (r *chatRepository) RemoveMember(roomID, userID string) error {
	return r.db.Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Delete(&model.ChatRoomMember{}).Error
}

func (r *chatRepository) FindMembers(roomID string) ([]model.ChatRoomMember, error) {
	var members []model.ChatRoomMember
	if err := r.db.Preload("User").Where("chat_room_id = ?", roomID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *chatRepository) IsMember(roomID, userID string) bool {
	var count int64
	r.db.Model(&model.ChatRoomMember{}).
		Where("chat_room_id = ? AND user_id = ?", roomID, userID).
		Count(&count)
	return count > 0
}

func (r *chatRepository) FindDirectRoom(userID1, userID2 string) (*model.ChatRoom, error) {
	var room model.ChatRoom
	if err := r.db.
		Joins("JOIN chat_room_members m1 ON m1.chat_room_id = chat_rooms.id AND m1.user_id = ?", userID1).
		Joins("JOIN chat_room_members m2 ON m2.chat_room_id = chat_rooms.id AND m2.user_id = ?", userID2).
		Where("chat_rooms.type = ?", model.RoomTypeDirect).
		First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}
