package service

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- fake friend repo ---

type fakeFriendRepo struct {
	friends map[string]*model.Friend
}

func newFakeFriendRepo() *fakeFriendRepo {
	return &fakeFriendRepo{friends: make(map[string]*model.Friend)}
}

func (r *fakeFriendRepo) Create(f *model.Friend) error {
	r.friends[f.ID] = f
	return nil
}
func (r *fakeFriendRepo) FindByUserID(userID string) ([]model.Friend, error) {
	var out []model.Friend
	for _, f := range r.friends {
		if f.UserID == userID {
			out = append(out, *f)
		}
	}
	return out, nil
}
func (r *fakeFriendRepo) FindByUserAndFriend(userID, friendID string) (*model.Friend, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeFriendRepo) UpdateStatus(id string, status model.FriendStatus) error { return nil }
func (r *fakeFriendRepo) Delete(id string) error {
	delete(r.friends, id)
	return nil
}
func (r *fakeFriendRepo) FindByID(id string) (*model.Friend, error) {
	f, ok := r.friends[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return f, nil
}

// ensure interface is satisfied
var _ repository.FriendRepository = (*fakeFriendRepo)(nil)

// --- fake chat repo ---

type fakeChatRepo struct{}

func (r *fakeChatRepo) CreateRoom(room *model.ChatRoom) error { return nil }
func (r *fakeChatRepo) FindRoomByID(id string) (*model.ChatRoom, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeChatRepo) FindRoomsByUserID(userID string) ([]model.ChatRoom, error) {
	return nil, nil
}
func (r *fakeChatRepo) AddMember(member *model.ChatRoomMember) error    { return nil }
func (r *fakeChatRepo) RemoveMember(roomID, userID string) error        { return nil }
func (r *fakeChatRepo) FindMembers(roomID string) ([]model.ChatRoomMember, error) {
	return nil, nil
}
func (r *fakeChatRepo) FindDirectRoom(userID1, userID2 string) (*model.ChatRoom, error) {
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeChatRepo) IsMember(roomID, userID string) bool { return false }

var _ repository.ChatRepository = (*fakeChatRepo)(nil)

func newUserSvc(repo *fakeUserRepo) UserService {
	return NewUserService(repo, newFakeFriendRepo(), &fakeChatRepo{})
}

// --- DeleteAccount tests ---

func TestDeleteAccount_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newUserSvc(repo)

	authSvc := newAuthSvc(repo)
	_, _, err := authSvc.Register("del@example.com", "password123", "delnick")
	require.NoError(t, err)

	user, err := repo.FindByEmail("del@example.com")
	require.NoError(t, err)

	err = svc.DeleteAccount(user.ID, "password123")
	assert.NoError(t, err)

	_, err = repo.FindByID(user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteAccount_WrongPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newUserSvc(repo)

	authSvc := newAuthSvc(repo)
	_, _, err := authSvc.Register("del2@example.com", "password123", "delnick2")
	require.NoError(t, err)

	user, err := repo.FindByEmail("del2@example.com")
	require.NoError(t, err)

	err = svc.DeleteAccount(user.ID, "wrongpassword")
	assert.ErrorIs(t, err, ErrWrongPassword)
}

func TestDeleteAccount_UserNotFound(t *testing.T) {
	repo := newFakeUserRepo()
	svc := newUserSvc(repo)

	err := svc.DeleteAccount("nonexistent-id", "password123")
	assert.Error(t, err)
}
