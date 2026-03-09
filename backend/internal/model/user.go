package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:varchar(36);primaryKey" json:"id"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"type:varchar(255);not null" json:"-"`
	Nickname      string         `gorm:"type:varchar(20);not null" json:"nickname"`
	AvatarURL     string         `gorm:"type:varchar(500)" json:"avatar_url"`
	StatusMessage string         `gorm:"type:varchar(255)" json:"status_message"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
