package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFriendTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}, &Friend{}))
	return db
}

func TestFriendStatus_Constants(t *testing.T) {
	assert.Equal(t, FriendStatus("pending"), FriendStatusPending)
	assert.Equal(t, FriendStatus("accepted"), FriendStatusAccepted)
	assert.Equal(t, FriendStatus("blocked"), FriendStatusBlocked)
}

func TestFriend_BeforeCreate_GeneratesUUIDWhenIDIsEmpty(t *testing.T) {
	db := setupFriendTestDB(t)

	userA := &User{Email: "a@example.com", Password: "pass", Nickname: "aa"}
	userB := &User{Email: "b@example.com", Password: "pass", Nickname: "bb"}
	require.NoError(t, db.Create(userA).Error)
	require.NoError(t, db.Create(userB).Error)

	friend := &Friend{
		UserID:   userA.ID,
		FriendID: userB.ID,
		Status:   FriendStatusPending,
	}

	require.NoError(t, db.Create(friend).Error)

	assert.NotEmpty(t, friend.ID)
	assert.Len(t, friend.ID, 36)
}

func TestFriend_BeforeCreate_PreservesExistingID(t *testing.T) {
	db := setupFriendTestDB(t)

	userA := &User{Email: "c@example.com", Password: "pass", Nickname: "cc"}
	userB := &User{Email: "d@example.com", Password: "pass", Nickname: "dd"}
	require.NoError(t, db.Create(userA).Error)
	require.NoError(t, db.Create(userB).Error)

	existingID := "fr-id-00000-1111-2222-333344445555"
	friend := &Friend{
		ID:       existingID,
		UserID:   userA.ID,
		FriendID: userB.ID,
		Status:   FriendStatusPending,
	}

	require.NoError(t, db.Create(friend).Error)
	assert.Equal(t, existingID, friend.ID)
}

func TestFriend_BeforeCreate_GeneratesUniqueIDs(t *testing.T) {
	db := setupFriendTestDB(t)

	userA := &User{Email: "e@example.com", Password: "pass", Nickname: "ee"}
	userB := &User{Email: "f@example.com", Password: "pass", Nickname: "ff"}
	userC := &User{Email: "g@example.com", Password: "pass", Nickname: "gg"}
	require.NoError(t, db.Create(userA).Error)
	require.NoError(t, db.Create(userB).Error)
	require.NoError(t, db.Create(userC).Error)

	f1 := &Friend{UserID: userA.ID, FriendID: userB.ID, Status: FriendStatusPending}
	f2 := &Friend{UserID: userA.ID, FriendID: userC.ID, Status: FriendStatusPending}

	require.NoError(t, db.Create(f1).Error)
	require.NoError(t, db.Create(f2).Error)

	assert.NotEqual(t, f1.ID, f2.ID)
}

func TestFriend_AllFieldsStoredAndRetrieved(t *testing.T) {
	db := setupFriendTestDB(t)

	userA := &User{Email: "h@example.com", Password: "pass", Nickname: "hh"}
	userB := &User{Email: "i@example.com", Password: "pass", Nickname: "ii"}
	require.NoError(t, db.Create(userA).Error)
	require.NoError(t, db.Create(userB).Error)

	friend := &Friend{
		UserID:   userA.ID,
		FriendID: userB.ID,
		Status:   FriendStatusAccepted,
	}
	require.NoError(t, db.Create(friend).Error)

	var found Friend
	require.NoError(t, db.First(&found, "id = ?", friend.ID).Error)

	assert.Equal(t, userA.ID, found.UserID)
	assert.Equal(t, userB.ID, found.FriendID)
	assert.Equal(t, FriendStatusAccepted, found.Status)
	assert.NotZero(t, found.CreatedAt)
}
