package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeImage MessageType = "image"
	MessageTypeFile  MessageType = "file"
)

type Message struct {
	ID          string      `gorm:"type:varchar(36);primaryKey" json:"id"`
	ChatRoomID  string      `gorm:"type:varchar(36);not null;index" json:"chat_room_id"`
	SenderID    string      `gorm:"type:varchar(36);not null;index" json:"sender_id"`
	Content     string      `gorm:"type:text;not null" json:"content"`
	MessageType MessageType `gorm:"type:varchar(10);not null;default:text" json:"message_type"`
	ReadBy      string      `gorm:"type:text" json:"read_by"`
	IsRead      bool        `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time   `json:"created_at"`
	ChatRoom    ChatRoom    `gorm:"foreignKey:ChatRoomID" json:"-"`
	Sender      User        `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}
