package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/godongmin/open_talk/backend/internal/model"
	"github.com/godongmin/open_talk/backend/internal/service"
	"github.com/stretchr/testify/assert"
)

// --- fake ChatService ---

type fakeChatService struct {
	createDirectRoom *model.ChatRoom
	createDirectErr  error

	createGroupRoom *model.ChatRoom
	createGroupErr  error

	getUserChats    []model.ChatRoom
	getUserChatsErr error

	getChatMessages    []model.Message
	getChatMessagesErr error

	sendMessageResult *model.Message
	sendMessageErr    error
}

func (f *fakeChatService) CreateDirectChat(userID, targetUserID string) (*model.ChatRoom, error) {
	return f.createDirectRoom, f.createDirectErr
}

func (f *fakeChatService) CreateGroupChat(userID string, name string, memberIDs []string) (*model.ChatRoom, error) {
	return f.createGroupRoom, f.createGroupErr
}

func (f *fakeChatService) GetUserChats(userID string) ([]model.ChatRoom, error) {
	return f.getUserChats, f.getUserChatsErr
}

func (f *fakeChatService) GetChatMessages(userID, roomID string, limit, offset int) ([]model.Message, error) {
	return f.getChatMessages, f.getChatMessagesErr
}

func (f *fakeChatService) SendMessage(userID, roomID, content string, messageType model.MessageType) (*model.Message, error) {
	return f.sendMessageResult, f.sendMessageErr
}

// --- fake MessageRepository ---

type fakeMessageRepo struct {
	unreadCount    int64
	unreadCountErr error

	lastMessage    *model.Message
	lastMessageErr error

	markAsReadErr error
}

func (r *fakeMessageRepo) Create(message *model.Message) error { return nil }

func (r *fakeMessageRepo) FindByRoomID(roomID string, limit, offset int) ([]model.Message, error) {
	return nil, nil
}

func (r *fakeMessageRepo) FindByID(id string) (*model.Message, error) { return nil, nil }

func (r *fakeMessageRepo) MarkAsRead(roomID, userID string) error {
	return r.markAsReadErr
}

func (r *fakeMessageRepo) GetUnreadCount(roomID, userID string) (int64, error) {
	return r.unreadCount, r.unreadCountErr
}

func (r *fakeMessageRepo) GetLastMessage(roomID string) (*model.Message, error) {
	return r.lastMessage, r.lastMessageErr
}

// --- router helpers ---

func newChatRouter(chatSvc service.ChatService, msgRepo *fakeMessageRepo) *gin.Engine {
	h := NewChatHandler(chatSvc, msgRepo)
	r := gin.New()
	r.POST("/api/chats", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.CreateChat(c)
	})
	r.GET("/api/chats", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.GetChats(c)
	})
	r.GET("/api/chats/:id/messages", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.GetMessages(c)
	})
	r.POST("/api/chats/:id/messages", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.SendMessage(c)
	})
	r.POST("/api/chats/:id/read", func(c *gin.Context) {
		c.Set("userID", "test-user-id")
		h.MarkMessagesRead(c)
	})
	return r
}

// --- CreateChat tests ---

func TestCreateChat_Direct_Success(t *testing.T) {
	chatSvc := &fakeChatService{
		createDirectRoom: &model.ChatRoom{ID: "room-1", Type: model.RoomTypeDirect},
	}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats", gin.H{
		"type":       "direct",
		"member_ids": []string{"other-user-id"},
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestCreateChat_Group_Success(t *testing.T) {
	chatSvc := &fakeChatService{
		createGroupRoom: &model.ChatRoom{ID: "room-2", Type: model.RoomTypeGroup, Name: "My Group"},
	}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats", gin.H{
		"type":       "group",
		"name":       "My Group",
		"member_ids": []string{"user-2", "user-3"},
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestCreateChat_BadRequest_MissingType(t *testing.T) {
	router := newChatRouter(&fakeChatService{}, &fakeMessageRepo{})

	// missing required "type" field
	w := postJSON(t, router, "/api/chats", gin.H{
		"member_ids": []string{"user-2"},
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateChat_BadRequest_InvalidType(t *testing.T) {
	router := newChatRouter(&fakeChatService{}, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats", gin.H{
		"type":       "invalid",
		"member_ids": []string{"user-2"},
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateChat_Direct_WrongMemberCount(t *testing.T) {
	router := newChatRouter(&fakeChatService{}, &fakeMessageRepo{})

	// direct chat requires exactly one member_id
	w := postJSON(t, router, "/api/chats", gin.H{
		"type":       "direct",
		"member_ids": []string{"user-2", "user-3"},
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- GetChats tests ---

func TestGetChats_Success(t *testing.T) {
	chatSvc := &fakeChatService{
		getUserChats: []model.ChatRoom{
			{ID: "room-1", Type: model.RoomTypeDirect},
			{ID: "room-2", Type: model.RoomTypeGroup, Name: "Group"},
		},
	}
	msgRepo := &fakeMessageRepo{unreadCount: 3, lastMessage: &model.Message{ID: "msg-1", Content: "hello"}}
	router := newChatRouter(chatSvc, msgRepo)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
	data := body["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestGetChats_Empty(t *testing.T) {
	chatSvc := &fakeChatService{getUserChats: []model.ChatRoom{}}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
	data := body["data"].([]interface{})
	assert.Len(t, data, 0)
}

// --- GetMessages tests ---

func TestGetMessages_Success(t *testing.T) {
	chatSvc := &fakeChatService{
		getChatMessages: []model.Message{
			{ID: "msg-1", Content: "hello"},
			{ID: "msg-2", Content: "world"},
		},
	}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats/room-1/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestGetMessages_NotMember(t *testing.T) {
	chatSvc := &fakeChatService{
		getChatMessagesErr: service.ErrNotRoomMember,
	}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats/room-1/messages", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

// --- SendMessage tests ---

func TestSendMessage_Success(t *testing.T) {
	chatSvc := &fakeChatService{
		sendMessageResult: &model.Message{ID: "msg-new", Content: "hello"},
	}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats/room-1/messages", gin.H{
		"content": "hello",
	})

	assert.Equal(t, http.StatusCreated, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestSendMessage_NotMember(t *testing.T) {
	chatSvc := &fakeChatService{
		sendMessageErr: service.ErrNotRoomMember,
	}
	router := newChatRouter(chatSvc, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats/room-1/messages", gin.H{
		"content": "hello",
	})

	assert.Equal(t, http.StatusForbidden, w.Code)
	body := decodeBody(t, w)
	assert.False(t, body["success"].(bool))
}

func TestSendMessage_BadRequest(t *testing.T) {
	router := newChatRouter(&fakeChatService{}, &fakeMessageRepo{})

	// missing required "content" field
	w := postJSON(t, router, "/api/chats/room-1/messages", gin.H{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// --- MarkMessagesRead tests ---

func TestMarkMessagesRead_Success(t *testing.T) {
	router := newChatRouter(&fakeChatService{}, &fakeMessageRepo{})

	w := postJSON(t, router, "/api/chats/room-1/read", gin.H{})

	assert.Equal(t, http.StatusOK, w.Code)
	body := decodeBody(t, w)
	assert.True(t, body["success"].(bool))
}

func TestMarkMessagesRead_Error(t *testing.T) {
	msgRepo := &fakeMessageRepo{markAsReadErr: service.ErrRoomNotFound}
	router := newChatRouter(&fakeChatService{}, msgRepo)

	w := postJSON(t, router, "/api/chats/room-1/read", gin.H{})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
