package repository

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupChatRepoDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.ChatRoom{}, &model.ChatRoomMember{}))
	return db
}

func createTestUser(t *testing.T, db *gorm.DB, email, nickname string) *model.User {
	t.Helper()
	user := &model.User{Email: email, Password: "pass", Nickname: nickname}
	require.NoError(t, db.Create(user).Error)
	return user
}

func TestChatRepository_CreateRoom_StoresRoom(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	room := &model.ChatRoom{Name: "Test Room", Type: model.RoomTypeGroup}
	err := repo.CreateRoom(room)

	require.NoError(t, err)
	assert.NotEmpty(t, room.ID)
}

func TestChatRepository_CreateRoom_AssignsUUID(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	room := &model.ChatRoom{Name: "UUID Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))

	assert.Len(t, room.ID, 36)
}

func TestChatRepository_FindRoomByID_ReturnsRoomWhenExists(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	room := &model.ChatRoom{Name: "Find Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))

	found, err := repo.FindRoomByID(room.ID)

	require.NoError(t, err)
	assert.Equal(t, room.ID, found.ID)
	assert.Equal(t, "Find Room", found.Name)
}

func TestChatRepository_FindRoomByID_ReturnsErrorWhenNotFound(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	_, err := repo.FindRoomByID("nonexistent-room-id")

	assert.Error(t, err)
}

func TestChatRepository_FindRoomsByUserID_ReturnsRoomsForUser(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "roomuser@example.com", "roomuser")
	room := &model.ChatRoom{Name: "User Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user.ID}))

	rooms, err := repo.FindRoomsByUserID(user.ID)

	require.NoError(t, err)
	assert.Len(t, rooms, 1)
	assert.Equal(t, room.ID, rooms[0].ID)
}

func TestChatRepository_FindRoomsByUserID_ReturnsEmptyWhenUserHasNoRooms(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "noroomuser@example.com", "noroom")

	rooms, err := repo.FindRoomsByUserID(user.ID)

	require.NoError(t, err)
	assert.Empty(t, rooms)
}

func TestChatRepository_AddMember_AddsMemberToRoom(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "member@example.com", "member")
	room := &model.ChatRoom{Name: "Member Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))

	member := &model.ChatRoomMember{ChatRoomID: room.ID, UserID: user.ID}
	err := repo.AddMember(member)

	require.NoError(t, err)
	assert.NotEmpty(t, member.ID)
}

func TestChatRepository_RemoveMember_RemovesMemberFromRoom(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "removeme@example.com", "removeme")
	room := &model.ChatRoom{Name: "Remove Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user.ID}))

	err := repo.RemoveMember(room.ID, user.ID)
	require.NoError(t, err)

	members, err := repo.FindMembers(room.ID)
	require.NoError(t, err)
	assert.Empty(t, members)
}

func TestChatRepository_FindMembers_ReturnsAllMembersInRoom(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user1 := createTestUser(t, db, "m1@example.com", "m1")
	user2 := createTestUser(t, db, "m2@example.com", "m2")
	room := &model.ChatRoom{Name: "Multi Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user1.ID}))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user2.ID}))

	members, err := repo.FindMembers(room.ID)

	require.NoError(t, err)
	assert.Len(t, members, 2)
}

func TestChatRepository_FindDirectRoom_ReturnsRoomWhenBothUsersAreMembersOfDirectRoom(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user1 := createTestUser(t, db, "direct1@example.com", "d1")
	user2 := createTestUser(t, db, "direct2@example.com", "d2")

	room := &model.ChatRoom{Type: model.RoomTypeDirect}
	require.NoError(t, repo.CreateRoom(room))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user1.ID}))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user2.ID}))

	found, err := repo.FindDirectRoom(user1.ID, user2.ID)

	require.NoError(t, err)
	assert.Equal(t, room.ID, found.ID)
}

func TestChatRepository_FindDirectRoom_ReturnsErrorWhenNoDirectRoomExists(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user1 := createTestUser(t, db, "nd1@example.com", "nd1")
	user2 := createTestUser(t, db, "nd2@example.com", "nd2")

	_, err := repo.FindDirectRoom(user1.ID, user2.ID)

	assert.Error(t, err)
}

func TestIsMember_True(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "ismember@example.com", "ismember")
	room := &model.ChatRoom{Name: "Member Check Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))
	require.NoError(t, repo.AddMember(&model.ChatRoomMember{ChatRoomID: room.ID, UserID: user.ID}))

	result := repo.IsMember(room.ID, user.ID)

	assert.True(t, result)
}

func TestIsMember_False(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "notmember@example.com", "notmember")
	room := &model.ChatRoom{Name: "Not Member Room", Type: model.RoomTypeGroup}
	require.NoError(t, repo.CreateRoom(room))

	result := repo.IsMember(room.ID, user.ID)

	assert.False(t, result)
}

func TestIsMember_InvalidRoom(t *testing.T) {
	db := setupChatRepoDB(t)
	repo := NewChatRepository(db)

	user := createTestUser(t, db, "anyuser@example.com", "anyuser")

	result := repo.IsMember("nonexistent-room-id", user.ID)

	assert.False(t, result)
}
