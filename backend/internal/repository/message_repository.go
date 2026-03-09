package repository

import (
	"github.com/godongmin/open_talk/backend/internal/model"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(message *model.Message) error
	FindByRoomID(roomID string, limit, offset int) ([]model.Message, error)
	FindByID(id string) (*model.Message, error)
	MarkAsRead(roomID, userID string) error
	GetUnreadCount(roomID, userID string) (int64, error)
	GetLastMessage(roomID string) (*model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) FindByRoomID(roomID string, limit, offset int) ([]model.Message, error) {
	var messages []model.Message
	if err := r.db.Preload("Sender").
		Where("chat_room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messageRepository) FindByID(id string) (*model.Message, error) {
	var message model.Message
	if err := r.db.Preload("Sender").Where("id = ?", id).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) MarkAsRead(roomID, userID string) error {
	return r.db.Model(&model.Message{}).
		Where("chat_room_id = ? AND sender_id != ? AND (read_by NOT LIKE ? OR read_by = '' OR read_by IS NULL)",
			roomID, userID, "%"+userID+"%").
		Updates(map[string]interface{}{
			"is_read": true,
			"read_by": gorm.Expr("CASE WHEN read_by = '' OR read_by IS NULL THEN ? ELSE read_by || ',' || ? END", userID, userID),
		}).Error
}

func (r *messageRepository) GetUnreadCount(roomID, userID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Message{}).
		Where("chat_room_id = ? AND sender_id != ? AND (read_by NOT LIKE ? OR read_by = '' OR read_by IS NULL)",
			roomID, userID, "%"+userID+"%").
		Count(&count).Error
	return count, err
}

func (r *messageRepository) GetLastMessage(roomID string) (*model.Message, error) {
	var message model.Message
	if err := r.db.Where("chat_room_id = ?", roomID).
		Order("created_at DESC").
		First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

