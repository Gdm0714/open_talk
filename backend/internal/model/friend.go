package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendStatus string

const (
	FriendStatusPending  FriendStatus = "pending"
	FriendStatusAccepted FriendStatus = "accepted"
	FriendStatusBlocked  FriendStatus = "blocked"
)

type Friend struct {
	ID        string       `gorm:"type:varchar(36);primaryKey" json:"id"`
	UserID    string       `gorm:"type:varchar(36);not null;index" json:"user_id"`
	FriendID  string       `gorm:"type:varchar(36);not null;index" json:"friend_id"`
	Status    FriendStatus `gorm:"type:varchar(10);not null;default:pending" json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	User      User         `gorm:"foreignKey:UserID" json:"-"`
	FriendUser User        `gorm:"foreignKey:FriendID" json:"friend,omitempty"`
}

func (f *Friend) BeforeCreate(tx *gorm.DB) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return nil
}
