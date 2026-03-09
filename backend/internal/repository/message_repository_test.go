package repository

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMessageRepoDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&model.User{}, &model.ChatRoom{}, &model.ChatRoomMember{}, &model.Message{}))
	return db
}

func seedRoomAndUser(t *testing.T, db *gorm.DB) (*model.User, *model.ChatRoom) {
	t.Helper()
	user := &model.User{Email: "msg@example.com", Password: "pass", Nickname: "sender"}
	require.NoError(t, db.Create(user).Error)
	room := &model.ChatRoom{Name: "Room", Type: model.RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)
	return user, room
}

func TestMessageRepository_Create_StoresMessage(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	msg := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "Hello",
		MessageType: model.MessageTypeText,
	}
	err := repo.Create(msg)

	require.NoError(t, err)
	assert.NotEmpty(t, msg.ID)
}

func TestMessageRepository_Create_AssignsUUID(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	msg := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "UUID test",
		MessageType: model.MessageTypeText,
	}
	require.NoError(t, repo.Create(msg))

	assert.Len(t, msg.ID, 36)
}

func TestMessageRepository_FindByID_ReturnsMessageWhenExists(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	msg := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "Find me",
		MessageType: model.MessageTypeText,
	}
	require.NoError(t, repo.Create(msg))

	found, err := repo.FindByID(msg.ID)

	require.NoError(t, err)
	assert.Equal(t, msg.ID, found.ID)
	assert.Equal(t, "Find me", found.Content)
}

func TestMessageRepository_FindByID_ReturnsErrorWhenNotFound(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)

	_, err := repo.FindByID("nonexistent-message-id")

	assert.Error(t, err)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestMessageRepository_FindByRoomID_ReturnsMessagesForRoom(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	for i := 0; i < 3; i++ {
		msg := &model.Message{
			ChatRoomID:  room.ID,
			SenderID:    user.ID,
			Content:     "msg",
			MessageType: model.MessageTypeText,
		}
		require.NoError(t, repo.Create(msg))
	}

	messages, err := repo.FindByRoomID(room.ID, 10, 0)

	require.NoError(t, err)
	assert.Len(t, messages, 3)
}

func TestMessageRepository_FindByRoomID_RespectsLimit(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	for i := 0; i < 5; i++ {
		msg := &model.Message{
			ChatRoomID:  room.ID,
			SenderID:    user.ID,
			Content:     "msg",
			MessageType: model.MessageTypeText,
		}
		require.NoError(t, repo.Create(msg))
	}

	messages, err := repo.FindByRoomID(room.ID, 2, 0)

	require.NoError(t, err)
	assert.Len(t, messages, 2)
}

func TestMessageRepository_FindByRoomID_RespectsOffset(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	for i := 0; i < 4; i++ {
		msg := &model.Message{
			ChatRoomID:  room.ID,
			SenderID:    user.ID,
			Content:     "msg",
			MessageType: model.MessageTypeText,
		}
		require.NoError(t, repo.Create(msg))
	}

	messages, err := repo.FindByRoomID(room.ID, 10, 2)

	require.NoError(t, err)
	assert.Len(t, messages, 2)
}

func TestMessageRepository_FindByRoomID_ReturnsEmptyForUnknownRoom(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)

	messages, err := repo.FindByRoomID("unknown-room-id", 10, 0)

	require.NoError(t, err)
	assert.Empty(t, messages)
}

func TestMarkAsRead_Success(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	sender, room := seedRoomAndUser(t, db)

	reader := &model.User{Email: "reader@example.com", Password: "pass", Nickname: "reader"}
	require.NoError(t, db.Create(reader).Error)

	msg := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    sender.ID,
		Content:     "Hello reader",
		MessageType: model.MessageTypeText,
	}
	require.NoError(t, repo.Create(msg))

	err := repo.MarkAsRead(room.ID, reader.ID)
	require.NoError(t, err)

	var updated model.Message
	require.NoError(t, db.Where("id = ?", msg.ID).First(&updated).Error)
	assert.True(t, updated.IsRead)
	assert.Contains(t, updated.ReadBy, reader.ID)
}

func TestGetUnreadCount_Success(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	sender, room := seedRoomAndUser(t, db)

	reader := &model.User{Email: "unread@example.com", Password: "pass", Nickname: "unread"}
	require.NoError(t, db.Create(reader).Error)

	for i := 0; i < 3; i++ {
		msg := &model.Message{
			ChatRoomID:  room.ID,
			SenderID:    sender.ID,
			Content:     "msg",
			MessageType: model.MessageTypeText,
		}
		require.NoError(t, repo.Create(msg))
	}

	count, err := repo.GetUnreadCount(room.ID, reader.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestGetUnreadCount_AllRead(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	sender, room := seedRoomAndUser(t, db)

	reader := &model.User{Email: "allread@example.com", Password: "pass", Nickname: "allread"}
	require.NoError(t, db.Create(reader).Error)

	msg := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    sender.ID,
		Content:     "already read",
		MessageType: model.MessageTypeText,
		IsRead:      true,
		ReadBy:      reader.ID,
	}
	require.NoError(t, repo.Create(msg))

	count, err := repo.GetUnreadCount(room.ID, reader.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestGetLastMessage_Success(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	user, room := seedRoomAndUser(t, db)

	first := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "first",
		MessageType: model.MessageTypeText,
	}
	require.NoError(t, repo.Create(first))

	last := &model.Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "last message",
		MessageType: model.MessageTypeText,
	}
	require.NoError(t, repo.Create(last))

	found, err := repo.GetLastMessage(room.ID)
	require.NoError(t, err)
	assert.Equal(t, last.ID, found.ID)
	assert.Equal(t, "last message", found.Content)
}

func TestGetLastMessage_EmptyRoom(t *testing.T) {
	db := setupMessageRepoDB(t)
	repo := NewMessageRepository(db)
	_, room := seedRoomAndUser(t, db)

	found, err := repo.GetLastMessage(room.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}
