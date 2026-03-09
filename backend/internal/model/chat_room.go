package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoomType string

const (
	RoomTypeDirect RoomType = "direct"
	RoomTypeGroup  RoomType = "group"
)

type ChatRoom struct {
	ID        string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	Type      RoomType       `gorm:"type:varchar(10);not null" json:"type"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Members   []ChatRoomMember `gorm:"foreignKey:ChatRoomID" json:"members,omitempty"`
}

func (r *ChatRoom) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

type ChatRoomMember struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ChatRoomID string    `gorm:"type:varchar(36);not null;index" json:"chat_room_id"`
	UserID     string    `gorm:"type:varchar(36);not null;index" json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
	ChatRoom   ChatRoom  `gorm:"foreignKey:ChatRoomID" json:"-"`
	User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (m *ChatRoomMember) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.JoinedAt.IsZero() {
		m.JoinedAt = time.Now()
	}
	return nil
}
