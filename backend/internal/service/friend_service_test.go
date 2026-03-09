package service

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// fakeFriendRepoFull is a richer fake that supports FindByUserAndFriend properly
// (the one in user_service_test.go always returns ErrRecordNotFound).

type fakeFriendRepoFull struct {
	friends map[string]*model.Friend // keyed by ID
	nextID  int
}

func newFakeFriendRepoFull() *fakeFriendRepoFull {
	return &fakeFriendRepoFull{
		friends: make(map[string]*model.Friend),
	}
}

func (r *fakeFriendRepoFull) nextIDStr() string {
	r.nextID++
	return "fr-" + string(rune('0'+r.nextID))
}

func (r *fakeFriendRepoFull) Create(f *model.Friend) error {
	if f.ID == "" {
		f.ID = r.nextIDStr()
	}
	r.friends[f.ID] = f
	return nil
}

func (r *fakeFriendRepoFull) FindByUserID(userID string) ([]model.Friend, error) {
	var out []model.Friend
	for _, f := range r.friends {
		if f.UserID == userID {
			out = append(out, *f)
		}
	}
	return out, nil
}

func (r *fakeFriendRepoFull) FindByUserAndFriend(userID, friendID string) (*model.Friend, error) {
	for _, f := range r.friends {
		if f.UserID == userID && f.FriendID == friendID {
			return f, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeFriendRepoFull) FindByID(id string) (*model.Friend, error) {
	f, ok := r.friends[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return f, nil
}

func (r *fakeFriendRepoFull) UpdateStatus(id string, status model.FriendStatus) error {
	f, ok := r.friends[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	f.Status = status
	return nil
}

func (r *fakeFriendRepoFull) Delete(id string) error {
	if _, ok := r.friends[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(r.friends, id)
	return nil
}

// --- helpers ---

func newFriendSvc(repo *fakeFriendRepoFull) FriendService {
	return NewFriendService(repo)
}

// --- SendFriendRequest tests ---

func TestSendRequest_Success(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	friend, err := svc.SendFriendRequest("user-1", "user-2")

	require.NoError(t, err)
	assert.NotNil(t, friend)
	assert.Equal(t, "user-1", friend.UserID)
	assert.Equal(t, "user-2", friend.FriendID)
	assert.Equal(t, model.FriendStatusPending, friend.Status)
}

func TestSendRequest_Self(t *testing.T) {
	svc := newFriendSvc(newFakeFriendRepoFull())

	_, err := svc.SendFriendRequest("user-1", "user-1")

	assert.ErrorIs(t, err, ErrCannotFriendSelf)
}

func TestSendRequest_AlreadyExists(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	_, err := svc.SendFriendRequest("user-1", "user-2")
	require.NoError(t, err)

	// sending again should return ErrFriendRequestExists
	_, err = svc.SendFriendRequest("user-1", "user-2")
	assert.ErrorIs(t, err, ErrFriendRequestExists)
}

// --- AcceptFriendRequest tests ---

func TestAcceptRequest_Success(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	// user-1 sends request to user-2; user-2 accepts
	pending, err := svc.SendFriendRequest("user-1", "user-2")
	require.NoError(t, err)

	accepted, err := svc.AcceptFriendRequest("user-2", pending.ID)
	require.NoError(t, err)
	assert.NotNil(t, accepted)
	assert.Equal(t, model.FriendStatusAccepted, accepted.Status)
}

func TestAcceptRequest_NotFound(t *testing.T) {
	svc := newFriendSvc(newFakeFriendRepoFull())

	_, err := svc.AcceptFriendRequest("user-2", "nonexistent-id")
	assert.ErrorIs(t, err, ErrFriendNotFound)
}

func TestAcceptRequest_WrongUser(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	pending, err := svc.SendFriendRequest("user-1", "user-2")
	require.NoError(t, err)

	// user-3 tries to accept a request meant for user-2
	_, err = svc.AcceptFriendRequest("user-3", pending.ID)
	assert.ErrorIs(t, err, ErrFriendNotFound)
}

// --- RejectFriendRequest tests ---

func TestRejectRequest_Success(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	pending, err := svc.SendFriendRequest("user-1", "user-2")
	require.NoError(t, err)

	err = svc.RejectFriendRequest("user-2", pending.ID)
	require.NoError(t, err)

	// should be deleted
	_, err = repo.FindByID(pending.ID)
	assert.Error(t, err)
}

func TestRejectRequest_NotFound(t *testing.T) {
	svc := newFriendSvc(newFakeFriendRepoFull())

	err := svc.RejectFriendRequest("user-2", "nonexistent")
	assert.ErrorIs(t, err, ErrFriendNotFound)
}

// --- GetFriends tests ---

func TestGetFriends_Success(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	pending, err := svc.SendFriendRequest("user-1", "user-2")
	require.NoError(t, err)
	_, err = svc.AcceptFriendRequest("user-2", pending.ID)
	require.NoError(t, err)

	friends, err := svc.GetFriends("user-1")
	require.NoError(t, err)
	// user-1 has one accepted friendship
	assert.NotEmpty(t, friends)
}

func TestGetFriends_Empty(t *testing.T) {
	svc := newFriendSvc(newFakeFriendRepoFull())

	friends, err := svc.GetFriends("user-nobody")
	require.NoError(t, err)
	assert.Empty(t, friends)
}

// --- BlockUser tests ---

func TestBlockFriend_Success_NoExisting(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	err := svc.BlockUser("user-1", "user-2")
	require.NoError(t, err)

	blocked, err := repo.FindByUserAndFriend("user-1", "user-2")
	require.NoError(t, err)
	assert.Equal(t, model.FriendStatusBlocked, blocked.Status)
}

func TestBlockFriend_Success_ExistingRelationship(t *testing.T) {
	repo := newFakeFriendRepoFull()
	svc := newFriendSvc(repo)

	// Create pending request first
	pending, err := svc.SendFriendRequest("user-1", "user-2")
	require.NoError(t, err)

	// Now block - should update existing record
	err = svc.BlockUser("user-1", "user-2")
	require.NoError(t, err)

	f, err := repo.FindByID(pending.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FriendStatusBlocked, f.Status)
}
