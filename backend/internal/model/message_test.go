package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMessageTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&User{}, &ChatRoom{}, &ChatRoomMember{}, &Message{}))
	return db
}

func TestMessageType_Constants(t *testing.T) {
	assert.Equal(t, MessageType("text"), MessageTypeText)
	assert.Equal(t, MessageType("image"), MessageTypeImage)
	assert.Equal(t, MessageType("file"), MessageTypeFile)
}

func TestMessage_BeforeCreate_GeneratesUUIDWhenIDIsEmpty(t *testing.T) {
	db := setupMessageTestDB(t)

	user := &User{Email: "msg@example.com", Password: "pass", Nickname: "sender"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	msg := &Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "Hello",
		MessageType: MessageTypeText,
	}

	require.NoError(t, db.Create(msg).Error)

	assert.NotEmpty(t, msg.ID)
	assert.Len(t, msg.ID, 36)
}

func TestMessage_BeforeCreate_PreservesExistingID(t *testing.T) {
	db := setupMessageTestDB(t)

	user := &User{Email: "msg2@example.com", Password: "pass", Nickname: "s2"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room2", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	existingID := "msg-id-0000-1111-2222-333344445555"
	msg := &Message{
		ID:          existingID,
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "Hi",
		MessageType: MessageTypeText,
	}

	require.NoError(t, db.Create(msg).Error)
	assert.Equal(t, existingID, msg.ID)
}

func TestMessage_BeforeCreate_GeneratesUniqueIDs(t *testing.T) {
	db := setupMessageTestDB(t)

	user := &User{Email: "msg3@example.com", Password: "pass", Nickname: "s3"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room3", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	msg1 := &Message{ChatRoomID: room.ID, SenderID: user.ID, Content: "A", MessageType: MessageTypeText}
	msg2 := &Message{ChatRoomID: room.ID, SenderID: user.ID, Content: "B", MessageType: MessageTypeText}

	require.NoError(t, db.Create(msg1).Error)
	require.NoError(t, db.Create(msg2).Error)

	assert.NotEqual(t, msg1.ID, msg2.ID)
}

func TestMessage_AllFieldsStoredAndRetrieved(t *testing.T) {
	db := setupMessageTestDB(t)

	user := &User{Email: "msg4@example.com", Password: "pass", Nickname: "s4"}
	require.NoError(t, db.Create(user).Error)

	room := &ChatRoom{Name: "Room4", Type: RoomTypeGroup}
	require.NoError(t, db.Create(room).Error)

	msg := &Message{
		ChatRoomID:  room.ID,
		SenderID:    user.ID,
		Content:     "Image content",
		MessageType: MessageTypeImage,
	}
	require.NoError(t, db.Create(msg).Error)

	var found Message
	require.NoError(t, db.First(&found, "id = ?", msg.ID).Error)

	assert.Equal(t, room.ID, found.ChatRoomID)
	assert.Equal(t, user.ID, found.SenderID)
	assert.Equal(t, "Image content", found.Content)
	assert.Equal(t, MessageTypeImage, found.MessageType)
	assert.NotZero(t, found.CreatedAt)
}
