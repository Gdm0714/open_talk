package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupChatRoomTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}, &ChatRoom{}, &ChatRoomMember{}))
	return db
}

func TestChatRoom_BeforeCreate_GeneratesUUIDWhenIDIsEmpty(t *testing.T) {
	db := setupChatRoomTestDB(t)

	room := &ChatRoom{
		Name: "Test Room",
		Type: RoomTypeGroup,
	}

	err := db.Create(room).Error
	require.NoError(t, err)

	assert.NotEmpty(t, room.ID)
	assert.Len(t, room.ID, 36)
}

func TestChatRoom_BeforeCreate_PreservesExistingID(t *testing.T) {
	db := setupChatRoomTestDB(t)

	existingID := "room-id-0000-1111-2222-333344445555"
	room := &ChatRoom{
		ID:   existingID,
		Name: "Preset Room",
		Type: RoomTypeDirect,
	}

	require.NoError(t, db.Create(room).Error)
	assert.Equal(t, existingID, room.ID)
}

func TestChatRoom_BeforeCreate_GeneratesUniqueIDs(t *testing.T) {
	db := setupChatRoomTestDB(t)

	room1 := &ChatRoom{Name: "Room A", Type: RoomTypeGroup}
	room2 := &ChatRoom{Name: "Room B", Type: RoomTypeGroup}

	require.NoError(t, db.Create(room1).Error)
	require.NoError(t, db.Create(room2).Error)

	assert.NotEqual(t, room1.ID, room2.ID)
}

func TestRoomType_Constants(t *testing.T) {
	assert.Equal(t, RoomType("direct"), RoomTypeDirect)
	assert.Equal(t, RoomType("group"), RoomTypeGroup)
}

func TestChatRoomMember_BeforeCreate_GeneratesUUID(t *testing.T) {
	db := setupChatRoomTestDB(t)

	user := &User{Email: "m@example.com", Password: "pass", Nickname: "mem"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	member := &ChatRoomMember{
		ChatRoomID: room.ID,
		UserID:     user.ID,
	}

	require.NoError(t, db.Create(member).Error)

	assert.NotEmpty(t, member.ID)
	assert.Len(t, member.ID, 36)
}

func TestChatRoomMember_BeforeCreate_SetsJoinedAtWhenZero(t *testing.T) {
	db := setupChatRoomTestDB(t)

	user := &User{Email: "j@example.com", Password: "pass", Nickname: "joiner"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	before := time.Now()
	member := &ChatRoomMember{
		ChatRoomID: room.ID,
		UserID:     user.ID,
	}
	require.NoError(t, db.Create(member).Error)
	after := time.Now()

	assert.False(t, member.JoinedAt.IsZero(), "JoinedAt should be set by BeforeCreate")
	assert.True(t, member.JoinedAt.After(before) || member.JoinedAt.Equal(before))
	assert.True(t, member.JoinedAt.Before(after) || member.JoinedAt.Equal(after))
}

func TestChatRoomMember_BeforeCreate_PreservesExistingJoinedAt(t *testing.T) {
	db := setupChatRoomTestDB(t)

	user := &User{Email: "p@example.com", Password: "pass", Nickname: "preset"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	presetTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	member := &ChatRoomMember{
		ChatRoomID: room.ID,
		UserID:     user.ID,
		JoinedAt:   presetTime,
	}
	require.NoError(t, db.Create(member).Error)

	assert.Equal(t, presetTime.Unix(), member.JoinedAt.Unix(), "BeforeCreate should not overwrite existing JoinedAt")
}
