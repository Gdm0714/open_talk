package repository

import (
	"github.com/godongmin/open_talk/backend/internal/model"
	"gorm.io/gorm"
)

type FriendRepository interface {
	Create(friend *model.Friend) error
	FindByUserID(userID string) ([]model.Friend, error)
	FindByUserAndFriend(userID, friendID string) (*model.Friend, error)
	UpdateStatus(id string, status model.FriendStatus) error
	Delete(id string) error
	FindByID(id string) (*model.Friend, error)
}

type friendRepository struct {
	db *gorm.DB
}

func NewFriendRepository(db *gorm.DB) FriendRepository {
	return &friendRepository{db: db}
}

func (r *friendRepository) Create(friend *model.Friend) error {
	return r.db.Create(friend).Error
}

func (r *friendRepository) FindByUserID(userID string) ([]model.Friend, error) {
	var friends []model.Friend
	if err := r.db.Preload("FriendUser").
		Where("user_id = ? AND status = ?", userID, model.FriendStatusAccepted).
		Find(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}

func (r *friendRepository) FindByUserAndFriend(userID, friendID string) (*model.Friend, error) {
	var friend model.Friend
	if err := r.db.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend).Error; err != nil {
		return nil, err
	}
	return &friend, nil
}

func (r *friendRepository) UpdateStatus(id string, status model.FriendStatus) error {
	return r.db.Model(&model.Friend{}).Where("id = ?", id).Update("status", status).Error
}

func (r *friendRepository) Delete(id string) error {
	return r.db.Delete(&model.Friend{}, "id = ?", id).Error
}

func (r *friendRepository) FindByID(id string) (*model.Friend, error) {
	var friend model.Friend
	if err := r.db.Preload("FriendUser").Where("id = ?", id).First(&friend).Error; err != nil {
		return nil, err
	}
	return &friend, nil
}
