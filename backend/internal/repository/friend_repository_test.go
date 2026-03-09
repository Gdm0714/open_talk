package repository

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFriendRepoDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.Friend{}))
	return db
}

func createFriendTestUser(t *testing.T, db *gorm.DB, email, nickname string) *model.User {
	t.Helper()
	user := &model.User{Email: email, Password: "pass", Nickname: nickname}
	require.NoError(t, db.Create(user).Error)
	return user
}

func TestFriendRepository_Create_StoresFriend(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fa@example.com", "fa")
	userB := createFriendTestUser(t, db, "fb@example.com", "fb")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusPending}
	err := repo.Create(friend)

	require.NoError(t, err)
	assert.NotEmpty(t, friend.ID)
}

func TestFriendRepository_Create_AssignsUUID(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fc@example.com", "fc")
	userB := createFriendTestUser(t, db, "fd@example.com", "fd")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusPending}
	require.NoError(t, repo.Create(friend))

	assert.Len(t, friend.ID, 36)
}

func TestFriendRepository_FindByUserID_ReturnsOnlyAcceptedFriends(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fe@example.com", "fe")
	userB := createFriendTestUser(t, db, "ff@example.com", "ff")
	userC := createFriendTestUser(t, db, "fg@example.com", "fg")

	accepted := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusAccepted}
	pending := &model.Friend{UserID: userA.ID, FriendID: userC.ID, Status: model.FriendStatusPending}
	require.NoError(t, repo.Create(accepted))
	require.NoError(t, repo.Create(pending))

	friends, err := repo.FindByUserID(userA.ID)

	require.NoError(t, err)
	assert.Len(t, friends, 1)
	assert.Equal(t, userB.ID, friends[0].FriendID)
}

func TestFriendRepository_FindByUserID_ReturnsEmptyWhenNoAcceptedFriends(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fh@example.com", "fh")

	friends, err := repo.FindByUserID(userA.ID)

	require.NoError(t, err)
	assert.Empty(t, friends)
}

func TestFriendRepository_FindByUserAndFriend_ReturnsFriendRecord(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fi@example.com", "fi")
	userB := createFriendTestUser(t, db, "fj@example.com", "fj")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusPending}
	require.NoError(t, repo.Create(friend))

	found, err := repo.FindByUserAndFriend(userA.ID, userB.ID)

	require.NoError(t, err)
	assert.Equal(t, friend.ID, found.ID)
}

func TestFriendRepository_FindByUserAndFriend_ReturnsErrorWhenNotFound(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	_, err := repo.FindByUserAndFriend("no-user", "no-friend")

	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFriendRepository_UpdateStatus_ChangesStatusToAccepted(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fk@example.com", "fk")
	userB := createFriendTestUser(t, db, "fl@example.com", "fl")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusPending}
	require.NoError(t, repo.Create(friend))

	err := repo.UpdateStatus(friend.ID, model.FriendStatusAccepted)
	require.NoError(t, err)

	found, err := repo.FindByID(friend.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FriendStatusAccepted, found.Status)
}

func TestFriendRepository_UpdateStatus_ChangesStatusToBlocked(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fm@example.com", "fm")
	userB := createFriendTestUser(t, db, "fn@example.com", "fn")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusAccepted}
	require.NoError(t, repo.Create(friend))

	err := repo.UpdateStatus(friend.ID, model.FriendStatusBlocked)
	require.NoError(t, err)

	found, err := repo.FindByID(friend.ID)
	require.NoError(t, err)
	assert.Equal(t, model.FriendStatusBlocked, found.Status)
}

func TestFriendRepository_Delete_RemovesFriend(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fo@example.com", "fo")
	userB := createFriendTestUser(t, db, "fp@example.com", "fp")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusAccepted}
	require.NoError(t, repo.Create(friend))

	err := repo.Delete(friend.ID)
	require.NoError(t, err)

	_, err = repo.FindByID(friend.ID)
	assert.Error(t, err)
}

func TestFriendRepository_FindByID_ReturnsFriendWhenExists(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	userA := createFriendTestUser(t, db, "fq@example.com", "fq")
	userB := createFriendTestUser(t, db, "fr@example.com", "fr")

	friend := &model.Friend{UserID: userA.ID, FriendID: userB.ID, Status: model.FriendStatusPending}
	require.NoError(t, repo.Create(friend))

	found, err := repo.FindByID(friend.ID)

	require.NoError(t, err)
	assert.Equal(t, friend.ID, found.ID)
	assert.Equal(t, userA.ID, found.UserID)
}

func TestFriendRepository_FindByID_ReturnsErrorWhenNotFound(t *testing.T) {
	db := setupFriendRepoDB(t)
	repo := NewFriendRepository(db)

	_, err := repo.FindByID("nonexistent-friend-id")

	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
