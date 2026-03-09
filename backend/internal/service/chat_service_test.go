package service

import (
	"testing"

	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// --- fake chat repository ---

type fakeChatRepoSvc struct {
	rooms   map[string]*model.ChatRoom
	members map[string][]model.ChatRoomMember // roomID -> members
}

func newFakeChatRepoSvc() *fakeChatRepoSvc {
	return &fakeChatRepoSvc{
		rooms:   make(map[string]*model.ChatRoom),
		members: make(map[string][]model.ChatRoomMember),
	}
}

func (r *fakeChatRepoSvc) CreateRoom(room *model.ChatRoom) error {
	if room.ID == "" {
		roomID := "room-" + string(room.Type) + "-" + room.Name
		if roomID == "room--" {
			roomID = "room-direct-new"
		}
		room.ID = roomID
	}
	r.rooms[room.ID] = room
	return nil
}

func (r *fakeChatRepoSvc) FindRoomByID(id string) (*model.ChatRoom, error) {
	room, ok := r.rooms[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	room.Members = r.members[id]
	return room, nil
}

func (r *fakeChatRepoSvc) FindRoomsByUserID(userID string) ([]model.ChatRoom, error) {
	var result []model.ChatRoom
	for _, room := range r.rooms {
		for _, m := range r.members[room.ID] {
			if m.UserID == userID {
				result = append(result, *room)
				break
			}
		}
	}
	return result, nil
}

func (r *fakeChatRepoSvc) AddMember(member *model.ChatRoomMember) error {
	r.members[member.ChatRoomID] = append(r.members[member.ChatRoomID], *member)
	return nil
}

func (r *fakeChatRepoSvc) RemoveMember(roomID, userID string) error {
	members := r.members[roomID]
	filtered := make([]model.ChatRoomMember, 0, len(members))
	for _, m := range members {
		if m.UserID != userID {
			filtered = append(filtered, m)
		}
	}
	r.members[roomID] = filtered
	return nil
}

func (r *fakeChatRepoSvc) FindMembers(roomID string) ([]model.ChatRoomMember, error) {
	members, ok := r.members[roomID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return members, nil
}

func (r *fakeChatRepoSvc) FindDirectRoom(userID1, userID2 string) (*model.ChatRoom, error) {
	for _, room := range r.rooms {
		if room.Type != model.RoomTypeDirect {
			continue
		}
		members := r.members[room.ID]
		hasUser1, hasUser2 := false, false
		for _, m := range members {
			if m.UserID == userID1 {
				hasUser1 = true
			}
			if m.UserID == userID2 {
				hasUser2 = true
			}
		}
		if hasUser1 && hasUser2 {
			return room, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeChatRepoSvc) IsMember(roomID, userID string) bool {
	for _, m := range r.members[roomID] {
		if m.UserID == userID {
			return true
		}
	}
	return false
}

var _ repository.ChatRepository = (*fakeChatRepoSvc)(nil)

// --- fake message repository ---

type fakeMessageRepoSvc struct {
	messages map[string]*model.Message // keyed by ID
}

func newFakeMessageRepoSvc() *fakeMessageRepoSvc {
	return &fakeMessageRepoSvc{messages: make(map[string]*model.Message)}
}

func (r *fakeMessageRepoSvc) Create(message *model.Message) error {
	if message.ID == "" {
		message.ID = "msg-" + message.ChatRoomID + "-" + message.SenderID
	}
	r.messages[message.ID] = message
	return nil
}

func (r *fakeMessageRepoSvc) FindByRoomID(roomID string, limit, offset int) ([]model.Message, error) {
	var result []model.Message
	for _, m := range r.messages {
		if m.ChatRoomID == roomID {
			result = append(result, *m)
		}
	}
	return result, nil
}

func (r *fakeMessageRepoSvc) FindByID(id string) (*model.Message, error) {
	m, ok := r.messages[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return m, nil
}

func (r *fakeMessageRepoSvc) MarkAsRead(roomID, userID string) error { return nil }

func (r *fakeMessageRepoSvc) GetUnreadCount(roomID, userID string) (int64, error) { return 0, nil }

func (r *fakeMessageRepoSvc) GetLastMessage(roomID string) (*model.Message, error) {
	return nil, gorm.ErrRecordNotFound
}

var _ repository.MessageRepository = (*fakeMessageRepoSvc)(nil)

// --- helpers ---

func newChatSvc(chatRepo *fakeChatRepoSvc, msgRepo *fakeMessageRepoSvc) ChatService {
	return NewChatService(chatRepo, msgRepo)
}

// --- CreateDirectChat tests ---

func TestCreateDirectChat_Success(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	room, err := svc.CreateDirectChat("user-1", "user-2")

	require.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, model.RoomTypeDirect, room.Type)
	// both members should be added
	members := chatRepo.members[room.ID]
	assert.Len(t, members, 2)
}

func TestCreateDirectChat_ReturnsExistingRoom(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	// Create first time
	room1, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	// Create again with same pair - should return existing room
	room2, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)
	assert.Equal(t, room1.ID, room2.ID)
}

// --- CreateGroupChat tests ---

func TestCreateGroupChat_Success(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	room, err := svc.CreateGroupChat("user-1", "My Group", []string{"user-2", "user-3"})

	require.NoError(t, err)
	assert.NotNil(t, room)
	assert.Equal(t, model.RoomTypeGroup, room.Type)
	assert.Equal(t, "My Group", room.Name)
	// creator + 2 members = 3
	members := chatRepo.members[room.ID]
	assert.Len(t, members, 3)
}

// --- GetUserChats tests ---

func TestGetUserChats_Success(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	// Create a couple of rooms for user-1
	_, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)
	_, err = svc.CreateGroupChat("user-1", "Group", []string{"user-2"})
	require.NoError(t, err)

	rooms, err := svc.GetUserChats("user-1")
	require.NoError(t, err)
	assert.Len(t, rooms, 2)
}

func TestGetUserChats_Empty(t *testing.T) {
	svc := newChatSvc(newFakeChatRepoSvc(), newFakeMessageRepoSvc())

	rooms, err := svc.GetUserChats("user-nobody")
	require.NoError(t, err)
	assert.Empty(t, rooms)
}

// --- SendMessage tests ---

func TestSendMessage_Success(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	msgRepo := newFakeMessageRepoSvc()
	svc := newChatSvc(chatRepo, msgRepo)

	// Set up room with member
	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	msg, err := svc.SendMessage("user-1", room.ID, "hello", model.MessageTypeText)

	require.NoError(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, "hello", msg.Content)
}

func TestSendMessage_NotMember(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	// Create room for user-1 and user-2, but user-3 is not a member
	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	_, err = svc.SendMessage("user-3", room.ID, "hello", model.MessageTypeText)

	assert.ErrorIs(t, err, ErrNotRoomMember)
}

// --- GetMessages tests ---

func TestGetMessages_Success(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	msgRepo := newFakeMessageRepoSvc()
	svc := newChatSvc(chatRepo, msgRepo)

	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	_, err = svc.SendMessage("user-1", room.ID, "msg1", model.MessageTypeText)
	require.NoError(t, err)
	_, err = svc.SendMessage("user-2", room.ID, "msg2", model.MessageTypeText)
	require.NoError(t, err)

	messages, err := svc.GetChatMessages("user-1", room.ID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, messages, 2)
}

func TestGetMessages_NotMember(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	_, err = svc.GetChatMessages("user-3", room.ID, 50, 0)
	assert.ErrorIs(t, err, ErrNotRoomMember)
}

func TestGetMessages_DefaultLimit(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	msgRepo := newFakeMessageRepoSvc()
	svc := newChatSvc(chatRepo, msgRepo)

	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	// limit=0 should default to 50 (no error)
	_, err = svc.GetChatMessages("user-1", room.ID, 0, 0)
	require.NoError(t, err)
}

// --- IsMember / verifyMembership tests (via SendMessage) ---

func TestIsMember_True(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	assert.True(t, chatRepo.IsMember(room.ID, "user-1"))
	assert.True(t, chatRepo.IsMember(room.ID, "user-2"))
}

func TestIsMember_False(t *testing.T) {
	chatRepo := newFakeChatRepoSvc()
	svc := newChatSvc(chatRepo, newFakeMessageRepoSvc())

	room, err := svc.CreateDirectChat("user-1", "user-2")
	require.NoError(t, err)

	assert.False(t, chatRepo.IsMember(room.ID, "user-99"))
}
